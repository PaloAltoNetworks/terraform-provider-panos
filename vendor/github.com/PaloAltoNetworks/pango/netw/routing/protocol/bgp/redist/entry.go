package redist

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// BGP redistribution rule.
type Entry struct {
	Name                 string
	Enable               bool
	AddressFamily        string
	RouteTable           string // 8.0+
	Metric               int
	SetOrigin            string
	SetMed               string
	SetLocalPreference   string
	SetAsPathLimit       int
	SetCommunity         []string
	SetExtendedCommunity []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.AddressFamily = s.AddressFamily
	o.RouteTable = s.RouteTable
	o.Metric = s.Metric
	o.SetOrigin = s.SetOrigin
	o.SetMed = s.SetMed
	o.SetLocalPreference = s.SetLocalPreference
	o.SetAsPathLimit = s.SetAsPathLimit
	o.SetCommunity = s.SetCommunity
	o.SetExtendedCommunity = s.SetExtendedCommunity
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
		Name:                 o.Name,
		Enable:               util.AsBool(o.Enable),
		Metric:               o.Metric,
		SetOrigin:            o.SetOrigin,
		SetMed:               o.SetMed,
		SetLocalPreference:   o.SetLocalPreference,
		SetAsPathLimit:       o.SetAsPathLimit,
		SetCommunity:         util.MemToStr(o.SetCommunity),
		SetExtendedCommunity: util.MemToStr(o.SetExtendedCommunity),
	}

	return ans
}

type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:                 o.Name,
		Enable:               util.AsBool(o.Enable),
		AddressFamily:        o.AddressFamily,
		RouteTable:           o.RouteTable,
		Metric:               o.Metric,
		SetOrigin:            o.SetOrigin,
		SetMed:               o.SetMed,
		SetLocalPreference:   o.SetLocalPreference,
		SetAsPathLimit:       o.SetAsPathLimit,
		SetCommunity:         util.MemToStr(o.SetCommunity),
		SetExtendedCommunity: util.MemToStr(o.SetExtendedCommunity),
	}

	return ans
}

type entry_v1 struct {
	XMLName              xml.Name         `xml:"entry"`
	Name                 string           `xml:"name,attr"`
	Enable               string           `xml:"enable"`
	Metric               int              `xml:"metric,omitempty"`
	SetOrigin            string           `xml:"set-origin,omitempty"`
	SetMed               string           `xml:"set-med,omitempty"`
	SetLocalPreference   string           `xml:"set-local-preference,omitempty"`
	SetAsPathLimit       int              `xml:"set-as-path-limit,omitempty"`
	SetCommunity         *util.MemberType `xml:"set-community"`
	SetExtendedCommunity *util.MemberType `xml:"set-extended-community"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                 e.Name,
		Enable:               util.YesNo(e.Enable),
		Metric:               e.Metric,
		SetOrigin:            e.SetOrigin,
		SetMed:               e.SetMed,
		SetLocalPreference:   e.SetLocalPreference,
		SetAsPathLimit:       e.SetAsPathLimit,
		SetCommunity:         util.StrToMem(e.SetCommunity),
		SetExtendedCommunity: util.StrToMem(e.SetExtendedCommunity),
	}

	return ans
}

type entry_v2 struct {
	XMLName              xml.Name         `xml:"entry"`
	Name                 string           `xml:"name,attr"`
	Enable               string           `xml:"enable"`
	AddressFamily        string           `xml:"address-family-identifier"`
	RouteTable           string           `xml:"route-table,omitempty"`
	Metric               int              `xml:"metric,omitempty"`
	SetOrigin            string           `xml:"set-origin,omitempty"`
	SetMed               string           `xml:"set-med,omitempty"`
	SetLocalPreference   string           `xml:"set-local-preference,omitempty"`
	SetAsPathLimit       int              `xml:"set-as-path-limit,omitempty"`
	SetCommunity         *util.MemberType `xml:"set-community"`
	SetExtendedCommunity *util.MemberType `xml:"set-extended-community"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                 e.Name,
		Enable:               util.YesNo(e.Enable),
		AddressFamily:        e.AddressFamily,
		RouteTable:           e.RouteTable,
		Metric:               e.Metric,
		SetOrigin:            e.SetOrigin,
		SetMed:               e.SetMed,
		SetLocalPreference:   e.SetLocalPreference,
		SetAsPathLimit:       e.SetAsPathLimit,
		SetCommunity:         util.StrToMem(e.SetCommunity),
		SetExtendedCommunity: util.StrToMem(e.SetExtendedCommunity),
	}

	return ans
}
