package eth

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of an ethernet
// interface.
type Entry struct {
    Name string
    Mode string
    StaticIps []string // ordered
    EnableDhcp bool
    CreateDhcpDefaultRoute bool
    DhcpDefaultRouteMetric int
    Ipv6Enabled bool
    Ipv6InterfaceId string
    ManagementProfile string
    Mtu int
    AdjustTcpMss bool
    NetflowProfile string
    LldpEnabled bool
    LldpProfile string
    LinkSpeed string
    LinkDuplex string
    LinkState string
    AggregateGroup string
    Comment string
    Ipv4MssAdjust int // 7.1+
    Ipv6MssAdjust int // 7.1+
    EnableUntaggedSubinterface bool // 7.1+
    DecryptForward bool // 8.1+
    RxPolicingRate int // 8.1+
    TxPolicingRate int // 8.1+
    DhcpSendHostnameEnable bool // 9.0+
    DhcpSendHostnameValue string // 9.0+

    raw map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Mode = s.Mode
    o.StaticIps = s.StaticIps
    o.EnableDhcp = s.EnableDhcp
    o.CreateDhcpDefaultRoute = s.CreateDhcpDefaultRoute
    o.DhcpDefaultRouteMetric = s.DhcpDefaultRouteMetric
    o.Ipv6Enabled = s.Ipv6Enabled
    o.ManagementProfile = s.ManagementProfile
    o.Mtu = s.Mtu
    o.AdjustTcpMss = s.AdjustTcpMss
    o.NetflowProfile = s.NetflowProfile
    o.LldpEnabled = s.LldpEnabled
    o.LldpProfile = s.LldpProfile
    o.LinkSpeed = s.LinkSpeed
    o.LinkDuplex = s.LinkDuplex
    o.LinkState = s.LinkState
    o.AggregateGroup = s.AggregateGroup
    o.Comment = s.Comment
    o.Ipv4MssAdjust = s.Ipv4MssAdjust
    o.Ipv6MssAdjust = s.Ipv6MssAdjust
    o.EnableUntaggedSubinterface = s.EnableUntaggedSubinterface
    o.DecryptForward = s.DecryptForward
    o.RxPolicingRate = s.RxPolicingRate
    o.TxPolicingRate = s.TxPolicingRate
    o.DhcpSendHostnameEnable = s.DhcpSendHostnameEnable
    o.DhcpSendHostnameValue = s.DhcpSendHostnameValue
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
        LinkSpeed: o.Answer.LinkSpeed,
        LinkDuplex: o.Answer.LinkDuplex,
        LinkState: o.Answer.LinkState,
        Comment: o.Answer.Comment,
    }
    ans.raw = make(map[string] string)
    switch {
        case o.Answer.ModeL3 != nil:
            ans.Mode = "layer3"
            ans.ManagementProfile = o.Answer.ModeL3.ManagementProfile
            ans.Mtu = o.Answer.ModeL3.Mtu
            ans.NetflowProfile = o.Answer.ModeL3.NetflowProfile
            ans.AdjustTcpMss = util.AsBool(o.Answer.ModeL3.AdjustTcpMss)
            ans.StaticIps = util.EntToStr(o.Answer.ModeL3.StaticIps)
            if o.Answer.ModeL3.Dhcp != nil {
                ans.EnableDhcp = util.AsBool(o.Answer.ModeL3.Dhcp.Enable)
                ans.CreateDhcpDefaultRoute = util.AsBool(o.Answer.ModeL3.Dhcp.CreateDefaultRoute)
                ans.DhcpDefaultRouteMetric = o.Answer.ModeL3.Dhcp.Metric
            }

            if o.Answer.ModeL3.Ipv6 != nil {
                ans.Ipv6Enabled = util.AsBool(o.Answer.ModeL3.Ipv6.Enabled)
                ans.Ipv6InterfaceId = o.Answer.ModeL3.Ipv6.Ipv6InterfaceId
                if o.Answer.ModeL3.Ipv6.Address != nil {
                    ans.raw["v6adr"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Address.Text)
                }
                if o.Answer.ModeL3.Ipv6.Neighbor != nil {
                    ans.raw["v6nd"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Neighbor.Text)
                }
            }

            if o.Answer.ModeL3.Arp != nil {
                ans.raw["arp"] = util.CleanRawXml(o.Answer.ModeL3.Arp.Text)
            }
            if o.Answer.ModeL3.Subinterface != nil {
                ans.raw["l3subinterface"] = util.CleanRawXml(o.Answer.ModeL3.Subinterface.Text)
            }
        case o.Answer.ModeL2 != nil:
            ans.Mode = "layer2"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeL2.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeL2.LldpProfile
            ans.NetflowProfile = o.Answer.ModeL2.NetflowProfile
            if o.Answer.ModeL2.Subinterface != nil {
                ans.raw["l2subinterface"] = util.CleanRawXml(o.Answer.ModeL2.Subinterface.Text)
            }
        case o.Answer.ModeVwire != nil:
            ans.Mode = "virtual-wire"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeVwire.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeVwire.LldpProfile
            ans.NetflowProfile = o.Answer.ModeVwire.NetflowProfile
        case o.Answer.TapMode != nil:
            ans.Mode = "tap"
        case o.Answer.HaMode != nil:
            ans.Mode = "ha"
        case o.Answer.DecryptMirrorMode != nil:
            ans.Mode = "decrypt-mirror"
        case o.Answer.AggregateGroupMode != nil:
            ans.Mode = "aggregate-group"
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }
    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    ModeL2 *otherMode `xml:"layer2"`
    ModeL3 *l3Mode_v1 `xml:"layer3"`
    ModeVwire *otherMode `xml:"virtual-wire"`
    TapMode *emptyMode `xml:"tap"`
    HaMode *emptyMode `xml:"ha"`
    DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
    AggregateGroupMode *emptyMode `xml:"aggregate-group"`
    LinkSpeed string `xml:"link-speed,omitempty"`
    LinkDuplex string `xml:"link-duplex,omitempty"`
    LinkState string `xml:"link-state,omitempty"`
    Comment string `xml:"comment"`
}

