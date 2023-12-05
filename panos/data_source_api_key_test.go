package panos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPanosApiKey(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccApiKeyConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.panos_api_key.test", "api_key"),
				),
			},
		},
	})
}

func testAccApiKeyConfig() string {
	return `
data "panos_api_key" "test" {}
`
}
