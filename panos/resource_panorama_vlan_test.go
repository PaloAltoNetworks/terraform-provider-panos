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

func TestAccPanosPanoramaVlan_basic(t *testing.T) {
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
		CheckDestroy: testAccPanosPanoramaVlanDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaVlanConfig(tmpl, name, "x"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVlanExists("panos_panorama_vlan.test", &o),
					testAccCheckPanosPanoramaVlanAttributes(&o, name, "x"),
				),
			},
			{
				Config: testAccPanoramaVlanConfig(tmpl, name, "y"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVlanExists("panos_panorama_vlan.test", &o),
					testAccCheckPanosPanoramaVlanAttributes(&o, name, "y"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaVlanExists(n string, o *vlan.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, _, name := parsePanoramaVlanId(rs.Primary.ID)
		v, err := pano.Network.Vlan.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaVlanAttributes(o *vlan.Entry, name, dep string) resource.TestCheckFunc {
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

func testAccPanosPanoramaVlanDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_vlan" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, _, name := parsePanoramaVlanId(rs.Primary.ID)
			_, err := pano.Network.Vlan.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaVlanConfig(tmpl, name, dep string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "vlan"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/5"
    mode = "layer2"
    vsys = "vsys1"
}

resource "panos_panorama_ethernet_interface" "y" {
    template = panos_panorama_template.x.name
    name = "ethernet1/6"
    mode = "layer2"
    vsys = "vsys1"
}

resource "panos_panorama_vlan_interface" "x" {
    template = panos_panorama_template.x.name
    name = "vlan.5"
    vsys = "vsys1"
}

resource "panos_panorama_vlan_interface" "y" {
    template = panos_panorama_template.x.name
    name = "vlan.6"
    vsys = "vsys1"
}

resource "panos_panorama_vlan" "test" {
    template = panos_panorama_template.x.name
    name = %q
    vlan_interface = panos_panorama_vlan_interface.%s.name
    interfaces = [panos_panorama_ethernet_interface.%s.name]
}
`, tmpl, name, dep, dep)
}
