package addrgrp

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of an address
// group.  The value set in DynamicMatch should be something like the following:
//
//  * 'tag1'
//  * 'tag1' or 'tag2' and 'tag3'
//
// The Tags param is for administrative tags for this address object
// group itself.
type Entry struct {
    Name string
    Description string
    StaticAddresses []string
    DynamicMatch string
    Tags []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.StaticAddresses = s.StaticAddresses
    o.DynamicMatch = s.DynamicMatch
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
        Description: o.Answer.Description,
        StaticAddresses: util.MemToStr(o.Answer.StaticAddresses),
        Tags: util.MemToStr(o.Answer.Tags),
    }
    if o.Answer.DynamicMatch != nil {
        ans.DynamicMatch = *o.Answer.DynamicMatch
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Description string `xml:"description"`
    StaticAddresses *util.Member `xml:"static"`
    DynamicMatch *string `xml:"dynamic>filter"`
    Tags *util.Member `xml:"tag"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        StaticAddresses: util.StrToMem(e.StaticAddresses),
        Tags: util.StrToMem(e.Tags),
    }
    if e.DynamicMatch != "" {
        ans.DynamicMatch = &e.DynamicMatch
    }

    return ans
}
