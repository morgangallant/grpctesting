#!/bin/bash

rm -f pb/*.pb.go
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
for proto in pb/protobufs/*.proto; do
	protoc -I . -I ~/go/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis $proto --go_out=pb/ --go-grpc_out=pb/ --grpc-gateway_out=logtostderr=true:pb/
done