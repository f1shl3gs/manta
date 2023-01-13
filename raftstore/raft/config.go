package raft

import "path/filepath"

var (
	maxRequestBytes uint = 4 * 1024 * 1024
)

type Config struct {
	Listen string

	Peers []string

	DataDir string
}

func (c *Config) MemberDir() string {
	return filepath.Join(c.DataDir, "members")
}
