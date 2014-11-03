mb_osname := $(shell uname)
mb_version := $(shell cat VERSION)
mb_count := $(shell git rev-list --count $(shell cat VERSION)..)
mb_hash := $(shell git rev-parse --short HEAD)

.PHONY: build doc fmt lint run test vendor_clean vendor_get vendor_update vet

default: build

build: vet
	go build -ldflags "-X main.Version $(mb_version)-$(mb_count).$(mb_hash)" -v -o ./dist/mediabase ./server/boot.go

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
	rm -rf build
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

release: clean build
	mkdir build
	cp dist/mediabase build
	cp CHANGES build
	cp LICENSE build
	cp README.md build
ifeq ($(mb_osname), Darwin)
	cd build && zip mediabase-$(mb_version)-darwin-amd64.zip *
else
	cd build && tar czf mediabase-$(mb_version)-linux-amd64.tar.gz *
endif


# http://godoc.org/code.google.com/p/go.tools/cmd/vet
# go get code.google.com/p/go.tools/cmd/vet
vet:
	go vet ./server/...
