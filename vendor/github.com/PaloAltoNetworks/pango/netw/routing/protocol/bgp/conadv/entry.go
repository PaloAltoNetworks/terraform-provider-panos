package conadv

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a BGP
// conditional advertisement.
type Entry struct {
	Name   string
	Enable bool
	UsedBy []string

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.UsedBy = s.UsedBy
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

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:   o.Name,
		Enable: util.AsBool(o.Enable),
		UsedBy: util.MemToStr(o.UsedBy),
	}

	m := make(map[string]string)
	if o.NonExistFilters != nil {
		m["nf"] = util.CleanRawXml(o.NonExistFilters.Text)
	}
	if o.AdvertiseFilters != nil {
		m["af"] = util.CleanRawXml(o.AdvertiseFilters.Text)
	}
	if len(m) > 0 {
		ans.raw = m
	}

	return ans
}

type entry_v1 struct {
	XMLName          xml.Name         `xml:"entry"`
	Name             string           `xml:"name,attr"`
	Enable           string           `xml:"enable"`
	UsedBy           *util.MemberType `xml:"used-by"`
	NonExistFilters  *util.RawXml     `xml:"non-exist-filters"`
	AdvertiseFilters *util.RawXml     `xml:"advertise-filters"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:   e.Name,
		Enable: util.YesNo(e.Enable),
		UsedBy: util.StrToMem(e.UsedBy),
	}

	if text, present := e.raw["nf"]; present {
		ans.NonExistFilters = &util.RawXml{text}
	}
	if text, present := e.raw["af"]; present {
		ans.AdvertiseFilters = &util.RawXml{text}
	}

	return ans
}
