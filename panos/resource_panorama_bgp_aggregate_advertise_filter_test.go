package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate/filter/advertise"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaBgpAggregateAdvertiseFilter_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o advertise.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	ag := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpAggregateAdvertiseFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpAggregateAdvertiseFilterConfig(tmpl, vr, ag, name, "path1", "com1", "ecom1", "21", "192.168.5.0/24", "10.1.1.0/24", true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpAggregateAdvertiseFilterExists("panos_panorama_bgp_aggregate_advertise_filter.test", &o),
					testAccCheckPanosPanoramaBgpAggregateAdvertiseFilterAttributes(&o, name, "path1", "com1", "ecom1", "21", "192.168.5.0/24", "10.1.1.0/24", true, false),
				),
			},
			{
				Config: testAccPanoramaBgpAggregateAdvertiseFilterConfig(tmpl, vr, ag, name, "path2", "com2", "ecom2", "22", "192.168.6.0/24", "10.1.2.0/24", false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpAggregateAdvertiseFilterExists("panos_panorama_bgp_aggregate_advertise_filter.test", &o),
					testAccCheckPanosPanoramaBgpAggregateAdvertiseFilterAttributes(&o, name, "path2", "com2", "ecom2", "22", "192.168.6.0/24", "10.1.2.0/24", false, true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpAggregateAdvertiseFilterExists(n string, o *advertise.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, ag, name := parsePanoramaBgpAggregateAdvertiseFilterId(rs.Primary.ID)
		v, err := pano.Network.BgpAggAdvertiseFilter.Get(tmpl, ts, vr, ag, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpAggregateAdvertiseFilterAttributes(o *advertise.Entry, name, apr, cr, ecr, med, nh, ap string, ex, en bool) resource.TestCheckFunc {
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

func testAccPanosPanoramaBgpAggregateAdvertiseFilterDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_aggregate_advertise_filter" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, ag, name := parsePanoramaBgpAggregateAdvertiseFilterId(rs.Primary.ID)
			_, err := pano.Network.BgpAggAdvertiseFilter.Get(tmpl, ts, vr, ag, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpAggregateAdvertiseFilterConfig(tmpl, vr, ag, name, apr, cr, ecr, med, nh, ap string, ex, en bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "t" {
    name = %q
}

resource "panos_panorama_virtual_router" "vr" {
    template = panos_panorama_template.t.name
    name = %q
}

resource "panos_panorama_bgp" "x" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "55"
    enable = false
}

resource "panos_panorama_bgp_aggregate" "x" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    name = %q
    prefix = "192.168.1.0/24"
}

resource "panos_panorama_bgp_aggregate_advertise_filter" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp_aggregate.x.virtual_router
    bgp_aggregate = panos_panorama_bgp_aggregate.x.name
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
`, tmpl, vr, ag, name, apr, cr, ecr, med, nh, ap, ex, en)
}
