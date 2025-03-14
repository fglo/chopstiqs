GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=chopstiqs
PATH_TO_MAIN_GO=./_examples/demo/main.go
PATH_TO_MAIN_GO2=./_examples/simple/main.go
OUT_PATH=out/bin/

GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
CYAN   := $(shell tput -Txterm setaf 6)
RESET  := $(shell tput -Txterm sgr0)

.PHONY: all test build vendor

all: help

## Build:
run-demo: ## Run project
	go run -race $(PATH_TO_MAIN_GO)

run-simple-demo: ## Run project
	go run -race $(PATH_TO_MAIN_GO2)

build-demo: ## Build demo project and put the output binary in out/bin/
	mkdir -p $(OUT_PATH)
#    GO111MODULE=on $(GOCMD) build -mod vendor -o $(OUT_PATH)/$(BINARY_NAME) $(PATH_TO_MAIN_GO)
	GO111MODULE=on $(GOCMD) build -o $(OUT_PATH)$(BINARY_NAME) $(PATH_TO_MAIN_GO)

build-demo-wasm:
	GOOS=js GOARCH=wasm $(GOCMD) build -tags wasm -o static/main.wasm $(PATH_TO_MAIN_GO)

wasmserve-demo: ## Run demo app as webapp and expose it under http://localhost:8080
	wasmserve -tags wasm $(PATH_TO_MAIN_GO)

clean: ## Remove build related file
	rm -fr ./bin
	rm -fr ./out

vendor: ## Copy of all packages needed to support builds and tests in the vendor directory
	$(GOCMD) mod vendor

watch-demo: ## Run the code with cosmtrek/air to have automatic reload on changes
	air  --build.cmd "go build -o out/bin/$(BINARY_NAME) $(PATH_TO_MAIN_GO)" --build.bin "./out/bin/chopstiqs"

## Test:
test: ## Run the tests of the project
	$(GOTEST) -p=1 -v -race ./... $(OUTPUT_OPTIONS)

## Format:
tidy: ## go mod tidy
	go mod tidy

fmt: ## go mod fmt
	go fmt ./...

vet: ## go mod vet
	go vet ./...

format: tidy fmt vet ## format

## Help:
help: ## Show this help.
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)

