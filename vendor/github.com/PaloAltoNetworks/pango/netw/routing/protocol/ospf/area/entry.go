package area

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an OSPF
// area.
type Entry struct {
	Name                  string
	Type                  string  // normal, stub, nssa
	AcceptSummary         bool    // stub, nssa
	DefaultRouteAdvertise bool    // stub, nssa
	AdvertiseMetric       int     // stub, nssa
	AdvertiseType         string  // nssa
	ExtRanges             []Range // nssa
	Ranges                []Range // area

	raw map[string]string
}

type Range struct {
	Network string
	Action  string
}

func (o *Entry) Copy(s Entry) {
	o.Type = s.Type
	o.AcceptSummary = s.AcceptSummary
	o.DefaultRouteAdvertise = s.DefaultRouteAdvertise
	o.AdvertiseMetric = s.AdvertiseMetric
	o.AdvertiseType = s.AdvertiseType
	if s.ExtRanges == nil {
		o.ExtRanges = nil
	} else {
		o.ExtRanges = make([]Range, len(s.ExtRanges))
		copy(o.ExtRanges, s.ExtRanges)
	}
	if s.Ranges == nil {
		o.Ranges = nil
	} else {
		o.Ranges = make([]Range, len(s.Ranges))
		copy(o.Ranges, s.Ranges)
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
		Name: o.Name,
	}

	if o.Ranges != nil && len(o.Ranges.Entry) > 0 {
		ans.Ranges = make([]Range, 0, len(o.Ranges.Entry))
		for i := range o.Ranges.Entry {
			entry := Range{
				Network: o.Ranges.Entry[i].Name,
			}
			switch {
			case o.Ranges.Entry[i].Advertise != nil:
				entry.Action = ActionAdvertise
			case o.Ranges.Entry[i].Suppress != nil:
				entry.Action = ActionSuppress
			}
			ans.Ranges = append(ans.Ranges, entry)
		}
	}

	switch {
	case o.Type.Normal != nil:
		ans.Type = TypeNormal
	case o.Type.Stub != nil:
		ans.Type = TypeStub
		ans.AcceptSummary = util.AsBool(o.Type.Stub.AcceptSummary)
		if o.Type.Stub.DefaultRoute != nil {
			if o.Type.Stub.DefaultRoute.Disable != nil {
				ans.DefaultRouteAdvertise = false
			}
			if o.Type.Stub.DefaultRoute.Advertise != nil {
				ans.DefaultRouteAdvertise = true
				ans.AdvertiseMetric = o.Type.Stub.DefaultRoute.Advertise.Metric
			}
		}
	case o.Type.Nssa != nil:
		ans.Type = TypeNssa
		ans.AcceptSummary = util.AsBool(o.Type.Nssa.AcceptSummary)
		if o.Type.Nssa.DefaultRoute != nil {
			if o.Type.Nssa.DefaultRoute.Disable != nil {
				ans.DefaultRouteAdvertise = false
			}
			if o.Type.Nssa.DefaultRoute.Advertise != nil {
				ans.DefaultRouteAdvertise = true
				ans.AdvertiseMetric = o.Type.Nssa.DefaultRoute.Advertise.Metric
				ans.AdvertiseType = o.Type.Nssa.DefaultRoute.Advertise.Type
			}
		}
		if o.Type.Nssa.Ranges != nil && len(o.Type.Nssa.Ranges.Entry) > 0 {
			ans.ExtRanges = make([]Range, 0, len(o.Type.Nssa.Ranges.Entry))
			for i := range o.Type.Nssa.Ranges.Entry {
				entry := Range{
					Network: o.Type.Nssa.Ranges.Entry[i].Name,
				}
				switch {
				case o.Type.Nssa.Ranges.Entry[i].Advertise != nil:
					entry.Action = ActionAdvertise
				case o.Type.Nssa.Ranges.Entry[i].Suppress != nil:
					entry.Action = ActionSuppress
				}
				ans.ExtRanges = append(ans.ExtRanges, entry)
			}
		}
	}

	raw := make(map[string]string)
	if o.Interface != nil {
		raw["int"] = util.CleanRawXml(o.Interface.Text)
	}
	if o.VirtualLink != nil {
		raw["link"] = util.CleanRawXml(o.VirtualLink.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name       `xml:"entry"`
	Name    string         `xml:"name,attr"`
	Type    areaType       `xml:"type"`
	Ranges  *advrangeEntry `xml:"range"`

	Interface   *util.RawXml `xml:"interface"`
	VirtualLink *util.RawXml `xml:"virtual-link"`
}

type areaType struct {
	Normal *string `xml:"normal"`
	Stub   *stub   `xml:"stub"`
	Nssa   *nssa   `xml:"nssa"`
}

type stub struct {
	DefaultRoute  *defroute `xml:"default-route"`
	AcceptSummary string    `xml:"accept-summary"`
}

type defroute struct {
	Advertise *advertise `xml:"advertise"`
	Disable   *string    `xml:"disable"`
}

type advertise struct {
	Metric int    `xml:"metric,omitempty"`
	Type   string `xml:"type,omitempty"`
}

type nssa struct {
	DefaultRoute  *defroute      `xml:"default-route"`
	AcceptSummary string         `xml:"accept-summary"`
	Ranges        *advrangeEntry `xml:"nssa-ext-range""`
}

type advrangeEntry struct {
	Entry []advrange `xml:"entry""`
}

type advrange struct {
	Name      string  `xml:"name,attr"`
	Advertise *string `xml:"advertise"`
	Suppress  *string `xml:"suppress"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	s := ""

	if len(e.Ranges) > 0 {
		ans.Ranges = &advrangeEntry{}
		ans.Ranges.Entry = make([]advrange, 0, len(e.Ranges))
		for i := range e.Ranges {
			entry := advrange{
				Name: e.Ranges[i].Network,
			}
			switch e.Ranges[i].Action {
			case ActionAdvertise:
				entry.Advertise = &s
			case ActionSuppress:
				entry.Suppress = &s
			}
			ans.Ranges.Entry = append(ans.Ranges.Entry, entry)
		}
	}

	switch e.Type {
	case TypeNormal:
		ans.Type = areaType{Normal: &s}
	case TypeStub:
		ans.Type = areaType{}
		ans.Type.Stub = &stub{
			DefaultRoute:  &defroute{},
			AcceptSummary: util.YesNo(e.AcceptSummary),
		}
		if e.DefaultRouteAdvertise {
			ans.Type.Stub.DefaultRoute.Advertise = &advertise{
				Metric: e.AdvertiseMetric,
			}
		} else {
			ans.Type.Stub.DefaultRoute.Disable = &s
		}
	case TypeNssa:
		ans.Type = areaType{}
		ans.Type.Nssa = &nssa{
			DefaultRoute:  &defroute{},
			AcceptSummary: util.YesNo(e.AcceptSummary),
		}
		if e.DefaultRouteAdvertise {
			ans.Type.Nssa.DefaultRoute.Advertise = &advertise{
				Metric: e.AdvertiseMetric,
				Type:   e.AdvertiseType,
			}
		} else {
			ans.Type.Nssa.DefaultRoute.Disable = &s
		}
		if len(e.ExtRanges) > 0 {
			ans.Type.Nssa.Ranges = &advrangeEntry{}
			ans.Type.Nssa.Ranges.Entry = make([]advrange, 0, len(e.ExtRanges))
			for i := range e.ExtRanges {
				entry := advrange{
					Name: e.ExtRanges[i].Network,
				}
				switch e.ExtRanges[i].Action {
				case ActionAdvertise:
					entry.Advertise = &s
				case ActionSuppress:
					entry.Suppress = &s
				}
				ans.Type.Nssa.Ranges.Entry = append(ans.Type.Nssa.Ranges.Entry, entry)
			}
		}
	}

	return ans
}
