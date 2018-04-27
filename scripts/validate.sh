#!/bin/bash
# set -e

cd $1


echo "Running validation in $1"

PACKAGES="$(find . -name '*.go' | grep -Ev '.pb.go' | xargs -I{} dirname {} | \
	sort -u | grep -Ev '(.git|.trash-cache|vendor|bin|k8s/pkg)')"

echo Packages: ${PACKAGES}
echo Running: go vet
vetout=$(go tool vet ${PACKAGES} 2>&1 | tee /dev/stderr )
if [ -n "$vetout" ]; then
    echo go vet failed
    exit 1
fi

echo Running: golint
for i in ${PACKAGES}; do
    if [ -n "$(golint $i | grep -v 'should have comment.*or be unexported' | grep -v 'just return error instead. ' | tee /dev/stderr)" ]; then
        failed=true
    fi
done

if test -n "$failed"
then
    echo go golint failed
    exit 1
fi

echo Running: go fmt
if test -n "$(gofmt -l ${PACKAGES} | tee /dev/stderr)"
then
    echo go fmt failed
    exit 1
fi
