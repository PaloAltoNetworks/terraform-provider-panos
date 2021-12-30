package group

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/plugin"
)

// Entry is a normalized, version independent representation of a GKE cluster group.
type Entry struct {
	Name                 string
	Description          string
	GcpProjectCredential string
	DeviceGroup          string
	TemplateStack        string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.GcpProjectCredential = s.GcpProjectCredential
	o.DeviceGroup = s.DeviceGroup
	o.TemplateStack = s.TemplateStack
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(list []plugin.Info) (string, interface{}, error) {
	_, fn, err := versioning(list)
	if err != nil {
		return o.Name, nil, err
	}

	return o.Name, fn(o), nil
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
		Name:                 o.Name,
		Description:          o.Description,
		GcpProjectCredential: o.GcpProjectCredential,
		DeviceGroup:          o.DeviceGroup,
		TemplateStack:        o.TemplateStack,
	}

	return ans
}

type entry_v1 struct {
	XMLName              xml.Name `xml:"entry"`
	Name                 string   `xml:"name,attr"`
	Description          string   `xml:"description,omitempty"`
	GcpProjectCredential string   `xml:"gcp-creds"`
	DeviceGroup          string   `xml:"device-group"`
	TemplateStack        string   `xml:"template-stack"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                 e.Name,
		Description:          e.Description,
		GcpProjectCredential: e.GcpProjectCredential,
		DeviceGroup:          e.DeviceGroup,
		TemplateStack:        e.TemplateStack,
	}

	return ans
}
