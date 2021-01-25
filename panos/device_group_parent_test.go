package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source test.
func TestAccPanosDsDeviceGroupParent(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	parent := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsDeviceGroupParentConfig(name, parent),
				Check: checkDataSource("panos_device_group_parent", []string{
					"total",
				}),
			},
		},
	})
}

func testAccDsDeviceGroupParentConfig(name, parent string) string {
	return fmt.Sprintf(`
data "panos_device_group_parent" "test" {}

resource "panos_device_group_parent" "x" {
    device_group = panos_panorama_device_group.dg.name
    parent = panos_panorama_device_group.parent.name
}

resource "panos_panorama_device_group" "dg" {
    name = %q
}

resource "panos_panorama_device_group" "parent" {
    name = %q
}
`, name, parent)
}

// Resource test.
func TestAccPanosDeviceGroupParent(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o map[string]string
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	p1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	p2 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosDeviceGroupParentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupParentConfig(name, p1, p1, p2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDeviceGroupParentExists("panos_device_group_parent.test", &o),
					testAccCheckPanosDeviceGroupParentAttributes(&o, name, p1),
				),
			},
			{
				Config: testAccDeviceGroupParentConfig(name, p2, p1, p2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosDeviceGroupParentExists("panos_device_group_parent.test", &o),
					testAccCheckPanosDeviceGroupParentAttributes(&o, name, p2),
				),
			},
		},
	})
}

func testAccCheckPanosDeviceGroupParentExists(n string, o *map[string]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v map[string]string

		switch con := testAccProvider.Meta().(type) {
		case *pango.Panorama:
			v, err = con.Panorama.DeviceGroup.GetParents()
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosDeviceGroupParentAttributes(o *map[string]string, name, parent string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if (*o)[name] != parent {
			return fmt.Errorf("Expected %q to have parent of %q, has %q", name, parent, (*o)[name])
		}

		return nil
	}
}

func testAccPanosDeviceGroupParentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_device_group_parent" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error
			var name string
			var info map[string]string

			switch con := testAccProvider.Meta().(type) {
			case *pango.Panorama:
				name = rs.Primary.ID
				info, err = con.Panorama.DeviceGroup.GetParents()
				if err != nil {
					return err
				}
			}
			if info[name] != "" {
				return fmt.Errorf("Device group %q still has a parent: %s", name, info[name])
			}
		}
		return nil
	}

	return nil
}

func testAccDeviceGroupParentConfig(name, parent, p1, p2 string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
}

resource "panos_panorama_device_group" %q {
    name = %q
}

resource "panos_panorama_device_group" %q {
    name = %q
}

resource "panos_device_group_parent" "test" {
    device_group = panos_panorama_device_group.x.name
    parent = panos_panorama_device_group.%s.name
}
`, name, p1, p1, p2, p2, parent)
}
