package security

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a security
// rule.
//
// Targets is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
	Name                            string
	Type                            string
	Description                     string
	Tags                            []string // ordered
	SourceZones                     []string // unordered
	SourceAddresses                 []string // unordered
	NegateSource                    bool
	SourceUsers                     []string // unordered
	HipProfiles                     []string // unordered
	DestinationZones                []string // unordered
	DestinationAddresses            []string // unordered
	NegateDestination               bool
	Applications                    []string // unordered
	Services                        []string // unordered
	Categories                      []string // unordered
	Action                          string
	LogSetting                      string
	LogStart                        bool
	LogEnd                          bool
	Disabled                        bool
	Schedule                        string
	IcmpUnreachable                 bool
	DisableServerResponseInspection bool
	Group                           string
	Targets                         map[string][]string
	NegateTarget                    bool
	Virus                           string
	Spyware                         string
	Vulnerability                   string
	UrlFiltering                    string
	FileBlocking                    string
	WildFireAnalysis                string
	DataFiltering                   string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * Type: "universal"
//      * SourceZones: ["any"]
//      * SourceAddresses: ["any"]
//      * SourceUsers: ["any"]
//      * HipProfiles: ["any"]
//      * DestinationZones: ["any"]
//      * DestinationAddresses: ["any"]
//      * Applications: ["any"]
//      * Services: ["application-default"]
//      * Categories: ["any"]
//      * Action: "allow"
//      * LogEnd: true
func (o *Entry) Defaults() {
	if o.Type == "" {
		o.Type = "universal"
	}

	if len(o.SourceZones) == 0 {
		o.SourceZones = []string{"any"}
	}

	if len(o.DestinationZones) == 0 {
		o.DestinationZones = []string{"any"}
	}

	if len(o.SourceAddresses) == 0 {
		o.SourceAddresses = []string{"any"}
	}

	if len(o.SourceUsers) == 0 {
		o.SourceUsers = []string{"any"}
	}

	if len(o.HipProfiles) == 0 {
		o.HipProfiles = []string{"any"}
	}

	if len(o.DestinationAddresses) == 0 {
		o.DestinationAddresses = []string{"any"}
	}

	if len(o.Applications) == 0 {
		o.Applications = []string{"any"}
	}

	if len(o.Services) == 0 {
		o.Services = []string{"application-default"}
	}

	if len(o.Categories) == 0 {
		o.Categories = []string{"any"}
	}

	if o.Action == "" {
		o.Action = "allow"
	}

	if !o.LogEnd {
		o.LogEnd = true
	}
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Type = s.Type
	o.Description = s.Description
	o.Tags = s.Tags
	o.SourceZones = s.SourceZones
	o.SourceAddresses = s.SourceAddresses
	o.NegateSource = s.NegateSource
	o.SourceUsers = s.SourceUsers
	o.HipProfiles = s.HipProfiles
	o.DestinationZones = s.DestinationZones
	o.DestinationAddresses = s.DestinationAddresses
	o.NegateDestination = s.NegateDestination
	o.Applications = s.Applications
	o.Services = s.Services
	o.Categories = s.Categories
	o.Action = s.Action
	o.LogSetting = s.LogSetting
	o.LogStart = s.LogStart
	o.LogEnd = s.LogEnd
	o.Disabled = s.Disabled
	o.Schedule = s.Schedule
	o.IcmpUnreachable = s.IcmpUnreachable
	o.DisableServerResponseInspection = s.DisableServerResponseInspection
	o.Group = s.Group
	o.Targets = s.Targets
	o.NegateTarget = s.NegateTarget
	o.Virus = s.Virus
	o.Spyware = s.Spyware
	o.Vulnerability = s.Vulnerability
	o.UrlFiltering = s.UrlFiltering
	o.FileBlocking = s.FileBlocking
	o.WildFireAnalysis = s.WildFireAnalysis
	o.DataFiltering = s.DataFiltering
}

