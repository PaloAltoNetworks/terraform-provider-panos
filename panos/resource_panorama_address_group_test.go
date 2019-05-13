package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addrgrp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaAddressGroup_static(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o addrgrp.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	o1 := fmt.Sprintf("ao%s", acctest.RandString(6))
	o2 := fmt.Sprintf("ao%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaAddressGroupStaticConfig(o1, o2, name, "first desc", 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAddressGroupExists("panos_panorama_address_group.test", &o),
					testAccCheckPanosPanoramaAddressGroupAttributes(&o, name, "first desc", o1, ""),
				),
			},
			{
				Config: testAccPanoramaAddressGroupStaticConfig(o1, o2, name, "second desc", 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAddressGroupExists("panos_panorama_address_group.test", &o),
					testAccCheckPanosPanoramaAddressGroupAttributes(&o, name, "second desc", o2, ""),
				),
			},
		},
	})
}

func TestAccPanosPanoramaAddressGroup_dynamic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o addrgrp.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaAddressGroupDynamicConfig(name, "first desc", "jack", "and", "burton"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAddressGroupExists("panos_panorama_address_group.test", &o),
					testAccCheckPanosPanoramaAddressGroupAttributes(&o, name, "first desc", "", buildPanoramaDagString("jack", "and", "burton")),
				),
			},
			{
				Config: testAccPanoramaAddressGroupDynamicConfig(name, "second desc", "foo", "or", "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAddressGroupExists("panos_panorama_address_group.test", &o),
					testAccCheckPanosPanoramaAddressGroupAttributes(&o, name, "second desc", "", buildPanoramaDagString("foo", "or", "bar")),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaAddressGroupExists(n string, o *addrgrp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		dg, name := parsePanoramaAddressGroupId(rs.Primary.ID)
		v, err := pano.Objects.AddressGroup.Get(dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaAddressGroupAttributes(o *addrgrp.Entry, name, desc, sv, dv string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if sv == "" {
			if o.StaticAddresses != nil {
				return fmt.Errorf("StaticAddresses is %#v, expected nil", o.StaticAddresses)
			}
		} else {
			if len(o.StaticAddresses) != 1 || o.StaticAddresses[0] != sv {
				return fmt.Errorf("StaticAddresses is %#v, expected [%s]", o.StaticAddresses, sv)
			}
		}

		if o.DynamicMatch != dv {
			return fmt.Errorf("DynamicMatch is %q, expected %q", o.DynamicMatch, dv)
		}

		return nil
	}
}

func testAccPanosPanoramaAddressGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_address_group" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, name := parsePanoramaAddressGroupId(rs.Primary.ID)
			_, err := pano.Objects.AddressGroup.Get(dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func buildPanoramaDagString(t1, op, t2 string) string {
	return fmt.Sprintf("'%s' %s '%s'", t1, op, t2)
}

func testAccPanoramaAddressGroupStaticConfig(o1, o2, name, desc string, sv int) string {
	return fmt.Sprintf(`
resource "panos_panorama_address_object" "o1" {
    name = %q
    value = "10.20.30.0/24"
}

resource "panos_panorama_address_object" "o2" {
    name = %q
    value = "10.25.35.0/24"
}

resource "panos_panorama_address_group" "test" {
    name = "%s"
    description = "%s"
    static_addresses = [panos_panorama_address_object.o%d.name]
}
`, o1, o2, name, desc, sv)
}

func testAccPanoramaAddressGroupDynamicConfig(name, desc, t1, op, t2 string) string {
	return fmt.Sprintf(`
resource "panos_panorama_address_group" "test" {
    name = %q
    description = %q
    dynamic_match = %q
}
`, name, desc, buildPanoramaDagString(t1, op, t2))
}
