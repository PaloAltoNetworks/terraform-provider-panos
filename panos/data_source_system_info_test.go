package panos

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccPanosSystemInfo_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSystemConfig(),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.panos_system_info.test", "version_major"),
					resource.TestCheckResourceAttrSet("data.panos_system_info.test", "version_minor"),
					resource.TestCheckResourceAttrSet("data.panos_system_info.test", "version_patch"),
				),
			},
		},
	})
}

func testAccSystemConfig() string {
	return `
data "panos_system_info" "test" {}
`
}
