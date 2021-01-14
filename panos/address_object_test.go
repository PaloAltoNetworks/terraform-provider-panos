package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/addr"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsAddressObjectList(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAddressObjectConfig(name),
				Check:  checkDataSourceListing("panos_address_objects"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsAddressObject_basic(t *testing.T) {
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAddressObjectConfig(name),
				Check: checkDataSource("panos_address_object", []string{
					"name", "type", "value", "description",
				}),
			},
		},
	})
}

func testAccDsAddressObjectConfig(name string) string {
	return fmt.Sprintf(`
data "panos_address_objects" "test" {}

data "panos_address_object" "test" {
    name = panos_address_object.x.name
}

resource "panos_address_object" "x" {
    name = %q
    description = "address object acctest"
    value = "10.59.48.37"
}
`, name)
}

// Resource tests.
func TestAccPanosAddressObject_basic(t *testing.T) {
	var o addr.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAddressObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAddressObjectConfig(name, "10.1.1.1-10.1.1.250", "ip-range", "new desc"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAddressObjectExists("panos_address_object.test", &o),
					testAccCheckPanosAddressObjectAttributes(&o, name, "10.1.1.1-10.1.1.250", "ip-range", "new desc"),
				),
			},
			{
				Config: testAccAddressObjectConfig(name, "10.1.1.1", "ip-netmask", "foobar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAddressObjectExists("panos_address_object.test", &o),
					testAccCheckPanosAddressObjectAttributes(&o, name, "10.1.1.1", "ip-netmask", "foobar"),
				),
			},
		},
	})
}

func testAccCheckPanosAddressObjectExists(n string, o *addr.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v addr.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseAddressObjectId(rs.Primary.ID)
			v, err = con.Objects.Address.Get(vsys, name)
		case *pango.Panorama:
			dg, name := parseAddressObjectId(rs.Primary.ID)
			v, err = con.Objects.Address.Get(dg, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAddressObjectAttributes(o *addr.Entry, n, v, t, d string) resource.TestCheckFunc {
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

func testAccPanosAddressObjectDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_address_object" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseAddressObjectId(rs.Primary.ID)
				_, err = con.Objects.Address.Get(vsys, name)
			case *pango.Panorama:
				dg, name := parseAddressObjectId(rs.Primary.ID)
				_, err = con.Objects.Address.Get(dg, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAddressObjectConfig(n, v, t, d string) string {
	return fmt.Sprintf(`
resource "panos_address_object" "test" {
    name = %q
    value = %q
    type = %q
    description = %q
}
`, n, v, t, d)
}
