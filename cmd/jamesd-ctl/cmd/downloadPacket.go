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
)

// downloadPacketCmd represents the downloadPacket command
var downloadPacketCmd = &cobra.Command{
	Use:   "download",
	Short: "download a packet",
	Long:  `this downloads a specific packet from the db and saves it to file.`,
	Run: func(cmd *cobra.Command, args []string) {
		pack, err := getPacketByID(cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			file = pack.Name
			if version, ok := pack.Labels["version"]; ok {
				file += "_" + version
			}
			file += ".jpk"
		}
		packData, err := pack.ToData()
		if err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(file, packData, 0655); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	packetCmd.AddCommand(downloadPacketCmd)
	downloadPacketCmd.Flags().StringP("file", "f", "", "output file")
}
