package ipv4

import (
    "encoding/xml"
)


// Entry is a normalized, version independent representation of an interface
// management profile.
type Entry struct {
    Name string
    Local string
    Remote string
    ProtocolAny bool
    ProtocolNumber int
    ProtocolTcpLocal int
    ProtocolTcpRemote int
    ProtocolUdpLocal int
    ProtocolUdpRemote int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Local = s.Local
    o.Remote = s.Remote
    o.ProtocolAny = s.ProtocolAny
    o.ProtocolNumber = s.ProtocolNumber
    o.ProtocolTcpLocal = s.ProtocolTcpLocal
    o.ProtocolTcpRemote = s.ProtocolTcpRemote
    o.ProtocolUdpLocal = s.ProtocolUdpLocal
    o.ProtocolUdpRemote = s.ProtocolUdpRemote
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
        Local: o.Answer.Local,
        Remote: o.Answer.Remote,
    }

    if o.Answer.Protocol != nil {
        if o.Answer.Protocol.Any != nil {
            ans.ProtocolAny = true
        } else if o.Answer.Protocol.Number != 0 {
            ans.ProtocolNumber = o.Answer.Protocol.Number
        } else if o.Answer.Protocol.Tcp != nil {
            ans.ProtocolTcpLocal = o.Answer.Protocol.Tcp.Local
            ans.ProtocolTcpRemote = o.Answer.Protocol.Tcp.Remote
        } else if o.Answer.Protocol.Udp != nil {
            ans.ProtocolUdpLocal = o.Answer.Protocol.Udp.Local
            ans.ProtocolUdpRemote = o.Answer.Protocol.Udp.Remote
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Local string `xml:"local,omitempty"`
    Remote string `xml:"remote,omitempty"`
    Protocol *proto `xml:"protocol"`
}

type proto struct {
    Any *string `xml:"any"`
    Number int `xml:"number,omitempty"`
    Tcp *subProto `xml:"tcp"`
    Udp *subProto `xml:"udp"`
}

type subProto struct {
    Local int `xml:"local-port,omitempty"`
    Remote int `xml:"remote-port,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Local: e.Local,
        Remote: e.Remote,
    }

    var p *proto
    if e.ProtocolAny {
        sp := ""
        p = &proto{Any: &sp}
    } else if e.ProtocolNumber != 0 {
        p = &proto{Number: e.ProtocolNumber}
    } else if e.ProtocolTcpLocal != 0 || e.ProtocolTcpRemote != 0 {
        p = &proto{Tcp: &subProto{
            Local: e.ProtocolTcpLocal,
            Remote: e.ProtocolTcpRemote,
        }}
    } else if e.ProtocolUdpLocal != 0 || e.ProtocolUdpRemote != 0 {
        p = &proto{Udp: &subProto{
            Local: e.ProtocolUdpLocal,
            Remote: e.ProtocolUdpRemote,
        }}
    }
    if p != nil {
        ans.Protocol = p
    }

    return ans
}
