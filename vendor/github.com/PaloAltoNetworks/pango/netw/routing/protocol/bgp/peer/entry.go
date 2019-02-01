package peer

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a BGP
// peer group peer.
type Entry struct {
    Name string
    Enable bool
    PeerAs string
    LocalAddressInterface string
    LocalAddressIp string
    PeerAddressIp string
    ReflectorClient string
    PeeringType string
    MaxPrefixes string
    AuthProfile string
    KeepAliveInterval int
    MultiHop int
    OpenDelayTime int
    HoldTime int
    IdleHoldTime int
    AllowIncomingConnections bool
    IncomingConnectionsRemotePort int
    AllowOutgoingConnections bool
    OutgoingConnectionsLocalPort int
    BfdProfile string // 7.1+
    EnableMpBgp bool // 8.0+
    AddressFamilyType string // 8.0+
    SubsequentAddressFamilyUnicast bool // 8.0+
    SubsequentAddressFamilyMulticast bool // 8.0+
    EnableSenderSideLoopDetection bool // 8.0+
    MinRouteAdvertisementInterval int // 8.1+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Enable = s.Enable
    o.PeerAs = s.PeerAs
    o.LocalAddressInterface = s.LocalAddressInterface
    o.LocalAddressIp = s.LocalAddressIp
    o.PeerAddressIp = s.PeerAddressIp
    o.ReflectorClient = s.ReflectorClient
    o.PeeringType = s.PeeringType
    o.MaxPrefixes = s.MaxPrefixes
    o.AuthProfile = s.AuthProfile
    o.KeepAliveInterval = s.KeepAliveInterval
    o.MultiHop = s.MultiHop
    o.OpenDelayTime = s.OpenDelayTime
    o.HoldTime = s.HoldTime
    o.IdleHoldTime = s.IdleHoldTime
    o.AllowIncomingConnections = s.AllowIncomingConnections
    o.IncomingConnectionsRemotePort = s.IncomingConnectionsRemotePort
    o.AllowOutgoingConnections = s.AllowOutgoingConnections
    o.OutgoingConnectionsLocalPort = s.OutgoingConnectionsLocalPort
    o.BfdProfile = s.BfdProfile
    o.EnableMpBgp = s.EnableMpBgp
    o.AddressFamilyType = s.AddressFamilyType
    o.SubsequentAddressFamilyUnicast = s.SubsequentAddressFamilyUnicast
    o.SubsequentAddressFamilyMulticast = s.SubsequentAddressFamilyMulticast
    o.EnableSenderSideLoopDetection = s.EnableSenderSideLoopDetection
    o.MinRouteAdvertisementInterval = s.MinRouteAdvertisementInterval
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
        PeerAs: o.Answer.PeerAs,
        LocalAddressInterface: o.Answer.LocalAddressInterface,
        LocalAddressIp: o.Answer.LocalAddressIp,
        PeerAddressIp: o.Answer.PeerAddressIp,
        ReflectorClient: o.Answer.ReflectorClient,
        PeeringType: o.Answer.PeeringType,
        MaxPrefixes: o.Answer.MaxPrefixes,
    }

    if o.Answer.Options != nil {
        ans.AuthProfile = o.Answer.Options.AuthProfile
        ans.KeepAliveInterval = o.Answer.Options.KeepAliveInterval
        ans.MultiHop = o.Answer.Options.MultiHop
        ans.OpenDelayTime = o.Answer.Options.OpenDelayTime
        ans.HoldTime = o.Answer.Options.HoldTime
        ans.IdleHoldTime = o.Answer.Options.IdleHoldTime

        if o.Answer.Options.BgpIn != nil {
            ans.AllowIncomingConnections = util.AsBool(o.Answer.Options.BgpIn.AllowIncomingConnections)
            ans.IncomingConnectionsRemotePort = o.Answer.Options.BgpIn.IncomingConnectionsRemotePort
        }

        if o.Answer.Options.BgpOut != nil {
            ans.AllowOutgoingConnections = util.AsBool(o.Answer.Options.BgpOut.AllowOutgoingConnections)
            ans.OutgoingConnectionsLocalPort = o.Answer.Options.BgpOut.OutgoingConnectionsLocalPort
        }
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
        PeerAs: o.Answer.PeerAs,
        LocalAddressInterface: o.Answer.LocalAddressInterface,
        LocalAddressIp: o.Answer.LocalAddressIp,
        PeerAddressIp: o.Answer.PeerAddressIp,
        ReflectorClient: o.Answer.ReflectorClient,
        PeeringType: o.Answer.PeeringType,
        MaxPrefixes: o.Answer.MaxPrefixes,
    }

    if o.Answer.Bfd != nil {
        ans.BfdProfile = o.Answer.Bfd.BfdProfile
    }

    if o.Answer.Options != nil {
        ans.AuthProfile = o.Answer.Options.AuthProfile
        ans.KeepAliveInterval = o.Answer.Options.KeepAliveInterval
        ans.MultiHop = o.Answer.Options.MultiHop
        ans.OpenDelayTime = o.Answer.Options.OpenDelayTime
        ans.HoldTime = o.Answer.Options.HoldTime
        ans.IdleHoldTime = o.Answer.Options.IdleHoldTime

        if o.Answer.Options.BgpIn != nil {
            ans.AllowIncomingConnections = util.AsBool(o.Answer.Options.BgpIn.AllowIncomingConnections)
            ans.IncomingConnectionsRemotePort = o.Answer.Options.BgpIn.IncomingConnectionsRemotePort
        }

        if o.Answer.Options.BgpOut != nil {
            ans.AllowOutgoingConnections = util.AsBool(o.Answer.Options.BgpOut.AllowOutgoingConnections)
            ans.OutgoingConnectionsLocalPort = o.Answer.Options.BgpOut.OutgoingConnectionsLocalPort
        }
    }

    return ans
}

