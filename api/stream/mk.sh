#! /bin/sh


 protoc -I . stream.proto --go_out=plugins=grpc:.
