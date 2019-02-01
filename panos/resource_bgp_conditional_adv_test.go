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

func TestAccPanosBgpConditionalAdv_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o conadv.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpConditionalAdvDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpConditionalAdvConfig(vr, name, "jack", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpConditionalAdvExists("panos_bgp_conditional_adv.test", &o),
					testAccCheckPanosBgpConditionalAdvAttributes(&o, name, "jack", false),
				),
			},
			{
				Config: testAccBgpConditionalAdvConfig(vr, name, "wang", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpConditionalAdvExists("panos_bgp_conditional_adv.test", &o),
					testAccCheckPanosBgpConditionalAdvAttributes(&o, name, "wang", true),
				),
			},
		},
	})
}

func testAccCheckPanosBgpConditionalAdvExists(n string, o *conadv.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, name := parseBgpConditionalAdvId(rs.Primary.ID)
		v, err := fw.Network.BgpConditionalAdv.Get(vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpConditionalAdvAttributes(o *conadv.Entry, name, ub string, en bool) resource.TestCheckFunc {
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

func testAccPanosBgpConditionalAdvDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_conditional_adv" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, name := parseBgpConditionalAdvId(rs.Primary.ID)
			_, err := fw.Network.BgpConditionalAdv.Get(vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpConditionalAdvConfig(vr, name, ub string, en bool) string {
	return fmt.Sprintf(`
resource "panos_virtual_router" "vr" {
    name = %q
}

resource "panos_bgp" "x" {
    virtual_router = "${panos_virtual_router.vr.name}"
    router_id = "5.5.5.5"
    as_number = "55"
    enable = false
}

resource "panos_bgp_peer_group" "jack" {
    virtual_router = "${panos_bgp.x.virtual_router}"
    name = "jack"
    type = "ibgp"
    export_next_hop = "original"
}

resource "panos_bgp_peer_group" "wang" {
    virtual_router = "${panos_bgp.x.virtual_router}"
    name = "wang"
    type = "ibgp"
    export_next_hop = "use-self"
}

resource "panos_bgp_conditional_adv" "test" {
    virtual_router = "${panos_bgp.x.virtual_router}"
    name = %q
    used_by = ["${panos_bgp_peer_group.%s.name}"]
    enable = %t
}

resource "panos_bgp_conditional_adv_non_exist_filter" "x" {
    virtual_router = "${panos_bgp.x.virtual_router}"
    bgp_conditional_adv = "${panos_bgp_conditional_adv.test.name}"
    name = "nef"
    community_regex = "*foo*"
    address_prefixes = ["9.8.7.0/24"]
}

resource "panos_bgp_conditional_adv_advertise_filter" "x" {
    virtual_router = "${panos_bgp.x.virtual_router}"
    bgp_conditional_adv = "${panos_bgp_conditional_adv.test.name}"
    name = "af"
    community_regex = "*bar*"
    address_prefixes = ["7.8.9.0/24"]
}
`, vr, name, ub, en)
}
