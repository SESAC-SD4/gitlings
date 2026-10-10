.PHONY: build test clean

build:
	bash build.sh

test:
	go test ./... -timeout 120s

clean:
	rm -f gitlings gitlings.exe gitlings-darwin-arm64 gitlings-darwin-amd64
