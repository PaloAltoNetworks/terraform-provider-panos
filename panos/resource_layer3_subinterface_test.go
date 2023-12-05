package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosLayer3Subinterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o layer3.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("ethernet1/5.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosLayer3SubinterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLayer3SubinterfaceConfig(name, "x", "desc1", "192.168.55.1/24", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLayer3SubinterfaceExists("panos_layer3_subinterface.test", &o),
					testAccCheckPanosLayer3SubinterfaceAttributes(&o, name, "x", "desc1", "192.168.55.1/24", 5),
				),
			},
			{
				Config: testAccLayer3SubinterfaceConfig(name, "y", "desc2", "192.168.66.1/24", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosLayer3SubinterfaceExists("panos_layer3_subinterface.test", &o),
					testAccCheckPanosLayer3SubinterfaceAttributes(&o, name, "y", "desc2", "192.168.66.1/24", 5),
				),
			},
		},
	})
}

func testAccCheckPanosLayer3SubinterfaceExists(n string, o *layer3.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		iType, eth, _, name := parseLayer3SubinterfaceId(rs.Primary.ID)
		v, err := fw.Network.Layer3Subinterface.Get(iType, eth, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosLayer3SubinterfaceAttributes(o *layer3.Entry, name, mp, com, ip string, tag int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.ManagementProfile != mp {
			return fmt.Errorf("Management profile is %q, expected %q", o.ManagementProfile, mp)
		}

		if o.Comment != com {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, com)
		}

		if len(o.StaticIps) != 1 || o.StaticIps[0] != ip {
			return fmt.Errorf("Static IPs is %#v, not [%s]", o.StaticIps, ip)
		}

		if o.Tag != tag {
			return fmt.Errorf("Tag is %d, not %d", o.Tag, tag)
		}

		return nil
	}
}

func testAccPanosLayer3SubinterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_layer3_subinterface" {
			continue
		}

		if rs.Primary.ID != "" {
			iType, eth, _, name := parseLayer3SubinterfaceId(rs.Primary.ID)
			_, err := fw.Network.Layer3Subinterface.Get(iType, eth, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccLayer3SubinterfaceConfig(name, mp, com, ip string, tag int) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "x" {
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer3"
    comment = "for layer3 test"
}

resource "panos_management_profile" "x" {
    name = "x"
    ping = true
}

resource "panos_management_profile" "y" {
    name = "y"
    ssh = true
}

resource "panos_layer3_subinterface" "test" {
    name = %q
    parent_interface = panos_ethernet_interface.x.name
    management_profile = panos_management_profile.%s.name
    comment = %q
    static_ips = [%q]
    tag = %d
}
`, name, mp, com, ip, tag)
}
