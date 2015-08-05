GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/ogier/pflag github.com/jmcvetta/napping

COMMIT = $(git log | head -n 1 | cut  -f 2 -d ' ')

all: glblsmsr

update: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

glblsmsr: util.go main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
    # binary
		GOPATH=$(GOPATH) go build -ldflags "-X main.commit $(COMMIT)" -o $@ -v $^
		touch $@

windows:
	  gox -os="windows"

.PHONY: $(DEPS) clean

clean:
	rm -f glblsmsr
