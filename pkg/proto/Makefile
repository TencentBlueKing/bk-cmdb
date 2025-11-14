# get all server proto directories
SERVER_DIRS := $(wildcard *-server)

# define the protoc related tools and dependencies path
PROTOC := $(shell which protoc)
PROTOC_VERSION=$(shell protoc --version)
PROTOC_GEN_GO := $(GOPATH)/bin/protoc-gen-go
PROTOC_GEN_GO_GRPC := $(GOPATH)/bin/protoc-gen-go-grpc
PROTOC_GEN_GO_GATEWAY := $(GOPATH)/bin/protoc-gen-go-grpc-gateway
PROTOC_GEN_OPENAPI := $(GOPATH)/bin/protoc-gen-openapiv2
GOOGLE_API_DIR := ./github.com/googleapis/googleapis

# clean previous generated code, install dependencies and generate new code
.DEFAULT_GOAL := default
.PHONY: default
default: clean install-tools generate

# install protoc related tools and dependencies
.PHONY: install-tools
install-tools:
	# check if protoc is installed and version is correct
	@if [ -z "$(PROTOC)" ]; then \
		echo "protoc is not installed，please install protoc first"; \
		exit 1; \
	fi
	@if [ "$(PROTOC_VERSION)" != "libprotoc 32.1" ]; then \
     	echo "protoc version is $(PROTOC_VERSION), expected libprotoc 32.1"; \
		exit 1; \
	fi
	# install protoc-gen-go
	@if [ ! -f "$(PROTOC_GEN_GO)" ]; then \
		go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.36.10; \
	fi
	# install protoc-gen-go-grpc
	@if [ ! -f "$(PROTOC_GEN_GO_GRPC)" ]; then \
		go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1; \
	fi
	# install protoc-gen-go-grpc-gateway
	@if [ ! -f "$(PROTOC_GEN_GO_GATEWAY)" ]; then \
		go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.27.3; \
	fi
	# clone googleapis dependency code
	@if [ ! -d "$(GOOGLE_API_DIR)" ]; then \
		git clone --depth=1 https://github.com/googleapis/googleapis.git $(GOOGLE_API_DIR); \
	fi

# generate code
.PHONY: generate
generate:
	@for dir in $(SERVER_DIRS); do \
		echo "generate code for: $$dir"; \
		protoc \
			--proto_path=$$dir \
			--proto_path=$(GOOGLE_API_DIR) \
			--go_out=paths=source_relative:$$dir \
			--go-grpc_out=paths=source_relative:$$dir \
			--grpc-gateway_out=paths=source_relative:$$dir \
			$$dir/*.proto; \
	done

# clean generated code
.PHONY: clean
clean:
	@for dir in $(SERVER_DIRS); do \
		echo "clean code for: $$dir"; \
		rm -f $$dir/*.pb.go $$dir/*.pb.gw.go; \
	done

# clean installed tools and dependencies
.PHONY: clean-tools
clean-tools:
	@rm -f $(PROTOC_GEN_GO);
	@rm -f $(PROTOC_GEN_GO_GRPC);
	@rm -f $(PROTOC_GEN_GO_GATEWAY);
	@rm -rf $(PROTOC_GEN_OPENAPI);

