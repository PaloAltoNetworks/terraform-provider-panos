package param

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of an http param.
//
// PAN-OS 7.1+.
type Entry struct {
	Name  string
	Value string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Value = s.Value
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
		Name:  o.Answer.Name,
		Value: o.Answer.Value,
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:  e.Name,
		Value: e.Value,
	}

	return ans
}
