package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosZoneEntry_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o zone.Entry
	eth_name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%7+1)
	zone_name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosZoneEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneEntryConfig(eth_name, zone_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosZoneEntryExists("panos_zone_entry.test", &o),
					testAccCheckPanosZoneEntryAttributes(&o, eth_name),
				),
			},
		},
	})
}

func testAccCheckPanosZoneEntryExists(n string, o *zone.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, _, vsys, zone_name, _, _ := parseZoneEntryId(rs.Primary.ID)
		v, err := fw.Network.Zone.Get(vsys, zone_name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosZoneEntryAttributes(o *zone.Entry, eth_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(o.Interfaces) != 1 || o.Interfaces[0] != eth_name {
			return fmt.Errorf("Zone interfaces is %#v, not [%s]", o.Interfaces, eth_name)
		}

		return nil
	}
}

func testAccPanosZoneEntryDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_zone_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			_, _, vsys, zone_name, _, _ := parseZoneEntryId(rs.Primary.ID)
			_, err := fw.Network.Zone.Get(vsys, zone_name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccZoneEntryConfig(eth_name, zone_name string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "eth" {
    name = %q
    mode = "layer3"
}

resource "panos_zone" "z" {
    name = %q
    mode = "layer3"
}

resource "panos_zone_entry" "test" {
    zone = panos_zone.z.name
    mode = panos_zone.z.mode
    interface = panos_ethernet_interface.eth.name
}
`, eth_name, zone_name)
}
