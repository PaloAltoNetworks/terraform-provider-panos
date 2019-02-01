package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaBgp_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o bgp.Config
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpConfig(tmpl, name, "5.5.5.5", "42", true, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpExists("panos_panorama_bgp.test", &o),
					testAccCheckPanosPanoramaBgpAttributes(&o, "5.5.5.5", "42", true, false, true, false),
				),
			},
			{
				Config: testAccPanoramaBgpConfig(tmpl, name, "6.6.6.6", "420", false, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpExists("panos_panorama_bgp.test", &o),
					testAccCheckPanosPanoramaBgpAttributes(&o, "6.6.6.6", "420", false, true, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpExists(n string, o *bgp.Config) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaBgpId(rs.Primary.ID)
		v, err := pano.Network.BgpConfig.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpAttributes(o *bgp.Config, ri, an string, en, rdr, ir, am bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.RouterId != ri {
			return fmt.Errorf("Router Id is %q, expected %q", o.RouterId, ri)
		}

		if o.AsNumber != an {
			return fmt.Errorf("AS number is %q, expected %q", o.AsNumber, an)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.RejectDefaultRoute != rdr {
			return fmt.Errorf("Reject default route is %t, not %t", o.RejectDefaultRoute, rdr)
		}

		if o.InstallRoute != ir {
			return fmt.Errorf("Install route is %t, not %t", o.InstallRoute, ir)
		}

		if o.AggregateMed != am {
			return fmt.Errorf("Aggregate MED is %t, not %t", o.AggregateMed, am)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaBgpId(rs.Primary.ID)
			_, err := pano.Network.BgpConfig.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpConfig(tmpl, name, ri, an string, en, rdr, ir, am bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "t" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "e1" {
    template = "${panos_panorama_template.t.name}"
    name = "ethernet1/1"
    mode = "layer3"
}

resource "panos_panorama_virtual_router" "vr" {
    template = "${panos_panorama_template.t.name}"
    name = %q
    interfaces = ["${panos_panorama_ethernet_interface.e1.name}"]
}

resource "panos_panorama_bgp" "test" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_virtual_router.vr.name}"
    router_id = %q
    as_number = %q
    enable = %t
    reject_default_route = %t
    install_route = %t
    aggregate_med = %t
}
`, tmpl, name, ri, an, en, rdr, ir, am)
}
