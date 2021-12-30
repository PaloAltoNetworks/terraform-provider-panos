package vsys

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a Vsys.
type Entry struct {
	Name           string
	NetworkImports *NetworkImports
}

type NetworkImports struct {
	Interfaces     []string
	VirtualRouters []string
	VirtualWires   []string
	Vlans          []string
	LogicalRouters []string // PAN-OS 10.0+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	if s.NetworkImports == nil {
		o.NetworkImports = nil
	} else {
		o.NetworkImports = &NetworkImports{}
		if s.NetworkImports.Interfaces != nil {
			o.NetworkImports.Interfaces = make([]string, len(s.NetworkImports.Interfaces))
			copy(o.NetworkImports.Interfaces, s.NetworkImports.Interfaces)
		}
		if s.NetworkImports.VirtualRouters != nil {
			o.NetworkImports.VirtualRouters = make([]string, len(s.NetworkImports.VirtualRouters))
			copy(o.NetworkImports.VirtualRouters, s.NetworkImports.VirtualRouters)
		}
		if s.NetworkImports.VirtualWires != nil {
			o.NetworkImports.VirtualWires = make([]string, len(s.NetworkImports.VirtualWires))
			copy(o.NetworkImports.VirtualWires, s.NetworkImports.VirtualWires)
		}
		if s.NetworkImports.Vlans != nil {
			o.NetworkImports.Vlans = make([]string, len(s.NetworkImports.Vlans))
			copy(o.NetworkImports.Vlans, s.NetworkImports.Vlans)
		}
		if s.NetworkImports.LogicalRouters != nil {
			o.NetworkImports.LogicalRouters = make([]string, len(s.NetworkImports.LogicalRouters))
			copy(o.NetworkImports.LogicalRouters, s.NetworkImports.LogicalRouters)
		}
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

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	Import  *imp_v1  `xml:"import"`
}

type imp_v1 struct {
	Network *impNetwork_v1 `xml:"network"`
}

type impNetwork_v1 struct {
	Interfaces     *util.MemberType `xml:"interface"`
	VirtualRouters *util.MemberType `xml:"virtual-router"`
	VirtualWires   *util.MemberType `xml:"virtual-wire"`
	Vlans          *util.MemberType `xml:"vlan"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	if e.NetworkImports != nil {
		ans.Import = &imp_v1{
			Network: &impNetwork_v1{
				Interfaces:     util.StrToMem(e.NetworkImports.Interfaces),
				VirtualRouters: util.StrToMem(e.NetworkImports.VirtualRouters),
				VirtualWires:   util.StrToMem(e.NetworkImports.VirtualWires),
				Vlans:          util.StrToMem(e.NetworkImports.Vlans),
			},
		}
	}

	return ans
}

func (e *entry_v1) normalize() Entry {
	ans := Entry{
		Name: e.Name,
	}

	if e.Import != nil {
		if e.Import.Network != nil {
			ans.NetworkImports = &NetworkImports{
				Interfaces:     util.MemToStr(e.Import.Network.Interfaces),
				VirtualRouters: util.MemToStr(e.Import.Network.VirtualRouters),
				VirtualWires:   util.MemToStr(e.Import.Network.VirtualWires),
				Vlans:          util.MemToStr(e.Import.Network.Vlans),
			}
		}
	}

	return ans
}

// PAN-OS 10.0
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

type entry_v2 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	Import  *imp_v2  `xml:"import"`
}

type imp_v2 struct {
	Network *impNetwork_v2 `xml:"network"`
}

type impNetwork_v2 struct {
	Interfaces     *util.MemberType `xml:"interface"`
	VirtualRouters *util.MemberType `xml:"virtual-router"`
	VirtualWires   *util.MemberType `xml:"virtual-wire"`
	Vlans          *util.MemberType `xml:"vlan"`
	LogicalRouters *util.MemberType `xml:"logical-router"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name: e.Name,
	}

	if e.NetworkImports != nil {
		ans.Import = &imp_v2{
			Network: &impNetwork_v2{
				Interfaces:     util.StrToMem(e.NetworkImports.Interfaces),
				VirtualRouters: util.StrToMem(e.NetworkImports.VirtualRouters),
				VirtualWires:   util.StrToMem(e.NetworkImports.VirtualWires),
				Vlans:          util.StrToMem(e.NetworkImports.Vlans),
				LogicalRouters: util.StrToMem(e.NetworkImports.LogicalRouters),
			},
		}
	}

	return ans
}

func (e *entry_v2) normalize() Entry {
	ans := Entry{
		Name: e.Name,
	}

	if e.Import != nil {
		if e.Import.Network != nil {
			ans.NetworkImports = &NetworkImports{
				Interfaces:     util.MemToStr(e.Import.Network.Interfaces),
				VirtualRouters: util.MemToStr(e.Import.Network.VirtualRouters),
				VirtualWires:   util.MemToStr(e.Import.Network.VirtualWires),
				Vlans:          util.MemToStr(e.Import.Network.Vlans),
				LogicalRouters: util.MemToStr(e.Import.Network.LogicalRouters),
			}
		}
	}

	return ans
}
