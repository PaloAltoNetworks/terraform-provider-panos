package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	agg "github.com/fpluchorg/pango/netw/interface/aggregate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosAggregateInterface_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccSupportsAggregateInterfaces {
		t.Skip(SkipAggregateTest)
	}

	var o agg.Entry
	num := (acctest.RandInt() % 5) + 1
	name := fmt.Sprintf("ae%d", num)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAggregateInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateInterfaceConfig(name, "10.5.5.1/24", "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAggregateInterfaceExists("panos_aggregate_interface.test", &o),
					testAccCheckPanosAggregateInterfaceAttributes(&o, name, "10.5.5.1/24", "foo"),
				),
			},
			{
				Config: testAccAggregateInterfaceConfig(name, "10.6.6.1/24", "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAggregateInterfaceExists("panos_aggregate_interface.test", &o),
					testAccCheckPanosAggregateInterfaceAttributes(&o, name, "10.6.6.1/24", "bar"),
				),
			},
		},
	})
}

func testAccCheckPanosAggregateInterfaceExists(n string, o *agg.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		_, name := parseAggregateInterfaceId(rs.Primary.ID)
		v, err := fw.Network.AggregateInterface.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAggregateInterfaceAttributes(o *agg.Entry, name, sip, com string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.StaticIps) != 1 || o.StaticIps[0] != sip {
			return fmt.Errorf("Static IPs is %#v, not [%s]", o.StaticIps, sip)
		}

		if o.Comment != com {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, com)
		}

		return nil
	}
}

func testAccPanosAggregateInterfaceDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_aggregate_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			_, name := parseAggregateInterfaceId(rs.Primary.ID)
			_, err := fw.Network.AggregateInterface.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAggregateInterfaceConfig(name, sip, com string) string {
	return fmt.Sprintf(`
resource "panos_aggregate_interface" "test" {
    name = %q
    mode = "layer3"
    static_ips = [%q]
    comment = %q
}
`, name, sip, com)
}
