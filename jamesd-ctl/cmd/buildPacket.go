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
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/tatar"
)

// buildPacketCmd represents the buildPacket command
var buildPacketCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a packet from a specified directory",
	Long:  `This builds a packet from a specified directory suitable to be imported into jamesd.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			if len(args) == 0 {
				dir = "."
			} else {
				dir = args[0]
			}
		}
		file, _ := cmd.Flags().GetString("file")
		build(dir, file)
	},
}

func build(dir, file string) {
	data, err := tatar.NewFromDirectory(filepath.Join(dir, "data"))
	if err != nil {
		log.Fatal(err)
	}
	data.Compression = tatar.LZMA

	preInst, err := ioutil.ReadFile(filepath.Join(dir, "preinst"))
	if err != nil {
		log.Fatal(err)
	}

	postInst, err := ioutil.ReadFile(filepath.Join(dir, "postinst"))
	if err != nil {
		log.Fatal(err)
	}

	preRm, err := ioutil.ReadFile(filepath.Join(dir, "prerm"))
	if err != nil {
		log.Fatal(err)
	}

	postRm, err := ioutil.ReadFile(filepath.Join(dir, "postrm"))
	if err != nil {
		log.Fatal(err)
	}

	control, err := ioutil.ReadFile(filepath.Join(dir, "control"))
	if err != nil {
		log.Fatal(err)
	}

	info := &packet.ControlInfo{}
	err = info.FromYaml(control)
	if err != nil {
		log.Fatal(err)
	}

	pack := packet.Packet{
		ControlInfo: packet.ControlInfo{
			Name:    info.Name,
			Tags:    info.Tags,
			Depends: info.Depends,
			Scripts: packet.Scripts{
				PreInst:  string(preInst),
				PostInst: string(postInst),
				PreRm:    string(preRm),
				PostRm:   string(postRm),
			},
		},
		Data: data,
	}
	resultData, err := pack.ToData()
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(file, resultData, 0755)
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	packetsCmd.AddCommand(buildPacketCmd)
	buildPacketCmd.Flags().StringP("dir", "d", "", "packet directory")
	buildPacketCmd.Flags().StringP("file", "f", "packet.jpk", "output packet file")
}
