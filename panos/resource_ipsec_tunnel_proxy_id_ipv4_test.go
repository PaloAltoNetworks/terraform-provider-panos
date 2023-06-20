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

func TestAccPanosIpsecTunnelProxyId_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o ipv4.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosIpsecTunnelProxyIdDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecTunnelProxyIdConfig(name, "10.1.1.1", "10.2.1.1", 7, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecTunnelProxyIdExists("panos_ipsec_tunnel_proxy_id_ipv4.test", &o),
					testAccCheckPanosIpsecTunnelProxyIdAttributes(&o, name, "10.1.1.1", "10.2.1.1", 7, false),
				),
			},
			{
				Config: testAccIpsecTunnelProxyIdConfig(name, "10.3.1.1", "10.4.1.1", 0, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecTunnelProxyIdExists("panos_ipsec_tunnel_proxy_id_ipv4.test", &o),
					testAccCheckPanosIpsecTunnelProxyIdAttributes(&o, name, "10.3.1.1", "10.4.1.1", 0, true),
				),
			},
		},
	})
}

func testAccCheckPanosIpsecTunnelProxyIdExists(n string, o *ipv4.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		tun, name := parseIpsecTunnelProxyIdIpv4Id(rs.Primary.ID)
		v, err := fw.Network.IpsecTunnelProxyId.Get(tun, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosIpsecTunnelProxyIdAttributes(o *ipv4.Entry, name, loc, rem string, pn int, pa bool) resource.TestCheckFunc {
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

func testAccPanosIpsecTunnelProxyIdDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ipsec_tunnel_proxy_id_ipv4" {
			continue
		}

		if rs.Primary.ID != "" {
			tun, name := parseIpsecTunnelProxyIdIpv4Id(rs.Primary.ID)
			_, err := fw.Network.IpsecTunnelProxyId.Get(tun, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccIpsecTunnelProxyIdConfig(name, loc, rem string, pn int, pa bool) string {
	return fmt.Sprintf(`
resource "panos_tunnel_interface" "x" {
    name = "tunnel.7"
    comment = "For ipsec tunnel proxyid test"
}

resource "panos_ethernet_interface" "x" {
    name = "ethernet1/3"
    static_ips = ["10.2.3.1/24"]
    mode = "layer3"
}

resource "panos_ike_gateway" "x" {
    name = "accTestProxyIdGw"
    version = "ikev1"
    peer_ip_type = "ip"
    peer_ip_value = "10.2.4.6"
    interface = panos_ethernet_interface.x.name
    auth_type = "pre-shared-key"
    pre_shared_key = "secret"
    local_id_type = "ipaddr"
    local_id_value = "10.1.1.1"
    peer_id_type = "ipaddr"
    peer_id_value = "10.2.1.1"
}

resource "panos_ipsec_tunnel" "x" {
    name = "tfAccProxy"
    tunnel_interface = panos_tunnel_interface.x.name
    type = %q
    ak_ike_gateway = panos_ike_gateway.x.name
}

resource "panos_ipsec_tunnel_proxy_id_ipv4" "test" {
    ipsec_tunnel = panos_ipsec_tunnel.x.name
    name = %q
    local = %q
    remote = %q
    protocol_number = %d
    protocol_any = %t
}
`, ipsectunnel.TypeAutoKey, name, loc, rem, pn, pa)
}
