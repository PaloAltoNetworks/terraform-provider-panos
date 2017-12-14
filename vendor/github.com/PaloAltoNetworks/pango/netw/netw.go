// Package netw is the client.Network namespace.
package netw


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/netw/eth"
    "github.com/PaloAltoNetworks/pango/netw/mngtprof"
    "github.com/PaloAltoNetworks/pango/netw/vlan"
    "github.com/PaloAltoNetworks/pango/netw/zone"
    "github.com/PaloAltoNetworks/pango/netw/router"
)


// Netw is the client.Network namespace.
type Netw struct {
    EthernetInterface *eth.Eth
    ManagementProfile *mngtprof.MngtProf
    Vlan *vlan.Vlan
    Zone *zone.Zone
    VirtualRouter *router.Router
}

// Initialize is invoked on client.Initialize().
func (c *Netw) Initialize(i util.XapiClient) {
    c.EthernetInterface = &eth.Eth{}
    c.EthernetInterface.Initialize(i)

    c.ManagementProfile = &mngtprof.MngtProf{}
    c.ManagementProfile.Initialize(i)

    c.Vlan = &vlan.Vlan{}
    c.Vlan.Initialize(i)

    c.Zone = &zone.Zone{}
    c.Zone.Initialize(i)

    c.VirtualRouter = &router.Router{}
    c.VirtualRouter.Initialize(i)
}
