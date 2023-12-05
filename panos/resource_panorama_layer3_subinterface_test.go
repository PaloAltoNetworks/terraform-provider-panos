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

func TestAccPanosPanoramaLayer3Subinterface_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o layer3.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("ethernet1/5.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaLayer3SubinterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaLayer3SubinterfaceConfig(tmpl, name, "x", "desc1", "192.168.55.1/24", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLayer3SubinterfaceExists("panos_panorama_layer3_subinterface.test", &o),
					testAccCheckPanosPanoramaLayer3SubinterfaceAttributes(&o, name, "x", "desc1", "192.168.55.1/24", 5),
				),
			},
			{
				Config: testAccPanoramaLayer3SubinterfaceConfig(tmpl, name, "y", "desc2", "192.168.66.1/24", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLayer3SubinterfaceExists("panos_panorama_layer3_subinterface.test", &o),
					testAccCheckPanosPanoramaLayer3SubinterfaceAttributes(&o, name, "y", "desc2", "192.168.66.1/24", 5),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaLayer3SubinterfaceExists(n string, o *layer3.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, iType, eth, _, name := parsePanoramaLayer3SubinterfaceId(rs.Primary.ID)
		v, err := pano.Network.Layer3Subinterface.Get(tmpl, ts, iType, eth, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaLayer3SubinterfaceAttributes(o *layer3.Entry, name, mp, com, ip string, tag int) resource.TestCheckFunc {
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

func testAccPanosPanoramaLayer3SubinterfaceDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_layer3_subinterface" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, iType, eth, _, name := parsePanoramaLayer3SubinterfaceId(rs.Primary.ID)
			_, err := pano.Network.Layer3Subinterface.Get(tmpl, ts, iType, eth, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaLayer3SubinterfaceConfig(tmpl, name, mp, com, ip string, tag int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "layer3 subinterface test"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer3"
    comment = "for layer3 test"
}

resource "panos_panorama_management_profile" "x" {
    template = panos_panorama_template.x.name
    name = "x"
    ping = true
}

resource "panos_panorama_management_profile" "y" {
    template = panos_panorama_template.x.name
    name = "y"
    ssh = true
}

resource "panos_panorama_layer3_subinterface" "test" {
    template = panos_panorama_template.x.name
    name = %q
    parent_interface = panos_panorama_ethernet_interface.x.name
    management_profile = panos_panorama_management_profile.%s.name
    comment = %q
    static_ips = [%q]
    tag = %d
}
`, tmpl, name, mp, com, ip, tag)
}
