package panos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPanosPanoramaPlugin(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaPluginConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttrSet("data.panos_panorama_plugin.test", "installed"),
					resource.TestCheckResourceAttrSet("data.panos_panorama_plugin.test", "total"),
				),
			},
		},
	})
}

func testAccPanoramaPluginConfig() string {
	return `
data "panos_panorama_plugin" "test" {}
`
}
