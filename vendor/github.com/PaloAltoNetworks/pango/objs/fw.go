package objs


import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/objs/addr"
    "github.com/PaloAltoNetworks/pango/objs/addrgrp"
    "github.com/PaloAltoNetworks/pango/objs/app"
    appgrp "github.com/PaloAltoNetworks/pango/objs/app/group"
    "github.com/PaloAltoNetworks/pango/objs/app/signature"
    "github.com/PaloAltoNetworks/pango/objs/app/signature/andcond"
    "github.com/PaloAltoNetworks/pango/objs/app/signature/orcond"
    "github.com/PaloAltoNetworks/pango/objs/edl"
    "github.com/PaloAltoNetworks/pango/objs/profile/logfwd"
    "github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist"
    "github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist/action"
    "github.com/PaloAltoNetworks/pango/objs/srvc"
    "github.com/PaloAltoNetworks/pango/objs/srvcgrp"
    "github.com/PaloAltoNetworks/pango/objs/tags"
)


// FwObjs is the client.Objects namespace.
type FwObjs struct {
    Address *addr.FwAddr
    AddressGroup *addrgrp.FwAddrGrp
    Application *app.FwApp
    AppGroup *appgrp.FwGroup
    AppSignature *signature.FwSignature
    AppSigAndCond *andcond.FwAndCond
    AppSigAndCondOrCond *orcond.FwOrCond
    Edl *edl.FwEdl
    LogForwardingProfile *logfwd.FwLogFwd
    LogForwardingProfileMatchList *matchlist.FwMatchList
    LogForwardingProfileMatchListAction *action.FwAction
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

    c.Application = &app.FwApp{}
    c.Application.Initialize(i)

    c.AppGroup = &appgrp.FwGroup{}
    c.AppGroup.Initialize(i)

    c.AppSignature = &signature.FwSignature{}
    c.AppSignature.Initialize(i)

    c.AppSigAndCond = &andcond.FwAndCond{}
    c.AppSigAndCond.Initialize(i)

    c.AppSigAndCondOrCond = &orcond.FwOrCond{}
    c.AppSigAndCondOrCond.Initialize(i)

    c.Edl = &edl.FwEdl{}
    c.Edl.Initialize(i)

    c.LogForwardingProfile = &logfwd.FwLogFwd{}
    c.LogForwardingProfile.Initialize(i)

    c.LogForwardingProfileMatchList = &matchlist.FwMatchList{}
    c.LogForwardingProfileMatchList.Initialize(i)

    c.LogForwardingProfileMatchListAction = &action.FwAction{}
    c.LogForwardingProfileMatchListAction.Initialize(i)

    c.Services = &srvc.FwSrvc{}
    c.Services.Initialize(i)

    c.ServiceGroup = &srvcgrp.FwSrvcGrp{}
    c.ServiceGroup.Initialize(i)

    c.Tags = &tags.FwTags{}
    c.Tags.Initialize(i)
}
