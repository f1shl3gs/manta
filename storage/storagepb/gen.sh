#!/usr/bin/env bash

set -euo pipefail

protoc \
  -I=. \
  -I="${GOPATH}/src/github.com/gogo/protobuf" \
  -I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
  --gogofaster_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:./ \
		*.proto

cd prompb

protoc  \
  -I=. \
  -I=../ \
  -I="${GOPATH}/src/github.com/gogo/protobuf" \
  -I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
  --gogofaster_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:./ \
		*.proto