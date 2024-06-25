#!/bin/bash

# Run tests and generate coverage profile
go test -coverprofile=coverage.out ./...
# get the coverage percentage without the percentage sign
coverage=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "Captured coverage percentage: $coverage"
# replace first line of README.md with new coverage percentage
sed -i '' "1s/.*/[![Coverage Status](https:\/\/img.shields.io\/badge\/coverage-$coverage%25-brightgreen)](https:\/\/github.com\/username\/repo)/" README.md