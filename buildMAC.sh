#!/bin/sh
export GOARCH="amd64"
export GOOS="darwin"
go build code.go
echo "Done!"
