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

	"gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/db"
)

// currentStateCmd represents the currentState command
var currentStateCmd = &cobra.Command{
	Use:   "get-state",
	Short: "returns the current state of the specified device",
	Long:  `This returns the current state if the specified device.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		getSystemState(db, name)
	},
}

func getSystemState(db *db.DB, name string) {
	if name == "" {
		log.Fatal("specify --name")
	}
	state, err := db.GetCurrentSystemState(name)
	if err != nil {
		log.Fatal(err)
	}
	d, err := yaml.Marshal(state)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(d))
}

func init() {
	devicesCmd.AddCommand(currentStateCmd)
}