type emptyMode struct {}

type otherMode struct {
    LldpEnabled string `xml:"lldp>enable"`
    LldpProfile string `xml:"lldp>profile"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    Subinterface *util.RawXml `xml:"units"`
}

type l3Mode_v1 struct {
    Ipv6 *ipv6 `xml:"ipv6"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`
    Mtu int `xml:"mtu,omitempty"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    AdjustTcpMss string `xml:"adjust-tcp-mss"`
    StaticIps *util.EntryType `xml:"ip"`
    Dhcp *dhcpSettings_v1 `xml:"dhcp-client"`
    Arp *util.RawXml `xml:"arp"`
    Subinterface *util.RawXml `xml:"units"`
}

type ipv6 struct {
    Enabled string `xml:"enabled"`
    Ipv6InterfaceId string `xml:"interface-id,omitempty"`
    Address *util.RawXml `xml:"address"`
    Neighbor *util.RawXml `xml:"neighbor-discovery"`
}

type dhcpSettings_v1 struct {
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
        LinkSpeed: o.Answer.LinkSpeed,
        LinkDuplex: o.Answer.LinkDuplex,
        LinkState: o.Answer.LinkState,
        Comment: o.Answer.Comment,
    }
    ans.raw = make(map[string] string)
    switch {
        case o.Answer.ModeL3 != nil:
            ans.Mode = "layer3"
            ans.ManagementProfile = o.Answer.ModeL3.ManagementProfile
            ans.Mtu = o.Answer.ModeL3.Mtu
            ans.NetflowProfile = o.Answer.ModeL3.NetflowProfile
            ans.AdjustTcpMss = util.AsBool(o.Answer.ModeL3.AdjustTcpMss)
            ans.Ipv4MssAdjust = o.Answer.ModeL3.Ipv4MssAdjust
            ans.Ipv6MssAdjust = o.Answer.ModeL3.Ipv6MssAdjust
            ans.StaticIps = util.EntToStr(o.Answer.ModeL3.StaticIps)
            ans.EnableUntaggedSubinterface = util.AsBool(o.Answer.ModeL3.EnableUntaggedSubinterface)

            if o.Answer.ModeL3.Dhcp != nil {
                ans.EnableDhcp = util.AsBool(o.Answer.ModeL3.Dhcp.Enable)
                ans.CreateDhcpDefaultRoute = util.AsBool(o.Answer.ModeL3.Dhcp.CreateDefaultRoute)
                ans.DhcpDefaultRouteMetric = o.Answer.ModeL3.Dhcp.Metric
            }

            if o.Answer.ModeL3.Ipv6 != nil {
                ans.Ipv6Enabled = util.AsBool(o.Answer.ModeL3.Ipv6.Enabled)
                ans.Ipv6InterfaceId = o.Answer.ModeL3.Ipv6.Ipv6InterfaceId
                if o.Answer.ModeL3.Ipv6.Address != nil {
                    ans.raw["v6adr"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Address.Text)
                }
                if o.Answer.ModeL3.Ipv6.Neighbor != nil {
                    ans.raw["v6nd"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Neighbor.Text)
                }
            }

            if o.Answer.ModeL3.Arp != nil {
                ans.raw["arp"] = util.CleanRawXml(o.Answer.ModeL3.Arp.Text)
            }
            if o.Answer.ModeL3.Subinterface != nil {
                ans.raw["l3subinterface"] = util.CleanRawXml(o.Answer.ModeL3.Subinterface.Text)
            }
            if o.Answer.ModeL3.Pppoe != nil {
                ans.raw["pppoe"] = util.CleanRawXml(o.Answer.ModeL3.Pppoe.Text)
            }
            if o.Answer.ModeL3.Ndp != nil {
                ans.raw["ndp"] = util.CleanRawXml(o.Answer.ModeL3.Ndp.Text)
            }
        case o.Answer.ModeL2 != nil:
            ans.Mode = "layer2"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeL2.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeL2.LldpProfile
            ans.NetflowProfile = o.Answer.ModeL2.NetflowProfile
            if o.Answer.ModeL2.Subinterface != nil {
                ans.raw["l2subinterface"] = util.CleanRawXml(o.Answer.ModeL2.Subinterface.Text)
            }
        case o.Answer.ModeVwire != nil:
            ans.Mode = "virtual-wire"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeVwire.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeVwire.LldpProfile
            ans.NetflowProfile = o.Answer.ModeVwire.NetflowProfile
        case o.Answer.TapMode != nil:
            ans.Mode = "tap"
        case o.Answer.HaMode != nil:
            ans.Mode = "ha"
        case o.Answer.DecryptMirrorMode != nil:
            ans.Mode = "decrypt-mirror"
        case o.Answer.AggregateGroupMode != nil:
            ans.Mode = "aggregate-group"
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }
    return ans
}

