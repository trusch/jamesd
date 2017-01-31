// Copyright © 2017 Tino Rusch
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/jamesd/cli"
	"github.com/trusch/jamesd/spec"
)

// updateSpecCmd represents the updateSpec command
var updateSpecCmd = &cobra.Command{
	Use:   "update",
	Short: "update a spec",
	Long:  `This updates a spec in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString("address")
		client := cli.New(addr)
		file, _ := cmd.Flags().GetString("file")
		if file == "" && len(args) > 0 {
			file = args[0]
		}
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		var s spec.Spec
		if err := yaml.Unmarshal(bs, &s); err != nil {
			log.Fatal(err)
		}
		if err := client.PutSpec(&s); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	specCmd.AddCommand(updateSpecCmd)
	updateSpecCmd.Flags().StringP("file", "f", "/dev/stdin", "spec file")
}
