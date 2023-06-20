package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/pnrm/template/stack"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaTemplateStack_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o stack.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaTemplateStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaTemplateStackConfig(name, "first description", "a"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateStackExists("panos_panorama_template_stack.test", &o),
					testAccCheckPanosPanoramaTemplateStackAttributes(&o, "first description", "a"),
				),
			},
			{
				Config: testAccPanoramaTemplateStackConfig(name, "second description", "b"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateStackExists("panos_panorama_template_stack.test", &o),
					testAccCheckPanosPanoramaTemplateStackAttributes(&o, "second description", "b"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaTemplateStackExists(n string, o *stack.Entry) resource.TestCheckFunc {
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
		v, err := pano.Panorama.TemplateStack.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaTemplateStackAttributes(o *stack.Entry, desc, tmpl string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Description != desc {
			return fmt.Errorf("Description is %q, expected %q", o.Description, desc)
		}

		ts := fmt.Sprintf("tfAccTmpl%s", tmpl)
		if len(o.Templates) != 1 || o.Templates[0] != ts {
			return fmt.Errorf("Templates is %v, not [%s]", o.Templates, ts)
		}

		if len(o.Devices) != 0 {
			return fmt.Errorf("Number of devices is %d, not 0", len(o.Devices))
		}

		return nil
	}
}

func testAccPanosPanoramaTemplateStackDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_template_stack" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := pano.Panorama.TemplateStack.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaTemplateStackConfig(name, desc, tmpl string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "a" {
    name = "tfAccTmpla"
}

resource "panos_panorama_template" "b" {
    name = "tfAccTmplb"
}

resource "panos_panorama_template_stack" "test" {
    name = %q
    description = %q
    templates = [panos_panorama_template.%s.name]
}
`, name, desc, tmpl)
}
