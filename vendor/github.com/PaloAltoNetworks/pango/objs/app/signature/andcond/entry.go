package andcond

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of an application signature and-condition.
type Entry struct {
	Name string

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
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
		Name: o.Answer.Name,
	}

	ans.raw = make(map[string]string)

	if o.Answer.Sigs != nil {
		ans.raw["sigs"] = util.CleanRawXml(o.Answer.Sigs.Text)
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name     `xml:"entry"`
	Name    string       `xml:"name,attr"`
	Sigs    *util.RawXml `xml:"or-condition"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	if text := e.raw["sigs"]; text != "" {
		ans.Sigs = &util.RawXml{text}
	}

	return ans
}
