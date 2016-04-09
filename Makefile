GOROOT := /usr/local/go
GOPATH := $(shell pwd)
GOBIN  := $(GOPATH)/bin
PATH   := $(GOROOT)/bin:$(PATH)
DEPS   := github.com/ogier/pflag github.com/jmcvetta/napping

COMMIT = $(git log | head -n 1 | cut  -f 2 -d ' ')

all: glblsmsr

deps: $(DEPS)
	GOPATH=$(GOPATH) go get -u $^

glblsmsr: util.go main.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOPATH=$(GOPATH) go build -o $@ -v $^
		touch $@

linux64: main.go commands.go stack.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOOS=linux GOARCH=amd64 GOPATH=$(GOPATH) go build -o glblsmsr-linux-amd64.bin -v $^
		touch glblsmsr-linux-amd64.bin

win64: main.go commands.go stack.go
    # always format code
		GOPATH=$(GOPATH) go fmt $^
		# vet it
		GOPATH=$(GOPATH) go tool vet $^
    # binary
		GOOS=windows GOARCH=amd64 GOPATH=$(GOPATH) go build -o glblsmsr-win-amd64.exe -v $^
		touch glblsmsr-win-amd64.exe

.PHONY: $(DEPS) clean

clean:
	rm -f glblsmsr glblsmsr-linux-amd64.bin glblsmsr-win-amd64.exe
