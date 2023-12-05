package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/template/variable"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaTemplateVariable_basic(t *testing.T) {
	versionAdded := version.Number{
		Major: 8,
		Minor: 1,
	}

	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if !testAccPanosVersion.Gte(versionAdded) {
		t.Skip("Template variables are available in PAN-OS 8.1+")
	}

	var o variable.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("$tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaTemplateVariableDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaTemplateVariableConfig(tmpl, name, variable.TypeIpNetmask, "10.1.1.1/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateVariableExists("panos_panorama_template_variable.test", &o),
					testAccCheckPanosPanoramaTemplateVariableAttributes(&o, name, variable.TypeIpNetmask, "10.1.1.1/24"),
				),
			},
			{
				Config: testAccPanoramaTemplateVariableConfig(tmpl, name, variable.TypeIpRange, "10.1.1.1-10.1.1.255"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateVariableExists("panos_panorama_template_variable.test", &o),
					testAccCheckPanosPanoramaTemplateVariableAttributes(&o, name, variable.TypeIpRange, "10.1.1.1-10.1.1.255"),
				),
			},
			{
				Config: testAccPanoramaTemplateVariableConfig(tmpl, name, variable.TypeFqdn, "example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateVariableExists("panos_panorama_template_variable.test", &o),
					testAccCheckPanosPanoramaTemplateVariableAttributes(&o, name, variable.TypeFqdn, "example.com"),
				),
			},
			{
				Config: testAccPanoramaTemplateVariableConfig(tmpl, name, variable.TypeGroupId, "42"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateVariableExists("panos_panorama_template_variable.test", &o),
					testAccCheckPanosPanoramaTemplateVariableAttributes(&o, name, variable.TypeGroupId, "42"),
				),
			},
			{
				Config: testAccPanoramaTemplateVariableConfig(tmpl, name, variable.TypeInterface, "ethernet1/1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTemplateVariableExists("panos_panorama_template_variable.test", &o),
					testAccCheckPanosPanoramaTemplateVariableAttributes(&o, name, variable.TypeInterface, "ethernet1/1"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaTemplateVariableExists(n string, o *variable.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaTemplateVariableId(rs.Primary.ID)
		v, err := pano.Panorama.TemplateVariable.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaTemplateVariableAttributes(o *variable.Entry, name, typ, val string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Type != typ {
			return fmt.Errorf("Type is %s, expected %s", o.Type, typ)
		}

		if o.Value != val {
			return fmt.Errorf("Value is %s, expected %s", o.Value, val)
		}

		return nil
	}
}

func testAccPanosPanoramaTemplateVariableDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_template_variable" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaTemplateVariableId(rs.Primary.ID)
			if _, err := pano.Panorama.TemplateVariable.Get(tmpl, ts, name); err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaTemplateVariableConfig(tmpl, name, typ, val string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "For template variable acctest"
}

resource "panos_panorama_template_variable" "test" {
    template = panos_panorama_template.x.name
    name = %q
    type = %q
    value = %q
}
`, tmpl, name, typ, val)
}
