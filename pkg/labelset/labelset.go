package labelset

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// SeparatorByte is a byte that cannot occur in valid UTF-8 sequences and is
// used to separate label names, label values, and other strings from each other
// when calculating their combined hash value (aka signature aka fingerprint).
const SeparatorByte byte = 255

var (
	// cache the signature of an empty label set.
	emptyLabelSignature = hashNew()
)

// A LabelSet is a collection of LabelName and LabelValue pairs.  The LabelSet
// may be fully-qualified down to the point where it may resolve to a single
// Metric in the data store or not.  All operations that occur within the realm
// of a LabelSet can emit a vector of Metric entities to which the LabelSet may
// match.
type LabelSet map[LabelName]LabelValue

// Validate checks whether all names and values in the label set
// are valid.
func (ls LabelSet) Validate() error {
	for ln, lv := range ls {
		if !ln.IsValid() {
			return fmt.Errorf("invalid name %q", ln)
		}
		if !lv.IsValid() {
			return fmt.Errorf("invalid value %q", lv)
		}
	}
	return nil
}

func (ls LabelSet) Set(k, v string) {
	ls[LabelName(k)] = LabelValue(v)
}

// Equal returns true iff both label sets have exactly the same key/value pairs.
func (ls LabelSet) Equal(o LabelSet) bool {
	if len(ls) != len(o) {
		return false
	}
	for ln, lv := range ls {
		olv, ok := o[ln]
		if !ok {
			return false
		}
		if olv != lv {
			return false
		}
	}
	return true
}

// Before compares the metrics, using the following criteria:
//
// If m has fewer labels than o, it is before o. If it has more, it is not.
//
// If the number of labels is the same, the superset of all label names is
// sorted alphanumerically. The first differing label pair found in that order
// determines the outcome: If the label does not exist at all in m, then m is
// before o, and vice versa. Otherwise the label value is compared
// alphanumerically.
//
// If m and o are equal, the method returns false.
func (ls LabelSet) Before(o LabelSet) bool {
	if len(ls) < len(o) {
		return true
	}
	if len(ls) > len(o) {
		return false
	}

	lns := make(LabelNames, 0, len(ls)+len(o))
	for ln := range ls {
		lns = append(lns, ln)
	}
	for ln := range o {
		lns = append(lns, ln)
	}
	// It's probably not worth it to de-dup lns.
	sort.Sort(lns)
	for _, ln := range lns {
		mlv, ok := ls[ln]
		if !ok {
			return true
		}
		olv, ok := o[ln]
		if !ok {
			return false
		}
		if mlv < olv {
			return true
		}
		if mlv > olv {
			return false
		}
	}
	return false
}

// Clone returns a copy of the label set.
func (ls LabelSet) Clone() LabelSet {
	lsn := make(LabelSet, len(ls))
	for ln, lv := range ls {
		lsn[ln] = lv
	}
	return lsn
}

// Merge is a helper function to non-destructively merge two label sets.
func (l LabelSet) Merge(other LabelSet) LabelSet {
	result := make(LabelSet, len(l))

	for k, v := range l {
		result[k] = v
	}

	for k, v := range other {
		result[k] = v
	}

	return result
}

func (l LabelSet) String() string {
	lstrs := make([]string, 0, len(l))
	for l, v := range l {
		lstrs = append(lstrs, fmt.Sprintf("%s=%q", l, v))
	}

	sort.Strings(lstrs)
	return fmt.Sprintf("{%s}", strings.Join(lstrs, ", "))
}

// Fingerprint returns the LabelSet's fingerprint.
func (ls LabelSet) Fingerprint() Fingerprint {
	return labelSetToFingerprint(ls)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *LabelSet) UnmarshalJSON(b []byte) error {
	var m map[LabelName]LabelValue
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	// encoding/json only unmarshals maps of the form map[string]T. It treats
	// LabelName as a string and does not call its UnmarshalJSON method.
	// Thus, we have to replicate the behavior here.
	for ln := range m {
		if !ln.IsValid() {
			return fmt.Errorf("%q is not a valid label name", ln)
		}
	}
	*l = LabelSet(m)
	return nil
}

// labelSetToFingerprint works exactly as LabelsToSignature but takes a LabelSet as
// parameter (rather than a label map) and returns a Fingerprint.
func labelSetToFingerprint(ls LabelSet) Fingerprint {
	if len(ls) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	labelNames := make(LabelNames, 0, len(ls))
	for labelName := range ls {
		labelNames = append(labelNames, labelName)
	}
	sort.Sort(labelNames)

	sum := hashNew()
	for _, labelName := range labelNames {
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(ls[labelName]))
		sum = hashAddByte(sum, SeparatorByte)
	}
	return Fingerprint(sum)
}

// labelSetToFastFingerprint works similar to labelSetToFingerprint but uses a
// faster and less allocation-heavy hash function, which is more susceptible to
// create hash collisions. Therefore, collision detection should be applied.
func labelSetToFastFingerprint(ls LabelSet) Fingerprint {
	if len(ls) == 0 {
		return Fingerprint(emptyLabelSignature)
	}

	var result uint64
	for labelName, labelValue := range ls {
		sum := hashNew()
		sum = hashAdd(sum, string(labelName))
		sum = hashAddByte(sum, SeparatorByte)
		sum = hashAdd(sum, string(labelValue))
		result ^= sum
	}
	return Fingerprint(result)
}
