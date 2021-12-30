package snmp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a snmptrap profile.
//
// PAN-OS 7.1+.
type Entry struct {
	Name       string
	V2cServers []V2cServer
	V3Servers  []V3Server
}

// V2cServer is a snmpv2 server.
type V2cServer struct {
	Name      string
	Manager   string
	Community string
}

// V3Server is a snmpv3 server.
type V3Server struct {
	Name         string
	Manager      string
	User         string
	EngineId     string
	AuthPassword string // encrypted
	PrivPassword string // encrypted
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	if s.V2cServers == nil {
		o.V2cServers = nil
	} else {
		o.V2cServers = make([]V2cServer, 0, len(s.V2cServers))
		for _, x := range s.V2cServers {
			o.V2cServers = append(o.V2cServers, V2cServer{
				Name:      x.Name,
				Manager:   x.Manager,
				Community: x.Community,
			})
		}
	}
	if s.V3Servers == nil {
		o.V3Servers = nil
	} else {
		o.V3Servers = make([]V3Server, 0, len(s.V3Servers))
		for _, x := range s.V3Servers {
			o.V3Servers = append(o.V3Servers, V3Server{
				Name:         x.Name,
				Manager:      x.Manager,
				User:         x.User,
				EngineId:     x.EngineId,
				AuthPassword: x.AuthPassword,
				PrivPassword: x.PrivPassword,
			})
		}
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
	ans := Entry{
		Name: o.Name,
	}

	if o.V2c != nil {
		list := make([]V2cServer, 0, len(o.V2c.Entries))
		for _, x := range o.V2c.Entries {
			list = append(list, V2cServer{
				Name:      x.Name,
				Manager:   x.Manager,
				Community: x.Community,
			})
		}
		ans.V2cServers = list
	}

	if o.V3 != nil {
		list := make([]V3Server, 0, len(o.V3.Entries))
		for _, x := range o.V3.Entries {
			list = append(list, V3Server{
				Name:         x.Name,
				Manager:      x.Manager,
				User:         x.User,
				EngineId:     x.EngineId,
				AuthPassword: x.AuthPassword,
				PrivPassword: x.PrivPassword,
			})
		}
		ans.V3Servers = list
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	V2c     *v2cs    `xml:"version>v2c"`
	V3      *v3s     `xml:"version>v3"`
}

type v2cs struct {
	Entries []v2cEntry `xml:"server>entry"`
}

type v2cEntry struct {
	XMLName   xml.Name `xml:"entry"`
	Name      string   `xml:"name,attr"`
	Manager   string   `xml:"manager"`
	Community string   `xml:"community"`
}

type v3s struct {
	Entries []v3Entry `xml:"server>entry"`
}

type v3Entry struct {
	XMLName      xml.Name `xml:"entry"`
	Name         string   `xml:"name,attr"`
	Manager      string   `xml:"manager"`
	User         string   `xml:"user"`
	EngineId     string   `xml:"engineid,omitempty"`
	AuthPassword string   `xml:"authpwd"`
	PrivPassword string   `xml:"privpwd"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	if len(e.V2cServers) > 0 {
		list := make([]v2cEntry, 0, len(e.V2cServers))
		for _, x := range e.V2cServers {
			list = append(list, v2cEntry{
				Name:      x.Name,
				Manager:   x.Manager,
				Community: x.Community,
			})
		}
		ans.V2c = &v2cs{Entries: list}
	}

	if len(e.V3Servers) > 0 {
		list := make([]v3Entry, 0, len(e.V3Servers))
		for _, x := range e.V3Servers {
			list = append(list, v3Entry{
				Name:         x.Name,
				Manager:      x.Manager,
				User:         x.User,
				EngineId:     x.EngineId,
				AuthPassword: x.AuthPassword,
				PrivPassword: x.PrivPassword,
			})
		}
		ans.V3 = &v3s{Entries: list}
	}

	return ans
}
