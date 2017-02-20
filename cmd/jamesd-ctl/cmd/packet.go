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
	"github.com/trusch/jamesd/cli"
	"github.com/trusch/jamesd/packet"
)

// packetCmd represents the packet command
var packetCmd = &cobra.Command{
	Use:   "packet",
	Short: "packet related commands",
}

func init() {
	RootCmd.AddCommand(packetCmd)
	packetCmd.PersistentFlags().String("id", "", "id of the packet")
}

func getPacketByID(cmd *cobra.Command, args []string) (*packet.Packet, error) {
	addr := viper.GetString("address")
	id, _ := cmd.Flags().GetString("id")
	if addr == "" {
		log.Fatal("specify address!")
	}
	if id == "" {
		if len(args) > 0 {
			id = args[0]
		} else {
			log.Fatal("specify an id")
		}
	}
	client := cli.New(addr)
	token := viper.GetString("token")
	if token != "" {
		client.SetToken(token)
	}
	return client.GetPacketData(id)
}
