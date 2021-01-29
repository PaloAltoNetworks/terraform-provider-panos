package variable

import (
	"encoding/xml"
)

// These are the constants for the Type field.
const (
	TypeIpNetmask = "ip-netmask"
	TypeIpRange   = "ip-range"
	TypeFqdn      = "fqdn"
	TypeGroupId   = "group-id"
	TypeInterface = "interface"
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

	if o.Answer.IpNetmask != "" {
		ans.Type = TypeIpNetmask
		ans.Value = o.Answer.IpNetmask
	} else if o.Answer.IpRange != "" {
		ans.Type = TypeIpRange
		ans.Value = o.Answer.IpRange
	} else if o.Answer.Fqdn != "" {
		ans.Type = TypeFqdn
		ans.Value = o.Answer.Fqdn
	} else if o.Answer.GroupId != "" {
		ans.Type = TypeGroupId
		ans.Value = o.Answer.GroupId
	} else if o.Answer.Interface != "" {
		ans.Type = TypeInterface
		ans.Value = o.Answer.Interface
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
