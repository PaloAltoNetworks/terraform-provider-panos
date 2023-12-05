package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/tags"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosAdministrativeTag_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o tags.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAdministrativeTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAdministrativeTagConfig(name, "color1", "old comment"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAdministrativeTagExists("panos_administrative_tag.test", &o),
					testAccCheckPanosAdministrativeTagAttributes(&o, name, "color1", "old comment"),
				),
			},
			{
				Config: testAccAdministrativeTagConfig(name, "color12", "new comment"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAdministrativeTagExists("panos_administrative_tag.test", &o),
					testAccCheckPanosAdministrativeTagAttributes(&o, name, "color12", "new comment"),
				),
			},
		},
	})
}

func testAccCheckPanosAdministrativeTagExists(n string, o *tags.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseAdministrativeTagId(rs.Primary.ID)
		v, err := fw.Objects.Tags.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAdministrativeTagAttributes(o *tags.Entry, name, color, comment string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Color != color {
			return fmt.Errorf("Color is %s, expected %s", o.Color, color)
		}

		if o.Comment != comment {
			return fmt.Errorf("Comment is %s, expected %s", o.Comment, comment)
		}

		return nil
	}
}

func testAccPanosAdministrativeTagDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_administrative_tag" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseAdministrativeTagId(rs.Primary.ID)
			_, err := fw.Objects.Tags.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAdministrativeTagConfig(name, color, comment string) string {
	return fmt.Sprintf(`
resource "panos_administrative_tag" "test" {
    name = "%s"
    vsys = "vsys1"
    color = "%s"
    comment = "%s"
}
`, name, color, comment)
}
