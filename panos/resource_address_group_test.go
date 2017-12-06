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

func TestPanosAddressGroup_static(t *testing.T) {
	var o addrgrp.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	o1 := fmt.Sprintf("ao%s", acctest.RandString(6))
	o2 := fmt.Sprintf("ao%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressGroupStaticConfig(o1, o2, name, "first desc", 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAddressGroupExists("panos_address_group.test", &o),
					testAccCheckPanosAddressGroupAttributes(&o, name, "first desc", o1, ""),
				),
			},
			{
				Config: testAccAddressGroupStaticConfig(o1, o2, name, "second desc", 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAddressGroupExists("panos_address_group.test", &o),
					testAccCheckPanosAddressGroupAttributes(&o, name, "second desc", o2, ""),
				),
			},
		},
	})
}

func TestPanosAddressGroup_dynamic(t *testing.T) {
	var o addrgrp.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressGroupDynamicConfig(name, "first desc", "jack", "and", "burton"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAddressGroupExists("panos_address_group.test", &o),
					testAccCheckPanosAddressGroupAttributes(&o, name, "first desc", "", buildDagString("jack", "and", "burton")),
				),
			},
			{
				Config: testAccAddressGroupDynamicConfig(name, "second desc", "foo", "or", "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAddressGroupExists("panos_address_group.test", &o),
					testAccCheckPanosAddressGroupAttributes(&o, name, "second desc", "", buildDagString("foo", "or", "bar")),
				),
			},
		},
	})
}

func testAccCheckPanosAddressGroupExists(n string, o *addrgrp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseAddressGroupId(rs.Primary.ID)
		v, err := fw.Objects.AddressGroup.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAddressGroupAttributes(o *addrgrp.Entry, name, desc, sv, dv string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if sv == "" {
			if o.Static != nil {
				return fmt.Errorf("Static is %#v, expected nil", o.Static)
			}
		} else {
			if len(o.Static) != 1 || o.Static[0] != sv {
				return fmt.Errorf("Static is %#v, expected [%s]", o.Static, sv)
			}
		}

		if o.Dynamic != dv {
			return fmt.Errorf("Dynamic is %q, expected %q", o.Dynamic, dv)
		}

		return nil
	}
}

func testAccPanosAddressGroupDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_address_group" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseAddressGroupId(rs.Primary.ID)
			_, err := fw.Objects.AddressGroup.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func buildDagString(t1, op, t2 string) string {
	return fmt.Sprintf("'%s' %s '%s'", t1, op, t2)
}

func testAccAddressGroupStaticConfig(o1, o2, name, desc string, sv int) string {
	return fmt.Sprintf(`
resource "panos_address_object" "o1" {
    name = "%s"
    value = "10.20.30.0/24"
}

resource "panos_address_object" "o2" {
    name = "%s"
    value = "10.25.35.0/24"
}

resource "panos_address_group" "test" {
    name = "%s"
    description = "%s"
    static = ["${panos_address_object.o%d.name}"]
}
`, o1, o2, name, desc, sv)
}

func testAccAddressGroupDynamicConfig(name, desc, t1, op, t2 string) string {
	return fmt.Sprintf(`
resource "panos_address_group" "test" {
    name = "%s"
    description = "%s"
    dynamic = "%s"
}
`, name, desc, buildDagString(t1, op, t2))
}
