package netw

import (
	"github.com/PaloAltoNetworks/pango/netw/ikegw"
	aggeth "github.com/PaloAltoNetworks/pango/netw/interface/aggregate"
	"github.com/PaloAltoNetworks/pango/netw/interface/arp"
	"github.com/PaloAltoNetworks/pango/netw/interface/eth"
	ipv6a "github.com/PaloAltoNetworks/pango/netw/interface/ipv6/address"
	ipv6n "github.com/PaloAltoNetworks/pango/netw/interface/ipv6/neighbor"
	"github.com/PaloAltoNetworks/pango/netw/interface/loopback"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer2"
	"github.com/PaloAltoNetworks/pango/netw/interface/subinterface/layer3"
	"github.com/PaloAltoNetworks/pango/netw/interface/tunnel"
	vli "github.com/PaloAltoNetworks/pango/netw/interface/vlan"
	"github.com/PaloAltoNetworks/pango/netw/ipsectunnel"
	tpiv4 "github.com/PaloAltoNetworks/pango/netw/ipsectunnel/proxyid/ipv4"
	"github.com/PaloAltoNetworks/pango/netw/profile/bfd"
	"github.com/PaloAltoNetworks/pango/netw/profile/ike"
	"github.com/PaloAltoNetworks/pango/netw/profile/ipsec"
	"github.com/PaloAltoNetworks/pango/netw/profile/mngtprof"
	"github.com/PaloAltoNetworks/pango/netw/profile/monitor"
	redist4 "github.com/PaloAltoNetworks/pango/netw/routing/profile/redist/ipv4"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate"
	agaf "github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate/filter/advertise"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/aggregate/filter/suppress"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv/filter/advertise"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/conadv/filter/nonexist"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/exp"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/imp"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/peer/group"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/profile/auth"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/profile/dampening"
	bgpredist "github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/redist"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf"
	ospfarea "github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area"
	ospfint "github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area/iface"
	ospfvlink "github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area/vlink"
	ospfexp "github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/exp"
	ospfauth "github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/profile/auth"
	"github.com/PaloAltoNetworks/pango/netw/routing/route/static/ipv4"
	ipv6sr "github.com/PaloAltoNetworks/pango/netw/routing/route/static/ipv6"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"
	"github.com/PaloAltoNetworks/pango/netw/tunnel/gre"
	"github.com/PaloAltoNetworks/pango/netw/vlan"
	"github.com/PaloAltoNetworks/pango/netw/zone"
	"github.com/PaloAltoNetworks/pango/util"
)

// Panorama is the client.Network namespace.
type Panorama struct {
	AggregateInterface       *aggeth.Panorama
	Arp                      *arp.Panorama
	BfdProfile               *bfd.Panorama
	BgpAggregate             *aggregate.Panorama
	BgpAggAdvertiseFilter    *agaf.Panorama
	BgpAggSuppressFilter     *suppress.Panorama
	BgpAuthProfile           *auth.Panorama
	BgpConAdvAdvertiseFilter *advertise.Panorama
	BgpConAdvNonExistFilter  *nonexist.Panorama
	BgpConditionalAdv        *conadv.Panorama
	BgpConfig                *bgp.Panorama
	BgpDampeningProfile      *dampening.Panorama
	BgpExport                *exp.Panorama
	BgpImport                *imp.Panorama
	BgpPeer                  *peer.Panorama
	BgpPeerGroup             *group.Panorama
	BgpRedistRule            *bgpredist.Panorama
	EthernetInterface        *eth.Panorama
	GreTunnel                *gre.Panorama
	IkeCryptoProfile         *ike.Panorama
	IkeGateway               *ikegw.Panorama
	IpsecCryptoProfile       *ipsec.Panorama
	IpsecTunnel              *ipsectunnel.Panorama
	IpsecTunnelProxyId       *tpiv4.Panorama
	Ipv6Address              *ipv6a.Panorama
	Ipv6NeighborDiscovery    *ipv6n.Panorama
	Ipv6StaticRoute          *ipv6sr.Panorama
	Layer2Subinterface       *layer2.Panorama
	Layer3Subinterface       *layer3.Panorama
	LoopbackInterface        *loopback.Panorama
	ManagementProfile        *mngtprof.Panorama
	MonitorProfile           *monitor.Panorama
	OspfArea                 *ospfarea.Panorama
	OspfAreaInterface        *ospfint.Panorama
	OspfAreaVirtualLink      *ospfvlink.Panorama
	OspfAuthProfile          *ospfauth.Panorama
	OspfConfig               *ospf.Panorama
	OspfExport               *ospfexp.Panorama
	RedistributionProfile    *redist4.Panorama
	StaticRoute              *ipv4.Panorama
	TunnelInterface          *tunnel.Panorama
	VirtualRouter            *router.Panorama
	Vlan                     *vlan.Panorama
	VlanInterface            *vli.Panorama
	Zone                     *zone.Panorama
}

