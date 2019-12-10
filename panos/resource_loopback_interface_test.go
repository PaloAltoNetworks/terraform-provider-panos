package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/loopback"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosLoopbackInterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o loopback.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("loopback.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosLoopbackInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLoopbackInterfaceConfig(name, "first comment", "10.8.9.1", 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLoopbackInterfaceExists("panos_loopback_interface.test", &o),
					testAccCheckPanosLoopbackInterfaceAttributes(&o, name, "first comment", "10.8.9.1", 600),
				),
			},
			{
				Config: testAccLoopbackInterfaceConfig(name, "second comment", "10.9.10.1", 700),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLoopbackInterfaceExists("panos_loopback_interface.test", &o),
					testAccCheckPanosLoopbackInterfaceAttributes(&o, name, "second comment", "10.9.10.1", 700),
				),
			},
		},
	})
}

func testAccCheckPanosLoopbackInterfaceExists(n string, o *loopback.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, name := parseLoopbackInterfaceId(rs.Primary.ID)
		v, err := fw.Network.LoopbackInterface.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosLoopbackInterfaceAttributes(o *loopback.Entry, name, cmt, ip string, mtu int) resource.TestCheckFunc {
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

func testAccPanosLoopbackInterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_loopback_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			_, name := parseLoopbackInterfaceId(rs.Primary.ID)
			_, err := fw.Network.LoopbackInterface.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccLoopbackInterfaceConfig(name, cmt, ip string, mtu int) string {
	return fmt.Sprintf(`
resource "panos_loopback_interface" "test" {
    name = %q
    comment = %q
    static_ips = [%q]
    mtu = %d
}
`, name, cmt, ip, mtu)
}
