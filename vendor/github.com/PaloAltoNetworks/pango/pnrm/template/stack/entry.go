package stack

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
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
	if s.Templates == nil {
		o.Templates = nil
	} else {
		o.Templates = make([]string, len(s.Templates))
		copy(o.Templates, s.Templates)
	}
	if s.Devices == nil {
		o.Devices = nil
	} else {
		o.Devices = make([]string, len(s.Devices))
		copy(o.Devices, s.Devices)
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
		Name:        o.Name,
		Description: o.Description,
		DefaultVsys: o.DefaultVsys,
		Templates:   util.MemToStr(o.Templates),
		Devices:     util.EntToStr(o.Devices),
	}

	raw := make(map[string]string)

	if o.Variables != nil {
		raw["var"] = util.CleanRawXml(o.Variables.Text)
	}

	if o.Config != nil {
		raw["conf"] = util.CleanRawXml(o.Config.Text)
	}

	if len(raw) > 0 {
		ans.raw = raw
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
