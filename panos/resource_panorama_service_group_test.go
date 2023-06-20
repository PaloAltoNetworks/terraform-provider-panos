package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/objs/srvcgrp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaServiceGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o srvcgrp.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	so1 := fmt.Sprintf("so%s", acctest.RandString(6))
	so2 := fmt.Sprintf("so%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaServiceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaServiceGroupConfig(so1, so2, name, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaServiceGroupExists("panos_panorama_service_group.test", &o),
					testAccCheckPanosPanoramaServiceGroupAttributes(&o, name, so1),
				),
			},
			{
				Config: testAccPanoramaServiceGroupConfig(so1, so2, name, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaServiceGroupExists("panos_panorama_service_group.test", &o),
					testAccCheckPanosPanoramaServiceGroupAttributes(&o, name, so2),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaServiceGroupExists(n string, o *srvcgrp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaServiceGroupId(rs.Primary.ID)
		v, err := pano.Objects.ServiceGroup.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaServiceGroupAttributes(o *srvcgrp.Entry, name, so string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.Services) != 1 || o.Services[0] != so {
			return fmt.Errorf("Services is %#v, expected [%s]", o.Services, so)
		}

		return nil
	}
}

func testAccPanosPanoramaServiceGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_service_group" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaServiceGroupId(rs.Primary.ID)
			_, err := pano.Objects.ServiceGroup.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaServiceGroupConfig(so1, so2, name string, sn int) string {
	return fmt.Sprintf(`
resource "panos_panorama_service_object" "so1" {
    name = "%s"
    source_port = 1111
    destination_port = 2222
    protocol = "tcp"
}

resource "panos_panorama_service_object" "so2" {
    name = "%s"
    source_port = 1111
    destination_port = 2222
    protocol = "tcp"
}

resource "panos_panorama_service_group" "test" {
    name = "%s"
    services = [panos_panorama_service_object.so%d.name]
}
`, so1, so2, name, sn)
}
