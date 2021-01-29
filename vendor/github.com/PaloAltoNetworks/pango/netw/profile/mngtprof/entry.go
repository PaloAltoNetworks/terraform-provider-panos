package mngtprof

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
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
	o.PermittedIps = s.PermittedIps
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
		Name:                    o.Answer.Name,
		Ping:                    util.AsBool(o.Answer.Ping),
		Telnet:                  util.AsBool(o.Answer.Telnet),
		Ssh:                     util.AsBool(o.Answer.Ssh),
		Http:                    util.AsBool(o.Answer.Http),
		HttpOcsp:                util.AsBool(o.Answer.HttpOcsp),
		Https:                   util.AsBool(o.Answer.Https),
		Snmp:                    util.AsBool(o.Answer.Snmp),
		ResponsePages:           util.AsBool(o.Answer.ResponsePages),
		UseridService:           util.AsBool(o.Answer.UseridService),
		UseridSyslogListenerSsl: util.AsBool(o.Answer.UseridSyslogListenerSsl),
		UseridSyslogListenerUdp: util.AsBool(o.Answer.UseridSyslogListenerUdp),
		PermittedIps:            util.EntToStr(o.Answer.PermittedIps),
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
