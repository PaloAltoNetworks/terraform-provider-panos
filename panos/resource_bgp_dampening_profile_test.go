package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/profile/dampening"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosBgpDampeningProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o dampening.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpDampeningProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpDampeningProfileConfig(vr, name, true, 1.5, 1.5, 800, 400, 800),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpDampeningProfileExists("panos_bgp_dampening_profile.test", &o),
					testAccCheckPanosBgpDampeningProfileAttributes(&o, true, 1.5, 1.5, 800, 400, 800),
				),
			},
			{
				Config: testAccBgpDampeningProfileConfig(vr, name, false, 1.25, 0.5, 900, 300, 900),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpDampeningProfileExists("panos_bgp_dampening_profile.test", &o),
					testAccCheckPanosBgpDampeningProfileAttributes(&o, false, 1.25, 0.5, 900, 300, 900),
				),
			},
		},
	})
}

func testAccCheckPanosBgpDampeningProfileExists(n string, o *dampening.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, name := parseBgpDampeningProfileId(rs.Primary.ID)
		v, err := fw.Network.BgpDampeningProfile.Get(vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpDampeningProfileAttributes(o *dampening.Entry, en bool, cut, ru float64, mht, dhlr, dhlu int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.Cutoff != cut {
			return fmt.Errorf("Cutoff is %f, not %f", o.Cutoff, cut)
		}

		if o.Reuse != ru {
			return fmt.Errorf("Reuse is %f, not %f", o.Reuse, ru)
		}

		if o.MaxHoldTime != mht {
			return fmt.Errorf("Max hold time is %d, not %d", o.MaxHoldTime, mht)
		}

		if o.DecayHalfLifeReachable != dhlr {
			return fmt.Errorf("Decay half life reachable is %d, not %d", o.DecayHalfLifeReachable, dhlr)
		}

		if o.DecayHalfLifeUnreachable != dhlu {
			return fmt.Errorf("Decay half life unreachable is %d, not %d", o.DecayHalfLifeUnreachable, dhlu)
		}

		return nil
	}
}

func testAccPanosBgpDampeningProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_dampening_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, name := parseBgpDampeningProfileId(rs.Primary.ID)
			_, err := fw.Network.BgpDampeningProfile.Get(vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpDampeningProfileConfig(vr, name string, en bool, cut, ru float64, mht, dhlr, dhlu int) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "e1" {
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_virtual_router" "vr" {
    name = %q
    interfaces = ["${panos_ethernet_interface.e1.name}"]
}

resource "panos_bgp" "conf" {
    virtual_router = "${panos_virtual_router.vr.name}"
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_bgp_dampening_profile" "test" {
    virtual_router = "${panos_bgp.conf.virtual_router}"
    name = %q
    enable = %t
    cutoff = %f
    reuse = %f
    max_hold_time = %d
    decay_half_life_reachable = %d
    decay_half_life_unreachable = %d
}
`, vr, name, en, cut, ru, mht, dhlr, dhlu)
}
