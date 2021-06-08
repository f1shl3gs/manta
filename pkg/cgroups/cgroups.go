// Copyright 2020 The Cockroach Authors.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0, included in the file
// licenses/APL.txt.

package cgroups

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	cgroupV1MemLimitFilename     = "memory.stat"
	cgroupV2MemLimitFilename     = "memory.max"
	cgroupV1CPUQuotaFilename     = "cpu.cfs_quota_us"
	cgroupV1CPUPeriodFilename    = "cpu.cfs_period_us"
	cgroupV1CPUSysUsageFilename  = "cpuacct.usage_sys"
	cgroupV1CPUUserUsageFilename = "cpuacct.usage_user"
	cgroupV2CPUMaxFilename       = "cpu.max"
	cgroupV2CPUStatFilename      = "cpu.stat"
)

type CPUUsage struct {
	// System time and user time taken by this cgroup or process.
	// In nanoseconds
	SysTime  uint64
	UserTime uint64

	// CPU period and quota for this process, in microseconds.
	// This cgroup has access to up to (quota/period) proportion
	// of CPU resources on the system. For instance, if there
	// are 4 CPUs, quota = 150000, period = 100000, this cgroup
	// can use around ~1.5 CPUs, or 37.5% of total scheduler
	// time. If quota is -1, it's unlimited
	Period int64
	Quota  int64

	// NumCPU is the number of CPUs in the system. Always returned
	// even if not called from a cgroup
	NumCPU int
}

// CPUShares returns the number of CPUs this cgroup can be exepected
// to max out. If there's no limit, NumCPU is returned
func (c CPUUsage) CPUShares() float64 {
	if c.Period <= 0 || c.Quota <= 0 {
		return float64(c.NumCPU)
	}

	return float64(c.Quota) / float64(c.Period)
}

// GetCgroupCPU returns the cpu usage and quota for the
// current cgroup
func GetCgroupCPU() (CPUUsage, error) {
	cu, err := getCgroupCPU("/")
	if err != nil {
		return CPUUsage{}, err
	}

	cu.NumCPU = runtime.NumCPU()
	return cu, nil
}

func getCgroupCPU(root string) (CPUUsage, error) {
	path, err := detectControllerPath(filepath.Join(root, "/proc/self/cgroup"), "cpu,cpuacct")
	if err != nil {
		return CPUUsage{}, err
	}

	// No CPU controller detected
	if path == "" {
		return CPUUsage{}, errors.New("no cpu controller detected")
	}

	mount, ver, err := getCgroupDetails(filepath.Join(root, "/proc/self/mountinfo"), path, "cpu,cpuacct")
	if err != nil {
		return CPUUsage{}, err
	}

	var res CPUUsage

	switch ver {
	case 1:
		res.Period, res.Quota, err = detectCPUQuotaInV1(filepath.Join(root, mount))
		if err != nil {
			return CPUUsage{}, err
		}

		res.SysTime, res.UserTime, err = detectCPUUsageInV1(filepath.Join(root, mount))
		if err != nil {
			return CPUUsage{}, err
		}
	case 2:
		res.Period, res.Quota, err = detectCPUQuotaInV2(filepath.Join(root, mount, path))
		if err != nil {
			return CPUUsage{}, err
		}

		res.SysTime, res.UserTime, err = detectCPUUsageInV2(filepath.Join(root, mount, path))
		if err != nil {
			return CPUUsage{}, err
		}

	default:
		return CPUUsage{}, fmt.Errorf("detected unknown cgroup version index: %d", ver)
	}

	return res, nil
}

