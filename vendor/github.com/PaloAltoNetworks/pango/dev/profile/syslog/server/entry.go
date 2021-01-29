package server

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of an syslog server.
//
// PAN-OS 7.1+.
type Entry struct {
	Name         string
	Server       string
	Transport    string
	Port         int
	SyslogFormat string
	Facility     string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Server = s.Server
	o.Transport = s.Transport
	o.Port = s.Port
	o.SyslogFormat = s.SyslogFormat
	o.Facility = s.Facility
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
		Name:         o.Answer.Name,
		Server:       o.Answer.Server,
		Transport:    o.Answer.Transport,
		Port:         o.Answer.Port,
		SyslogFormat: o.Answer.SyslogFormat,
		Facility:     o.Answer.Facility,
	}

	return ans
}

type entry_v1 struct {
	XMLName      xml.Name `xml:"entry"`
	Name         string   `xml:"name,attr"`
	Server       string   `xml:"server"`
	Transport    string   `xml:"transport,omitempty"`
	Port         int      `xml:"port,omitempty"`
	SyslogFormat string   `xml:"format,omitempty"`
	Facility     string   `xml:"facility,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:         e.Name,
		Server:       e.Server,
		Transport:    e.Transport,
		Port:         e.Port,
		SyslogFormat: e.SyslogFormat,
		Facility:     e.Facility,
	}

	return ans
}