type container_v3 struct {
    Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Enable: util.AsBool(o.Answer.Enable),
        PeerAs: o.Answer.PeerAs,
        EnableMpBgp: util.AsBool(o.Answer.EnableMpBgp),
        AddressFamilyType: o.Answer.AddressFamilyType,
        EnableSenderSideLoopDetection: util.AsBool(o.Answer.EnableSenderSideLoopDetection),
        LocalAddressInterface: o.Answer.LocalAddressInterface,
        LocalAddressIp: o.Answer.LocalAddressIp,
        PeerAddressIp: o.Answer.PeerAddressIp,
        ReflectorClient: o.Answer.ReflectorClient,
        PeeringType: o.Answer.PeeringType,
        MaxPrefixes: o.Answer.MaxPrefixes,
    }

    if o.Answer.Safi != nil {
        ans.SubsequentAddressFamilyUnicast = util.AsBool(o.Answer.Safi.SubsequentAddressFamilyUnicast)
        ans.SubsequentAddressFamilyMulticast = util.AsBool(o.Answer.Safi.SubsequentAddressFamilyMulticast)
    }

    if o.Answer.Bfd != nil {
        ans.BfdProfile = o.Answer.Bfd.BfdProfile
    }

    if o.Answer.Options != nil {
        ans.AuthProfile = o.Answer.Options.AuthProfile
        ans.KeepAliveInterval = o.Answer.Options.KeepAliveInterval
        ans.MultiHop = o.Answer.Options.MultiHop
        ans.OpenDelayTime = o.Answer.Options.OpenDelayTime
        ans.HoldTime = o.Answer.Options.HoldTime
        ans.IdleHoldTime = o.Answer.Options.IdleHoldTime

        if o.Answer.Options.BgpIn != nil {
            ans.AllowIncomingConnections = util.AsBool(o.Answer.Options.BgpIn.AllowIncomingConnections)
            ans.IncomingConnectionsRemotePort = o.Answer.Options.BgpIn.IncomingConnectionsRemotePort
        }

        if o.Answer.Options.BgpOut != nil {
            ans.AllowOutgoingConnections = util.AsBool(o.Answer.Options.BgpOut.AllowOutgoingConnections)
            ans.OutgoingConnectionsLocalPort = o.Answer.Options.BgpOut.OutgoingConnectionsLocalPort
        }
    }

    return ans
}

