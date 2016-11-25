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
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/trusch/jamesd/packet"
)

// initPacketCmd represents the initPacket command
var initPacketCmd = &cobra.Command{
	Use:   "init",
	Short: "Inititalize a directory to build a packet",
	Long: `This command initializes a directory to build a packet.

	The specified directory contains a (empty) data directory, a optionally prepoulated control file and empty scripts for preinst, postinst, prerm and postrm`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		if name == "" && len(args) > 0 {
			name = args[0]
		}
		dir, _ := cmd.Flags().GetString("dir")
		if dir == "" {
			dir = name
		}
		tags, _ := cmd.Flags().GetStringSlice("tags")
		initEmptyPacket(dir, name, tags)
	},
}

func initEmptyPacket(dir, name string, tags []string) {
	if err := checkDir(dir); err != nil {
		log.Fatal(err)
	}
	os.MkdirAll(dir, 0755)
	os.MkdirAll(filepath.Join(dir, "data"), 0755)
	ioutil.WriteFile(filepath.Join(dir, "preinst"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(dir, "postinst"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(dir, "prerm"), []byte{}, 0755)
	ioutil.WriteFile(filepath.Join(dir, "postrm"), []byte{}, 0755)
	ctrl := packet.ControlInfo{
		Name: name,
		Tags: tags,
	}
	ioutil.WriteFile(filepath.Join(dir, "control"), ctrl.ToYaml(), 0755)
}

func checkDir(dir string) error {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil // no stat -> no dir -> safe to write
	}
	if !stat.IsDir() {
		return errors.New("target directory already exists (and is a file oO)")
	}
	f, err := os.Open(dir)
	if err != nil {
		return errors.New("target directory is not writeable")
	}
	defer f.Close()
	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return nil // dir exists, but no content
	}
	return errors.New("target directory already exists and is not empty")
}

func init() {
	packetsCmd.AddCommand(initPacketCmd)
	initPacketCmd.Flags().StringP("dir", "d", "", "packet directory")
}
