package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaZoneEntry_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o zone.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	eth_name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%7+1)
	zone_name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaZoneEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaZoneEntryConfig(tmpl, eth_name, zone_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaZoneEntryExists("panos_panorama_zone_entry.test", &o),
					testAccCheckPanosPanoramaZoneEntryAttributes(&o, eth_name),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaZoneEntryExists(n string, o *zone.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vsys, zone_name, _, _ := parseZoneEntryId(rs.Primary.ID)
		v, err := pano.Network.Zone.Get(tmpl, ts, vsys, zone_name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaZoneEntryAttributes(o *zone.Entry, eth_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(o.Interfaces) != 1 || o.Interfaces[0] != eth_name {
			return fmt.Errorf("Zone interfaces is %#v, not [%s]", o.Interfaces, eth_name)
		}

		return nil
	}
}

func testAccPanosPanoramaZoneEntryDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_zone_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vsys, zone_name, _, _ := parseZoneEntryId(rs.Primary.ID)
			_, err := pano.Network.Zone.Get(tmpl, ts, vsys, zone_name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaZoneEntryConfig(tmpl, eth_name, zone_name string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "tmpl" {
    name = %q
    description = "for template acc test"
}

resource "panos_panorama_ethernet_interface" "eth" {
    template = panos_panorama_template.tmpl.name
    name = %q
    mode = "layer3"
}

resource "panos_panorama_zone" "z" {
    template = panos_panorama_template.tmpl.name
    name = %q
    mode = "layer3"
}

resource "panos_panorama_zone_entry" "test" {
    template = panos_panorama_template.tmpl.name
    zone = panos_panorama_zone.z.name
    mode = panos_panorama_zone.z.mode
    interface = panos_panorama_ethernet_interface.eth.name
}
`, tmpl, eth_name, zone_name)
}
