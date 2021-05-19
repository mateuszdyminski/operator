GOLANG_VERSION := 1.16.4
ALPINE_VERSION := 3.12

VERSION ?= $(shell git describe --tags --always)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
LAST_COMMIT_USER ?= $(shell git log -1 --format='%cn <%ce>')
LAST_COMMIT_HASH ?= $(shell git log -1 --format=%H)
LAST_COMMIT_TIME ?= $(shell git log -1 --format=%cd --date=format:'%Y-%m-%d %H:%M:%S')

GIT_REPO := github.com/mateuszdyminski/operator/usersDB
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*" -not -path "./swagger/*" -not -path "./.go/*")
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)

NAME ?= users-db

.DEFAULT_GOAL := all
.PHONY: all misspell lint test ci-check race vet swag fmt build-healthcheck build docker-build docker-push deps download run version help

all: deps fmt vet misspell lint race build ## Combines `deps` `fmt` `vet` `misspell` `lint` `race` `build` commands

misspell: ## Runs misspell
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo 'Please install "misspell" tool: https://github.com/client9/misspell'; \
		exit 1; \
	fi
	misspell -w $(GOFILES)

lint: ## Runs golangci-lint
	@hash golangci-lint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo 'Please install "golangci-lint" tool: https://github.com/golangci/golangci-lint'; \
		exit 1; \
	fi
	golangci-lint run -v

test: ## Runs all unit tests and generates coverage report
	go test -coverprofile coverage.out -v ./...
	go tool cover -html=coverage.out -o coverage_report.html

ci-check: ## Checks if files in repository are modified
	@[ -z "$$(git status --porcelain)" ] || \
		(echo 'Repository is unclean - please fix fmt/vet/misspell/go mod issues and commit again'; git diff; exit 1)

race: ## Runs all unit tests with -race flag
	go test -race -coverprofile coverage.out -v ./...
	go tool cover -html=coverage.out -o coverage_report.html

vet: ## Runs `go vet` for all packages
	go vet $(PACKAGES)

fmt: ## Formats all Golang files
	gofmt -l -s -w $(GOFILES)

build: ## Builds app binary
	CGO_ENABLED=0 \
	go build \
	-v \
	-ldflags "-s -w \
	-X '$(GIT_REPO)/domain/appstate.AppName=$(NAME)' \
	-X '$(GIT_REPO)/domain/appstate.AppVersion=$(VERSION)' \
	-X '$(GIT_REPO)/domain/appstate.APIVersion=$(VERSION)' \
	-X '$(GIT_REPO)/domain/appstate.BuildTime=$(BUILD_TIME)' \
	-X '$(GIT_REPO)/domain/appstate.LastCommitUser=$(LAST_COMMIT_USER)' \
	-X '$(GIT_REPO)/domain/appstate.LastCommitHash=$(LAST_COMMIT_HASH)' \
	-X '$(GIT_REPO)/domain/appstate.LastCommitTime=$(LAST_COMMIT_TIME)'" \
	-o $(NAME)-$(VERSION) .

docker-build: ## Builds Docker image with App
	docker build \
	--pull \
	--build-arg http_proxy \
	--build-arg https_proxy \
	--build-arg no_proxy \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg NAME="$(NAME)" \
	--build-arg VERSION="$(VERSION)" \
	--build-arg BUILD_TIME="$(BUILD_TIME)" \
	--build-arg LAST_COMMIT_USER="$(LAST_COMMIT_USER)" \
	--build-arg LAST_COMMIT_HASH="$(LAST_COMMIT_HASH)" \
	--build-arg LAST_COMMIT_TIME="$(LAST_COMMIT_TIME)" \
	--label="build.version=$(VERSION)" \
	--label="build.time=$(BUILD_TIME)" \
	--label="commit.user=$(LAST_COMMIT_USER)" \
	--label="commit.hash=$(LAST_COMMIT_HASH)" \
	--label="commit.time=$(LAST_COMMIT_TIME)" \
	--tag="$(NAME):$(VERSION)" \
	.

docker-push: ## Logins, tags and pushes current version of Docker image to the registry
	docker tag "$(NAME):$(VERSION)" "mateuszdyminski/$(NAME):$(VERSION)"
	docker push "mateuszdyminski/$(NAME):$(VERSION)"

deps: ## Synchronises all dependencies
	go mod tidy

download: ## Downloads all dependencies
	go mod download

run: ## Runs App in development mode locally
	DEBUG="false" \
	HTTP_PORT="8080" \
	DB_HOST="127.0.0.1:30000" \
	DB_USER="root" \
	DB_PASS="password" \
	DB_NAME="users" \
	DB_MIGRATIONS_PATH="migrations" \
	go run .

version: ## Displays current project version
	@echo $(VERSION)

help: ## Displays this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
