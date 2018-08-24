package router

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a virtual
// router.
type Entry struct {
    Name string
    Interfaces []string
    StaticDist int
    StaticIpv6Dist int
    OspfIntDist int
    OspfExtDist int
    Ospfv3IntDist int
    Ospfv3ExtDist int
    IbgpDist int
    EbgpDist int
    RipDist int

    raw map[string] string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * StaticDist: 10
//      * StaticIpv6Dist: 10
//      * OspfIntDist: 30
//      * OspfExtDist: 110
//      * Ospfv3IntDist: 30
//      * Ospfv3ExtDist: 110
//      * IbgpDist: 200
//      * EbgpDist: 20
//      * RipDist: 120
func (o *Entry) Defaults() {
    if o.StaticDist == 0 {
        o.StaticDist = 10
    }

    if o.StaticIpv6Dist == 0 {
        o.StaticIpv6Dist = 10
    }

    if o.OspfIntDist == 0 {
        o.OspfIntDist = 30
    }

    if o.OspfExtDist == 0 {
        o.OspfExtDist = 110
    }

    if o.Ospfv3IntDist == 0 {
        o.Ospfv3IntDist = 30
    }

    if o.Ospfv3ExtDist == 0 {
        o.Ospfv3ExtDist = 110
    }

    if o.IbgpDist == 0 {
        o.IbgpDist = 200
    }

    if o.EbgpDist == 0 {
        o.EbgpDist = 20
    }

    if o.RipDist == 0 {
        o.RipDist = 120
    }
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Interfaces = s.Interfaces
    o.StaticDist = s.StaticDist
    o.StaticIpv6Dist = s.StaticIpv6Dist
    o.OspfIntDist = s.OspfIntDist
    o.OspfExtDist = s.OspfExtDist
    o.Ospfv3IntDist = s.Ospfv3IntDist
    o.Ospfv3ExtDist = s.Ospfv3ExtDist
    o.IbgpDist = s.IbgpDist
    o.EbgpDist = s.EbgpDist
    o.RipDist = s.RipDist
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
        StaticDist: o.Answer.Dist.StaticDist,
        StaticIpv6Dist: o.Answer.Dist.StaticIpv6Dist,
        OspfIntDist: o.Answer.Dist.OspfIntDist,
        OspfExtDist: o.Answer.Dist.OspfExtDist,
        Ospfv3IntDist: o.Answer.Dist.Ospfv3IntDist,
        Ospfv3ExtDist: o.Answer.Dist.Ospfv3ExtDist,
        IbgpDist: o.Answer.Dist.IbgpDist,
        EbgpDist: o.Answer.Dist.EbgpDist,
        RipDist: o.Answer.Dist.RipDist,
    }
    ans.raw = make(map[string] string)
    if o.Answer.Ecmp != nil {
        ans.raw["ecmp"] = util.CleanRawXml(o.Answer.Ecmp.Text)
    }
    if o.Answer.Multicast != nil {
        ans.raw["multicast"] = util.CleanRawXml(o.Answer.Multicast.Text)
    }
    if o.Answer.Protocol != nil {
        ans.raw["protocol"] = util.CleanRawXml(o.Answer.Protocol.Text)
    }
    if o.Answer.Routing != nil {
        ans.raw["routing"] = util.CleanRawXml(o.Answer.Routing.Text)
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }
    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Interfaces *util.MemberType `xml:"interface"`
    Dist dist `xml:"admin-dists"`
    Ecmp *util.RawXml `xml:"ecmp"`
    Multicast *util.RawXml `xml:"multicast"`
    Protocol *util.RawXml `xml:"protocol"`
    Routing *util.RawXml `xml:"routing-table"`
}

type dist struct {
    StaticDist int `xml:"static,omitempty"`
    StaticIpv6Dist int `xml:"static-ipv6,omitempty"`
    OspfIntDist int `xml:"ospf-int,omitempty"`
    OspfExtDist int `xml:"ospf-ext,omitempty"`
    Ospfv3IntDist int `xml:"ospfv3-int,omitempty"`
    Ospfv3ExtDist int `xml:"ospfv3-ext,omitempty"`
    IbgpDist int `xml:"ibgp,omitempty"`
    EbgpDist int `xml:"ebgp,omitempty"`
    RipDist int `xml:"rip,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Interfaces: util.StrToMem(e.Interfaces),
        Dist: dist{
            StaticDist: e.StaticDist,
            StaticIpv6Dist: e.StaticIpv6Dist,
            OspfIntDist: e.OspfIntDist,
            OspfExtDist: e.OspfExtDist,
            Ospfv3IntDist: e.Ospfv3IntDist,
            Ospfv3ExtDist: e.Ospfv3ExtDist,
            IbgpDist: e.IbgpDist,
            EbgpDist: e.EbgpDist,
            RipDist: e.RipDist,
        },
    }
    if text, present := e.raw["ecmp"]; present {
        ans.Ecmp = &util.RawXml{text}
    }
    if text, present := e.raw["multicast"]; present {
        ans.Multicast = &util.RawXml{text}
    }
    if text, present := e.raw["protocol"]; present {
        ans.Protocol = &util.RawXml{text}
    }
    if text, present := e.raw["routing"]; present {
        ans.Routing = &util.RawXml{text}
    }

    return ans
}
