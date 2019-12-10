package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaBgpAggregate_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o aggregate.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpAggregateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpAggregateConfig(tmpl, vr, name, "192.168.1.0/24", "123", "234", "192.168.5.5", aggregate.OriginIncomplete, aggregate.AsPathTypeNone, "", false, true, true, 21, 12),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpAggregateExists("panos_panorama_bgp_aggregate.test", &o),
					testAccCheckPanosPanoramaBgpAggregateAttributes(&o, name, "192.168.1.0/24", "123", "234", "192.168.5.5", aggregate.OriginIncomplete, aggregate.AsPathTypeNone, "", false, true, true, 21, 12),
				),
			},
			{
				Config: testAccPanoramaBgpAggregateConfig(tmpl, vr, name, "192.168.2.0/24", "321", "432", "192.168.6.6", aggregate.OriginEgp, aggregate.AsPathTypePrepend, "7", true, false, false, 73, 37),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpAggregateExists("panos_panorama_bgp_aggregate.test", &o),
					testAccCheckPanosPanoramaBgpAggregateAttributes(&o, name, "192.168.2.0/24", "321", "432", "192.168.6.6", aggregate.OriginEgp, aggregate.AsPathTypePrepend, "7", true, false, false, 73, 37),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpAggregateExists(n string, o *aggregate.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, name := parsePanoramaBgpAggregateId(rs.Primary.ID)
		v, err := pano.Network.BgpAggregate.Get(tmpl, ts, vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpAggregateAttributes(o *aggregate.Entry, name, pre, lp, med, nh, ori, apt, apv string, en, sum, as bool, wei, apl int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, expected %q", o.Name, name)
		}

		if o.Prefix != pre {
			return fmt.Errorf("Prefix is %q, not %q", o.Prefix, pre)
		}

		if o.LocalPreference != lp {
			return fmt.Errorf("Local preference is %q, not %q", o.LocalPreference, lp)
		}

		if o.Med != med {
			return fmt.Errorf("MED is %q, not %q", o.Med, med)
		}

		if o.NextHop != nh {
			return fmt.Errorf("Next hop is %q, not %q", o.NextHop, nh)
		}

		if o.Origin != ori {
			return fmt.Errorf("Origin is %q, not %q", o.Origin, ori)
		}

		if o.AsPathType != apt {
			return fmt.Errorf("AS path type is %q, not %q", o.AsPathType, apt)
		}

		if o.AsPathValue != apv {
			return fmt.Errorf("AS path value is %q, not %q", o.AsPathValue, apv)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.Summary != sum {
			return fmt.Errorf("Summary is %t, not %t", o.Summary, sum)
		}

		if o.AsSet != as {
			return fmt.Errorf("AS set is %t, not %t", o.AsSet, as)
		}

		if o.Weight != wei {
			return fmt.Errorf("Weight is %d, not %d", o.Weight, wei)
		}

		if o.AsPathLimit != apl {
			return fmt.Errorf("AS path limit is %d, not %d", o.AsPathLimit, apl)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpAggregateDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_aggregate" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, name := parsePanoramaBgpAggregateId(rs.Primary.ID)
			_, err := pano.Network.BgpAggregate.Get(tmpl, ts, vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpAggregateConfig(tmpl, vr, name, pre, lp, med, nh, ori, apt, apv string, en, sum, as bool, wei, apl int) string {
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

resource "panos_panorama_bgp_aggregate" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    name = %q
    prefix = %q
    local_preference = %q
    med = %q
    next_hop = %q
    origin = %q
    as_path_type = %q
    as_path_value = %q
    enable = %t
    summary = %t
    as_set = %t
    weight = %d
    as_path_limit = %d
}
`, tmpl, vr, name, pre, lp, med, nh, ori, apt, apv, en, sum, as, wei, apl)
}
