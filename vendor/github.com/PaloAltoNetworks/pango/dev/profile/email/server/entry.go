package server

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of an email server.
//
// PAN-OS 7.1+.
type Entry struct {
	Name         string
	DisplayName  string
	From         string
	To           string
	AlsoTo       string
	EmailGateway string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.DisplayName = s.DisplayName
	o.From = s.From
	o.To = s.To
	o.AlsoTo = s.AlsoTo
	o.EmailGateway = s.EmailGateway
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
		DisplayName:  o.Answer.DisplayName,
		From:         o.Answer.From,
		To:           o.Answer.To,
		AlsoTo:       o.Answer.AlsoTo,
		EmailGateway: o.Answer.EmailGateway,
	}

	return ans
}

type entry_v1 struct {
	XMLName      xml.Name `xml:"entry"`
	Name         string   `xml:"name,attr"`
	DisplayName  string   `xml:"display-name,omitempty"`
	From         string   `xml:"from"`
	To           string   `xml:"to"`
	AlsoTo       string   `xml:"and-also-to,omitempty"`
	EmailGateway string   `xml:"gateway"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:         e.Name,
		DisplayName:  e.DisplayName,
		From:         e.From,
		To:           e.To,
		AlsoTo:       e.AlsoTo,
		EmailGateway: e.EmailGateway,
	}

	return ans
}
