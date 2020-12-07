#!/bin/sh
go build -o main cmd/app/main.go
go build -o authms cmd/sessionblock/main.go
go build -o userblockms cmd/userblock/main.go
