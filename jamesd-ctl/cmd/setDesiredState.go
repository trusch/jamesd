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
	"github.com/trusch/jamesd/systemstate"
	"gopkg.in/yaml.v2"
)

// setDesiredStateCmd represents the setDesiredState command
var setDesiredStateCmd = &cobra.Command{
	Use:   "set-desired",
	Short: "set the desired state of the specified device",
	Long:  `This set the desired state of the specified device.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbUrl, _ := cmd.Flags().GetString("db")
		file, _ := cmd.Flags().GetString("file")
		db, err := db.New(dbUrl)
		if err != nil {
			log.Fatal(err)
		}
		setDesiredSystemState(db, file)
	},
}

func setDesiredSystemState(db *db.DB, file string) {
	if file == "" {
		log.Fatal("specify --file")
	}
	bs, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	systemState := &systemstate.SystemState{}
	err = yaml.Unmarshal(bs, systemState)
	if err != nil {
		log.Fatal(err)
	}
	err = db.SaveDesiredSystemState(systemState)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	devicesCmd.AddCommand(setDesiredStateCmd)
	setDesiredStateCmd.Flags().StringP("file", "f", "", "statefile to use")
}
