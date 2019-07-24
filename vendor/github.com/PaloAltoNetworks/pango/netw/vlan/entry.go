package vlan

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a VLAN.
//
// Static MAC addresses are given as a map[string] string, where the key is
// the MAC address and the value is the interface it should be associated with.
type Entry struct {
    Name string
    VlanInterface string
    Interfaces []string // unordered
    StaticMacs map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry, copyMacs bool) {
    o.VlanInterface = s.VlanInterface
    o.Interfaces = s.Interfaces

    if copyMacs {
        o.StaticMacs = s.StaticMacs
    }
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
        Interfaces: util.MemToStr(o.Answer.Interfaces),
    }

    if o.Answer.Vi != nil {
        ans.VlanInterface = o.Answer.Vi.VlanInterface
    }

    if len(o.Answer.Mac.Entry) > 0 {
        ans.StaticMacs = make(map[string] string, len(o.Answer.Mac.Entry))
        for i := range o.Answer.Mac.Entry {
            ans.StaticMacs[o.Answer.Mac.Entry[i].Mac] = o.Answer.Mac.Entry[i].Interface
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Vi *vi `xml:"virtual-interface"`
    Interfaces *util.MemberType `xml:"interface"`
    Mac mac `xml:"mac"`
}

type vi struct {
    VlanInterface string `xml:"interface,omitempty"`
}

type mac struct {
    Entry []macList `xml:"entry"`
}

type macList struct {
    XMLName xml.Name `xml:"entry"`
    Mac string `xml:"name,attr"`
    Interface string `xml:"interface"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Interfaces: util.StrToMem(e.Interfaces),
    }

    if e.VlanInterface != "" {
        ans.Vi = &vi{
            VlanInterface: e.VlanInterface,
        }
    }

    i := 0
    ans.Mac.Entry = make([]macList, len(e.StaticMacs))
    for key := range e.StaticMacs {
        ans.Mac.Entry[i] = macList{Mac: key, Interface: e.StaticMacs[key]}
        i = i + 1
    }

    return ans
}