type container_v3 struct {
    Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        LinkSpeed: o.Answer.LinkSpeed,
        LinkDuplex: o.Answer.LinkDuplex,
        LinkState: o.Answer.LinkState,
        Comment: o.Answer.Comment,
    }
    ans.raw = make(map[string] string)
    switch {
        case o.Answer.ModeL3 != nil:
            ans.Mode = "layer3"
            ans.ManagementProfile = o.Answer.ModeL3.ManagementProfile
            ans.Mtu = o.Answer.ModeL3.Mtu
            ans.NetflowProfile = o.Answer.ModeL3.NetflowProfile
            ans.AdjustTcpMss = util.AsBool(o.Answer.ModeL3.AdjustTcpMss)
            ans.Ipv4MssAdjust = o.Answer.ModeL3.Ipv4MssAdjust
            ans.Ipv6MssAdjust = o.Answer.ModeL3.Ipv6MssAdjust
            ans.StaticIps = util.EntToStr(o.Answer.ModeL3.StaticIps)
            ans.EnableUntaggedSubinterface = util.AsBool(o.Answer.ModeL3.EnableUntaggedSubinterface)
            ans.DecryptForward = util.AsBool(o.Answer.ModeL3.DecryptForward)

            if o.Answer.ModeL3.Dhcp != nil {
                ans.EnableDhcp = util.AsBool(o.Answer.ModeL3.Dhcp.Enable)
                ans.CreateDhcpDefaultRoute = util.AsBool(o.Answer.ModeL3.Dhcp.CreateDefaultRoute)
                ans.DhcpDefaultRouteMetric = o.Answer.ModeL3.Dhcp.Metric
            }

            if o.Answer.ModeL3.Policing != nil {
                ans.RxPolicingRate = o.Answer.ModeL3.Policing.RxPolicingRate
                ans.TxPolicingRate = o.Answer.ModeL3.Policing.TxPolicingRate
            }

            if o.Answer.ModeL3.Ipv6 != nil {
                ans.Ipv6Enabled = util.AsBool(o.Answer.ModeL3.Ipv6.Enabled)
                ans.Ipv6InterfaceId = o.Answer.ModeL3.Ipv6.Ipv6InterfaceId
                if o.Answer.ModeL3.Ipv6.Address != nil {
                    ans.raw["v6adr"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Address.Text)
                }
                if o.Answer.ModeL3.Ipv6.Neighbor != nil {
                    ans.raw["v6nd"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Neighbor.Text)
                }
            }

            if o.Answer.ModeL3.Arp != nil {
                ans.raw["arp"] = util.CleanRawXml(o.Answer.ModeL3.Arp.Text)
            }
            if o.Answer.ModeL3.Subinterface != nil {
                ans.raw["l3subinterface"] = util.CleanRawXml(o.Answer.ModeL3.Subinterface.Text)
            }
            if o.Answer.ModeL3.Pppoe != nil {
                ans.raw["pppoe"] = util.CleanRawXml(o.Answer.ModeL3.Pppoe.Text)
            }
            if o.Answer.ModeL3.Ndp != nil {
                ans.raw["ndp"] = util.CleanRawXml(o.Answer.ModeL3.Ndp.Text)
            }
        case o.Answer.ModeL2 != nil:
            ans.Mode = "layer2"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeL2.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeL2.LldpProfile
            ans.NetflowProfile = o.Answer.ModeL2.NetflowProfile
            if o.Answer.ModeL2.Subinterface != nil {
                ans.raw["l2subinterface"] = util.CleanRawXml(o.Answer.ModeL2.Subinterface.Text)
            }
        case o.Answer.ModeVwire != nil:
            ans.Mode = "virtual-wire"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeVwire.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeVwire.LldpProfile
            ans.NetflowProfile = o.Answer.ModeVwire.NetflowProfile
        case o.Answer.TapMode != nil:
            ans.Mode = "tap"
        case o.Answer.HaMode != nil:
            ans.Mode = "ha"
        case o.Answer.DecryptMirrorMode != nil:
            ans.Mode = "decrypt-mirror"
        case o.Answer.AggregateGroupMode != nil:
            ans.Mode = "aggregate-group"
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }
    return ans
}

