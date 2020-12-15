#!/bin/sh
go test ./... -coverprofile cover.out.tmp
cat cover.out.tmp | grep -v ".pb.go" >cover.out
go tool cover -func cover.out
