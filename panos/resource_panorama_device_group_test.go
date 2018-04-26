package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/dg"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosDeviceGroup_basic(t *testing.T) {
	var o dg.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosDeviceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupConfig(name, "first description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDeviceGroupExists("panos_device_group.test", &o),
					testAccCheckPanosDeviceGroupAttributes(&o, "first description"),
				),
			},
			{
				Config: testAccDeviceGroupConfig(name, "second description"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDeviceGroupExists("panos_device_group.test", &o),
					testAccCheckPanosDeviceGroupAttributes(&o, "second description"),
				),
			},
		},
	})
}

func testAccCheckPanosDeviceGroupExists(n string, o *dg.Entry) resource.TestCheckFunc {
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
		v, err := pano.Panorama.DeviceGroup.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosDeviceGroupAttributes(o *dg.Entry, desc string) resource.TestCheckFunc {
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

func testAccPanosDeviceGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_device_group" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := pano.Panorama.DeviceGroup.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccDeviceGroupConfig(n, d string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "test" {
    name = "%s"
    description = "%s"
}
`, n, d)
}
