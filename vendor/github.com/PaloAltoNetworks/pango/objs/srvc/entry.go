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
    Tags []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.Protocol = s.Protocol
    o.SourcePort = s.SourcePort
    o.DestinationPort = s.DestinationPort
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
    case o.Answer.TcpProto != nil:
        ans.Protocol = "tcp"
        ans.SourcePort = o.Answer.TcpProto.SourcePort
        ans.DestinationPort = o.Answer.TcpProto.DestinationPort
    case o.Answer.UdpProto != nil:
        ans.Protocol = "udp"
        ans.SourcePort = o.Answer.UdpProto.SourcePort
        ans.DestinationPort = o.Answer.UdpProto.DestinationPort
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TcpProto *protoDef `xml:"protocol>tcp"`
    UdpProto *protoDef `xml:"protocol>udp"`
    Description string `xml:"description"`
    Tags *util.Member `xml:"tag"`
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
    case "tcp":
        ans.TcpProto = &protoDef{
            e.SourcePort,
            e.DestinationPort,
        }
    case "udp":
        ans.UdpProto = &protoDef{
            e.SourcePort,
            e.DestinationPort,
        }
    }

    return ans
}
