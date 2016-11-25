// Copyright © 2016 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/db"
	"github.com/trusch/jamesd/packet"
)

// importPacketCmd represents the importPacket command
var importPacketCmd = &cobra.Command{
	Use:   "import",
	Short: "import lets you import .jpk files.",
	Long: `import lets you import .jpk files.

	These files could be created with jamesd-pkg or handcrafted. They will be imported in the jamesd repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		file, _ := cmd.Flags().GetString("file")
		if file == "" && len(args) > 0 {
			file = args[0]
		}
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		addPacket(db, file)
	},
}

func addPacket(db *db.DB, filename string) {
	if filename == "" {
		log.Fatal("specify --file")
	}
	bs, err := ioutil.ReadFile(filename)
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

func init() {
	packetsCmd.AddCommand(importPacketCmd)
	importPacketCmd.Flags().StringP("file", "f", "", "path to the .jpk file")
}
