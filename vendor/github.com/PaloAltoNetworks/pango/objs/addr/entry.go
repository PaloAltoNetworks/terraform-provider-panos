package addr

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an address
// object.
type Entry struct {
	Name        string
	Value       string
	Type        string
	Description string
	Tags        []string // ordered
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Value = s.Value
	o.Type = s.Type
	o.Description = s.Description
	o.Tags = s.Tags
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

type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name         `xml:"entry"`
	Name        string           `xml:"name,attr"`
	IpNetmask   *valType         `xml:"ip-netmask"`
	IpRange     *valType         `xml:"ip-range"`
	Fqdn        *valType         `xml:"fqdn"`
	Description string           `xml:"description"`
	Tags        *util.MemberType `xml:"tag"`
}

type valType struct {
	Value string `xml:",chardata"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		Tags:        util.StrToMem(e.Tags),
	}
	vt := &valType{e.Value}
	switch e.Type {
	case IpNetmask:
		ans.IpNetmask = vt
	case IpRange:
		ans.IpRange = vt
	case Fqdn:
		ans.Fqdn = vt
	}

	return ans
}

func (e *entry_v1) normalize() Entry {
	ans := Entry{
		Name:        e.Name,
		Description: e.Description,
		Tags:        util.MemToStr(e.Tags),
	}

	switch {
	case e.IpNetmask != nil:
		ans.Type = IpNetmask
		ans.Value = e.IpNetmask.Value
	case e.IpRange != nil:
		ans.Type = IpRange
		ans.Value = e.IpRange.Value
	case e.Fqdn != nil:
		ans.Type = Fqdn
		ans.Value = e.Fqdn.Value
	}

	return ans
}

type entry_v2 struct {
	XMLName     xml.Name         `xml:"entry"`
	Name        string           `xml:"name,attr"`
	IpNetmask   *valType         `xml:"ip-netmask"`
	IpRange     *valType         `xml:"ip-range"`
	Fqdn        *valType         `xml:"fqdn"`
	IpWildcard  *valType         `xml:"ip-wildcard"`
	Description string           `xml:"description"`
	Tags        *util.MemberType `xml:"tag"`
}

func (e *entry_v2) normalize() Entry {
	ans := Entry{
		Name:        e.Name,
		Description: e.Description,
		Tags:        util.MemToStr(e.Tags),
	}

	switch {
	case e.IpNetmask != nil:
		ans.Type = IpNetmask
		ans.Value = e.IpNetmask.Value
	case e.IpRange != nil:
		ans.Type = IpRange
		ans.Value = e.IpRange.Value
	case e.Fqdn != nil:
		ans.Type = Fqdn
		ans.Value = e.Fqdn.Value
	case e.IpWildcard != nil:
		ans.Type = IpWildcard
		ans.Value = e.IpWildcard.Value
	}

	return ans
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:        e.Name,
		Description: e.Description,
		Tags:        util.StrToMem(e.Tags),
	}
	vt := &valType{e.Value}
	switch e.Type {
	case IpNetmask:
		ans.IpNetmask = vt
	case IpRange:
		ans.IpRange = vt
	case Fqdn:
		ans.Fqdn = vt
	case IpWildcard:
		ans.IpWildcard = vt
	}

	return ans
}
