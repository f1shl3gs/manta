build:
	CGO_ENABLED=0 go build --ldflags "-s -w"

image: build
	podman build -t logproducer .