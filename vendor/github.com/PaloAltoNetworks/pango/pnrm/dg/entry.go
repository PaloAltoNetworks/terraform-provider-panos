package dg

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a virtual
// router.
type Entry struct {
    Name string
    Description string
    Devices []string
}

// Copy copies the information from source's Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Devices = s.Devices
}

/** Structs / functions for normalization. **/

type normalizer interface {
    Normalize() Entry
}

type container_v1 struct {
    Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Description: o.Answer.Description,
        Devices: util.EntToStr(o.Answer.Devices),
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description"`
    Devices *util.Entry `xml:"devices"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Devices: util.StrToEnt(e.Devices),
    }

    return ans
}
