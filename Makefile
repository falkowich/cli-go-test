# Define project-specific variables
PROJECT_NAME := cli-go-test
VERSION := $(shell git describe --tags --always)
BUILD_DIR := releases/$(VERSION)

# Define Go-related variables
GO_FILES := $(shell find . -type f -name '*.go')

# Platforms to build for
PLATFORMS := linux/amd64 windows/amd64 darwin/amd64

.PHONY: all build bump clean release tag

# Default target
all: build

# Build binaries for specified platforms
build:
	@mkdir -p $(BUILD_DIR)
	@for platform in $(PLATFORMS); do \
		OS=`echo $$platform | cut -d'/' -f1`; \
		ARCH=`echo $$platform | cut -d'/' -f2`; \
		EXT=""; \
		if [ "$$OS" = "windows" ]; then EXT=".exe"; fi; \
		OUTPUT=$(BUILD_DIR)/$(PROJECT_NAME)-$$OS-$$ARCH$$EXT; \
		echo "Building $$OUTPUT"; \
		GOOS=$$OS GOARCH=$$ARCH go build -o $$OUTPUT .; \
	done

# Bump version (major, minor, or patch)
bump:
	@./scripts/bump_version.sh $(TYPE)

# Clean up generated files
clean:
	@rm -rf $(BUILD_DIR)

# Create a release: build, generate checksums, and tag
release: build checksums tag

# Generate checksums for the binaries
checksums:
	@cd $(BUILD_DIR) && sha256sum * > checksums.txt

# Tag the release in git
tag:
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)

# Utility target to show help
help:
	@echo "Usage:"
	@echo "  make build     - Build the project for multiple platforms"
	@echo "  make clean     - Clean the build artifacts"
	@echo "  make release   - Build, generate checksums, and tag the release"
	@echo "  make tag       - Tag the current version in git"
	@echo "  make help      - Show this help message"

