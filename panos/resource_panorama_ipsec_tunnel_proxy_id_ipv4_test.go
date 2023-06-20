package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/ipsectunnel"
	"github.com/fpluchorg/pango/netw/ipsectunnel/proxyid/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaIpsecTunnelProxyId_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o ipv4.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tfTemplate%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaIpsecTunnelProxyIdDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaIpsecTunnelProxyIdConfig(tmpl, name, "10.1.1.1", "10.2.1.1", 7, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIpsecTunnelProxyIdExists("panos_panorama_ipsec_tunnel_proxy_id_ipv4.test", &o),
					testAccCheckPanosPanoramaIpsecTunnelProxyIdAttributes(&o, name, "10.1.1.1", "10.2.1.1", 7, false),
				),
			},
			{
				Config: testAccPanoramaIpsecTunnelProxyIdConfig(tmpl, name, "10.3.1.1", "10.4.1.1", 0, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaIpsecTunnelProxyIdExists("panos_panorama_ipsec_tunnel_proxy_id_ipv4.test", &o),
					testAccCheckPanosPanoramaIpsecTunnelProxyIdAttributes(&o, name, "10.3.1.1", "10.4.1.1", 0, true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaIpsecTunnelProxyIdExists(n string, o *ipv4.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, tun, name := parsePanoramaIpsecTunnelProxyIdIpv4Id(rs.Primary.ID)
		v, err := pano.Network.IpsecTunnelProxyId.Get(tmpl, ts, tun, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaIpsecTunnelProxyIdAttributes(o *ipv4.Entry, name, loc, rem string, pn int, pa bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Local != loc {
			return fmt.Errorf("Local is %q, not %q", o.Local, loc)
		}

		if o.Remote != rem {
			return fmt.Errorf("Remote is %q, not %q", o.Remote, rem)
		}

		if o.ProtocolNumber != pn {
			return fmt.Errorf("Protocol number is %d, not %d", o.ProtocolNumber, pn)
		}

		if o.ProtocolAny != pa {
			return fmt.Errorf("Protocol any is %t, not %t", o.ProtocolAny, pa)
		}

		return nil
	}
}

func testAccPanosPanoramaIpsecTunnelProxyIdDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_ipsec_tunnel_proxy_id_ipv4" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, tun, name := parsePanoramaIpsecTunnelProxyIdIpv4Id(rs.Primary.ID)
			_, err := pano.Network.IpsecTunnelProxyId.Get(tmpl, ts, tun, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaIpsecTunnelProxyIdConfig(tmpl, name, loc, rem string, pn int, pa bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_tunnel_interface" "x" {
    template = panos_panorama_template.x.name
    name = "tunnel.7"
    comment = "For ipsec tunnel proxyid test"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/3"
    static_ips = ["10.2.3.1/24"]
    mode = "layer3"
}

resource "panos_panorama_ike_gateway" "x" {
    template = panos_panorama_template.x.name
    name = "accTestProxyIdGw"
    version = "ikev1"
    peer_ip_type = "ip"
    peer_ip_value = "10.2.4.6"
    interface = panos_panorama_ethernet_interface.x.name
    auth_type = "pre-shared-key"
    pre_shared_key = "secret"
    local_id_type = "ipaddr"
    local_id_value = "10.1.1.1"
    peer_id_type = "ipaddr"
    peer_id_value = "10.2.1.1"
}

resource "panos_panorama_ipsec_tunnel" "x" {
    template = panos_panorama_template.x.name
    name = "tfAccProxy"
    tunnel_interface = panos_panorama_tunnel_interface.x.name
    type = %q
    ak_ike_gateway = panos_panorama_ike_gateway.x.name
}

resource "panos_panorama_ipsec_tunnel_proxy_id_ipv4" "test" {
    template = panos_panorama_template.x.name
    name = %q
    ipsec_tunnel = panos_panorama_ipsec_tunnel.x.name
    local = %q
    remote = %q
    protocol_number = %d
    protocol_any = %t
}
`, tmpl, ipsectunnel.TypeAutoKey, name, loc, rem, pn, pa)
}
