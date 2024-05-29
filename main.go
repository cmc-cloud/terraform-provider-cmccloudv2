package main

import (
	"github.com/cmc-cloud/terraform-provider-cmccloudv2/cmccloudv2"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		// -provider-name flag or set its value to the updated provider name.
		ProviderFunc: cmccloudv2.Provider,
	})
}
