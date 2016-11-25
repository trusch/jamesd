package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/trusch/jamesd/client"
)

var addr = flag.String("addr", "localhost:2761", "address to connect to on")
var keyFile = flag.String("key", "jamesc.key", "key to use")
var certFile = flag.String("cert", "jamesc.crt", "cert to use")
var caFile = flag.String("ca", "ca.crt", "CA to use")
var installRoot = flag.String("install-root", "/", "CA to use")
var stateFile = flag.String("state-file", "/var/lib/jamesc/state.gob", "statefile location")
var tagList = flag.String("system-tags", "", "comma separated list of system tags")

var tags []string

func init() {
	flag.Parse()
	if *tagList != "" {
		tags = strings.Split(*tagList, ",")
	}
}

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	for {
		cli, err := client.New(*addr, *certFile, *keyFile, *caFile, *installRoot, *stateFile, tags)
		if err != nil {
			log.Print("Error while establishing jamesd connection, retry in 5 seconds...")
			time.Sleep(5 * time.Second)
			continue
		}
		log.Print("got jamesd connection, starting normal operation.")
		err = cli.Run()
		if err != nil {
			log.Print("Error in jamesd connection, retry in 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
}
