package stack

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a template stack.
//
// Devices is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
	Name        string
	Description string
	DefaultVsys string
	Templates   []string
	Devices     []string

	raw map[string]string
}

// Copy copies the information from source's Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.DefaultVsys = s.DefaultVsys
	o.Templates = s.Templates
	o.Devices = s.Devices
}

/** Structs / functions for normalization. **/

type normalizer interface {
	Normalize() Entry
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
	ans := Entry{
		Name:        o.Answer.Name,
		Description: o.Answer.Description,
		DefaultVsys: o.Answer.DefaultVsys,
		Templates:   util.MemToStr(o.Answer.Templates),
		Devices:     util.EntToStr(o.Answer.Devices),
	}

	ans.raw = make(map[string]string)

	if o.Answer.Variables != nil {
		ans.raw["var"] = util.CleanRawXml(o.Answer.Variables.Text)
	}

	if o.Answer.Config != nil {
		ans.raw["conf"] = util.CleanRawXml(o.Answer.Config.Text)
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name         `xml:"entry"`
	Name        string           `xml:"name,attr"`
	Description string           `xml:"description,omitempty"`
	DefaultVsys string           `xml:"settings>default-vsys,omitempty"`
	Templates   *util.MemberType `xml:"templates"`
	Devices     *util.EntryType  `xml:"devices"`
	Config      *util.RawXml     `xml:"config"`
	Variables   *util.RawXml     `xml:"variable"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		Devices:     util.StrToEnt(e.Devices),
		Templates:   util.StrToMem(e.Templates),
		DefaultVsys: e.DefaultVsys,
	}

	if text, present := e.raw["var"]; present {
		ans.Variables = &util.RawXml{text}
	}

	if text, present := e.raw["conf"]; present {
		ans.Config = &util.RawXml{text}
	}

	return ans
}
