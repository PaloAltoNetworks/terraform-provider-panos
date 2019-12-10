package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-panos/panos"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: panos.Provider,
	})
}
