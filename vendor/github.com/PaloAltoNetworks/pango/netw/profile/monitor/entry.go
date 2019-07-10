package monitor

import (
    "encoding/xml"
)


// Entry is a normalized, version independent representation of a peer.
type Entry struct {
    Name string
    Interval int
    Threshold int
    Action string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Interval = s.Interval
    o.Threshold = s.Threshold
    o.Action = s.Action
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
        Interval: o.Answer.Interval,
        Threshold: o.Answer.Threshold,
        Action: o.Answer.Action,
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Interval int `xml:"interval,omitempty"`
    Threshold int `xml:"threshold,omitempty"`
    Action string `xml:"action,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Interval: e.Interval,
        Threshold: e.Threshold,
        Action: e.Action,
    }

    return ans
}
