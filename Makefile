.PHONY: help run test

.DEFAULT_GOAL := help

help: ## display help
	@echo "Makefile available targets:"
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  * \033[36m%-15s\033[0m %s\n", $$1, $$2}'


run: ## run nakama server
	docker-compose up --build

test: ## run all tests
	go clean -testcache
	go test -race ./...
