#!/bin/bash

go test . -cover -coverpkg configcenter/src/sence-server/host-server/host-service/actions/... -coverprofile=cover.prof
go tool cover -html cover.prof -o cover.html
