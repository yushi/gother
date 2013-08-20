.PHONY: deps gother all
export GOPATH:=$(shell pwd)

all: gother

deps:
	go get -d -v gother/...

gother: deps
	go install gother/main/gother
