package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/ikegw"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosIkeGateway_basic(t *testing.T) {
	fqdnOptionVersion := version.Number{8, 1, 0, ""}

	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccPanosVersion.Gte(fqdnOptionVersion) {
		t.Skip("Peer IP type of FQDN available on PAN-OS 8.1+")
	}

	var mp ikegw.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosIkeGatewayDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIkeGatewayConfig(name, ikegw.PeerTypeIp, "192.168.1.1", ikegw.LocalTypeIp, "10.1.21.1", "secret1", ikegw.IdTypeIpAddress, "10.5.5.5", ikegw.IdTypeFqdn, "example.com", 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIkeGatewayExists("panos_ike_gateway.test", &mp),
					testAccCheckPanosIkeGatewayAttributes(&mp, name, ikegw.PeerTypeIp, "192.168.1.1", ikegw.LocalTypeIp, "10.1.21.1", "secret1", ikegw.IdTypeIpAddress, "10.5.5.5", ikegw.IdTypeFqdn, "example.com", 1),
				),
			},
			{
				Config: testAccIkeGatewayConfig(name, ikegw.PeerTypeFqdn, "foobar.com", ikegw.LocalTypeIp, "10.2.21.1", "secret2", ikegw.IdTypeFqdn, "acctest.org", ikegw.IdTypeKeyId, "beef", 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosIkeGatewayExists("panos_ike_gateway.test", &mp),
					testAccCheckPanosIkeGatewayAttributes(&mp, name, ikegw.PeerTypeFqdn, "foobar.com", ikegw.LocalTypeIp, "10.2.21.1", "secret2", ikegw.IdTypeFqdn, "acctest.org", ikegw.IdTypeKeyId, "beef", 2),
				),
			},
		},
	})
}

func testAccCheckPanosIkeGatewayExists(n string, o *ikegw.Entry) resource.TestCheckFunc {
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
		v, err := fw.Network.IkeGateway.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosIkeGatewayAttributes(o *ikegw.Entry, name, pipt, pipv, liat, liav, psk, lit, liv, pidt, pidv string, prof int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.PeerIpType != pipt {
			return fmt.Errorf("Peer ip type is %q, not %q", o.PeerIpType, pipt)
		}

		if o.PeerIpValue != pipv {
			return fmt.Errorf("Peer ip value is %q, not %q", o.PeerIpValue, pipv)
		}

		if o.LocalIpAddressType != liat {
			return fmt.Errorf("Local ip address type is %q, not %q", o.LocalIpAddressType, liat)
		}

		if o.LocalIpAddressValue != liav {
			return fmt.Errorf("Local ip address value is %q, not %q", o.LocalIpAddressValue, liav)
		}

		// Skip pre_shared_key, as it's encrypted.

		if o.LocalIdType != lit {
			return fmt.Errorf("Local id type is %q, not %q", o.LocalIdType, lit)
		}

		if o.LocalIdValue != liv {
			return fmt.Errorf("Local id value is %q, not %q", o.LocalIdValue, liv)
		}

		if o.PeerIdType != pidt {
			return fmt.Errorf("Peer id type is %q, not %q", o.PeerIdType, pidt)
		}

		if o.PeerIdValue != pidv {
			return fmt.Errorf("Peer id value is %q, not %q", o.PeerIdValue, pidv)
		}

		if o.Ikev1CryptoProfile != fmt.Sprintf("prof%d", prof) {
			return fmt.Errorf("Ikev1 crypto profile is %q, not \"prof%d\"", o.Ikev1CryptoProfile, prof)
		}

		return nil
	}
}

func testAccPanosIkeGatewayDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ike_gateway" {
			continue
		}

		if rs.Primary.ID != "" {
			_, err := fw.Network.IkeGateway.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccIkeGatewayConfig(name, pipt, pipv, liat, liav, psk, lit, liv, pidt, pidv string, prof int) string {
	return fmt.Sprintf(`
resource "panos_loopback_interface" "lo" {
    name = "loopback.42"
    static_ips = ["10.1.21.1", "10.2.21.1"]
}

resource "panos_ike_crypto_profile" "prof1" {
    name = "prof1"
    dh_groups = ["group1"]
    authentications = ["md5"]
    encryptions = ["3des"]
    lifetime_type = "hours"
    lifetime_value = 8
}

resource "panos_ike_crypto_profile" "prof2" {
    name = "prof2"
    dh_groups = ["group1"]
    authentications = ["md5"]
    encryptions = ["3des"]
    lifetime_type = "hours"
    lifetime_value = 8
}

resource "panos_ike_gateway" "test" {
    name = %q
    version = "ikev1"
    peer_ip_type = %q
    peer_ip_value = %q
    interface = panos_loopback_interface.lo.name
    local_ip_address_type = %q
    local_ip_address_value = %q
    auth_type = %q
    pre_shared_key = %q
    local_id_type = %q
    local_id_value = %q
    peer_id_type = %q
    peer_id_value = %q
    ikev1_crypto_profile = panos_ike_crypto_profile.prof%d.name
}
`, name, pipt, pipv, liat, liav, ikegw.AuthPreSharedKey, psk, lit, liv, pidt, pidv, prof)
}
