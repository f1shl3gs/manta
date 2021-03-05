UISOURCES  := $(shell find ./ui -type f -not \( -path ./ui/dist/\* -o -path ./ui/src/.umi/\* -o -path ./ui/node_modules/\* -o -path ./ui/.cache/\* -o -name Makefile -prune \) )
GOSROUCES  := $(shell find . -type f -name '*.go') go.mod go.sum
DOCSOURCES := $(shell find ./docs -not \( -path ./docs/resources/\* -o -prune \) )

VERSION 	:= 0.9.0
GIT_SHA 	:= $(shell git rev-parse --short HEAD)
GIT_BRANCE 	:= $(shell git rev-parse --abbrev-ref HEAD)
BUILD_DATE 	:= $(shell date --rfc-3339 ns)
LDFLAGS    	:= "-extldflags '-static' -s -w -X manta/version.Version=${VERSION} -X manta/version.GitSHA=${GIT_SHA} -X manta/version.GitBranch=${GIT_BRANCE}"
GO			:= go

export GOOS=$(shell go env GOOS)
export GOBUILD=${GO} build -ldflags ${LDFLAGS}

.PHONY: proto tidy deps

deps:
	go mod download

lint:
	go fmt ./...
	golangci-lint run ./...

test:
	go test ./...

proto:
	# protobuf 3.12.3 is required
	protoc \
		--experimental_allow_proto3_optional \
		-I=./ \
        -I="${GOPATH}/src" 	\
		-I="${GOPATH}/src/github.com/gogo/protobuf/protobuf" \
		--gogofaster_out=plugins=grpc,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/struct.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/wrappers.proto=github.com/gogo/protobuf/types:./ \
		*.proto

swagger:
	wget https://codeload.github.com/swagger-api/swagger-ui/tar.gz/v3.44.1 -O swagger.tgz

mantad: $(GOSROUCES) $(UISOURCES) deps
	CGO_ENABLED=0 $(GOBUILD) -o bin/$@ ./cmd/$(shell basename "$@")


