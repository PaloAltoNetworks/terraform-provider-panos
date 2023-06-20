package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/interface/tunnel"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosTunnelInterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o tunnel.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("tunnel.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosTunnelInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelInterfaceConfig(name, "first comment", "10.8.9.1/24", 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTunnelInterfaceExists("panos_tunnel_interface.test", &o),
					testAccCheckPanosTunnelInterfaceAttributes(&o, name, "first comment", "10.8.9.1/24", 600),
				),
			},
			{
				Config: testAccTunnelInterfaceConfig(name, "second comment", "10.9.10.1/24", 700),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTunnelInterfaceExists("panos_tunnel_interface.test", &o),
					testAccCheckPanosTunnelInterfaceAttributes(&o, name, "second comment", "10.9.10.1/24", 700),
				),
			},
		},
	})
}

func testAccCheckPanosTunnelInterfaceExists(n string, o *tunnel.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, name := parseTunnelInterfaceId(rs.Primary.ID)
		v, err := fw.Network.TunnelInterface.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosTunnelInterfaceAttributes(o *tunnel.Entry, name, cmt, ip string, mtu int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Comment != cmt {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, cmt)
		}

		if len(o.StaticIps) != 1 || o.StaticIps[0] != ip {
			return fmt.Errorf("Static IPs is %#v, expected [%q]", o.StaticIps, ip)
		}

		if o.Mtu != mtu {
			return fmt.Errorf("MTU is %d, expected %d", o.Mtu, mtu)
		}

		return nil
	}
}

func testAccPanosTunnelInterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_tunnel_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			_, name := parseTunnelInterfaceId(rs.Primary.ID)
			_, err := fw.Network.TunnelInterface.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccTunnelInterfaceConfig(name, cmt, ip string, mtu int) string {
	return fmt.Sprintf(`
resource "panos_tunnel_interface" "test" {
    name = %q
    comment = %q
    static_ips = [%q]
    mtu = %d
}
`, name, cmt, ip, mtu)
}
