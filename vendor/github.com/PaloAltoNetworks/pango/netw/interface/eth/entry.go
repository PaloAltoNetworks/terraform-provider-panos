package eth

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an ethernet
// interface.
type Entry struct {
	Name                        string
	Mode                        string
	StaticIps                   []string // ordered
	EnableDhcp                  bool
	CreateDhcpDefaultRoute      bool
	DhcpDefaultRouteMetric      int
	Ipv6Enabled                 bool
	Ipv6InterfaceId             string
	ManagementProfile           string
	Mtu                         int
	AdjustTcpMss                bool
	NetflowProfile              string
	LldpEnabled                 bool
	LldpProfile                 string
	LldpHaPassivePreNegotiation bool
	LacpHaPassivePreNegotiation bool
	LinkSpeed                   string
	LinkDuplex                  string
	LinkState                   string
	AggregateGroup              string
	Comment                     string
	LacpPortPriority            int
	Ipv4MssAdjust               int    // 7.1+
	Ipv6MssAdjust               int    // 7.1+
	EnableUntaggedSubinterface  bool   // 7.1+
	DecryptForward              bool   // 8.1+
	RxPolicingRate              int    // 8.1+
	TxPolicingRate              int    // 8.1+
	DhcpSendHostnameEnable      bool   // 9.0+
	DhcpSendHostnameValue       string // 9.0+

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Mode = s.Mode
	o.StaticIps = s.StaticIps
	o.EnableDhcp = s.EnableDhcp
	o.CreateDhcpDefaultRoute = s.CreateDhcpDefaultRoute
	o.DhcpDefaultRouteMetric = s.DhcpDefaultRouteMetric
	o.Ipv6Enabled = s.Ipv6Enabled
	o.ManagementProfile = s.ManagementProfile
	o.Mtu = s.Mtu
	o.AdjustTcpMss = s.AdjustTcpMss
	o.NetflowProfile = s.NetflowProfile
	o.LldpEnabled = s.LldpEnabled
	o.LldpProfile = s.LldpProfile
	o.LldpHaPassivePreNegotiation = s.LldpHaPassivePreNegotiation
	o.LacpHaPassivePreNegotiation = s.LacpHaPassivePreNegotiation
	o.LinkSpeed = s.LinkSpeed
	o.LinkDuplex = s.LinkDuplex
	o.LinkState = s.LinkState
	o.AggregateGroup = s.AggregateGroup
	o.Comment = s.Comment
	o.LacpPortPriority = s.LacpPortPriority
	o.Ipv4MssAdjust = s.Ipv4MssAdjust
	o.Ipv6MssAdjust = s.Ipv6MssAdjust
	o.EnableUntaggedSubinterface = s.EnableUntaggedSubinterface
	o.DecryptForward = s.DecryptForward
	o.RxPolicingRate = s.RxPolicingRate
	o.TxPolicingRate = s.TxPolicingRate
	o.DhcpSendHostnameEnable = s.DhcpSendHostnameEnable
	o.DhcpSendHostnameValue = s.DhcpSendHostnameValue
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, string, interface{}) {
	var iName string
	if o.Mode != ModeHa && o.Mode != ModeAggregateGroup {
		iName = o.Name
	}
	_, fn := versioning(v)

	return o.Name, iName, fn(o)
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
		Name:       o.Name,
		LinkSpeed:  o.LinkSpeed,
		LinkDuplex: o.LinkDuplex,
		LinkState:  o.LinkState,
		Comment:    o.Comment,
	}

	if o.Lacp != nil {
		ans.LacpPortPriority = o.Lacp.LacpPortPriority
	}

	ans.raw = make(map[string]string)
	switch {
	case o.ModeL3 != nil:
		ans.Mode = ModeLayer3
		ans.ManagementProfile = o.ModeL3.ManagementProfile
		ans.Mtu = o.ModeL3.Mtu
		ans.NetflowProfile = o.ModeL3.NetflowProfile
		ans.AdjustTcpMss = util.AsBool(o.ModeL3.AdjustTcpMss)
		ans.StaticIps = util.EntToStr(o.ModeL3.StaticIps)
		if o.ModeL3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.ModeL3.Dhcp.Enable)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.ModeL3.Dhcp.CreateDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.ModeL3.Dhcp.Metric
		}

		if o.ModeL3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.ModeL3.Ipv6.Enabled)
			ans.Ipv6InterfaceId = o.ModeL3.Ipv6.Ipv6InterfaceId
			if o.ModeL3.Ipv6.Address != nil {
				ans.raw["v6adr"] = util.CleanRawXml(o.ModeL3.Ipv6.Address.Text)
			}
			if o.ModeL3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.ModeL3.Ipv6.Neighbor.Text)
			}
		}

		if o.ModeL3.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL3.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL3.Lldp.LldpProfile

			if o.ModeL3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeL3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.ModeL3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.ModeL3.Arp.Text)
		}
		if o.ModeL3.Subinterface != nil {
			ans.raw["l3subinterface"] = util.CleanRawXml(o.ModeL3.Subinterface.Text)
		}
	case o.ModeL2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.ModeL2.NetflowProfile
		if o.ModeL2.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL2.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL2.Lldp.LldpProfile
		}
		if o.ModeL2.Subinterface != nil {
			ans.raw["l2subinterface"] = util.CleanRawXml(o.ModeL2.Subinterface.Text)
		}
	case o.ModeVwire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.ModeVwire.NetflowProfile
		if o.ModeVwire.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeVwire.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeVwire.Lldp.LldpProfile
			if o.ModeVwire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Lacp != nil {
			if o.ModeVwire.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lacp.Ha.LacpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Subinterface != nil {
			ans.raw["vwsub"] = util.CleanRawXml(o.ModeVwire.Subinterface.Text)
		}
	case o.TapMode != nil:
		ans.Mode = ModeTap
	case o.HaMode != nil:
		ans.Mode = ModeHa
	case o.DecryptMirrorMode != nil:
		ans.Mode = ModeDecryptMirror
	case o.AggregateGroup != "":
		ans.Mode = ModeAggregateGroup
		ans.AggregateGroup = o.AggregateGroup
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName           xml.Name   `xml:"entry"`
	Name              string     `xml:"name,attr"`
	ModeL2            *otherMode `xml:"layer2"`
	ModeL3            *l3Mode_v1 `xml:"layer3"`
	ModeVwire         *otherMode `xml:"virtual-wire"`
	TapMode           *emptyMode `xml:"tap"`
	HaMode            *emptyMode `xml:"ha"`
	DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
	AggregateGroup    string     `xml:"aggregate-group,omitempty"`
	LinkSpeed         string     `xml:"link-speed,omitempty"`
	LinkDuplex        string     `xml:"link-duplex,omitempty"`
	LinkState         string     `xml:"link-state,omitempty"`
	Comment           string     `xml:"comment"`
	Lacp              *lacp      `xml:"lacp"`
}

