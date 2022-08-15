help:: ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | sort | fgrep -v fgrep | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

test: ## Go test
	go test ./...

test_coverage: ## Go test with coverage file
	go test ./... -coverprofile=coverage.out

get_dependencies: ## Go get dependencies
	 go get -v -t -d ./...