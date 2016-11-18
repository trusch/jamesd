package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/trusch/jamesd/db"
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/systemstate"
)

var dbUrl = flag.String("db", "localhost", "mongodb url")
var command = flag.String("cmd", "", "one of add-packet, get-packet, remove-packet, list-packets, list-systems, get-state, get-desired-state, set-desired-state")
var packetFile = flag.String("packet", "", "packet file to upload")
var id = flag.String("id", "", "system id")
var file = flag.String("file", "", "system state file")
var packetName = flag.String("name", "", "packet name")
var tagList = flag.String("tags", "", "tags of packet")

var tags []string

func init() {
	flag.Parse()
	if *tagList != "" {
		tags = strings.Split(*tagList, ",")
	}
}

func addPacket(db *db.DB) {
	if *packetFile == "" {
		log.Fatal("specify --packet")
	}
	bs, err := ioutil.ReadFile(*packetFile)
	if err != nil {
		log.Fatal(err)
	}
	pack, err := packet.NewFromData(bs)
	if err != nil {
		log.Fatal(err)
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
	fmt.Printf("name: %v tags: %v\n", pack.Name, pack.Tags)
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
