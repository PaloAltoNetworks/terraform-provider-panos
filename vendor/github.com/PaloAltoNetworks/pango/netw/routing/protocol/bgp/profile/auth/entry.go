package auth

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an auth profile.
type Entry struct {
	Name   string
	Secret string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Secret = s.Secret
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
		Name:   o.Name,
		Secret: o.Secret,
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	Secret  string   `xml:"secret,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:   e.Name,
		Secret: e.Secret,
	}

	return ans
}
