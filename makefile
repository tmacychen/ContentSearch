## golang makefile

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=cs

VERSION=0.0.1
BUILD_TIME=`date +%FT%T%z`


#go parameters
GOCMD=$(shell which go)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOFMT=gofmt -w

#LDFLAGS=-ldflags 
#"-X github.com/ariejan/roll/core.Version=${VERSION} -X github.com/ariejan/roll/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	    $(GOBUILD)  -o ${BINARY} main.go

run :
	go run main.go
install:
#	    go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	    @if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