type container_v4 struct {
    Answer entry_v4 `xml:"result>entry"`
}

func (o *container_v4) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Enable: util.AsBool(o.Answer.Enable),
        PeerAs: o.Answer.PeerAs,
        EnableMpBgp: util.AsBool(o.Answer.EnableMpBgp),
        AddressFamilyType: o.Answer.AddressFamilyType,
        EnableSenderSideLoopDetection: util.AsBool(o.Answer.EnableSenderSideLoopDetection),
        LocalAddressInterface: o.Answer.LocalAddressInterface,
        LocalAddressIp: o.Answer.LocalAddressIp,
        PeerAddressIp: o.Answer.PeerAddressIp,
        ReflectorClient: o.Answer.ReflectorClient,
        PeeringType: o.Answer.PeeringType,
        MaxPrefixes: o.Answer.MaxPrefixes,
    }

    if o.Answer.Safi != nil {
        ans.SubsequentAddressFamilyUnicast = util.AsBool(o.Answer.Safi.SubsequentAddressFamilyUnicast)
        ans.SubsequentAddressFamilyMulticast = util.AsBool(o.Answer.Safi.SubsequentAddressFamilyMulticast)
    }

    if o.Answer.Bfd != nil {
        ans.BfdProfile = o.Answer.Bfd.BfdProfile
    }

    if o.Answer.Options != nil {
        ans.AuthProfile = o.Answer.Options.AuthProfile
        ans.KeepAliveInterval = o.Answer.Options.KeepAliveInterval
        ans.MultiHop = o.Answer.Options.MultiHop
        ans.OpenDelayTime = o.Answer.Options.OpenDelayTime
        ans.HoldTime = o.Answer.Options.HoldTime
        ans.IdleHoldTime = o.Answer.Options.IdleHoldTime
        ans.MinRouteAdvertisementInterval = o.Answer.Options.MinRouteAdvertisementInterval

        if o.Answer.Options.BgpIn != nil {
            ans.AllowIncomingConnections = util.AsBool(o.Answer.Options.BgpIn.AllowIncomingConnections)
            ans.IncomingConnectionsRemotePort = o.Answer.Options.BgpIn.IncomingConnectionsRemotePort
        }

        if o.Answer.Options.BgpOut != nil {
            ans.AllowOutgoingConnections = util.AsBool(o.Answer.Options.BgpOut.AllowOutgoingConnections)
            ans.OutgoingConnectionsLocalPort = o.Answer.Options.BgpOut.OutgoingConnectionsLocalPort
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    PeerAs string `xml:"peer-as,omitempty"`
    LocalAddressInterface string `xml:"local-address>interface"`
    LocalAddressIp string `xml:"local-address>ip,omitempty"`
    PeerAddressIp string `xml:"peer-address>ip"`
    ReflectorClient string `xml:"reflector-client,omitempty"`
    PeeringType string `xml:"peering-type,omitempty"`
    MaxPrefixes string `xml:"max-prefixes,omitempty"`
    Options *opts_v1 `xml:"connection-options"`
}

type opts_v1 struct {
    AuthProfile string `xml:"authentication,omitempty"`
    KeepAliveInterval int `xml:"keep-alive-interval,omitempty"`
    MultiHop int `xml:"multihop,omitempty"`
    OpenDelayTime int `xml:"open-delay-time,omitempty"`
    HoldTime int `xml:"hold-time,omitempty"`
    IdleHoldTime int `xml:"idle-hold-time,omitempty"`
    BgpIn *bgpIn `xml:"incoming-bgp-connection"`
    BgpOut *bgpOut `xml:"outgoing-bgp-connection"`
}

type bgpIn struct {
    AllowIncomingConnections string `xml:"allow"`
    IncomingConnectionsRemotePort int `xml:"remote-port,omitempty"`
}

type bgpOut struct {
    AllowOutgoingConnections string `xml:"allow"`
    OutgoingConnectionsLocalPort int `xml:"local-port,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        PeerAs: e.PeerAs,
        LocalAddressInterface: e.LocalAddressInterface,
        LocalAddressIp: e.LocalAddressIp,
        PeerAddressIp: e.PeerAddressIp,
        ReflectorClient: e.ReflectorClient,
        PeeringType: e.PeeringType,
        MaxPrefixes: e.MaxPrefixes,
    }

    hasIn := e.AllowIncomingConnections || e.IncomingConnectionsRemotePort != 0
    hasOut := e.AllowOutgoingConnections || e.OutgoingConnectionsLocalPort != 0

    if hasIn || hasOut || e.AuthProfile != "" || e.KeepAliveInterval != 0 || e.MultiHop != 0 || e.OpenDelayTime != 0 || e.HoldTime != 0 || e.IdleHoldTime != 0 {
        ans.Options = &opts_v1{
            AuthProfile: e.AuthProfile,
            KeepAliveInterval: e.KeepAliveInterval,
            MultiHop: e.MultiHop,
            OpenDelayTime: e.OpenDelayTime,
            HoldTime: e.HoldTime,
            IdleHoldTime: e.IdleHoldTime,
        }

        if hasIn {
            ans.Options.BgpIn = &bgpIn{
                AllowIncomingConnections: util.YesNo(e.AllowIncomingConnections),
                IncomingConnectionsRemotePort: e.IncomingConnectionsRemotePort,
            }
        }

        if hasOut {
            ans.Options.BgpOut = &bgpOut{
                AllowOutgoingConnections: util.YesNo(e.AllowOutgoingConnections),
                OutgoingConnectionsLocalPort: e.OutgoingConnectionsLocalPort,
            }
        }
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    PeerAs string `xml:"peer-as,omitempty"`
    LocalAddressInterface string `xml:"local-address>interface"`
    LocalAddressIp string `xml:"local-address>ip,omitempty"`
    PeerAddressIp string `xml:"peer-address>ip"`
    ReflectorClient string `xml:"reflector-client,omitempty"`
    PeeringType string `xml:"peering-type,omitempty"`
    MaxPrefixes string `xml:"max-prefixes,omitempty"`
    Bfd *bfd `xml:"bfd"`
    Options *opts_v1 `xml:"connection-options"`
}

type bfd struct {
    BfdProfile string `xml:"profile,omitempty"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Enable: util.YesNo(e.Enable),
        PeerAs: e.PeerAs,
        LocalAddressInterface: e.LocalAddressInterface,
        LocalAddressIp: e.LocalAddressIp,
        PeerAddressIp: e.PeerAddressIp,
        ReflectorClient: e.ReflectorClient,
        PeeringType: e.PeeringType,
        MaxPrefixes: e.MaxPrefixes,
    }

    if e.BfdProfile != "" {
        ans.Bfd = &bfd{
            BfdProfile: e.BfdProfile,
        }
    }

    hasIn := e.AllowIncomingConnections || e.IncomingConnectionsRemotePort != 0
    hasOut := e.AllowOutgoingConnections || e.OutgoingConnectionsLocalPort != 0

    if hasIn || hasOut || e.AuthProfile != "" || e.KeepAliveInterval != 0 || e.MultiHop != 0 || e.OpenDelayTime != 0 || e.HoldTime != 0 || e.IdleHoldTime != 0 {
        ans.Options = &opts_v1{
            AuthProfile: e.AuthProfile,
            KeepAliveInterval: e.KeepAliveInterval,
            MultiHop: e.MultiHop,
            OpenDelayTime: e.OpenDelayTime,
            HoldTime: e.HoldTime,
            IdleHoldTime: e.IdleHoldTime,
        }

        if hasIn {
            ans.Options.BgpIn = &bgpIn{
                AllowIncomingConnections: util.YesNo(e.AllowIncomingConnections),
                IncomingConnectionsRemotePort: e.IncomingConnectionsRemotePort,
            }
        }

        if hasOut {
            ans.Options.BgpOut = &bgpOut{
                AllowOutgoingConnections: util.YesNo(e.AllowOutgoingConnections),
                OutgoingConnectionsLocalPort: e.OutgoingConnectionsLocalPort,
            }
        }
    }

    return ans
}

