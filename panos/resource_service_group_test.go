package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosServiceGroup_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o srvcgrp.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	so1 := fmt.Sprintf("so%s", acctest.RandString(6))
	so2 := fmt.Sprintf("so%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosServiceGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceGroupConfig(so1, so2, name, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosServiceGroupExists("panos_service_group.test", &o),
					testAccCheckPanosServiceGroupAttributes(&o, name, so1),
				),
			},
			{
				Config: testAccServiceGroupConfig(so1, so2, name, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosServiceGroupExists("panos_service_group.test", &o),
					testAccCheckPanosServiceGroupAttributes(&o, name, so2),
				),
			},
		},
	})
}

func testAccCheckPanosServiceGroupExists(n string, o *srvcgrp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseServiceGroupId(rs.Primary.ID)
		v, err := fw.Objects.ServiceGroup.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosServiceGroupAttributes(o *srvcgrp.Entry, name, so string) resource.TestCheckFunc {
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

func testAccPanosServiceGroupDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_service_group" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseServiceGroupId(rs.Primary.ID)
			_, err := fw.Objects.ServiceGroup.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccServiceGroupConfig(so1, so2, name string, sn int) string {
	return fmt.Sprintf(`
resource "panos_service_object" "so1" {
    name = "%s"
    source_port = 1111
    destination_port = 2222
    protocol = "tcp"
}

resource "panos_service_object" "so2" {
    name = "%s"
    source_port = 1111
    destination_port = 2222
    protocol = "tcp"
}

resource "panos_service_group" "test" {
    name = "%s"
    services = [panos_service_object.so%d.name]
}
`, so1, so2, name, sn)
}