// getCgroupDetails reads /proc/[pid]/mountinfo for cgroup or cgroup2 mount which defines the used
// version. See http://man7.org/linux/man-pages/man5/proc.5.html for `mountinfo` format.
func getCgroupDetails(mountinfoPath string, root string, controller string) (string, int, error) {
	info, err := os.Open(mountinfoPath)
	if err != nil {
		return "", 0, errors.Wrapf(err, "failed to read mounts info from file: %s", mountinfoPath)
	}
	defer func() {
		_ = info.Close()
	}()

	scanner := bufio.NewScanner(info)
	for scanner.Scan() {
		fields := bytes.Fields(scanner.Bytes())
		if len(fields) < 10 {
			continue
		}

		ver, ok := detectCgroupVersion(fields, controller)
		if ok {
			mountPoint := string(fields[4])
			if ver == 2 {
				return mountPoint, ver, nil
			}
			// It is possible that the controller mount and the cgroup path are not the same (both are relative to the NS root).
			// So start with the mount and construct the relative path of the cgroup.
			// To test:
			//   cgcreate -t $USER:$USER -a $USER:$USER -g memory:crdb_test
			//   echo 3999997952 > /sys/fs/cgroup/memory/crdb_test/memory.limit_in_bytes
			//   cgexec -g memory:crdb_test ./cockroach start-single-node
			// cockroach.log -> server/config.go:433 ⋮ system total memory: ‹3.7 GiB›
			//   without constructing the relative path
			// cockroach.log -> server/config.go:433 ⋮ system total memory: ‹63 GiB›
			nsRelativePath := string(fields[3])
			if !strings.Contains(nsRelativePath, "..") {
				// We don't expect to see err here ever but in case that it happens
				// the best action is to ignore the line and hope that the rest of the lines
				// will allow us to extract a valid path.
				if relPath, err := filepath.Rel(nsRelativePath, root); err == nil {
					return filepath.Join(mountPoint, relPath), ver, nil
				}
			}
		}
	}

	return "", 0, fmt.Errorf("failed to detect cgroup root mount and version")
}

// Return version of cgroup mount for memory controller if found
func detectCgroupVersion(fields [][]byte, controller string) (_ int, found bool) {
	if len(fields) < 10 {
		return 0, false
	}

	// Due to strange format there can be optional fields in the middle of the set, starting
	// from the field #7. The end of the fields is marked with "-" field
	var pos = 6
	for pos < len(fields) {
		if bytes.Equal(fields[pos], []byte{'-'}) {
			break
		}

		pos++
	}

	// No optional fields separator found or there is less than 3 fields after it which is wrong
	if (len(fields) - pos - 1) < 3 {
		return 0, false
	}

	pos++

	// Check for controller specifically in cgroup v1 (it is listed in super options field),
	// as the limit can't be found if it is not enforced
	if bytes.Equal(fields[pos], []byte("cgroup")) && bytes.Contains(fields[pos+2], []byte(controller)) {
		return 1, true
	} else if bytes.Equal(fields[pos], []byte("cgroup2")) {
		return 2, true
	}

	return 0, false
}

// detectControllerPath
// the controller is defined via either type `memory` for cgroup v1 or via empty type
// for cgroup v2, where the type is the second failed in /proc/[pid]/cgroup file
func detectControllerPath(path string, controller string) (string, error) {
	cgroup, err := os.Open(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read %q cgroup from cgroups file %s", controller, path)
	}

	defer cgroup.Close()

	scanner := bufio.NewScanner(cgroup)
	var unifiedPathIfFound string
	for scanner.Scan() {
		fields := bytes.Split(scanner.Bytes(), []byte{':'})
		if len(fields) != 3 {
			// The lines should always have three fields, there's something fishy here.
			continue
		}

		f0, f1 := string(fields[0]), string(fields[1])
		// First case if v2, second - v1. We give v2 the priority here.
		// There is also a `hybrid` mode when both  versions are enabled,
		// but no known container solutions support it afaik
		if f0 == "0" && f1 == "" {
			unifiedPathIfFound = string(fields[2])
		} else if f1 == controller {
			return string(fields[2]), nil
		}
	}

	return unifiedPathIfFound, nil
}

func detectCPUQuotaInV1(cRoot string) (period, quota int64, err error) {
	quotaFilePath := filepath.Join(cRoot, cgroupV1CPUQuotaFilename)
	periodFilePath := filepath.Join(cRoot, cgroupV1CPUPeriodFilename)
	quota, err = cgroupFileToInt64(quotaFilePath, "cpu quota")
	if err != nil {
		return 0, 0, err
	}
	period, err = cgroupFileToInt64(periodFilePath, "cpu period")
	if err != nil {
		return 0, 0, err
	}

	return period, quota, err
}

