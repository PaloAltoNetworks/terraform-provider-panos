package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/srvc"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestPanosServiceObject_basic(t *testing.T) {
	var o srvc.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosServiceObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceObjectConfig(name, "description one", "2000-5000", "5432", "tcp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosServiceObjectExists("panos_service_object.test", &o),
					testAccCheckPanosServiceObjectAttributes(&o, name, "description one", "2000-5000", "5432", "tcp"),
				),
			},
			{
				Config: testAccServiceObjectConfig(name, "description two", "1025-65535", "12345", "udp"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosServiceObjectExists("panos_service_object.test", &o),
					testAccCheckPanosServiceObjectAttributes(&o, name, "description two", "1025-65535", "12345", "udp"),
				),
			},
		},
	})
}

func testAccCheckPanosServiceObjectExists(n string, o *srvc.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseServiceObjectId(rs.Primary.ID)
		v, err := fw.Objects.Services.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosServiceObjectAttributes(o *srvc.Entry, n, desc, sp, dp, proto string) resource.TestCheckFunc {
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

func testAccPanosServiceObjectDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_service_object" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseServiceObjectId(rs.Primary.ID)
			_, err := fw.Objects.Services.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccServiceObjectConfig(n, desc, sp, dp, proto string) string {
	return fmt.Sprintf(`
resource "panos_service_object" "test" {
    name = "%s"
    vsys = "vsys1"
    description = "%s"
    source_port = "%s"
    destination_port = "%s"
    protocol = "%s"
}
`, n, desc, sp, dp, proto)
}
