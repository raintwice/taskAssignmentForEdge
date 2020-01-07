all:
	cd client && make
	cd master && make
	cd node && make
.PHONY: clean
clean:
	rm master/master
	rm node/node
	rm client/client
	rm proto/*.pb.go
