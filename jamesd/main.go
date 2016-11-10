package main

import (
	"flag"
	"log"

	"github.com/trusch/jamesd/server"
)

var addr = flag.String("addr", ":2761", "address to listen on")
var dbAddr = flag.String("db", "localhost", "mongodb address")
var keyFile = flag.String("key", "jamesd.key", "key to use")
var certFile = flag.String("cert", "jamesd.crt", "cert to use")
var caFile = flag.String("ca", "ca.crt", "CA to use")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	server := server.New(*addr, *certFile, *keyFile, *caFile, *dbAddr)
	server.Run()
}
