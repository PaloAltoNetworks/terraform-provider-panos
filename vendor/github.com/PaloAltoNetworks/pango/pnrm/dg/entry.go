package dg

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a device group.
//
// Devices is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
	Name        string
	Description string
	Devices     map[string][]string

	raw map[string]string
}

// Copy copies the information from source's Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	if s.Devices == nil {
		o.Devices = nil
	} else {
		o.Devices = make(map[string][]string)
		for key, val := range s.Devices {
			if val == nil {
				o.Devices[key] = nil
			} else {
				list := make([]string, len(val))
				copy(list, val)
				o.Devices[key] = list
			}
		}
	}
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
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:        o.Name,
		Description: o.Description,
		Devices:     util.VsysEntToMap(o.Devices),
	}

	ans.raw = make(map[string]string)

	if o.Address != nil {
		ans.raw["address"] = util.CleanRawXml(o.Address.Text)
	}
	if o.AddressGroup != nil {
		ans.raw["addressGroup"] = util.CleanRawXml(o.AddressGroup.Text)
	}
	if o.Application != nil {
		ans.raw["application"] = util.CleanRawXml(o.Application.Text)
	}
	if o.ApplicationFilter != nil {
		ans.raw["applicationFilter"] = util.CleanRawXml(o.ApplicationFilter.Text)
	}
	if o.ApplicationGroup != nil {
		ans.raw["applicationGroup"] = util.CleanRawXml(o.ApplicationGroup.Text)
	}
	if o.ApplicationStatus != nil {
		ans.raw["applicationStatus"] = util.CleanRawXml(o.ApplicationStatus.Text)
	}
	if o.ApplicationTag != nil {
		ans.raw["applicationTag"] = util.CleanRawXml(o.ApplicationTag.Text)
	}
	if o.AuthenticationObject != nil {
		ans.raw["authenticationObject"] = util.CleanRawXml(o.AuthenticationObject.Text)
	}
	if o.AuthorizationCode != nil {
		ans.raw["authorizationCode"] = util.CleanRawXml(o.AuthorizationCode.Text)
	}
	if o.DynamicUserGroup != nil {
		ans.raw["dynamicUserGroup"] = util.CleanRawXml(o.DynamicUserGroup.Text)
	}
	if o.EmailScheduler != nil {
		ans.raw["emailScheduler"] = util.CleanRawXml(o.EmailScheduler.Text)
	}
	if o.Edl != nil {
		ans.raw["edl"] = util.CleanRawXml(o.Edl.Text)
	}
	if o.LogSettings != nil {
		ans.raw["logSettings"] = util.CleanRawXml(o.LogSettings.Text)
	}
	if o.MasterDevice != nil {
		ans.raw["masterDevice"] = util.CleanRawXml(o.MasterDevice.Text)
	}
	if o.PdfSummaryReport != nil {
		ans.raw["pdfSummaryReport"] = util.CleanRawXml(o.PdfSummaryReport.Text)
	}
	if o.PostRulebase != nil {
		ans.raw["postRulebase"] = util.CleanRawXml(o.PostRulebase.Text)
	}
	if o.PreRulebase != nil {
		ans.raw["preRulebase"] = util.CleanRawXml(o.PreRulebase.Text)
	}
	if o.ProfileGroup != nil {
		ans.raw["profileGroup"] = util.CleanRawXml(o.ProfileGroup.Text)
	}
	if o.Profiles != nil {
		ans.raw["profiles"] = util.CleanRawXml(o.Profiles.Text)
	}
	if o.Region != nil {
		ans.raw["region"] = util.CleanRawXml(o.Region.Text)
	}
	if o.ReportGroup != nil {
		ans.raw["reportGroup"] = util.CleanRawXml(o.ReportGroup.Text)
	}
	if o.Reports != nil {
		ans.raw["reports"] = util.CleanRawXml(o.Reports.Text)
	}
	if o.Schedule != nil {
		ans.raw["schedule"] = util.CleanRawXml(o.Schedule.Text)
	}
	if o.Service != nil {
		ans.raw["service"] = util.CleanRawXml(o.Service.Text)
	}
	if o.ServiceGroup != nil {
		ans.raw["serviceGroup"] = util.CleanRawXml(o.ServiceGroup.Text)
	}
	if o.Tag != nil {
		ans.raw["tag"] = util.CleanRawXml(o.Tag.Text)
	}
	if o.Threats != nil {
		ans.raw["threats"] = util.CleanRawXml(o.Threats.Text)
	}
	if o.ToSwVersion != nil {
		ans.raw["toSwVersion"] = util.CleanRawXml(o.ToSwVersion.Text)
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name            `xml:"entry"`
	Name        string              `xml:"name,attr"`
	Description string              `xml:"description,omitempty"`
	Devices     *util.VsysEntryType `xml:"devices"`

	Address              *util.RawXml `xml:"address"`
	AddressGroup         *util.RawXml `xml:"address-group"`
	Application          *util.RawXml `xml:"application"`
	ApplicationFilter    *util.RawXml `xml:"application-filter"`
	ApplicationGroup     *util.RawXml `xml:"application-group"`
	ApplicationStatus    *util.RawXml `xml:"application-status"`
	ApplicationTag       *util.RawXml `xml:"application-tag"`
	AuthenticationObject *util.RawXml `xml:"authentication-object"`
	AuthorizationCode    *util.RawXml `xml:"authorization-code"`
	DynamicUserGroup     *util.RawXml `xml:"dynamic-user-group"`
	EmailScheduler       *util.RawXml `xml:"email-scheduler"`
	Edl                  *util.RawXml `xml:"external-list"`
	LogSettings          *util.RawXml `xml:"log-settings"`
	MasterDevice         *util.RawXml `xml:"master-device"`
	PdfSummaryReport     *util.RawXml `xml:"pdf-summary-report"`
	PostRulebase         *util.RawXml `xml:"post-rulebase"`
	PreRulebase          *util.RawXml `xml:"pre-rulebase"`
	ProfileGroup         *util.RawXml `xml:"profile-group"`
	Profiles             *util.RawXml `xml:"profiles"`
	Region               *util.RawXml `xml:"region"`
	ReportGroup          *util.RawXml `xml:"report-group"`
	Reports              *util.RawXml `xml:"reports"`
	Schedule             *util.RawXml `xml:"schedule"`
	Service              *util.RawXml `xml:"service"`
	ServiceGroup         *util.RawXml `xml:"service-group"`
	Tag                  *util.RawXml `xml:"tag"`
	Threats              *util.RawXml `xml:"threats"`
	ToSwVersion          *util.RawXml `xml:"to-sw-version"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		Devices:     util.MapToVsysEnt(e.Devices),
	}

	if t, p := e.raw["address"]; p {
		ans.Address = &util.RawXml{t}
	}

	if t, p := e.raw["addressGroup"]; p {
		ans.AddressGroup = &util.RawXml{t}
	}

	if t, p := e.raw["application"]; p {
		ans.Application = &util.RawXml{t}
	}

	if t, p := e.raw["applicationFilter"]; p {
		ans.ApplicationFilter = &util.RawXml{t}
	}

	if t, p := e.raw["applicationGroup"]; p {
		ans.ApplicationGroup = &util.RawXml{t}
	}

	if t, p := e.raw["applicationStatus"]; p {
		ans.ApplicationStatus = &util.RawXml{t}
	}

	if t, p := e.raw["applicationTag"]; p {
		ans.ApplicationTag = &util.RawXml{t}
	}

	if t, p := e.raw["authenticationObject"]; p {
		ans.AuthenticationObject = &util.RawXml{t}
	}

	if t, p := e.raw["authorizationCode"]; p {
		ans.AuthorizationCode = &util.RawXml{t}
	}

	if t, p := e.raw["dynamicUserGroup"]; p {
		ans.DynamicUserGroup = &util.RawXml{t}
	}

	if t, p := e.raw["emailScheduler"]; p {
		ans.EmailScheduler = &util.RawXml{t}
	}

	if t, p := e.raw["edl"]; p {
		ans.Edl = &util.RawXml{t}
	}

	if t, p := e.raw["logSettings"]; p {
		ans.LogSettings = &util.RawXml{t}
	}

	if t, p := e.raw["masterDevice"]; p {
		ans.MasterDevice = &util.RawXml{t}
	}

	if t, p := e.raw["pdfSummaryReport"]; p {
		ans.PdfSummaryReport = &util.RawXml{t}
	}

	if t, p := e.raw["postRulebase"]; p {
		ans.PostRulebase = &util.RawXml{t}
	}

	if t, p := e.raw["preRulebase"]; p {
		ans.PreRulebase = &util.RawXml{t}
	}

	if t, p := e.raw["profileGroup"]; p {
		ans.ProfileGroup = &util.RawXml{t}
	}

	if t, p := e.raw["profiles"]; p {
		ans.Profiles = &util.RawXml{t}
	}

	if t, p := e.raw["region"]; p {
		ans.Region = &util.RawXml{t}
	}

	if t, p := e.raw["reportGroup"]; p {
		ans.ReportGroup = &util.RawXml{t}
	}

	if t, p := e.raw["reports"]; p {
		ans.Reports = &util.RawXml{t}
	}

	if t, p := e.raw["schedule"]; p {
		ans.Schedule = &util.RawXml{t}
	}

	if t, p := e.raw["service"]; p {
		ans.Service = &util.RawXml{t}
	}

	if t, p := e.raw["serviceGroup"]; p {
		ans.ServiceGroup = &util.RawXml{t}
	}

	if t, p := e.raw["tag"]; p {
		ans.Tag = &util.RawXml{t}
	}

	if t, p := e.raw["threats"]; p {
		ans.Threats = &util.RawXml{t}
	}

	if t, p := e.raw["toSwVersion"]; p {
		ans.ToSwVersion = &util.RawXml{t}
	}

	return ans
}
