package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosEthernetInterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o eth.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("ethernet1/%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosEthernetInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEthernetInterfaceConfig(name, "down", "first comment", "10.1.1.1/24", "192.168.1.1/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosEthernetInterfaceExists("panos_ethernet_interface.test", &o),
					testAccCheckPanosEthernetInterfaceAttributes(&o, name, "down", "first comment", "10.1.1.1/24", "192.168.1.1/24"),
				),
			},
			{
				Config: testAccEthernetInterfaceConfig(name, "up", "second comment", "10.1.2.1/24", "192.168.2.1/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosEthernetInterfaceExists("panos_ethernet_interface.test", &o),
					testAccCheckPanosEthernetInterfaceAttributes(&o, name, "up", "second comment", "10.1.2.1/24", "192.168.2.1/24"),
				),
			},
		},
	})
}

func testAccCheckPanosEthernetInterfaceExists(n string, o *eth.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, name := parseEthernetInterfaceId(rs.Primary.ID)
		v, err := fw.Network.EthernetInterface.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosEthernetInterfaceAttributes(o *eth.Entry, n, ls, c, i1, i2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %s, expected %s", o.Name, n)
		}

		if o.LinkState != ls {
			return fmt.Errorf("Link state is %s, expected %s", o.LinkState, ls)
		}

		if o.Comment != c {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, c)
		}

		if len(o.StaticIps) != 2 {
			return fmt.Errorf("len(StaticIps) is %d, expected 2", len(o.StaticIps))
		}

		if o.StaticIps[0] != i1 {
			return fmt.Errorf("StaticIps[0] is %s, expected %s", o.StaticIps[0], i1)
		}

		if o.StaticIps[1] != i2 {
			return fmt.Errorf("StaticIps[1] is %s, expected %s", o.StaticIps[1], i2)
		}

		return nil
	}
}

func testAccPanosEthernetInterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ethernet_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			_, name := parseEthernetInterfaceId(rs.Primary.ID)
			_, err := fw.Network.EthernetInterface.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccEthernetInterfaceConfig(n, ls, c, i1, i2 string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "test" {
    name = "%s"
    vsys = "vsys1"
    mode = "layer3"
    link_state = "%s"
    comment = "%s"
    static_ips = ["%s", "%s"]
}
`, n, ls, c, i1, i2)
}
