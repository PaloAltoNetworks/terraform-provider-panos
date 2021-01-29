package url

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// URL filtering security profile.
type Entry struct {
	Name                      string
	Description               string
	DynamicUrl                bool     // Removed in 9.0
	ExpiredLicenseAction      bool     // Removed in 9.0
	BlockListAction           string   // Removed in 9.0
	BlockList                 []string // Removed in 9.0
	AllowList                 []string // Removed in 9.0
	AllowCategories           []string // ordered
	AlertCategories           []string // ordered
	BlockCategories           []string // ordered
	ContinueCategories        []string // ordered
	OverrideCategories        []string // ordered
	TrackContainerPage        bool
	LogContainerPageOnly      bool
	SafeSearchEnforcement     bool
	LogHttpHeaderXff          bool
	LogHttpHeaderUserAgent    bool
	LogHttpHeaderReferer      bool
	UcdMode                   string                 // 8.0
	UcdModeGroupMapping       string                 // 8.0
	UcdLogSeverity            string                 // 8.0
	UcdAllowCategories        []string               // 8.0, ordered
	UcdAlertCategories        []string               // 8.0, ordered
	UcdBlockCategories        []string               // 8.0, ordered
	UcdContinueCategories     []string               // 8.0, ordered
	HttpHeaderInsertions      []HttpHeaderInsertion  // 8.1
	MachineLearningModels     []MachineLearningModel // 10.0
	MachineLearningExceptions []string               // 10.0
}

type HttpHeaderInsertion struct {
	Name        string
	Type        string
	Domains     []string // ordered
	HttpHeaders []HttpHeader
}

// HttpHeader is an individual HTTP header.  In PAN-OS, the Name param is
// auto generated and look like "header-0", "header-1"..  If the Name param
// is an empty string, the name will be auto populated as appropriate.
type HttpHeader struct {
	Name   string
	Header string
	Value  string
	Log    bool
}

