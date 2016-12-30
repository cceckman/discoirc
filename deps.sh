#! /bin/sh
#
# Checks for / gets dependencies, and builds the binaries.

go get \
  github.com/golang/protobuf/proto \
  github.com/golang/protobuf/protoc-gen-go \
  github.com/Shopify/go-lua  \
  github.com/fluffle/goirc  \
  github.com/jroimartin/gocui