type container_v4 struct {
    Answer entry_v4 `xml:"result>entry"`
}

func (o *container_v4) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        LinkSpeed: o.Answer.LinkSpeed,
        LinkDuplex: o.Answer.LinkDuplex,
        LinkState: o.Answer.LinkState,
        Comment: o.Answer.Comment,
    }
    ans.raw = make(map[string] string)
    switch {
        case o.Answer.ModeL3 != nil:
            ans.Mode = "layer3"
            ans.ManagementProfile = o.Answer.ModeL3.ManagementProfile
            ans.Mtu = o.Answer.ModeL3.Mtu
            ans.NetflowProfile = o.Answer.ModeL3.NetflowProfile
            ans.AdjustTcpMss = util.AsBool(o.Answer.ModeL3.AdjustTcpMss)
            ans.Ipv4MssAdjust = o.Answer.ModeL3.Ipv4MssAdjust
            ans.Ipv6MssAdjust = o.Answer.ModeL3.Ipv6MssAdjust
            ans.StaticIps = util.EntToStr(o.Answer.ModeL3.StaticIps)
            ans.EnableUntaggedSubinterface = util.AsBool(o.Answer.ModeL3.EnableUntaggedSubinterface)
            ans.DecryptForward = util.AsBool(o.Answer.ModeL3.DecryptForward)

            if o.Answer.ModeL3.Dhcp != nil {
                ans.EnableDhcp = util.AsBool(o.Answer.ModeL3.Dhcp.Enable)
                ans.CreateDhcpDefaultRoute = util.AsBool(o.Answer.ModeL3.Dhcp.CreateDefaultRoute)
                ans.DhcpDefaultRouteMetric = o.Answer.ModeL3.Dhcp.Metric
                if o.Answer.ModeL3.Dhcp.Hostname != nil {
                    ans.DhcpSendHostnameEnable = util.AsBool(o.Answer.ModeL3.Dhcp.Hostname.DhcpSendHostnameEnable)
                    ans.DhcpSendHostnameValue = o.Answer.ModeL3.Dhcp.Hostname.DhcpSendHostnameValue
                }
            }

            if o.Answer.ModeL3.Policing != nil {
                ans.RxPolicingRate = o.Answer.ModeL3.Policing.RxPolicingRate
                ans.TxPolicingRate = o.Answer.ModeL3.Policing.TxPolicingRate
            }

            if o.Answer.ModeL3.Ipv6 != nil {
                ans.Ipv6Enabled = util.AsBool(o.Answer.ModeL3.Ipv6.Enabled)
                ans.Ipv6InterfaceId = o.Answer.ModeL3.Ipv6.Ipv6InterfaceId
                if o.Answer.ModeL3.Ipv6.Address != nil {
                    ans.raw["v6adr"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Address.Text)
                }
                if o.Answer.ModeL3.Ipv6.Neighbor != nil {
                    ans.raw["v6nd"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6.Neighbor.Text)
                }
            }

            if o.Answer.ModeL3.Arp != nil {
                ans.raw["arp"] = util.CleanRawXml(o.Answer.ModeL3.Arp.Text)
            }
            if o.Answer.ModeL3.Subinterface != nil {
                ans.raw["l3subinterface"] = util.CleanRawXml(o.Answer.ModeL3.Subinterface.Text)
            }
            if o.Answer.ModeL3.Pppoe != nil {
                ans.raw["pppoe"] = util.CleanRawXml(o.Answer.ModeL3.Pppoe.Text)
            }
            if o.Answer.ModeL3.Ndp != nil {
                ans.raw["ndp"] = util.CleanRawXml(o.Answer.ModeL3.Ndp.Text)
            }
            if o.Answer.ModeL3.Ipv6Client != nil {
                ans.raw["v6client"] = util.CleanRawXml(o.Answer.ModeL3.Ipv6Client.Text)
            }
        case o.Answer.ModeL2 != nil:
            ans.Mode = "layer2"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeL2.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeL2.LldpProfile
            ans.NetflowProfile = o.Answer.ModeL2.NetflowProfile
            if o.Answer.ModeL2.Subinterface != nil {
                ans.raw["l2subinterface"] = util.CleanRawXml(o.Answer.ModeL2.Subinterface.Text)
            }
        case o.Answer.ModeVwire != nil:
            ans.Mode = "virtual-wire"
            ans.LldpEnabled = util.AsBool(o.Answer.ModeVwire.LldpEnabled)
            ans.LldpProfile = o.Answer.ModeVwire.LldpProfile
            ans.NetflowProfile = o.Answer.ModeVwire.NetflowProfile
        case o.Answer.TapMode != nil:
            ans.Mode = "tap"
        case o.Answer.HaMode != nil:
            ans.Mode = "ha"
        case o.Answer.DecryptMirrorMode != nil:
            ans.Mode = "decrypt-mirror"
        case o.Answer.AggregateGroupMode != nil:
            ans.Mode = "aggregate-group"
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }
    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    ModeL3 *l3Mode_v2 `xml:"layer3"`
    ModeL2 *otherMode `xml:"layer2"`
    ModeVwire *otherMode `xml:"virtual-wire"`
    TapMode *emptyMode `xml:"tap"`
    HaMode *emptyMode `xml:"ha"`
    DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
    AggregateGroupMode *emptyMode `xml:"aggregate-group"`
    LinkSpeed string `xml:"link-speed,omitempty"`
    LinkDuplex string `xml:"link-duplex,omitempty"`
    LinkState string `xml:"link-state,omitempty"`
    Comment string `xml:"comment"`
}

