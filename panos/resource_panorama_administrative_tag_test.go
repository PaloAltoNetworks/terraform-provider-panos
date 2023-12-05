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

func TestAccPanosPanoramaAdministrativeTag_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o tags.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaAdministrativeTagDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaAdministrativeTagConfig(name, "color1", "old comment"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAdministrativeTagExists("panos_panorama_administrative_tag.test", &o),
					testAccCheckPanosPanoramaAdministrativeTagAttributes(&o, name, "color1", "old comment"),
				),
			},
			{
				Config: testAccPanoramaAdministrativeTagConfig(name, "color12", "new comment"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAdministrativeTagExists("panos_panorama_administrative_tag.test", &o),
					testAccCheckPanosPanoramaAdministrativeTagAttributes(&o, name, "color12", "new comment"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaAdministrativeTagExists(n string, o *tags.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaAdministrativeTagId(rs.Primary.ID)
		v, err := pano.Objects.Tags.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaAdministrativeTagAttributes(o *tags.Entry, name, color, comment string) resource.TestCheckFunc {
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

func testAccPanosPanoramaAdministrativeTagDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_administrative_tag" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaAdministrativeTagId(rs.Primary.ID)
			_, err := pano.Objects.Tags.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaAdministrativeTagConfig(name, color, comment string) string {
	return fmt.Sprintf(`
resource "panos_panorama_administrative_tag" "test" {
    name = "%s"
    color = "%s"
    comment = "%s"
}
`, name, color, comment)
}
