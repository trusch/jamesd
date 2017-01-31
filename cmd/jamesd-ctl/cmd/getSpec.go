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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/jamesd2/cli"
)

// getSpecCmd represents the getSpec command
var getSpecCmd = &cobra.Command{
	Use:   "get",
	Short: "get a single spec",
	Long:  `This gives you a single spec.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString("address")
		id, _ := cmd.Flags().GetString("id")
		client := cli.New(addr)
		spec, err := client.GetSpec(id)
		if err != nil {
			log.Fatal(err)
		}
		dumpAsYaml(spec)
	},
}

func init() {
	specCmd.AddCommand(getSpecCmd)
	getSpecCmd.Flags().String("id", "", "id of the spec")
}
