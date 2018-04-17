package objs


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/objs/addr"
    "github.com/PaloAltoNetworks/pango/objs/addrgrp"
    "github.com/PaloAltoNetworks/pango/objs/srvc"
    "github.com/PaloAltoNetworks/pango/objs/srvcgrp"
    "github.com/PaloAltoNetworks/pango/objs/tags"
)


// PanoObjs is the client.Objects namespace.
type PanoObjs struct {
    Address *addr.PanoAddr
    AddressGroup *addrgrp.PanoAddrGrp
    Services *srvc.PanoSrvc
    ServiceGroup *srvcgrp.PanoSrvcGrp
    Tags *tags.PanoTags
}

// Initialize is invoked on client.Initialize().
func (c *PanoObjs) Initialize(i util.XapiClient) {
    c.Address = &addr.PanoAddr{}
    c.Address.Initialize(i)

    c.AddressGroup = &addrgrp.PanoAddrGrp{}
    c.AddressGroup.Initialize(i)

    c.Services = &srvc.PanoSrvc{}
    c.Services.Initialize(i)

    c.ServiceGroup = &srvcgrp.PanoSrvcGrp{}
    c.ServiceGroup.Initialize(i)

    c.Tags = &tags.PanoTags{}
    c.Tags.Initialize(i)
}
