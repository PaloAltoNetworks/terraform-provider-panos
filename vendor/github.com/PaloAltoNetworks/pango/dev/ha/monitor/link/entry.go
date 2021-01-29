package link

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an HA
// link monitor group.
type Entry struct {
	Name             string
	Enable           bool
	FailureCondition string
	Interfaces       []string // unordered
}

func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.FailureCondition = s.FailureCondition
	if s.Interfaces == nil {
		o.Interfaces = nil
	} else {
		o.Interfaces = make([]string, len(s.Interfaces))
		copy(o.Interfaces, s.Interfaces)
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
		FailureCondition: o.FailureCondition,
		Interfaces:       util.MemToStr(o.Interfaces),
	}

	return ans
}

type entry_v1 struct {
	XMLName          xml.Name         `xml:"entry"`
	Name             string           `xml:"name,attr"`
	Enable           string           `xml:"enabled"`
	FailureCondition string           `xml:"failure-condition,omitempty"`
	Interfaces       *util.MemberType `xml:"interface,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:             e.Name,
		Enable:           util.YesNo(e.Enable),
		FailureCondition: e.FailureCondition,
		Interfaces:       util.StrToMem(e.Interfaces),
	}

	return ans
}
