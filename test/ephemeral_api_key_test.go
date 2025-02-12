package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralApiKey(t *testing.T) {
	t.Parallel()

	resource.Test(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},

		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ephemeralApiKeyTmpl,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"echo.test_api_key",
						tfjsonpath.New("data").
							AtMapKey("api_key"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const ephemeralApiKeyTmpl = `
ephemeral "panos_api_key" "apikey" {
  username = "api-admin"
  password = "test-password"
}

provider "echo" {
  data = ephemeral.panos_api_key.apikey
}

resource "echo" "test_api_key" {}
`
