package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/app/group"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosApplicationGroup_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o group.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	g1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	g2 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosApplicationGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationGroupConfig(g1, g2, name, g1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosApplicationGroupExists("panos_application_group.test", &o),
					testAccCheckPanosApplicationGroupAttributes(&o, name, g1),
				),
			},
			{
				Config: testAccApplicationGroupConfig(g1, g2, name, g2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosApplicationGroupExists("panos_application_group.test", &o),
					testAccCheckPanosApplicationGroupAttributes(&o, name, g2),
				),
			},
		},
	})
}

func testAccCheckPanosApplicationGroupExists(n string, o *group.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseApplicationGroupId(rs.Primary.ID)
		v, err := fw.Objects.AppGroup.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosApplicationGroupAttributes(o *group.Entry, name, grp string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.Applications) != 1 || o.Applications[0] != grp {
			return fmt.Errorf("Applications is %#v, expected [%s]", o.Applications, grp)
		}

		return nil
	}
}

func testAccPanosApplicationGroupDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_application_group" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseApplicationGroupId(rs.Primary.ID)
			_, err := fw.Objects.AppGroup.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccApplicationGroupConfig(g1, g2, name, grp string) string {
	return fmt.Sprintf(`
resource "panos_application_object" %q {
    name = %q
    description = "appgroup test"
    category = "media"
    subcategory = "gaming"
    technology = "client-server"
}

resource "panos_application_object" %q {
    name = %q
    description = "appgroup test"
    category = "media"
    subcategory = "gaming"
    technology = "client-server"
}

resource "panos_application_group" "test" {
    name = %q
    applications = [panos_application_object.%s.name]
}
`, g1, g1, g2, g2, name, grp)
}
