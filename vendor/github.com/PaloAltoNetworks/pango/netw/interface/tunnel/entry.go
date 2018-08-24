package tunnel

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of
// a VLAN interface.
type Entry struct {
    Name string
    Comment string
    NetflowProfile string
    StaticIps []string
    ManagementProfile string
    Mtu int

    raw map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Comment = s.Comment
    o.NetflowProfile = s.NetflowProfile
    o.StaticIps = s.StaticIps
    o.ManagementProfile = s.ManagementProfile
    o.Mtu = s.Mtu
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
        Comment: o.Answer.Comment,
        NetflowProfile: o.Answer.NetflowProfile,
        StaticIps: util.EntToStr(o.Answer.StaticIps),
        Mtu: int(o.Answer.Mtu),
        ManagementProfile: o.Answer.ManagementProfile,
    }

    ans.raw = make(map[string] string)
    if o.Answer.Ipv6 != nil {
        ans.raw["ipv6"] = util.CleanRawXml(o.Answer.Ipv6.Text)
    }
    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Comment string `xml:"comment,omitempty"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    StaticIps *util.EntryType `xml:"ip"`
    Mtu int `xml:"mtu,omitempty"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`

    Ipv6 *util.RawXml `xml:"ipv6"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Comment: e.Comment,
        NetflowProfile: e.NetflowProfile,
        StaticIps: util.StrToEnt(e.StaticIps),
        Mtu: e.Mtu,
        ManagementProfile: e.ManagementProfile,
    }

    if text, ok := e.raw["ipv6"]; ok {
        ans.Ipv6 = &util.RawXml{text}
    }

    return ans
}
