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
	"github.com/trusch/jamesd/cli"
	"github.com/trusch/jamesd/packet"
)

// uploadPacketCmd represents the uploadPacket command
var uploadPacketCmd = &cobra.Command{
	Use:   "upload",
	Short: "upload a packet",
	Long:  `This uploads a packet to the database.`,
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
		pack, err := packet.NewFromData(bs)
		if err != nil {
			log.Fatal(err)
		}
		if err := client.UploadPacket(pack); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	packetCmd.AddCommand(uploadPacketCmd)
	uploadPacketCmd.Flags().StringP("file", "f", "", "packet file")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// uploadPacketCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// uploadPacketCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
