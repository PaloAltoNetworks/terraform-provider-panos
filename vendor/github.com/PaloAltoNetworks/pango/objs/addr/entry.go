package addr

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of an address
// object.
type Entry struct {
    Name string
    Value string
    Type string
    Description string
    Tags []string // ordered
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Value = s.Value
    o.Type = s.Type
    o.Description = s.Description
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
        Tags: util.MemToStr(o.Answer.Tags),
    }
    switch {
    case o.Answer.IpNetmask != nil:
        ans.Type = IpNetmask
        ans.Value = o.Answer.IpNetmask.Value
    case o.Answer.IpRange != nil:
        ans.Type = IpRange
        ans.Value = o.Answer.IpRange.Value
    case o.Answer.Fqdn != nil:
        ans.Type = Fqdn
        ans.Value = o.Answer.Fqdn.Value
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Description: o.Answer.Description,
        Tags: util.MemToStr(o.Answer.Tags),
    }
    switch {
    case o.Answer.IpNetmask != nil:
        ans.Type = IpNetmask
        ans.Value = o.Answer.IpNetmask.Value
    case o.Answer.IpRange != nil:
        ans.Type = IpRange
        ans.Value = o.Answer.IpRange.Value
    case o.Answer.Fqdn != nil:
        ans.Type = Fqdn
        ans.Value = o.Answer.Fqdn.Value
    case o.Answer.IpWildcard != nil:
        ans.Type = IpWildcard
        ans.Value = o.Answer.IpWildcard.Value
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    IpNetmask *valType `xml:"ip-netmask"`
    IpRange *valType `xml:"ip-range"`
    Fqdn *valType `xml:"fqdn"`
    Description string `xml:"description"`
    Tags *util.MemberType `xml:"tag"`
}

type valType struct {
    Value string `xml:",chardata"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Tags: util.StrToMem(e.Tags),
    }
    vt := &valType{e.Value}
    switch e.Type {
    case IpNetmask:
        ans.IpNetmask = vt
    case IpRange:
        ans.IpRange = vt
    case Fqdn:
        ans.Fqdn = vt
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    IpNetmask *valType `xml:"ip-netmask"`
    IpRange *valType `xml:"ip-range"`
    Fqdn *valType `xml:"fqdn"`
    IpWildcard *valType `xml:"ip-wildcard"`
    Description string `xml:"description"`
    Tags *util.MemberType `xml:"tag"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Description: e.Description,
        Tags: util.StrToMem(e.Tags),
    }
    vt := &valType{e.Value}
    switch e.Type {
    case IpNetmask:
        ans.IpNetmask = vt
    case IpRange:
        ans.IpRange = vt
    case Fqdn:
        ans.Fqdn = vt
    case IpWildcard:
        ans.IpWildcard = vt
    }

    return ans
}
