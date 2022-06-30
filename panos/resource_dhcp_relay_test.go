package panos

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccPanosDhcpRelay_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,

		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "panos_dhcp_relay" "test" {
					name = "test-eth0"
					ipv4_enabled = true
					ipv4_servers = [
						"10.1.0.1",
						"10.2.0.1",
					]
				}`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("panos_dhcp_relay.test", "name", "test-eth0"),
					resource.TestCheckResourceAttr("panos_dhcp_relay.ipv4_servers", "#", "2"),
				),
			},
		},
	})
}
