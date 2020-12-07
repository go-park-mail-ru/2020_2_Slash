#!/bin/sh
go run cmd/app/main.go &
go run cmd/userblock/main.go &
go run cmd/sessionblock/main.go &
