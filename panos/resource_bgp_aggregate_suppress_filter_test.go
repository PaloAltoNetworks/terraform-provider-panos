package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/aggregate/filter/suppress"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosBgpAggregateSuppressFilter_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o suppress.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	ag := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpAggregateSuppressFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpAggregateSuppressFilterConfig(vr, ag, name, "path1", "com1", "ecom1", "21", "192.168.5.0/24", "10.1.1.0/24", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpAggregateSuppressFilterExists("panos_bgp_aggregate_suppress_filter.test", &o),
					testAccCheckPanosBgpAggregateSuppressFilterAttributes(&o, name, "path1", "com1", "ecom1", "21", "192.168.5.0/24", "10.1.1.0/24", true, false),
				),
			},
			{
				Config: testAccBgpAggregateSuppressFilterConfig(vr, ag, name, "path2", "com2", "ecom2", "22", "192.168.6.0/24", "10.1.2.0/24", false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpAggregateSuppressFilterExists("panos_bgp_aggregate_suppress_filter.test", &o),
					testAccCheckPanosBgpAggregateSuppressFilterAttributes(&o, name, "path2", "com2", "ecom2", "22", "192.168.6.0/24", "10.1.2.0/24", false, true),
				),
			},
		},
	})
}

func testAccCheckPanosBgpAggregateSuppressFilterExists(n string, o *suppress.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, ag, name := parseBgpAggregateSuppressFilterId(rs.Primary.ID)
		v, err := fw.Network.BgpAggSuppressFilter.Get(vr, ag, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpAggregateSuppressFilterAttributes(o *suppress.Entry, name, apr, cr, ecr, med, nh, ap string, ex, en bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, expected %q", o.Name, name)
		}

		if o.AsPathRegex != apr {
			return fmt.Errorf("AS path regex is %q, not %q", o.AsPathRegex, apr)
		}

		if o.CommunityRegex != cr {
			return fmt.Errorf("Community regex is %q, not %q", o.CommunityRegex, cr)
		}
		if o.ExtendedCommunityRegex != ecr {
			return fmt.Errorf("Extended community regex is %q, not %q", o.ExtendedCommunityRegex, ecr)
		}

		if o.Med != med {
			return fmt.Errorf("MED is %q, not %q", o.Med, med)
		}

		if len(o.NextHop) != 1 || o.NextHop[0] != nh {
			return fmt.Errorf("Next hop is %#v, not [%q]", o.NextHop, nh)
		}

		if len(o.AddressPrefix) != 1 {
			return fmt.Errorf("Address prefix should be len of 1, is %d", len(o.AddressPrefix))
		}

		v, ok := o.AddressPrefix[ap]
		if !ok {
			return fmt.Errorf("Address prefix %q is not present", ap)
		} else if v != ex {
			return fmt.Errorf("Address prefix exact is %t, not %t", v, ex)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		return nil
	}
}

func testAccPanosBgpAggregateSuppressFilterDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_aggregate_suppress_filter" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, ag, name := parseBgpAggregateSuppressFilterId(rs.Primary.ID)
			_, err := fw.Network.BgpAggSuppressFilter.Get(vr, ag, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpAggregateSuppressFilterConfig(vr, ag, name, apr, cr, ecr, med, nh, ap string, ex, en bool) string {
	return fmt.Sprintf(`
resource "panos_virtual_router" "vr" {
    name = %q
}

resource "panos_bgp" "x" {
    virtual_router = panos_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "55"
    enable = false
}

resource "panos_bgp_aggregate" "x" {
    virtual_router = panos_bgp.x.virtual_router
    name = %q
    prefix = "192.168.1.0/24"
}

resource "panos_bgp_aggregate_suppress_filter" "test" {
    virtual_router = panos_bgp_aggregate.x.virtual_router
    bgp_aggregate = panos_bgp_aggregate.x.name
    name = %q
    as_path_regex = %q
    community_regex = %q
    extended_community_regex = %q
    med = %q
    next_hops = [%q]
    address_prefix {
        prefix = %q
        exact = %t
    }
    enable = %t
}
`, vr, ag, name, apr, cr, ecr, med, nh, ap, ex, en)
}
