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
	datapat "github.com/PaloAltoNetworks/pango/objs/custom/data"
	"github.com/PaloAltoNetworks/pango/objs/dug"
	"github.com/PaloAltoNetworks/pango/objs/edl"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist/action"
	dfsp "github.com/PaloAltoNetworks/pango/objs/profile/security/data"
	dpsp "github.com/PaloAltoNetworks/pango/objs/profile/security/dos"
	fprof "github.com/PaloAltoNetworks/pango/objs/profile/security/file"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/spyware"
	ufsp "github.com/PaloAltoNetworks/pango/objs/profile/security/url"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/virus"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/vulnerability"
	wfasp "github.com/PaloAltoNetworks/pango/objs/profile/security/wildfire"
	"github.com/PaloAltoNetworks/pango/objs/srvc"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"
	"github.com/PaloAltoNetworks/pango/objs/tags"
)

// FwObjs is the client.Objects namespace.
type FwObjs struct {
	Address                             *addr.Firewall
	AddressGroup                        *addrgrp.FwAddrGrp
	AntiSpywareProfile                  *spyware.Firewall
	AntivirusProfile                    *virus.Firewall
	Application                         *app.FwApp
	AppGroup                            *appgrp.FwGroup
	AppSignature                        *signature.FwSignature
	AppSigAndCond                       *andcond.FwAndCond
	AppSigOrCond                        *orcond.FwOrCond
	DataPattern                         *datapat.Firewall
	DataFilteringProfile                *dfsp.Firewall
	DosProtectionProfile                *dpsp.Firewall
	DynamicUserGroup                    *dug.Firewall
	Edl                                 *edl.FwEdl
	FileBlockingProfile                 *fprof.Firewall
	LogForwardingProfile                *logfwd.Firewall
	LogForwardingProfileMatchList       *matchlist.FwMatchList
	LogForwardingProfileMatchListAction *action.FwAction
	Services                            *srvc.FwSrvc
	ServiceGroup                        *srvcgrp.FwSrvcGrp
	Tags                                *tags.FwTags
	UrlFilteringProfile                 *ufsp.Firewall
	VulnerabilityProfile                *vulnerability.Firewall
	WildfireAnalysisProfile             *wfasp.Firewall
}

// Initialize is invoked on client.Initialize().
func (c *FwObjs) Initialize(i util.XapiClient) {
	c.Address = addr.FirewallNamespace(i)

	c.AddressGroup = &addrgrp.FwAddrGrp{}
	c.AddressGroup.Initialize(i)

	c.AntiSpywareProfile = spyware.FirewallNamespace(i)
	c.AntivirusProfile = virus.FirewallNamespace(i)

	c.Application = &app.FwApp{}
	c.Application.Initialize(i)

	c.AppGroup = &appgrp.FwGroup{}
	c.AppGroup.Initialize(i)

	c.AppSignature = &signature.FwSignature{}
	c.AppSignature.Initialize(i)

	c.AppSigAndCond = &andcond.FwAndCond{}
	c.AppSigAndCond.Initialize(i)

	c.AppSigOrCond = &orcond.FwOrCond{}
	c.AppSigOrCond.Initialize(i)

	c.DataFilteringProfile = dfsp.FirewallNamespace(i)
	c.DataPattern = datapat.FirewallNamespace(i)
	c.DosProtectionProfile = dpsp.FirewallNamespace(i)
	c.DynamicUserGroup = dug.FirewallNamespace(i)

	c.Edl = &edl.FwEdl{}
	c.Edl.Initialize(i)

	c.FileBlockingProfile = fprof.FirewallNamespace(i)
	c.LogForwardingProfile = logfwd.FirewallNamespace(i)

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

	c.UrlFilteringProfile = ufsp.FirewallNamespace(i)
	c.VulnerabilityProfile = vulnerability.FirewallNamespace(i)
	c.WildfireAnalysisProfile = wfasp.FirewallNamespace(i)
}
