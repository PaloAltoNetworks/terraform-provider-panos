package vlink

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an OSPF
// area virtual link.
type Entry struct {
	Name               string
	Enable             bool
	NeighborId         string
	TransitAreaId      string
	HelloInterval      int
	DeadCounts         int
	RetransmitInterval int
	TransitDelay       int
	AuthProfile        string
	BfdProfile         string
}

func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.NeighborId = s.NeighborId
	o.TransitAreaId = s.TransitAreaId
	o.HelloInterval = s.HelloInterval
	o.DeadCounts = s.DeadCounts
	o.RetransmitInterval = s.RetransmitInterval
	o.TransitDelay = s.TransitDelay
	o.AuthProfile = s.AuthProfile
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
		NeighborId:         o.NeighborId,
		TransitAreaId:      o.TransitAreaId,
		HelloInterval:      o.HelloInterval,
		DeadCounts:         o.DeadCounts,
		RetransmitInterval: o.RetransmitInterval,
		TransitDelay:       o.TransitDelay,
		AuthProfile:        o.AuthProfile,
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.BfdProfile
	}

	return ans
}

type entry_v1 struct {
	XMLName            xml.Name `xml:"entry"`
	Name               string   `xml:"name,attr"`
	Enable             string   `xml:"enable"`
	NeighborId         string   `xml:"neighbor-id"`
	TransitAreaId      string   `xml:"transit-area-id"`
	HelloInterval      int      `xml:"hello-interval,omitempty"`
	DeadCounts         int      `xml:"dead-counts,omitempty"`
	RetransmitInterval int      `xml:"retransmit-interval,omitempty"`
	TransitDelay       int      `xml:"transit-delay,omitempty"`
	AuthProfile        string   `xml:"authentication,omitempty"`
	Bfd                *bfd     `xml:"bfd"`
}

type bfd struct {
	BfdProfile string `xml:"profile,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:               e.Name,
		Enable:             util.YesNo(e.Enable),
		NeighborId:         e.NeighborId,
		TransitAreaId:      e.TransitAreaId,
		HelloInterval:      e.HelloInterval,
		DeadCounts:         e.DeadCounts,
		RetransmitInterval: e.RetransmitInterval,
		TransitDelay:       e.TransitDelay,
		AuthProfile:        e.AuthProfile,
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{BfdProfile: e.BfdProfile}
	}

	return ans
}
