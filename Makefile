PROJECT = rdap-exporter
VERSION = $(shell cat VERSION)
LDFLAGS=-ldflags "-w -s -X github.com/hsn723/rdap-exporter/cmd.version=${VERSION}"

CST_VERSION = 1.16.0

WORKDIR = /tmp/$(PROJECT)/work
BINDIR = /tmp/$(PROJECT)/bin
CONTAINER_STRUCTURE_TEST = $(BINDIR)/container-structure-test

PATH := $(PATH):$(BINDIR)

export PATH

all: build

.PHONY: clean
clean:
	@if [ -f $(PROJECT) ]; then rm $(PROJECT); fi

.PHONY: lint
lint:
	@if [ -z "$(shell which pre-commit)" ]; then pip3 install pre-commit; fi
	pre-commit install
	pre-commit run --all-files

.PHONY: test
test:
	go test --tags=test -coverprofile cover.out -count=1 -race -p 4 -v ./...

.PHONY: $(CONTAINER_STRUCTURE_TEST)
$(CONTAINER_STRUCTURE_TEST): $(BINDIR)
	curl -sSLf -o $(CONTAINER_STRUCTURE_TEST) -O https://storage.googleapis.com/container-structure-test/v$(CST_VERSION)/container-structure-test-linux-amd64 && chmod +x $(CONTAINER_STRUCTURE_TEST)

.PHONY: container-structure-test
container-structure-test: $(CONTAINER_STRUCTURE_TEST)
	printf "amd64\narm64" | xargs -n1 -I {} $(CONTAINER_STRUCTURE_TEST) test --image ghcr.io/hsn723/$(PROJECT):$(shell git describe --tags --abbrev=0 --match "v*" || echo v0.0.0)-next-{} --config cst.yaml

.PHONY: verify
verify:
	go mod download
	go mod verify

.PHONY: build
build: clean
	env CGO_ENABLED=0 go build $(LDFLAGS) .

$(BINDIR):
	mkdir -p $(BINDIR)

$(WORKDIR):
	mkdir -p $(WORKDIR)
