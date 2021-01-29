package v3

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of a snmptrap v3 server.
//
// PAN-OS 7.1+.
type Entry struct {
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
	o.Manager = s.Manager
	o.User = s.User
	o.EngineId = s.EngineId
	o.AuthPassword = s.AuthPassword
	o.PrivPassword = s.PrivPassword
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
		Manager:      o.Answer.Manager,
		User:         o.Answer.User,
		EngineId:     o.Answer.EngineId,
		AuthPassword: o.Answer.AuthPassword,
		PrivPassword: o.Answer.PrivPassword,
	}

	return ans
}

type entry_v1 struct {
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
		Name:         e.Name,
		Manager:      e.Manager,
		User:         e.User,
		EngineId:     e.EngineId,
		AuthPassword: e.AuthPassword,
		PrivPassword: e.PrivPassword,
	}

	return ans
}
