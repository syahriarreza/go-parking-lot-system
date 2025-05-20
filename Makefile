# Makefile for go-parking-lot-system

.PHONY: test coverage coverage-html test-coverage

# Run unit tests with verbose output
test:
	go test -v ./...

# Run unit tests with coverage report
coverage:
	go test -coverprofile=coverage.out ./...

# Generate HTML coverage report
coverage-html:
	go tool cover -html=coverage.out

# Run tests with coverage and open HTML report
test-coverage: coverage coverage-html
