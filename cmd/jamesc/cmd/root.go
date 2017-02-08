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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/trusch/jamesd/cli"
	"github.com/trusch/jamesd/installer"
	"github.com/trusch/jamesd/packet"
	"github.com/trusch/jamesd/state"
)

var cfgFile string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "jamesc",
	Short: "A client for automatic software installation via jamesd",
	Long:  `This client polls jamesd periodically, asking for packages to be installed and installs them accordingly.`,
	Run: func(cmd *cobra.Command, args []string) {
		jamesdAddr := viper.GetString("addr")
		token := viper.GetString("token")
		interval := viper.GetDuration("interval")
		installRoot := viper.GetString("root")
		packetDir := viper.GetString("packets")

		labels := getLabels(cmd)
		client := cli.New(jamesdAddr)
		if token != "" {
			client.SetToken(token)
		}
		for {
			if state, err := client.GetDesiredState(labels); err == nil {
				log.Print("got new state")
				e := uninstall(packetDir, installRoot, state.Apps)
				if e != nil {
					log.Printf("ERROR in UNINSTALL: %v", e)
				}
				e = install(client, packetDir, installRoot, state.Apps)
				if e != nil {
					log.Printf("ERROR in INSTALL: %v", e)
				}
			} else {
				log.Print(err)
			}
			time.Sleep(interval)
		}
	},
}

// Execute is the main entry point
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.jamesc.yaml)")
	RootCmd.Flags().StringP("addr", "a", "http://localhost", "jamesd address")
	RootCmd.Flags().StringP("token", "t", "", "authorization token")
	RootCmd.Flags().StringP("labels", "l", "", "system labels")
	RootCmd.Flags().StringP("root", "r", "/", "install root")
	RootCmd.Flags().StringP("packets", "p", "/var/lib/jamesc/packets", "packet directory")
	RootCmd.Flags().DurationP("interval", "i", 30*time.Second, "check interval")

	viper.BindPFlag("config", RootCmd.Flags().Lookup("config"))
	viper.BindPFlag("addr", RootCmd.Flags().Lookup("addr"))
	viper.BindPFlag("token", RootCmd.Flags().Lookup("token"))
	viper.BindPFlag("root", RootCmd.Flags().Lookup("root"))
	viper.BindPFlag("packets", RootCmd.Flags().Lookup("packets"))
	viper.BindPFlag("interval", RootCmd.Flags().Lookup("interval"))

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigName(".jamesc") // name of config file (without extension)
	viper.AddConfigPath("$HOME")   // adding home directory as first search path
	viper.AutomaticEnv()           // read in environment variables that match
	if cfgFile != "" {             // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}
	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		log.Print("Using config file: ", viper.ConfigFileUsed())
	}
}

func getLabels(cmd *cobra.Command) map[string]string {
	labelStr, _ := cmd.Flags().GetString("labels")
	if labelStr == "" {
		labelStr = os.Getenv("LABELS")
		if labelStr == "" {
			return viper.GetStringMapString("labels")
		}
	}
	labelSlice := strings.Split(labelStr, ",")
	res := make(map[string]string)
	for _, labelStr := range labelSlice {
		parts := strings.Split(labelStr, "=")
		if len(parts) == 2 {
			res[parts[0]] = parts[1]
		}
	}
	return res
}

func checkIfInstalled(root, id string) bool {
	installed := false
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if info.Name() == id+".jpk" {
			installed = true
			return errors.New("") //just to break the walk
		}
		return nil
	})
	return installed
}

func collectUnneededPackets(root string, desired []*state.App) ([]string, error) {
	res := make([]string, 0, 32)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if err != nil {
			return err
		}
		needUninstall := true
		for _, app := range desired {
			if info.Name() == app.Hash+".jpk" {
				needUninstall = false
				break
			}
		}
		if needUninstall {
			res = append(res, info.Name())
		}
		return nil
	})
	return res, err
}

func uninstall(packetRoot, installRoot string, desired []*state.App) error {
	uninstallPackets, err := collectUnneededPackets(packetRoot, desired)
	if err != nil {
		return err
	}
	for _, p := range uninstallPackets {
		pack, err := packet.NewFromFile(filepath.Join(packetRoot, p))
		if err != nil {
			return err
		}
		err = installer.Uninstall(pack, installRoot)
		if err != nil {
			return err
		}
		err = os.Remove(filepath.Join(packetRoot, p))
		if err != nil {
			return err
		}
		log.Printf("uninstalled %v (%v)", pack.Name, pack.ControlInfo.Hash)
	}
	return nil
}

func install(cli *cli.Client, packetRoot, installRoot string, desired []*state.App) error {
	for _, app := range desired {
		if !checkIfInstalled(packetRoot, app.Hash) {
			pack, err := cli.GetPacketData(app.Hash)
			if err != nil {
				return err
			}
			bs, err := pack.ToData()
			if err != nil {
				return err
			}
			err = installer.Install(pack, installRoot)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(filepath.Join(packetRoot, app.Hash+".jpk"), bs, 0655)
			if err != nil {
				return err
			}
			log.Printf("installed %v (%+v)", app.Name, app.Labels)
		}
	}
	return nil
}
