APP_NAME  := perf-agent
MODULE    := github.com/paoloelefante/perf-agent
BINDIR    := bin
BIN       := $(BINDIR)/$(APP_NAME)
MAIN_PKG  := ./cmd/perf-agent
PKGS      := ./...

VERSION   ?= 0.1.0
IMAGE     ?= ghcr.io/paoloelefante/perf-agent:$(VERSION)
CHART_DIR := ./charts/perf-agent

GO        ?= go
DOCKER    ?= docker

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make fmt           - Format Go code"
	@echo "  make vet           - Run go vet"
	@echo "  make test          - Run tests"
	@echo "  make build         - Build binary"
	@echo "  make run           - Run locally"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make tidy          - Tidy go modules"
	@echo "  make image         - Build container image with docker"
	@echo "  make push          - Push image to GHCR"
	@echo "  make image-push    - Build and push in one step"
	@echo "  make helm-lint     - Lint Helm chart"
	@echo "  make helm-package  - Package Helm chart"
	@echo "  make all           - fmt + vet + test + build"

.PHONY: fmt
fmt:
	$(GO) fmt $(PKGS)

.PHONY: vet
vet:
	$(GO) vet $(PKGS)

.PHONY: test
test:
	$(GO) test $(PKGS)

.PHONY: build
build:
	mkdir -p $(BINDIR)
	CGO_ENABLED=0 $(GO) build \
		-ldflags="-s -w -X $(MODULE)/internal/version.Version=$(VERSION)" \
		-o $(BIN) $(MAIN_PKG)

.PHONY: run
run:
	$(GO) run $(MAIN_PKG)

.PHONY: tidy
tidy:
	$(GO) mod tidy

.PHONY: clean
clean:
	rm -rf $(BINDIR)

.PHONY: image
image:
	$(DOCKER) build \
		--build-arg VERSION=$(VERSION) \
		-f Dockerfile \
		-t $(IMAGE) .

.PHONY: push
push:
	$(DOCKER) push $(IMAGE)

.PHONY: image-push
image-push: image push

.PHONY: helm-lint
helm-lint:
	helm lint $(CHART_DIR)

.PHONY: helm-package
helm-package:
	helm package $(CHART_DIR) --destination ./dist

.PHONY: all
all: fmt vet test build
