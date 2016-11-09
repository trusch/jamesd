package main

import (
	"flag"
	"log"

	"github.com/trusch/jamesd"
)

var config = flag.String("config", "config.yaml", "client config file (yaml or json)")

func main() {
	log.SetFlags(log.Lshortfile)
	flag.Parse()
	cfg := jamesd.ClientConfig{}
	err := cfg.ParseFile(*config)
	if err != nil {
		log.Fatal(err)
	}
	client, err := jamesd.NewClient(cfg.Endpoint, cfg.Cert, cfg.Key, cfg.CA, cfg.InstallRoot)
	if err != nil {
		log.Fatal(err)
	}
	for _, service := range cfg.Services {
		request := &jamesd.Request{
			Package: service.Name,
			Version: service.Version,
			Arch:    cfg.Architecture,
		}
		response, err := client.GetPackage(request)
		if err != nil {
			log.Fatal(err)
		}
		err = client.InstallPackage(response)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("successfully installed %v:%v", response.Name, response.Arch)
	}
}
