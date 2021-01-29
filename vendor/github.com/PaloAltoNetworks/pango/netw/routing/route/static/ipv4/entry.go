package ipv4

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an IPv4
// static route.
type Entry struct {
	Name          string
	Destination   string
	Interface     string
	Type          string
	NextHop       string
	AdminDistance int
	Metric        int
	RouteTable    string
	BfdProfile    string
}

func (o *Entry) Copy(s Entry) {
	o.Destination = s.Destination
	o.Interface = s.Interface
	o.Type = s.Type
	o.NextHop = s.NextHop
	o.AdminDistance = s.AdminDistance
	o.Metric = s.Metric
	o.RouteTable = s.RouteTable
	o.BfdProfile = s.BfdProfile
}

/** Structs / functions for this namespace. **/

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
		Name:          o.Name,
		Destination:   o.Destination,
		Interface:     o.Interface,
		AdminDistance: o.AdminDistance,
		Metric:        o.Metric,
	}

	if o.NextHop == nil {
		ans.Type = ""
	} else if o.NextHop.Discard != nil {
		ans.Type = NextHopDiscard
	} else if o.NextHop.IpAddress != nil {
		ans.Type = NextHopIpAddress
		ans.NextHop = *o.NextHop.IpAddress
	} else if o.NextHop.NextVr != nil {
		ans.Type = NextHopNextVr
		ans.NextHop = *o.NextHop.NextVr
	}

	if o.Option != nil && o.Option.NoInstall != nil {
		ans.RouteTable = RouteTableNoInstall
	}

	return ans
}

type entry_v1 struct {
	XMLName       xml.Name     `xml:"entry"`
	Name          string       `xml:"name,attr"`
	Destination   string       `xml:"destination"`
	Interface     string       `xml:"interface,omitempty"`
	NextHop       *nextHop     `xml:"nexthop"`
	AdminDistance int          `xml:"admin-dist,omitempty"`
	Metric        int          `xml:"metric,omitempty"`
	Option        *rtOption_v1 `xml:"option"`
}

type nextHop struct {
	Discard   *string `xml:"discard"`
	IpAddress *string `xml:"ip-address"`
	NextVr    *string `xml:"next-vr"`
}

type rtOption_v1 struct {
	NoInstall *string `xml:"no-install"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:          e.Name,
		Destination:   e.Destination,
		Interface:     e.Interface,
		AdminDistance: e.AdminDistance,
		Metric:        e.Metric,
	}

	switch e.Type {
	case NextHopDiscard:
		var sp string
		ans.NextHop = &nextHop{Discard: &sp}
	case NextHopIpAddress:
		sp := e.NextHop
		ans.NextHop = &nextHop{IpAddress: &sp}
	case NextHopNextVr:
		sp := e.NextHop
		ans.NextHop = &nextHop{NextVr: &sp}
	}

	if e.RouteTable == RouteTableNoInstall {
		sp := ""
		ans.Option = &rtOption_v1{NoInstall: &sp}
	}

	return ans
}

// PAN-OS 7.1, adds BfdProfile
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
		Name:          o.Name,
		Destination:   o.Destination,
		Interface:     o.Interface,
		AdminDistance: o.AdminDistance,
		Metric:        o.Metric,
	}

	if o.NextHop == nil {
		ans.Type = ""
	} else if o.NextHop.Discard != nil {
		ans.Type = NextHopDiscard
	} else if o.NextHop.IpAddress != nil {
		ans.Type = NextHopIpAddress
		ans.NextHop = *o.NextHop.IpAddress
	} else if o.NextHop.NextVr != nil {
		ans.Type = NextHopNextVr
		ans.NextHop = *o.NextHop.NextVr
	}

	if o.Option != nil && o.Option.NoInstall != nil {
		ans.RouteTable = RouteTableNoInstall
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.Profile
	}

	return ans
}

type entry_v2 struct {
	XMLName       xml.Name     `xml:"entry"`
	Name          string       `xml:"name,attr"`
	Destination   string       `xml:"destination"`
	Interface     string       `xml:"interface,omitempty"`
	NextHop       *nextHop     `xml:"nexthop"`
	AdminDistance int          `xml:"admin-dist,omitempty"`
	Metric        int          `xml:"metric,omitempty"`
	Option        *rtOption_v1 `xml:"option"`
	Bfd           *bfd         `xml:"bfd"`
}

type bfd struct {
	Profile string `xml:"profile"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:          e.Name,
		Destination:   e.Destination,
		Interface:     e.Interface,
		AdminDistance: e.AdminDistance,
		Metric:        e.Metric,
	}

	switch e.Type {
	case NextHopDiscard:
		var sp string
		ans.NextHop = &nextHop{Discard: &sp}
	case NextHopIpAddress:
		sp := e.NextHop
		ans.NextHop = &nextHop{IpAddress: &sp}
	case NextHopNextVr:
		sp := e.NextHop
		ans.NextHop = &nextHop{NextVr: &sp}
	}

	if e.RouteTable == RouteTableNoInstall {
		sp := ""
		ans.Option = &rtOption_v1{NoInstall: &sp}
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{Profile: e.BfdProfile}
	}

	return ans
}