type emptyMode struct{}

type otherMode struct {
	NetflowProfile string       `xml:"netflow-profile,omitempty"`
	Lldp           *lldp        `xml:"lldp"`
	Subinterface   *util.RawXml `xml:"units"`
	Lacp           *omLacp      `xml:"lacp"`
}

type lldp struct {
	LldpEnabled string  `xml:"enable"`
	LldpProfile string  `xml:"profile,omitempty"`
	Ha          *lldpHa `xml:"high-availability"`
}

type lldpHa struct {
	LldpHaPassivePreNegotiation string `xml:"passive-pre-negotiation"`
}

type omLacp struct {
	Ha *omLacpHa `xml:"high-availability"`
}

type omLacpHa struct {
	LacpHaPassivePreNegotiation string `xml:"passive-pre-negotiation"`
}

type l3Mode_v1 struct {
	Ipv6              *ipv6            `xml:"ipv6"`
	ManagementProfile string           `xml:"interface-management-profile,omitempty"`
	Mtu               int              `xml:"mtu,omitempty"`
	NetflowProfile    string           `xml:"netflow-profile,omitempty"`
	AdjustTcpMss      string           `xml:"adjust-tcp-mss"`
	StaticIps         *util.EntryType  `xml:"ip"`
	Dhcp              *dhcpSettings_v1 `xml:"dhcp-client"`
	Lldp              *lldp            `xml:"lldp"`
	Arp               *util.RawXml     `xml:"arp"`
	Subinterface      *util.RawXml     `xml:"units"`
}

type ipv6 struct {
	Enabled         string       `xml:"enabled"`
	Ipv6InterfaceId string       `xml:"interface-id,omitempty"`
	Address         *util.RawXml `xml:"address"`
	Neighbor        *util.RawXml `xml:"neighbor-discovery"`
}

type dhcpSettings_v1 struct {
	Enable             string `xml:"enable"`
	CreateDefaultRoute string `xml:"create-default-route"`
	Metric             int    `xml:"default-route-metric,omitempty"`
}

