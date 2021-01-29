package group

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of an application group.
type Entry struct {
	Name         string
	Applications []string // ordered
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Applications = s.Applications
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
		Applications: util.MemToStr(o.Answer.Applications),
	}

	return ans
}

type entry_v1 struct {
	XMLName      xml.Name         `xml:"entry"`
	Name         string           `xml:"name,attr"`
	Applications *util.MemberType `xml:"members"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:         e.Name,
		Applications: util.StrToMem(e.Applications),
	}

	return ans
}
