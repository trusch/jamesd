package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/trusch/jamesd"
)

var dbUrl = flag.String("db", "localhost", "mongodb url")
var command = flag.String("cmd", "", "one of add-permission, del-permission, check-permission, add-package, get-package")

var user = flag.String("user", "", "affected user")
var permission = flag.String("permission", "", "permission to add/delete")

var packageName = flag.String("package-name", "", "name of the package")
var packageDataFile = flag.String("package-data", "", "compressed tar archive with package data")
var scriptFile = flag.String("script-file", "", "post-install script")
var version = flag.String("version", "", "package version")
var arch = flag.String("architecture", "", "package architecture")

func addPermission(db *jamesd.DB) {
	if *user == "" || *permission == "" {
		log.Fatal("specify --user and --permission")
	}
	err := db.AddPermission(*user, *permission)
	if err != nil {
		log.Fatal(err)
	}
}

func delPermission(db *jamesd.DB) {
	if *user == "" || *permission == "" {
		log.Fatal("specify --user and --permission")
	}
	err := db.RemovePermission(*user, *permission)
	if err != nil {
		log.Fatal(err)
	}
}

func checkPermission(db *jamesd.DB) {
	if *user == "" || *permission == "" {
		log.Fatal("specify --user and --permission")
	}
	if db.CheckPermission(*user, *permission) {
		log.Printf("%v has the permission %v", *user, *permission)
	} else {
		log.Printf("%v has NOT the permission %v", *user, *permission)
		os.Exit(1)
	}
}

func addPackage(db *jamesd.DB) {
	if *packageName == "" || *packageDataFile == "" || *version == "" || *arch == "" {
		log.Fatal("specify --package-name, --package-data, --version and --architecture")
	}
	pack := &jamesd.Package{
		Name:    *packageName,
		Arch:    *arch,
		Version: *version,
	}
	switch filepath.Ext(*packageDataFile) {
	case ".gz":
		pack.Compression = jamesd.GZIP
	case ".bzip2":
		pack.Compression = jamesd.BZIP2
	case ".bz2":
		pack.Compression = jamesd.BZIP2
	case ".lzma":
		pack.Compression = jamesd.LZMA
	case ".xz":
		pack.Compression = jamesd.LZMA
	}
	bs, err := ioutil.ReadFile(*packageDataFile)
	if err != nil {
		log.Fatal(err)
	}
	pack.Data = bs
	if *scriptFile != "" {
		bs, err = ioutil.ReadFile(*scriptFile)
		if err != nil {
			log.Fatal(err)
		}
		pack.Script = bs
	}
	err = db.AddPackage(pack)
	if err != nil {
		log.Fatal(err)
	}
}

func getPackage(db *jamesd.DB) {
	if *packageName == "" || *version == "" || *arch == "" {
		log.Fatal("specify --package-name, --version and --architecture")
	}
	pack, err := db.GetPackage(*packageName, *arch, *version)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("found package %v_%v:%v data-size: %v script-size: %v", *packageName, *arch, *version, len(pack.Data), len(pack.Script))
}

func main() {
	flag.Parse()
	db, err := jamesd.NewDB(*dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	switch *command {
	case "add-permission":
		{
			addPermission(db)
		}
	case "del-permission":
		{
			delPermission(db)
		}
	case "check-permission":
		{
			checkPermission(db)
		}
	case "add-package":
		{
			addPackage(db)
		}
	case "get-package":
		{
			getPackage(db)
		}
	default:
		{
			log.Fatal("no such command")
		}
	}
}
