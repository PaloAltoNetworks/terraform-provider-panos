package variable

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a template
// variable.
//
// Template variables are a new addition to PAN-OS 8.1.
type Entry struct {
	Name  string
	Type  string
	Value string
}

// Copy copies the information from source's Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Type = s.Type
	o.Value = s.Value
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
		Name: o.Name,
	}

	if o.IpNetmask != "" {
		ans.Type = TypeIpNetmask
		ans.Value = o.IpNetmask
	} else if o.IpRange != "" {
		ans.Type = TypeIpRange
		ans.Value = o.IpRange
	} else if o.Fqdn != "" {
		ans.Type = TypeFqdn
		ans.Value = o.Fqdn
	} else if o.GroupId != "" {
		ans.Type = TypeGroupId
		ans.Value = o.GroupId
	} else if o.Interface != "" {
		ans.Type = TypeInterface
		ans.Value = o.Interface
	}

	return ans
}

type entry_v1 struct {
	XMLName   xml.Name `xml:"entry"`
	Name      string   `xml:"name,attr"`
	IpNetmask string   `xml:"type>ip-netmask,omitempty"`
	IpRange   string   `xml:"type>ip-range,omitempty"`
	Fqdn      string   `xml:"type>fqdn,omitempty"`
	GroupId   string   `xml:"type>group-id,omitempty"`
	Interface string   `xml:"type>interface,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	switch e.Type {
	case TypeIpNetmask:
		ans.IpNetmask = e.Value
	case TypeIpRange:
		ans.IpRange = e.Value
	case TypeFqdn:
		ans.Fqdn = e.Value
	case TypeGroupId:
		ans.GroupId = e.Value
	case TypeInterface:
		ans.Interface = e.Value
	}

	return ans
}
