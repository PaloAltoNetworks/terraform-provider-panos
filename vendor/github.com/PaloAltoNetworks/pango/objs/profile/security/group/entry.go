package group

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// security profile group.
type Entry struct {
	Name                    string
	AntivirusProfile        string
	AntiSpywareProfile      string
	VulnerabilityProfile    string
	UrlFilteringProfile     string
	FileBlockingProfile     string
	DataFilteringProfile    string
	WildfireAnalysisProfile string // PAN-OS 7.0
	GtpProfile              string // PAN-OS 8.0
	SctpProfile             string // PAN-OS 9.0
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.AntivirusProfile = s.AntivirusProfile
	o.AntiSpywareProfile = s.AntiSpywareProfile
	o.VulnerabilityProfile = s.VulnerabilityProfile
	o.UrlFilteringProfile = s.UrlFilteringProfile
	o.FileBlockingProfile = s.FileBlockingProfile
	o.DataFilteringProfile = s.DataFilteringProfile
	o.WildfireAnalysisProfile = s.WildfireAnalysisProfile
	o.GtpProfile = s.GtpProfile
	o.SctpProfile = s.SctpProfile
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
		AntivirusProfile:     util.MemToOneStr(o.AntivirusProfile),
		AntiSpywareProfile:   util.MemToOneStr(o.AntiSpywareProfile),
		VulnerabilityProfile: util.MemToOneStr(o.VulnerabilityProfile),
		UrlFilteringProfile:  util.MemToOneStr(o.UrlFilteringProfile),
		FileBlockingProfile:  util.MemToOneStr(o.FileBlockingProfile),
		DataFilteringProfile: util.MemToOneStr(o.DataFilteringProfile),
	}

	return ans
}

type entry_v1 struct {
	XMLName              xml.Name         `xml:"entry"`
	Name                 string           `xml:"name,attr"`
	AntivirusProfile     *util.MemberType `xml:"virus"`
	AntiSpywareProfile   *util.MemberType `xml:"spyware"`
	VulnerabilityProfile *util.MemberType `xml:"vulnerability"`
	UrlFilteringProfile  *util.MemberType `xml:"url-filtering"`
	FileBlockingProfile  *util.MemberType `xml:"file-blocking"`
	DataFilteringProfile *util.MemberType `xml:"data-filtering"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                 e.Name,
		AntivirusProfile:     util.OneStrToMem(e.AntivirusProfile),
		AntiSpywareProfile:   util.OneStrToMem(e.AntiSpywareProfile),
		VulnerabilityProfile: util.OneStrToMem(e.VulnerabilityProfile),
		UrlFilteringProfile:  util.OneStrToMem(e.UrlFilteringProfile),
		FileBlockingProfile:  util.OneStrToMem(e.FileBlockingProfile),
		DataFilteringProfile: util.OneStrToMem(e.DataFilteringProfile),
	}

	return ans
}

// PAN-OS 7.0
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
		Name:                    o.Name,
		AntivirusProfile:        util.MemToOneStr(o.AntivirusProfile),
		AntiSpywareProfile:      util.MemToOneStr(o.AntiSpywareProfile),
		VulnerabilityProfile:    util.MemToOneStr(o.VulnerabilityProfile),
		UrlFilteringProfile:     util.MemToOneStr(o.UrlFilteringProfile),
		FileBlockingProfile:     util.MemToOneStr(o.FileBlockingProfile),
		DataFilteringProfile:    util.MemToOneStr(o.DataFilteringProfile),
		WildfireAnalysisProfile: util.MemToOneStr(o.WildfireAnalysisProfile),
	}

	return ans
}

type entry_v2 struct {
	XMLName                 xml.Name         `xml:"entry"`
	Name                    string           `xml:"name,attr"`
	AntivirusProfile        *util.MemberType `xml:"virus"`
	AntiSpywareProfile      *util.MemberType `xml:"spyware"`
	VulnerabilityProfile    *util.MemberType `xml:"vulnerability"`
	UrlFilteringProfile     *util.MemberType `xml:"url-filtering"`
	FileBlockingProfile     *util.MemberType `xml:"file-blocking"`
	DataFilteringProfile    *util.MemberType `xml:"data-filtering"`
	WildfireAnalysisProfile *util.MemberType `xml:"wildfire-analysis"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                    e.Name,
		AntivirusProfile:        util.OneStrToMem(e.AntivirusProfile),
		AntiSpywareProfile:      util.OneStrToMem(e.AntiSpywareProfile),
		VulnerabilityProfile:    util.OneStrToMem(e.VulnerabilityProfile),
		UrlFilteringProfile:     util.OneStrToMem(e.UrlFilteringProfile),
		FileBlockingProfile:     util.OneStrToMem(e.FileBlockingProfile),
		DataFilteringProfile:    util.OneStrToMem(e.DataFilteringProfile),
		WildfireAnalysisProfile: util.OneStrToMem(e.WildfireAnalysisProfile),
	}

	return ans
}

