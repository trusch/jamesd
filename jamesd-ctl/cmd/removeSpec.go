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

// removeSpecCmd represents the removeSpec command
var removeSpecCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove a spec",
	Long:  `remove a spec by name`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		name, _ := cmd.Flags().GetString("name")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		err = db.RemoveSpec(name)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	specsCmd.AddCommand(removeSpecCmd)
	removeSpecCmd.Flags().StringP("name", "n", "", "spec name")

}
