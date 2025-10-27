GOBUILD                  = GOEXPERIMENT=jsonv2 CGO_ENABLED=0 go build -trimpath
PREFIX                  ?= $(shell pwd)
BUILD_PATH               = ${PREFIX}/build
lint                     = ${PREFIX}/bin/lint
GOLANGCI_VERSION         = v2.4.0

ifdef VERSION
	VERSION=$(VERSION)
else
	VERSION=$(shell git describe --tags --always)
endif

BUILDTIME=$(shell date +%Y-%m-%dT%T%z)
GITHASH=$(shell git rev-parse HEAD)

LDFLAG=-ldflags "-X github.com/TencentBlueKing/bk-cmdb/pkg/version.Version=${VERSION} \
-X github.com/TencentBlueKing/bk-cmdb/pkg/version.BuildTime=${BUILDTIME} \
-X github.com/TencentBlueKing/bk-cmdb/pkg/version.GitHash=${GITHASH}"

${lint}:
	@echo ">> downloading golangci"
	@mkdir -p ${PREFIX}/bin/lint
	@curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/$(GOLANGCI_VERSION)/install.sh | sh -s -- -b ${lint} ${GOLANGCI_VERSION} && chmod +x ${lint}

.PHONY: lint
lint: ${lint}
	@echo ${lint}/golangci-lint
	GOEXPERIMENT=jsonv2 ${lint}/golangci-lint run ./...

.PHONY: test
test:
	GOEXPERIMENT=jsonv2 go test -count=1 -cover ./...

.PHONY: build
build:
	@echo ">> building apiserver"
	${GOBUILD} ${LDFLAG} -o ${BUILD_PATH}/cmdb_apiserver cmd/api_server/*.go
