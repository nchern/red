.PHONY: install install-deps build


install-deps:
	go get -u github.com/jteeuwen/go-bindata/...

bindata:
	go-bindata -pkg app -o app/bindata.go assets

install: bindata
	go install ./...

build: bindata
	go build ./...

test: bindata
	go test ./...
