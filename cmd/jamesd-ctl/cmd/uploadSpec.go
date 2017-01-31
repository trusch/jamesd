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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/jamesd2/cli"
	"github.com/trusch/jamesd2/spec"
	yaml "gopkg.in/yaml.v2"
)

// uploadSpecCmd represents the uploadSpec command
var uploadSpecCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload a spec",
	Long:  `This uploads a new spec to the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := viper.GetString("address")
		client := cli.New(addr)
		file, _ := cmd.Flags().GetString("file")
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		var s spec.Spec
		if err := yaml.Unmarshal(bs, &s); err != nil {
			log.Fatal(err)
		}
		if err := client.UploadSpec(&s); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	specCmd.AddCommand(uploadSpecCmd)
	uploadSpecCmd.Flags().StringP("file", "f", "/dev/stdin", "spec file")
}