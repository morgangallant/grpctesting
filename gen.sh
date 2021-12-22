#!/bin/bash

rm -f pb/*.pb.go
go install google.golang.org/protobuf/cmd/protoc-gen-go
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
for proto in pb/azad/*.proto; do
	protoc -I pb/ $proto --go_out=pb/ --go-grpc_out=pb/ --grpc-gateway_out=logtostderr=true:pb/
done