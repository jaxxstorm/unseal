package main

import (
	"fmt"
	"github.com/bgentry/speakeasy"
	vault "github.com/hashicorp/vault/api"
	v "github.com/jaxxstorm/unseal/vault"
	"github.com/urfave/cli"
	"os"
)

func main() {

	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "host, H", Usage: "The vault host to unseal"},
		cli.IntFlag{Name: "port, P", Usage: "Vault Port", Value: 8200},
		cli.StringFlag{Name: "key, K", Usage: "Vault key to use to unseal (will prompt if not provided)", EnvVar: "VAULT_KEY"},
	}

	app.Name = "unseal"
	app.Version = "0.1"
	app.Usage = "Safely unseal a vault server"
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Lee Briggs",
		},
	}

	app.Action = func(c *cli.Context) error {

		if !c.IsSet("host") {
			cli.ShowAppHelp(c)
			return cli.NewExitError("Error: Please specify a host to unseal", -1)
		}

		var unsealkey string
		var err error

		if !c.IsSet("key") {
			unsealkey, err = speakeasy.Ask("Please specify your unseal key: ")
			if err != nil {
				return cli.NewExitError("Error: Please specify an unseal key", -1)
			}
		} else {
			unsealkey = c.String("key")
		}

		// init client
		url := fmt.Sprintf("https://%s:%v", c.String("host"), c.Int("port"))
		client, err := vault.NewClient(&vault.Config{Address: url})

		// determine if our vault server is ready for us
		ready := v.Ready(client)

		// if we are ready, unseal
		if ready.Ready == true {

			unseal, err := client.Sys().Unseal(unsealkey)
			if err != nil {
				return cli.NewExitError("Error: Error unsealing vault", -1)
			}

			fmt.Println(fmt.Sprintf("Unseal operation complete. Required: %v Progress: %v", unseal.T, unseal.Progress))

		} else {
			// else, die
			return cli.NewExitError(fmt.Sprintf("Error: %s", ready.Reason), 2)
		}
		return nil
	}
	app.Run(os.Args)
}
