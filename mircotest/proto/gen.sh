#!/bin/bash


#go get github.com/micro/protoc-gen-micro/v3
#go get github.com/micro/protobuf/{proto,protoc-gen-go}
#老版本
#protoc -I . --proto_path=$GOPATH/src:.  --micro_out=. --go_out=. test.proto
protoc --go_out=. --micro_out=. test.proto