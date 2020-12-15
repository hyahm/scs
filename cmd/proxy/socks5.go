package main

import (
	"flag"

	"github.com/armon/go-socks5"
)

var l = flag.String("l", ":8080", "listen")

func main() {
	flag.Parse()
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port :8080
	if err := server.ListenAndServe("tcp", *l); err != nil {
		panic(err)
	}
}