type l3Mode_v2 struct {
    Ipv6 *ipv6 `xml:"ipv6"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`
    Mtu int `xml:"mtu,omitempty"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    AdjustTcpMss string `xml:"adjust-tcp-mss>enable"`
    Ipv4MssAdjust int `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
    Ipv6MssAdjust int `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`
    StaticIps *util.EntryType `xml:"ip"`
    Dhcp *dhcpSettings_v1 `xml:"dhcp-client"`
    EnableUntaggedSubinterface string `xml:"untagged-sub-interface,omitempty"`
    Arp *util.RawXml `xml:"arp"`
    Pppoe *util.RawXml `xml:"pppoe"`
    Ndp *util.RawXml `xml:"ndp-proxy"`
    Subinterface *util.RawXml `xml:"units"`
}

type entry_v3 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    ModeL3 *l3Mode_v3 `xml:"layer3"`
    ModeL2 *otherMode `xml:"layer2"`
    ModeVwire *otherMode `xml:"virtual-wire"`
    TapMode *emptyMode `xml:"tap"`
    HaMode *emptyMode `xml:"ha"`
    DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
    AggregateGroupMode *emptyMode `xml:"aggregate-group"`
    LinkSpeed string `xml:"link-speed,omitempty"`
    LinkDuplex string `xml:"link-duplex,omitempty"`
    LinkState string `xml:"link-state,omitempty"`
    Comment string `xml:"comment"`
}

