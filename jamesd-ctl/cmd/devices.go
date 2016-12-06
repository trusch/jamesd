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

import "github.com/spf13/cobra"

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:     "devices",
	Aliases: []string{"systems", "machines", "device", "system", "machine"},
	Short:   "device related tasks",
	Long:    `This allows you to inspect and interact with the connected devices packet state`,
}

func init() {
	RootCmd.AddCommand(devicesCmd)
	devicesCmd.PersistentFlags().StringP("name", "n", "", "name of the device to operate on")
}