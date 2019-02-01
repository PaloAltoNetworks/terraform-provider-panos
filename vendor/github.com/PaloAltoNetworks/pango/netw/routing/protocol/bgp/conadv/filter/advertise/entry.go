package advertise

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a BGP
// conditional advertisement advertise filter.
type Entry struct {
    Name string
    Enable bool
    AsPathRegex string
    CommunityRegex string
    ExtendedCommunityRegex string
    Med string
    RouteTable string // 8.0+
    AddressPrefix []string
    NextHop []string
    FromPeer []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Enable = s.Enable
    o.AsPathRegex = s.AsPathRegex
    o.CommunityRegex = s.CommunityRegex
    o.ExtendedCommunityRegex = s.ExtendedCommunityRegex
    o.Med = s.Med
    o.RouteTable = s.RouteTable
    o.AddressPrefix = s.AddressPrefix
    o.NextHop = s.NextHop
    o.FromPeer = s.FromPeer
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
        Med: o.Answer.Med,
        AddressPrefix: util.EntToStr(o.Answer.AddressPrefix),
        NextHop: util.MemToStr(o.Answer.NextHop),
        FromPeer: util.MemToStr(o.Answer.FromPeer),
    }

    if o.Answer.AsPathRegex != nil {
        ans.AsPathRegex = o.Answer.AsPathRegex.Regex
    }

    if o.Answer.CommunityRegex != nil {
        ans.CommunityRegex = o.Answer.CommunityRegex.Regex
    }

    if o.Answer.ExtendedCommunityRegex != nil {
        ans.ExtendedCommunityRegex = o.Answer.ExtendedCommunityRegex.Regex
    }

    return ans
}


type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Enable: util.AsBool(o.Answer.Enable),
        Med: o.Answer.Med,
        RouteTable: o.Answer.RouteTable,
        AddressPrefix: util.EntToStr(o.Answer.AddressPrefix),
        NextHop: util.MemToStr(o.Answer.NextHop),
        FromPeer: util.MemToStr(o.Answer.FromPeer),
    }

    if o.Answer.AsPathRegex != nil {
        ans.AsPathRegex = o.Answer.AsPathRegex.Regex
    }

    if o.Answer.CommunityRegex != nil {
        ans.CommunityRegex = o.Answer.CommunityRegex.Regex
    }

    if o.Answer.ExtendedCommunityRegex != nil {
        ans.ExtendedCommunityRegex = o.Answer.ExtendedCommunityRegex.Regex
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    AsPathRegex *regex `xml:"match>as-path"`
    CommunityRegex *regex `xml:"match>community"`
    ExtendedCommunityRegex *regex `xml:"match>extended-community"`
    Med string `xml:"match>med,omitempty"`
    AddressPrefix *util.EntryType `xml:"match>address-prefix"`
    NextHop *util.MemberType `xml:"match>nexthop"`
    FromPeer *util.MemberType `xml:"match>from-peer"`
}

type regex struct {
    Regex string `xml:"regex,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        Med: e.Med,
        AddressPrefix: util.StrToEnt(e.AddressPrefix),
        NextHop: util.StrToMem(e.NextHop),
        FromPeer: util.StrToMem(e.FromPeer),
    }

    if e.AsPathRegex != "" {
        ans.AsPathRegex = &regex{
            Regex: e.AsPathRegex,
        }
    }

    if e.CommunityRegex != "" {
        ans.CommunityRegex = &regex{
            Regex: e.CommunityRegex,
        }
    }

    if e.ExtendedCommunityRegex != "" {
        ans.ExtendedCommunityRegex = &regex{
            Regex: e.ExtendedCommunityRegex,
        }
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    AsPathRegex *regex `xml:"match>as-path"`
    CommunityRegex *regex `xml:"match>community"`
    ExtendedCommunityRegex *regex `xml:"match>extended-community"`
    Med string `xml:"match>med,omitempty"`
    RouteTable string `xml:"match>route-table,omitempty"`
    AddressPrefix *util.EntryType `xml:"match>address-prefix"`
    NextHop *util.MemberType `xml:"match>nexthop"`
    FromPeer *util.MemberType `xml:"match>from-peer"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        Med: e.Med,
        RouteTable: e.RouteTable,
        AddressPrefix: util.StrToEnt(e.AddressPrefix),
        NextHop: util.StrToMem(e.NextHop),
        FromPeer: util.StrToMem(e.FromPeer),
    }

    if e.AsPathRegex != "" {
        ans.AsPathRegex = &regex{
            Regex: e.AsPathRegex,
        }
    }

    if e.CommunityRegex != "" {
        ans.CommunityRegex = &regex{
            Regex: e.CommunityRegex,
        }
    }

    if e.ExtendedCommunityRegex != "" {
        ans.ExtendedCommunityRegex = &regex{
            Regex: e.ExtendedCommunityRegex,
        }
    }

    return ans
}
