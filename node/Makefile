PROTO_PATH = ../proto

all:build

build:$(PROTO_PATH)/message.pb.go
	go build node.go
	GOARCH=arm GOOS=linux go build -o nodearm node.go

$(PROTO_PATH)/%.pb.go:$(PROTO_PATH)/%.proto
	cd $(PROTO_PATH) && protoc --go_out=plugins=grpc:. *.proto

.PHONY: clean
clean:
	rm node

.PHONY: arm
arm:
	GOARCH=arm GOOS=linux go build -o nodearm node.go
