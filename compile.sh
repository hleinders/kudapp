#!/bin/bash

go mod tidy

go vet ./...

which govulncheck >/dev/null 2>&1
[[ $? -eq 0 ]] && { govulncheck ./... || exit $?; }

mkdir -p build/linux/{amd64,arm64}
# build static to avoid problems e.g. with alpine
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/linux/amd64/kudapp
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o build/linux/arm64/kudapp
