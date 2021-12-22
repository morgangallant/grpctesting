#!/bin/bash

rm -f pb/*.pb.go
go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
for proto in pb/azad/*.proto; do
	protoc -I pb/ $proto --go_out=pb/ --go-grpc_out=pb/ --grpc-gateway_out=logtostderr=true:pb/
done