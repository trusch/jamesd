package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/trusch/jamesd/db"
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/systemstate"
)

var dbUrl = flag.String("db", "localhost", "mongodb url")
var command = flag.String("cmd", "", "one of add-packet, get-packet, list-packets")

var packetName = flag.String("packet-name", "", "name of the packet")
var packetDataFile = flag.String("packet-data", "", "compressed tar archive with packet data")
var preInst = flag.String("preinst", "", "pre-install script")
var postInst = flag.String("postinst", "", "post-install script")
var version = flag.String("version", "", "packet version")
var arch = flag.String("architecture", "", "packet architecture")

var id = flag.String("id", "", "system id")
var file = flag.String("file", "", "system state file")

func addPacket(db *db.DB) {
	if *packetName == "" || *packetDataFile == "" || *version == "" || *arch == "" {
		log.Fatal("specify --packet-name, --packet-data, --version and --architecture")
	}
	pack := &packet.Packet{
		Name:    *packetName,
		Arch:    *arch,
		Version: *version,
	}
	switch filepath.Ext(*packetDataFile) {
	case ".gz":
		pack.Compression = packet.GZIP
	case ".bzip2":
		pack.Compression = packet.BZIP2
	case ".bz2":
		pack.Compression = packet.BZIP2
	case ".lzma":
		pack.Compression = packet.LZMA
	case ".xz":
		pack.Compression = packet.LZMA
	}
	bs, err := ioutil.ReadFile(*packetDataFile)
	if err != nil {
		log.Fatal(err)
	}
	pack.Data = bs
	if *preInst != "" {
		bs, err = ioutil.ReadFile(*preInst)
		if err != nil {
			log.Fatal(err)
		}
		pack.PreInstallScript = string(bs)
	}
	if *postInst != "" {
		bs, err = ioutil.ReadFile(*postInst)
		if err != nil {
			log.Fatal(err)
		}
		pack.PostInstallScript = string(bs)
	}
	err = db.AddPacket(pack)
	if err != nil {
		log.Fatal(err)
	}
}

func getPacket(db *db.DB) {
	if *packetName == "" || *version == "" || *arch == "" {
		log.Fatal("specify --packet-name, --version and --architecture")
	}
	pack, err := db.GetPacket(*packetName, *arch, *version)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("found packet %v_%v:%v data-size: %v preinst-size: %v postinst-size: %v\n",
		*packetName, *arch, *version, len(pack.Data), len(pack.PreInstallScript), len(pack.PostInstallScript))
}

func listPackets(db *db.DB) {
	packets, err := db.ListPackets()
	if err != nil {
		log.Fatal(err)
	}
	for _, packet := range packets {
		fmt.Printf("%v\t%v\t%v\n", packet.Name, packet.Version, packet.Arch)
	}
}

func listSystems(db *db.DB) {
	systems, err := db.ListSystems()
	if err != nil {
		log.Fatal(err)
	}
	for _, system := range systems {
		fmt.Println(system)
	}
}

func getSystemState(db *db.DB) {
	if *id == "" {
		log.Fatal("specify --id")
	}
	state, err := db.GetCurrentSystemState(*id)
	if err != nil {
		log.Fatal(err)
	}
	d, err := yaml.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(d))
}

func getDesiredSystemState(db *db.DB) {
	if *id == "" {
		log.Fatal("specify --id")
	}
	state, err := db.GetDesiredSystemState(*id)
	if err != nil {
		log.Fatal(err)
	}
	d, err := yaml.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(d))
}

func setDesiredSystemState(db *db.DB) {
	if *file == "" {
		log.Fatal("specify --file")
	}
	bs, err := ioutil.ReadFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	systemState := &systemstate.SystemState{}
	err = yaml.Unmarshal(bs, systemState)
	if err != nil {
		log.Fatal(err)
	}
	err = db.SaveDesiredSystemState(systemState)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()
	db, err := db.New(*dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	switch *command {
	case "add-packet":
		{
			addPacket(db)
		}
	case "get-packet":
		{
			getPacket(db)
		}
	case "list-packets":
		{
			listPackets(db)
		}
	case "list-systems":
		{
			listSystems(db)
		}
	case "get-systemstate":
		{
			getSystemState(db)
		}
	case "get-desired-systemstate":
		{
			getDesiredSystemState(db)
		}
	case "set-desired-systemstate":
		{
			setDesiredSystemState(db)
		}
	default:
		{
			log.Fatal("no such command")
		}
	}
}
