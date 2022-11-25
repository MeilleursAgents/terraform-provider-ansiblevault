package main

import (
	"github.com/MeilleursAgents/terraform-provider-ansiblevault/v2/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.Provider,
	})
}