type l3Mode_v3 struct {
    Ipv6 *ipv6 `xml:"ipv6"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`
    Mtu int `xml:"mtu,omitempty"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    AdjustTcpMss string `xml:"adjust-tcp-mss>enable"`
    Ipv4MssAdjust int `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
    Ipv6MssAdjust int `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`
    StaticIps *util.EntryType `xml:"ip"`
    Dhcp *dhcpSettings_v1 `xml:"dhcp-client"`
    EnableUntaggedSubinterface string `xml:"untagged-sub-interface,omitempty"`
    DecryptForward string `xml:"decrypt-forward,omitempty"`
    Policing *policing `xml:"policing"`
    Arp *util.RawXml `xml:"arp"`
    Pppoe *util.RawXml `xml:"pppoe"`
    Ndp *util.RawXml `xml:"ndp-proxy"`
    Subinterface *util.RawXml `xml:"units"`
}

type policing struct {
    RxPolicingRate int `xml:"rx-rate,omitempty"`
    TxPolicingRate int `xml:"tx-rate,omitempty"`
}

type entry_v4 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    ModeL3 *l3Mode_v4 `xml:"layer3"`
    ModeL2 *otherMode `xml:"layer2"`
    ModeVwire *otherMode `xml:"virtual-wire"`
    TapMode *emptyMode `xml:"tap"`
    HaMode *emptyMode `xml:"ha"`
    DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
    AggregateGroupMode *emptyMode `xml:"aggregate-group"`
    LinkSpeed string `xml:"link-speed,omitempty"`
    LinkDuplex string `xml:"link-duplex,omitempty"`
    LinkState string `xml:"link-state,omitempty"`
    Comment string `xml:"comment"`
}

type l3Mode_v4 struct {
    Ipv6 *ipv6 `xml:"ipv6"`
    ManagementProfile string `xml:"interface-management-profile,omitempty"`
    Mtu int `xml:"mtu,omitempty"`
    NetflowProfile string `xml:"netflow-profile,omitempty"`
    AdjustTcpMss string `xml:"adjust-tcp-mss>enable"`
    Ipv4MssAdjust int `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
    Ipv6MssAdjust int `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`
    StaticIps *util.EntryType `xml:"ip"`
    Dhcp *dhcpSettings_v2 `xml:"dhcp-client"`
    EnableUntaggedSubinterface string `xml:"untagged-sub-interface,omitempty"`
    DecryptForward string `xml:"decrypt-forward,omitempty"`
    Policing *policing `xml:"policing"`
    Arp *util.RawXml `xml:"arp"`
    Pppoe *util.RawXml `xml:"pppoe"`
    Ndp *util.RawXml `xml:"ndp-proxy"`
    Ipv6Client *util.RawXml `xml:"ipv6-client"`
    Subinterface *util.RawXml `xml:"units"`
}

type dhcpSettings_v2 struct {
    Enable string `xml:"enable"`
    CreateDefaultRoute string `xml:"create-default-route"`
    Metric int `xml:"default-route-metric,omitempty"`
    Hostname *dhcpHostname `xml:"send-hostname"`
}

type dhcpHostname struct {
    DhcpSendHostnameEnable string `xml:"enable,omitempty"`
    DhcpSendHostnameValue string `xml:"hostname,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        LinkSpeed: e.LinkSpeed,
        LinkDuplex: e.LinkDuplex,
        LinkState: e.LinkState,
        Comment: e.Comment,
    }

    switch e.Mode {
    case "layer3":
        i := &l3Mode_v1{
            StaticIps: util.StrToEnt(e.StaticIps),
            ManagementProfile: e.ManagementProfile,
            Mtu: e.Mtu,
            NetflowProfile: e.NetflowProfile,
            AdjustTcpMss: util.YesNo(e.AdjustTcpMss),
        }

        if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
            i.Dhcp = &dhcpSettings_v1{
                Enable: util.YesNo(e.EnableDhcp),
                CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
                Metric: e.DhcpDefaultRouteMetric,
            }
        }

        v6adr := e.raw["v6adr"]
        v6nd := e.raw["v6nd"]
        if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
            v6 := ipv6{
                Enabled: util.YesNo(e.Ipv6Enabled),
                Ipv6InterfaceId: e.Ipv6InterfaceId,
            }
            if v6adr != "" {
                v6.Address = &util.RawXml{v6adr}
            }
            if v6nd != "" {
                v6.Neighbor = &util.RawXml{v6nd}
            }
            i.Ipv6 = &v6
        }

        if text, present := e.raw["arp"]; present {
            i.Arp = &util.RawXml{text}
        }
        if text, present := e.raw["l3subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        ans.ModeL3 = i
    case "layer2":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        if text, present := e.raw["l2subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        ans.ModeL2 = i
    case "virtual-wire":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        ans.ModeVwire = i
    case "tap":
        ans.TapMode = &emptyMode{}
    case "ha":
        ans.HaMode = &emptyMode{}
    case "decrypt-mirror":
        ans.DecryptMirrorMode = &emptyMode{}
    case "aggregate-group":
        ans.AggregateGroupMode = &emptyMode{}
    }

    return ans
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        LinkSpeed: e.LinkSpeed,
        LinkDuplex: e.LinkDuplex,
        LinkState: e.LinkState,
        Comment: e.Comment,
    }

    switch e.Mode {
    case "layer3":
        i := &l3Mode_v2{
            StaticIps: util.StrToEnt(e.StaticIps),
            ManagementProfile: e.ManagementProfile,
            Mtu: e.Mtu,
            NetflowProfile: e.NetflowProfile,
            AdjustTcpMss: util.YesNo(e.AdjustTcpMss),
            Ipv4MssAdjust: e.Ipv4MssAdjust,
            Ipv6MssAdjust: e.Ipv6MssAdjust,
        }

        if e.EnableUntaggedSubinterface {
            i.EnableUntaggedSubinterface = util.YesNo(e.EnableUntaggedSubinterface)
        }

        if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
            i.Dhcp = &dhcpSettings_v1{
                Enable: util.YesNo(e.EnableDhcp),
                CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
                Metric: e.DhcpDefaultRouteMetric,
            }
        }

        v6adr := e.raw["v6adr"]
        v6nd := e.raw["v6nd"]
        if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
            v6 := ipv6{
                Enabled: util.YesNo(e.Ipv6Enabled),
                Ipv6InterfaceId: e.Ipv6InterfaceId,
            }
            if v6adr != "" {
                v6.Address = &util.RawXml{v6adr}
            }
            if v6nd != "" {
                v6.Neighbor = &util.RawXml{v6nd}
            }
            i.Ipv6 = &v6
        }

        if text, present := e.raw["arp"]; present {
            i.Arp = &util.RawXml{text}
        }
        if text, present := e.raw["l3subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        if text := e.raw["pppoe"]; text != "" {
            i.Pppoe = &util.RawXml{text}
        }
        if text := e.raw["ndp"]; text != "" {
            i.Ndp = &util.RawXml{text}
        }
        ans.ModeL3 = i
    case "layer2":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        if text, present := e.raw["l2subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        ans.ModeL2 = i
    case "virtual-wire":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        ans.ModeVwire = i
    case "tap":
        ans.TapMode = &emptyMode{}
    case "ha":
        ans.HaMode = &emptyMode{}
    case "decrypt-mirror":
        ans.DecryptMirrorMode = &emptyMode{}
    case "aggregate-group":
        ans.AggregateGroupMode = &emptyMode{}
    }

    return ans
}