type entry_v3 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    PeerAs string `xml:"peer-as,omitempty"`
    EnableMpBgp string `xml:"enable-mp-bgp"`
    AddressFamilyType string `xml:"address-family-identifier,omitempty"`
    Safi *safi `xml:"subsequent-address-family-identifier"`
    EnableSenderSideLoopDetection string `xml:"enable-sender-side-loop-detection"`
    LocalAddressInterface string `xml:"local-address>interface"`
    LocalAddressIp string `xml:"local-address>ip,omitempty"`
    PeerAddressIp string `xml:"peer-address>ip"`
    ReflectorClient string `xml:"reflector-client,omitempty"`
    PeeringType string `xml:"peering-type,omitempty"`
    MaxPrefixes string `xml:"max-prefixes,omitempty"`
    Bfd *bfd `xml:"bfd"`
    Options *opts_v1 `xml:"connection-options"`
}

type safi struct {
    SubsequentAddressFamilyUnicast string `xml:"unicast"`
    SubsequentAddressFamilyMulticast string `xml:"multicast"`
}

func specify_v3(e Entry) interface{} {
    ans := entry_v3{
        Name: e.Name,
        EnableMpBgp: util.YesNo(e.EnableMpBgp),
        AddressFamilyType: e.AddressFamilyType,
        EnableSenderSideLoopDetection: util.YesNo(e.EnableSenderSideLoopDetection),
        Enable: util.YesNo(e.Enable),
        PeerAs: e.PeerAs,
        LocalAddressInterface: e.LocalAddressInterface,
        LocalAddressIp: e.LocalAddressIp,
        PeerAddressIp: e.PeerAddressIp,
        ReflectorClient: e.ReflectorClient,
        PeeringType: e.PeeringType,
        MaxPrefixes: e.MaxPrefixes,
    }

    if e.SubsequentAddressFamilyUnicast || e.SubsequentAddressFamilyMulticast {
        ans.Safi = &safi{
            SubsequentAddressFamilyUnicast: util.YesNo(e.SubsequentAddressFamilyUnicast),
            SubsequentAddressFamilyMulticast: util.YesNo(e.SubsequentAddressFamilyMulticast),
        }
    }

    if e.BfdProfile != "" {
        ans.Bfd = &bfd{
            BfdProfile: e.BfdProfile,
        }
    }

    hasIn := e.AllowIncomingConnections || e.IncomingConnectionsRemotePort != 0
    hasOut := e.AllowOutgoingConnections || e.OutgoingConnectionsLocalPort != 0

    if hasIn || hasOut || e.AuthProfile != "" || e.KeepAliveInterval != 0 || e.MultiHop != 0 || e.OpenDelayTime != 0 || e.HoldTime != 0 || e.IdleHoldTime != 0 {
        ans.Options = &opts_v1{
            AuthProfile: e.AuthProfile,
            KeepAliveInterval: e.KeepAliveInterval,
            MultiHop: e.MultiHop,
            OpenDelayTime: e.OpenDelayTime,
            HoldTime: e.HoldTime,
            IdleHoldTime: e.IdleHoldTime,
        }

        if hasIn {
            ans.Options.BgpIn = &bgpIn{
                AllowIncomingConnections: util.YesNo(e.AllowIncomingConnections),
                IncomingConnectionsRemotePort: e.IncomingConnectionsRemotePort,
            }
        }

        if hasOut {
            ans.Options.BgpOut = &bgpOut{
                AllowOutgoingConnections: util.YesNo(e.AllowOutgoingConnections),
                OutgoingConnectionsLocalPort: e.OutgoingConnectionsLocalPort,
            }
        }
    }

    return ans
}

