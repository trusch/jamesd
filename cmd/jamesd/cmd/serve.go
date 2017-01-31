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
	"github.com/trusch/jamesd/db"
	"github.com/trusch/jamesd/http"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the jamesd daemon",
	Long:  `start the jamesd daemon.`,
	Run: func(cmd *cobra.Command, args []string) {
		dbURI := viper.GetString("database")
		addr := viper.GetString("listen")
		log.Printf("connecting to %v...", dbURI)
		db, err := db.New(dbURI)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("successfully connected to %v.", dbURI)
		log.Printf("start listening on %v...", addr)
		err = http.ListenAndServe(db, addr)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringP("listen", "l", ":80", "REST server address")
	viper.BindPFlag("listen", serveCmd.Flags().Lookup("listen"))
}
