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

// listDevicesCmd represents the listDevices command
var listDevicesCmd = &cobra.Command{
	Use:   "list",
	Short: "list all known devices",
	Long: `List all known devices.

	Devices which successfully authenticated agains jamesd should be listed here.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		listDevices(db)
	},
}

func listDevices(db *db.DB) {
	systems, err := db.ListSystems()
	if err != nil {
		log.Fatal(err)
	}
	for _, system := range systems {
		fmt.Println(system)
	}
}

func init() {
	devicesCmd.AddCommand(listDevicesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listDevicesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listDevicesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
