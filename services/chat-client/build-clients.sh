#!/bin/bash

protoc -I=./../../protos ../../protos/*.proto \
    --grpc-web_out=import_style=typescript,mode=grpcweb:./grpc-clients