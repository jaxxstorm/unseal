package vault

import (
	"github.com/hashicorp/vault/api"
)

func InitStatus(client *api.Client) Status {

	// statuses
	init, err := client.Sys().InitStatus()

	if err != nil {
		panic(err)
	}
	if init == false {
		return Status{
			Ready:  false,
			Reason: "Vault is not initialized",
		}
	}

	seal, err := client.Sys().SealStatus()

	if err != nil {
		panic(err)
	}

	if seal.Sealed != true {
		return Status{
			Ready:  false,
			Reason: "Vault is already unsealed",
		}
	}

	return Status{
		Ready: true,
	}

}
