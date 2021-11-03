#!/bin/sh

go build -o=s2c ../s2c

rm -rf sfx-exe/ || true

GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/linux-amd64 ./_unpack/main.go
GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/linux-arm64 ./_unpack/main.go
GOOS=linux GOARCH=arm go build -trimpath -ldflags="-s -w" -o ./sfx-exe/linux-arm ./_unpack/main.go
GOOS=linux GOARCH=ppc64le go build -trimpath -ldflags="-s -w" -o ./sfx-exe/linux-ppc64le ./_unpack/main.go
GOOS=linux GOARCH=mips64 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/linux-mips64 ./_unpack/main.go

GOOS=darwin GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/darwin-amd64 ./_unpack/main.go
GOOS=darwin GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/darwin-arm64 ./_unpack/main.go

GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/windows-amd64 ./_unpack/main.go
GOOS=windows GOARCH=386 go build -trimpath -ldflags="-s -w" -o ./sfx-exe/windows-386 ./_unpack/main.go

./s2c -rm -slower sfx-exe/*

rm s2c
