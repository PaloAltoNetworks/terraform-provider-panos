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


// PanoNetw is the client.Network namespace.
type PanoNetw struct {
    BfdProfile *bfd.PanoBfd
    EthernetInterface *eth.PanoEth
    IkeCryptoProfile *ike.PanoIke
    IkeGateway *ikegw.PanoIkeGw
    IpsecCryptoProfile *ipsec.PanoIpsec
    IpsecTunnel *ipsectunnel.PanoIpsecTunnel
    IpsecTunnelProxyId *tpiv4.PanoIpv4
    LoopbackInterface *loopback.PanoLoopback
    ManagementProfile *mngtprof.PanoMngtProf
    RedistributionProfile *redist4.PanoIpv4
    StaticRoute *ipv4.PanoIpv4
    TunnelInterface *tunnel.PanoTunnel
    VirtualRouter *router.PanoRouter
    Vlan *vlan.PanoVlan
    VlanInterface *vli.PanoVlan
    Zone *zone.PanoZone
}

// Initialize is invoked on client.Initialize().
func (c *PanoNetw) Initialize(i util.XapiClient) {
    c.BfdProfile = &bfd.PanoBfd{}
    c.BfdProfile.Initialize(i)

    c.EthernetInterface = &eth.PanoEth{}
    c.EthernetInterface.Initialize(i)

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

    c.LoopbackInterface = &loopback.PanoLoopback{}
    c.LoopbackInterface.Initialize(i)

    c.ManagementProfile = &mngtprof.PanoMngtProf{}
    c.ManagementProfile.Initialize(i)

    c.RedistributionProfile = &redist4.PanoIpv4{}
    c.RedistributionProfile.Initialize(i)

    c.StaticRoute = &ipv4.PanoIpv4{}
    c.StaticRoute.Initialize(i)

    c.TunnelInterface = &tunnel.PanoTunnel{}
    c.TunnelInterface.Initialize(i)

    c.VirtualRouter = &router.PanoRouter{}
    c.VirtualRouter.Initialize(i)

    c.Vlan = &vlan.PanoVlan{}
    c.Vlan.Initialize(i)

    c.VlanInterface = &vli.PanoVlan{}
    c.VlanInterface.Initialize(i)

    c.Zone = &zone.PanoZone{}
    c.Zone.Initialize(i)
}
