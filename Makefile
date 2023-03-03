.PHONY: help
help:
	@echo 'Targets:'
	@echo '	lint-install	- installs golangci-lint'
	@echo '	lint		- run linter (make sure that the linter is installed before executing this command)'
	@echo '	unittest	- run all unittests (creates coverage report file in ./test)'

.PHONY: lint-install
lint-install:
	@echo "Installing golangci-lint from source"
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: lint
lint:
	golangci-lint run --max-issues-per-linter=0 --max-same-issues=0 --timeout=10m

.PHONY: unittest
unittest:
	if [ -d .test ]; then echo "Removing .test dir" && rm -rf .test; fi
	mkdir .test
	go test ./... -v -coverprofile=.test/coverage.out | grep -v 'no test files'
	go tool cover -html=.test/coverage.out -o .test/coverage.html