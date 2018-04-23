package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/nat"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaNatPolicy_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o nat.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaNatPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaNatPolicyConfig(name, "first description", "192.168.1.1", 5555),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaNatPolicyExists("panos_panorama_nat_policy.test", &o),
					testAccCheckPanosPanoramaNatPolicyAttributes(&o, name, "first description", "192.168.1.1", 5555),
				),
			},
			{
				Config: testAccPanoramaNatPolicyConfig(name, "second desc", "192.168.3.1", 6666),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaNatPolicyExists("panos_panorama_nat_policy.test", &o),
					testAccCheckPanosPanoramaNatPolicyAttributes(&o, name, "second desc", "192.168.3.1", 6666),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaNatPolicyExists(n string, o *nat.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, rb, name := parsePanoramaNatPolicyId(rs.Primary.ID)
		v, err := pano.Policies.Nat.Get(dg, rb, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaNatPolicyAttributes(o *nat.Entry, n, de, da string, dp int) resource.TestCheckFunc {
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

func testAccPanosPanoramaNatPolicyDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_nat_policy" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, rb, name := parsePanoramaNatPolicyId(rs.Primary.ID)
			_, err := pano.Policies.Nat.Get(dg, rb, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaNatPolicyConfig(n, de, da string, dp int) string {
	return fmt.Sprintf(`
resource "panos_panorama_nat_policy" "test" {
    name = "%s"
    description = "%s"
    source_zones = ["any"]
    destination_zone = "any"
    to_interface = "any"
    service = "any"
    source_addresses = ["any"]
    destination_addresses = ["any"]
    sat_type = "none"
    dat_address = "%s"
    dat_port = "%d"
}
`, n, de, da, dp)
}
