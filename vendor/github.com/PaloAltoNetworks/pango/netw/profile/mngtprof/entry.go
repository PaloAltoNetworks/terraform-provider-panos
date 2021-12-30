package mngtprof

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an interface
// management profile.
type Entry struct {
	Name                    string
	Ping                    bool
	Telnet                  bool
	Ssh                     bool
	Http                    bool
	HttpOcsp                bool
	Https                   bool
	Snmp                    bool
	ResponsePages           bool
	UseridService           bool
	UseridSyslogListenerSsl bool
	UseridSyslogListenerUdp bool
	PermittedIps            []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Ping = s.Ping
	o.Telnet = s.Telnet
	o.Ssh = s.Ssh
	o.Http = s.Http
	o.HttpOcsp = s.HttpOcsp
	o.Https = s.Https
	o.Snmp = s.Snmp
	o.ResponsePages = s.ResponsePages
	o.UseridService = s.UseridService
	o.UseridSyslogListenerSsl = s.UseridSyslogListenerSsl
	o.UseridSyslogListenerUdp = s.UseridSyslogListenerUdp
	if s.PermittedIps == nil {
		o.PermittedIps = nil
	} else {
		o.PermittedIps = make([]string, len(s.PermittedIps))
		copy(o.PermittedIps, s.PermittedIps)
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
		Name:                    o.Name,
		Ping:                    util.AsBool(o.Ping),
		Telnet:                  util.AsBool(o.Telnet),
		Ssh:                     util.AsBool(o.Ssh),
		Http:                    util.AsBool(o.Http),
		HttpOcsp:                util.AsBool(o.HttpOcsp),
		Https:                   util.AsBool(o.Https),
		Snmp:                    util.AsBool(o.Snmp),
		ResponsePages:           util.AsBool(o.ResponsePages),
		UseridService:           util.AsBool(o.UseridService),
		UseridSyslogListenerSsl: util.AsBool(o.UseridSyslogListenerSsl),
		UseridSyslogListenerUdp: util.AsBool(o.UseridSyslogListenerUdp),
		PermittedIps:            util.EntToStr(o.PermittedIps),
	}

	return ans
}

type entry_v1 struct {
	XMLName                 xml.Name        `xml:"entry"`
	Name                    string          `xml:"name,attr"`
	Ping                    string          `xml:"ping"`
	Telnet                  string          `xml:"telnet"`
	Ssh                     string          `xml:"ssh"`
	Http                    string          `xml:"http"`
	HttpOcsp                string          `xml:"http-ocsp"`
	Https                   string          `xml:"https"`
	Snmp                    string          `xml:"snmp"`
	ResponsePages           string          `xml:"response-pages"`
	UseridService           string          `xml:"userid-service"`
	UseridSyslogListenerSsl string          `xml:"userid-syslog-listener-ssl"`
	UseridSyslogListenerUdp string          `xml:"userid-syslog-listener-udp"`
	PermittedIps            *util.EntryType `xml:"permitted-ip"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                    e.Name,
		Ping:                    util.YesNo(e.Ping),
		Telnet:                  util.YesNo(e.Telnet),
		Ssh:                     util.YesNo(e.Ssh),
		Http:                    util.YesNo(e.Http),
		HttpOcsp:                util.YesNo(e.HttpOcsp),
		Https:                   util.YesNo(e.Https),
		Snmp:                    util.YesNo(e.Snmp),
		ResponsePages:           util.YesNo(e.ResponsePages),
		UseridService:           util.YesNo(e.UseridService),
		UseridSyslogListenerSsl: util.YesNo(e.UseridSyslogListenerSsl),
		UseridSyslogListenerUdp: util.YesNo(e.UseridSyslogListenerUdp),
		PermittedIps:            util.StrToEnt(e.PermittedIps),
	}

	return ans
}
