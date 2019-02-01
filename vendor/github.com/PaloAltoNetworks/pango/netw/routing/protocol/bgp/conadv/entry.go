package conadv

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a BGP
// conditional advertisement.
type Entry struct {
    Name string
    Enable bool
    UsedBy []string

    raw map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Enable = s.Enable
    o.UsedBy = s.UsedBy
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
        UsedBy: util.MemToStr(o.Answer.UsedBy),
    }

    m := make(map[string] string)
    if o.Answer.NonExistFilters != nil {
        m["nf"] = util.CleanRawXml(o.Answer.NonExistFilters.Text)
    }
    if o.Answer.AdvertiseFilters != nil {
        m["af"] = util.CleanRawXml(o.Answer.AdvertiseFilters.Text)
    }
    if len(m) > 0 {
        ans.raw = m
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    UsedBy *util.MemberType `xml:"used-by"`
    NonExistFilters *util.RawXml `xml:"non-exist-filters"`
    AdvertiseFilters *util.RawXml `xml:"advertise-filters"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        UsedBy: util.StrToMem(e.UsedBy),
    }

    if text, present := e.raw["nf"]; present {
        ans.NonExistFilters = &util.RawXml{text}
    }
    if text, present := e.raw["af"]; present {
        ans.AdvertiseFilters = &util.RawXml{text}
    }

    return ans
}
