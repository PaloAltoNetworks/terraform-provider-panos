package ipv6

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an
// IPv6 static route.
type Entry struct {
	Name                string
	Destination         string
	Interface           string
	Type                string
	NextHop             string
	AdminDistance       int
	Metric              int
	RouteTable          string
	BfdProfile          string
	EnablePathMonitor   bool                 // 8.0
	PmFailureCondition  string               // 8.0
	PmHoldTime          int                  // 8.0
	MonitorDestinations []MonitorDestination // 8.0
}

type MonitorDestination struct {
	Name          string
	Enable        bool
	SourceIp      string
	DestinationIp string
	PingInterval  int
	PingCount     int
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
	o.EnablePathMonitor = s.EnablePathMonitor
	o.PmFailureCondition = s.PmFailureCondition
	o.PmHoldTime = s.PmHoldTime
	o.MonitorDestinations = s.MonitorDestinations
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

	if o.NextHop != nil {
		switch {
		case o.NextHop.Discard != nil:
			ans.Type = NextHopDiscard
		case o.NextHop.Ipv6Address != "":
			ans.Type = NextHopIpv6Address
			ans.NextHop = o.NextHop.Ipv6Address
		case o.NextHop.NextVr != "":
			ans.Type = NextHopNextVr
			ans.NextHop = o.NextHop.NextVr
		}
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.BfdProfile
	}

	if o.Option != nil && o.Option.NoInstall != nil {
		ans.RouteTable = RouteTableNoInstall
	}

	return ans
}

type entry_v1 struct {
	XMLName       xml.Name    `xml:"entry"`
	Name          string      `xml:"name,attr"`
	Destination   string      `xml:"destination"`
	Interface     string      `xml:"interface,omitempty"`
	NextHop       *nextHop_v1 `xml:"nexthop"`
	AdminDistance int         `xml:"admin-dist,omitempty"`
	Metric        int         `xml:"metric,omitempty"`
	Bfd           *bfd        `xml:"bfd"`
	Option        *option_v1  `xml:"option"`
}

type nextHop_v1 struct {
	Discard     *string `xml:"discard"`
	Ipv6Address string  `xml:"ipv6-address,omitempty"`
	NextVr      string  `xml:"next-vr,omitempty"`
}

type bfd struct {
	BfdProfile string `xml:"profile,omitempty"`
}

type option_v1 struct {
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

	s := ""
	switch e.Type {
	case NextHopDiscard:
		ans.NextHop = &nextHop_v1{Discard: &s}
	case NextHopIpv6Address:
		ans.NextHop = &nextHop_v1{Ipv6Address: e.NextHop}
	case NextHopNextVr:
		ans.NextHop = &nextHop_v1{NextVr: e.NextHop}
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{
			BfdProfile: e.BfdProfile,
		}
	}

	if e.RouteTable == RouteTableNoInstall {
		ans.Option = &option_v1{
			NoInstall: &s,
		}
	}

	return ans
}

// 8.0
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

	if o.NextHop != nil {
		switch {
		case o.NextHop.Discard != nil:
			ans.Type = NextHopDiscard
		case o.NextHop.Ipv6Address != "":
			ans.Type = NextHopIpv6Address
			ans.NextHop = o.NextHop.Ipv6Address
		case o.NextHop.NextVr != "":
			ans.Type = NextHopNextVr
			ans.NextHop = o.NextHop.NextVr
		}
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.BfdProfile
	}

	if o.Rt != nil {
		switch {
		case o.Rt.NoInstall != nil:
			ans.RouteTable = RouteTableNoInstall
		case o.Rt.Unicast != nil:
			ans.RouteTable = RouteTableUnicast
		}
	}

	if o.Monitor != nil {
		ans.EnablePathMonitor = util.AsBool(o.Monitor.EnablePathMonitor)
		ans.PmFailureCondition = o.Monitor.PmFailureCondition
		ans.PmHoldTime = o.Monitor.PmHoldTime

		if o.Monitor.Destinations != nil {
			list := make([]MonitorDestination, 0, len(o.Monitor.Destinations.Entries))
			for _, v := range o.Monitor.Destinations.Entries {
				list = append(list, MonitorDestination{
					Name:          v.Name,
					Enable:        util.AsBool(v.Enable),
					SourceIp:      v.SourceIp,
					DestinationIp: v.DestinationIp,
					PingInterval:  v.PingInterval,
					PingCount:     v.PingCount,
				})
			}

			ans.MonitorDestinations = list
		}
	}

	return ans
}

