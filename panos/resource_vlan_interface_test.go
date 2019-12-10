package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/vlan"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosVlanInterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o vlan.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("vlan.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosVlanInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanInterfaceConfig(name, "first comment", "10.8.9.1/24", 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVlanInterfaceExists("panos_vlan_interface.test", &o),
					testAccCheckPanosVlanInterfaceAttributes(&o, name, "first comment", "10.8.9.1/24", 600),
				),
			},
			{
				Config: testAccVlanInterfaceConfig(name, "second comment", "10.9.10.1/24", 700),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVlanInterfaceExists("panos_vlan_interface.test", &o),
					testAccCheckPanosVlanInterfaceAttributes(&o, name, "second comment", "10.9.10.1/24", 700),
				),
			},
		},
	})
}

func testAccCheckPanosVlanInterfaceExists(n string, o *vlan.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, name := parseVlanInterfaceId(rs.Primary.ID)
		v, err := fw.Network.VlanInterface.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosVlanInterfaceAttributes(o *vlan.Entry, name, cmt, ip string, mtu int) resource.TestCheckFunc {
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

func testAccPanosVlanInterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_vlan_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			_, name := parseVlanInterfaceId(rs.Primary.ID)
			_, err := fw.Network.VlanInterface.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccVlanInterfaceConfig(name, cmt, ip string, mtu int) string {
	return fmt.Sprintf(`
resource "panos_vlan_interface" "test" {
    name = %q
    comment = %q
    static_ips = [%q]
    mtu = %d
}
`, name, cmt, ip, mtu)
}
