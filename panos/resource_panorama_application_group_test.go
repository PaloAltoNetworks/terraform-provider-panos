package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app/group"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaApplicationGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o group.Entry
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	g1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	g2 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaApplicationGroupConfig(dg, g1, g2, name, g1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationGroupExists("panos_panorama_application_group.test", &o),
					testAccCheckPanosPanoramaApplicationGroupAttributes(&o, name, g1),
				),
			},
			{
				Config: testAccPanoramaApplicationGroupConfig(dg, g1, g2, name, g2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaApplicationGroupExists("panos_panorama_application_group.test", &o),
					testAccCheckPanosPanoramaApplicationGroupAttributes(&o, name, g2),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaApplicationGroupExists(n string, o *group.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaApplicationGroupId(rs.Primary.ID)
		v, err := pano.Objects.AppGroup.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaApplicationGroupAttributes(o *group.Entry, name, grp string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.Applications) != 1 || o.Applications[0] != grp {
			return fmt.Errorf("Applications is %#v, expected [%s]", o.Applications, grp)
		}

		return nil
	}
}

func testAccPanosPanoramaApplicationGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_application_group" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaApplicationGroupId(rs.Primary.ID)
			_, err := pano.Objects.AppGroup.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaApplicationGroupConfig(dg, g1, g2, name, grp string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
    description = "app group test"
}

resource "panos_panorama_application_object" %q {
    device_group = panos_panorama_device_group.x.name
    name = %q
    description = "appgroup test"
    category = "media"
    subcategory = "gaming"
    technology = "client-server"
}

resource "panos_panorama_application_object" %q {
    device_group = panos_panorama_device_group.x.name
    name = %q
    description = "appgroup test"
    category = "media"
    subcategory = "gaming"
    technology = "client-server"
}

resource "panos_panorama_application_group" "test" {
    device_group = panos_panorama_device_group.x.name
    name = %q
    applications = [panos_panorama_application_object.%s.name]
}
`, dg, g1, g1, g2, g2, name, grp)
}
