PROTO_PATH = ../proto

all: master

master:$(PROTO_PATH)/message.pb.go
	go build master.go

$(PROTO_PATH)/%.pb.go:$(PROTO_PATH)/%.proto
	cd $(PROTO_PATH) && protoc --go_out=plugins=grpc:. *.proto

.PHONY: clean
clean:
	rm master
