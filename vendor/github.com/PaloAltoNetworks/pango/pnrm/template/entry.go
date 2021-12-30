package template

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a template.
//
// Devices is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
type Entry struct {
	Name           string
	Description    string
	DefaultVsys    string
	MultiVsys      bool
	Mode           string
	VpnDisableMode bool
	Devices        map[string][]string

	raw map[string]string
}

// Copy copies the information from source's Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.DefaultVsys = s.DefaultVsys
	o.MultiVsys = s.MultiVsys
	o.Mode = s.Mode
	o.VpnDisableMode = s.VpnDisableMode
	if s.Devices == nil {
		o.Devices = nil
	} else {
		o.Devices = make(map[string][]string)
		for key, list := range s.Devices {
			if list == nil {
				o.Devices[key] = nil
			} else {
				o.Devices[key] = make([]string, len(list))
				copy(o.Devices[key], list)
			}
		}
	}
}

// SetConfTree sets the conf internal variable such that the XML contains
// the mandatory "/config" subelement tree.
//
// If a template is missing this, then it does not behave properly when
// referenced from a template stack.
func (o *Entry) SetConfTree() {
	if _, present := o.raw["conf"]; !present {
		if o.raw == nil {
			o.raw = make(map[string]string)
		}
		o.raw["conf"] = "<devices><entry name='localhost.localdomain'><vsys><entry name='vsys1' /></vsys></entry></devices>"
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
		Devices:     util.VsysEntToMap(o.Devices),
	}

	if o.Settings != nil {
		ans.MultiVsys = util.AsBool(o.Settings.MultiVsys)
		ans.Mode = o.Settings.Mode
		ans.VpnDisableMode = util.AsBool(o.Settings.VpnDisableMode)
	}

	raw := make(map[string]string)

	if o.Config != nil {
		raw["conf"] = util.CleanRawXml(o.Config.Text)
	}

	if len(raw) > 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name            `xml:"entry"`
	Name        string              `xml:"name,attr"`
	Description string              `xml:"description,omitempty"`
	Devices     *util.VsysEntryType `xml:"devices"`
	Settings    *settings_v1        `xml:"settings"`
	Config      *util.RawXml        `xml:"config"`
}

type settings_v1 struct {
	MultiVsys      string `xml:"multi-vsys"`
	Mode           string `xml:"operational-mode,omitempty"`
	VpnDisableMode string `xml:"vpn-disable-mode"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		Devices:     util.MapToVsysEnt(e.Devices),
	}

	if e.MultiVsys || e.VpnDisableMode || e.Mode != "" {
		ans.Settings = &settings_v1{
			MultiVsys:      util.YesNo(e.MultiVsys),
			Mode:           e.Mode,
			VpnDisableMode: util.YesNo(e.VpnDisableMode),
		}
	}

	if text, present := e.raw["conf"]; present {
		ans.Config = &util.RawXml{text}
	}

	return ans
}

// PAN-OS 7.0.
type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:        o.Name,
		Description: o.Description,
		Devices:     util.VsysEntToMap(o.Devices),
	}

	if o.Settings != nil {
		ans.DefaultVsys = o.Settings.DefaultVsys
	}

	raw := make(map[string]string)

	if o.Config != nil {
		raw["conf"] = util.CleanRawXml(o.Config.Text)
	}

	if len(raw) > 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v2 struct {
	XMLName     xml.Name            `xml:"entry"`
	Name        string              `xml:"name,attr"`
	Description string              `xml:"description,omitempty"`
	Devices     *util.VsysEntryType `xml:"devices"`
	Settings    *settings_v2        `xml:"settings"`
	Config      *util.RawXml        `xml:"config"`
}

type settings_v2 struct {
	DefaultVsys string `xml:"default-vsys"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:        e.Name,
		Description: e.Description,
		Devices:     util.MapToVsysEnt(e.Devices),
	}

	if e.DefaultVsys != "" {
		ans.Settings = &settings_v2{e.DefaultVsys}
	}

	if text, present := e.raw["conf"]; present {
		ans.Config = &util.RawXml{text}
	}

	return ans
}

// PAN-OS 8.1.
type container_v3 struct {
	Answer []entry_v3 `xml:"entry"`
}

func (o *container_v3) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v3) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v3) normalize() Entry {
	ans := Entry{
		Name:        o.Name,
		Description: o.Description,
		// TODO(gfreeman) - seems like devices are removed in 8.1..?
		Devices: util.VsysEntToMap(o.Devices),
	}

	if o.Settings != nil {
		ans.DefaultVsys = o.Settings.DefaultVsys
	}

	raw := make(map[string]string)

	if o.Variables != nil {
		raw["vars"] = util.CleanRawXml(o.Variables.Text)
	}

	if o.Config != nil {
		raw["conf"] = util.CleanRawXml(o.Config.Text)
	}

	if len(raw) > 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v3 struct {
	XMLName     xml.Name            `xml:"entry"`
	Name        string              `xml:"name,attr"`
	Description string              `xml:"description,omitempty"`
	Devices     *util.VsysEntryType `xml:"devices"`
	Settings    *settings_v2        `xml:"settings"`
	Config      *util.RawXml        `xml:"config"`
	Variables   *util.RawXml        `xml:"variable"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:        e.Name,
		Description: e.Description,
		Devices:     util.MapToVsysEnt(e.Devices),
	}

	if e.DefaultVsys != "" {
		ans.Settings = &settings_v2{e.DefaultVsys}
	}

	if text, present := e.raw["conf"]; present {
		ans.Config = &util.RawXml{text}
	}

	if text, present := e.raw["vars"]; present {
		ans.Variables = &util.RawXml{text}
	}

	return ans
}