func detectCPUUsageInV1(cRoot string) (stime, utime uint64, err error) {
	sysFilePath := filepath.Join(cRoot, cgroupV1CPUSysUsageFilename)
	userFilePath := filepath.Join(cRoot, cgroupV1CPUUserUsageFilename)
	stime, err = cgroupFileToUint64(sysFilePath, "cpu system time")
	if err != nil {
		return 0, 0, err
	}
	utime, err = cgroupFileToUint64(userFilePath, "cpu user time")
	if err != nil {
		return 0, 0, err
	}

	return stime, utime, err
}

func detectCPUQuotaInV2(cRoot string) (period, quota int64, err error) {
	maxFilePath := filepath.Join(cRoot, cgroupV2CPUMaxFilename)
	contents, err := readFile(maxFilePath)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "error when read cpu quota from cgroup v2 at %s", maxFilePath)
	}
	fields := strings.Fields(string(contents))
	if len(fields) > 2 || len(fields) == 0 {
		return 0, 0, errors.Errorf("unexpected format when reading cpu quota from cgroup v2 at %s: %s", maxFilePath, contents)
	}
	if fields[0] == "max" {
		// Negative quota denotes no limit.
		quota = -1
	} else {
		quota, err = strconv.ParseInt(fields[0], 10, 64)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "error when reading cpu quota from cgroup v2 at %s", maxFilePath)
		}
	}
	if len(fields) == 2 {
		period, err = strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "error when reading cpu period from cgroup v2 at %s", maxFilePath)
		}
	}
	return period, quota, nil
}

func detectCPUUsageInV2(cRoot string) (stime, utime uint64, err error) {
	statFilePath := filepath.Join(cRoot, cgroupV2CPUStatFilename)
	var stat *os.File
	stat, err = os.Open(statFilePath)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "can't read cpu usage from cgroup v2 at %s", statFilePath)
	}
	defer func() {
		if err == nil {
			err = stat.Close()
		}
	}()

	scanner := bufio.NewScanner(stat)
	for scanner.Scan() {
		fields := bytes.Fields(scanner.Bytes())
		if len(fields) != 2 || (string(fields[0]) != "user_usec" && string(fields[0]) != "system_usec") {
			continue
		}
		keyField := string(fields[0])

		trimmed := string(bytes.TrimSpace(fields[1]))
		usageVar := &stime
		if keyField == "user_usec" {
			usageVar = &utime
		}
		*usageVar, err = strconv.ParseUint(trimmed, 10, 64)
		if err != nil {
			return 0, 0, errors.Wrapf(err, "can't read cpu usage %s from cgroup v1 at %s", keyField, statFilePath)
		}
	}

	return stime, utime, err
}

func cgroupFileToUint64(filepath, desc string) (res uint64, err error) {
	contents, err := readFile(filepath)
	if err != nil {
		return 0, errors.Wrapf(err, "error when reading %s from cgroup v1 at %s", desc, filepath)
	}
	res, err = strconv.ParseUint(string(bytes.TrimSpace(contents)), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "error when parsing %s from cgroup v1 at %s", desc, filepath)
	}
	return res, err
}

func cgroupFileToInt64(filepath, desc string) (res int64, err error) {
	contents, err := readFile(filepath)
	if err != nil {
		return 0, errors.Wrapf(err, "error when reading %s from cgroup v1 at %s", desc, filepath)
	}
	res, err = strconv.ParseInt(string(bytes.TrimSpace(contents)), 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "error when parsing %s from cgroup v1 at %s", desc, filepath)
	}
	return res, nil
}

func readFile(filepath string) (res []byte, err error) {
	var f *os.File
	f, err = os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			err = f.Close()
		}
	}()
	res, err = ioutil.ReadAll(f)
	return res, err
}