package dampening

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a dampening
// profile.
type Entry struct {
    Name string
    Enable bool
    Cutoff float64
    Reuse float64
    MaxHoldTime int
    DecayHalfLifeReachable int
    DecayHalfLifeUnreachable int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Enable = s.Enable
    o.Cutoff = s.Cutoff
    o.Reuse = s.Reuse
    o.MaxHoldTime = s.MaxHoldTime
    o.DecayHalfLifeReachable = s.DecayHalfLifeReachable
    o.DecayHalfLifeUnreachable = s.DecayHalfLifeUnreachable
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
        Name: o.Answer.Name,
        Enable: util.AsBool(o.Answer.Enable),
        Cutoff: o.Answer.Cutoff,
        Reuse: o.Answer.Reuse,
        MaxHoldTime: o.Answer.MaxHoldTime,
        DecayHalfLifeReachable: o.Answer.DecayHalfLifeReachable,
        DecayHalfLifeUnreachable: o.Answer.DecayHalfLifeUnreachable,
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    Cutoff float64 `xml:"cutoff,omitempty"`
    Reuse float64 `xml:"reuse,omitempty"`
    MaxHoldTime int `xml:"max-hold-time,omitempty"`
    DecayHalfLifeReachable int `xml:"decay-half-life-reachable,omitempty"`
    DecayHalfLifeUnreachable int `xml:"decay-half-life-unreachable,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        Cutoff: e.Cutoff,
        Reuse: e.Reuse,
        MaxHoldTime: e.MaxHoldTime,
        DecayHalfLifeReachable: e.DecayHalfLifeReachable,
        DecayHalfLifeUnreachable: e.DecayHalfLifeUnreachable,
    }

    return ans
}
