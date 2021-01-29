package ipv4

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
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
	o.Types = s.Types
	o.Interfaces = s.Interfaces
	o.Destinations = s.Destinations
	o.NextHops = s.NextHops
	o.OspfPathTypes = s.OspfPathTypes
	o.OspfAreas = s.OspfAreas
	o.OspfTags = s.OspfTags
	o.BgpCommunities = s.BgpCommunities
	o.BgpExtendedCommunities = s.BgpExtendedCommunities
}

/** Structs / functions for this namespace. **/

type normalizer interface {
	Normalize() Entry
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
	ans := Entry{
		Name:     o.Answer.Name,
		Priority: o.Answer.Priority,
	}

	if o.Answer.Action.Redist != nil {
		ans.Action = ActionRedist
	} else if o.Answer.Action.NoRedist != nil {
		ans.Action = ActionNoRedist
	}

	if o.Answer.Filter != nil {
		ans.Types = util.MemToStr(o.Answer.Filter.Types)
		ans.Interfaces = util.MemToStr(o.Answer.Filter.Interfaces)
		ans.Destinations = util.MemToStr(o.Answer.Filter.Destinations)
		ans.NextHops = util.MemToStr(o.Answer.Filter.NextHops)

		if o.Answer.Filter.Ospf != nil {
			ans.OspfPathTypes = util.MemToStr(o.Answer.Filter.Ospf.OspfPathTypes)
			ans.OspfAreas = util.MemToStr(o.Answer.Filter.Ospf.OspfAreas)
			ans.OspfTags = util.MemToStr(o.Answer.Filter.Ospf.OspfTags)
		}

		if o.Answer.Filter.Bgp != nil {
			ans.BgpCommunities = util.MemToStr(o.Answer.Filter.Bgp.BgpCommunities)
			ans.BgpExtendedCommunities = util.MemToStr(o.Answer.Filter.Bgp.BgpExtendedCommunities)
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
