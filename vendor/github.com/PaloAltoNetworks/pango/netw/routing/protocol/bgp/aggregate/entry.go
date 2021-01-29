package aggregate

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a BGP
// address aggregation policy.
type Entry struct {
	Name                   string
	Prefix                 string
	Enable                 bool
	Summary                bool
	AsSet                  bool
	LocalPreference        string
	Med                    string
	Weight                 int
	NextHop                string
	Origin                 string
	AsPathLimit            int
	AsPathType             string
	AsPathValue            string
	CommunityType          string
	CommunityValue         string
	ExtendedCommunityType  string
	ExtendedCommunityValue string

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Prefix = s.Prefix
	o.Enable = s.Enable
	o.Summary = s.Summary
	o.AsSet = s.AsSet
	o.LocalPreference = s.LocalPreference
	o.Med = s.Med
	o.Weight = s.Weight
	o.NextHop = s.NextHop
	o.Origin = s.Origin
	o.AsPathLimit = s.AsPathLimit
	o.AsPathType = s.AsPathType
	o.AsPathValue = s.AsPathValue
	o.CommunityType = s.CommunityType
	o.CommunityValue = s.CommunityValue
	o.ExtendedCommunityType = s.ExtendedCommunityType
	o.ExtendedCommunityValue = s.ExtendedCommunityValue
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
		Name:    o.Name,
		Prefix:  o.Prefix,
		Enable:  util.AsBool(o.Enable),
		Summary: util.AsBool(o.Summary),
		AsSet:   util.AsBool(o.AsSet),
	}

	if o.Options != nil {
		ans.LocalPreference = o.Options.LocalPreference
		ans.Med = o.Options.Med
		ans.Weight = o.Options.Weight
		ans.NextHop = o.Options.NextHop
		ans.Origin = o.Options.Origin
		ans.AsPathLimit = o.Options.AsPathLimit

		if o.Options.AsPath != nil {
			if o.Options.AsPath.None != nil {
				ans.AsPathType = AsPathTypeNone
			} else if o.Options.AsPath.Remove != nil {
				ans.AsPathType = AsPathTypeRemove
			} else if o.Options.AsPath.Prepend != "" {
				ans.AsPathType = AsPathTypePrepend
				ans.AsPathValue = o.Options.AsPath.Prepend
			} else if o.Options.AsPath.RemoveAndPrepend != "" {
				ans.AsPathType = AsPathTypeRemoveAndPrepend
				ans.AsPathValue = o.Options.AsPath.RemoveAndPrepend
			}
		}

		if o.Options.Community != nil {
			if o.Options.Community.None != nil {
				ans.CommunityType = CommunityTypeNone
			} else if o.Options.Community.RemoveAll != nil {
				ans.CommunityType = CommunityTypeRemoveAll
			} else if o.Options.Community.RemoveRegex != "" {
				ans.CommunityType = CommunityTypeRemoveRegex
				ans.CommunityValue = o.Options.Community.RemoveRegex
			} else if o.Options.Community.Append != nil {
				ans.CommunityType = CommunityTypeAppend
				ans.CommunityValue = util.MemToOneStr(o.Options.Community.Append)
			} else if o.Options.Community.Overwrite != nil {
				ans.CommunityType = CommunityTypeOverwrite
				ans.CommunityValue = util.MemToOneStr(o.Options.Community.Overwrite)
			}
		}

		if o.Options.ExtendedCommunity != nil {
			if o.Options.ExtendedCommunity.None != nil {
				ans.ExtendedCommunityType = CommunityTypeNone
			} else if o.Options.ExtendedCommunity.RemoveAll != nil {
				ans.ExtendedCommunityType = CommunityTypeRemoveAll
			} else if o.Options.ExtendedCommunity.RemoveRegex != "" {
				ans.ExtendedCommunityType = CommunityTypeRemoveRegex
				ans.ExtendedCommunityValue = o.Options.ExtendedCommunity.RemoveRegex
			} else if o.Options.ExtendedCommunity.Append != nil {
				ans.ExtendedCommunityType = CommunityTypeAppend
				ans.ExtendedCommunityValue = util.MemToOneStr(o.Options.ExtendedCommunity.Append)
			} else if o.Options.ExtendedCommunity.Overwrite != nil {
				ans.ExtendedCommunityType = CommunityTypeOverwrite
				ans.ExtendedCommunityValue = util.MemToOneStr(o.Options.ExtendedCommunity.Overwrite)
			}
		}
	}

	m := make(map[string]string)
	if o.SuppressFilters != nil {
		m["sf"] = util.CleanRawXml(o.SuppressFilters.Text)
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
	XMLName          xml.Name     `xml:"entry"`
	Name             string       `xml:"name,attr"`
	Prefix           string       `xml:"prefix"`
	Enable           string       `xml:"enable"`
	Summary          string       `xml:"summary"`
	AsSet            string       `xml:"as-set"`
	Options          *options     `xml:"aggregate-route-attributes"`
	SuppressFilters  *util.RawXml `xml:"suppress-filters"`
	AdvertiseFilters *util.RawXml `xml:"advertise-filters"`
}