type MachineLearningModel struct {
	Model  string
	Action string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.DynamicUrl = s.DynamicUrl
	o.ExpiredLicenseAction = s.ExpiredLicenseAction
	o.BlockListAction = s.BlockListAction
	if s.BlockList == nil {
		o.BlockList = nil
	} else {
		o.BlockList = make([]string, len(s.BlockList))
		copy(o.BlockList, s.BlockList)
	}
	if s.AllowList == nil {
		o.AllowList = nil
	} else {
		o.AllowList = make([]string, len(s.AllowList))
		copy(o.AllowList, s.AllowList)
	}
	if s.AllowCategories == nil {
		o.AllowCategories = nil
	} else {
		o.AllowCategories = make([]string, len(s.AllowCategories))
		copy(o.AllowCategories, s.AllowCategories)
	}
	if s.AlertCategories == nil {
		o.AlertCategories = nil
	} else {
		o.AlertCategories = make([]string, len(s.AlertCategories))
		copy(o.AlertCategories, s.AlertCategories)
	}
	if s.BlockCategories == nil {
		o.BlockCategories = nil
	} else {
		o.BlockCategories = make([]string, len(s.BlockCategories))
		copy(o.BlockCategories, s.BlockCategories)
	}
	if s.ContinueCategories == nil {
		o.ContinueCategories = nil
	} else {
		o.ContinueCategories = make([]string, len(s.ContinueCategories))
		copy(o.ContinueCategories, s.ContinueCategories)
	}
	if s.OverrideCategories == nil {
		o.OverrideCategories = nil
	} else {
		o.OverrideCategories = make([]string, len(s.OverrideCategories))
		copy(o.OverrideCategories, s.OverrideCategories)
	}
	o.TrackContainerPage = s.TrackContainerPage
	o.LogContainerPageOnly = s.LogContainerPageOnly
	o.SafeSearchEnforcement = s.SafeSearchEnforcement
	o.LogHttpHeaderXff = s.LogHttpHeaderXff
	o.LogHttpHeaderUserAgent = s.LogHttpHeaderUserAgent
	o.LogHttpHeaderReferer = s.LogHttpHeaderReferer
	o.UcdMode = s.UcdMode
	o.UcdModeGroupMapping = s.UcdModeGroupMapping
	o.UcdLogSeverity = s.UcdLogSeverity
	if s.UcdAllowCategories == nil {
		o.UcdAllowCategories = nil
	} else {
		o.UcdAllowCategories = make([]string, len(s.UcdAllowCategories))
		copy(o.UcdAllowCategories, s.UcdAllowCategories)
	}
	if s.UcdAlertCategories == nil {
		o.UcdAlertCategories = nil
	} else {
		o.UcdAlertCategories = make([]string, len(s.UcdAlertCategories))
		copy(o.UcdAlertCategories, s.UcdAlertCategories)
	}
	if s.UcdBlockCategories == nil {
		o.UcdBlockCategories = nil
	} else {
		o.UcdBlockCategories = make([]string, len(s.UcdBlockCategories))
		copy(o.UcdBlockCategories, s.UcdBlockCategories)
	}
	if s.UcdContinueCategories == nil {
		o.UcdContinueCategories = nil
	} else {
		o.UcdContinueCategories = make([]string, len(s.UcdContinueCategories))
		copy(o.UcdContinueCategories, s.UcdContinueCategories)
	}
	if s.HttpHeaderInsertions == nil {
		o.HttpHeaderInsertions = nil
	} else {
		o.HttpHeaderInsertions = make([]HttpHeaderInsertion, 0, len(s.HttpHeaderInsertions))
		for _, hi := range s.HttpHeaderInsertions {
			item := HttpHeaderInsertion{
				Name: hi.Name,
				Type: hi.Type,
			}
			if hi.Domains != nil {
				item.Domains = make([]string, len(hi.Domains))
				copy(item.Domains, hi.Domains)
			}
			if hi.HttpHeaders != nil {
				item.HttpHeaders = make([]HttpHeader, 0, len(hi.HttpHeaders))
				for _, hh := range hi.HttpHeaders {
					item.HttpHeaders = append(item.HttpHeaders, HttpHeader{
						Name:   hh.Name,
						Header: hh.Header,
						Value:  hh.Value,
						Log:    hh.Log,
					})
				}
			}
			o.HttpHeaderInsertions = append(o.HttpHeaderInsertions, item)
		}
	}
	if s.MachineLearningModels == nil {
		o.MachineLearningModels = nil
	} else {
		o.MachineLearningModels = make([]MachineLearningModel, 0, len(s.MachineLearningModels))
		for _, mlm := range s.MachineLearningModels {
			o.MachineLearningModels = append(o.MachineLearningModels, MachineLearningModel{
				Model:  mlm.Model,
				Action: mlm.Action,
			})
		}
	}
	if s.MachineLearningExceptions == nil {
		o.MachineLearningExceptions = nil
	} else {
		o.MachineLearningExceptions = make([]string, len(s.MachineLearningExceptions))
		copy(o.MachineLearningExceptions, s.MachineLearningExceptions)
	}
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

// 7.1 and lower.
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
		Name:                   o.Name,
		Description:            o.Description,
		DynamicUrl:             util.AsBool(o.DynamicUrl),
		ExpiredLicenseAction:   util.AsBool(o.ExpiredLicenseAction),
		BlockListAction:        o.BlockListAction,
		BlockList:              util.MemToStr(o.BlockList),
		AllowList:              util.MemToStr(o.AllowList),
		AllowCategories:        util.MemToStr(o.AllowCategories),
		AlertCategories:        util.MemToStr(o.AlertCategories),
		BlockCategories:        util.MemToStr(o.BlockCategories),
		ContinueCategories:     util.MemToStr(o.ContinueCategories),
		OverrideCategories:     util.MemToStr(o.OverrideCategories),
		TrackContainerPage:     util.AsBool(o.TrackContainerPage),
		LogContainerPageOnly:   util.AsBool(o.LogContainerPageOnly),
		SafeSearchEnforcement:  util.AsBool(o.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.AsBool(o.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.AsBool(o.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.AsBool(o.LogHttpHeaderReferer),
	}

	return ans
}

type entry_v1 struct {
	XMLName                xml.Name         `xml:"entry"`
	Name                   string           `xml:"name,attr"`
	Description            string           `xml:"description,omitempty"`
	DynamicUrl             string           `xml:"dynamic-url"`
	ExpiredLicenseAction   string           `xml:"license-expired"`
	BlockListAction        string           `xml:"action"`
	BlockList              *util.MemberType `xml:"block-list"`
	AllowList              *util.MemberType `xml:"allow-list"`
	AllowCategories        *util.MemberType `xml:"allow"`
	AlertCategories        *util.MemberType `xml:"alert"`
	BlockCategories        *util.MemberType `xml:"block"`
	ContinueCategories     *util.MemberType `xml:"continue"`
	OverrideCategories     *util.MemberType `xml:"override"`
	TrackContainerPage     string           `xml:"enable-container-page"`
	LogContainerPageOnly   string           `xml:"log-container-page-only"`
	SafeSearchEnforcement  string           `xml:"safe-search-enforcement"`
	LogHttpHeaderXff       string           `xml:"log-http-hdr-xff"`
	LogHttpHeaderUserAgent string           `xml:"log-http-hdr-user-agent"`
	LogHttpHeaderReferer   string           `xml:"log-http-hdr-referer"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                   e.Name,
		Description:            e.Description,
		DynamicUrl:             util.YesNo(e.DynamicUrl),
		ExpiredLicenseAction:   util.YesNo(e.ExpiredLicenseAction),
		BlockListAction:        e.BlockListAction,
		BlockList:              util.StrToMem(e.BlockList),
		AllowList:              util.StrToMem(e.AllowList),
		AllowCategories:        util.StrToMem(e.AllowCategories),
		AlertCategories:        util.StrToMem(e.AlertCategories),
		BlockCategories:        util.StrToMem(e.BlockCategories),
		ContinueCategories:     util.StrToMem(e.ContinueCategories),
		OverrideCategories:     util.StrToMem(e.OverrideCategories),
		TrackContainerPage:     util.YesNo(e.TrackContainerPage),
		LogContainerPageOnly:   util.YesNo(e.LogContainerPageOnly),
		SafeSearchEnforcement:  util.YesNo(e.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.YesNo(e.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.YesNo(e.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.YesNo(e.LogHttpHeaderReferer),
	}

	return ans
}

type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v2) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:                   o.Name,
		Description:            o.Description,
		DynamicUrl:             util.AsBool(o.DynamicUrl),
		ExpiredLicenseAction:   util.AsBool(o.ExpiredLicenseAction),
		BlockListAction:        o.BlockListAction,
		BlockList:              util.MemToStr(o.BlockList),
		AllowList:              util.MemToStr(o.AllowList),
		AllowCategories:        util.MemToStr(o.AllowCategories),
		AlertCategories:        util.MemToStr(o.AlertCategories),
		BlockCategories:        util.MemToStr(o.BlockCategories),
		ContinueCategories:     util.MemToStr(o.ContinueCategories),
		OverrideCategories:     util.MemToStr(o.OverrideCategories),
		TrackContainerPage:     util.AsBool(o.TrackContainerPage),
		LogContainerPageOnly:   util.AsBool(o.LogContainerPageOnly),
		SafeSearchEnforcement:  util.AsBool(o.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.AsBool(o.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.AsBool(o.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.AsBool(o.LogHttpHeaderReferer),
	}

	if o.Ucd != nil {
		switch {
		case o.Ucd.Mode.UcdModeDisabled != nil:
			ans.UcdMode = UcdModeDisabled
		case o.Ucd.Mode.UcdModeIpUser != nil:
			ans.UcdMode = UcdModeIpUser
		case o.Ucd.Mode.UcdModeDomainCredentials != nil:
			ans.UcdMode = UcdModeDomainCredentials
		case o.Ucd.Mode.UcdModeGroupMapping != "":
			ans.UcdMode = UcdModeGroupMapping
			ans.UcdModeGroupMapping = o.Ucd.Mode.UcdModeGroupMapping
		}

		ans.UcdLogSeverity = o.Ucd.UcdLogSeverity
		ans.UcdAllowCategories = util.MemToStr(o.Ucd.UcdAllowCategories)
		ans.UcdAlertCategories = util.MemToStr(o.Ucd.UcdAlertCategories)
		ans.UcdBlockCategories = util.MemToStr(o.Ucd.UcdBlockCategories)
		ans.UcdContinueCategories = util.MemToStr(o.Ucd.UcdContinueCategories)
	}

	return ans
}

// 8.0
type entry_v2 struct {
	XMLName                xml.Name         `xml:"entry"`
	Name                   string           `xml:"name,attr"`
	Description            string           `xml:"description,omitempty"`
	DynamicUrl             string           `xml:"dynamic-url"`
	ExpiredLicenseAction   string           `xml:"license-expired"`
	BlockListAction        string           `xml:"action"`
	BlockList              *util.MemberType `xml:"block-list"`
	AllowList              *util.MemberType `xml:"allow-list"`
	AllowCategories        *util.MemberType `xml:"allow"`
	AlertCategories        *util.MemberType `xml:"alert"`
	BlockCategories        *util.MemberType `xml:"block"`
	ContinueCategories     *util.MemberType `xml:"continue"`
	OverrideCategories     *util.MemberType `xml:"override"`
	Ucd                    *creds           `xml:"credential-enforcement"`
	TrackContainerPage     string           `xml:"enable-container-page"`
	LogContainerPageOnly   string           `xml:"log-container-page-only"`
	SafeSearchEnforcement  string           `xml:"safe-search-enforcement"`
	LogHttpHeaderXff       string           `xml:"log-http-hdr-xff"`
	LogHttpHeaderUserAgent string           `xml:"log-http-hdr-user-agent"`
	LogHttpHeaderReferer   string           `xml:"log-http-hdr-referer"`
}

type creds struct {
	Mode                  credMode         `xml:"mode"`
	UcdLogSeverity        string           `xml:"log-severity"`
	UcdAllowCategories    *util.MemberType `xml:"allow"`
	UcdAlertCategories    *util.MemberType `xml:"alert"`
	UcdBlockCategories    *util.MemberType `xml:"block"`
	UcdContinueCategories *util.MemberType `xml:"continue"`
}

type credMode struct {
	UcdModeDisabled          *string `xml:"disabled"`
	UcdModeIpUser            *string `xml:"ip-user"`
	UcdModeDomainCredentials *string `xml:"domain-credentials"`
	UcdModeGroupMapping      string  `xml:"group-mapping,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                   e.Name,
		Description:            e.Description,
		DynamicUrl:             util.YesNo(e.DynamicUrl),
		ExpiredLicenseAction:   util.YesNo(e.ExpiredLicenseAction),
		BlockListAction:        e.BlockListAction,
		BlockList:              util.StrToMem(e.BlockList),
		AllowList:              util.StrToMem(e.AllowList),
		AllowCategories:        util.StrToMem(e.AllowCategories),
		AlertCategories:        util.StrToMem(e.AlertCategories),
		BlockCategories:        util.StrToMem(e.BlockCategories),
		ContinueCategories:     util.StrToMem(e.ContinueCategories),
		OverrideCategories:     util.StrToMem(e.OverrideCategories),
		TrackContainerPage:     util.YesNo(e.TrackContainerPage),
		LogContainerPageOnly:   util.YesNo(e.LogContainerPageOnly),
		SafeSearchEnforcement:  util.YesNo(e.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.YesNo(e.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.YesNo(e.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.YesNo(e.LogHttpHeaderReferer),
	}

	if e.UcdMode != "" || e.UcdLogSeverity != "" || len(e.UcdAllowCategories) != 0 || len(e.UcdAlertCategories) != 0 || len(e.UcdBlockCategories) != 0 || len(e.UcdContinueCategories) != 0 {
		s := ""
		var m credMode
		switch e.UcdMode {
		case UcdModeDisabled:
			m = credMode{
				UcdModeDisabled: &s,
			}
		case UcdModeIpUser:
			m = credMode{
				UcdModeIpUser: &s,
			}
		case UcdModeDomainCredentials:
			m = credMode{
				UcdModeDomainCredentials: &s,
			}
		case UcdModeGroupMapping:
			m = credMode{
				UcdModeGroupMapping: e.UcdModeGroupMapping,
			}
		}

		ans.Ucd = &creds{
			Mode:                  m,
			UcdLogSeverity:        e.UcdLogSeverity,
			UcdAllowCategories:    util.StrToMem(e.UcdAllowCategories),
			UcdAlertCategories:    util.StrToMem(e.UcdAlertCategories),
			UcdBlockCategories:    util.StrToMem(e.UcdBlockCategories),
			UcdContinueCategories: util.StrToMem(e.UcdContinueCategories),
		}
	}

	return ans
}

// 8.1
type container_v3 struct {
	Answer []entry_v3 `xml:"entry"`
}

func (o *container_v3) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v3) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v3) normalize() Entry {
	ans := Entry{
		Name:                   o.Name,
		Description:            o.Description,
		DynamicUrl:             util.AsBool(o.DynamicUrl),
		ExpiredLicenseAction:   util.AsBool(o.ExpiredLicenseAction),
		BlockListAction:        o.BlockListAction,
		BlockList:              util.MemToStr(o.BlockList),
		AllowList:              util.MemToStr(o.AllowList),
		AllowCategories:        util.MemToStr(o.AllowCategories),
		AlertCategories:        util.MemToStr(o.AlertCategories),
		BlockCategories:        util.MemToStr(o.BlockCategories),
		ContinueCategories:     util.MemToStr(o.ContinueCategories),
		OverrideCategories:     util.MemToStr(o.OverrideCategories),
		TrackContainerPage:     util.AsBool(o.TrackContainerPage),
		LogContainerPageOnly:   util.AsBool(o.LogContainerPageOnly),
		SafeSearchEnforcement:  util.AsBool(o.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.AsBool(o.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.AsBool(o.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.AsBool(o.LogHttpHeaderReferer),
	}

	if o.Ucd != nil {
		switch {
		case o.Ucd.Mode.UcdModeDisabled != nil:
			ans.UcdMode = UcdModeDisabled
		case o.Ucd.Mode.UcdModeIpUser != nil:
			ans.UcdMode = UcdModeIpUser
		case o.Ucd.Mode.UcdModeDomainCredentials != nil:
			ans.UcdMode = UcdModeDomainCredentials
		case o.Ucd.Mode.UcdModeGroupMapping != "":
			ans.UcdMode = UcdModeGroupMapping
			ans.UcdModeGroupMapping = o.Ucd.Mode.UcdModeGroupMapping
		}

		ans.UcdLogSeverity = o.Ucd.UcdLogSeverity
		ans.UcdAllowCategories = util.MemToStr(o.Ucd.UcdAllowCategories)
		ans.UcdAlertCategories = util.MemToStr(o.Ucd.UcdAlertCategories)
		ans.UcdBlockCategories = util.MemToStr(o.Ucd.UcdBlockCategories)
		ans.UcdContinueCategories = util.MemToStr(o.Ucd.UcdContinueCategories)
	}

	if o.Hhi != nil {
		ins := make([]HttpHeaderInsertion, 0, len(o.Hhi.Entries))
		for _, hhiObj := range o.Hhi.Entries {
			var headerList []HttpHeader
			if hhiObj.Types.Entry.Headers != nil {
				headerList = make([]HttpHeader, 0, len(hhiObj.Types.Entry.Headers.Entries))
				for _, hle := range hhiObj.Types.Entry.Headers.Entries {
					headerList = append(headerList, HttpHeader{
						Name:   hle.Name,
						Header: hle.Header,
						Value:  hle.Value,
						Log:    util.AsBool(hle.Log),
					})
				}
			}

			ins = append(ins, HttpHeaderInsertion{
				Name:        hhiObj.Name,
				Type:        hhiObj.Types.Entry.Type,
				Domains:     util.MemToStr(hhiObj.Types.Entry.Domains),
				HttpHeaders: headerList,
			})
		}

		ans.HttpHeaderInsertions = ins
	}

	return ans
}

type entry_v3 struct {
	XMLName                xml.Name         `xml:"entry"`
	Name                   string           `xml:"name,attr"`
	Description            string           `xml:"description,omitempty"`
	DynamicUrl             string           `xml:"dynamic-url"`
	ExpiredLicenseAction   string           `xml:"license-expired"`
	BlockListAction        string           `xml:"action"`
	BlockList              *util.MemberType `xml:"block-list"`
	AllowList              *util.MemberType `xml:"allow-list"`
	AllowCategories        *util.MemberType `xml:"allow"`
	AlertCategories        *util.MemberType `xml:"alert"`
	BlockCategories        *util.MemberType `xml:"block"`
	ContinueCategories     *util.MemberType `xml:"continue"`
	OverrideCategories     *util.MemberType `xml:"override"`
	Ucd                    *creds           `xml:"credential-enforcement"`
	TrackContainerPage     string           `xml:"enable-container-page"`
	LogContainerPageOnly   string           `xml:"log-container-page-only"`
	SafeSearchEnforcement  string           `xml:"safe-search-enforcement"`
	LogHttpHeaderXff       string           `xml:"log-http-hdr-xff"`
	LogHttpHeaderUserAgent string           `xml:"log-http-hdr-user-agent"`
	LogHttpHeaderReferer   string           `xml:"log-http-hdr-referer"`
	Hhi                    *hhi             `xml:"http-header-insertion"`
}

type hhi struct {
	Entries []hhiEntry `xml:"entry"`
}

type hhiEntry struct {
	Name  string  `xml:"name,attr"`
	Types hhiType `xml:"type"`
}

type hhiType struct {
	Entry hhiTypeEntry `xml:"entry"`
}

type hhiTypeEntry struct {
	Type    string           `xml:"name,attr"`
	Domains *util.MemberType `xml:"domains"`
	Headers *headers         `xml:"headers"`
}

type headers struct {
	Entries []headerEntry `xml:"entry"`
}

type headerEntry struct {
	Name   string `xml:"name,attr"`
	Header string `xml:"header"`
	Value  string `xml:"value"`
	Log    string `xml:"log"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:                   e.Name,
		Description:            e.Description,
		DynamicUrl:             util.YesNo(e.DynamicUrl),
		ExpiredLicenseAction:   util.YesNo(e.ExpiredLicenseAction),
		BlockListAction:        e.BlockListAction,
		BlockList:              util.StrToMem(e.BlockList),
		AllowList:              util.StrToMem(e.AllowList),
		AllowCategories:        util.StrToMem(e.AllowCategories),
		AlertCategories:        util.StrToMem(e.AlertCategories),
		BlockCategories:        util.StrToMem(e.BlockCategories),
		ContinueCategories:     util.StrToMem(e.ContinueCategories),
		OverrideCategories:     util.StrToMem(e.OverrideCategories),
		TrackContainerPage:     util.YesNo(e.TrackContainerPage),
		LogContainerPageOnly:   util.YesNo(e.LogContainerPageOnly),
		SafeSearchEnforcement:  util.YesNo(e.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.YesNo(e.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.YesNo(e.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.YesNo(e.LogHttpHeaderReferer),
	}

	if e.UcdMode != "" || e.UcdLogSeverity != "" || len(e.UcdAllowCategories) != 0 || len(e.UcdAlertCategories) != 0 || len(e.UcdBlockCategories) != 0 || len(e.UcdContinueCategories) != 0 {
		s := ""
		var m credMode
		switch e.UcdMode {
		case UcdModeDisabled:
			m = credMode{
				UcdModeDisabled: &s,
			}
		case UcdModeIpUser:
			m = credMode{
				UcdModeIpUser: &s,
			}
		case UcdModeDomainCredentials:
			m = credMode{
				UcdModeDomainCredentials: &s,
			}
		case UcdModeGroupMapping:
			m = credMode{
				UcdModeGroupMapping: e.UcdModeGroupMapping,
			}
		}

		ans.Ucd = &creds{
			Mode:                  m,
			UcdLogSeverity:        e.UcdLogSeverity,
			UcdAllowCategories:    util.StrToMem(e.UcdAllowCategories),
			UcdAlertCategories:    util.StrToMem(e.UcdAlertCategories),
			UcdBlockCategories:    util.StrToMem(e.UcdBlockCategories),
			UcdContinueCategories: util.StrToMem(e.UcdContinueCategories),
		}
	}

	if len(e.HttpHeaderInsertions) > 0 {
		hhiEntries := make([]hhiEntry, 0, len(e.HttpHeaderInsertions))

		for _, hhiObject := range e.HttpHeaderInsertions {
			var headersInst *headers

			if len(hhiObject.HttpHeaders) > 0 {
				list := make([]headerEntry, 0, len(hhiObject.HttpHeaders))
				for i := range hhiObject.HttpHeaders {
					var name string
					if hhiObject.HttpHeaders[i].Name == "" {
						name = fmt.Sprintf("header-%d", i)
					} else {
						name = hhiObject.HttpHeaders[i].Name
					}
					list = append(list, headerEntry{
						Name:   name,
						Header: hhiObject.HttpHeaders[i].Header,
						Value:  hhiObject.HttpHeaders[i].Value,
						Log:    util.YesNo(hhiObject.HttpHeaders[i].Log),
					})
				}

				headersInst = &headers{
					Entries: list,
				}
			}

			hhiEntries = append(hhiEntries, hhiEntry{
				Name: hhiObject.Name,
				Types: hhiType{
					Entry: hhiTypeEntry{
						Type:    hhiObject.Type,
						Domains: util.StrToMem(hhiObject.Domains),
						Headers: headersInst,
					},
				},
			})
		}

		ans.Hhi = &hhi{
			Entries: hhiEntries,
		}
	}

	return ans
}

// 9.0
type container_v4 struct {
	Answer []entry_v4 `xml:"entry"`
}

func (o *container_v4) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v4) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v4) normalize() Entry {
	ans := Entry{
		Name:                   o.Name,
		Description:            o.Description,
		AllowCategories:        util.MemToStr(o.AllowCategories),
		AlertCategories:        util.MemToStr(o.AlertCategories),
		BlockCategories:        util.MemToStr(o.BlockCategories),
		ContinueCategories:     util.MemToStr(o.ContinueCategories),
		OverrideCategories:     util.MemToStr(o.OverrideCategories),
		TrackContainerPage:     util.AsBool(o.TrackContainerPage),
		LogContainerPageOnly:   util.AsBool(o.LogContainerPageOnly),
		SafeSearchEnforcement:  util.AsBool(o.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.AsBool(o.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.AsBool(o.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.AsBool(o.LogHttpHeaderReferer),
	}

	if o.Ucd != nil {
		switch {
		case o.Ucd.Mode.UcdModeDisabled != nil:
			ans.UcdMode = UcdModeDisabled
		case o.Ucd.Mode.UcdModeIpUser != nil:
			ans.UcdMode = UcdModeIpUser
		case o.Ucd.Mode.UcdModeDomainCredentials != nil:
			ans.UcdMode = UcdModeDomainCredentials
		case o.Ucd.Mode.UcdModeGroupMapping != "":
			ans.UcdMode = UcdModeGroupMapping
			ans.UcdModeGroupMapping = o.Ucd.Mode.UcdModeGroupMapping
		}

		ans.UcdLogSeverity = o.Ucd.UcdLogSeverity
		ans.UcdAllowCategories = util.MemToStr(o.Ucd.UcdAllowCategories)
		ans.UcdAlertCategories = util.MemToStr(o.Ucd.UcdAlertCategories)
		ans.UcdBlockCategories = util.MemToStr(o.Ucd.UcdBlockCategories)
		ans.UcdContinueCategories = util.MemToStr(o.Ucd.UcdContinueCategories)
	}

	if o.Hhi != nil {
		ins := make([]HttpHeaderInsertion, 0, len(o.Hhi.Entries))
		for _, hhiObj := range o.Hhi.Entries {
			var headerList []HttpHeader
			if hhiObj.Types.Entry.Headers != nil {
				headerList = make([]HttpHeader, 0, len(hhiObj.Types.Entry.Headers.Entries))
				for _, hle := range hhiObj.Types.Entry.Headers.Entries {
					headerList = append(headerList, HttpHeader{
						Name:   hle.Name,
						Header: hle.Header,
						Value:  hle.Value,
						Log:    util.AsBool(hle.Log),
					})
				}
			}

			ins = append(ins, HttpHeaderInsertion{
				Name:        hhiObj.Name,
				Type:        hhiObj.Types.Entry.Type,
				Domains:     util.MemToStr(hhiObj.Types.Entry.Domains),
				HttpHeaders: headerList,
			})
		}

		ans.HttpHeaderInsertions = ins
	}

	return ans
}

type entry_v4 struct {
	XMLName                xml.Name         `xml:"entry"`
	Name                   string           `xml:"name,attr"`
	Description            string           `xml:"description,omitempty"`
	AllowCategories        *util.MemberType `xml:"allow"`
	AlertCategories        *util.MemberType `xml:"alert"`
	BlockCategories        *util.MemberType `xml:"block"`
	ContinueCategories     *util.MemberType `xml:"continue"`
	OverrideCategories     *util.MemberType `xml:"override"`
	Ucd                    *creds           `xml:"credential-enforcement"`
	TrackContainerPage     string           `xml:"enable-container-page"`
	LogContainerPageOnly   string           `xml:"log-container-page-only"`
	SafeSearchEnforcement  string           `xml:"safe-search-enforcement"`
	LogHttpHeaderXff       string           `xml:"log-http-hdr-xff"`
	LogHttpHeaderUserAgent string           `xml:"log-http-hdr-user-agent"`
	LogHttpHeaderReferer   string           `xml:"log-http-hdr-referer"`
	Hhi                    *hhi             `xml:"http-header-insertion"`
}

func specify_v4(e Entry) interface{} {
	ans := entry_v4{
		Name:                   e.Name,
		Description:            e.Description,
		AllowCategories:        util.StrToMem(e.AllowCategories),
		AlertCategories:        util.StrToMem(e.AlertCategories),
		BlockCategories:        util.StrToMem(e.BlockCategories),
		ContinueCategories:     util.StrToMem(e.ContinueCategories),
		OverrideCategories:     util.StrToMem(e.OverrideCategories),
		TrackContainerPage:     util.YesNo(e.TrackContainerPage),
		LogContainerPageOnly:   util.YesNo(e.LogContainerPageOnly),
		SafeSearchEnforcement:  util.YesNo(e.SafeSearchEnforcement),
		LogHttpHeaderXff:       util.YesNo(e.LogHttpHeaderXff),
		LogHttpHeaderUserAgent: util.YesNo(e.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:   util.YesNo(e.LogHttpHeaderReferer),
	}

	if e.UcdMode != "" || e.UcdLogSeverity != "" || len(e.UcdAllowCategories) != 0 || len(e.UcdAlertCategories) != 0 || len(e.UcdBlockCategories) != 0 || len(e.UcdContinueCategories) != 0 {
		s := ""
		var m credMode
		switch e.UcdMode {
		case UcdModeDisabled:
			m = credMode{
				UcdModeDisabled: &s,
			}
		case UcdModeIpUser:
			m = credMode{
				UcdModeIpUser: &s,
			}
		case UcdModeDomainCredentials:
			m = credMode{
				UcdModeDomainCredentials: &s,
			}
		case UcdModeGroupMapping:
			m = credMode{
				UcdModeGroupMapping: e.UcdModeGroupMapping,
			}
		}

		ans.Ucd = &creds{
			Mode:                  m,
			UcdLogSeverity:        e.UcdLogSeverity,
			UcdAllowCategories:    util.StrToMem(e.UcdAllowCategories),
			UcdAlertCategories:    util.StrToMem(e.UcdAlertCategories),
			UcdBlockCategories:    util.StrToMem(e.UcdBlockCategories),
			UcdContinueCategories: util.StrToMem(e.UcdContinueCategories),
		}
	}

	if len(e.HttpHeaderInsertions) > 0 {
		hhiEntries := make([]hhiEntry, 0, len(e.HttpHeaderInsertions))

		for _, hhiObject := range e.HttpHeaderInsertions {
			var headersInst *headers

			if len(hhiObject.HttpHeaders) > 0 {
				list := make([]headerEntry, 0, len(hhiObject.HttpHeaders))
				for i := range hhiObject.HttpHeaders {
					var name string
					if hhiObject.HttpHeaders[i].Name == "" {
						name = fmt.Sprintf("header-%d", i)
					} else {
						name = hhiObject.HttpHeaders[i].Name
					}
					list = append(list, headerEntry{
						Name:   name,
						Header: hhiObject.HttpHeaders[i].Header,
						Value:  hhiObject.HttpHeaders[i].Value,
						Log:    util.YesNo(hhiObject.HttpHeaders[i].Log),
					})
				}

				headersInst = &headers{
					Entries: list,
				}
			}

			hhiEntries = append(hhiEntries, hhiEntry{
				Name: hhiObject.Name,
				Types: hhiType{
					Entry: hhiTypeEntry{
						Type:    hhiObject.Type,
						Domains: util.StrToMem(hhiObject.Domains),
						Headers: headersInst,
					},
				},
			})
		}

		ans.Hhi = &hhi{
			Entries: hhiEntries,
		}
	}

	return ans
}

// 10.0
type container_v5 struct {
	Answer []entry_v5 `xml:"entry"`
}

func (o *container_v5) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v5) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v5) normalize() Entry {
	ans := Entry{
		Name:                      o.Name,
		Description:               o.Description,
		AllowCategories:           util.MemToStr(o.AllowCategories),
		AlertCategories:           util.MemToStr(o.AlertCategories),
		BlockCategories:           util.MemToStr(o.BlockCategories),
		ContinueCategories:        util.MemToStr(o.ContinueCategories),
		OverrideCategories:        util.MemToStr(o.OverrideCategories),
		TrackContainerPage:        util.AsBool(o.TrackContainerPage),
		LogContainerPageOnly:      util.AsBool(o.LogContainerPageOnly),
		SafeSearchEnforcement:     util.AsBool(o.SafeSearchEnforcement),
		LogHttpHeaderXff:          util.AsBool(o.LogHttpHeaderXff),
		LogHttpHeaderUserAgent:    util.AsBool(o.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:      util.AsBool(o.LogHttpHeaderReferer),
		MachineLearningExceptions: util.MemToStr(o.MachineLearningExceptions),
	}

	if o.Ucd != nil {
		switch {
		case o.Ucd.Mode.UcdModeDisabled != nil:
			ans.UcdMode = UcdModeDisabled
		case o.Ucd.Mode.UcdModeIpUser != nil:
			ans.UcdMode = UcdModeIpUser
		case o.Ucd.Mode.UcdModeDomainCredentials != nil:
			ans.UcdMode = UcdModeDomainCredentials
		case o.Ucd.Mode.UcdModeGroupMapping != "":
			ans.UcdMode = UcdModeGroupMapping
			ans.UcdModeGroupMapping = o.Ucd.Mode.UcdModeGroupMapping
		}

		ans.UcdLogSeverity = o.Ucd.UcdLogSeverity
		ans.UcdAllowCategories = util.MemToStr(o.Ucd.UcdAllowCategories)
		ans.UcdAlertCategories = util.MemToStr(o.Ucd.UcdAlertCategories)
		ans.UcdBlockCategories = util.MemToStr(o.Ucd.UcdBlockCategories)
		ans.UcdContinueCategories = util.MemToStr(o.Ucd.UcdContinueCategories)
	}

	if o.Hhi != nil {
		ins := make([]HttpHeaderInsertion, 0, len(o.Hhi.Entries))
		for _, hhiObj := range o.Hhi.Entries {
			var headerList []HttpHeader
			if hhiObj.Types.Entry.Headers != nil {
				headerList = make([]HttpHeader, 0, len(hhiObj.Types.Entry.Headers.Entries))
				for _, hle := range hhiObj.Types.Entry.Headers.Entries {
					headerList = append(headerList, HttpHeader{
						Name:   hle.Name,
						Header: hle.Header,
						Value:  hle.Value,
						Log:    util.AsBool(hle.Log),
					})
				}
			}

			ins = append(ins, HttpHeaderInsertion{
				Name:        hhiObj.Name,
				Type:        hhiObj.Types.Entry.Type,
				Domains:     util.MemToStr(hhiObj.Types.Entry.Domains),
				HttpHeaders: headerList,
			})
		}

		ans.HttpHeaderInsertions = ins
	}

	if o.MlModels != nil {
		listing := make([]MachineLearningModel, 0, len(o.MlModels.Entries))
		for _, model := range o.MlModels.Entries {
			listing = append(listing, MachineLearningModel{
				Model:  model.Model,
				Action: model.Action,
			})
		}

		ans.MachineLearningModels = listing
	}

	return ans
}

type entry_v5 struct {
	XMLName                   xml.Name         `xml:"entry"`
	Name                      string           `xml:"name,attr"`
	Description               string           `xml:"description,omitempty"`
	AllowCategories           *util.MemberType `xml:"allow"`
	AlertCategories           *util.MemberType `xml:"alert"`
	BlockCategories           *util.MemberType `xml:"block"`
	ContinueCategories        *util.MemberType `xml:"continue"`
	OverrideCategories        *util.MemberType `xml:"override"`
	Ucd                       *creds           `xml:"credential-enforcement"`
	TrackContainerPage        string           `xml:"enable-container-page"`
	LogContainerPageOnly      string           `xml:"log-container-page-only"`
	SafeSearchEnforcement     string           `xml:"safe-search-enforcement"`
	LogHttpHeaderXff          string           `xml:"log-http-hdr-xff"`
	LogHttpHeaderUserAgent    string           `xml:"log-http-hdr-user-agent"`
	LogHttpHeaderReferer      string           `xml:"log-http-hdr-referer"`
	Hhi                       *hhi             `xml:"http-header-insertion"`
	MachineLearningExceptions *util.MemberType `xml:"mlav-category-exception"`
	MlModels                  *mlmodels        `xml:"mlav-engine-urlbased-enabled"`
}

type mlmodels struct {
	Entries []mlmodel `xml:"entry"`
}

type mlmodel struct {
	Model  string `xml:"name,attr"`
	Action string `xml:"mlav-policy-action"`
}

func specify_v5(e Entry) interface{} {
	ans := entry_v5{
		Name:                      e.Name,
		Description:               e.Description,
		AllowCategories:           util.StrToMem(e.AllowCategories),
		AlertCategories:           util.StrToMem(e.AlertCategories),
		BlockCategories:           util.StrToMem(e.BlockCategories),
		ContinueCategories:        util.StrToMem(e.ContinueCategories),
		OverrideCategories:        util.StrToMem(e.OverrideCategories),
		TrackContainerPage:        util.YesNo(e.TrackContainerPage),
		LogContainerPageOnly:      util.YesNo(e.LogContainerPageOnly),
		SafeSearchEnforcement:     util.YesNo(e.SafeSearchEnforcement),
		LogHttpHeaderXff:          util.YesNo(e.LogHttpHeaderXff),
		LogHttpHeaderUserAgent:    util.YesNo(e.LogHttpHeaderUserAgent),
		LogHttpHeaderReferer:      util.YesNo(e.LogHttpHeaderReferer),
		MachineLearningExceptions: util.StrToMem(e.MachineLearningExceptions),
	}

	if e.UcdMode != "" || e.UcdLogSeverity != "" || len(e.UcdAllowCategories) != 0 || len(e.UcdAlertCategories) != 0 || len(e.UcdBlockCategories) != 0 || len(e.UcdContinueCategories) != 0 {
		s := ""
		var m credMode
		switch e.UcdMode {
		case UcdModeDisabled:
			m = credMode{
				UcdModeDisabled: &s,
			}
		case UcdModeIpUser:
			m = credMode{
				UcdModeIpUser: &s,
			}
		case UcdModeDomainCredentials:
			m = credMode{
				UcdModeDomainCredentials: &s,
			}
		case UcdModeGroupMapping:
			m = credMode{
				UcdModeGroupMapping: e.UcdModeGroupMapping,
			}
		}

		ans.Ucd = &creds{
			Mode:                  m,
			UcdLogSeverity:        e.UcdLogSeverity,
			UcdAllowCategories:    util.StrToMem(e.UcdAllowCategories),
			UcdAlertCategories:    util.StrToMem(e.UcdAlertCategories),
			UcdBlockCategories:    util.StrToMem(e.UcdBlockCategories),
			UcdContinueCategories: util.StrToMem(e.UcdContinueCategories),
		}
	}

	if len(e.HttpHeaderInsertions) > 0 {
		hhiEntries := make([]hhiEntry, 0, len(e.HttpHeaderInsertions))

		for _, hhiObject := range e.HttpHeaderInsertions {
			var headersInst *headers

			if len(hhiObject.HttpHeaders) > 0 {
				list := make([]headerEntry, 0, len(hhiObject.HttpHeaders))
				for i := range hhiObject.HttpHeaders {
					var name string
					if hhiObject.HttpHeaders[i].Name == "" {
						name = fmt.Sprintf("header-%d", i)
					} else {
						name = hhiObject.HttpHeaders[i].Name
					}
					list = append(list, headerEntry{
						Name:   name,
						Header: hhiObject.HttpHeaders[i].Header,
						Value:  hhiObject.HttpHeaders[i].Value,
						Log:    util.YesNo(hhiObject.HttpHeaders[i].Log),
					})
				}

				headersInst = &headers{
					Entries: list,
				}
			}

			hhiEntries = append(hhiEntries, hhiEntry{
				Name: hhiObject.Name,
				Types: hhiType{
					Entry: hhiTypeEntry{
						Type:    hhiObject.Type,
						Domains: util.StrToMem(hhiObject.Domains),
						Headers: headersInst,
					},
				},
			})
		}

		ans.Hhi = &hhi{
			Entries: hhiEntries,
		}
	}

	if len(e.MachineLearningModels) > 0 {
		listing := make([]mlmodel, 0, len(e.MachineLearningModels))
		for _, model := range e.MachineLearningModels {
			listing = append(listing, mlmodel{
				Model:  model.Model,
				Action: model.Action,
			})
		}

		ans.MlModels = &mlmodels{
			Entries: listing,
		}
	}

	return ans
}
