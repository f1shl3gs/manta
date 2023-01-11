package raftstore

import "github.com/f1shl3gs/manta/raftstore/raft"

type Config struct {
	raft.Config

	// Listen is the address, grpc server will listen to
	Listen       string
	DefragOnBoot bool
}
