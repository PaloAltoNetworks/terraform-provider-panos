package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/peer/group"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaBgpPeerGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o group.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpPeerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpPeerGroupConfig(tmpl, vr, name, group.TypeIbgp, group.NextHopOriginal, "", true, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpPeerGroupExists("panos_panorama_bgp_peer_group.test", &o),
					testAccCheckPanosPanoramaBgpPeerGroupAttributes(&o, group.TypeIbgp, group.NextHopOriginal, "", true, false, true, false),
				),
			},
			{
				Config: testAccPanoramaBgpPeerGroupConfig(tmpl, vr, name, group.TypeEbgp, group.NextHopResolve, group.NextHopOriginal, true, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpPeerGroupExists("panos_panorama_bgp_peer_group.test", &o),
					testAccCheckPanosPanoramaBgpPeerGroupAttributes(&o, group.TypeEbgp, group.NextHopResolve, group.NextHopOriginal, true, true, true, false),
				),
			},
			{
				Config: testAccPanoramaBgpPeerGroupConfig(tmpl, vr, name, group.TypeEbgp, group.NextHopUseSelf, group.NextHopUsePeer, false, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpPeerGroupExists("panos_panorama_bgp_peer_group.test", &o),
					testAccCheckPanosPanoramaBgpPeerGroupAttributes(&o, group.TypeEbgp, group.NextHopUseSelf, group.NextHopUsePeer, false, false, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpPeerGroupExists(n string, o *group.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, name := parsePanoramaBgpPeerGroupId(rs.Primary.ID)
		v, err := pano.Network.BgpPeerGroup.Get(tmpl, ts, vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpPeerGroupAttributes(o *group.Entry, typ, enh, inh string, en, acap, srwsi, rpa bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Type != typ {
			return fmt.Errorf("Type is %q, expected %q", o.Type, typ)
		}

		if o.ExportNextHop != enh {
			return fmt.Errorf("Export next hop is %q, expected %q", o.ExportNextHop, enh)
		}

		if o.ImportNextHop != inh {
			return fmt.Errorf("Import next hop is %q, expected %q", o.ImportNextHop, inh)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.AggregatedConfedAsPath != acap {
			return fmt.Errorf("Aggregated confed AS path is %t, expected %t", o.AggregatedConfedAsPath, acap)
		}

		if o.SoftResetWithStoredInfo != srwsi {
			return fmt.Errorf("Soft reset with stored info is %t, expected %t", o.SoftResetWithStoredInfo, srwsi)
		}

		if o.RemovePrivateAs != rpa {
			return fmt.Errorf("Remove private AS is %t, expected %t", o.RemovePrivateAs, rpa)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpPeerGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_peer_group" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, name := parsePanoramaBgpPeerGroupId(rs.Primary.ID)
			_, err := pano.Network.BgpPeerGroup.Get(tmpl, ts, vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpPeerGroupConfig(tmpl, vr, name, typ, enh, inh string, en, acap, srwsi, rpa bool) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "t" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "e1" {
    template = panos_panorama_template.t.name
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_panorama_virtual_router" "vr" {
    template = panos_panorama_template.t.name
    name = %q
    interfaces = [panos_panorama_ethernet_interface.e1.name]
}

resource "panos_panorama_bgp" "conf" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_panorama_bgp_peer_group" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.conf.virtual_router
    name = %q
    type = %q
    export_next_hop = %q
    import_next_hop = %q
    enable = %t
    aggregated_confed_as_path = %t
    soft_reset_with_stored_info = %t
    remove_private_as = %t
}
`, tmpl, vr, name, typ, enh, inh, en, acap, srwsi, rpa)
}
