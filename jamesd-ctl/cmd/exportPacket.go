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

// exportPacketCmd represents the importPacket command
var exportPacketCmd = &cobra.Command{
	Use:   "export",
	Short: "export lets you export .jpk files.",
	Long: `export lets you export .jpk files.

	These files could then be transfered to another jamesd repo, or used for backup and stuff.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		file, _ := cmd.Flags().GetString("file")
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		pack := getPacket(db, name, tags)
		bs, err := pack.ToData()
		if err != nil {
			log.Fatal(err)
		}
		if file == "" {
			log.Fatal("specify --file")
		}
		err = ioutil.WriteFile(file, bs, 0755)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func getPacket(db *db.DB, name string, tags []string) *packet.Packet {
	if name == "" {
		log.Fatal("specify --name")
	}
	pack, err := db.GetPacket(name, tags)
	if err != nil {
		log.Fatal(err)
	}
	return pack
}

func init() {
	packetsCmd.AddCommand(exportPacketCmd)
	exportPacketCmd.Flags().StringP("file", "f", "", "path to the .jpk file")
}
