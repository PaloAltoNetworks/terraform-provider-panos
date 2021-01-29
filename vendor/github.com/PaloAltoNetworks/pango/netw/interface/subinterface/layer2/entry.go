package layer2

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a layer2
// subinterface.
type Entry struct {
	Name           string
	Tag            int
	NetflowProfile string
	Comment        string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Tag = s.Tag
	o.NetflowProfile = s.NetflowProfile
	o.Comment = s.Comment
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, string, interface{}) {
	_, fn := versioning(v)
	return o.Name, o.Name, fn(o)
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
		Name:           o.Name,
		Tag:            o.Tag,
		NetflowProfile: o.NetflowProfile,
		Comment:        o.Comment,
	}

	return ans
}

type entry_v1 struct {
	XMLName        xml.Name `xml:"entry"`
	Name           string   `xml:"name,attr"`
	Tag            int      `xml:"tag,omitempty"`
	NetflowProfile string   `xml:"netflow-profile,omitempty"`
	Comment        string   `xml:"comment,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:           e.Name,
		Tag:            e.Tag,
		NetflowProfile: e.NetflowProfile,
		Comment:        e.Comment,
	}

	return ans
}
