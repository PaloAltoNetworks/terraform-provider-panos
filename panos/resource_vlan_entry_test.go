package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/vlan"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosVlanEntry_basic(t *testing.T) {
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
		CheckDestroy: testAccPanosVlanEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVlanEntryConfig(name, "00:30:48:52:aa:bb"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVlanEntryExists("panos_vlan_entry.test", &o),
					testAccCheckPanosVlanEntryAttributes(&o, name, "00:30:48:52:aa:bb"),
				),
			},
			{
				Config: testAccVlanEntryConfig(name, "00:30:48:52:cc:dd"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVlanEntryExists("panos_vlan_entry.test", &o),
					testAccCheckPanosVlanEntryAttributes(&o, name, "00:30:48:52:cc:dd"),
				),
			},
		},
	})
}

func testAccCheckPanosVlanEntryExists(n string, o *vlan.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vName, _ := parseVlanEntryId(rs.Primary.ID)
		v, err := fw.Network.Vlan.Get(vName)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosVlanEntryAttributes(o *vlan.Entry, name, mac string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if val := o.StaticMacs[mac]; val != "ethernet1/5" {
			return fmt.Errorf("MAC is %s, not 'ethernet1/5'", val)
		}

		return nil
	}
}

func testAccPanosVlanEntryDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_vlan_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			vName, _ := parseVlanEntryId(rs.Primary.ID)
			_, err := fw.Network.Vlan.Get(vName)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccVlanEntryConfig(name, mac string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "x" {
    name = "ethernet1/5"
    mode = "layer2"
    vsys = "vsys1"
}

resource "panos_vlan" "x" {
    name = %q
    vsys = "vsys1"
}

resource "panos_vlan_entry" "test" {
    vlan = panos_vlan.x.name
    interface = panos_ethernet_interface.x.name
    mac_addresses = [%q]
}
`, name, mac)
}
