package advertise

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a BGP
// aggregation advertisement filter.
type Entry struct {
	Name                   string
	Enable                 bool
	AsPathRegex            string
	CommunityRegex         string
	ExtendedCommunityRegex string
	Med                    string
	RouteTable             string // 8.0+
	AddressPrefix          map[string]bool
	NextHop                []string
	FromPeer               []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.AsPathRegex = s.AsPathRegex
	o.CommunityRegex = s.CommunityRegex
	o.ExtendedCommunityRegex = s.ExtendedCommunityRegex
	o.Med = s.Med
	o.RouteTable = s.RouteTable
	o.AddressPrefix = s.AddressPrefix
	o.NextHop = s.NextHop
	o.FromPeer = s.FromPeer
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
	}

	if o.Match != nil {
		ans.Med = o.Match.Med
		ans.NextHop = util.MemToStr(o.Match.NextHop)
		ans.FromPeer = util.MemToStr(o.Match.FromPeer)

		if o.Match.AsPathRegex != nil {
			ans.AsPathRegex = o.Match.AsPathRegex.Regex
		}

		if o.Match.CommunityRegex != nil {
			ans.CommunityRegex = o.Match.CommunityRegex.Regex
		}

		if o.Match.ExtendedCommunityRegex != nil {
			ans.ExtendedCommunityRegex = o.Match.ExtendedCommunityRegex.Regex
		}

		if o.Match.AddressPrefix != nil {
			m := make(map[string]bool)
			for _, v := range o.Match.AddressPrefix.Entry {
				m[v.Name] = util.AsBool(v.Exact)
			}
			ans.AddressPrefix = m
		}
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
		Name:   o.Name,
		Enable: util.AsBool(o.Enable),
	}

	if o.Match != nil {
		ans.Med = o.Match.Med
		ans.NextHop = util.MemToStr(o.Match.NextHop)
		ans.FromPeer = util.MemToStr(o.Match.FromPeer)
		ans.RouteTable = o.Match.RouteTable

		if o.Match.AsPathRegex != nil {
			ans.AsPathRegex = o.Match.AsPathRegex.Regex
		}

		if o.Match.CommunityRegex != nil {
			ans.CommunityRegex = o.Match.CommunityRegex.Regex
		}

		if o.Match.ExtendedCommunityRegex != nil {
			ans.ExtendedCommunityRegex = o.Match.ExtendedCommunityRegex.Regex
		}

		if o.Match.AddressPrefix != nil {
			m := make(map[string]bool)
			for _, v := range o.Match.AddressPrefix.Entry {
				m[v.Name] = util.AsBool(v.Exact)
			}
			ans.AddressPrefix = m
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name  `xml:"entry"`
	Name    string    `xml:"name,attr"`
	Enable  string    `xml:"enable"`
	Match   *match_v1 `xml:"match"`
}

type match_v1 struct {
	AsPathRegex            *regex           `xml:"as-path"`
	CommunityRegex         *regex           `xml:"community"`
	ExtendedCommunityRegex *regex           `xml:"extended-community"`
	Med                    string           `xml:"med,omitempty"`
	AddressPrefix          *addPre          `xml:"address-prefix"`
	NextHop                *util.MemberType `xml:"nexthop"`
	FromPeer               *util.MemberType `xml:"from-peer"`
}

type addPre struct {
	Entry []apEntry `xml:"entry"`
}

type apEntry struct {
	Name  string `xml:"name,attr"`
	Exact string `xml:"exact"`
}

type regex struct {
	Regex string `xml:"regex,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:   e.Name,
		Enable: util.YesNo(e.Enable),
	}

	if e.AsPathRegex != "" || e.CommunityRegex != "" || e.ExtendedCommunityRegex != "" || e.Med != "" || len(e.AddressPrefix) > 0 || len(e.NextHop) > 0 || len(e.FromPeer) > 0 {
		ans.Match = &match_v1{
			Med:      e.Med,
			NextHop:  util.StrToMem(e.NextHop),
			FromPeer: util.StrToMem(e.FromPeer),
		}

		if e.AsPathRegex != "" {
			ans.Match.AsPathRegex = &regex{
				Regex: e.AsPathRegex,
			}
		}

		if e.CommunityRegex != "" {
			ans.Match.CommunityRegex = &regex{
				Regex: e.CommunityRegex,
			}
		}

		if e.ExtendedCommunityRegex != "" {
			ans.Match.ExtendedCommunityRegex = &regex{
				Regex: e.ExtendedCommunityRegex,
			}
		}

		if len(e.AddressPrefix) > 0 {
			apList := make([]apEntry, 0, len(e.AddressPrefix))
			for k, v := range e.AddressPrefix {
				apList = append(apList, apEntry{
					Name:  k,
					Exact: util.YesNo(v),
				})
			}
			ans.Match.AddressPrefix = &addPre{apList}
		}
	}

	return ans
}

type entry_v2 struct {
	XMLName xml.Name  `xml:"entry"`
	Name    string    `xml:"name,attr"`
	Enable  string    `xml:"enable"`
	Match   *match_v2 `xml:"match"`
}

type match_v2 struct {
	AsPathRegex            *regex           `xml:"as-path"`
	CommunityRegex         *regex           `xml:"community"`
	ExtendedCommunityRegex *regex           `xml:"extended-community"`
	Med                    string           `xml:"med,omitempty"`
	RouteTable             string           `xml:"route-table,omitempty"`
	AddressPrefix          *addPre          `xml:"address-prefix"`
	NextHop                *util.MemberType `xml:"nexthop"`
	FromPeer               *util.MemberType `xml:"from-peer"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:   e.Name,
		Enable: util.YesNo(e.Enable),
	}

	if e.AsPathRegex != "" || e.CommunityRegex != "" || e.ExtendedCommunityRegex != "" || e.Med != "" || e.RouteTable != "" || len(e.AddressPrefix) > 0 || len(e.NextHop) > 0 || len(e.FromPeer) > 0 {
		ans.Match = &match_v2{
			Med:        e.Med,
			RouteTable: e.RouteTable,
			NextHop:    util.StrToMem(e.NextHop),
			FromPeer:   util.StrToMem(e.FromPeer),
		}

		if e.AsPathRegex != "" {
			ans.Match.AsPathRegex = &regex{
				Regex: e.AsPathRegex,
			}
		}

		if e.CommunityRegex != "" {
			ans.Match.CommunityRegex = &regex{
				Regex: e.CommunityRegex,
			}
		}

		if e.ExtendedCommunityRegex != "" {
			ans.Match.ExtendedCommunityRegex = &regex{
				Regex: e.ExtendedCommunityRegex,
			}
		}

		if len(e.AddressPrefix) > 0 {
			apList := make([]apEntry, 0, len(e.AddressPrefix))
			for k, v := range e.AddressPrefix {
				apList = append(apList, apEntry{
					Name:  k,
					Exact: util.YesNo(v),
				})
			}
			ans.Match.AddressPrefix = &addPre{apList}
		}
	}

	return ans
}
