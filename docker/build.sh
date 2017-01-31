#!/bin/bash

go build -ldflags '-linkmode external -extldflags -static' -o docker/jamesd cmd/jamesd/main.go
go build -ldflags '-linkmode external -extldflags -static' -o docker/jamesd-ctl cmd/jamesd-ctl/main.go

docker build -t trusch/jamesd:latest docker/
