SHELL := /bin/bash

.PHONY: help
help:
	@echo "Usage make <TARGET>"
	@echo ""
	@echo "  Targets:"
	@echo "	   gen-certs			Generate Certificates"
	@echo "	   build				Build a server executable"
	@echo "	   run					Run the server from code"
	@echo "	   docker-server		Containerize the server"
	@echo "	   docker-run			Run server container"
	@echo "	   kind					Add to kind so that we can setup in ingress controller"
	@echo "	   test-server			Test server via curl"


.PHONY: gen-certs
gen-certs:
# https://centrifugal.dev/blog/2020/10/16/experimenting-with-quic-transport
	openssl genrsa -des3 -passout pass:x -out ./certs/server.pass.key 2048
	openssl rsa -passin pass:x -in ./certs/server.pass.key -out ./certs/server.key
	openssl req -new -key ./certs/server.key \
			-out ./certs/server.csr \
			-subj "/C=US/ST=PA/L=PHL/O=DREXEL/OU=ComputerScience/CN=LocalQuicTesting"

.PHONY: build
build:
	go build 

.PHONY: run
run:
	go run *.go