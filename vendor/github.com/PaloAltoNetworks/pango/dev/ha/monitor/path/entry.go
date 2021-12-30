package path

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a HA
// path monitor group.
type Entry struct {
	Name             string
	Enable           bool
	SrcIp            string // virtual-wire, vlan
	FailureCondition string
	PingInterval     int
	PingCount        int
	DstIpGroups      []DstIpGroup
}

type DstIpGroup struct {
	Name             string
	Enable           bool
	FailureCondition string
	DstIps           []string // unordered
}

func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.SrcIp = s.SrcIp
	o.FailureCondition = s.FailureCondition
	o.PingInterval = s.PingInterval
	o.PingCount = s.PingCount
	if s.DstIpGroups == nil {
		o.DstIpGroups = nil
	} else {
		o.DstIpGroups = make([]DstIpGroup, len(s.DstIpGroups))
		for i := range s.DstIpGroups {
			o.DstIpGroups[i] = s.DstIpGroups[i]
			if s.DstIpGroups[i].DstIps != nil {
				o.DstIpGroups[i].DstIps = make([]string,
					len(s.DstIpGroups[i].DstIps))
				copy(o.DstIpGroups[i].DstIps, s.DstIpGroups[i].DstIps)
			}
		}
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
		Name:             o.Name,
		Enable:           util.AsBool(o.Enable),
		SrcIp:            o.SrcIp,
		FailureCondition: o.FailureCondition,
		PingInterval:     o.PingInterval,
		PingCount:        o.PingCount,
	}

	if o.DstIpGroups != nil {
		groups := make([]DstIpGroup, 0, len(o.DstIpGroups.Entries))
		for _, group := range o.DstIpGroups.Entries {
			groups = append(groups, DstIpGroup{
				Name:             group.Name,
				Enable:           util.AsBool(group.Enable),
				FailureCondition: group.FailureCondition,
				DstIps:           util.MemToStr(group.DstIps),
			})
		}
		ans.DstIpGroups = groups
	}

	return ans
}

type entry_v1 struct {
	XMLName          xml.Name    `xml:"entry"`
	Name             string      `xml:"name,attr"`
	Enable           string      `xml:"enabled"`
	SrcIp            string      `xml:"source-ip,omitempty"`
	FailureCondition string      `xml:"failure-condition,omitempty"`
	PingInterval     int         `xml:"ping-interval,omitempty"`
	PingCount        int         `xml:"ping-count,omitempty"`
	DstIpGroups      *dstIpGroup `xml:"destination-ip-group"`
}

type dstIpGroup struct {
	Entries []dstIpGroupEntry `xml:"entry"`
}

type dstIpGroupEntry struct {
	Name             string           `xml:"name,attr"`
	Enable           string           `xml:"enabled"`
	FailureCondition string           `xml:"failure-condition,omitempty"`
	DstIps           *util.MemberType `xml:"destination-ip"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:             e.Name,
		Enable:           util.YesNo(e.Enable),
		SrcIp:            e.SrcIp,
		FailureCondition: e.FailureCondition,
		PingInterval:     e.PingInterval,
		PingCount:        e.PingCount,
	}

	if len(e.DstIpGroups) > 0 {
		ans.DstIpGroups = &dstIpGroup{}
		groups := make([]dstIpGroupEntry, 0, len(e.DstIpGroups))
		for _, group := range e.DstIpGroups {
			groups = append(groups, dstIpGroupEntry{
				Name:             group.Name,
				Enable:           util.YesNo(group.Enable),
				FailureCondition: group.FailureCondition,
				DstIps:           util.StrToMem(group.DstIps),
			})
		}
		ans.DstIpGroups.Entries = groups
	}

	return ans
}
