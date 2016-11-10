package main

import (
	"flag"
	"log"

	"github.com/trusch/jamesd/client"
)

var addr = flag.String("addr", "localhost:2761", "address to connect to on")
var keyFile = flag.String("key", "jamesc.key", "key to use")
var certFile = flag.String("cert", "jamesc.crt", "cert to use")
var caFile = flag.String("ca", "ca.crt", "CA to use")
var installRoot = flag.String("install-root", "/", "CA to use")
var stateFile = flag.String("state-file", "/var/lib/jamesc/state.gob", "statefile location")
var arch = flag.String("arch", "armv7l", "host architecture")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	cli, err := client.New(*addr, *certFile, *keyFile, *caFile, *installRoot, *stateFile, *arch)
	if err != nil {
		log.Fatal(err)
	}
	cli.Run()
}