type entry_v2 struct {
	XMLName       xml.Name    `xml:"entry"`
	Name          string      `xml:"name,attr"`
	Destination   string      `xml:"destination"`
	Interface     string      `xml:"interface,omitempty"`
	NextHop       *nextHop_v1 `xml:"nexthop"`
	AdminDistance int         `xml:"admin-dist,omitempty"`
	Metric        int         `xml:"metric,omitempty"`
	Bfd           *bfd        `xml:"bfd"`
	Rt            *rt_v1      `xml:"route-table"`
	Monitor       *pm         `xml:"path-monitor"`
}

type rt_v1 struct {
	NoInstall *string `xml:"no-install"`
	Unicast   *string `xml:"unicast"`
}

type pm struct {
	EnablePathMonitor  string   `xml:"enable"`
	PmFailureCondition string   `xml:"failure-condition,omitempty"`
	PmHoldTime         int      `xml:"hold-time"`
	Destinations       *monitor `xml:"monitor-destinations"`
}

type monitor struct {
	Entries []pmEntry `xml:"entry"`
}

type pmEntry struct {
	Name          string `xml:"name,attr"`
	Enable        string `xml:"enable"`
	SourceIp      string `xml:"source"`
	DestinationIp string `xml:"destination"`
	PingInterval  int    `xml:"interval,omitempty"`
	PingCount     int    `xml:"count,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:          e.Name,
		Destination:   e.Destination,
		Interface:     e.Interface,
		AdminDistance: e.AdminDistance,
		Metric:        e.Metric,
	}

	s := ""

	switch e.Type {
	case NextHopDiscard:
		ans.NextHop = &nextHop_v1{Discard: &s}
	case NextHopIpv6Address:
		ans.NextHop = &nextHop_v1{Ipv6Address: e.NextHop}
	case NextHopNextVr:
		ans.NextHop = &nextHop_v1{NextVr: e.NextHop}
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{
			BfdProfile: e.BfdProfile,
		}
	}

	switch e.RouteTable {
	case RouteTableNoInstall:
		ans.Rt = &rt_v1{NoInstall: &s}
	case RouteTableUnicast:
		ans.Rt = &rt_v1{Unicast: &s}
	}

	if e.EnablePathMonitor || e.PmFailureCondition != "" || e.PmHoldTime != 2 || len(e.MonitorDestinations) > 0 {
		ans.Monitor = &pm{
			EnablePathMonitor:  util.YesNo(e.EnablePathMonitor),
			PmFailureCondition: e.PmFailureCondition,
			PmHoldTime:         e.PmHoldTime,
		}

		if len(e.MonitorDestinations) > 0 {
			list := make([]pmEntry, 0, len(e.MonitorDestinations))
			for _, v := range e.MonitorDestinations {
				list = append(list, pmEntry{
					Name:          v.Name,
					Enable:        util.YesNo(v.Enable),
					SourceIp:      v.SourceIp,
					DestinationIp: v.DestinationIp,
					PingInterval:  v.PingInterval,
					PingCount:     v.PingCount,
				})
			}

			ans.Monitor.Destinations = &monitor{Entries: list}
		}
	}

	return ans
}

// 9.0
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

	if o.NextHop != nil {
		switch {
		case o.NextHop.Discard != nil:
			ans.Type = NextHopDiscard
		case o.NextHop.Ipv6Address != "":
			ans.Type = NextHopIpv6Address
			ans.NextHop = o.NextHop.Ipv6Address
		case o.NextHop.Fqdn != "":
			ans.Type = NextHopFqdn
			ans.NextHop = o.NextHop.Fqdn
		case o.NextHop.NextVr != "":
			ans.Type = NextHopNextVr
			ans.NextHop = o.NextHop.NextVr
		}
	}

	if o.Bfd != nil {
		ans.BfdProfile = o.Bfd.BfdProfile
	}

	if o.Rt != nil {
		switch {
		case o.Rt.NoInstall != nil:
			ans.RouteTable = RouteTableNoInstall
		case o.Rt.Unicast != nil:
			ans.RouteTable = RouteTableUnicast
		}
	}

	if o.Monitor != nil {
		ans.EnablePathMonitor = util.AsBool(o.Monitor.EnablePathMonitor)
		ans.PmFailureCondition = o.Monitor.PmFailureCondition
		ans.PmHoldTime = o.Monitor.PmHoldTime

		if o.Monitor.Destinations != nil {
			list := make([]MonitorDestination, 0, len(o.Monitor.Destinations.Entries))
			for _, v := range o.Monitor.Destinations.Entries {
				list = append(list, MonitorDestination{
					Name:          v.Name,
					Enable:        util.AsBool(v.Enable),
					SourceIp:      v.SourceIp,
					DestinationIp: v.DestinationIp,
					PingInterval:  v.PingInterval,
					PingCount:     v.PingCount,
				})
			}

			ans.MonitorDestinations = list
		}
	}

	return ans
}

type entry_v3 struct {
	XMLName       xml.Name    `xml:"entry"`
	Name          string      `xml:"name,attr"`
	Destination   string      `xml:"destination"`
	Interface     string      `xml:"interface,omitempty"`
	NextHop       *nextHop_v2 `xml:"nexthop"`
	AdminDistance int         `xml:"admin-dist,omitempty"`
	Metric        int         `xml:"metric,omitempty"`
	Bfd           *bfd        `xml:"bfd"`
	Rt            *rt_v1      `xml:"route-table"`
	Monitor       *pm         `xml:"path-monitor"`
}

type nextHop_v2 struct {
	Discard     *string `xml:"discard"`
	Ipv6Address string  `xml:"ipv6-address,omitempty"`
	Fqdn        string  `xml:"fqdn,omitempty"`
	NextVr      string  `xml:"next-vr,omitempty"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:          e.Name,
		Destination:   e.Destination,
		Interface:     e.Interface,
		AdminDistance: e.AdminDistance,
		Metric:        e.Metric,
	}

	s := ""

	switch e.Type {
	case NextHopDiscard:
		ans.NextHop = &nextHop_v2{Discard: &s}
	case NextHopIpv6Address:
		ans.NextHop = &nextHop_v2{Ipv6Address: e.NextHop}
	case NextHopFqdn:
		ans.NextHop = &nextHop_v2{Fqdn: e.NextHop}
	case NextHopNextVr:
		ans.NextHop = &nextHop_v2{NextVr: e.NextHop}
	}

	if e.BfdProfile != "" {
		ans.Bfd = &bfd{
			BfdProfile: e.BfdProfile,
		}
	}

	switch e.RouteTable {
	case RouteTableNoInstall:
		ans.Rt = &rt_v1{NoInstall: &s}
	case RouteTableUnicast:
		ans.Rt = &rt_v1{Unicast: &s}
	}

	if e.EnablePathMonitor || e.PmFailureCondition != "" || e.PmHoldTime != 2 || len(e.MonitorDestinations) > 0 {
		ans.Monitor = &pm{
			EnablePathMonitor:  util.YesNo(e.EnablePathMonitor),
			PmFailureCondition: e.PmFailureCondition,
			PmHoldTime:         e.PmHoldTime,
		}

		if len(e.MonitorDestinations) > 0 {
			list := make([]pmEntry, 0, len(e.MonitorDestinations))
			for _, v := range e.MonitorDestinations {
				list = append(list, pmEntry{
					Name:          v.Name,
					Enable:        util.YesNo(v.Enable),
					SourceIp:      v.SourceIp,
					DestinationIp: v.DestinationIp,
					PingInterval:  v.PingInterval,
					PingCount:     v.PingCount,
				})
			}

			ans.Monitor.Destinations = &monitor{Entries: list}
		}
	}

	return ans
}
