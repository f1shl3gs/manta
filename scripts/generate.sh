#!/bin/bash

set -euo pipefail

# gogo
go get github.com/gogo/protobuf@f67b8970b736e53dbd7d0a27146c8f1ac52f74e5
go get go.etcd.io/etcd/raft/v3@cecbe35ce0703cd0f8d2063dad4a9e541ae317e5

GOGOPROTO_ROOT="${GOPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.3-0.20221024144010-f67b8970b736"
GOGOPROTO_PATH="${GOGOPROTO_ROOT}/protobuf"
RAFTPROTO_PATH="${GOPATH}/pkg/mod/go.etcd.io/etcd/raft/v3@v3.5.6"

function gen_proto {
    local DIR="$1"

    protoc \
        -I="${DIR}" \
        -I="${GOGOPROTO_ROOT}" \
        -I="${GOGOPROTO_PATH}" \
        -I="${RAFTPROTO_PATH}" \
        --gogofaster_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:${DIR} ${DIR}/*.proto
}

# generate proto
gen_proto ./raftstore/kvpb
gen_proto ./raftstore/mvcc
