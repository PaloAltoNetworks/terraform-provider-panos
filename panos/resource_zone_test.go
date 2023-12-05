package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/zone"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosZone_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o zone.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccZoneConfig(name, "10.1.1.0/24", "10.1.2.0/24", "192.168.1.0/24", "192.168.2.0/24", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosZoneExists("panos_zone.test", &o),
					testAccCheckPanosZoneAttributes(&o, name, "10.1.1.0/24", "10.1.2.0/24", "192.168.1.0/24", "192.168.2.0/24", true),
				),
			},
			{
				Config: testAccZoneConfig(name, "192.168.3.0/24", "192.168.4.0/24", "10.1.3.0/24", "10.1.4.0/24", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosZoneExists("panos_zone.test", &o),
					testAccCheckPanosZoneAttributes(&o, name, "192.168.3.0/24", "192.168.4.0/24", "10.1.3.0/24", "10.1.4.0/24", false),
				),
			},
		},
	})
}

func testAccCheckPanosZoneExists(n string, o *zone.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, _, vsys, name := parseZoneId(rs.Primary.ID)
		v, err := fw.Network.Zone.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosZoneAttributes(o *zone.Entry, n, ia1, ia2, ea1, ea2 string, eui bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %s, expected %s", o.Name, n)
		}

		if len(o.IncludeAcls) != 2 || o.IncludeAcls[0] != ia1 || o.IncludeAcls[1] != ia2 {
			return fmt.Errorf("Include ACLs is %v, expected [%s, %s]", o.IncludeAcls, ia1, ia2)
		}

		if len(o.ExcludeAcls) != 2 || o.ExcludeAcls[0] != ea1 || o.ExcludeAcls[1] != ea2 {
			return fmt.Errorf("Exclude ACLs is %v, expected [%s, %s]", o.ExcludeAcls, ea1, ea2)
		}

		if o.EnableUserId != eui {
			return fmt.Errorf("Enable User Id is %t, expected %t", o.EnableUserId, eui)
		}

		return nil
	}
}

func testAccPanosZoneDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_zone" {
			continue
		}

		if rs.Primary.ID != "" {
			_, _, vsys, name := parseZoneId(rs.Primary.ID)
			_, err := fw.Network.Zone.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccZoneConfig(n, ia1, ia2, ea1, ea2 string, eui bool) string {
	return fmt.Sprintf(`
resource "panos_zone" "test" {
    name = "%s"
    vsys = "vsys1"
    mode = "layer3"
    include_acls = ["%s", "%s"]
    exclude_acls = ["%s", "%s"]
    enable_user_id = %t
}
`, n, ia1, ia2, ea1, ea2, eui)
}
