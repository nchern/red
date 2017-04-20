.PHONY: install install-deps build


install-deps:
	go get -u github.com/jteeuwen/go-bindata/...

bindata:
	go-bindata assets

install: bindata
	go get ./...

build: bindata
	go build ./...

test: bindata
	go test ./...
