UISOURCES  := $(shell find ./ui -type f -not \( -path ./ui/dist/\* -o -path ./ui/src/.umi/\* -o -path ./ui/node_modules/\* -o -path ./ui/.cache/\* -o -name Makefile -prune \) )
GOSROUCES  := $(shell find . -type f -name '*.go') go.mod go.sum
LDFLAGS    	:= "-extldflags '-static' -s -w -X manta/version.Version=${VERSION} -X manta/version.GitSHA=${GIT_SHA} -X manta/version.GitBranch=${GIT_BRANCE}"
GO			:= go

export GOOS=$(shell go env GOOS)
export GOBUILD=${GO} build -tags assets -ldflags ${LDFLAGS}

build: ui $(GOSROUCES) assets.tgz
	tar czf assets.tgz ui/build http/openapi.yaml
	CGO_ENABLED=0 $(GOBUILD) -o bin/mantad ./cmd/mantad

ui: $(UISOURCES)
	cd ui && yarn && yarn build

.PHONY: clean
clean:
	rm -rf bin
	rm -rf ui/build

.PHONY: dep
dep:
	go mod download

genproto: dep
	./scripts/genproto.sh

.PHONY: fmt
fmt: $(UISOURCES) $(GOSROUCES)
	go fmt
	cd ui && yarn prettier:fix

.PHONY: test
test: $(GOSROUCES) $(UISOURCES)
	go test ./...

.PHONY: cypress
cypress:
	cd ui && yarn cypress:run

.PHONY: lint
lint:
	golangci-lint run ./...
	cd ui && yarn prettier
