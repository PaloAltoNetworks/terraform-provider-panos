package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/tunnel/gre"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaGreTunnel_basic(t *testing.T) {
	versionAdded := version.Number{
		Major: 9,
		Minor: 0,
	}

	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if !testAccPanosVersion.Gte(versionAdded) {
		t.Skip("GRE tunnels are available in PAN-OS 9.0+")
	}

	var o gre.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaGreTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaGreTunnelConfig(tmpl, name, "5.5.5.5", 0, 42, 3, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGreTunnelExists("panos_panorama_gre_tunnel.test", &o),
					testAccCheckPanosPanoramaGreTunnelAttributes(&o, name, "5.5.5.5", 0, 42, 3, true, false),
				),
			},
			{
				Config: testAccPanoramaGreTunnelConfig(tmpl, name, "6.6.6.6", 1, 43, 4, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaGreTunnelExists("panos_panorama_gre_tunnel.test", &o),
					testAccCheckPanosPanoramaGreTunnelAttributes(&o, name, "6.6.6.6", 1, 43, 4, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaGreTunnelExists(n string, o *gre.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaGreTunnelId(rs.Primary.ID)
		v, err := pano.Network.GreTunnel.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaGreTunnelAttributes(o *gre.Entry, name, pa string, lav, ttl, kai int, tos, dis bool) resource.TestCheckFunc {
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

func testAccPanosPanoramaGreTunnelDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_gre_tunnel" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaGreTunnelId(rs.Primary.ID)
			_, err := pano.Network.GreTunnel.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaGreTunnelConfig(tmpl, name, pa string, lav, ttl, kai int, tos, dis bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/2"
    mode = "layer3"
    vsys = "vsys1"
    static_ips = ["10.5.0.1/24", "10.5.1.1/24"]
}

resource "panos_panorama_tunnel_interface" "x" {
    template = panos_panorama_template.x.name
    name = "tunnel.7"
    vsys = "vsys1"
}

resource "panos_panorama_gre_tunnel" "test" {
    template = panos_panorama_template.x.name
    name = %q
    interface = panos_panorama_ethernet_interface.x.name
    tunnel_interface = panos_panorama_tunnel_interface.x.name
    enable_keep_alive = true
    peer_address = %q
    local_address_value = panos_panorama_ethernet_interface.x.static_ips.%d
    ttl = %d
    keep_alive_interval = %d
    copy_tos = %t
    disabled = %t
}
`, tmpl, name, pa, lav, ttl, kai, tos, dis)
}
