package raft

import "path/filepath"

var (
	maxRequestBytes uint = 4 * 1024 * 1024
)

type Config struct {
	Peers []string

	MaxRequestBytes uint
	DataDir         string

	ElectionTicks int
	// PreVote is true to enable Raft Pre-Vote
	PreVote bool
}

func (c *Config) WalDir() string {
	return filepath.Join(c.DataDir, "wals")
}

func (c *Config) SnapDir() string {
	return filepath.Join(c.DataDir, "snapshots")
}

func (c *Config) MemberDir() string {
	return filepath.Join(c.DataDir, "members")
}
