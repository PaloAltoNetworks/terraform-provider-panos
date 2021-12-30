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
	cusspy "github.com/PaloAltoNetworks/pango/objs/custom/spyware"
	cusurl "github.com/PaloAltoNetworks/pango/objs/custom/url"
	cusvuln "github.com/PaloAltoNetworks/pango/objs/custom/vulnerability"
	"github.com/PaloAltoNetworks/pango/objs/dug"
	"github.com/PaloAltoNetworks/pango/objs/edl"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist"
	"github.com/PaloAltoNetworks/pango/objs/profile/logfwd/matchlist/action"
	dfsp "github.com/PaloAltoNetworks/pango/objs/profile/security/data"
	dpsp "github.com/PaloAltoNetworks/pango/objs/profile/security/dos"
	fprof "github.com/PaloAltoNetworks/pango/objs/profile/security/file"
	spg "github.com/PaloAltoNetworks/pango/objs/profile/security/group"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/spyware"
	ufsp "github.com/PaloAltoNetworks/pango/objs/profile/security/url"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/virus"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/vulnerability"
	wfasp "github.com/PaloAltoNetworks/pango/objs/profile/security/wildfire"
	"github.com/PaloAltoNetworks/pango/objs/srvc"
	"github.com/PaloAltoNetworks/pango/objs/srvcgrp"
	"github.com/PaloAltoNetworks/pango/objs/tags"
)

// PanoObjs is the client.Objects namespace.
type PanoObjs struct {
	Address                             *addr.Panorama
	AddressGroup                        *addrgrp.Panorama
	AntiSpywareProfile                  *spyware.Panorama
	AntivirusProfile                    *virus.Panorama
	Application                         *app.Panorama
	AppGroup                            *appgrp.PanoGroup
	AppSignature                        *signature.PanoSignature
	AppSigAndCond                       *andcond.PanoAndCond
	AppSigOrCond                        *orcond.PanoOrCond
	CustomSpyware                       *cusspy.Panorama
	CustomUrlCategory                   *cusurl.Panorama
	CustomVulnerability                 *cusvuln.Panorama
	DataFilteringProfile                *dfsp.Panorama
	DataPattern                         *datapat.Panorama
	DosProtectionProfile                *dpsp.Panorama
	DynamicUserGroup                    *dug.Panorama
	Edl                                 *edl.Panorama
	FileBlockingProfile                 *fprof.Panorama
	LogForwardingProfile                *logfwd.Panorama
	LogForwardingProfileMatchList       *matchlist.PanoMatchList
	LogForwardingProfileMatchListAction *action.PanoAction
	SecurityProfileGroup                *spg.Panorama
	Services                            *srvc.Panorama
	ServiceGroup                        *srvcgrp.Panorama
	Tags                                *tags.Panorama
	UrlFilteringProfile                 *ufsp.Panorama
	VulnerabilityProfile                *vulnerability.Panorama
	WildfireAnalysisProfile             *wfasp.Panorama
}

// Initialize is invoked on client.Initialize().
func (c *PanoObjs) Initialize(i util.XapiClient) {
	c.Address = addr.PanoramaNamespace(i)
	c.AddressGroup = addrgrp.PanoramaNamespace(i)
	c.AntiSpywareProfile = spyware.PanoramaNamespace(i)
	c.AntivirusProfile = virus.PanoramaNamespace(i)
	c.Application = app.PanoramaNamespace(i)

	c.AppGroup = &appgrp.PanoGroup{}
	c.AppGroup.Initialize(i)

	c.AppSignature = &signature.PanoSignature{}
	c.AppSignature.Initialize(i)

	c.AppSigAndCond = &andcond.PanoAndCond{}
	c.AppSigAndCond.Initialize(i)

	c.AppSigOrCond = &orcond.PanoOrCond{}
	c.AppSigOrCond.Initialize(i)

	c.CustomSpyware = cusspy.PanoramaNamespace(i)
	c.CustomUrlCategory = cusurl.PanoramaNamespace(i)
	c.CustomVulnerability = cusvuln.PanoramaNamespace(i)
	c.DataFilteringProfile = dfsp.PanoramaNamespace(i)
	c.DataPattern = datapat.PanoramaNamespace(i)
	c.DosProtectionProfile = dpsp.PanoramaNamespace(i)
	c.DynamicUserGroup = dug.PanoramaNamespace(i)
	c.Edl = edl.PanoramaNamespace(i)
	c.FileBlockingProfile = fprof.PanoramaNamespace(i)
	c.LogForwardingProfile = logfwd.PanoramaNamespace(i)

	c.LogForwardingProfileMatchList = &matchlist.PanoMatchList{}
	c.LogForwardingProfileMatchList.Initialize(i)

	c.LogForwardingProfileMatchListAction = &action.PanoAction{}
	c.LogForwardingProfileMatchListAction.Initialize(i)

	c.SecurityProfileGroup = spg.PanoramaNamespace(i)
	c.Services = srvc.PanoramaNamespace(i)
	c.ServiceGroup = srvcgrp.PanoramaNamespace(i)
	c.Tags = tags.PanoramaNamespace(i)
	c.UrlFilteringProfile = ufsp.PanoramaNamespace(i)
	c.VulnerabilityProfile = vulnerability.PanoramaNamespace(i)
	c.WildfireAnalysisProfile = wfasp.PanoramaNamespace(i)
}
