package panos

import (
	"fmt"
	"os"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/pnrm/dg"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosDeviceGroupEntry_basic(t *testing.T) {
	/*
	   In order to run this test you'll need:

	   * panorama as the device (PANOS_HOSTNAME, PANOS_USERNAME, PANOS_PASSWORD)
	   * a multi-vsys firewall that's already added as a managed device
	     (PANOS_MANAGED_SERIAL_NUMBER)
	   * and two vsys that are currently unassigned to any device group
	     (PANOS_MANAGED_VSYS1 and PANOS_MANAGED_VSYS2)
	*/

	serial := os.Getenv("PANOS_MANAGED_SERIAL_NUMBER")
	vsys1 := os.Getenv("PANOS_MANAGED_VSYS1")
	vsys2 := os.Getenv("PANOS_MANAGED_VSYS2")

	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if serial == "" || vsys1 == "" || vsys2 == "" {
		t.Skip("One or more required env variables are unset (PANOS_MANAGED_SERIAL_NUMBER, PANOS_MANAGED_VSYS1, PANOS_MANAGED_VSYS2")
	}

	var o dg.Entry
	group := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosDeviceGroupEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupEntryConfig(group, serial, vsys1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDeviceGroupEntryExists("panos_panorama_device_group_entry.test", &o),
					testAccCheckPanosDeviceGroupEntryAttributes(&o, serial, vsys1),
				),
			},
			{
				Config: testAccDeviceGroupEntryConfig(group, serial, vsys2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDeviceGroupEntryExists("panos_panorama_device_group_entry.test", &o),
					testAccCheckPanosDeviceGroupEntryAttributes(&o, serial, vsys2),
				),
			},
		},
	})
}

func testAccCheckPanosDeviceGroupEntryExists(n string, o *dg.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		group, serial := parseDeviceGroupEntryId(rs.Primary.ID)
		v, err := pano.Panorama.DeviceGroup.Get(group)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		} else if _, ok = v.Devices[serial]; !ok {
			return fmt.Errorf("Serial %s not in devices list", serial)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosDeviceGroupEntryAttributes(o *dg.Entry, serial, vsys string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(o.Devices[serial]) != 1 || o.Devices[serial][0] != vsys {
			return fmt.Errorf("Vsys list for serial %q is %#v, not [%s]", serial, o.Devices[serial], vsys)
		}
		return nil
	}
}

func testAccPanosDeviceGroupEntryDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_device_group_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			group, serial := parseDeviceGroupEntryId(rs.Primary.ID)
			o, err := pano.Panorama.DeviceGroup.Get(group)
			if err == nil {
				if _, ok := o.Devices[serial]; ok {
					return fmt.Errorf("Object %q still exists", rs.Primary.ID)
				}
			}
		}
		return nil
	}

	return nil
}

func testAccDeviceGroupEntryConfig(group, serial, vsys string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "dg" {
    name = %q
    description = "for device group entry test"
}

resource "panos_panorama_device_group_entry" "test" {
    device_group = panos_panorama_device_group.dg.name
    serial = %q
    vsys_list = [%q]
}
`, group, serial, vsys)
}
