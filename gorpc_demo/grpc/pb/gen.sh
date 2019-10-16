#!/bin/bash
PROTOC=`which protoc`
$PROTOC  --go_out=plugins=grpc:. search.proto