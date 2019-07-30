package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/vlan"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaVlanEntry_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o vlan.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaVlanEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaVlanEntryConfig(tmpl, name, "00:30:48:52:aa:bb"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVlanEntryExists("panos_panorama_vlan_entry.test", &o),
					testAccCheckPanosPanoramaVlanEntryAttributes(&o, name, "00:30:48:52:aa:bb"),
				),
			},
			{
				Config: testAccPanoramaVlanEntryConfig(tmpl, name, "00:30:48:52:cc:dd"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVlanEntryExists("panos_panorama_vlan_entry.test", &o),
					testAccCheckPanosPanoramaVlanEntryAttributes(&o, name, "00:30:48:52:cc:dd"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaVlanEntryExists(n string, o *vlan.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vName, _ := parsePanoramaVlanEntryId(rs.Primary.ID)
		v, err := pano.Network.Vlan.Get(tmpl, ts, vName)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaVlanEntryAttributes(o *vlan.Entry, name, mac string) resource.TestCheckFunc {
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

func testAccPanosPanoramaVlanEntryDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_vlan_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vName, _ := parsePanoramaVlanEntryId(rs.Primary.ID)
			_, err := pano.Network.Vlan.Get(tmpl, ts, vName)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaVlanEntryConfig(tmpl, name, mac string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "t" {
    name = %q
    description = "vlan entry test"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.t.name
    name = "ethernet1/5"
    mode = "layer2"
    vsys = "vsys1"
}

resource "panos_panorama_vlan" "x" {
    template = panos_panorama_template.t.name
    name = %q
    vsys = "vsys1"
}

resource "panos_panorama_vlan_entry" "test" {
    template = panos_panorama_template.t.name
    vlan = panos_panorama_vlan.x.name
    interface = panos_panorama_ethernet_interface.x.name
    mac_addresses = [%q]
}
`, tmpl, name, mac)
}
