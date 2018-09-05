#!/bin/sh

PACKAGES=$(go list ../src/...)

#test:
	echo "go test"
	go test -cover=true ../src/...

#collect-cover-data:
	echo "collect-cover-data"
	echo "mode: count" > coverage-all.out
	for pkg in $PACKAGES;do
		echo ${pkg}
		go test -v -coverprofile=coverage.out -covermode=count ${pkg} || exit $$?;\
		if [ -f coverage.out ]; then\
			tail -n +2 coverage.out >> coverage-all.out;\
		fi\
	done

#test-cover-html:
	echo "test-cover-html"
	go tool cover -html=coverage-all.out -o coverage.html

#test-cover-func:
	echo "test-cover-func"
	go tool cover -func=coverage-all.out 
