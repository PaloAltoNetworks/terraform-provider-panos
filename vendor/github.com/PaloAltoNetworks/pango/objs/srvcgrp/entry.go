package srvcgrp

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a service
// group.
type Entry struct {
    Name string
    Services []string
    Tags []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Services = s.Services
    o.Tags = s.Tags
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
        Services: util.MemToStr(o.Answer.Services),
        Tags: util.MemToStr(o.Answer.Tags),
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Services *util.MemberType `xml:"members"`
    Tags *util.MemberType `xml:"tag"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Services: util.StrToMem(e.Services),
        Tags: util.StrToMem(e.Tags),
    }

    return ans
}
