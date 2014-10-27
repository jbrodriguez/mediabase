.PHONY: build doc fmt lint run test vendor_clean vendor_get vendor_update vet

default: build

build: vet
#	CMT_VERSION=$(shell cat VERSION)
#	CMT_COUNT=$(shell git rev-list --count $(shell cat VERSION)..)
#	CMT_HASH=$(shell git rev-parse --short HEAD)
	go build -ldflags "-X main.Version $(shell cat VERSION)-$(shell git rev-list --count $(shell cat VERSION)..).$(shell git rev-parse --short HEAD)" -v -o ./dist/mediabase ./server/boot.go

doc:
	godoc -http=:6060 -index

# http://golang.org/cmd/go/#hdr-Run_gofmt_on_package_sources
fmt:
	go fmt ./...

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	golint ./src

test:
	go test ./...

clean:
	rm -rf dist/web/app
	rm -rf dist/web/css
	rm -rf dist/web/js
	rm -rf dist/web/index.html
	unlink dist/web/img || true
	unlink dist/db || true

run: clean build
	cp -r client/* dist/web
	ln -s "$(shell echo $$HOME)/.mediabase/db" dist/db
	ln -s "$(shell echo $$HOME)/.mediabase/web/img" dist/web/img
	cd dist && ./mediabase

vendor_clean:
	rm -dRf ./_vendor/src

# We have to set GOPATH to just the _vendor
# directory to ensure that `go get` doesn't
# update packages in our primary GOPATH instead.
# This will happen if you already have the package
# installed in GOPATH since `go get` will use
# that existing location as the destination.
vendor_get: vendor_clean
	GOPATH=${PWD}/_vendor go get -d -u -v \
	github.com/jpoehls/gophermail \
	github.com/codegangsta/martini

vendor_update: vendor_get
	rm -rf `find ./_vendor/src -type d -name .git` \
	&& rm -rf `find ./_vendor/src -type d -name .hg` \
	&& rm -rf `find ./_vendor/src -type d -name .bzr` \
	&& rm -rf `find ./_vendor/src -type d -name .svn`

# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./server/...
