package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/edl"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosEdl_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o edl.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosEdlDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEdlConfig(name, edl.TypeIp, "First desc", "https://first.paloaltonetworks.com", edl.RepeatEveryFiveMinutes, "10.0.0.0/8"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosEdlExists("panos_edl.test", &o),
					testAccCheckPanosEdlAttributes(&o, name, edl.TypeIp, "First desc", "https://first.paloaltonetworks.com", edl.RepeatEveryFiveMinutes, "10.0.0.0/8"),
				),
			},
			{
				Config: testAccEdlConfig(name, edl.TypeDomain, "Second desc", "https://second.paloaltonetworks.com", edl.RepeatHourly, "foobar.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosEdlExists("panos_edl.test", &o),
					testAccCheckPanosEdlAttributes(&o, name, edl.TypeDomain, "Second desc", "https://second.paloaltonetworks.com", edl.RepeatHourly, "foobar.com"),
				),
			},
		},
	})
}

func testAccCheckPanosEdlExists(n string, o *edl.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseEdlId(rs.Primary.ID)
		v, err := fw.Objects.Edl.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosEdlAttributes(o *edl.Entry, name, ty, desc, src, rpt, exc string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Type != ty {
			return fmt.Errorf("Type is %s, expected %s", o.Type, ty)
		}

		if o.Description != desc {
			return fmt.Errorf("Description is %s, expected %s", o.Description, desc)
		}

		if o.Source != src {
			return fmt.Errorf("Source is %q, expected %q", o.Source, src)
		}

		if o.Repeat != rpt {
			return fmt.Errorf("Repeat is %q, expected %q", o.Repeat, rpt)
		}

		if len(o.Exceptions) != 1 || o.Exceptions[0] != exc {
			return fmt.Errorf("Exceptions is %v, expected [\"%s\"]", o.Exceptions, exc)
		}

		return nil
	}
}

func testAccPanosEdlDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_edl" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseEdlId(rs.Primary.ID)
			_, err := fw.Objects.Edl.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccEdlConfig(name, ty, desc, src, rpt, exc string) string {
	return fmt.Sprintf(`
resource "panos_edl" "test" {
    name = %q
    type = %q
    description = %q
    source = %q
    repeat = %q
    exceptions = [%q]
}
`, name, ty, desc, src, rpt, exc)
}