type lacp struct {
	LacpPortPriority int `xml:"omitempty"`
}

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
		Name:       o.Name,
		LinkSpeed:  o.LinkSpeed,
		LinkDuplex: o.LinkDuplex,
		LinkState:  o.LinkState,
		Comment:    o.Comment,
	}

	if o.Lacp != nil {
		ans.LacpPortPriority = o.Lacp.LacpPortPriority
	}

	ans.raw = make(map[string]string)
	switch {
	case o.ModeL3 != nil:
		ans.Mode = ModeLayer3
		ans.ManagementProfile = o.ModeL3.ManagementProfile
		ans.Mtu = o.ModeL3.Mtu
		ans.NetflowProfile = o.ModeL3.NetflowProfile
		ans.AdjustTcpMss = util.AsBool(o.ModeL3.AdjustTcpMss)
		ans.Ipv4MssAdjust = o.ModeL3.Ipv4MssAdjust
		ans.Ipv6MssAdjust = o.ModeL3.Ipv6MssAdjust
		ans.StaticIps = util.EntToStr(o.ModeL3.StaticIps)
		ans.EnableUntaggedSubinterface = util.AsBool(o.ModeL3.EnableUntaggedSubinterface)

		if o.ModeL3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.ModeL3.Dhcp.Enable)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.ModeL3.Dhcp.CreateDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.ModeL3.Dhcp.Metric
		}

		if o.ModeL3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.ModeL3.Ipv6.Enabled)
			ans.Ipv6InterfaceId = o.ModeL3.Ipv6.Ipv6InterfaceId
			if o.ModeL3.Ipv6.Address != nil {
				ans.raw["v6adr"] = util.CleanRawXml(o.ModeL3.Ipv6.Address.Text)
			}
			if o.ModeL3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.ModeL3.Ipv6.Neighbor.Text)
			}
		}

		if o.ModeL3.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL3.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL3.Lldp.LldpProfile

			if o.ModeL3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeL3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.ModeL3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.ModeL3.Arp.Text)
		}
		if o.ModeL3.Subinterface != nil {
			ans.raw["l3subinterface"] = util.CleanRawXml(o.ModeL3.Subinterface.Text)
		}
		if o.ModeL3.Pppoe != nil {
			ans.raw["pppoe"] = util.CleanRawXml(o.ModeL3.Pppoe.Text)
		}
		if o.ModeL3.Ndp != nil {
			ans.raw["ndp"] = util.CleanRawXml(o.ModeL3.Ndp.Text)
		}
	case o.ModeL2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.ModeL2.NetflowProfile
		if o.ModeL2.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL2.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL2.Lldp.LldpProfile
		}
		if o.ModeL2.Subinterface != nil {
			ans.raw["l2subinterface"] = util.CleanRawXml(o.ModeL2.Subinterface.Text)
		}
	case o.ModeVwire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.ModeVwire.NetflowProfile
		if o.ModeVwire.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeVwire.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeVwire.Lldp.LldpProfile
			if o.ModeVwire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Lacp != nil {
			if o.ModeVwire.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lacp.Ha.LacpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Subinterface != nil {
			ans.raw["vwsub"] = util.CleanRawXml(o.ModeVwire.Subinterface.Text)
		}
	case o.TapMode != nil:
		ans.Mode = ModeTap
	case o.HaMode != nil:
		ans.Mode = ModeHa
	case o.DecryptMirrorMode != nil:
		ans.Mode = ModeDecryptMirror
	case o.AggregateGroup != "":
		ans.Mode = ModeAggregateGroup
		ans.AggregateGroup = o.AggregateGroup
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

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
		Name:       o.Name,
		LinkSpeed:  o.LinkSpeed,
		LinkDuplex: o.LinkDuplex,
		LinkState:  o.LinkState,
		Comment:    o.Comment,
	}

	if o.Lacp != nil {
		ans.LacpPortPriority = o.Lacp.LacpPortPriority
	}

	ans.raw = make(map[string]string)
	switch {
	case o.ModeL3 != nil:
		ans.Mode = ModeLayer3
		ans.ManagementProfile = o.ModeL3.ManagementProfile
		ans.Mtu = o.ModeL3.Mtu
		ans.NetflowProfile = o.ModeL3.NetflowProfile
		ans.AdjustTcpMss = util.AsBool(o.ModeL3.AdjustTcpMss)
		ans.Ipv4MssAdjust = o.ModeL3.Ipv4MssAdjust
		ans.Ipv6MssAdjust = o.ModeL3.Ipv6MssAdjust
		ans.StaticIps = util.EntToStr(o.ModeL3.StaticIps)
		ans.EnableUntaggedSubinterface = util.AsBool(o.ModeL3.EnableUntaggedSubinterface)
		ans.DecryptForward = util.AsBool(o.ModeL3.DecryptForward)

		if o.ModeL3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.ModeL3.Dhcp.Enable)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.ModeL3.Dhcp.CreateDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.ModeL3.Dhcp.Metric
		}

		if o.ModeL3.Policing != nil {
			ans.RxPolicingRate = o.ModeL3.Policing.RxPolicingRate
			ans.TxPolicingRate = o.ModeL3.Policing.TxPolicingRate
		}

		if o.ModeL3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.ModeL3.Ipv6.Enabled)
			ans.Ipv6InterfaceId = o.ModeL3.Ipv6.Ipv6InterfaceId
			if o.ModeL3.Ipv6.Address != nil {
				ans.raw["v6adr"] = util.CleanRawXml(o.ModeL3.Ipv6.Address.Text)
			}
			if o.ModeL3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.ModeL3.Ipv6.Neighbor.Text)
			}
		}

		if o.ModeL3.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL3.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL3.Lldp.LldpProfile

			if o.ModeL3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeL3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.ModeL3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.ModeL3.Arp.Text)
		}
		if o.ModeL3.Subinterface != nil {
			ans.raw["l3subinterface"] = util.CleanRawXml(o.ModeL3.Subinterface.Text)
		}
		if o.ModeL3.Pppoe != nil {
			ans.raw["pppoe"] = util.CleanRawXml(o.ModeL3.Pppoe.Text)
		}
		if o.ModeL3.Ndp != nil {
			ans.raw["ndp"] = util.CleanRawXml(o.ModeL3.Ndp.Text)
		}
	case o.ModeL2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.ModeL2.NetflowProfile
		if o.ModeL2.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL2.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL2.Lldp.LldpProfile
		}
		if o.ModeL2.Subinterface != nil {
			ans.raw["l2subinterface"] = util.CleanRawXml(o.ModeL2.Subinterface.Text)
		}
	case o.ModeVwire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.ModeVwire.NetflowProfile
		if o.ModeVwire.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeVwire.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeVwire.Lldp.LldpProfile
			if o.ModeVwire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Lacp != nil {
			if o.ModeVwire.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lacp.Ha.LacpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Subinterface != nil {
			ans.raw["vwsub"] = util.CleanRawXml(o.ModeVwire.Subinterface.Text)
		}
	case o.TapMode != nil:
		ans.Mode = ModeTap
	case o.HaMode != nil:
		ans.Mode = ModeHa
	case o.DecryptMirrorMode != nil:
		ans.Mode = ModeDecryptMirror
	case o.AggregateGroup != "":
		ans.Mode = ModeAggregateGroup
		ans.AggregateGroup = o.AggregateGroup
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type container_v4 struct {
	Answer []entry_v4 `xml:"entry"`
}

func (o *container_v4) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v4) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v4) normalize() Entry {
	ans := Entry{
		Name:       o.Name,
		LinkSpeed:  o.LinkSpeed,
		LinkDuplex: o.LinkDuplex,
		LinkState:  o.LinkState,
		Comment:    o.Comment,
	}

	if o.Lacp != nil {
		ans.LacpPortPriority = o.Lacp.LacpPortPriority
	}

	ans.raw = make(map[string]string)
	switch {
	case o.ModeL3 != nil:
		ans.Mode = ModeLayer3
		ans.ManagementProfile = o.ModeL3.ManagementProfile
		ans.Mtu = o.ModeL3.Mtu
		ans.NetflowProfile = o.ModeL3.NetflowProfile
		ans.AdjustTcpMss = util.AsBool(o.ModeL3.AdjustTcpMss)
		ans.Ipv4MssAdjust = o.ModeL3.Ipv4MssAdjust
		ans.Ipv6MssAdjust = o.ModeL3.Ipv6MssAdjust
		ans.StaticIps = util.EntToStr(o.ModeL3.StaticIps)
		ans.EnableUntaggedSubinterface = util.AsBool(o.ModeL3.EnableUntaggedSubinterface)
		ans.DecryptForward = util.AsBool(o.ModeL3.DecryptForward)

		if o.ModeL3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.ModeL3.Dhcp.Enable)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.ModeL3.Dhcp.CreateDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.ModeL3.Dhcp.Metric
			if o.ModeL3.Dhcp.Hostname != nil {
				ans.DhcpSendHostnameEnable = util.AsBool(o.ModeL3.Dhcp.Hostname.DhcpSendHostnameEnable)
				ans.DhcpSendHostnameValue = o.ModeL3.Dhcp.Hostname.DhcpSendHostnameValue
			}
		}

		if o.ModeL3.Policing != nil {
			ans.RxPolicingRate = o.ModeL3.Policing.RxPolicingRate
			ans.TxPolicingRate = o.ModeL3.Policing.TxPolicingRate
		}

		if o.ModeL3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.ModeL3.Ipv6.Enabled)
			ans.Ipv6InterfaceId = o.ModeL3.Ipv6.Ipv6InterfaceId
			if o.ModeL3.Ipv6.Address != nil {
				ans.raw["v6adr"] = util.CleanRawXml(o.ModeL3.Ipv6.Address.Text)
			}
			if o.ModeL3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.ModeL3.Ipv6.Neighbor.Text)
			}
		}

		if o.ModeL3.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL3.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL3.Lldp.LldpProfile

			if o.ModeL3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeL3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.ModeL3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.ModeL3.Arp.Text)
		}
		if o.ModeL3.Subinterface != nil {
			ans.raw["l3subinterface"] = util.CleanRawXml(o.ModeL3.Subinterface.Text)
		}
		if o.ModeL3.Pppoe != nil {
			ans.raw["pppoe"] = util.CleanRawXml(o.ModeL3.Pppoe.Text)
		}
		if o.ModeL3.Ndp != nil {
			ans.raw["ndp"] = util.CleanRawXml(o.ModeL3.Ndp.Text)
		}
		if o.ModeL3.Ipv6Client != nil {
			ans.raw["v6client"] = util.CleanRawXml(o.ModeL3.Ipv6Client.Text)
		}
		if o.ModeL3.Ddns != nil {
			ans.raw["ddns"] = util.CleanRawXml(o.ModeL3.Ddns.Text)
		}
	case o.ModeL2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.ModeL2.NetflowProfile
		if o.ModeL2.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeL2.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeL2.Lldp.LldpProfile
		}
		if o.ModeL2.Subinterface != nil {
			ans.raw["l2subinterface"] = util.CleanRawXml(o.ModeL2.Subinterface.Text)
		}
	case o.ModeVwire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.ModeVwire.NetflowProfile
		if o.ModeVwire.Lldp != nil {
			ans.LldpEnabled = util.AsBool(o.ModeVwire.Lldp.LldpEnabled)
			ans.LldpProfile = o.ModeVwire.Lldp.LldpProfile
			if o.ModeVwire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Lacp != nil {
			if o.ModeVwire.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.ModeVwire.Lacp.Ha.LacpHaPassivePreNegotiation)
			}
		}
		if o.ModeVwire.Subinterface != nil {
			ans.raw["vwsub"] = util.CleanRawXml(o.ModeVwire.Subinterface.Text)
		}
	case o.TapMode != nil:
		ans.Mode = ModeTap
	case o.HaMode != nil:
		ans.Mode = ModeHa
	case o.DecryptMirrorMode != nil:
		ans.Mode = ModeDecryptMirror
	case o.AggregateGroup != "":
		ans.Mode = ModeAggregateGroup
		ans.AggregateGroup = o.AggregateGroup
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v2 struct {
	XMLName           xml.Name   `xml:"entry"`
	Name              string     `xml:"name,attr"`
	ModeL3            *l3Mode_v2 `xml:"layer3"`
	ModeL2            *otherMode `xml:"layer2"`
	ModeVwire         *otherMode `xml:"virtual-wire"`
	TapMode           *emptyMode `xml:"tap"`
	HaMode            *emptyMode `xml:"ha"`
	DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
	AggregateGroup    string     `xml:"aggregate-group,omitempty"`
	LinkSpeed         string     `xml:"link-speed,omitempty"`
	LinkDuplex        string     `xml:"link-duplex,omitempty"`
	LinkState         string     `xml:"link-state,omitempty"`
	Comment           string     `xml:"comment"`
	Lacp              *lacp      `xml:"lacp"`
}

