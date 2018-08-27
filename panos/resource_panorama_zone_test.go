package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaZone_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o zone.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	ts := fmt.Sprintf("tfStack%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanosPanoramaZoneConfig(ts, name, "10.1.1.0/24", "192.168.1.0/24", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaZoneExists("panos_panorama_zone.test", &o),
					testAccCheckPanosPanoramaZoneAttributes(&o, name, "10.1.1.0/24", "192.168.1.0/24", true),
				),
			},
			{
				Config: testAccPanosPanoramaZoneConfig(ts, name, "192.168.3.0/24", "10.1.3.0/24", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaZoneExists("panos_panorama_zone.test", &o),
					testAccCheckPanosPanoramaZoneAttributes(&o, name, "192.168.3.0/24", "10.1.3.0/24", false),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaZoneExists(n string, o *zone.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vsys, name := parsePanoramaZoneId(rs.Primary.ID)
		v, err := pano.Network.Zone.Get(tmpl, ts, vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaZoneAttributes(o *zone.Entry, name, inc, exc string, eui bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.IncludeAcls) != 1 || o.IncludeAcls[0] != inc {
			return fmt.Errorf("Include ACLs is %v, expected [%s]", o.IncludeAcls, inc)
		}

		if len(o.ExcludeAcls) != 1 || o.ExcludeAcls[0] != exc {
			return fmt.Errorf("Exclude ACLs is %v, expected [%s]", o.ExcludeAcls, exc)
		}

		if o.EnableUserId != eui {
			return fmt.Errorf("Enable User Id is %t, expected %t", o.EnableUserId, eui)
		}

		return nil
	}
}

func testAccPanosPanoramaZoneDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_zone" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vsys, name := parsePanoramaZoneId(rs.Primary.ID)
			_, err := pano.Network.Zone.Get(tmpl, ts, vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanosPanoramaZoneConfig(ts, name, inc, exc string, eui bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_template_stack" "x" {
    name = %q
}

resource "panos_panorama_zone" "test" {
    template_stack = "${panos_panorama_template_stack.x.name}"
    name = %q
    mode = "layer3"
    include_acls = [%q]
    exclude_acls = [%q]
    enable_user_id = %t
}
`, ts, name, inc, exc, eui)
}
