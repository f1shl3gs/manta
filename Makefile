UISOURCES   := $(shell find ./ui -type f -not \( -path ./ui/dist/\* -o -path ./ui/node_modules/\* -o -name Makefile -prune \) )
GOSROUCES   := $(shell find . -type f -name '*.go') go.mod go.sum
LDFLAGS    	:= "-extldflags '-static' -s -w -X manta/version.Version=${VERSION} -X manta/version.GitSHA=${GIT_SHA} -X manta/version.GitBranch=${GIT_BRANCE}"
GO			:= go

export GOOS=$(shell go env GOOS)
export GOBUILD=${GO} build -tags assets -ldflags ${LDFLAGS}

.PHONY: build
build: ui $(GOSROUCES)
	tar czf assets.tgz ui/build
	CGO_ENABLED=0 $(GOBUILD) -o bin/mantad ./cmd/mantad

genproto:
	@bash ./scripts/genproto.sh

.PHONY: ui
ui: $(UISOURCES)
	cd ui && yarn && yarn build

.PHONY: clean
clean:
	rm -rf bin
	rm -rf ui/build
	rm -rf ui/node_modules
	rm -rf ui/cypress/downloads
	rm -rf ui/cypress/screenshots
	rm -rf ui/cypress/videos

.PHONY: dep
dep:
	go mod download

.PHONY: fmt
fmt: $(UISOURCES) $(GOSROUCES)
	go fmt ./...
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

.PHONY: lines
lines:
	@./scripts/lines.sh