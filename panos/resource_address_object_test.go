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


func TestPanosAddressObject_basic(t *testing.T) {
    var o addr.Entry
    name := fmt.Sprintf("tf%s", acctest.RandString(6))

    resource.Test(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        Providers: testAccProviders,
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
            return fmt.Errorf("Management profile label ID is not set")
        }

        fw := testAccProvider.Meta().(*pango.Firewall)
        vsys, name := parseAddressObjectId(rs.Primary.ID)
        v, err := fw.Objects.Address.Get(vsys, name)
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
    fw := testAccProvider.Meta().(*pango.Firewall)

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "panos_address_object" {
            continue
        }

        if rs.Primary.ID != "" {
            vsys, name := parseAddressObjectId(rs.Primary.ID)
            _, err := fw.Objects.Address.Get(vsys, name)
            if err == nil {
                return fmt.Errorf("Address object %q still exists", rs.Primary.ID)
            }
        }
        return nil
    }

    return nil
}

func testAccAddressObjectConfig(n, v, t, d string) string {
    return fmt.Sprintf(`
resource "panos_address_object" "test" {
    name = "%s"
    vsys = "vsys1"
    value = "%s"
    type = "%s"
    description = "%s"
}
`, n, v, t, d)
}
