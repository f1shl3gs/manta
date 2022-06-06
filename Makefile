UISOURCES  := $(shell find ./ui -type f -not \( -path ./ui/dist/\* -o -path ./ui/src/.umi/\* -o -path ./ui/node_modules/\* -o -path ./ui/.cache/\* -o -name Makefile -prune \) )
GOSROUCES  := $(shell find . -type f -name '*.go') go.mod go.sum
LDFLAGS    	:= "-extldflags '-static' -s -w -X manta/version.Version=${VERSION} -X manta/version.GitSHA=${GIT_SHA} -X manta/version.GitBranch=${GIT_BRANCE}"
GO			:= go

export GOOS=$(shell go env GOOS)
export GOBUILD=${GO} build -tags assets -ldflags ${LDFLAGS}

dep:
	go mod download

genproto: dep
	./scripts/genproto.sh

.PHONY: ui
ui: $(UISOURCES)
	cd ui && yarn build
