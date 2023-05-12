# Some Hacking with QUIC in Go

Playing with QUIC transport in Go

Checkout the Makefile

`make build` : This build 2 binaries, `./cli` and `./svr`

Start the server first, and then you can run the client

Its pretty basic but shows using the quic transport in golang using the quic-go library

Note, currently the server uses go code to manage the required TLS crypto for QUIC, down the road ill likely enable option to use generated certs (e.g., `make gen-certs`)