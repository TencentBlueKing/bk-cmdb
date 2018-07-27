.PHONY: all install test testrace vet int

all: test vet

install:
	go install

test:
	go test ./...

testrace:
	go test -race ./...

vet:
	go vet ./...
	golint ./...

int:
	${MAKE} -C integration all
