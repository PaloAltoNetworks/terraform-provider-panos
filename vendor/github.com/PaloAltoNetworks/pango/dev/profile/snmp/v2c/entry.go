package v2c

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of a snmptrap v2c server.
//
// PAN-OS 7.1+.
type Entry struct {
	Name      string
	Manager   string
	Community string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Manager = s.Manager
	o.Community = s.Community
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
		Name:      o.Answer.Name,
		Manager:   o.Answer.Manager,
		Community: o.Answer.Community,
	}

	return ans
}

type entry_v1 struct {
	XMLName   xml.Name `xml:"entry"`
	Name      string   `xml:"name,attr"`
	Manager   string   `xml:"manager"`
	Community string   `xml:"community"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:      e.Name,
		Manager:   e.Manager,
		Community: e.Community,
	}

	return ans
}
