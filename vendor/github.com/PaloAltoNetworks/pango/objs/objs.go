// Package objs is the client.Objects namespace.
package objs


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/objs/addr"
    "github.com/PaloAltoNetworks/pango/objs/addrgrp"
    "github.com/PaloAltoNetworks/pango/objs/srvc"
    "github.com/PaloAltoNetworks/pango/objs/srvcgrp"
    "github.com/PaloAltoNetworks/pango/objs/tags"
)


// Objs is the client.Objects namespace.
type Objs struct {
    Address *addr.Addr
    AddressGroup *addrgrp.AddrGrp
    Services *srvc.Srvc
    ServiceGroup *srvcgrp.SrvcGrp
    Tags *tags.Tags
}

// Initialize is invoked on client.Initialize().
func (c *Objs) Initialize(i util.XapiClient) {
    c.Address = &addr.Addr{}
    c.Address.Initialize(i)

    c.AddressGroup = &addrgrp.AddrGrp{}
    c.AddressGroup.Initialize(i)

    c.Services = &srvc.Srvc{}
    c.Services.Initialize(i)

    c.ServiceGroup = &srvcgrp.SrvcGrp{}
    c.ServiceGroup.Initialize(i)

    c.Tags = &tags.Tags{}
    c.Tags.Initialize(i)
}
