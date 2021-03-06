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
	"strings"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/installer"
	"github.com/trusch/jamesd/packet"
)

// installPacketCmd represents the installPacket command
var installPacketCmd = &cobra.Command{
	Use:   "install",
	Short: "install a packet",
	Long:  `This installs a packet on your system.`,
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := cmd.Flags().GetString("file")
		if file == "" {
			if len(args) > 0 && strings.HasSuffix(args[0], ".jpk") {
				file = args[0]
			}
		}
		var pack *packet.Packet
		if file == "" {
			p, err := getPacketByID(cmd, args)
			if err != nil {
				log.Fatal(err)
			}
			pack = p
		} else {
			bs, err := ioutil.ReadFile(file)
			if err != nil {
				log.Fatal(err)
			}
			p, err := packet.NewFromData(bs)
			if err != nil {
				log.Fatal(err)
			}
			pack = p
		}
		pack.Hash()
		root, _ := cmd.Flags().GetString("root")
		if err := installer.Install(pack, root); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	packetCmd.AddCommand(installPacketCmd)
	installPacketCmd.Flags().StringP("file", "f", "", "packet filename")
	installPacketCmd.Flags().StringP("root", "r", "/", "install root")
}