type entry_v4 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Enable string `xml:"enable"`
    PeerAs string `xml:"peer-as,omitempty"`
    EnableMpBgp string `xml:"enable-mp-bgp"`
    AddressFamilyType string `xml:"address-family-identifier,omitempty"`
    Safi *safi `xml:"subsequent-address-family-identifier"`
    EnableSenderSideLoopDetection string `xml:"enable-sender-side-loop-detection"`
    LocalAddressInterface string `xml:"local-address>interface"`
    LocalAddressIp string `xml:"local-address>ip,omitempty"`
    PeerAddressIp string `xml:"peer-address>ip"`
    ReflectorClient string `xml:"reflector-client,omitempty"`
    PeeringType string `xml:"peering-type,omitempty"`
    MaxPrefixes string `xml:"max-prefixes,omitempty"`
    Bfd *bfd `xml:"bfd"`
    Options *opts_v2 `xml:"connection-options"`
}

type opts_v2 struct {
    AuthProfile string `xml:"authentication,omitempty"`
    KeepAliveInterval int `xml:"keep-alive-interval,omitempty"`
    MultiHop int `xml:"multihop,omitempty"`
    OpenDelayTime int `xml:"open-delay-time,omitempty"`
    HoldTime int `xml:"hold-time,omitempty"`
    IdleHoldTime int `xml:"idle-hold-time,omitempty"`
    MinRouteAdvertisementInterval int `xml:"min-route-adv-interval,omitempty"`
    BgpIn *bgpIn `xml:"incoming-bgp-connection"`
    BgpOut *bgpOut `xml:"outgoing-bgp-connection"`
}