type l3Mode_v2 struct {
	Ipv6                       *ipv6            `xml:"ipv6"`
	ManagementProfile          string           `xml:"interface-management-profile,omitempty"`
	Mtu                        int              `xml:"mtu,omitempty"`
	NetflowProfile             string           `xml:"netflow-profile,omitempty"`
	AdjustTcpMss               string           `xml:"adjust-tcp-mss>enable"`
	Ipv4MssAdjust              int              `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
	Ipv6MssAdjust              int              `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`
	StaticIps                  *util.EntryType  `xml:"ip"`
	Dhcp                       *dhcpSettings_v1 `xml:"dhcp-client"`
	Lldp                       *lldp            `xml:"lldp"`
	EnableUntaggedSubinterface string           `xml:"untagged-sub-interface,omitempty"`
	Arp                        *util.RawXml     `xml:"arp"`
	Pppoe                      *util.RawXml     `xml:"pppoe"`
	Ndp                        *util.RawXml     `xml:"ndp-proxy"`
	Subinterface               *util.RawXml     `xml:"units"`
}

type entry_v3 struct {
	XMLName           xml.Name   `xml:"entry"`
	Name              string     `xml:"name,attr"`
	ModeL3            *l3Mode_v3 `xml:"layer3"`
	ModeL2            *otherMode `xml:"layer2"`
	ModeVwire         *otherMode `xml:"virtual-wire"`
	TapMode           *emptyMode `xml:"tap"`
	HaMode            *emptyMode `xml:"ha"`
	DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
	AggregateGroup    string     `xml:"aggregate-group,omitempty"`
	LinkSpeed         string     `xml:"link-speed,omitempty"`
	LinkDuplex        string     `xml:"link-duplex,omitempty"`
	LinkState         string     `xml:"link-state,omitempty"`
	Comment           string     `xml:"comment"`
	Lacp              *lacp      `xml:"lacp"`
}

type l3Mode_v3 struct {
	Ipv6                       *ipv6            `xml:"ipv6"`
	ManagementProfile          string           `xml:"interface-management-profile,omitempty"`
	Mtu                        int              `xml:"mtu,omitempty"`
	NetflowProfile             string           `xml:"netflow-profile,omitempty"`
	AdjustTcpMss               string           `xml:"adjust-tcp-mss>enable"`
	Ipv4MssAdjust              int              `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
	Ipv6MssAdjust              int              `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`
	StaticIps                  *util.EntryType  `xml:"ip"`
	Dhcp                       *dhcpSettings_v1 `xml:"dhcp-client"`
	Lldp                       *lldp            `xml:"lldp"`
	EnableUntaggedSubinterface string           `xml:"untagged-sub-interface,omitempty"`
	DecryptForward             string           `xml:"decrypt-forward,omitempty"`
	Policing                   *policing        `xml:"policing"`
	Arp                        *util.RawXml     `xml:"arp"`
	Pppoe                      *util.RawXml     `xml:"pppoe"`
	Ndp                        *util.RawXml     `xml:"ndp-proxy"`
	Subinterface               *util.RawXml     `xml:"units"`
}

