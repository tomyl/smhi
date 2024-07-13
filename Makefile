
.PHONY: all install lint

build:
	go build ./cmd/smhi

install:
	go install ./cmd/smhi

all: lint install

lint:
	go mod verify
	go vet ./...
	go run honnef.co/go/tools/cmd/staticcheck@latest ./...
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
