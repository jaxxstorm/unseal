package cmd

import (
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	v "github.com/jaxxstorm/unseal/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Print the vault status for all configured vaults",
	Long:  `Print the current status for all configured vaults, without passing a key`,
	Run: func(cmd *cobra.Command, args []string) {

		err := viper.UnmarshalKey("hosts", &hosts)

		caPath = viper.GetString("capath")

		// check for valid config file
		if err != nil {
			log.Fatal("Unable to read hosts key in config file: %s", err)
		}

		if os.Getenv("VAULT_ADDR") != "" {
			log.Warning("Warning VAULT_ADDR environment variable is set. This will override the hostname in your config file, it's probably not what you want")
		}

		var wg sync.WaitGroup

		// loop through the hosts
		for _, h := range hosts {

			// set hostnames for waitgroup
			hostName := h.Name
			hostPort := h.Port

			wg.Add(1)

			go func(hostName string, hostPort int) {
				defer wg.Done()

				// create a vault client
				client, err := v.VaultClient(hostName, hostPort, caPath)

				// issue creating vault client for this host
				if err != nil {
					log.WithFields(log.Fields{"host": hostName}).Error("Error creating vault client: ", err)
				}

				// get the seal status
				result, err := client.Sys().SealStatus()

				if err != nil {
					log.WithFields(log.Fields{"host": hostName}).Error("Error getting seal status: ", err)
				} else {
					// only check the seal status if we have a client
					if result.Sealed == true {
						log.WithFields(log.Fields{"host": hostName, "progress": result.Progress, "threshold": result.T}).Error("Vault is sealed!")
					} else {
						log.WithFields(log.Fields{"host": hostName, "progress": result.Progress, "threshold": result.T}).Info("Vault is unsealed!")
					}
				}

			}(hostName, hostPort)
		}
		wg.Wait()

	},
}
