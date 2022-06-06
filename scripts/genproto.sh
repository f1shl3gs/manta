#!/usr/bin/env bash

set -euo pipefail

# Requirements
if [[ $(protoc --version | cut -f2 -d' ') != "3.14.0" ]]; then
  echo "could not find protoc 3.14.0, is it installed + in PATH?"
  exit 255
fi

# Install tools
GOBIN="${PWD}/bin"
export GOBIN=${GOBIN}
if [ ! -d ${GOBIN} ]; then
  mkdir "${GOBIN}"
fi

# go install github.com/gogo/protobuf/protoc-gen-gogofaster

# gogo protobuf
GOGOPROTO_ROOT="${GOPATH}/pkg/mod/$(go list -m github.com/gogo/protobuf | sed 's/\ /@/g')"

export PATH=${PATH}:${GOBIN}
echo
echo "Resolved binary and packages versions:"
echo "  - protoc                    $(which protoc)"
echo "  - protoc-gen-gogofaster:    $(which protoc-gen-gogofaster)"

# Generate
protoc -I. \
  -I "${GOGOPROTO_ROOT}/protobuf" \
  -I "${GOGOPROTO_ROOT}" \
  --gogofaster_out=Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:./ \
  *.proto