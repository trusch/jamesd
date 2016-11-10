package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/trusch/jamesd/db"
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/systemstate"
)

var dbUrl = flag.String("db", "localhost", "mongodb url")
var command = flag.String("cmd", "", "one of add-packet, get-packet, remove-packet, list-packets, list-systems, get-state, get-desired-state, set-desired-state")

var packetName = flag.String("name", "", "name of the packet")
var tagList = flag.String("tags", "", "comma separated list of tags")
var packetDataFile = flag.String("data", "", "compressed tar archive with packet data")
var preInst = flag.String("preinst", "", "pre-install script")
var postInst = flag.String("postinst", "", "post-install script")
var preRemove = flag.String("prerm", "", "pre-install script")
var postRemove = flag.String("postrm", "", "post-install script")

var id = flag.String("id", "", "system id")
var file = flag.String("file", "", "system state file")

var tags []string

func init() {
	flag.Parse()
	if *tagList != "" {
		tags = strings.Split(*tagList, ",")
	}
}

func addPacket(db *db.DB) {
	if *packetName == "" || *packetDataFile == "" {
		log.Fatal("specify --name, --data")
	}
	pack := &packet.Packet{
		Name: *packetName,
		Tags: tags,
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
	if *preRemove != "" {
		bs, err = ioutil.ReadFile(*preRemove)
		if err != nil {
			log.Fatal(err)
		}
		pack.PreRemoveScript = string(bs)
	}
	if *postRemove != "" {
		bs, err = ioutil.ReadFile(*postRemove)
		if err != nil {
			log.Fatal(err)
		}
		pack.PostRemoveScript = string(bs)
	}
	err = db.AddPacket(pack)
	if err != nil {
		log.Fatal(err)
	}
}

func getPacket(db *db.DB) {
	if *packetName == "" {
		log.Fatal("specify --name")
	}
	pack, err := db.GetPacket(*packetName, tags)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v tags: %v data-size: %v preinst-size: %v postinst-size: %v prerm-size: %v postrm-size: %v\n",
		pack.Name, pack.Tags, len(pack.Data), len(pack.PreInstallScript), len(pack.PostInstallScript), len(pack.PreRemoveScript), len(pack.PostRemoveScript))
}

func removePacket(db *db.DB) {
	if *packetName == "" {
		log.Fatal("specify --name")
	}
	err := db.RemovePacket(*packetName, tags)
	if err != nil {
		log.Fatal(err)
	}
}

func listPackets(db *db.DB) {
	packets, err := db.ListPackets(*packetName, tags)
	if err != nil {
		log.Fatal(err)
	}
	for _, packet := range packets {
		fmt.Printf("%v\t%v\n", packet.Name, packet.Tags)
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
	case "remove-packet":
		{
			removePacket(db)
		}
	case "list-packets":
		{
			listPackets(db)
		}
	case "list-systems":
		{
			listSystems(db)
		}
	case "get-state":
		{
			getSystemState(db)
		}
	case "get-desired-state":
		{
			getDesiredSystemState(db)
		}
	case "set-desired-state":
		{
			setDesiredSystemState(db)
		}
	default:
		{
			log.Fatal("no such command")
		}
	}
}