type options struct {
	LocalPreference   string  `xml:"local-preference,omitempty"`
	Med               string  `xml:"med,omitempty"`
	Weight            int     `xml:"weight,omitempty"`
	NextHop           string  `xml:"nexthop,omitempty"`
	Origin            string  `xml:"origin,omitempty"`
	AsPathLimit       int     `xml:"as-path-limit,omitempty"`
	AsPath            *asPath `xml:"as-path"`
	Community         *comm   `xml:"community"`
	ExtendedCommunity *comm   `xml:"extended-community"`
}

type asPath struct {
	None             *string `xml:"none"`
	Remove           *string `xml:"remove"`
	Prepend          string  `xml:"prepend,omitempty"`
	RemoveAndPrepend string  `xml:"remove-and-prepend,omitempty"`
}

type comm struct {
	None        *string          `xml:"none"`
	RemoveAll   *string          `xml:"remove-all"`
	RemoveRegex string           `xml:"remove-regex,omitempty"`
	Append      *util.MemberType `xml:"append"`
	Overwrite   *util.MemberType `xml:"overwrite"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:    e.Name,
		Prefix:  e.Prefix,
		Enable:  util.YesNo(e.Enable),
		Summary: util.YesNo(e.Summary),
		AsSet:   util.YesNo(e.AsSet),
	}
	s := ""

	if e.LocalPreference != "" || e.Med != "" || e.Weight != 0 || e.NextHop != "" || e.Origin != "" || e.AsPathLimit != 0 || e.AsPathType != "" || e.CommunityType != "" || e.ExtendedCommunityType != "" {
		ans.Options = &options{
			LocalPreference: e.LocalPreference,
			Med:             e.Med,
			Weight:          e.Weight,
			NextHop:         e.NextHop,
			Origin:          e.Origin,
			AsPathLimit:     e.AsPathLimit,
		}

		switch e.AsPathType {
		case AsPathTypeNone:
			ans.Options.AsPath = &asPath{
				None: &s,
			}
		case AsPathTypeRemove:
			ans.Options.AsPath = &asPath{
				Remove: &s,
			}
		case AsPathTypePrepend:
			ans.Options.AsPath = &asPath{
				Prepend: e.AsPathValue,
			}
		case AsPathTypeRemoveAndPrepend:
			ans.Options.AsPath = &asPath{
				RemoveAndPrepend: e.AsPathValue,
			}
		}

		switch e.CommunityType {
		case CommunityTypeNone:
			ans.Options.Community = &comm{
				None: &s,
			}
		case CommunityTypeRemoveAll:
			ans.Options.Community = &comm{
				RemoveAll: &s,
			}
		case CommunityTypeRemoveRegex:
			ans.Options.Community = &comm{
				RemoveRegex: e.CommunityValue,
			}
		case CommunityTypeAppend:
			ans.Options.Community = &comm{
				Append: util.OneStrToMem(e.CommunityValue),
			}
		case CommunityTypeOverwrite:
			ans.Options.Community = &comm{
				Overwrite: util.OneStrToMem(e.CommunityValue),
			}
		}

		switch e.ExtendedCommunityType {
		case CommunityTypeNone:
			ans.Options.ExtendedCommunity = &comm{
				None: &s,
			}
		case CommunityTypeRemoveAll:
			ans.Options.ExtendedCommunity = &comm{
				RemoveAll: &s,
			}
		case CommunityTypeRemoveRegex:
			ans.Options.ExtendedCommunity = &comm{
				RemoveRegex: e.ExtendedCommunityValue,
			}
		case CommunityTypeAppend:
			ans.Options.ExtendedCommunity = &comm{
				Append: util.OneStrToMem(e.ExtendedCommunityValue),
			}
		case CommunityTypeOverwrite:
			ans.Options.ExtendedCommunity = &comm{
				Overwrite: util.OneStrToMem(e.ExtendedCommunityValue),
			}
		}
	}

	if text, present := e.raw["sf"]; present {
		ans.SuppressFilters = &util.RawXml{text}
	}
	if text, present := e.raw["af"]; present {
		ans.AdvertiseFilters = &util.RawXml{text}
	}

	return ans
}
