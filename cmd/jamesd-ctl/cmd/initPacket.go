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
	"github.com/trusch/jamesd2/packet"
)

// initPacketCmd represents the initPacket command
var initPacketCmd = &cobra.Command{
	Use:   "init",
	Short: "init a new packet directory",
	Long:  `This command initializes a directory structure to create a new jamesd-packet`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			log.Fatal("specify a name")
		}
		labels := getLabels(cmd)
		if err := packet.InitDirectory(dir, name, labels); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	packetCmd.AddCommand(initPacketCmd)
	initPacketCmd.Flags().StringP("dir", "d", ".", "packet root directory")
	initPacketCmd.Flags().StringP("name", "n", "", "packet name")
	initPacketCmd.Flags().StringSliceP("labels", "l", []string{}, "comma separated list of labels: foo=bar,baz=quy...")
}
