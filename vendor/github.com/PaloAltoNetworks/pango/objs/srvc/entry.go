package srvc

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a service
// object.
//
// Protocol should be either "tcp" or "udp".
type Entry struct {
    Name string
    Description string
    Protocol string
    SourcePort string
    DestinationPort string
    Tags []string // ordered
    OverrideSessionTimeout bool // 8.1+
    OverrideTimeout int // 8.1+
    OverrideHalfClosedTimeout int // 8.1+
    OverrideTimeWaitTimeout int // 8.1+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Protocol = s.Protocol
    o.SourcePort = s.SourcePort
    o.DestinationPort = s.DestinationPort
    o.Tags = s.Tags
    o.OverrideSessionTimeout = s.OverrideSessionTimeout
    o.OverrideTimeout = s.OverrideTimeout
    o.OverrideHalfClosedTimeout = s.OverrideHalfClosedTimeout
    o.OverrideTimeWaitTimeout = s.OverrideTimeWaitTimeout
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
    case o.Answer.TcpProto != nil:
        ans.Protocol = ProtocolTcp
        ans.SourcePort = o.Answer.TcpProto.SourcePort
        ans.DestinationPort = o.Answer.TcpProto.DestinationPort
    case o.Answer.UdpProto != nil:
        ans.Protocol = ProtocolUdp
        ans.SourcePort = o.Answer.UdpProto.SourcePort
        ans.DestinationPort = o.Answer.UdpProto.DestinationPort
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
    case o.Answer.TcpProto != nil:
        ans.Protocol = ProtocolTcp
        ans.SourcePort = o.Answer.TcpProto.SourcePort
        ans.DestinationPort = o.Answer.TcpProto.DestinationPort
        if o.Answer.TcpProto.Override != nil && o.Answer.TcpProto.Override.Yes != nil {
            ans.OverrideSessionTimeout = true
            ans.OverrideTimeout = o.Answer.TcpProto.Override.Yes.OverrideTimeout
            ans.OverrideHalfClosedTimeout = o.Answer.TcpProto.Override.Yes.OverrideHalfClosedTimeout
            ans.OverrideTimeWaitTimeout = o.Answer.TcpProto.Override.Yes.OverrideTimeWaitTimeout
        }
    case o.Answer.UdpProto != nil:
        ans.Protocol = ProtocolUdp
        ans.SourcePort = o.Answer.UdpProto.SourcePort
        ans.DestinationPort = o.Answer.UdpProto.DestinationPort
        if o.Answer.UdpProto.Override != nil && o.Answer.UdpProto.Override.Yes != nil {
            ans.OverrideSessionTimeout = true
            ans.OverrideTimeout = o.Answer.UdpProto.Override.Yes.OverrideTimeout
        }
    case o.Answer.SctpProto != nil:
        ans.Protocol = ProtocolSctp
        ans.SourcePort = o.Answer.SctpProto.SourcePort
        ans.DestinationPort = o.Answer.SctpProto.DestinationPort
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TcpProto *protoDef `xml:"protocol>tcp"`
    UdpProto *protoDef `xml:"protocol>udp"`
    Description string `xml:"description"`
    Tags *util.MemberType `xml:"tag"`
}

type protoDef struct {
    SourcePort string `xml:"source-port,omitempty"`
    DestinationPort string `xml:"port"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Description: e.Description,
        Tags: util.StrToMem(e.Tags),
    }
    switch e.Protocol {
    case ProtocolTcp:
        ans.TcpProto = &protoDef{
            e.SourcePort,
            e.DestinationPort,
        }
    case ProtocolUdp:
        ans.UdpProto = &protoDef{
            e.SourcePort,
            e.DestinationPort,
        }
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TcpProto *tcpProto `xml:"protocol>tcp"`
    UdpProto *udpProto `xml:"protocol>udp"`
    SctpProto *protoDef `xml:"protocol>sctp"`
    Description string `xml:"description"`
    Tags *util.MemberType `xml:"tag"`
}

type tcpProto struct {
    SourcePort string `xml:"source-port,omitempty"`
    DestinationPort string `xml:"port"`
    Override *tcpOverride `xml:"override"`
}

type tcpOverride struct {
    No *string `xml:"no"`
    Yes *yesTcpOverride `xml:"yes"`
}

type yesTcpOverride struct {
    OverrideTimeout int `xml:"timeout,omitempty"`
    OverrideHalfClosedTimeout int `xml:"halfclose-timeout,omitempty"`
    OverrideTimeWaitTimeout int `xml:"timewait-timeout,omitempty"`
}

type udpProto struct {
    SourcePort string `xml:"source-port,omitempty"`
    DestinationPort string `xml:"port"`
    Override *udpOverride `xml:"override"`
}

type udpOverride struct {
    No *string `xml:"no"`
    Yes *yesUdpOverride `xml:"yes"`
}

type yesUdpOverride struct {
    OverrideTimeout int `xml:"timeout,omitempty"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Description: e.Description,
        Tags: util.StrToMem(e.Tags),
    }

    switch e.Protocol {
    case ProtocolTcp:
        ans.TcpProto = &tcpProto{
            SourcePort: e.SourcePort,
            DestinationPort: e.DestinationPort,
        }
        if e.OverrideSessionTimeout {
            ans.TcpProto.Override = &tcpOverride{
                Yes: &yesTcpOverride{
                    OverrideTimeout: e.OverrideTimeout,
                    OverrideHalfClosedTimeout: e.OverrideHalfClosedTimeout,
                    OverrideTimeWaitTimeout: e.OverrideTimeWaitTimeout,
                },
            }
        }
    case ProtocolUdp:
        ans.UdpProto = &udpProto{
            SourcePort: e.SourcePort,
            DestinationPort: e.DestinationPort,
        }
        if e.OverrideSessionTimeout {
            ans.UdpProto.Override = &udpOverride{
                Yes: &yesUdpOverride{
                    OverrideTimeout: e.OverrideTimeout,
                },
            }
        }
    case ProtocolSctp:
        ans.SctpProto = &protoDef{
            SourcePort: e.SourcePort,
            DestinationPort: e.DestinationPort,
        }
    }

    return ans
}
