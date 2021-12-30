package monitor

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of
// a monitor profile.
type Entry struct {
	Name      string
	Interval  int
	Threshold int
	Action    string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Interval = s.Interval
	o.Threshold = s.Threshold
	o.Action = s.Action
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
		Name:      o.Name,
		Interval:  o.Interval,
		Threshold: o.Threshold,
		Action:    o.Action,
	}

	return ans
}

type entry_v1 struct {
	XMLName   xml.Name `xml:"entry"`
	Name      string   `xml:"name,attr"`
	Interval  int      `xml:"interval,omitempty"`
	Threshold int      `xml:"threshold,omitempty"`
	Action    string   `xml:"action,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:      e.Name,
		Interval:  e.Interval,
		Threshold: e.Threshold,
		Action:    e.Action,
	}

	return ans
}
