package filetype

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// DLP file type.
type Entry struct {
	Name       string
	Properties []Property
}

type Property struct {
	Name  string
	Label string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	if s.Properties == nil {
		o.Properties = nil
	} else {
		o.Properties = make([]Property, 0, len(s.Properties))
		for _, x := range s.Properties {
			o.Properties = append(o.Properties, Property{
				Name:  x.Name,
				Label: x.Label,
			})
		}
	}
}

/** Structs / functions for normalization. **/

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

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v1 struct {
	XMLName  xml.Name `xml:"entry"`
	Name     string   `xml:"name,attr"`
	Property *prop    `xml:"file-property"`
}

type prop struct {
	Entries []propEntry `xml:"entry"`
}

type propEntry struct {
	Name  string `xml:"name,attr"`
	Label string `xml:"label"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	if len(e.Properties) > 0 {
		list := make([]propEntry, 0, len(e.Properties))
		for _, x := range e.Properties {
			list = append(list, propEntry{
				Name:  x.Name,
				Label: x.Label,
			})
		}
		ans.Property = &prop{Entries: list}
	}

	return ans
}

func (e *entry_v1) normalize() Entry {
	ans := Entry{
		Name: e.Name,
	}

	if e.Property != nil {
		ans.Properties = make([]Property, 0, len(e.Property.Entries))
		for _, x := range e.Property.Entries {
			ans.Properties = append(ans.Properties, Property{
				Name:  x.Name,
				Label: x.Label,
			})
		}
	}

	return ans
}
