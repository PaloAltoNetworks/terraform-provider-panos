package panos

import (
	"fmt"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template/stack"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosTemplateStackEntry_basic(t *testing.T) {
	/*
	   In order to run this test you'll need:

	   * panorama as the device (PANOS_HOSTNAME, PANOS_USERNAME, PANOS_PASSWORD)
	   * a firewall that's already added as a managed device
	     (PANOS_MANAGED_SERIAL_NUMBER)
	*/

	dev := os.Getenv("PANOS_MANAGED_SERIAL_NUMBER")

	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if dev == "" {
		t.Skip("A required env variable is unset (PANOS_MANAGED_SERIAL_NUMBER)")
	}

	var o stack.Entry
	ts := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosTemplateStackEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTemplateStackEntryConfig(ts, dev),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTemplateStackEntryExists("panos_panorama_template_stack_entry.test", &o),
					testAccCheckPanosTemplateStackEntryAttributes(&o, dev),
				),
			},
		},
	})
}

func testAccCheckPanosTemplateStackEntryExists(n string, o *stack.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		ts, _ := parsePanoramaTemplateStackEntryId(rs.Primary.ID)
		v, err := pano.Panorama.TemplateStack.Get(ts)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}
		/*
			        ok = false
			        for i := range v.Devices {
			            if v.Devices[i] == serial {
			                ok = true
			                break
			            }
			        }
			        if !ok {
						return fmt.Errorf("Serial %s not in devices list", serial)
					}
		*/

		*o = v

		return nil
	}
}

func testAccCheckPanosTemplateStackEntryAttributes(o *stack.Entry, dev string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for i := range o.Devices {
			if o.Devices[i] == dev {
				return nil
			}
		}

		return fmt.Errorf("Serial %q is not in devices list: %#v", dev, o.Devices)
	}
}

func testAccPanosTemplateStackEntryDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_template_stack_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			ts, dev := parsePanoramaTemplateStackEntryId(rs.Primary.ID)
			o, err := pano.Panorama.TemplateStack.Get(ts)
			if err == nil {
				for i := range o.Devices {
					if o.Devices[i] == dev {
						return fmt.Errorf("Device entry %q still exists", dev)
					}
				}
			}
		}
		return nil
	}

	return nil
}

func testAccTemplateStackEntryConfig(ts, dev string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template_stack" "x" {
    name = %q
    description = "for template acc test"
}

resource "panos_panorama_template_stack_entry" "test" {
    template_stack = panos_panorama_template_stack.x.name
    device = %q
}
`, ts, dev)
}
