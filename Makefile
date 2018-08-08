## golang makefile

SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=cs

BUILD_TIME=`date +%FT%T%z`


#go parameters
GOCMD=$(shell which go)
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install
GOTEST=$(GOCMD) test
GOFMT=gofmt -w

#LDFLAGS=-ldflags 

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	    $(GOBUILD)  -o ${BINARY} main.go

run :
	go run main.go
install:
	sudo mv ./cs /usr/local/bin 

.PHONY: clean
clean:
	    @if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
