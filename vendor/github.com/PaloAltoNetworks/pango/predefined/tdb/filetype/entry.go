package filetype

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// threat.
type Entry struct {
	Name          string
	Id            int
	ThreatName    string
	FullName      string
	DataIdent     bool
	FileTypeIdent bool
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.ThreatName = s.ThreatName
	o.Id = s.Id
	o.ThreatName = s.ThreatName
	o.FullName = s.FullName
	o.DataIdent = s.DataIdent
	o.FileTypeIdent = s.FileTypeIdent
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
	XMLName       xml.Name `xml:"entry"`
	Name          string   `xml:"name,attr"`
	Id            int      `xml:"id,attr"`
	DataIdent     string   `xml:"data-ident,omitempty"`
	FileTypeIdent string   `xml:"file-type-ident"`
	ThreatName    string   `xml:"threat-name,omitempty"`
	FullName      string   `xml:"full-name"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:          e.Name,
		Id:            e.Id,
		FileTypeIdent: util.YesNo(e.FileTypeIdent),
		ThreatName:    e.ThreatName,
		FullName:      e.FullName,
	}

	if e.DataIdent {
		ans.DataIdent = "yes"
	}

	return ans
}

func (e *entry_v1) normalize() Entry {
	ans := Entry{
		Name:          e.Name,
		Id:            e.Id,
		DataIdent:     util.AsBool(e.DataIdent),
		FileTypeIdent: util.AsBool(e.FileTypeIdent),
		ThreatName:    e.ThreatName,
		FullName:      e.FullName,
	}

	return ans
}
