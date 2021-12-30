package iface

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an OSPF
// area interface.
type Entry struct {
	Name               string
	Enable             bool
	Passive            bool
	LinkType           string
	Metric             int
	Priority           int
	HelloInterval      int
	DeadCounts         int
	RetransmitInterval int
	TransitDelay       int
	GraceRestartDelay  int
	AuthProfile        string
	Neighbors          []string // unordered; p2mp link type only
	BfdProfile         string
}

func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.Passive = s.Passive
	o.LinkType = s.LinkType
	o.Metric = s.Metric
	o.Priority = s.Priority
	o.HelloInterval = s.HelloInterval
	o.DeadCounts = s.DeadCounts
	o.RetransmitInterval = s.RetransmitInterval
	o.TransitDelay = s.TransitDelay
	o.AuthProfile = s.AuthProfile
	o.GraceRestartDelay = s.GraceRestartDelay
	if s.Neighbors == nil {
		o.Neighbors = nil
	} else {
		o.Neighbors = make([]string, len(s.Neighbors))
		copy(o.Neighbors, s.Neighbors)
	}
	o.BfdProfile = s.BfdProfile
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
		Name:               o.Name,
		Enable:             util.AsBool(o.Enable),
		Passive:            util.AsBool(o.Passive),
		Metric:             o.Metric,
		Priority:           o.Priority,
		HelloInterval:      o.HelloInterval,
		DeadCounts:         o.DeadCounts,
		RetransmitInterval: o.RetransmitInterval,
		TransitDelay:       o.TransitDelay,
		GraceRestartDelay:  o.GraceRestartDelay,
		AuthProfile:        o.AuthProfile,
		Neighbors:          util.EntToStr(o.Neighbors),
	}

	if o.LinkType != nil {
		if o.LinkType.Broadcast != nil {
			ans.LinkType = LinkTypeBroadcast
		} else if o.LinkType.PointToPoint != nil {
			ans.LinkType = LinkTypePointToPoint
		} else if o.LinkType.PointToMultiPoint != nil {
			ans.LinkType = LinkTypePointToMultiPoint
		}
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.BfdProfile
	}

	return ans
}

type entry_v1 struct {
	XMLName            xml.Name        `xml:"entry"`
	Name               string          `xml:"name,attr"`
	Enable             string          `xml:"enable"`
	Passive            string          `xml:"passive"`
	LinkType           *linktype       `xml:"link-type"`
	Metric             int             `xml:"metric,omitempty"`
	Priority           int             `xml:"priority"`
	HelloInterval      int             `xml:"hello-interval"`
	DeadCounts         int             `xml:"dead-counts,omitempty"`
	RetransmitInterval int             `xml:"retransmit-interval,omitempty"`
	TransitDelay       int             `xml:"transit-delay,omitempty"`
	GraceRestartDelay  int             `xml:"gr-delay,omitempty"`
	AuthProfile        string          `xml:"authentication,omitempty"`
	Neighbors          *util.EntryType `xml:"neighbor"`
	Bfd                *bfd            `xml:"bfd"`
}

func (e *entry_v1) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type localEntry_v1 entry_v1
	ans := localEntry_v1{
		Priority:           1,
		Metric:             10,
		HelloInterval:      10,
		DeadCounts:         4,
		RetransmitInterval: 5,
		TransitDelay:       1,
		GraceRestartDelay:  10,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v1(ans)
	return nil
}

type linktype struct {
	Broadcast         *string `xml:"broadcast"`
	PointToPoint      *string `xml:"p2p"`
	PointToMultiPoint *string `xml:"p2mp"`
}

type bfd struct {
	BfdProfile string `xml:"profile,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:               e.Name,
		Enable:             util.YesNo(e.Enable),
		Passive:            util.YesNo(e.Passive),
		Metric:             e.Metric,
		Priority:           e.Priority,
		HelloInterval:      e.HelloInterval,
		DeadCounts:         e.DeadCounts,
		RetransmitInterval: e.RetransmitInterval,
		TransitDelay:       e.TransitDelay,
		GraceRestartDelay:  e.GraceRestartDelay,
		AuthProfile:        e.AuthProfile,
		Neighbors:          util.StrToEnt(e.Neighbors),
	}

	s := ""

	switch e.LinkType {
	case LinkTypeBroadcast:
		ans.LinkType = &linktype{Broadcast: &s}
	case LinkTypePointToPoint:
		ans.LinkType = &linktype{PointToPoint: &s}
	case LinkTypePointToMultiPoint:
		ans.LinkType = &linktype{PointToMultiPoint: &s}
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{BfdProfile: e.BfdProfile}
	}

	return ans
}
