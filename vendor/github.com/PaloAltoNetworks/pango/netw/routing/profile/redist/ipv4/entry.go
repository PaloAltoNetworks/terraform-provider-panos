package ipv4

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a redist profile.
type Entry struct {
	Name                   string
	Priority               int
	Action                 string
	Types                  []string
	Interfaces             []string
	Destinations           []string
	NextHops               []string
	OspfPathTypes          []string
	OspfAreas              []string
	OspfTags               []string
	BgpCommunities         []string
	BgpExtendedCommunities []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Priority = s.Priority
	o.Action = s.Action
	if s.Types == nil {
		o.Types = nil
	} else {
		o.Types = make([]string, len(s.Types))
		copy(o.Types, s.Types)
	}
	if s.Interfaces == nil {
		o.Interfaces = nil
	} else {
		o.Interfaces = make([]string, len(s.Interfaces))
		copy(o.Interfaces, s.Interfaces)
	}
	if s.Destinations == nil {
		o.Destinations = nil
	} else {
		o.Destinations = make([]string, len(s.Destinations))
		copy(o.Destinations, s.Destinations)
	}
	if s.NextHops == nil {
		o.NextHops = nil
	} else {
		o.NextHops = make([]string, len(s.NextHops))
		copy(o.NextHops, s.NextHops)
	}
	if s.OspfPathTypes == nil {
		o.OspfPathTypes = nil
	} else {
		o.OspfPathTypes = make([]string, len(s.OspfPathTypes))
		copy(o.OspfPathTypes, s.OspfPathTypes)
	}
	if s.OspfAreas == nil {
		o.OspfAreas = nil
	} else {
		o.OspfAreas = make([]string, len(s.OspfAreas))
		copy(o.OspfAreas, s.OspfAreas)
	}
	if s.OspfTags == nil {
		o.OspfTags = nil
	} else {
		o.OspfTags = make([]string, len(s.OspfTags))
		copy(o.OspfTags, s.OspfTags)
	}
	if s.BgpCommunities == nil {
		o.BgpCommunities = nil
	} else {
		o.BgpCommunities = make([]string, len(s.BgpCommunities))
		copy(o.BgpCommunities, s.BgpCommunities)
	}
	if s.BgpExtendedCommunities == nil {
		o.BgpExtendedCommunities = nil
	} else {
		o.BgpExtendedCommunities = make([]string, len(s.BgpExtendedCommunities))
		copy(o.BgpExtendedCommunities, s.BgpExtendedCommunities)
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
		Name:     o.Name,
		Priority: o.Priority,
	}

	if o.Action.Redist != nil {
		ans.Action = ActionRedist
	} else if o.Action.NoRedist != nil {
		ans.Action = ActionNoRedist
	}

	if o.Filter != nil {
		ans.Types = util.MemToStr(o.Filter.Types)
		ans.Interfaces = util.MemToStr(o.Filter.Interfaces)
		ans.Destinations = util.MemToStr(o.Filter.Destinations)
		ans.NextHops = util.MemToStr(o.Filter.NextHops)

		if o.Filter.Ospf != nil {
			ans.OspfPathTypes = util.MemToStr(o.Filter.Ospf.OspfPathTypes)
			ans.OspfAreas = util.MemToStr(o.Filter.Ospf.OspfAreas)
			ans.OspfTags = util.MemToStr(o.Filter.Ospf.OspfTags)
		}

		if o.Filter.Bgp != nil {
			ans.BgpCommunities = util.MemToStr(o.Filter.Bgp.BgpCommunities)
			ans.BgpExtendedCommunities = util.MemToStr(o.Filter.Bgp.BgpExtendedCommunities)
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName  xml.Name `xml:"entry"`
	Name     string   `xml:"name,attr"`
	Priority int      `xml:"priority"`
	Action   act      `xml:"action"`
	Filter   *filter  `xml:"filter"`
}

type act struct {
	Redist   *string `xml:"redist"`
	NoRedist *string `xml:"no-redist"`
}

type filter struct {
	Types        *util.MemberType `xml:"type"`
	Interfaces   *util.MemberType `xml:"interface"`
	Destinations *util.MemberType `xml:"destination"`
	NextHops     *util.MemberType `xml:"nexthop"`
	Ospf         *ospf            `xml:"ospf"`
	Bgp          *bgp             `xml:"bgp"`
}

type ospf struct {
	OspfPathTypes *util.MemberType `xml:"path-type"`
	OspfAreas     *util.MemberType `xml:"area"`
	OspfTags      *util.MemberType `xml:"tag"`
}

type bgp struct {
	BgpCommunities         *util.MemberType `xml:"community"`
	BgpExtendedCommunities *util.MemberType `xml:"extended-community"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:     e.Name,
		Priority: e.Priority,
	}

	s := ""
	switch e.Action {
	case ActionRedist:
		ans.Action.Redist = &s
	case ActionNoRedist:
		ans.Action.NoRedist = &s
	}

	if len(e.Types) != 0 || len(e.Interfaces) != 0 || len(e.Destinations) != 0 || len(e.NextHops) != 0 || len(e.OspfPathTypes) != 0 || len(e.OspfAreas) != 0 || len(e.OspfTags) != 0 || len(e.BgpCommunities) != 0 || len(e.BgpExtendedCommunities) != 0 {
		f := &filter{
			Types:        util.StrToMem(e.Types),
			Interfaces:   util.StrToMem(e.Interfaces),
			Destinations: util.StrToMem(e.Destinations),
			NextHops:     util.StrToMem(e.NextHops),
		}

		if len(e.OspfPathTypes) != 0 || len(e.OspfAreas) != 0 || len(e.OspfTags) != 0 {
			f.Ospf = &ospf{
				OspfPathTypes: util.StrToMem(e.OspfPathTypes),
				OspfAreas:     util.StrToMem(e.OspfAreas),
				OspfTags:      util.StrToMem(e.OspfTags),
			}
		}

		if len(e.BgpCommunities) != 0 || len(e.BgpExtendedCommunities) != 0 {
			f.Bgp = &bgp{
				BgpCommunities:         util.StrToMem(e.BgpCommunities),
				BgpExtendedCommunities: util.StrToMem(e.BgpExtendedCommunities),
			}
		}

		ans.Filter = f
	}

	return ans
}
