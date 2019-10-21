PROTO_PATH = ../proto

all: node

node:$(PROTO_PATH)/message.pb.go
	go build node.go

$(PROTO_PATH)/%.pb.go:$(PROTO_PATH)/%.proto
	cd $(PROTO_PATH) && protoc --go_out=plugins=grpc:. *.proto

.PHONY: clean
clean:
	rm node
