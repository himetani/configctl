.PHONY: all build install clean test

all: ;

NAME := configctl
VERSION  := 0.0.1
REVISION  := $(shell git rev-parse --short HEAD)
LDFLAGS := -ldflags="-s -w -X \"github.com/himetani/configctl/cmd.Version=$(VERSION)\" -X \"github.com/himetani/configctl/cmd.Revision=$(REVISION)\""

SRCS    := $(shell find . -path ./vendor -prune -o -name '*.go' -print)

bin/$(NAME): $(SRCS)
	go build $(LDFLAGS) -o bin/$(NAME)


$$GOPATH/bin/$(NAME):
	go install $(LDFLAGS)

build: bin/$(NAME)

install: $$GOPATH/bin/$(NAME)

clean:
	rm -rf bin/*

test: 
	go test -v github.com/himetani/configctl/...
