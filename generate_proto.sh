#!/bin/bash
go install google.golang.org/protobuf/cmd/protoc-gen-go
protoc --go_out=. message.proto