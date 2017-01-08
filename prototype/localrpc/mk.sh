#! /bin/sh
#
# mk.sh
# Copyright (C) 2017 cceckman <charles@cceckman.com>
#
# Distributed under terms of the MIT license.
#


 protoc -I proto/ proto/SimpleService.proto --go_out=plugins=grpc:proto/
