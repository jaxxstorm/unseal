// Copyright © 2017 Lee Briggs <lee@leebriggs.co.uk>
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
	"fmt"
	"os"
	"sync"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
	v "github.com/jaxxstorm/unseal/vault"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var unsealKey string
var vaultHost string
var vaultPort int

type Host struct {
	Name string
	Port int
	Key  string
}

var hosts []Host

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "unseal",
	Short: "Unseal a set of vault servers",
	Long:  `Unseal allows you to unseal a large set of vault servers using the HTTP API.`,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {

		// unmarshal config file
		err := viper.UnmarshalKey("hosts", &hosts)

		// check for valid config file
		if err != nil {
			panic("Unable to unmarshal hosts. Is your config file valid??")
		}

		var wg sync.WaitGroup

		for _, h := range hosts {

			hostName := h.Name
			hostPort := h.Port
			key := h.Key

			wg.Add(1)

			go func(hostName string, hostPort int) {
				defer wg.Done()
				httpClient := cleanhttp.DefaultPooledClient()

				// format the URL with the passed host and por
				url := fmt.Sprintf("https://%s:%v", hostName, hostPort)
				// create a vault client
				client, err := api.NewClient(&api.Config{Address: url, HttpClient: httpClient})
				if err != nil {
					log.WithFields(log.Fields{"host": hostName, "port": hostPort}).Error(err)
				}
				// get the current status
				init := v.InitStatus(client)
				if init.Ready == true {
					result, err := client.Sys().Unseal(key)
					// should we keep going here? Don't panic?
					if err != nil {
						log.WithFields(log.Fields{"host": hostName}).Error("Error running unseal operation")
					}

					// if it's still sealed, print the progress
					if result.Sealed == true {
						log.WithFields(log.Fields{"host": hostName, "progress": result.Progress, "threshold": result.T}).Info("Unseal operation performed")
						// otherwise, tell us it's unsealed!
					} else {
						log.WithFields(log.Fields{"host": hostName, "progress": result.Progress, "threshold": result.T}).Info("Vault is unsealed!")
					}
					// zero out the key
					// FIXME: is this the best way to do this?
					// Is it safe?
					key = ""
				} else {
					// sad times, not ready to be unsealed
					log.WithFields(log.Fields{"host": hostName}).Error("Vault is not ready to be unsealed")
				}

			}(hostName, hostPort)

		}
		wg.Wait()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// define flags
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.unseal/config.yaml)")
	//RootCmd.PersistentFlags().StringVar(&unsealKey, "key", "", "unseal key (will prompt if not provided)")
	//RootCmd.PersistentFlags().StringVar(&vaultHost, "host", "", "vault host to unseal")
	//RootCmd.PersistentFlags().IntVar(&vaultPort, "port", 8200, "port vault is listening on")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config") // name of config file (without extension)
		viper.AddConfigPath("/etc/unseal")
		viper.AddConfigPath("$HOME/.unseal") // adding home directory as first search path
		viper.AddConfigPath(".")
		viper.AutomaticEnv() // read in environment variables that match
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
