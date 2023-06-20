package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/ipsectunnel"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosIpsecTunnel_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o ipsectunnel.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosIpsecTunnelDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIpsecTunnelConfig(name, true, "00001111", "11112222", "10.1.1.1", ipsectunnel.MkAuthTypeMd5, "00000000-11111111-22222222-33333333"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecTunnelExists("panos_ipsec_tunnel.test", &o),
					testAccCheckPanosIpsecTunnelAttributes(&o, name, true, "00001111", "11112222", "10.1.1.1", ipsectunnel.MkAuthTypeMd5, "00000000-11111111-22222222-33333333"),
				),
			},
			{
				Config: testAccIpsecTunnelConfig(name, false, "11112222", "00001111", "10.2.3.4", ipsectunnel.MkAuthTypeSha384, "00000001-00000002-00000003-00000004-00000005-00000006-00000007-00000008-00000009-0000000a-0000000b-0000000c"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIpsecTunnelExists("panos_ipsec_tunnel.test", &o),
					testAccCheckPanosIpsecTunnelAttributes(&o, name, false, "11112222", "00001111", "10.2.3.4", ipsectunnel.MkAuthTypeSha384, "00000001-00000002-00000003-00000004-00000005-00000006-00000007-00000008-00000009-0000000a-0000000b-0000000c"),
				),
			},
		},
	})
}

func testAccCheckPanosIpsecTunnelExists(n string, o *ipsectunnel.Entry) resource.TestCheckFunc {
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
		v, err := fw.Network.IpsecTunnel.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosIpsecTunnelAttributes(o *ipsectunnel.Entry, name string, tos bool, lspi, rspi, ra, at, ak string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.CopyTos != tos {
			return fmt.Errorf("Copy tos is %t, expected %t", o.CopyTos, tos)
		}

		if o.MkLocalSpi != lspi {
			return fmt.Errorf("Local SPI is %q, expected %q", o.MkLocalSpi, lspi)
		}

		if o.MkRemoteSpi != rspi {
			return fmt.Errorf("Remote SPI is %q, expected %q", o.MkRemoteSpi, rspi)
		}

		if o.MkRemoteAddress != ra {
			return fmt.Errorf("Remote address is %q, expected %q", o.MkRemoteAddress, ra)
		}

		if o.MkAuthType != at {
			return fmt.Errorf("Auth type is %q, expected %q", o.MkAuthType, at)
		}

		return nil
	}
}

func testAccPanosIpsecTunnelDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ipsec_tunnel" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := fw.Network.IpsecTunnel.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccIpsecTunnelConfig(name string, tos bool, lspi, rspi, ra, at, ak string) string {
	return fmt.Sprintf(`
resource "panos_tunnel_interface" "t" {
    name = "tunnel.3"
    comment = "For ipsec tunnel test"
}

resource "panos_loopback_interface" "l" {
    name = "loopback.4"
    comment = "For ipsec tunnel test"
}

resource "panos_ipsec_tunnel" "test" {
    name = %q
    tunnel_interface = panos_tunnel_interface.t.name
    copy_tos = %t
    type = %q
    mk_local_spi = %q
    mk_remote_spi = %q
    mk_interface = panos_loopback_interface.l.name
    mk_remote_address = %q
    mk_protocol = "ah"
    mk_auth_type = %q
    mk_auth_key = %q
}
`, name, tos, ipsectunnel.TypeManualKey, lspi, rspi, ra, at, ak)
}