func specify_v3(e Entry) interface{} {
    ans := entry_v3{
        Name: e.Name,
        LinkSpeed: e.LinkSpeed,
        LinkDuplex: e.LinkDuplex,
        LinkState: e.LinkState,
        Comment: e.Comment,
    }

    switch e.Mode {
    case "layer3":
        i := &l3Mode_v3{
            StaticIps: util.StrToEnt(e.StaticIps),
            ManagementProfile: e.ManagementProfile,
            Mtu: e.Mtu,
            NetflowProfile: e.NetflowProfile,
            AdjustTcpMss: util.YesNo(e.AdjustTcpMss),
            Ipv4MssAdjust: e.Ipv4MssAdjust,
            Ipv6MssAdjust: e.Ipv6MssAdjust,
        }

        if e.EnableUntaggedSubinterface {
            i.EnableUntaggedSubinterface = util.YesNo(e.EnableUntaggedSubinterface)
        }

        if e.DecryptForward {
            i.DecryptForward = util.YesNo(e.DecryptForward)
        }

        if e.RxPolicingRate != 0 || e.TxPolicingRate != 0 {
            i.Policing = &policing{
                RxPolicingRate: e.RxPolicingRate,
                TxPolicingRate: e.TxPolicingRate,
            }
        }

        if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
            i.Dhcp = &dhcpSettings_v1{
                Enable: util.YesNo(e.EnableDhcp),
                CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
                Metric: e.DhcpDefaultRouteMetric,
            }
        }

        v6adr := e.raw["v6adr"]
        v6nd := e.raw["v6nd"]
        if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
            v6 := ipv6{
                Enabled: util.YesNo(e.Ipv6Enabled),
                Ipv6InterfaceId: e.Ipv6InterfaceId,
            }
            if v6adr != "" {
                v6.Address = &util.RawXml{v6adr}
            }
            if v6nd != "" {
                v6.Neighbor = &util.RawXml{v6nd}
            }
            i.Ipv6 = &v6
        }

        if text, present := e.raw["arp"]; present {
            i.Arp = &util.RawXml{text}
        }
        if text, present := e.raw["l3subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        if text := e.raw["pppoe"]; text != "" {
            i.Pppoe = &util.RawXml{text}
        }
        if text := e.raw["ndp"]; text != "" {
            i.Ndp = &util.RawXml{text}
        }
        ans.ModeL3 = i
    case "layer2":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        if text, present := e.raw["l2subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        ans.ModeL2 = i
    case "virtual-wire":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        ans.ModeVwire = i
    case "tap":
        ans.TapMode = &emptyMode{}
    case "ha":
        ans.HaMode = &emptyMode{}
    case "decrypt-mirror":
        ans.DecryptMirrorMode = &emptyMode{}
    case "aggregate-group":
        ans.AggregateGroupMode = &emptyMode{}
    }

    return ans
}

