package exp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an OSPF
// export rule.
type Entry struct {
	Name     string
	PathType string
	Tag      string
	Metric   int
}

func (o *Entry) Copy(s Entry) {
	o.PathType = s.PathType
	o.Tag = s.Tag
	o.Metric = s.Metric
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
		PathType: o.PathType,
		Tag:      o.Tag,
		Metric:   o.Metric,
	}

	return ans
}

type entry_v1 struct {
	XMLName  xml.Name `xml:"entry"`
	Name     string   `xml:"name,attr"`
	PathType string   `xml:"new-path-type,omitempty"`
	Tag      string   `xml:"new-tag,omitempty"`
	Metric   int      `xml:"metric,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:     e.Name,
		PathType: e.PathType,
		Tag:      e.Tag,
		Metric:   e.Metric,
	}

	return ans
}
