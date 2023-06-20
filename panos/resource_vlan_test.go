package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/vlan"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosVlan_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o vlan.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosVlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanConfig(name, "x"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVlanExists("panos_vlan.test", &o),
					testAccCheckPanosVlanAttributes(&o, name, "x"),
				),
			},
			{
				Config: testAccVlanConfig(name, "y"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVlanExists("panos_vlan.test", &o),
					testAccCheckPanosVlanAttributes(&o, name, "y"),
				),
			},
		},
	})
}

func testAccCheckPanosVlanExists(n string, o *vlan.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, name := parseVlanId(rs.Primary.ID)
		v, err := fw.Network.Vlan.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosVlanAttributes(o *vlan.Entry, name, dep string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		eths := map[string]string{
			"x": "ethernet1/5",
			"y": "ethernet1/6",
		}

		vis := map[string]string{
			"x": "vlan.5",
			"y": "vlan.6",
		}

		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.Interfaces) != 1 || o.Interfaces[0] != eths[dep] {
			return fmt.Errorf("Interfaces is %#v, not [%s]", o.Interfaces, eths[dep])
		}

		if o.VlanInterface != vis[dep] {
			return fmt.Errorf("VlanInterface is %s, not %s", o.VlanInterface, vis[dep])
		}

		return nil
	}
}

func testAccPanosVlanDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_vlan" {
			continue
		}

		if rs.Primary.ID != "" {
			_, name := parseVlanId(rs.Primary.ID)
			_, err := fw.Network.Vlan.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccVlanConfig(name, dep string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "x" {
    name = "ethernet1/5"
    mode = "layer2"
    vsys = "vsys1"
}

resource "panos_ethernet_interface" "y" {
    name = "ethernet1/6"
    mode = "layer2"
    vsys = "vsys1"
}

resource "panos_vlan_interface" "x" {
    name = "vlan.5"
    vsys = "vsys1"
}

resource "panos_vlan_interface" "y" {
    name = "vlan.6"
    vsys = "vsys1"
}

resource "panos_vlan" "test" {
    name = %q
    vlan_interface = panos_vlan_interface.%s.name
    interfaces = [panos_ethernet_interface.%s.name]
}
`, name, dep, dep)
}
