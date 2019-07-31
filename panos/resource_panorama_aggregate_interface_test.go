package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	agg "github.com/PaloAltoNetworks/pango/netw/interface/aggregate"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaAggregateInterface_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	} else if !testAccSupportsAggregateInterfaces {
		t.Skip(SkipAggregateTest)
	}

	var o agg.Entry
	num := (acctest.RandInt() % 5) + 1
	name := fmt.Sprintf("ae%d", num)
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(8))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaAggregateInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaAggregateInterfaceConfig(tmpl, name, "10.5.5.1/24", "foo"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAggregateInterfaceExists("panos_panorama_aggregate_interface.test", &o),
					testAccCheckPanosPanoramaAggregateInterfaceAttributes(&o, name, "10.5.5.1/24", "foo"),
				),
			},
			{
				Config: testAccPanoramaAggregateInterfaceConfig(tmpl, name, "10.6.6.1/24", "bar"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaAggregateInterfaceExists("panos_panorama_aggregate_interface.test", &o),
					testAccCheckPanosPanoramaAggregateInterfaceAttributes(&o, name, "10.6.6.1/24", "bar"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaAggregateInterfaceExists(n string, o *agg.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, _, name := parsePanoramaAggregateInterfaceId(rs.Primary.ID)
		v, err := pano.Network.AggregateInterface.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaAggregateInterfaceAttributes(o *agg.Entry, name, sip, com string) resource.TestCheckFunc {
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

func testAccPanosPanoramaAggregateInterfaceDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_aggregate_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, _, name := parsePanoramaAggregateInterfaceId(rs.Primary.ID)
			_, err := pano.Network.AggregateInterface.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaAggregateInterfaceConfig(tmpl, name, sip, com string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "for aggregate interface"
}

resource "panos_panorama_aggregate_interface" "test" {
    name = %q
    template = panos_panorama_template.x.name
    mode = "layer3"
    static_ips = [%q]
    comment = %q
}
`, tmpl, name, sip, com)
}
