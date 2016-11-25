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

// getSpecCmd represents the getSpec command
var getSpecCmd = &cobra.Command{
	Use:   "get",
	Short: "get a single spec",
	Long:  `Get returns a spec identified by name`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		if name == "" && len(args) > 0 {
			name = args[0]
		}
		getSpec(db, name)
	},
}

func getSpec(db *db.DB, name string) {
	if name == "" {
		log.Fatal("specify --name")
	}
	spec, err := db.GetSpec(name)
	if err != nil {
		log.Fatal(err)
	}
	d, err := yaml.Marshal(spec)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(d))
}

func init() {
	specsCmd.AddCommand(getSpecCmd)

	getSpecCmd.Flags().StringP("name", "n", "", "spec name")
}