/** Structs / functions for normalization. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v1) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                 o.Name,
		Type:                 o.Type,
		Description:          o.Description,
		Tags:                 util.MemToStr(o.Tags),
		SourceZones:          util.MemToStr(o.SourceZones),
		DestinationZones:     util.MemToStr(o.DestinationZones),
		SourceAddresses:      util.MemToStr(o.SourceAddresses),
		NegateSource:         util.AsBool(o.NegateSource),
		SourceUsers:          util.MemToStr(o.SourceUsers),
		HipProfiles:          util.MemToStr(o.HipProfiles),
		DestinationAddresses: util.MemToStr(o.DestinationAddresses),
		NegateDestination:    util.AsBool(o.NegateDestination),
		Applications:         util.MemToStr(o.Applications),
		Services:             util.MemToStr(o.Services),
		Categories:           util.MemToStr(o.Categories),
		Action:               o.Action,
		LogSetting:           o.LogSetting,
		LogStart:             util.AsBool(o.LogStart),
		LogEnd:               util.AsBool(o.LogEnd),
		Disabled:             util.AsBool(o.Disabled),
		Schedule:             o.Schedule,
		IcmpUnreachable:      util.AsBool(o.IcmpUnreachable),
	}
	if o.Options != nil {
		ans.DisableServerResponseInspection = util.AsBool(o.Options.DisableServerResponseInspection)
	}
	if o.TargetInfo != nil {
		ans.NegateTarget = util.AsBool(o.TargetInfo.NegateTarget)
		ans.Targets = util.VsysEntToMap(o.TargetInfo.Targets)
	}
	if o.ProfileSettings != nil {
		ans.Group = util.MemToOneStr(o.ProfileSettings.Group)
		if o.ProfileSettings.Profiles != nil {
			ans.Virus = util.MemToOneStr(o.ProfileSettings.Profiles.Virus)
			ans.Spyware = util.MemToOneStr(o.ProfileSettings.Profiles.Spyware)
			ans.Vulnerability = util.MemToOneStr(o.ProfileSettings.Profiles.Vulnerability)
			ans.UrlFiltering = util.MemToOneStr(o.ProfileSettings.Profiles.UrlFiltering)
			ans.FileBlocking = util.MemToOneStr(o.ProfileSettings.Profiles.FileBlocking)
			ans.WildFireAnalysis = util.MemToOneStr(o.ProfileSettings.Profiles.WildFireAnalysis)
			ans.DataFiltering = util.MemToOneStr(o.ProfileSettings.Profiles.DataFiltering)
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName              xml.Name         `xml:"entry"`
	Name                 string           `xml:"name,attr"`
	Type                 string           `xml:"rule-type,omitempty"`
	Description          string           `xml:"description,omitempty"`
	Tags                 *util.MemberType `xml:"tag"`
	SourceZones          *util.MemberType `xml:"from"`
	DestinationZones     *util.MemberType `xml:"to"`
	SourceAddresses      *util.MemberType `xml:"source"`
	NegateSource         string           `xml:"negate-source"`
	SourceUsers          *util.MemberType `xml:"source-user"`
	HipProfiles          *util.MemberType `xml:"hip-profiles"`
	DestinationAddresses *util.MemberType `xml:"destination"`
	NegateDestination    string           `xml:"negate-destination"`
	Applications         *util.MemberType `xml:"application"`
	Services             *util.MemberType `xml:"service"`
	Categories           *util.MemberType `xml:"category"`
	Action               string           `xml:"action"`
	LogSetting           string           `xml:"log-setting,omitempty"`
	LogStart             string           `xml:"log-start"`
	LogEnd               string           `xml:"log-end"`
	Disabled             string           `xml:"disabled"`
	Schedule             string           `xml:"schedule,omitempty"`
	IcmpUnreachable      string           `xml:"icmp-unreachable"`
	Options              *secOptions      `xml:"option"`
	TargetInfo           *targetInfo      `xml:"target"`
	ProfileSettings      *profileSettings `xml:"profile-setting"`
}

type secOptions struct {
	DisableServerResponseInspection string `xml:"disable-server-response-inspection,omitempty"`
}

type targetInfo struct {
	Targets      *util.VsysEntryType `xml:"devices"`
	NegateTarget string              `xml:"negate,omitempty"`
}

type profileSettings struct {
	Group    *util.MemberType        `xml:"group"`
	Profiles *profileSettingsProfile `xml:"profiles"`
}

type profileSettingsProfile struct {
	Virus            *util.MemberType `xml:"virus"`
	Spyware          *util.MemberType `xml:"spyware"`
	Vulnerability    *util.MemberType `xml:"vulnerability"`
	UrlFiltering     *util.MemberType `xml:"url-filtering"`
	FileBlocking     *util.MemberType `xml:"file-blocking"`
	WildFireAnalysis *util.MemberType `xml:"wildfire-analysis"`
	DataFiltering    *util.MemberType `xml:"data-filtering"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                 e.Name,
		Type:                 e.Type,
		Description:          e.Description,
		Tags:                 util.StrToMem(e.Tags),
		SourceZones:          util.StrToMem(e.SourceZones),
		DestinationZones:     util.StrToMem(e.DestinationZones),
		SourceAddresses:      util.StrToMem(e.SourceAddresses),
		NegateSource:         util.YesNo(e.NegateSource),
		SourceUsers:          util.StrToMem(e.SourceUsers),
		HipProfiles:          util.StrToMem(e.HipProfiles),
		DestinationAddresses: util.StrToMem(e.DestinationAddresses),
		NegateDestination:    util.YesNo(e.NegateDestination),
		Applications:         util.StrToMem(e.Applications),
		Services:             util.StrToMem(e.Services),
		Categories:           util.StrToMem(e.Categories),
		Action:               e.Action,
		LogSetting:           e.LogSetting,
		LogStart:             util.YesNo(e.LogStart),
		LogEnd:               util.YesNo(e.LogEnd),
		Disabled:             util.YesNo(e.Disabled),
		Schedule:             e.Schedule,
		IcmpUnreachable:      util.YesNo(e.IcmpUnreachable),
		Options:              &secOptions{util.YesNo(e.DisableServerResponseInspection)},
	}
	if e.Targets != nil || e.NegateTarget {
		nfo := &targetInfo{
			Targets:      util.MapToVsysEnt(e.Targets),
			NegateTarget: util.YesNo(e.NegateTarget),
		}
		ans.TargetInfo = nfo
	}
	gs := e.Virus != "" || e.Spyware != "" || e.Vulnerability != "" || e.UrlFiltering != "" || e.FileBlocking != "" || e.WildFireAnalysis != "" || e.DataFiltering != ""
	if e.Group != "" || gs {
		ps := &profileSettings{
			Group: util.OneStrToMem(e.Group),
		}
		if gs {
			ps.Profiles = &profileSettingsProfile{
				util.OneStrToMem(e.Virus),
				util.OneStrToMem(e.Spyware),
				util.OneStrToMem(e.Vulnerability),
				util.OneStrToMem(e.UrlFiltering),
				util.OneStrToMem(e.FileBlocking),
				util.OneStrToMem(e.WildFireAnalysis),
				util.OneStrToMem(e.DataFiltering),
			}
		}
		ans.ProfileSettings = ps
	}

	return ans
}
