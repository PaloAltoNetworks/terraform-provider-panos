package main

import (
	"github.com/PaloAltoNetworks/terraform-provider-panos/panos"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: panos.Provider,
	})
}
