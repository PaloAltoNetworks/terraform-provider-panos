package vlan

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
    StaticIps []string // ordered
    EnableDhcp bool
    CreateDhcpDefaultRoute bool
    DhcpDefaultRouteMetric int
    ManagementProfile string
    Mtu int
    AdjustTcpMss bool
    Ipv4MssAdjust int
    Ipv6MssAdjust int

    raw map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Comment = s.Comment
    o.NetflowProfile = s.NetflowProfile
    o.StaticIps = s.StaticIps
    o.EnableDhcp = s.EnableDhcp
    o.CreateDhcpDefaultRoute = s.CreateDhcpDefaultRoute
    o.DhcpDefaultRouteMetric = s.DhcpDefaultRouteMetric
    o.ManagementProfile = s.ManagementProfile
    o.Mtu = s.Mtu
    o.AdjustTcpMss = s.AdjustTcpMss
    o.Ipv4MssAdjust = s.Ipv4MssAdjust
    o.Ipv6MssAdjust = s.Ipv6MssAdjust
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
        AdjustTcpMss: util.AsBool(o.Answer.AdjustTcpMss),
    }
    if o.Answer.Dhcp != nil {
        ans.EnableDhcp = util.AsBool(o.Answer.Dhcp.Enable)
        ans.CreateDhcpDefaultRoute = util.AsBool(o.Answer.Dhcp.CreateDefaultRoute)
        ans.DhcpDefaultRouteMetric = int(o.Answer.Dhcp.Metric)
    }

    ans.raw = make(map[string] string)
    if o.Answer.Ipv6 != nil {
        ans.raw["ipv6"] = util.CleanRawXml(o.Answer.Ipv6.Text)
    }
    if o.Answer.Arp != nil {
        ans.raw["arp"] = util.CleanRawXml(o.Answer.Arp.Text)
    }
    if o.Answer.NdpProxy != nil {
        ans.raw["ndp"] = util.CleanRawXml(o.Answer.NdpProxy.Text)
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
    Dhcp *dhcpSettings `xml:"dhcp-client"`
    Mtu int `xml:"mtu,omitempty"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`
    AdjustTcpMss string `xml:"adjust-tcp-mss"`

    Ipv6 *util.RawXml `xml:"ipv6"`
    Arp *util.RawXml `xml:"arp"`
    NdpProxy *util.RawXml `xml:"ndp-proxy"`
}

type dhcpSettings struct {
    Enable string `xml:"enable"`
    CreateDefaultRoute string `xml:"create-default-route"`
    Metric int `xml:"default-route-metric,omitempty"`
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Comment: o.Answer.Comment,
        NetflowProfile: o.Answer.NetflowProfile,
        StaticIps: util.EntToStr(o.Answer.StaticIps),
        Mtu: int(o.Answer.Mtu),
        ManagementProfile: o.Answer.ManagementProfile,
        AdjustTcpMss: util.AsBool(o.Answer.AdjustTcpMss),
        Ipv4MssAdjust: int(o.Answer.Ipv4MssAdjust),
        Ipv6MssAdjust: int(o.Answer.Ipv6MssAdjust),
    }
    if o.Answer.Dhcp != nil {
        ans.EnableDhcp = util.AsBool(o.Answer.Dhcp.Enable)
        ans.CreateDhcpDefaultRoute = util.AsBool(o.Answer.Dhcp.CreateDefaultRoute)
        ans.DhcpDefaultRouteMetric = int(o.Answer.Dhcp.Metric)
    }

    ans.raw = make(map[string] string)
    if o.Answer.Ipv6 != nil {
        ans.raw["ipv6"] = util.CleanRawXml(o.Answer.Ipv6.Text)
    }
    if o.Answer.Arp != nil {
        ans.raw["arp"] = util.CleanRawXml(o.Answer.Arp.Text)
    }
    if o.Answer.NdpProxy != nil {
        ans.raw["ndp"] = util.CleanRawXml(o.Answer.NdpProxy.Text)
    }
    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Comment string `xml:"comment,omitempty"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    StaticIps *util.EntryType `xml:"ip"`
    Dhcp *dhcpSettings `xml:"dhcp-client"`
    Mtu int `xml:"mtu,omitempty"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`
    AdjustTcpMss string `xml:"adjust-tcp-mss>enable"`
    Ipv4MssAdjust int `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
    Ipv6MssAdjust int `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`

    Ipv6 *util.RawXml `xml:"ipv6"`
    Arp *util.RawXml `xml:"arp"`
    NdpProxy *util.RawXml `xml:"ndp-proxy"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Comment: e.Comment,
        NetflowProfile: e.NetflowProfile,
        StaticIps: util.StrToEnt(e.StaticIps),
        Mtu: e.Mtu,
        ManagementProfile: e.ManagementProfile,
        AdjustTcpMss: util.YesNo(e.AdjustTcpMss),
    }

    if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
        v := dhcpSettings{
            Enable: util.YesNo(e.EnableDhcp),
            CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
            Metric: e.DhcpDefaultRouteMetric,
        }
        ans.Dhcp = &v
    }

    if text, ok := e.raw["ipv6"]; ok {
        ans.Ipv6 = &util.RawXml{text}
    }
    if text, ok := e.raw["arp"]; ok {
        ans.Arp = &util.RawXml{text}
    }
    if text, ok := e.raw["ndp"]; ok {
        ans.NdpProxy = &util.RawXml{text}
    }

    return ans
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Comment: e.Comment,
        NetflowProfile: e.NetflowProfile,
        StaticIps: util.StrToEnt(e.StaticIps),
        Mtu: e.Mtu,
        ManagementProfile: e.ManagementProfile,
        AdjustTcpMss: util.YesNo(e.AdjustTcpMss),
        Ipv4MssAdjust: e.Ipv4MssAdjust,
        Ipv6MssAdjust: e.Ipv6MssAdjust,
    }

    if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
        v := dhcpSettings{
            Enable: util.YesNo(e.EnableDhcp),
            CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
            Metric: e.DhcpDefaultRouteMetric,
        }
        ans.Dhcp = &v
    }

    if text, ok := e.raw["ipv6"]; ok {
        ans.Ipv6 = &util.RawXml{text}
    }
    if text, ok := e.raw["arp"]; ok {
        ans.Arp = &util.RawXml{text}
    }
    if text, ok := e.raw["ndp"]; ok {
        ans.NdpProxy = &util.RawXml{text}
    }

    return ans
}
