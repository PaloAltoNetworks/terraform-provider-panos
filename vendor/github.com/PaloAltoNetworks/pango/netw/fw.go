package netw


import (
    "github.com/PaloAltoNetworks/pango/netw/ikegw"
    "github.com/PaloAltoNetworks/pango/netw/interface/eth"
    "github.com/PaloAltoNetworks/pango/netw/interface/loopback"
    "github.com/PaloAltoNetworks/pango/netw/interface/tunnel"
    vli "github.com/PaloAltoNetworks/pango/netw/interface/vlan"
    "github.com/PaloAltoNetworks/pango/netw/ipsectunnel"
    tpiv4 "github.com/PaloAltoNetworks/pango/netw/ipsectunnel/proxyid/ipv4"
    "github.com/PaloAltoNetworks/pango/netw/profile/bfd"
    "github.com/PaloAltoNetworks/pango/netw/profile/ike"
    "github.com/PaloAltoNetworks/pango/netw/profile/ipsec"
    "github.com/PaloAltoNetworks/pango/netw/profile/mngtprof"
    redist4 "github.com/PaloAltoNetworks/pango/netw/routing/profile/redist/ipv4"
    "github.com/PaloAltoNetworks/pango/netw/routing/router"
    "github.com/PaloAltoNetworks/pango/netw/routing/route/static/ipv4"
    "github.com/PaloAltoNetworks/pango/netw/vlan"
    "github.com/PaloAltoNetworks/pango/netw/zone"
    "github.com/PaloAltoNetworks/pango/util"
)


// Netw is the client.Network namespace.
type FwNetw struct {
    BfdProfile *bfd.FwBfd
    EthernetInterface *eth.FwEth
    IkeCryptoProfile *ike.FwIke
    IkeGateway *ikegw.FwIkeGw
    IpsecCryptoProfile *ipsec.FwIpsec
    IpsecTunnel *ipsectunnel.FwIpsecTunnel
    IpsecTunnelProxyId *tpiv4.FwIpv4
    LoopbackInterface *loopback.FwLoopback
    ManagementProfile *mngtprof.FwMngtProf
    RedistributionProfile *redist4.FwIpv4
    StaticRoute *ipv4.FwIpv4
    TunnelInterface *tunnel.FwTunnel
    VirtualRouter *router.FwRouter
    Vlan *vlan.FwVlan
    VlanInterface *vli.FwVlan
    Zone *zone.FwZone
}

// Initialize is invoked on client.Initialize().
func (c *FwNetw) Initialize(i util.XapiClient) {
    c.BfdProfile = &bfd.FwBfd{}
    c.BfdProfile.Initialize(i)

    c.EthernetInterface = &eth.FwEth{}
    c.EthernetInterface.Initialize(i)

    c.IkeCryptoProfile = &ike.FwIke{}
    c.IkeCryptoProfile.Initialize(i)

    c.IkeGateway = &ikegw.FwIkeGw{}
    c.IkeGateway.Initialize(i)

    c.IpsecCryptoProfile = &ipsec.FwIpsec{}
    c.IpsecCryptoProfile.Initialize(i)

    c.IpsecTunnel = &ipsectunnel.FwIpsecTunnel{}
    c.IpsecTunnel.Initialize(i)

    c.IpsecTunnelProxyId = &tpiv4.FwIpv4{}
    c.IpsecTunnelProxyId.Initialize(i)

    c.LoopbackInterface = &loopback.FwLoopback{}
    c.LoopbackInterface.Initialize(i)

    c.ManagementProfile = &mngtprof.FwMngtProf{}
    c.ManagementProfile.Initialize(i)

    c.RedistributionProfile = &redist4.FwIpv4{}
    c.RedistributionProfile.Initialize(i)

    c.StaticRoute = &ipv4.FwIpv4{}
    c.StaticRoute.Initialize(i)

    c.TunnelInterface = &tunnel.FwTunnel{}
    c.TunnelInterface.Initialize(i)

    c.VirtualRouter = &router.FwRouter{}
    c.VirtualRouter.Initialize(i)

    c.Vlan = &vlan.FwVlan{}
    c.Vlan.Initialize(i)

    c.VlanInterface = &vli.FwVlan{}
    c.VlanInterface.Initialize(i)

    c.Zone = &zone.FwZone{}
    c.Zone.Initialize(i)
}
