PROTO_PATH = ../proto

all:build

build:$(PROTO_PATH)/message.pb.go
	go build client.go

$(PROTO_PATH)/%.pb.go:$(PROTO_PATH)/%.proto
	cd $(PROTO_PATH) && protoc --go_out=plugins=grpc:. *.proto

.PHONY: clean
clean:
	rm client