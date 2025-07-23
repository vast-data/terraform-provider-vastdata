VERSION ?= $(shell cat version/VERSION | tr -d '[:space:]')
GOFMT_FILES ?= $$(find . -name '*.go')
PKG_NAME = vastdata
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

EXT :=
ifeq ($(GOOS),windows)
EXT := .exe
endif

BINARY_NAME := terraform-provider-$(PKG_NAME)
BUILD_DIR := build/$(GOOS)_$(GOARCH)
BINARY := $(BUILD_DIR)/$(BINARY_NAME)$(EXT)
GOFLAGS := -mod=readonly
LDFLAGS := -X main.version=$(VERSION)

.PHONY: build show-resources show-datasources

build:
	@echo "Building $(BINARY_NAME) for $(GOOS)_$(GOARCH)..."
	mkdir -p $(BUILD_DIR)
	env GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 \
	go build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY)
	@echo "Binary created at $(BINARY)"
	@chmod +x $(BINARY)


vet:
	@echo "go vet ."
	@go vet $$(go list ./...) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -s -w $(GOFMT_FILES)


ifeq (show, $(firstword $(MAKECMDGOALS)))
  runargs := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(foreach arg,$(runargs),$(eval $(arg):;@true))
endif

#   Run the Go schema visualization tool with dynamic arguments.
#   Usage:
#     make show <type> [filter] [automation]
#
#   Arguments:
#     <type>       - Required. The type of schema to show: 'resource', 'datasource', or short forms ('r', 'd').
#     [filter]     - Optional. Case-insensitive substring to filter by.
#     [automation] - Optional. If present (any value), enables automation mode (🔸 markers).
#
#   Examples:
#     make show resource
#     make show datasource user
#     make show r quota automation
show:
	@type=$(word 1, $(runargs)); \
	arg2=$(word 2, $(runargs)); \
	arg3=$(word 3, $(runargs)); \
	filter=""; \
	auto_flag=""; \
	[ "$$arg2" = "automation" ] && auto_flag="-automation" || filter="$$arg2"; \
	[ "$$arg3" = "automation" ] && auto_flag="-automation"; \
	echo "Running show with type=$$type, filter=$$filter, automation=$$auto_flag"; \
	go run $(CURDIR)/misc/show_schemas.go -type=$$type $${filter:+-filter=$$filter} $$auto_flag


test:
	go test -v -cover ./...


# Generate docs and copywrite headers
generate-docs:
	@echo "Generating documentation and headers..."
	@go run $(CURDIR)/misc/generate_examples.go
	cd tools; go generate ./...; cd -


ifeq (gen-openapi-tar, $(firstword $(MAKECMDGOALS)))
  runargs := $(wordlist 2, $(words $(MAKECMDGOALS)), $(MAKECMDGOALS))
  $(foreach arg,$(runargs),$(eval $(arg):;@true))
endif

#   Convert an OpenAPI YAML file to an OpenAPI v3 tarball and JSON file.
#   Usage:
#     make gen-openapi-tar runargs=<path>
#
#   Arguments:
#     <path> - Required. Path to the Swagger/OpenAPI YAML file.
#
#   This target:
#     - Converts the YAML file to Swagger JSON using Python.
#     - Converts Swagger JSON to OpenAPI v3 using Go.
#     - Creates a tarball containing the final JSON.
#     - Copies both `openapi.tar.gz` and `openapi.json` to the current directory.
#
#   Examples:
#     make gen-openapi-tar runargs=./specs/swagger.yaml
#     make gen-openapi-tar runargs=/tmp/vast-openapi.yaml
gen-openapi-tar:
	@set -e; \
	path=$(word 1, $(runargs)); \
	tmp_dir="/tmp/apiconv"; \
	json_out="$$tmp_dir/api.json"; \
	tarball="$$tmp_dir/api.tar.gz"; \
	dest_dir="$(CURDIR)"; \
	\
	echo "Converting YAML to JSON..."; \
	python3 $(CURDIR)/misc/yaml2json.py $$path -o $$tmp_dir/swagger.json; \
	\
	echo "Converting JSON to OpenAPI v3 and creating tarball..."; \
	go run $(CURDIR)/misc/convert_to_v3.go; \
	\
	echo "Copying outputs to current directory..."; \
	cp $$tarball $$dest_dir/openapi.tar.gz; \
	cp $$json_out $$dest_dir/openapi.json; \
	echo "Copied: openapi.tar.gz and openapi.json to $$dest_dir"