func specify_v4(e Entry) interface{} {
    ans := entry_v4{
        Name: e.Name,
        LinkSpeed: e.LinkSpeed,
        LinkDuplex: e.LinkDuplex,
        LinkState: e.LinkState,
        Comment: e.Comment,
    }

    switch e.Mode {
    case "layer3":
        i := &l3Mode_v4{
            StaticIps: util.StrToEnt(e.StaticIps),
            ManagementProfile: e.ManagementProfile,
            Mtu: e.Mtu,
            NetflowProfile: e.NetflowProfile,
            AdjustTcpMss: util.YesNo(e.AdjustTcpMss),
            Ipv4MssAdjust: e.Ipv4MssAdjust,
            Ipv6MssAdjust: e.Ipv6MssAdjust,
        }

        if e.EnableUntaggedSubinterface {
            i.EnableUntaggedSubinterface = util.YesNo(e.EnableUntaggedSubinterface)
        }

        if e.DecryptForward {
            i.DecryptForward = util.YesNo(e.DecryptForward)
        }

        if e.RxPolicingRate != 0 || e.TxPolicingRate != 0 {
            i.Policing = &policing{
                RxPolicingRate: e.RxPolicingRate,
                TxPolicingRate: e.TxPolicingRate,
            }
        }

        if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 || e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != "" {
            i.Dhcp = &dhcpSettings_v2{
                Enable: util.YesNo(e.EnableDhcp),
                CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
                Metric: e.DhcpDefaultRouteMetric,
            }

            if e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != "" {
                i.Dhcp.Hostname = &dhcpHostname{
                    DhcpSendHostnameEnable: util.YesNo(e.DhcpSendHostnameEnable),
                    DhcpSendHostnameValue: e.DhcpSendHostnameValue,
                }
            }
        }

        v6adr := e.raw["v6adr"]
        v6nd := e.raw["v6nd"]
        if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
            v6 := ipv6{
                Enabled: util.YesNo(e.Ipv6Enabled),
                Ipv6InterfaceId: e.Ipv6InterfaceId,
            }
            if v6adr != "" {
                v6.Address = &util.RawXml{v6adr}
            }
            if v6nd != "" {
                v6.Neighbor = &util.RawXml{v6nd}
            }
            i.Ipv6 = &v6
        }

        if text, present := e.raw["arp"]; present {
            i.Arp = &util.RawXml{text}
        }
        if text, present := e.raw["l3subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        if text := e.raw["pppoe"]; text != "" {
            i.Pppoe = &util.RawXml{text}
        }
        if text := e.raw["ndp"]; text != "" {
            i.Ndp = &util.RawXml{text}
        }
        if text := e.raw["v6client"]; text != "" {
            i.Ipv6Client = &util.RawXml{text}
        }
        ans.ModeL3 = i
    case "layer2":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        if text, present := e.raw["l2subinterface"]; present {
            i.Subinterface = &util.RawXml{text}
        }
        ans.ModeL2 = i
    case "virtual-wire":
        i := &otherMode{
            LldpEnabled: util.YesNo(e.LldpEnabled),
            LldpProfile: e.LldpProfile,
            NetflowProfile: e.NetflowProfile,
        }
        ans.ModeVwire = i
    case "tap":
        ans.TapMode = &emptyMode{}
    case "ha":
        ans.HaMode = &emptyMode{}
    case "decrypt-mirror":
        ans.DecryptMirrorMode = &emptyMode{}
    case "aggregate-group":
        ans.AggregateGroupMode = &emptyMode{}
    }

    return ans
}
