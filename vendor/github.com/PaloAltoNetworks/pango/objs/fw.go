package objs


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/objs/addr"
    "github.com/PaloAltoNetworks/pango/objs/addrgrp"
    "github.com/PaloAltoNetworks/pango/objs/srvc"
    "github.com/PaloAltoNetworks/pango/objs/srvcgrp"
    "github.com/PaloAltoNetworks/pango/objs/tags"
)


// FwObjs is the client.Objects namespace.
type FwObjs struct {
    Address *addr.FwAddr
    AddressGroup *addrgrp.FwAddrGrp
    Services *srvc.FwSrvc
    ServiceGroup *srvcgrp.FwSrvcGrp
    Tags *tags.FwTags
}

// Initialize is invoked on client.Initialize().
func (c *FwObjs) Initialize(i util.XapiClient) {
    c.Address = &addr.FwAddr{}
    c.Address.Initialize(i)

    c.AddressGroup = &addrgrp.FwAddrGrp{}
    c.AddressGroup.Initialize(i)

    c.Services = &srvc.FwSrvc{}
    c.Services.Initialize(i)

    c.ServiceGroup = &srvcgrp.FwSrvcGrp{}
    c.ServiceGroup.Initialize(i)

    c.Tags = &tags.FwTags{}
    c.Tags.Initialize(i)
}
