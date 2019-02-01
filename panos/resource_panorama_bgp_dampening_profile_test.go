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

func TestAccPanosPanoramaBgpDampeningProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o dampening.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpDampeningProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpDampeningProfileConfig(tmpl, vr, name, true, 1.5, 1.5, 800, 400, 800),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpDampeningProfileExists("panos_panorama_bgp_dampening_profile.test", &o),
					testAccCheckPanosPanoramaBgpDampeningProfileAttributes(&o, true, 1.5, 1.5, 800, 400, 800),
				),
			},
			{
				Config: testAccPanoramaBgpDampeningProfileConfig(tmpl, vr, name, false, 1.25, 0.5, 900, 300, 900),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpDampeningProfileExists("panos_panorama_bgp_dampening_profile.test", &o),
					testAccCheckPanosPanoramaBgpDampeningProfileAttributes(&o, false, 1.25, 0.5, 900, 300, 900),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpDampeningProfileExists(n string, o *dampening.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, name := parsePanoramaBgpDampeningProfileId(rs.Primary.ID)
		v, err := pano.Network.BgpDampeningProfile.Get(tmpl, ts, vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpDampeningProfileAttributes(o *dampening.Entry, en bool, cut, ru float64, mht, dhlr, dhlu int) resource.TestCheckFunc {
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

func testAccPanosPanoramaBgpDampeningProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_dampening_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, name := parsePanoramaBgpDampeningProfileId(rs.Primary.ID)
			_, err := pano.Network.BgpDampeningProfile.Get(tmpl, ts, vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpDampeningProfileConfig(tmpl, vr, name string, en bool, cut, ru float64, mht, dhlr, dhlu int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "t" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "e1" {
    template = "${panos_panorama_template.t.name}"
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_panorama_virtual_router" "vr" {
    template = "${panos_panorama_template.t.name}"
    name = %q
    interfaces = ["${panos_panorama_ethernet_interface.e1.name}"]
}

resource "panos_panorama_bgp" "conf" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_virtual_router.vr.name}"
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_panorama_bgp_dampening_profile" "test" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    name = %q
    enable = %t
    cutoff = %f
    reuse = %f
    max_hold_time = %d
    decay_half_life_reachable = %d
    decay_half_life_unreachable = %d
}
`, tmpl, vr, name, en, cut, ru, mht, dhlr, dhlu)
}
