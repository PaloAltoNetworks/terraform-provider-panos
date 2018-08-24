package panos

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccPanosDhcpInterfaceInfo(t *testing.T) {
	// This acctest requires that an interface already be configured as DHCP,
	// as this requires a commit and Terraform does not yet support commits.
	ev := "PANOS_DHCP_INTERFACE"
	di := os.Getenv(ev)

	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if di == "" {
		t.Skip("Env %s must be specified to run this acc test", ev)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDhcpInterfaceInfoConfig(di),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.panos_dhcp_interface_info.test", "ip"),
					resource.TestCheckResourceAttrSet("data.panos_dhcp_interface_info.test", "gateway"),
					resource.TestCheckResourceAttrSet("data.panos_dhcp_interface_info.test", "server"),
					resource.TestCheckResourceAttrSet("data.panos_dhcp_interface_info.test", "primary_dns"),
				),
			},
		},
	})
}

func testAccDhcpInterfaceInfoConfig(di string) string {
	return fmt.Sprintf(`
data "panos_dhcp_interface_info" "test" {
    interface = %q
}`, di)
}
