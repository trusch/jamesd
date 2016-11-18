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
	"log"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/db"
)

// deletePacketCmd represents the deletePacket command
var deletePacketCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a specific packet from the repository.",
	Long: `Delete a specific packet from the repository.

	Use with caution, this is destructive by nature ;)`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		removePacket(db, name, tags)
	},
}

func removePacket(db *db.DB, name string, tags []string) {
	if name == "" {
		log.Fatal("specify --name")
	}
	err := db.RemovePacket(name, tags)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	packetsCmd.AddCommand(deletePacketCmd)

}
