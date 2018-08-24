package panos

import (
	"fmt"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosTemplateEntry_basic(t *testing.T) {
	/*
			   In order to run this test you'll need:

			   * panorama as the device (PANOS_HOSTNAME, PANOS_USERNAME, PANOS_PASSWORD)
			   * a multi-vsys firewall that's already added as a managed device
			     (PANOS_MANAGED_SERIAL_NUMBER)
			   * and two vsys that are currently unassigned to any template
			     (PANOS_MANAGED_VSYS1 and PANOS_MANAGED_VSYS2)
		       * panorama has to be < 8.1.0
	*/

	serial := os.Getenv("PANOS_MANAGED_SERIAL_NUMBER")
	vsys1 := os.Getenv("PANOS_MANAGED_VSYS1")
	vsys2 := os.Getenv("PANOS_MANAGED_VSYS2")
    versionRemoved := version.Number{
        Major: 8,
        Minor: 1,
    }

	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if testAccPanosVersion.Gte(versionRemoved) {
		t.Skip("This test is only valid for PAN-OS < 8.1.0.")
	} else if serial == "" || vsys1 == "" || vsys2 == "" {
		t.Skip("One or more required env variables are unset (PANOS_MANAGED_SERIAL_NUMBER, PANOS_MANAGED_VSYS1, PANOS_MANAGED_VSYS2")
	}

	var o template.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosTemplateEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateEntryConfig(tmpl, serial, vsys1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTemplateEntryExists("panos_panorama_template_entry.test", &o),
					testAccCheckPanosTemplateEntryAttributes(&o, serial, vsys1),
				),
			},
			{
				Config: testAccTemplateEntryConfig(tmpl, serial, vsys2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTemplateEntryExists("panos_panorama_template_entry.test", &o),
					testAccCheckPanosTemplateEntryAttributes(&o, serial, vsys2),
				),
			},
		},
	})
}

func testAccCheckPanosTemplateEntryExists(n string, o *template.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, serial := parsePanoramaTemplateEntryId(rs.Primary.ID)
		v, err := pano.Panorama.Template.Get(tmpl)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		} else if _, ok = v.Devices[serial]; !ok {
			return fmt.Errorf("Serial %s not in devices list", serial)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosTemplateEntryAttributes(o *template.Entry, serial, vsys string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(o.Devices[serial]) != 1 || o.Devices[serial][0] != vsys {
			return fmt.Errorf("Vsys list for serial %q is %#v, not [%s]", serial, o.Devices[serial], vsys)
		}
		return nil
	}
}

func testAccPanosTemplateEntryDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_template_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, serial := parsePanoramaTemplateEntryId(rs.Primary.ID)
			o, err := pano.Panorama.Template.Get(tmpl)
			if err == nil {
				if _, ok := o.Devices[serial]; ok {
					return fmt.Errorf("Object %q still exists", rs.Primary.ID)
				}
			}
		}
		return nil
	}

	return nil
}

func testAccTemplateEntryConfig(tmpl, serial, vsys string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "tmpl" {
    name = %q
    description = "for template acc test"
}

resource "panos_panorama_template_entry" "test" {
    template = "${panos_panorama_template.tmpl.name}"
    serial = %q
    vsys_list = [%q]
}
`, tmpl, serial, vsys)
}
