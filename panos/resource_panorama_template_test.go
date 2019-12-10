package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaTemplate_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o template.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaTemplateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaTemplateConfig(name, "first description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateExists("panos_panorama_template.test", &o),
					testAccCheckPanosPanoramaTemplateAttributes(&o, "first description"),
				),
			},
			{
				Config: testAccPanoramaTemplateConfig(name, "second description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateExists("panos_panorama_template.test", &o),
					testAccCheckPanosPanoramaTemplateAttributes(&o, "second description"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaTemplateExists(n string, o *template.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		name := rs.Primary.ID
		v, err := pano.Panorama.Template.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaTemplateAttributes(o *template.Entry, desc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}
		if len(o.Devices) != 0 {
			return fmt.Errorf("Number of devices is %d, not 0", len(o.Devices))
		}
		return nil
	}
}

func testAccPanosPanoramaTemplateDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_template" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := pano.Panorama.Template.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaTemplateConfig(n, d string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "test" {
    name = %q
    description = %q
}
`, n, d)
}
