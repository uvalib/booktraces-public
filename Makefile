GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test

build: darwin web

all: darwin linux web

linux-full: linux web

darwin-full: darwin web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/btsrv.darwin backend/btsrv/*.go

web:
	mkdir -p bin/
	cd frontend/; yarn install; yarn build
	rm -rf bin/public
	mv frontend/dist bin/public

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/btsrv.linux backend/btsrv/*.go

clean:
	$(GOCLEAN)
	rm -rf bin
