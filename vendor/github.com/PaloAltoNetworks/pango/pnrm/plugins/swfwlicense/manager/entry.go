package manager

import (
	"encoding/xml"
	"github.com/PaloAltoNetworks/pango/plugin"
)

// Version Independent Data Structure
type Entry struct {
	Name                string
	Description         string
	DeviceGroup         string
	TemplateStack       string
	BootstrapDefinition string
	AutoDeactivate      int
}

// normalizer interface refers to container_v1
type normalizer interface {
	Names() []string
	Normalize() []Entry
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.DeviceGroup = s.DeviceGroup
	o.TemplateStack = s.TemplateStack
	o.BootstrapDefinition = s.BootstrapDefinition
	o.AutoDeactivate = s.AutoDeactivate
}

func (o Entry) Specify(list []plugin.Info) (string, interface{}, error) {
	_, fn, err := versioning(list)
	if err != nil {
		return o.Name, nil, err
	}

	return o.Name, fn(o), nil
}

// type container_v1 contains Answer which is a slice of entry_v1 from the XML entry.
type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

// Normalize function that returns a slice of Entry taking each item in the slice from the Answer array in the
// container_v1 struct
func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

// Names function that returns a slice of String containing all the Name fields  from the entries in the Answer array
// in the container_v1 struct
func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

// entry_v1 data structure from XML
type entry_v1 struct {
	XMLName             xml.Name `xml:"entry"`
	Name                string   `xml:"name,attr"`
	Description         string   `xml:"description,omitempty"`
	DeviceGroup         string   `xml:"device-group"`
	TemplateStack       string   `xml:"template-stack"`
	BootstrapDefinition string   `xml:"bootstrap-definition"`
	AutoDeactivate      int      `xml:"auto-deactivate"`
}

// entry_v1 function normalize returns version independent Entry
func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                o.Name,
		Description:         o.Description,
		DeviceGroup:         o.DeviceGroup,
		TemplateStack:       o.TemplateStack,
		BootstrapDefinition: o.BootstrapDefinition,
		AutoDeactivate:      o.AutoDeactivate,
	}

	return ans
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                e.Name,
		Description:         e.Description,
		DeviceGroup:         e.DeviceGroup,
		TemplateStack:       e.TemplateStack,
		BootstrapDefinition: e.BootstrapDefinition,
		AutoDeactivate:      e.AutoDeactivate,
	}

	return ans
}
