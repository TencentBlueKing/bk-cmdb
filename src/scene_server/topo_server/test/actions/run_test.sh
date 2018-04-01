#!/bin/bash

go test . -cover -coverpkg configcenter/src/sence-server/topo-server/topo-service/actions/... -coverprofile=cover.prof
go tool cover -html cover.prof -o cover.html