func specify_v4(e Entry) interface{} {
    ans := entry_v4{
        Name: e.Name,
        EnableMpBgp: util.YesNo(e.EnableMpBgp),
        AddressFamilyType: e.AddressFamilyType,
        EnableSenderSideLoopDetection: util.YesNo(e.EnableSenderSideLoopDetection),
        Enable: util.YesNo(e.Enable),
        PeerAs: e.PeerAs,
        LocalAddressInterface: e.LocalAddressInterface,
        LocalAddressIp: e.LocalAddressIp,
        PeerAddressIp: e.PeerAddressIp,
        ReflectorClient: e.ReflectorClient,
        PeeringType: e.PeeringType,
        MaxPrefixes: e.MaxPrefixes,
    }

    if e.SubsequentAddressFamilyUnicast || e.SubsequentAddressFamilyMulticast {
        ans.Safi = &safi{
            SubsequentAddressFamilyUnicast: util.YesNo(e.SubsequentAddressFamilyUnicast),
            SubsequentAddressFamilyMulticast: util.YesNo(e.SubsequentAddressFamilyMulticast),
        }
    }

    if e.BfdProfile != "" {
        ans.Bfd = &bfd{
            BfdProfile: e.BfdProfile,
        }
    }

    hasIn := e.AllowIncomingConnections || e.IncomingConnectionsRemotePort != 0
    hasOut := e.AllowOutgoingConnections || e.OutgoingConnectionsLocalPort != 0

    if hasIn || hasOut || e.AuthProfile != "" || e.KeepAliveInterval != 0 || e.MultiHop != 0 || e.OpenDelayTime != 0 || e.HoldTime != 0 || e.IdleHoldTime != 0 || e.MinRouteAdvertisementInterval != 0 {
        ans.Options = &opts_v2{
            AuthProfile: e.AuthProfile,
            KeepAliveInterval: e.KeepAliveInterval,
            MultiHop: e.MultiHop,
            OpenDelayTime: e.OpenDelayTime,
            HoldTime: e.HoldTime,
            IdleHoldTime: e.IdleHoldTime,
            MinRouteAdvertisementInterval: e.MinRouteAdvertisementInterval,
        }

        if hasIn {
            ans.Options.BgpIn = &bgpIn{
                AllowIncomingConnections: util.YesNo(e.AllowIncomingConnections),
                IncomingConnectionsRemotePort: e.IncomingConnectionsRemotePort,
            }
        }

        if hasOut {
            ans.Options.BgpOut = &bgpOut{
                AllowOutgoingConnections: util.YesNo(e.AllowOutgoingConnections),
                OutgoingConnectionsLocalPort: e.OutgoingConnectionsLocalPort,
            }
        }
    }

    return ans
}
