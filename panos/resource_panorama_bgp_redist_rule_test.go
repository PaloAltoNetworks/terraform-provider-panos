package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/redist"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaBgpRedistRule_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o redist.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("192.168.%d.0/24", acctest.RandInt()%240+5)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpRedistRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpRedistRuleConfig(tmpl, vr, name, redist.SetOriginIgp, "98", "76", "54:32", 21, 12, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpRedistRuleExists("panos_panorama_bgp_redist_rule.test", &o),
					testAccCheckPanosPanoramaBgpRedistRuleAttributes(&o, name, redist.SetOriginIgp, "98", "76", "54:32", 21, 12, false),
				),
			},
			{
				Config: testAccPanoramaBgpRedistRuleConfig(tmpl, vr, name, redist.SetOriginIncomplete, "76", "98", "32:54", 12, 21, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpRedistRuleExists("panos_panorama_bgp_redist_rule.test", &o),
					testAccCheckPanosPanoramaBgpRedistRuleAttributes(&o, name, redist.SetOriginIncomplete, "76", "98", "32:54", 12, 21, true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpRedistRuleExists(n string, o *redist.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, name := parsePanoramaBgpRedistRuleId(rs.Primary.ID)
		v, err := pano.Network.BgpRedistRule.Get(tmpl, ts, vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpRedistRuleAttributes(o *redist.Entry, name, ori, med, lp, com string, met, sapl int, en bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, expected %q", o.Name, name)
		}

		if o.SetOrigin != ori {
			return fmt.Errorf("Set origin is %q, not %q", o.SetOrigin, ori)
		}

		if o.SetMed != med {
			return fmt.Errorf("Set MED is %q, not %q", o.SetMed, med)
		}

		if o.SetLocalPreference != lp {
			return fmt.Errorf("Set local preference is %q, not %q", o.SetLocalPreference, lp)
		}

		if len(o.SetCommunity) != 1 || o.SetCommunity[0] != com {
			return fmt.Errorf("Set communities is %#v, not [%q]", o.SetCommunity, com)
		}

		if o.Metric != met {
			return fmt.Errorf("Metric is %d, not %d", o.Metric, met)
		}

		if o.SetAsPathLimit != sapl {
			return fmt.Errorf("Set AS path limit is %d, not %d", o.SetAsPathLimit, sapl)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpRedistRuleDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_redist_rule" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, name := parsePanoramaBgpRedistRuleId(rs.Primary.ID)
			_, err := pano.Network.BgpRedistRule.Get(tmpl, ts, vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpRedistRuleConfig(tmpl, vr, name, ori, med, lp, com string, met, sapl int, en bool) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

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

resource "panos_panorama_bgp_redist_rule" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    address_family = "ipv4"
    route_table = data.panos_system_info.x.version_major >= 8 ? "unicast" : ""
    name = %q
    set_origin = %q
    set_med = %q
    set_local_preference = %q
    set_communities = [%q]
    metric = %d
    set_as_path_limit = %d
    enable = %t
}
`, tmpl, vr, name, ori, med, lp, com, met, sapl, en)
}
