package vault

import (
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/api"
)

func VaultClient(hostName string, hostPort int) (*api.Client, error) {

	// init a clean httpClient
	httpClient := cleanhttp.DefaultPooledClient()

	// format the URL with the passed host and por
	url := fmt.Sprintf("https://%s:%v", hostName, hostPort)

	// create a vault client
	client, err := api.NewClient(&api.Config{Address: url, HttpClient: httpClient})
	if err != nil {
		return nil, err
	}

	return client, nil

}
