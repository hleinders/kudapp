#!/bin/bash

go mod tidy

go vet ./...

which govulncheck >/dev/null 2>&1
[[ $? -eq 0 ]] && { govulncheck ./... || exit $?; }

# build static to avoid problems e.g. with alpine
CGO_ENABLED=0 go build -o kudapp


# docker build . -t kudapp:red --build-arg DEAFULT_COLOR=red
# docker build . -t kudapp:blue --build-arg DEAFULT_COLOR=blue
# docker build . -t kudapp:green --build-arg DEAFULT_COLOR=green

# docker tag kudapp:red kudapp:latest
