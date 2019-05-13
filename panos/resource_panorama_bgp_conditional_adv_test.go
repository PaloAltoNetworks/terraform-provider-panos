package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaBgpConditionalAdv_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o conadv.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpConditionalAdvDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpConditionalAdvConfig(tmpl, vr, name, "jack", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpConditionalAdvExists("panos_panorama_bgp_conditional_adv.test", &o),
					testAccCheckPanosPanoramaBgpConditionalAdvAttributes(&o, name, "jack", false),
				),
			},
			{
				Config: testAccPanoramaBgpConditionalAdvConfig(tmpl, vr, name, "wang", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpConditionalAdvExists("panos_panorama_bgp_conditional_adv.test", &o),
					testAccCheckPanosPanoramaBgpConditionalAdvAttributes(&o, name, "wang", true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpConditionalAdvExists(n string, o *conadv.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, name := parsePanoramaBgpConditionalAdvId(rs.Primary.ID)
		v, err := pano.Network.BgpConditionalAdv.Get(tmpl, ts, vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpConditionalAdvAttributes(o *conadv.Entry, name, ub string, en bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, expected %q", o.Name, name)
		}

		if len(o.UsedBy) != 1 || o.UsedBy[0] != ub {
			return fmt.Errorf("Used by is %#v, not [%q]", o.UsedBy, ub)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpConditionalAdvDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_conditional_adv" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, name := parsePanoramaBgpConditionalAdvId(rs.Primary.ID)
			_, err := pano.Network.BgpConditionalAdv.Get(tmpl, ts, vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpConditionalAdvConfig(tmpl, vr, name, ub string, en bool) string {
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

resource "panos_panorama_bgp_peer_group" "jack" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    name = "jack"
    type = "ibgp"
    export_next_hop = "original"
}

resource "panos_panorama_bgp_peer_group" "wang" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    name = "wang"
    type = "ibgp"
    export_next_hop = "use-self"
}

resource "panos_panorama_bgp_conditional_adv" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    name = %q
    used_by = [panos_panorama_bgp_peer_group.%s.name]
    enable = %t
}

resource "panos_panorama_bgp_conditional_adv_non_exist_filter" "x" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    bgp_conditional_adv = panos_panorama_bgp_conditional_adv.test.name
    name = "nef"
    community_regex = "*foo*"
    address_prefixes = ["9.8.7.0/24"]
}

resource "panos_panorama_bgp_conditional_adv_advertise_filter" "x" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.x.virtual_router
    bgp_conditional_adv = panos_panorama_bgp_conditional_adv.test.name
    name = "af"
    community_regex = "*bar*"
    address_prefixes = ["7.8.9.0/24"]
}
`, tmpl, vr, name, ub, en)
}