func PanoramaNamespace(x util.XapiClient) *Panorama {
	return &Panorama{
		AggregateInterface:       aggeth.PanoramaNamespace(x),
		Arp:                      arp.PanoramaNamespace(x),
		BfdProfile:               bfd.PanoramaNamespace(x),
		BgpAggregate:             aggregate.PanoramaNamespace(x),
		BgpAggAdvertiseFilter:    agaf.PanoramaNamespace(x),
		BgpAggSuppressFilter:     suppress.PanoramaNamespace(x),
		BgpAuthProfile:           auth.PanoramaNamespace(x),
		BgpConAdvAdvertiseFilter: advertise.PanoramaNamespace(x),
		BgpConAdvNonExistFilter:  nonexist.PanoramaNamespace(x),
		BgpConditionalAdv:        conadv.PanoramaNamespace(x),
		BgpConfig:                bgp.PanoramaNamespace(x),
		BgpDampeningProfile:      dampening.PanoramaNamespace(x),
		BgpExport:                exp.PanoramaNamespace(x),
		BgpImport:                imp.PanoramaNamespace(x),
		BgpPeer:                  peer.PanoramaNamespace(x),
		BgpPeerGroup:             group.PanoramaNamespace(x),
		BgpRedistRule:            bgpredist.PanoramaNamespace(x),
		EthernetInterface:        eth.PanoramaNamespace(x),
		GreTunnel:                gre.PanoramaNamespace(x),
		IkeCryptoProfile:         ike.PanoramaNamespace(x),
		IkeGateway:               ikegw.PanoramaNamespace(x),
		IpsecCryptoProfile:       ipsec.PanoramaNamespace(x),
		IpsecTunnel:              ipsectunnel.PanoramaNamespace(x),
		IpsecTunnelProxyId:       tpiv4.PanoramaNamespace(x),
		Ipv6Address:              ipv6a.PanoramaNamespace(x),
		Ipv6NeighborDiscovery:    ipv6n.PanoramaNamespace(x),
		Ipv6StaticRoute:          ipv6sr.PanoramaNamespace(x),
		Layer2Subinterface:       layer2.PanoramaNamespace(x),
		Layer3Subinterface:       layer3.PanoramaNamespace(x),
		LoopbackInterface:        loopback.PanoramaNamespace(x),
		ManagementProfile:        mngtprof.PanoramaNamespace(x),
		MonitorProfile:           monitor.PanoramaNamespace(x),
		OspfArea:                 ospfarea.PanoramaNamespace(x),
		OspfAreaInterface:        ospfint.PanoramaNamespace(x),
		OspfAreaVirtualLink:      ospfvlink.PanoramaNamespace(x),
		OspfAuthProfile:          ospfauth.PanoramaNamespace(x),
		OspfConfig:               ospf.PanoramaNamespace(x),
		OspfExport:               ospfexp.PanoramaNamespace(x),
		RedistributionProfile:    redist4.PanoramaNamespace(x),
		StaticRoute:              ipv4.PanoramaNamespace(x),
		TunnelInterface:          tunnel.PanoramaNamespace(x),
		VirtualRouter:            router.PanoramaNamespace(x),
		Vlan:                     vlan.PanoramaNamespace(x),
		VlanInterface:            vli.PanoramaNamespace(x),
		Zone:                     zone.PanoramaNamespace(x),
	}
}
