package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer/group"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosBgpPeerGroup_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o group.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpPeerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpPeerGroupConfig(vr, name, group.TypeIbgp, group.NextHopOriginal, "", true, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpPeerGroupExists("panos_bgp_peer_group.test", &o),
					testAccCheckPanosBgpPeerGroupAttributes(&o, group.TypeIbgp, group.NextHopOriginal, "", true, false, true, false),
				),
			},
			{
				Config: testAccBgpPeerGroupConfig(vr, name, group.TypeEbgp, group.NextHopResolve, group.NextHopOriginal, true, true, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpPeerGroupExists("panos_bgp_peer_group.test", &o),
					testAccCheckPanosBgpPeerGroupAttributes(&o, group.TypeEbgp, group.NextHopResolve, group.NextHopOriginal, true, true, true, false),
				),
			},
			{
				Config: testAccBgpPeerGroupConfig(vr, name, group.TypeEbgp, group.NextHopUseSelf, group.NextHopUsePeer, false, false, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpPeerGroupExists("panos_bgp_peer_group.test", &o),
					testAccCheckPanosBgpPeerGroupAttributes(&o, group.TypeEbgp, group.NextHopUseSelf, group.NextHopUsePeer, false, false, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosBgpPeerGroupExists(n string, o *group.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, name := parseBgpPeerGroupId(rs.Primary.ID)
		v, err := fw.Network.BgpPeerGroup.Get(vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpPeerGroupAttributes(o *group.Entry, typ, enh, inh string, en, acap, srwsi, rpa bool) resource.TestCheckFunc {
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

func testAccPanosBgpPeerGroupDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_peer_group" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, name := parseBgpPeerGroupId(rs.Primary.ID)
			_, err := fw.Network.BgpPeerGroup.Get(vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpPeerGroupConfig(vr, name, typ, enh, inh string, en, acap, srwsi, rpa bool) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "e1" {
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_virtual_router" "vr" {
    name = %q
    interfaces = [panos_ethernet_interface.e1.name]
}

resource "panos_bgp" "conf" {
    virtual_router = panos_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_bgp_peer_group" "test" {
    virtual_router = panos_bgp.conf.virtual_router
    name = %q
    type = %q
    export_next_hop = %q
    import_next_hop = %q
    enable = %t
    aggregated_confed_as_path = %t
    soft_reset_with_stored_info = %t
    remove_private_as = %t
}
`, vr, name, typ, enh, inh, en, acap, srwsi, rpa)
}
