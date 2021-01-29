package group

import (
	"encoding/xml"
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

type normalizer interface {
	Normalize() Entry
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
	ans := Entry{
		Name:                 o.Answer.Name,
		Description:          o.Answer.Description,
		GcpProjectCredential: o.Answer.GcpProjectCredential,
		DeviceGroup:          o.Answer.DeviceGroup,
		TemplateStack:        o.Answer.TemplateStack,
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
