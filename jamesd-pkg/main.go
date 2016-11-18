package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/trusch/jamesd/packet"
	"github.com/trusch/tatar"
)

var cmd = flag.String("cmd", "", "one of init, info or build")
var dir = flag.String("dir", ".", "directory to operate in")
var name = flag.String("name", "", "name of the packet")
var tagList = flag.String("tags", "", "comma separated list of tags")
var file = flag.String("file", "packet.jpk", "packet file name")

var tags []string

func init() {
	flag.Parse()
	if *tagList != "" {
		tags = strings.Split(*tagList, ",")
	}
}

func initEmptyDir() {
	os.MkdirAll(*dir, 0755)
	os.MkdirAll(filepath.Join(*dir, "data"), 0755)
	ioutil.WriteFile(filepath.Join(*dir, "preinst"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(*dir, "postinst"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(*dir, "prerm"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(*dir, "postrm"), []byte{}, 0755)
	ctrl := packet.ControlInfo{
		Name: *name,
		Tags: tags,
	}
	ioutil.WriteFile(filepath.Join(*dir, "control"), ctrl.ToYaml(), 0755)
}

func build() {
	data, err := tatar.NewFromDirectory(filepath.Join(*dir, "data"))
	if err != nil {
		log.Fatal(err)
	}
	data.Compression = tatar.LZMA

	preInst, err := ioutil.ReadFile(filepath.Join(*dir, "preinst"))
	if err != nil {
		log.Fatal(err)
	}

	postInst, err := ioutil.ReadFile(filepath.Join(*dir, "postinst"))
	if err != nil {
		log.Fatal(err)
	}

	preRm, err := ioutil.ReadFile(filepath.Join(*dir, "prerm"))
	if err != nil {
		log.Fatal(err)
	}

	postRm, err := ioutil.ReadFile(filepath.Join(*dir, "postrm"))
	if err != nil {
		log.Fatal(err)
	}

	control, err := ioutil.ReadFile(filepath.Join(*dir, "control"))
	if err != nil {
		log.Fatal(err)
	}

	info := &packet.ControlInfo{}
	err = info.FromYaml(control)
	if err != nil {
		log.Fatal(err)
	}

	pack := packet.Packet{
		ControlInfo: packet.ControlInfo{
			Name:    info.Name,
			Tags:    info.Tags,
			Depends: info.Depends,
			Scripts: packet.Scripts{
				PreInst:  string(preInst),
				PostInst: string(postInst),
				PreRm:    string(preRm),
				PostRm:   string(postRm),
			},
		},
		Data: data,
	}
	resultData, err := pack.ToData()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*file, resultData, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func info() {
	bs, _ := ioutil.ReadFile(*file)
	pack := &packet.Packet{}
	err := pack.FromData(bs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(pack.ControlInfo.ToYaml()))
}

func main() {
	flag.Parse()
	switch *cmd {
	case "init":
		{
			initEmptyDir()
		}
	case "build":
		{
			build()
		}
	case "info":
		{
			info()
		}
	}
}
