all:
	cd master && make
	cd node && make
	cd client && make
.PHONY: clean
clean:
	rm master/master
	rm node/node
	rm proto/*.pb.go
