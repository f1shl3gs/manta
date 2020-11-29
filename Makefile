UISOURCES  := $(shell find ./ui -type f -not \( -path ./ui/dist/\* -o -path ./ui/src/.umi/\* -o -path ./ui/node_modules/\* -o -path ./ui/.cache/\* -o -name Makefile -prune \) )
GOSROUCES  := $(shell find . -type f -name '*.go') go.mod go.sum
DOCSOURCES := $(shell find ./docs -not \( -path ./docs/resources/\* -o -prune \) )

VERSION 	:= 0.9.0
GIT_SHA 	:= $(shell git rev-parse --short HEAD)
GIT_BRANCE 	:= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE 	:= $(shell date --rfc-3339 ns)
LDFLAGS    	:= "-s -w -X manta/version.Version=${VERSION} -X manta/version.GitSHA=${GIT_SHA} -X manta/version.GitBranch=${GIT_BRANCE}"

export GOOS=$(shell go env GOOS)
export GOBUILD=go build -ldflags ${LDFLAGS}

.PHONY: proto tidy deps

deps:
	go mod download

lint:
	go fmt ./...
	golangci-lint run ./...

proto:
	# protobuf 3.12.3 is required
	protoc \
		-I=./ \
        -I="${GOPATH}/src" 	\
		-I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
		--gogofaster_out=plugins=deepcopy+grpc+storeobject,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:./ \
		*.proto

mantad: $(GOSROUCES)
	CGO_ENABLED=0 $(GOBUILD) -o bin/$@ ./cmd/$(shell basename "$@")

manta: $(GOSROUCES)
	$(GOBUILD) -o bin/$@ ./cmd/$(shell basename "$@")
