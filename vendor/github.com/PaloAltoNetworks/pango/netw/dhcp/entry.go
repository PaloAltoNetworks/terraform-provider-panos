package dhcp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a DHCP
// relay and server.
type Entry struct {
	Name  string
	Relay *Relay
	//Server *Server

	raw map[string]string
}

type Relay struct {
	Ipv4Enabled bool
	Ipv4Servers []string // unordered
	Ipv6Enabled bool
	Ipv6Servers []Ipv6Server
}

type Ipv6Server struct {
	Server    string
	Interface string
}

//type Server struct{}

func (o *Entry) Copy(s Entry) {
	if s.Relay != nil {
		o.Relay = &Relay{
			Ipv4Enabled: s.Relay.Ipv4Enabled,
			Ipv6Enabled: s.Relay.Ipv6Enabled,
		}
		if s.Relay.Ipv4Servers == nil {
			o.Relay.Ipv4Servers = nil
		} else {
			o.Relay.Ipv4Servers = make([]string, len(s.Relay.Ipv4Servers))
			copy(o.Relay.Ipv4Servers, s.Relay.Ipv4Servers)
		}
		if s.Relay.Ipv6Servers == nil {
			o.Relay.Ipv6Servers = nil
		} else {
			o.Relay.Ipv6Servers = make([]Ipv6Server, len(s.Relay.Ipv6Servers))
			copy(o.Relay.Ipv6Servers, s.Relay.Ipv6Servers)
		}
	} else {
		o.Relay = nil
	}
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{Name: o.Name}

	if o.Relay != nil {
		ans.Relay = &Relay{}
		if o.Relay.RelayIpv4 != nil {
			ans.Relay.Ipv4Enabled = util.AsBool(o.Relay.RelayIpv4.Ipv4Enabled)
			ans.Relay.Ipv4Servers = util.MemToStr(o.Relay.RelayIpv4.Ipv4Servers)
		}
		if o.Relay.RelayIpv6 != nil {
			ans.Relay.Ipv6Enabled = util.AsBool(o.Relay.RelayIpv6.Ipv6Enabled)
			if o.Relay.RelayIpv6.Ipv6Servers != nil {
				entries := o.Relay.RelayIpv6.Ipv6Servers.Entries
				x := make([]Ipv6Server, len(entries))
				for i := range entries {
					x[i].Server = entries[i].Name
					x[i].Interface = entries[i].Interface
				}
				ans.Relay.Ipv6Servers = x
			}
		}
	}

	raw := make(map[string]string)
	if o.Server != nil {
		raw["server"] = util.CleanRawXml(o.Server.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name     `xml:"entry"`
	Name    string       `xml:"name,attr"`
	Relay   *relay       `xml:"relay"`
	Server  *util.RawXml `xml:"server"`
}

type relay struct {
	RelayIpv4 *relayIpv4 `xml:"ip"`
	RelayIpv6 *relayIpv6 `xml:"ipv6"`
}

type relayIpv4 struct {
	Ipv4Enabled string           `xml:"enabled"`
	Ipv4Servers *util.MemberType `xml:"server"`
}
type relayIpv6 struct {
	Ipv6Enabled string             `xml:"enabled"`
	Ipv6Servers *ipv6ServerEntries `xml:"server"`
}

type ipv6ServerEntries struct {
	Entries []ipv6ServerEntry `xml:"entry"`
}

type ipv6ServerEntry struct {
	XMLName   xml.Name `xml:"entry"`
	Name      string   `xml:"name,attr"`
	Interface string   `xml:"interface,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{Name: e.Name}

	if e.Relay != nil {
		if len(e.Relay.Ipv4Servers) > 0 || len(e.Relay.Ipv6Servers) > 0 {
			ans.Relay = &relay{}
			if len(e.Relay.Ipv4Servers) > 0 {
				ans.Relay.RelayIpv4 = &relayIpv4{
					Ipv4Enabled: util.YesNo(e.Relay.Ipv4Enabled),
					Ipv4Servers: util.StrToMem(e.Relay.Ipv4Servers),
				}
			}
			if len(e.Relay.Ipv6Servers) > 0 {
				ans.Relay.RelayIpv6 = &relayIpv6{
					Ipv6Enabled: util.YesNo(e.Relay.Ipv6Enabled),
					Ipv6Servers: &ipv6ServerEntries{},
				}
				x := make([]ipv6ServerEntry, len(e.Relay.Ipv6Servers))
				for i := range e.Relay.Ipv6Servers {
					x[i] = ipv6ServerEntry{
						Name:      e.Relay.Ipv6Servers[i].Server,
						Interface: e.Relay.Ipv6Servers[i].Interface,
					}
				}
				ans.Relay.RelayIpv6.Ipv6Servers.Entries = x
			}
		}
	}

	if text, present := e.raw["server"]; present {
		ans.Server = &util.RawXml{text}
	}

	return ans
}
