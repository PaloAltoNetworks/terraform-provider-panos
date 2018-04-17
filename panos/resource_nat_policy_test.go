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

func TestAccPanosNatPolicy_basic(t *testing.T) {
	var o nat.Entry
	z1 := fmt.Sprintf("z%s", acctest.RandString(7))
	z2 := fmt.Sprintf("z%s", acctest.RandString(7))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosNatPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNatPolicyConfig(z1, z2, name, 1, 2, "192.168.1.1", "192.168.2.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosNatPolicyExists("panos_nat_policy.test", &o),
					testAccCheckPanosNatPolicyAttributes(&o, name, z1, z2, "192.168.1.1", "192.168.2.1"),
				),
			},
			{
				Config: testAccNatPolicyConfig(z1, z2, name, 2, 1, "192.168.3.1", "192.168.4.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosNatPolicyExists("panos_nat_policy.test", &o),
					testAccCheckPanosNatPolicyAttributes(&o, name, z2, z1, "192.168.3.1", "192.168.4.1"),
				),
			},
		},
	})
}

func testAccCheckPanosNatPolicyExists(n string, o *nat.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseNatPolicyId(rs.Primary.ID)
		v, err := fw.Policies.Nat.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosNatPolicyAttributes(o *nat.Entry, n, sz, dz, sta1, sta2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %q, expected %q", o.Name, n)
		}

		if len(o.SourceZones) != 1 || o.SourceZones[0] != sz {
			return fmt.Errorf("Source zones is %#v, expected [%s]", o.SourceZones, sz)
		}

		if o.DestinationZone != dz {
			return fmt.Errorf("Destination zone is %q, expected %q", o.DestinationZone, dz)
		}

		if o.SatType != "dynamic-ip-and-port" {
			return fmt.Errorf("SatType is %s, expected dynamic-ip-and-port", o.SatType)
		}

		if o.SatAddressType != "translated-address" {
			return fmt.Errorf("SatAddressType is %s, expected translated-address", o.SatAddressType)
		}

		if len(o.SatTranslatedAddresses) != 2 || o.SatTranslatedAddresses[0] != sta1 || o.SatTranslatedAddresses[1] != sta2 {
			return fmt.Errorf("SatTranslatedAddresses is %#v, expected [%s %s]", o.SatTranslatedAddresses, sta1, sta2)
		}

		return nil
	}
}

func testAccPanosNatPolicyDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_nat_policy" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseNatPolicyId(rs.Primary.ID)
			_, err := fw.Policies.Nat.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccNatPolicyConfig(z1, z2, n string, sz, dz int, sta1, sta2 string) string {
	return fmt.Sprintf(`
resource "panos_zone" "z1" {
    name = "%s"
    mode = "layer3"
}

resource "panos_zone" "z2" {
    name = "%s"
    mode = "layer3"
}

resource "panos_nat_policy" "test" {
    name = "%s"
    source_addresses = ["any"]
    destination_addresses = ["any"]
    source_zones = ["${panos_zone.z%d.name}"]
    destination_zone = "${panos_zone.z%d.name}"
    sat_type = "dynamic-ip-and-port"
    sat_address_type = "translated-address"
    sat_translated_addresses = ["%s", "%s"]
}
`, z1, z2, n, sz, dz, sta1, sta2)
}