// PAN-OS 8.0
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
		Name:                    o.Name,
		AntivirusProfile:        util.MemToOneStr(o.AntivirusProfile),
		AntiSpywareProfile:      util.MemToOneStr(o.AntiSpywareProfile),
		VulnerabilityProfile:    util.MemToOneStr(o.VulnerabilityProfile),
		UrlFilteringProfile:     util.MemToOneStr(o.UrlFilteringProfile),
		FileBlockingProfile:     util.MemToOneStr(o.FileBlockingProfile),
		DataFilteringProfile:    util.MemToOneStr(o.DataFilteringProfile),
		WildfireAnalysisProfile: util.MemToOneStr(o.WildfireAnalysisProfile),
		GtpProfile:              util.MemToOneStr(o.GtpProfile),
	}

	return ans
}

type entry_v3 struct {
	XMLName                 xml.Name         `xml:"entry"`
	Name                    string           `xml:"name,attr"`
	AntivirusProfile        *util.MemberType `xml:"virus"`
	AntiSpywareProfile      *util.MemberType `xml:"spyware"`
	VulnerabilityProfile    *util.MemberType `xml:"vulnerability"`
	UrlFilteringProfile     *util.MemberType `xml:"url-filtering"`
	FileBlockingProfile     *util.MemberType `xml:"file-blocking"`
	DataFilteringProfile    *util.MemberType `xml:"data-filtering"`
	WildfireAnalysisProfile *util.MemberType `xml:"wildfire-analysis"`
	GtpProfile              *util.MemberType `xml:"gtp"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:                    e.Name,
		AntivirusProfile:        util.OneStrToMem(e.AntivirusProfile),
		AntiSpywareProfile:      util.OneStrToMem(e.AntiSpywareProfile),
		VulnerabilityProfile:    util.OneStrToMem(e.VulnerabilityProfile),
		UrlFilteringProfile:     util.OneStrToMem(e.UrlFilteringProfile),
		FileBlockingProfile:     util.OneStrToMem(e.FileBlockingProfile),
		DataFilteringProfile:    util.OneStrToMem(e.DataFilteringProfile),
		WildfireAnalysisProfile: util.OneStrToMem(e.WildfireAnalysisProfile),
		GtpProfile:              util.OneStrToMem(e.GtpProfile),
	}

	return ans
}

// PAN-OS 9.0
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
		Name:                    o.Name,
		AntivirusProfile:        util.MemToOneStr(o.AntivirusProfile),
		AntiSpywareProfile:      util.MemToOneStr(o.AntiSpywareProfile),
		VulnerabilityProfile:    util.MemToOneStr(o.VulnerabilityProfile),
		UrlFilteringProfile:     util.MemToOneStr(o.UrlFilteringProfile),
		FileBlockingProfile:     util.MemToOneStr(o.FileBlockingProfile),
		DataFilteringProfile:    util.MemToOneStr(o.DataFilteringProfile),
		WildfireAnalysisProfile: util.MemToOneStr(o.WildfireAnalysisProfile),
		GtpProfile:              util.MemToOneStr(o.GtpProfile),
		SctpProfile:             util.MemToOneStr(o.SctpProfile),
	}

	return ans
}

type entry_v4 struct {
	XMLName                 xml.Name         `xml:"entry"`
	Name                    string           `xml:"name,attr"`
	AntivirusProfile        *util.MemberType `xml:"virus"`
	AntiSpywareProfile      *util.MemberType `xml:"spyware"`
	VulnerabilityProfile    *util.MemberType `xml:"vulnerability"`
	UrlFilteringProfile     *util.MemberType `xml:"url-filtering"`
	FileBlockingProfile     *util.MemberType `xml:"file-blocking"`
	DataFilteringProfile    *util.MemberType `xml:"data-filtering"`
	WildfireAnalysisProfile *util.MemberType `xml:"wildfire-analysis"`
	GtpProfile              *util.MemberType `xml:"gtp"`
	SctpProfile             *util.MemberType `xml:"sctp"`
}

func specify_v4(e Entry) interface{} {
	ans := entry_v4{
		Name:                    e.Name,
		AntivirusProfile:        util.OneStrToMem(e.AntivirusProfile),
		AntiSpywareProfile:      util.OneStrToMem(e.AntiSpywareProfile),
		VulnerabilityProfile:    util.OneStrToMem(e.VulnerabilityProfile),
		UrlFilteringProfile:     util.OneStrToMem(e.UrlFilteringProfile),
		FileBlockingProfile:     util.OneStrToMem(e.FileBlockingProfile),
		DataFilteringProfile:    util.OneStrToMem(e.DataFilteringProfile),
		WildfireAnalysisProfile: util.OneStrToMem(e.WildfireAnalysisProfile),
		GtpProfile:              util.OneStrToMem(e.GtpProfile),
		SctpProfile:             util.OneStrToMem(e.SctpProfile),
	}

	return ans
}
