#!/bin/bash
protoc --go_out=plugins=grpc:. message.proto
ls -l
exit 
