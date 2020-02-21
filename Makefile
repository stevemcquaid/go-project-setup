SOURCE_FILES?=$$(go list ./... | grep -v /vendor/)
TEST_PATTERN?=.
TEST_OPTIONS?=

.PHONY: build test
.DEFAULT_GOAL := build

## Build docker container
docker-build:
	@docker build --target final -t stevemcquaid/hello-go:latest .

## Run docker container
docker-run: docker-build
	@docker run -it -p 5000:5000 stevemcquaid/hello-go:latest .

## Build the binary
build:
	@go build -o bin/app ./

## Run the binary
run:
	@go run main.go

## Install all the build and lint dependencies
setup:
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u github.com/pierrre/gotestcover
	go get -u golang.org/x/tools/cmd/cover
	go get -u golang.org/x/tools/cmd/goimports
	mkdir build

## Run all the tests
test: test-unit

test-unit: ## Run all unit tests
	go test -cover $(TEST_OPTIONS) -covermode=atomic -coverprofile=build/unit.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m


test-integration: ## Run all integration tests
	go test -cover $(TEST_OPTIONS) -tags=integration -covermode=atomic -coverprofile=build/integration.out ./... -run $(TEST_PATTERN) -timeout=2m

## Run all the tests and opens the coverage report
cover: test
	# gocovmerge build/unit.out build/integration.out > build/all.out
	gocovmerge build/unit.out > build/all.out
	go tool cover -html=build/all.out


fmt: ## gofmt and goimports all go files
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

## Run all the linters
lint:
	golangci-lint run  \
		--deadline=30m \
		--disable-all  \
		--no-config  \
		--issues-exit-code=0  \
		--enable=deadcode  \
		--enable=dupl  \
		--enable=errcheck  \
		--enable=goconst  \
		--enable=gocyclo \
		--enable=goimports \
		--enable=golint \
		--enable=gosec  \
		--enable=govet \
		--enable=ineffassign \
		--enable=interfacer  \
		--enable=maligned  \
		--enable=megacheck \
		--enable=misspell \
		--enable=staticcheck \
		--enable=structcheck  \
		--enable=typecheck \
		--enable=unconvert  \
		--enable=varcheck
#		--enable=whitespace


deps: ## Download deps & Runs `go mod tidy`
	@go mod tidy
	@go mod download
	@go mod vendor


vet: ## Verifies `go vet` passes
	@go vet $(shell go list ./... | grep -v vendor) | grep -v '.pb.go:' | tee /dev/stderr

## Prep for commit - run make fmt, vendor, tidy
clean: fmt deps vet

## Run all the tests and code checks
ci: fmt lint deps vet test

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
#help:
#	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

help:
	@printf "Available targets:\n\n"
	@awk '/^[a-zA-Z\-\_0-9%:\\]+:/ { \
	  helpMessage = match(lastLine, /^## (.*)/); \
	  if (helpMessage) { \
	    helpCommand = $$1; \
	    helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
      gsub("\\\\", "", helpCommand); \
      gsub(":+$$", "", helpCommand); \
	    printf "  \x1b[32;01m%-35s\x1b[0m %s\n", helpCommand, helpMessage; \
	  } \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST) | sort -u
	@printf "\n"

