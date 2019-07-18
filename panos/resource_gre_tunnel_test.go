package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/tunnel/gre"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosGreTunnel_basic(t *testing.T) {
	versionAdded := version.Number{
		Major: 9,
		Minor: 0,
	}

	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccPanosVersion.Gte(versionAdded) {
		t.Skip("GRE tunnels are available in PAN-OS 9.0+")
	}

	var o gre.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosGreTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGreTunnelConfig(name, "5.5.5.5", 0, 42, 3, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosGreTunnelExists("panos_gre_tunnel.test", &o),
					testAccCheckPanosGreTunnelAttributes(&o, name, "5.5.5.5", 0, 42, 3, true, false),
				),
			},
			{
				Config: testAccGreTunnelConfig(name, "6.6.6.6", 1, 43, 4, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosGreTunnelExists("panos_gre_tunnel.test", &o),
					testAccCheckPanosGreTunnelAttributes(&o, name, "6.6.6.6", 1, 43, 4, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosGreTunnelExists(n string, o *gre.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		name := rs.Primary.ID
		v, err := fw.Network.GreTunnel.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosGreTunnelAttributes(o *gre.Entry, name, pa string, lav, ttl, kai int, tos, dis bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.PeerAddress != pa {
			return fmt.Errorf("Peer address is %s, expected %s", o.PeerAddress, pa)
		}

		addy := fmt.Sprintf("10.5.%d.1/24", lav)
		if o.LocalAddressValue != addy {
			return fmt.Errorf("Local address value is %s, expected %s", o.LocalAddressValue, addy)
		}

		if o.Ttl != ttl {
			return fmt.Errorf("TTL is %d, expected %d", o.Ttl, ttl)
		}

		if o.KeepAliveInterval != kai {
			return fmt.Errorf("Keep alive interval is %d, expected %d", o.KeepAliveInterval, kai)
		}

		if o.CopyTos != tos {
			return fmt.Errorf("Copy TOS is %t, expected %t", o.CopyTos, tos)
		}

		if o.Disabled != dis {
			return fmt.Errorf("Disabled is %t, expected %t", o.Disabled, dis)
		}

		return nil
	}
}

func testAccPanosGreTunnelDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_gre_tunnel" {
			continue
		}

		if rs.Primary.ID != "" {
			_, err := fw.Network.GreTunnel.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccGreTunnelConfig(name, pa string, lav, ttl, kai int, tos, dis bool) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "x" {
    name = "ethernet1/2"
    mode = "layer3"
    vsys = "vsys1"
    static_ips = ["10.5.0.1/24", "10.5.1.1/24"]
}

resource "panos_tunnel_interface" "x" {
    name = "tunnel.7"
    vsys = "vsys1"
}

resource "panos_gre_tunnel" "test" {
    name = %q
    interface = panos_ethernet_interface.x.name
    tunnel_interface = panos_tunnel_interface.x.name
    enable_keep_alive = true
    peer_address = %q
    local_address_value = panos_ethernet_interface.x.static_ips.%d
    ttl = %d
    keep_alive_interval = %d
    copy_tos = %t
    disabled = %t
}
`, name, pa, lav, ttl, kai, tos, dis)
}
