#!/bin/bash

RELEASE=${1:-"v3.6.3"}
# docker build . -t bk-cmdb-dev:${RELEASE} --build-arg branch=release-${RELEASE}

id=$(docker create bk-cmdb-dev:${RELEASE})
docker cp $id:/data/bin/bk-cmdb .
docker rm -v $id
docker build . -t bk-cmdb:${RELEASE} -f Dockerfile.product
rm -rf ./bk-cmdb
