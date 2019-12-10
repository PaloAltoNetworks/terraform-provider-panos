package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvc"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaServiceObject_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o srvc.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaServiceObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaServiceObjectConfig(name, "description one", "2000-5000", "5432", "tcp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaServiceObjectExists("panos_panorama_service_object.test", &o),
					testAccCheckPanosPanoramaServiceObjectAttributes(&o, name, "description one", "2000-5000", "5432", "tcp"),
				),
			},
			{
				Config: testAccPanoramaServiceObjectConfig(name, "description two", "1025-65535", "12345", "udp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaServiceObjectExists("panos_panorama_service_object.test", &o),
					testAccCheckPanosPanoramaServiceObjectAttributes(&o, name, "description two", "1025-65535", "12345", "udp"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaServiceObjectExists(n string, o *srvc.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaServiceObjectId(rs.Primary.ID)
		v, err := pano.Objects.Services.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaServiceObjectAttributes(o *srvc.Entry, n, desc, sp, dp, proto string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %s, expected %s", o.Name, n)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if o.SourcePort != sp {
			return fmt.Errorf("Source port is %s, expected %s", o.SourcePort, sp)
		}

		if o.DestinationPort != dp {
			return fmt.Errorf("Destination port is %s, expected %s", o.DestinationPort, dp)
		}

		if o.Protocol != proto {
			return fmt.Errorf("Protocol is %s, expected %s", o.Protocol, proto)
		}

		return nil
	}
}

func testAccPanosPanoramaServiceObjectDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_service_object" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaServiceObjectId(rs.Primary.ID)
			_, err := pano.Objects.Services.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaServiceObjectConfig(n, desc, sp, dp, proto string) string {
	return fmt.Sprintf(`
resource "panos_panorama_service_object" "test" {
    name = "%s"
    description = "%s"
    source_port = "%s"
    destination_port = "%s"
    protocol = "%s"
}
`, n, desc, sp, dp, proto)
}
