.PHONY: build doc fmt lint run test vendor_clean vendor_get vendor_update vet

default: build

build: vet
	go build -v -o ./dist/mediabase ./server/boot.go

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
	rm -rf dist

run: clean build
	cp -r client/ dist
	ln -s "/Volumes/Users/kayak/Library/Application Support/net.apertoire.mediabase/web/img" dist/img
	ln -s "/Volumes/Users/kayak/Library/Application Support/net.apertoire.mediabase/db" dist/db
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