type policing struct {
	RxPolicingRate int `xml:"rx-rate,omitempty"`
	TxPolicingRate int `xml:"tx-rate,omitempty"`
}

type entry_v4 struct {
	XMLName           xml.Name   `xml:"entry"`
	Name              string     `xml:"name,attr"`
	ModeL3            *l3Mode_v4 `xml:"layer3"`
	ModeL2            *otherMode `xml:"layer2"`
	ModeVwire         *otherMode `xml:"virtual-wire"`
	TapMode           *emptyMode `xml:"tap"`
	HaMode            *emptyMode `xml:"ha"`
	DecryptMirrorMode *emptyMode `xml:"decrypt-mirror"`
	AggregateGroup    string     `xml:"aggregate-group,omitempty"`
	LinkSpeed         string     `xml:"link-speed,omitempty"`
	LinkDuplex        string     `xml:"link-duplex,omitempty"`
	LinkState         string     `xml:"link-state,omitempty"`
	Comment           string     `xml:"comment"`
	Lacp              *lacp      `xml:"lacp"`
}

type l3Mode_v4 struct {
	Ipv6                       *ipv6            `xml:"ipv6"`
	ManagementProfile          string           `xml:"interface-management-profile,omitempty"`
	Mtu                        int              `xml:"mtu,omitempty"`
	NetflowProfile             string           `xml:"netflow-profile,omitempty"`
	AdjustTcpMss               string           `xml:"adjust-tcp-mss>enable"`
	Ipv4MssAdjust              int              `xml:"adjust-tcp-mss>ipv4-mss-adjustment,omitempty"`
	Ipv6MssAdjust              int              `xml:"adjust-tcp-mss>ipv6-mss-adjustment,omitempty"`
	StaticIps                  *util.EntryType  `xml:"ip"`
	Dhcp                       *dhcpSettings_v2 `xml:"dhcp-client"`
	Lldp                       *lldp            `xml:"lldp"`
	EnableUntaggedSubinterface string           `xml:"untagged-sub-interface,omitempty"`
	DecryptForward             string           `xml:"decrypt-forward,omitempty"`
	Policing                   *policing        `xml:"policing"`
	Arp                        *util.RawXml     `xml:"arp"`
	Pppoe                      *util.RawXml     `xml:"pppoe"`
	Ndp                        *util.RawXml     `xml:"ndp-proxy"`
	Ipv6Client                 *util.RawXml     `xml:"ipv6-client"`
	Subinterface               *util.RawXml     `xml:"units"`
	Ddns                       *util.RawXml     `xml:"ddns-config"`
}

