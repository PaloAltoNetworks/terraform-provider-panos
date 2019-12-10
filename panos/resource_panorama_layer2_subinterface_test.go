package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer2"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaLayer2Subinterface_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if !testAccSupportsL2 {
		t.Skip(SkipL2AccTest)
	}

	var o layer2.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("ethernet1/5.%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaLayer2SubinterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaLayer2SubinterfaceConfig(tmpl, name, "desc1", 5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLayer2SubinterfaceExists("panos_panorama_layer2_subinterface.test", &o),
					testAccCheckPanosPanoramaLayer2SubinterfaceAttributes(&o, name, "desc1", 5),
				),
			},
			{
				Config: testAccPanoramaLayer2SubinterfaceConfig(tmpl, name, "desc2", 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLayer2SubinterfaceExists("panos_panorama_layer2_subinterface.test", &o),
					testAccCheckPanosPanoramaLayer2SubinterfaceAttributes(&o, name, "desc2", 7),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaLayer2SubinterfaceExists(n string, o *layer2.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, iType, eth, mType, _, name := parsePanoramaLayer2SubinterfaceId(rs.Primary.ID)
		v, err := pano.Network.Layer2Subinterface.Get(tmpl, ts, iType, eth, mType, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaLayer2SubinterfaceAttributes(o *layer2.Entry, name, com string, tag int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Comment != com {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, com)
		}

		if o.Tag != tag {
			return fmt.Errorf("Tag is %d, not %d", o.Tag, tag)
		}

		return nil
	}
}

func testAccPanosPanoramaLayer2SubinterfaceDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_layer2_subinterface" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, iType, eth, mType, _, name := parsePanoramaLayer2SubinterfaceId(rs.Primary.ID)
			_, err := pano.Network.Layer2Subinterface.Get(tmpl, ts, iType, eth, mType, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaLayer2SubinterfaceConfig(tmpl, name, com string, tag int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "layer2 subinterface test"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer2"
    comment = "for layer2 test"
}

resource "panos_panorama_layer2_subinterface" "test" {
    template = panos_panorama_template.x.name
    name = %q
    parent_interface = panos_panorama_ethernet_interface.x.name
    parent_mode = panos_panorama_ethernet_interface.x.mode
    comment = %q
    tag = %d
}
`, tmpl, name, com, tag)
}
