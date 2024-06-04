GO=GO111MODULE=on go
GOBUILD=$(GO) build

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: build

build:
	$(GOBUILD) ./cli/gocomply_xsd2go

.PHONY: pkger vendor
pkger:
ifeq ("$(wildcard $(GOBIN)/pkger)","")
	go install github.com/markbates/pkger/cmd/pkger@latest
endif

ci-update-bundled-deps: pkger
	$(GOBIN)/pkger -o pkg/template
	go fmt ./pkg/template

generate-test-schemas: build $(XSD_TEST_DIR)/*.xsd
	@for xsd in $(shell ls $(XSD_TEST_DIR)/*.xsd); do \
		./gocomply_xsd2go convert $${xsd} "" $(TEST_MODELS_DIR); \
	done

check: generate-test-schemas
	$(GO) test -v -cover -coverpkg=./... -tags=xml_test ./tests

vendor:
	$(GO) mod tidy
	$(GO) mod vendor
	$(GO) mod verify
