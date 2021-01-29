package bfd

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a BFD profile.
type Entry struct {
	Name                string
	Mode                string
	MinimumTxInterval   int
	MinimumRxInterval   int
	DetectionMultiplier int
	HoldTime            int
	MinimumRxTtl        int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Mode = s.Mode
	o.MinimumTxInterval = s.MinimumTxInterval
	o.MinimumRxInterval = s.MinimumRxInterval
	o.DetectionMultiplier = s.DetectionMultiplier
	o.HoldTime = s.HoldTime
	o.MinimumRxTtl = s.MinimumRxTtl
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
		Name:                o.Name,
		Mode:                o.Mode,
		MinimumTxInterval:   o.MinimumTxInterval,
		MinimumRxInterval:   o.MinimumRxInterval,
		DetectionMultiplier: o.DetectionMultiplier,
		HoldTime:            o.HoldTime,
	}

	if o.Multihop != nil {
		ans.MinimumRxTtl = o.Multihop.MinimumRxTtl
	}

	return ans
}

type entry_v1 struct {
	XMLName             xml.Name  `xml:"entry"`
	Name                string    `xml:"name,attr"`
	Mode                string    `xml:"mode,omitempty"`
	MinimumTxInterval   int       `xml:"min-tx-interval,omitempty"`
	MinimumRxInterval   int       `xml:"min-rx-interval,omitempty"`
	DetectionMultiplier int       `xml:"detection-multiplier,omitempty"`
	HoldTime            int       `xml:"hold-time,omitempty"`
	Multihop            *multihop `xml:"multihop"`
}

type multihop struct {
	MinimumRxTtl int `xml:"min-received-ttl"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                e.Name,
		Mode:                e.Mode,
		MinimumTxInterval:   e.MinimumTxInterval,
		MinimumRxInterval:   e.MinimumRxInterval,
		DetectionMultiplier: e.DetectionMultiplier,
		HoldTime:            e.HoldTime,
	}

	if e.MinimumRxTtl != 0 {
		ans.Multihop = &multihop{e.MinimumRxTtl}
	}

	return ans
}
