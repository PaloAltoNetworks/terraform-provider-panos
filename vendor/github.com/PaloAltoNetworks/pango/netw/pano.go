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

// PanoNetw is the client.Network namespace.
type PanoNetw struct {
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
	GreTunnel                *gre.PanoGre
	IkeCryptoProfile         *ike.PanoIke
	IkeGateway               *ikegw.PanoIkeGw
	IpsecCryptoProfile       *ipsec.PanoIpsec
	IpsecTunnel              *ipsectunnel.PanoIpsecTunnel
	IpsecTunnelProxyId       *tpiv4.PanoIpv4
	Ipv6Address              *ipv6a.Panorama
	Ipv6NeighborDiscovery    *ipv6n.Panorama
	Ipv6StaticRoute          *ipv6sr.Panorama
	Layer2Subinterface       *layer2.Panorama
	Layer3Subinterface       *layer3.Panorama
	LoopbackInterface        *loopback.Panorama
	ManagementProfile        *mngtprof.PanoMngtProf
	MonitorProfile           *monitor.PanoMonitor
	OspfArea                 *ospfarea.Panorama
	OspfAreaInterface        *ospfint.Panorama
	OspfAreaVirtualLink      *ospfvlink.Panorama
	OspfAuthProfile          *ospfauth.Panorama
	OspfConfig               *ospf.Panorama
	OspfExport               *ospfexp.Panorama
	RedistributionProfile    *redist4.PanoIpv4
	StaticRoute              *ipv4.Panorama
	TunnelInterface          *tunnel.Panorama
	VirtualRouter            *router.Panorama
	Vlan                     *vlan.Panorama
	VlanInterface            *vli.Panorama
	Zone                     *zone.Panorama
}

// Initialize is invoked on client.Initialize().
func (c *PanoNetw) Initialize(i util.XapiClient) {
	c.AggregateInterface = aggeth.PanoramaNamespace(i)
	c.Arp = arp.PanoramaNamespace(i)
	c.BfdProfile = bfd.PanoramaNamespace(i)
	c.BgpAggregate = aggregate.PanoramaNamespace(i)
	c.BgpAggAdvertiseFilter = agaf.PanoramaNamespace(i)
	c.BgpAggSuppressFilter = suppress.PanoramaNamespace(i)
	c.BgpAuthProfile = auth.PanoramaNamespace(i)
	c.BgpConAdvAdvertiseFilter = advertise.PanoramaNamespace(i)
	c.BgpConAdvNonExistFilter = nonexist.PanoramaNamespace(i)
	c.BgpConditionalAdv = conadv.PanoramaNamespace(i)
	c.BgpConfig = bgp.PanoramaNamespace(i)
	c.BgpDampeningProfile = dampening.PanoramaNamespace(i)
	c.BgpExport = exp.PanoramaNamespace(i)
	c.BgpImport = imp.PanoramaNamespace(i)
	c.BgpPeer = peer.PanoramaNamespace(i)
	c.BgpPeerGroup = group.PanoramaNamespace(i)
	c.BgpRedistRule = bgpredist.PanoramaNamespace(i)
	c.EthernetInterface = eth.PanoramaNamespace(i)

	c.GreTunnel = &gre.PanoGre{}
	c.GreTunnel.Initialize(i)

	c.IkeCryptoProfile = &ike.PanoIke{}
	c.IkeCryptoProfile.Initialize(i)

	c.IkeGateway = &ikegw.PanoIkeGw{}
	c.IkeGateway.Initialize(i)

	c.IpsecCryptoProfile = &ipsec.PanoIpsec{}
	c.IpsecCryptoProfile.Initialize(i)

	c.IpsecTunnel = &ipsectunnel.PanoIpsecTunnel{}
	c.IpsecTunnel.Initialize(i)

	c.IpsecTunnelProxyId = &tpiv4.PanoIpv4{}
	c.IpsecTunnelProxyId.Initialize(i)

	c.Ipv6Address = ipv6a.PanoramaNamespace(i)
	c.Ipv6NeighborDiscovery = ipv6n.PanoramaNamespace(i)
	c.Ipv6StaticRoute = ipv6sr.PanoramaNamespace(i)
	c.Layer2Subinterface = layer2.PanoramaNamespace(i)
	c.Layer3Subinterface = layer3.PanoramaNamespace(i)
	c.LoopbackInterface = loopback.PanoramaNamespace(i)

	c.ManagementProfile = &mngtprof.PanoMngtProf{}
	c.ManagementProfile.Initialize(i)

	c.MonitorProfile = &monitor.PanoMonitor{}
	c.MonitorProfile.Initialize(i)

	c.OspfArea = ospfarea.PanoramaNamespace(i)
	c.OspfAreaInterface = ospfint.PanoramaNamespace(i)
	c.OspfAreaVirtualLink = ospfvlink.PanoramaNamespace(i)
	c.OspfAuthProfile = ospfauth.PanoramaNamespace(i)
	c.OspfConfig = ospf.PanoramaNamespace(i)
	c.OspfExport = ospfexp.PanoramaNamespace(i)

	c.RedistributionProfile = &redist4.PanoIpv4{}
	c.RedistributionProfile.Initialize(i)

	c.StaticRoute = ipv4.PanoramaNamespace(i)
	c.TunnelInterface = tunnel.PanoramaNamespace(i)
	c.VirtualRouter = router.PanoramaNamespace(i)
	c.Vlan = vlan.PanoramaNamespace(i)
	c.VlanInterface = vli.PanoramaNamespace(i)
	c.Zone = zone.PanoramaNamespace(i)
}
