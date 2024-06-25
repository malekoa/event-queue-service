# Define variables
TEST_OUTPUT_FILE = coverage.out

.PHONY: all test test-output clean

# Default target
all: test

# Run tests
test:
	@echo "Running tests..."
	go test ./... -v

# Run tests and generates coverage report
test-output:
	@echo "Running tests and generating HTML output..."
	@go test -coverprofile=coverage.out
	@go tool cover -html=coverage.out

# update the coverage badge
update-coverage-badge:
	@echo "Updating coverage badge..."
	@chmod +x ./update-coverage-badge.sh
	@./update-coverage-badge.sh
	@rm -f $(TEST_OUTPUT_FILE)

# Clean up
clean:
	@echo "Cleaning up..."
	@rm -f $(TEST_OUTPUT_FILE)
