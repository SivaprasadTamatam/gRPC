#!/bin/bash

protoc greet/greetpb/greet.proto --go_out=plugins=grpc:.

protoc calculator/calculatorpb/calculator.proto --go_out=plugins=grpc:.

export GO_PATH=~/go
export PATH=$PATH:/$GO_PATH/bin