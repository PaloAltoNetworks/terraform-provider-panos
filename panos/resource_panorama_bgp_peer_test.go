package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaBgpPeer_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o peer.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	pg := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpPeerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpPeerConfig(tmpl, vr, pg, name, "5.5.6.6", "unlimited", false, 31, 1, 4, 89, 14, 4455, 443),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpPeerExists("panos_panorama_bgp_peer.test", &o),
					testAccCheckPanosPanoramaBgpPeerAttributes(&o, "5.5.6.6", "unlimited", false, 31, 1, 4, 89, 14, 4455, 443),
				),
			},
			{
				Config: testAccPanoramaBgpPeerConfig(tmpl, vr, pg, name, "6.5.6.5", "5000", true, 32, 2, 5, 90, 15, 4321, 554),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpPeerExists("panos_panorama_bgp_peer.test", &o),
					testAccCheckPanosPanoramaBgpPeerAttributes(&o, "6.5.6.5", "5000", true, 32, 2, 5, 90, 15, 4321, 554),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpPeerExists(n string, o *peer.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, pg, name := parsePanoramaBgpPeerId(rs.Primary.ID)
		v, err := pano.Network.BgpPeer.Get(tmpl, ts, vr, pg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpPeerAttributes(o *peer.Entry, pai, mp string, en bool, kai, mh, odt, ht, iht, icrp, oclp int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.PeerAddressIp != pai {
			return fmt.Errorf("Peer address ip is %q, not %q", o.PeerAddressIp, pai)
		}

		if o.MaxPrefixes != mp {
			return fmt.Errorf("Max prefixes is %q, not %q", o.MaxPrefixes, mp)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.KeepAliveInterval != kai {
			return fmt.Errorf("Keep alive interval is %d, not %d", o.KeepAliveInterval, kai)
		}

		if o.MultiHop != mh {
			return fmt.Errorf("Multi hop is %d, not %d", o.MultiHop, mh)
		}

		if o.OpenDelayTime != odt {
			return fmt.Errorf("Open delay time is %d, not %d", o.OpenDelayTime, odt)
		}

		if o.HoldTime != ht {
			return fmt.Errorf("Hold time is %d, not %d", o.HoldTime, ht)
		}

		if o.IdleHoldTime != iht {
			return fmt.Errorf("Idle hold time is %d, not %d", o.IdleHoldTime, iht)
		}

		if o.IncomingConnectionsRemotePort != icrp {
			return fmt.Errorf("Incoming connections remote port is %d, not %d", o.IncomingConnectionsRemotePort, icrp)
		}

		if o.OutgoingConnectionsLocalPort != oclp {
			return fmt.Errorf("Outgoing connections local port is %d, not %d", o.OutgoingConnectionsLocalPort, oclp)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpPeerDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_peer" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, pg, name := parsePanoramaBgpPeerId(rs.Primary.ID)
			_, err := pano.Network.BgpPeer.Get(tmpl, ts, vr, pg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpPeerConfig(tmpl, vr, pg, name, pai, mp string, en bool, kai, mh, odt, ht, iht, icrp, oclp int) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

resource "panos_panorama_template" "t" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "e" {
    template = panos_panorama_template.t.name
    name = "ethernet1/1"
    mode = "layer3"
    static_ips = ["10.5.5.1/24"]
}

resource "panos_panorama_virtual_router" "vr" {
    template = panos_panorama_template.t.name
    name = %q
    interfaces = [panos_panorama_ethernet_interface.e.name]
}

resource "panos_panorama_bgp" "conf" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_panorama_bgp_peer_group" "pg" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.conf.virtual_router
    name = %q
    enable = false
    type = "ibgp"
    export_next_hop = "use-self"
}

resource "panos_panorama_bgp_peer" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.conf.virtual_router
    bgp_peer_group = panos_panorama_bgp_peer_group.pg.name
    local_address_interface = panos_panorama_ethernet_interface.e.name
    local_address_ip = panos_panorama_ethernet_interface.e.static_ips.0
    peer_as = panos_panorama_bgp.conf.as_number
    name = %q
    peer_address_ip = %q
    max_prefixes = %q
    enable = %t
    keep_alive_interval = %d
    multi_hop = %d
    open_delay_time = %d
    hold_time = %d
    idle_hold_time = %d
    incoming_connections_remote_port = %d
    outgoing_connections_local_port = %d
    bfd_profile = (
        data.panos_system_info.x.version_major >= 7 ?
            data.panos_system_info.x.version_minor >= 1 ? "None" : ""
        : ""
    )
    address_family_type = data.panos_system_info.x.version_major >= 8 ? "ipv4" : ""
    reflector_client = data.panos_system_info.x.version_major >= 8 ? "non-client" : ""
    min_route_advertisement_interval = (
        data.panos_system_info.x.version_major >= 8 ?
            data.panos_system_info.x.version_minor >= 1 ? 30 : 0
        : 0
    )
}
`, tmpl, vr, pg, name, pai, mp, en, kai, mh, odt, ht, iht, icrp, oclp)
}
