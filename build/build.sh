#!/bin/sh
go build -o bin/main cmd/app/main.go
go build -o bin/authms cmd/sessionblock/main.go
go build -o bin/userblockms cmd/userblock/main.go
