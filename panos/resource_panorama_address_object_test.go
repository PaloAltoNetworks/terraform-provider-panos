package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaAddressObject_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o addr.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaAddressObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaAddressObjectConfig(name, "10.1.1.1-10.1.1.250", "ip-range", "new desc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAddressObjectExists("panos_panorama_address_object.test", &o),
					testAccCheckPanosPanoramaAddressObjectAttributes(&o, name, "10.1.1.1-10.1.1.250", "ip-range", "new desc"),
				),
			},
			{
				Config: testAccPanoramaAddressObjectConfig(name, "10.1.1.1", "ip-netmask", "foobar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAddressObjectExists("panos_panorama_address_object.test", &o),
					testAccCheckPanosPanoramaAddressObjectAttributes(&o, name, "10.1.1.1", "ip-netmask", "foobar"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaAddressObjectExists(n string, o *addr.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaAddressObjectId(rs.Primary.ID)
		v, err := pano.Objects.Address.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaAddressObjectAttributes(o *addr.Entry, n, v, t, d string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %s, expected %s", o.Name, n)
		}

		if o.Value != v {
			return fmt.Errorf("Value is %s, expected %s", o.Value, v)
		}

		if o.Type != t {
			return fmt.Errorf("Type is %s, expected %s", o.Type, t)
		}

		if o.Description != d {
			return fmt.Errorf("Description is %s, expected %s", o.Description, d)
		}

		return nil
	}
}

func testAccPanosPanoramaAddressObjectDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_address_object" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaAddressObjectId(rs.Primary.ID)
			_, err := pano.Objects.Address.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaAddressObjectConfig(n, v, t, d string) string {
	return fmt.Sprintf(`
resource "panos_panorama_address_object" "test" {
    name = "%s"
    device_group = "shared"
    value = "%s"
    type = "%s"
    description = "%s"
}
`, n, v, t, d)
}
