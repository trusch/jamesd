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

// listSpecsCmd represents the listSpecs command
var listSpecsCmd = &cobra.Command{
	Use:   "list",
	Short: "list your specs",
	Long:  `list dumps specs matching your request`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		tags, _ := cmd.Flags().GetStringSlice("tags")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		specs, err := db.GetSpecs(name, tags)
		if err != nil {
			log.Fatal(err)
		}
		bs, err := yaml.Marshal(specs)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bs))
	},
}

func init() {
	specsCmd.AddCommand(listSpecsCmd)
	listSpecsCmd.Flags().StringP("name", "n", "", "spec target name")
	listSpecsCmd.Flags().StringSliceP("tags", "t", []string{}, "spec target tag list")

}
