package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaNatRule_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o nat.Entry
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaNatRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaNatRuleConfig(dg, name, "first description", "192.168.1.1", 5555),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaNatRuleExists("panos_panorama_nat_rule.test", &o),
					testAccCheckPanosPanoramaNatRuleAttributes(&o, name, "first description", "192.168.1.1", 5555),
				),
			},
			{
				Config: testAccPanoramaNatRuleConfig(dg, name, "second desc", "192.168.3.1", 6666),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaNatRuleExists("panos_panorama_nat_rule.test", &o),
					testAccCheckPanosPanoramaNatRuleAttributes(&o, name, "second desc", "192.168.3.1", 6666),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaNatRuleExists(n string, o *nat.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, rb, name := parsePanoramaNatRuleId(rs.Primary.ID)
		v, err := pano.Policies.Nat.Get(dg, rb, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaNatRuleAttributes(o *nat.Entry, n, de, da string, dp int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %q, expected %q", o.Name, n)
		}

		if o.Description != de {
			return fmt.Errorf("Description is %q, expected %q", o.Description, de)
		}

		if o.DatAddress != da {
			return fmt.Errorf("DatAddress is %q, expected %q", o.DatAddress, da)
		}

		if o.DatPort != dp {
			return fmt.Errorf("DatPort is %d, expected %d", o.DatPort, dp)
		}

		return nil
	}
}

func testAccPanosPanoramaNatRuleDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_nat_rule" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, rb, name := parsePanoramaNatRuleId(rs.Primary.ID)
			_, err := pano.Policies.Nat.Get(dg, rb, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaNatRuleConfig(dg, n, de, da string, dp int) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
}

resource "panos_panorama_nat_rule" "test" {
    device_group = panos_panorama_device_group.x.name
    name = "%s"
    description = "%s"
    source_zones = ["any"]
    destination_zone = "myZone"
    to_interface = "any"
    service = "any"
    source_addresses = ["any"]
    destination_addresses = ["any"]
    sat_type = "none"
    dat_type = "static"
    dat_address = "%s"
    dat_port = "%d"
}
`, dg, n, de, da, dp)
}
