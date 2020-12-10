package panos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPanosPlugin(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPluginConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttrSet("data.panos_plugin.test", "installed"),
					resource.TestCheckResourceAttrSet("data.panos_plugin.test", "total"),
				),
			},
		},
	})
}

func testAccPluginConfig() string {
	return `
data "panos_plugin" "test" {}
`
}
