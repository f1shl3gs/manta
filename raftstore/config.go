package raftstore

type Config struct {
	Peers   []string
	DataDir string

	// Listen is the address, grpc server will listen to
	Listen       string
	DefragOnBoot bool
}
