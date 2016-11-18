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
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/db"
)

// listPacketsCmd represents the listPackets command
var listPacketsCmd = &cobra.Command{
	Use:   "list",
	Short: "List packets which matches the given name and tag list.",
	Long: `List packets which matches the given name and tag list.

	A packet is matched if it contains all tags in the taglist
	For example the packet {name: foo, tags: [a,b,c]} matches the request {name: foo, tags: [a]}`,

	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		listPackets(db, name, tags)
	},
}

func listPackets(db *db.DB, packetName string, tags []string) {
	packets, err := db.ListPackets(packetName, tags)
	if err != nil {
		log.Fatal(err)
	}
	for _, packet := range packets {
		fmt.Printf("%v\t%v\n", packet.Name, packet.Tags)
	}
}

func init() {
	packetsCmd.AddCommand(listPacketsCmd)
}
