GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET = $(GOCMD) get
GOMOD = $(GOCMD) mod
GOVET = $(GOCMD) vet
GOFMT = $(GOCMD) fmt
GOMOD = $(GOCMD) mod

build: darwin web

all: darwin linux web

linux-full: linux web

darwin-full: darwin web

darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -a -o bin/btsrv.darwin backend/btsrv/*.go

web:
	mkdir -p bin/
	cd frontend && npm install && npm run build
	rm -rf bin/public
	mv frontend/dist bin/public
	rm -rf bin/templates
	mkdir -p bin/templates
	cp ./templates/* bin/templates

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -a -installsuffix cgo -o bin/btsrv.linux backend/btsrv/*.go

clean:
	$(GOCLEAN) ./backend/btsrv/...
	rm -rf bin

dep:
	cd frontend && npm upgrade
	$(GOGET) -u ./backend/btsrv/...
	$(GOMOD) tidy
	$(GOMOD) verify

fmt:
	$(GOFMT) ./backend/...

vet:
	$(GOVET) ./backend/...

check:
	go install honnef.co/go/tools/cmd/staticcheck
	$(HOME)/go/bin/staticcheck -checks all,-S1002,-ST1003,-S1008 ./backend/...
	go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
	$(GOVET) -vettool=$(HOME)/go/bin/shadow ./backend/...