// PAN-OS 8.0, new routing table options
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
		Name:          o.Name,
		Destination:   o.Destination,
		Interface:     o.Interface,
		AdminDistance: o.AdminDistance,
		Metric:        o.Metric,
	}

	if o.NextHop == nil {
		ans.Type = ""
	} else if o.NextHop.Discard != nil {
		ans.Type = NextHopDiscard
	} else if o.NextHop.IpAddress != nil {
		ans.Type = NextHopIpAddress
		ans.NextHop = *o.NextHop.IpAddress
	} else if o.NextHop.NextVr != nil {
		ans.Type = NextHopNextVr
		ans.NextHop = *o.NextHop.NextVr
	}

	if o.Option != nil {
		if o.Option.Unicast != nil {
			ans.RouteTable = RouteTableUnicast
		} else if o.Option.Multicast != nil {
			ans.RouteTable = RouteTableMulticast
		} else if o.Option.Both != nil {
			ans.RouteTable = RouteTableBoth
		} else if o.Option.NoInstall != nil {
			ans.RouteTable = RouteTableNoInstall
		}
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.Profile
	}

	return ans
}

type entry_v3 struct {
	XMLName       xml.Name     `xml:"entry"`
	Name          string       `xml:"name,attr"`
	Destination   string       `xml:"destination"`
	Interface     string       `xml:"interface,omitempty"`
	NextHop       *nextHop     `xml:"nexthop"`
	AdminDistance int          `xml:"admin-dist,omitempty"`
	Metric        int          `xml:"metric,omitempty"`
	Option        *rtOption_v2 `xml:"route-table"`
	Bfd           *bfd         `xml:"bfd"`
}

type rtOption_v2 struct {
	Unicast   *string `xml:"unicast"`
	Multicast *string `xml:"multicast"`
	Both      *string `xml:"both"`
	NoInstall *string `xml:"no-install"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:          e.Name,
		Destination:   e.Destination,
		Interface:     e.Interface,
		AdminDistance: e.AdminDistance,
		Metric:        e.Metric,
	}

	switch e.Type {
	case NextHopDiscard:
		var sp string
		ans.NextHop = &nextHop{Discard: &sp}
	case NextHopIpAddress:
		sp := e.NextHop
		ans.NextHop = &nextHop{IpAddress: &sp}
	case NextHopNextVr:
		sp := e.NextHop
		ans.NextHop = &nextHop{NextVr: &sp}
	}

	switch e.RouteTable {
	case RouteTableUnicast:
		sp := ""
		ans.Option = &rtOption_v2{Unicast: &sp}
	case RouteTableMulticast:
		sp := ""
		ans.Option = &rtOption_v2{Multicast: &sp}
	case RouteTableBoth:
		sp := ""
		ans.Option = &rtOption_v2{Both: &sp}
	case RouteTableNoInstall:
		sp := ""
		ans.Option = &rtOption_v2{NoInstall: &sp}
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{Profile: e.BfdProfile}
	}

	return ans
}
