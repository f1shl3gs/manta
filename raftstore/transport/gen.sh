#!/usr/bin/env bash

set -euo pipefail

protoc \
  -I="./" \
  -I="${GOPATH}/src" 	\
	-I="${GOPATH}/src/github.com/gogo/protobuf" \
	-I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
	-I="${GOPATH}/pkg/mod/go.etcd.io/etcd/raft/v3@v3.5.0-alpha.0" \
  --gogofaster_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:./ \
  ./*.proto

sed -i 's/\"raftpb\"/\"go.etcd.io\/etcd\/raft\/v3\/raftpb\"/g' raft.pb.go