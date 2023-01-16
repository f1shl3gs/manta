package raft

import "path/filepath"

type Config struct {
	Listen string

	Peers []string

	DataDir string
}

func (c *Config) MemberDir() string {
	return filepath.Join(c.DataDir, "members")
}
