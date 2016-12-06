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
	"os"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/db"
)

// getMatchingPacketsCmd represents the getMatchingPackets command
var getMatchingPacketsCmd = &cobra.Command{
	Use:   "matching",
	Short: "List packets which matches given name and tags",
	Long: `List packets which matches given name and tags.

	A packet satifies a request if all its tags are in the request tag list.
	For example the packet {name: foo, tags: [a]} matches the request {name: foo, tags: [a,b,c]}`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		if name == "" && len(args) > 0 {
			name = args[0]
		}
		tags, _ := cmd.Flags().GetStringSlice("tags")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		listMatchingPackets(db, name, tags)
	},
}

func listMatchingPackets(db *db.DB, packetName string, tags []string) {
	packets, err := db.GetMatchingPackets(packetName, tags)
	if err != nil {
		log.Fatal(err)
	}
	sort.Sort(packets)
	const padding = 3
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "Name\t Tags")
	for _, p := range packets {
		fmt.Fprintf(w, "%v\t %v\n", p.Name, p.Tags)
	}
	w.Flush()
}

func init() {
	packetsCmd.AddCommand(getMatchingPacketsCmd)
}