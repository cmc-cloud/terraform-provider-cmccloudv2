package cmccloudv2

import (
	"log"

	gocmcapi "github.com/cmc-cloud/gocmcapiv2"
)

// Config object
type Config struct {
	APIKey      string
	APIEndpoint string
	ProjectId   string
	RegionId    string
}

// CombinedConfig struct
type CombinedConfig struct {
	client *gocmcapi.Client
}

func (c *CombinedConfig) goCMCClient() *gocmcapi.Client { return c.client }

// Client config
func (c *Config) Client() (*CombinedConfig, error) {
	myStruct := gocmcapi.ClientConfigs{
		APIKey:      c.APIKey,
		APIEndpoint: c.APIEndpoint,
		ProjectId:   c.ProjectId,
		RegionId:    c.RegionId,
	}

	client, err := gocmcapi.NewClient(myStruct)
	if err != nil {
		log.Fatal(err)
	}
	return &CombinedConfig{
		client: client,
	}, nil
}
