package url

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of
// a custom URL category.
type Entry struct {
	Name        string
	Description string
	Sites       []string // Ordered
	Type        string   // PAN-OS 9.0
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.Sites = util.CopyStringSlice(s.Sites)
	o.Type = s.Type
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
	XMLName     xml.Name         `xml:"entry"`
	Name        string           `xml:"name,attr"`
	Description string           `xml:"description,omitempty"`
	Sites       *util.MemberType `xml:"list"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		Sites:       util.StrToMem(e.Sites),
	}

	return ans
}

func (e *entry_v1) normalize() Entry {
	ans := Entry{
		Name:        e.Name,
		Description: e.Description,
		Sites:       util.MemToStr(e.Sites),
	}

	return ans
}

// PAN-OS 9.0
type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v2 struct {
	XMLName     xml.Name         `xml:"entry"`
	Name        string           `xml:"name,attr"`
	Description string           `xml:"description,omitempty"`
	Sites       *util.MemberType `xml:"list"`
	Type        string           `xml:"type"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:        e.Name,
		Description: e.Description,
		Sites:       util.StrToMem(e.Sites),
		Type:        e.Type,
	}

	return ans
}

func (e *entry_v2) normalize() Entry {
	ans := Entry{
		Name:        e.Name,
		Description: e.Description,
		Sites:       util.MemToStr(e.Sites),
		Type:        e.Type,
	}

	return ans
}