type dhcpSettings_v2 struct {
	Enable             string        `xml:"enable"`
	CreateDefaultRoute string        `xml:"create-default-route"`
	Metric             int           `xml:"default-route-metric,omitempty"`
	Hostname           *dhcpHostname `xml:"send-hostname"`
}

type dhcpHostname struct {
	DhcpSendHostnameEnable string `xml:"enable,omitempty"`
	DhcpSendHostnameValue  string `xml:"hostname,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:       e.Name,
		LinkSpeed:  e.LinkSpeed,
		LinkDuplex: e.LinkDuplex,
		LinkState:  e.LinkState,
		Comment:    e.Comment,
	}

	if e.LacpPortPriority > 0 {
		ans.Lacp = &lacp{
			LacpPortPriority: e.LacpPortPriority,
		}
	}

	switch e.Mode {
	case ModeLayer3:
		i := &l3Mode_v1{
			StaticIps:         util.StrToEnt(e.StaticIps),
			ManagementProfile: e.ManagementProfile,
			Mtu:               e.Mtu,
			NetflowProfile:    e.NetflowProfile,
			AdjustTcpMss:      util.YesNo(e.AdjustTcpMss),
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
			i.Dhcp = &dhcpSettings_v1{
				Enable:             util.YesNo(e.EnableDhcp),
				CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				Metric:             e.DhcpDefaultRouteMetric,
			}
		}

		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			i.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				i.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6adr := e.raw["v6adr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
			v6 := ipv6{
				Enabled:         util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}
			if v6adr != "" {
				v6.Address = &util.RawXml{v6adr}
			}
			if v6nd != "" {
				v6.Neighbor = &util.RawXml{v6nd}
			}
			i.Ipv6 = &v6
		}

		if text, present := e.raw["arp"]; present {
			i.Arp = &util.RawXml{text}
		}
		if text, present := e.raw["l3subinterface"]; present {
			i.Subinterface = &util.RawXml{text}
		}
		ans.ModeL3 = i
	case ModeLayer2:
		ans.ModeL2 = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" {
			ans.ModeL2.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}
		}
		if text := e.raw["l2subinterface"]; text != "" {
			ans.ModeL2.Subinterface = &util.RawXml{text}
		}
	case ModeVirtualWire:
		ans.ModeVwire = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.ModeVwire.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				ans.ModeVwire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}
		if e.LacpHaPassivePreNegotiation {
			ans.ModeVwire.Lacp = &omLacp{
				Ha: &omLacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				},
			}
		}
		if text := e.raw["vwsub"]; text != "" {
			ans.ModeVwire.Subinterface = &util.RawXml{text}
		}
	case ModeTap:
		ans.TapMode = &emptyMode{}
	case ModeHa:
		ans.HaMode = &emptyMode{}
	case ModeDecryptMirror:
		ans.DecryptMirrorMode = &emptyMode{}
	case ModeAggregateGroup:
		ans.AggregateGroup = e.AggregateGroup
	}

	return ans
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:       e.Name,
		LinkSpeed:  e.LinkSpeed,
		LinkDuplex: e.LinkDuplex,
		LinkState:  e.LinkState,
		Comment:    e.Comment,
	}

	if e.LacpPortPriority > 0 {
		ans.Lacp = &lacp{
			LacpPortPriority: e.LacpPortPriority,
		}
	}

	switch e.Mode {
	case ModeLayer3:
		i := &l3Mode_v2{
			StaticIps:         util.StrToEnt(e.StaticIps),
			ManagementProfile: e.ManagementProfile,
			Mtu:               e.Mtu,
			NetflowProfile:    e.NetflowProfile,
			AdjustTcpMss:      util.YesNo(e.AdjustTcpMss),
			Ipv4MssAdjust:     e.Ipv4MssAdjust,
			Ipv6MssAdjust:     e.Ipv6MssAdjust,
		}

		if e.EnableUntaggedSubinterface {
			i.EnableUntaggedSubinterface = util.YesNo(e.EnableUntaggedSubinterface)
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
			i.Dhcp = &dhcpSettings_v1{
				Enable:             util.YesNo(e.EnableDhcp),
				CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				Metric:             e.DhcpDefaultRouteMetric,
			}
		}

		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			i.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				i.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6adr := e.raw["v6adr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
			v6 := ipv6{
				Enabled:         util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}
			if v6adr != "" {
				v6.Address = &util.RawXml{v6adr}
			}
			if v6nd != "" {
				v6.Neighbor = &util.RawXml{v6nd}
			}
			i.Ipv6 = &v6
		}

		if text, present := e.raw["arp"]; present {
			i.Arp = &util.RawXml{text}
		}
		if text, present := e.raw["l3subinterface"]; present {
			i.Subinterface = &util.RawXml{text}
		}
		if text := e.raw["pppoe"]; text != "" {
			i.Pppoe = &util.RawXml{text}
		}
		if text := e.raw["ndp"]; text != "" {
			i.Ndp = &util.RawXml{text}
		}
		ans.ModeL3 = i
	case ModeLayer2:
		ans.ModeL2 = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" {
			ans.ModeL2.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}
		}
		if text := e.raw["l2subinterface"]; text != "" {
			ans.ModeL2.Subinterface = &util.RawXml{text}
		}
	case ModeVirtualWire:
		ans.ModeVwire = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.ModeVwire.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				ans.ModeVwire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}
		if e.LacpHaPassivePreNegotiation {
			ans.ModeVwire.Lacp = &omLacp{
				Ha: &omLacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				},
			}
		}
		if text := e.raw["vwsub"]; text != "" {
			ans.ModeVwire.Subinterface = &util.RawXml{text}
		}
	case ModeTap:
		ans.TapMode = &emptyMode{}
	case ModeHa:
		ans.HaMode = &emptyMode{}
	case ModeDecryptMirror:
		ans.DecryptMirrorMode = &emptyMode{}
	case ModeAggregateGroup:
		ans.AggregateGroup = e.AggregateGroup
	}

	return ans
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:       e.Name,
		LinkSpeed:  e.LinkSpeed,
		LinkDuplex: e.LinkDuplex,
		LinkState:  e.LinkState,
		Comment:    e.Comment,
	}

	if e.LacpPortPriority > 0 {
		ans.Lacp = &lacp{
			LacpPortPriority: e.LacpPortPriority,
		}
	}

	switch e.Mode {
	case ModeLayer3:
		i := &l3Mode_v3{
			StaticIps:         util.StrToEnt(e.StaticIps),
			ManagementProfile: e.ManagementProfile,
			Mtu:               e.Mtu,
			NetflowProfile:    e.NetflowProfile,
			AdjustTcpMss:      util.YesNo(e.AdjustTcpMss),
			Ipv4MssAdjust:     e.Ipv4MssAdjust,
			Ipv6MssAdjust:     e.Ipv6MssAdjust,
		}

		if e.EnableUntaggedSubinterface {
			i.EnableUntaggedSubinterface = util.YesNo(e.EnableUntaggedSubinterface)
		}

		if e.DecryptForward {
			i.DecryptForward = util.YesNo(e.DecryptForward)
		}

		if e.RxPolicingRate != 0 || e.TxPolicingRate != 0 {
			i.Policing = &policing{
				RxPolicingRate: e.RxPolicingRate,
				TxPolicingRate: e.TxPolicingRate,
			}
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
			i.Dhcp = &dhcpSettings_v1{
				Enable:             util.YesNo(e.EnableDhcp),
				CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				Metric:             e.DhcpDefaultRouteMetric,
			}
		}

		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			i.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				i.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6adr := e.raw["v6adr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
			v6 := ipv6{
				Enabled:         util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}
			if v6adr != "" {
				v6.Address = &util.RawXml{v6adr}
			}
			if v6nd != "" {
				v6.Neighbor = &util.RawXml{v6nd}
			}
			i.Ipv6 = &v6
		}

		if text, present := e.raw["arp"]; present {
			i.Arp = &util.RawXml{text}
		}
		if text, present := e.raw["l3subinterface"]; present {
			i.Subinterface = &util.RawXml{text}
		}
		if text := e.raw["pppoe"]; text != "" {
			i.Pppoe = &util.RawXml{text}
		}
		if text := e.raw["ndp"]; text != "" {
			i.Ndp = &util.RawXml{text}
		}
		ans.ModeL3 = i
	case ModeLayer2:
		ans.ModeL2 = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" {
			ans.ModeL2.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}
		}
		if text := e.raw["l2subinterface"]; text != "" {
			ans.ModeL2.Subinterface = &util.RawXml{text}
		}
	case ModeVirtualWire:
		ans.ModeVwire = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.ModeVwire.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				ans.ModeVwire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}
		if e.LacpHaPassivePreNegotiation {
			ans.ModeVwire.Lacp = &omLacp{
				Ha: &omLacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				},
			}
		}
		if text := e.raw["vwsub"]; text != "" {
			ans.ModeVwire.Subinterface = &util.RawXml{text}
		}
	case ModeTap:
		ans.TapMode = &emptyMode{}
	case ModeHa:
		ans.HaMode = &emptyMode{}
	case ModeDecryptMirror:
		ans.DecryptMirrorMode = &emptyMode{}
	case ModeAggregateGroup:
		ans.AggregateGroup = e.AggregateGroup
	}

	return ans
}

func specify_v4(e Entry) interface{} {
	ans := entry_v4{
		Name:       e.Name,
		LinkSpeed:  e.LinkSpeed,
		LinkDuplex: e.LinkDuplex,
		LinkState:  e.LinkState,
		Comment:    e.Comment,
	}

	if e.LacpPortPriority > 0 {
		ans.Lacp = &lacp{
			LacpPortPriority: e.LacpPortPriority,
		}
	}

	switch e.Mode {
	case ModeLayer3:
		i := &l3Mode_v4{
			StaticIps:         util.StrToEnt(e.StaticIps),
			ManagementProfile: e.ManagementProfile,
			Mtu:               e.Mtu,
			NetflowProfile:    e.NetflowProfile,
			AdjustTcpMss:      util.YesNo(e.AdjustTcpMss),
			Ipv4MssAdjust:     e.Ipv4MssAdjust,
			Ipv6MssAdjust:     e.Ipv6MssAdjust,
		}

		if e.EnableUntaggedSubinterface {
			i.EnableUntaggedSubinterface = util.YesNo(e.EnableUntaggedSubinterface)
		}

		if e.DecryptForward {
			i.DecryptForward = util.YesNo(e.DecryptForward)
		}

		if e.RxPolicingRate != 0 || e.TxPolicingRate != 0 {
			i.Policing = &policing{
				RxPolicingRate: e.RxPolicingRate,
				TxPolicingRate: e.TxPolicingRate,
			}
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 || e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != "" {
			i.Dhcp = &dhcpSettings_v2{
				Enable:             util.YesNo(e.EnableDhcp),
				CreateDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				Metric:             e.DhcpDefaultRouteMetric,
			}

			if e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != "" {
				i.Dhcp.Hostname = &dhcpHostname{
					DhcpSendHostnameEnable: util.YesNo(e.DhcpSendHostnameEnable),
					DhcpSendHostnameValue:  e.DhcpSendHostnameValue,
				}
			}
		}

		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			i.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				i.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6adr := e.raw["v6adr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nd != "" {
			v6 := ipv6{
				Enabled:         util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}
			if v6adr != "" {
				v6.Address = &util.RawXml{v6adr}
			}
			if v6nd != "" {
				v6.Neighbor = &util.RawXml{v6nd}
			}
			i.Ipv6 = &v6
		}

		if text, present := e.raw["arp"]; present {
			i.Arp = &util.RawXml{text}
		}
		if text, present := e.raw["l3subinterface"]; present {
			i.Subinterface = &util.RawXml{text}
		}
		if text := e.raw["pppoe"]; text != "" {
			i.Pppoe = &util.RawXml{text}
		}
		if text := e.raw["ndp"]; text != "" {
			i.Ndp = &util.RawXml{text}
		}
		if text := e.raw["v6client"]; text != "" {
			i.Ipv6Client = &util.RawXml{text}
		}
		if text := e.raw["ddns"]; text != "" {
			i.Ddns = &util.RawXml{text}
		}
		ans.ModeL3 = i
	case ModeLayer2:
		ans.ModeL2 = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" {
			ans.ModeL2.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}
		}
		if text := e.raw["l2subinterface"]; text != "" {
			ans.ModeL2.Subinterface = &util.RawXml{text}
		}
	case ModeVirtualWire:
		ans.ModeVwire = &otherMode{
			NetflowProfile: e.NetflowProfile,
		}
		if e.LldpEnabled || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.ModeVwire.Lldp = &lldp{
				LldpEnabled: util.YesNo(e.LldpEnabled),
				LldpProfile: e.LldpProfile,
			}

			if e.LldpHaPassivePreNegotiation {
				ans.ModeVwire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}
		if e.LacpHaPassivePreNegotiation {
			ans.ModeVwire.Lacp = &omLacp{
				Ha: &omLacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				},
			}
		}
		if text := e.raw["vwsub"]; text != "" {
			ans.ModeVwire.Subinterface = &util.RawXml{text}
		}
	case ModeTap:
		ans.TapMode = &emptyMode{}
	case ModeHa:
		ans.HaMode = &emptyMode{}
	case ModeDecryptMirror:
		ans.DecryptMirrorMode = &emptyMode{}
	case ModeAggregateGroup:
		ans.AggregateGroup = e.AggregateGroup
	}

	return ans
}
