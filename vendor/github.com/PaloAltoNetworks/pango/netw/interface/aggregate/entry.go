package aggregate

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an aggregate
// ethernet interface.
type Entry struct {
	Name                        string
	Mode                        string
	NetflowProfile              string
	Mtu                         int
	AdjustTcpMss                bool
	Ipv4MssAdjust               int
	Ipv6MssAdjust               int
	EnableUntaggedSubinterface  bool
	StaticIps                   []string // ordered
	Ipv6Enabled                 bool
	Ipv6InterfaceId             string
	ManagementProfile           string
	EnableDhcp                  bool
	CreateDhcpDefaultRoute      bool
	DhcpDefaultRouteMetric      int
	LacpEnable                  bool   // ha, l2, l3
	LacpFastFailover            bool   // ha, l2, l3
	LacpMode                    string // ha, l2, l3
	LacpTransmissionRate        string // ha, l2, l3
	LacpSystemPriority          int    // ha, l2, l3
	LacpMaxPorts                int    // ha, l2, l3
	LacpHaPassivePreNegotiation bool   // l2, l3
	LacpHaEnableSameSystemMac   bool   // l2, l3
	LacpHaSameSystemMacAddress  string // l2, l3
	LldpEnable                  bool
	LldpProfile                 string
	LldpHaPassivePreNegotiation bool
	Comment                     string
	DecryptForward              bool   // 8.1+
	DhcpSendHostnameEnable      bool   // 9.0+
	DhcpSendHostnameValue       string // 9.0+

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Mode = s.Mode
	o.NetflowProfile = s.NetflowProfile
	o.Mtu = s.Mtu
	o.AdjustTcpMss = s.AdjustTcpMss
	o.Ipv4MssAdjust = s.Ipv4MssAdjust
	o.Ipv6MssAdjust = s.Ipv6MssAdjust
	o.EnableUntaggedSubinterface = s.EnableUntaggedSubinterface
	o.StaticIps = s.StaticIps
	o.Ipv6Enabled = s.Ipv6Enabled
	o.Ipv6InterfaceId = s.Ipv6InterfaceId
	o.ManagementProfile = s.ManagementProfile
	o.EnableDhcp = s.EnableDhcp
	o.CreateDhcpDefaultRoute = s.CreateDhcpDefaultRoute
	o.DhcpDefaultRouteMetric = s.DhcpDefaultRouteMetric
	o.Comment = s.Comment
	o.LacpEnable = s.LacpEnable
	o.LacpFastFailover = s.LacpFastFailover
	o.LacpMode = s.LacpMode
	o.LacpTransmissionRate = s.LacpTransmissionRate
	o.LacpSystemPriority = s.LacpSystemPriority
	o.LacpMaxPorts = s.LacpMaxPorts
	o.LacpHaPassivePreNegotiation = s.LacpHaPassivePreNegotiation
	o.LacpHaEnableSameSystemMac = s.LacpHaEnableSameSystemMac
	o.LacpHaSameSystemMacAddress = s.LacpHaSameSystemMacAddress
	o.LldpEnable = s.LldpEnable
	o.LldpProfile = s.LldpProfile
	o.LldpHaPassivePreNegotiation = s.LldpHaPassivePreNegotiation
	o.DecryptForward = s.DecryptForward
	o.DhcpSendHostnameEnable = s.DhcpSendHostnameEnable
	o.DhcpSendHostnameValue = s.DhcpSendHostnameValue
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, string, interface{}) {
	var iName string
	if o.Mode != ModeHa {
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
		Name:    o.Name,
		Comment: o.Comment,
	}

	ans.raw = make(map[string]string)
	switch {
	case o.Ha != nil:
		ans.Mode = ModeHa
		if o.Ha.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.Ha.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.Ha.Lacp.LacpFastFailover)
			ans.LacpMode = o.Ha.Lacp.LacpMode
			ans.LacpTransmissionRate = o.Ha.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.Ha.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.Ha.Lacp.LacpMaxPorts
		}
	case o.DecryptMirror != nil:
		ans.Mode = ModeDecryptMirror
	case o.VirtualWire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.VirtualWire.NetflowProfile

		if o.VirtualWire.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.VirtualWire.Lldp.LldpEnable)
			ans.LldpProfile = o.VirtualWire.Lldp.LldpProfile
			if o.VirtualWire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.VirtualWire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.VirtualWire.Subinterfaces != nil {
			ans.raw["vwsi"] = util.CleanRawXml(o.VirtualWire.Subinterfaces.Text)
		}
	case o.L2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.L2.NetflowProfile

		if o.L2.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.L2.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.L2.Lacp.LacpFastFailover)
			ans.LacpMode = o.L2.Lacp.LacpMode
			ans.LacpTransmissionRate = o.L2.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.L2.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.L2.Lacp.LacpMaxPorts

			if o.L2.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.L2.Lacp.Ha.LacpHaPassivePreNegotiation)
				if o.L2.Lacp.Ha.Mac != nil {
					ans.LacpHaEnableSameSystemMac = util.AsBool(o.L2.Lacp.Ha.Mac.LacpHaEnableSameSystemMac)
					ans.LacpHaSameSystemMacAddress = o.L2.Lacp.Ha.Mac.LacpHaSameSystemMacAddress
				}
			}
		}

		if o.L2.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.L2.Lldp.LldpEnable)
			ans.LldpProfile = o.L2.Lldp.LldpProfile
			if o.L2.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.L2.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.L2.Subinterfaces != nil {
			ans.raw["l2si"] = util.CleanRawXml(o.L2.Subinterfaces.Text)
		}
	case o.L3 != nil:
		ans.Mode = ModeLayer3
		ans.Mtu = o.L3.Mtu
		ans.EnableUntaggedSubinterface = util.AsBool(o.L3.EnableUntaggedSubinterface)
		ans.StaticIps = util.EntToStr(o.L3.StaticIps)
		ans.ManagementProfile = o.L3.ManagementProfile
		ans.NetflowProfile = o.L3.NetflowProfile

		if o.L3.Mss != nil {
			ans.AdjustTcpMss = util.AsBool(o.L3.Mss.AdjustTcpMss)
			ans.Ipv4MssAdjust = o.L3.Mss.Ipv4MssAdjust
			ans.Ipv6MssAdjust = o.L3.Mss.Ipv6MssAdjust
		}

		if o.L3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.L3.Ipv6.Ipv6Enabled)
			ans.Ipv6InterfaceId = o.L3.Ipv6.Ipv6InterfaceId

			if o.L3.Ipv6.Address != nil {
				ans.raw["v6addr"] = util.CleanRawXml(o.L3.Ipv6.Address.Text)
			}
			if o.L3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.L3.Ipv6.Neighbor.Text)
			}
		}

		if o.L3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.L3.Dhcp.EnableDhcp)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.L3.Dhcp.CreateDhcpDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.L3.Dhcp.DhcpDefaultRouteMetric
		}

		if o.L3.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.L3.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.L3.Lacp.LacpFastFailover)
			ans.LacpMode = o.L3.Lacp.LacpMode
			ans.LacpTransmissionRate = o.L3.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.L3.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.L3.Lacp.LacpMaxPorts

			if o.L3.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.L3.Lacp.Ha.LacpHaPassivePreNegotiation)
				if o.L3.Lacp.Ha.Mac != nil {
					ans.LacpHaEnableSameSystemMac = util.AsBool(o.L3.Lacp.Ha.Mac.LacpHaEnableSameSystemMac)
					ans.LacpHaSameSystemMacAddress = o.L3.Lacp.Ha.Mac.LacpHaSameSystemMacAddress
				}
			}
		}

		if o.L3.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.L3.Lldp.LldpEnable)
			ans.LldpProfile = o.L3.Lldp.LldpProfile
			if o.L3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.L3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.L3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.L3.Arp.Text)
		}
		if o.L3.Ndp != nil {
			ans.raw["ndp"] = util.CleanRawXml(o.L3.Ndp.Text)
		}
		if o.L3.Subinterfaces != nil {
			ans.raw["l3si"] = util.CleanRawXml(o.L3.Subinterfaces.Text)
		}
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
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
		Name:    o.Name,
		Comment: o.Comment,
	}

	ans.raw = make(map[string]string)
	switch {
	case o.Ha != nil:
		ans.Mode = ModeHa
		if o.Ha.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.Ha.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.Ha.Lacp.LacpFastFailover)
			ans.LacpMode = o.Ha.Lacp.LacpMode
			ans.LacpTransmissionRate = o.Ha.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.Ha.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.Ha.Lacp.LacpMaxPorts
		}
	case o.DecryptMirror != nil:
		ans.Mode = ModeDecryptMirror
	case o.VirtualWire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.VirtualWire.NetflowProfile

		if o.VirtualWire.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.VirtualWire.Lldp.LldpEnable)
			ans.LldpProfile = o.VirtualWire.Lldp.LldpProfile
			if o.VirtualWire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.VirtualWire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.VirtualWire.Subinterfaces != nil {
			ans.raw["vwsi"] = util.CleanRawXml(o.VirtualWire.Subinterfaces.Text)
		}
	case o.L2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.L2.NetflowProfile

		if o.L2.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.L2.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.L2.Lacp.LacpFastFailover)
			ans.LacpMode = o.L2.Lacp.LacpMode
			ans.LacpTransmissionRate = o.L2.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.L2.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.L2.Lacp.LacpMaxPorts

			if o.L2.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.L2.Lacp.Ha.LacpHaPassivePreNegotiation)
				if o.L2.Lacp.Ha.Mac != nil {
					ans.LacpHaEnableSameSystemMac = util.AsBool(o.L2.Lacp.Ha.Mac.LacpHaEnableSameSystemMac)
					ans.LacpHaSameSystemMacAddress = o.L2.Lacp.Ha.Mac.LacpHaSameSystemMacAddress
				}
			}
		}

		if o.L2.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.L2.Lldp.LldpEnable)
			ans.LldpProfile = o.L2.Lldp.LldpProfile
			if o.L2.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.L2.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.L2.Subinterfaces != nil {
			ans.raw["l2si"] = util.CleanRawXml(o.L2.Subinterfaces.Text)
		}
	case o.L3 != nil:
		ans.Mode = ModeLayer3
		ans.Mtu = o.L3.Mtu
		ans.EnableUntaggedSubinterface = util.AsBool(o.L3.EnableUntaggedSubinterface)
		ans.StaticIps = util.EntToStr(o.L3.StaticIps)
		ans.ManagementProfile = o.L3.ManagementProfile
		ans.NetflowProfile = o.L3.NetflowProfile
		ans.DecryptForward = util.AsBool(o.L3.DecryptForward)

		if o.L3.Mss != nil {
			ans.AdjustTcpMss = util.AsBool(o.L3.Mss.AdjustTcpMss)
			ans.Ipv4MssAdjust = o.L3.Mss.Ipv4MssAdjust
			ans.Ipv6MssAdjust = o.L3.Mss.Ipv6MssAdjust
		}

		if o.L3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.L3.Ipv6.Ipv6Enabled)
			ans.Ipv6InterfaceId = o.L3.Ipv6.Ipv6InterfaceId

			if o.L3.Ipv6.Address != nil {
				ans.raw["v6addr"] = util.CleanRawXml(o.L3.Ipv6.Address.Text)
			}
			if o.L3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.L3.Ipv6.Neighbor.Text)
			}
		}

		if o.L3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.L3.Dhcp.EnableDhcp)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.L3.Dhcp.CreateDhcpDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.L3.Dhcp.DhcpDefaultRouteMetric
		}

		if o.L3.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.L3.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.L3.Lacp.LacpFastFailover)
			ans.LacpMode = o.L3.Lacp.LacpMode
			ans.LacpTransmissionRate = o.L3.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.L3.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.L3.Lacp.LacpMaxPorts

			if o.L3.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.L3.Lacp.Ha.LacpHaPassivePreNegotiation)
				if o.L3.Lacp.Ha.Mac != nil {
					ans.LacpHaEnableSameSystemMac = util.AsBool(o.L3.Lacp.Ha.Mac.LacpHaEnableSameSystemMac)
					ans.LacpHaSameSystemMacAddress = o.L3.Lacp.Ha.Mac.LacpHaSameSystemMacAddress
				}
			}
		}

		if o.L3.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.L3.Lldp.LldpEnable)
			ans.LldpProfile = o.L3.Lldp.LldpProfile
			if o.L3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.L3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.L3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.L3.Arp.Text)
		}
		if o.L3.Ndp != nil {
			ans.raw["ndp"] = util.CleanRawXml(o.L3.Ndp.Text)
		}
		if o.L3.Subinterfaces != nil {
			ans.raw["l3si"] = util.CleanRawXml(o.L3.Subinterfaces.Text)
		}
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
		Name:    o.Name,
		Comment: o.Comment,
	}

	ans.raw = make(map[string]string)
	switch {
	case o.Ha != nil:
		ans.Mode = ModeHa
		if o.Ha.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.Ha.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.Ha.Lacp.LacpFastFailover)
			ans.LacpMode = o.Ha.Lacp.LacpMode
			ans.LacpTransmissionRate = o.Ha.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.Ha.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.Ha.Lacp.LacpMaxPorts
		}
	case o.DecryptMirror != nil:
		ans.Mode = ModeDecryptMirror
	case o.VirtualWire != nil:
		ans.Mode = ModeVirtualWire
		ans.NetflowProfile = o.VirtualWire.NetflowProfile

		if o.VirtualWire.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.VirtualWire.Lldp.LldpEnable)
			ans.LldpProfile = o.VirtualWire.Lldp.LldpProfile
			if o.VirtualWire.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.VirtualWire.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.VirtualWire.Subinterfaces != nil {
			ans.raw["vwsi"] = util.CleanRawXml(o.VirtualWire.Subinterfaces.Text)
		}
	case o.L2 != nil:
		ans.Mode = ModeLayer2
		ans.NetflowProfile = o.L2.NetflowProfile

		if o.L2.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.L2.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.L2.Lacp.LacpFastFailover)
			ans.LacpMode = o.L2.Lacp.LacpMode
			ans.LacpTransmissionRate = o.L2.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.L2.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.L2.Lacp.LacpMaxPorts

			if o.L2.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.L2.Lacp.Ha.LacpHaPassivePreNegotiation)
				if o.L2.Lacp.Ha.Mac != nil {
					ans.LacpHaEnableSameSystemMac = util.AsBool(o.L2.Lacp.Ha.Mac.LacpHaEnableSameSystemMac)
					ans.LacpHaSameSystemMacAddress = o.L2.Lacp.Ha.Mac.LacpHaSameSystemMacAddress
				}
			}
		}

		if o.L2.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.L2.Lldp.LldpEnable)
			ans.LldpProfile = o.L2.Lldp.LldpProfile
			if o.L2.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.L2.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.L2.Subinterfaces != nil {
			ans.raw["l2si"] = util.CleanRawXml(o.L2.Subinterfaces.Text)
		}
	case o.L3 != nil:
		ans.Mode = ModeLayer3
		ans.Mtu = o.L3.Mtu
		ans.EnableUntaggedSubinterface = util.AsBool(o.L3.EnableUntaggedSubinterface)
		ans.StaticIps = util.EntToStr(o.L3.StaticIps)
		ans.ManagementProfile = o.L3.ManagementProfile
		ans.NetflowProfile = o.L3.NetflowProfile
		ans.DecryptForward = util.AsBool(o.L3.DecryptForward)

		if o.L3.Mss != nil {
			ans.AdjustTcpMss = util.AsBool(o.L3.Mss.AdjustTcpMss)
			ans.Ipv4MssAdjust = o.L3.Mss.Ipv4MssAdjust
			ans.Ipv6MssAdjust = o.L3.Mss.Ipv6MssAdjust
		}

		if o.L3.Ipv6 != nil {
			ans.Ipv6Enabled = util.AsBool(o.L3.Ipv6.Ipv6Enabled)
			ans.Ipv6InterfaceId = o.L3.Ipv6.Ipv6InterfaceId

			if o.L3.Ipv6.Address != nil {
				ans.raw["v6addr"] = util.CleanRawXml(o.L3.Ipv6.Address.Text)
			}
			if o.L3.Ipv6.Neighbor != nil {
				ans.raw["v6nd"] = util.CleanRawXml(o.L3.Ipv6.Neighbor.Text)
			}
		}

		if o.L3.Dhcp != nil {
			ans.EnableDhcp = util.AsBool(o.L3.Dhcp.EnableDhcp)
			ans.CreateDhcpDefaultRoute = util.AsBool(o.L3.Dhcp.CreateDhcpDefaultRoute)
			ans.DhcpDefaultRouteMetric = o.L3.Dhcp.DhcpDefaultRouteMetric

			if o.L3.Dhcp.Hostname != nil {
				ans.DhcpSendHostnameEnable = util.AsBool(o.L3.Dhcp.Hostname.DhcpSendHostnameEnable)
				ans.DhcpSendHostnameValue = o.L3.Dhcp.Hostname.DhcpSendHostnameValue
			}
		}

		if o.L3.Lacp != nil {
			ans.LacpEnable = util.AsBool(o.L3.Lacp.LacpEnable)
			ans.LacpFastFailover = util.AsBool(o.L3.Lacp.LacpFastFailover)
			ans.LacpMode = o.L3.Lacp.LacpMode
			ans.LacpTransmissionRate = o.L3.Lacp.LacpTransmissionRate
			ans.LacpSystemPriority = o.L3.Lacp.LacpSystemPriority
			ans.LacpMaxPorts = o.L3.Lacp.LacpMaxPorts

			if o.L3.Lacp.Ha != nil {
				ans.LacpHaPassivePreNegotiation = util.AsBool(o.L3.Lacp.Ha.LacpHaPassivePreNegotiation)
				if o.L3.Lacp.Ha.Mac != nil {
					ans.LacpHaEnableSameSystemMac = util.AsBool(o.L3.Lacp.Ha.Mac.LacpHaEnableSameSystemMac)
					ans.LacpHaSameSystemMacAddress = o.L3.Lacp.Ha.Mac.LacpHaSameSystemMacAddress
				}
			}
		}

		if o.L3.Lldp != nil {
			ans.LldpEnable = util.AsBool(o.L3.Lldp.LldpEnable)
			ans.LldpProfile = o.L3.Lldp.LldpProfile
			if o.L3.Lldp.Ha != nil {
				ans.LldpHaPassivePreNegotiation = util.AsBool(o.L3.Lldp.Ha.LldpHaPassivePreNegotiation)
			}
		}

		if o.L3.Arp != nil {
			ans.raw["arp"] = util.CleanRawXml(o.L3.Arp.Text)
		}
		if o.L3.Ndp != nil {
			ans.raw["ndp"] = util.CleanRawXml(o.L3.Ndp.Text)
		}
		if o.L3.Subinterfaces != nil {
			ans.raw["l3si"] = util.CleanRawXml(o.L3.Subinterfaces.Text)
		}
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName       xml.Name   `xml:"entry"`
	Name          string     `xml:"name,attr"`
	Ha            *ha        `xml:"ha"`
	DecryptMirror *string    `xml:"decrypt-mirror"`
	VirtualWire   *layer2    `xml:"virtual-wire"`
	L2            *layer2    `xml:"layer2"`
	L3            *layer3_v1 `xml:"layer3"`
	Comment       string     `xml:"comment,omitempty"`
}

type ha struct {
	Lacp *lacpNoHa `xml:"lacp"`
}

type lacpNoHa struct {
	LacpEnable           string `xml:"enable"`
	LacpFastFailover     string `xml:"fast-failover"`
	LacpMode             string `xml:"mode,omitempty"`
	LacpTransmissionRate string `xml:"transmission-rate,omitempty"`
	LacpSystemPriority   int    `xml:"system-priority,omitempty"`
	LacpMaxPorts         int    `xml:"max-ports,omitempty"`
}

type layer2 struct {
	NetflowProfile string       `xml:"netflow-profile,omitempty"`
	Subinterfaces  *util.RawXml `xml:"units"`
	Lacp           *lacpWithHa  `xml:"lacp"`
	Lldp           *lldp        `xml:"lldp"`
}

type lacpWithHa struct {
	LacpEnable           string  `xml:"enable"`
	LacpFastFailover     string  `xml:"fast-failover"`
	LacpMode             string  `xml:"mode,omitempty"`
	LacpTransmissionRate string  `xml:"transmission-rate,omitempty"`
	LacpSystemPriority   int     `xml:"system-priority,omitempty"`
	LacpMaxPorts         int     `xml:"max-ports,omitempty"`
	Ha                   *lacpHa `xml:"high-availability"`
}

type lacpHa struct {
	Mac                         *lacpHaMac `xml:"use-same-system-mac"`
	LacpHaPassivePreNegotiation string     `xml:"passive-pre-negotiation"`
}

type lacpHaMac struct {
	LacpHaEnableSameSystemMac  string `xml:"enable"`
	LacpHaSameSystemMacAddress string `xml:"mac-address,omitempty"`
}

type lldp struct {
	LldpEnable  string  `xml:"enable"`
	LldpProfile string  `xml:"profile,omitempty"`
	Ha          *lldpHa `xml:"high-availability"`
}

type lldpHa struct {
	LldpHaPassivePreNegotiation string `xml:"passive-pre-negotiation"`
}

type layer3_v1 struct {
	Mtu                        int              `xml:"mtu,omitempty"`
	Mss                        *mss             `xml:"adjust-tcp-mss"`
	EnableUntaggedSubinterface string           `xml:"untagged-sub-interface"`
	Lacp                       *lacpWithHa      `xml:"lacp"`
	Lldp                       *lldp            `xml:"lldp"`
	StaticIps                  *util.EntryType  `xml:"ip"`
	Ipv6                       *ipv6            `xml:"ipv6"`
	Arp                        *util.RawXml     `xml:"arp"`
	Ndp                        *util.RawXml     `xml:"ndp-proxy"`
	ManagementProfile          string           `xml:"interface-management-profile,omitempty"`
	Dhcp                       *dhcpSettings_v1 `xml:"dhcp-client"`
	Subinterfaces              *util.RawXml     `xml:"units"`
	NetflowProfile             string           `xml:"netflow-profile,omitempty"`
}

type mss struct {
	AdjustTcpMss  string `xml:"enable"`
	Ipv4MssAdjust int    `xml:"ipv4-mss-adjustment,omitempty"`
	Ipv6MssAdjust int    `xml:"ipv6-mss-adjustment,omitempty"`
}

type ipv6 struct {
	Ipv6Enabled     string       `xml:"enabled"`
	Ipv6InterfaceId string       `xml:"interface-id,omitempty"`
	Address         *util.RawXml `xml:"address"`
	Neighbor        *util.RawXml `xml:"neighbor-discovery"`
}

type dhcpSettings_v1 struct {
	EnableDhcp             string `xml:"enable"`
	CreateDhcpDefaultRoute string `xml:"create-default-route"`
	DhcpDefaultRouteMetric int    `xml:"default-route-metric,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:    e.Name,
		Comment: e.Comment,
	}

	switch e.Mode {
	case ModeHa:
		ans.Ha = &ha{}
		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 {
			ans.Ha.Lacp = &lacpNoHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}
		}
	case ModeDecryptMirror:
		s := ""
		ans.DecryptMirror = &s
	case ModeVirtualWire:
		ans.VirtualWire = &layer2{
			NetflowProfile: e.NetflowProfile,
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.VirtualWire.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.VirtualWire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		if text := e.raw["vwsi"]; text != "" {
			ans.VirtualWire.Subinterfaces = &util.RawXml{text}
		}
	case ModeLayer2:
		ans.L2 = &layer2{
			NetflowProfile: e.NetflowProfile,
		}

		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 || e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
			ans.L2.Lacp = &lacpWithHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}

			if e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
				ans.L2.Lacp.Ha = &lacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				}

				if e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
					ans.L2.Lacp.Ha.Mac = &lacpHaMac{
						LacpHaEnableSameSystemMac:  util.YesNo(e.LacpHaEnableSameSystemMac),
						LacpHaSameSystemMacAddress: e.LacpHaSameSystemMacAddress,
					}
				}
			}
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.L2.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.L2.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		if text := e.raw["l2si"]; text != "" {
			ans.L2.Subinterfaces = &util.RawXml{text}
		}
	case ModeLayer3:
		ans.L3 = &layer3_v1{
			Mtu:                        e.Mtu,
			EnableUntaggedSubinterface: util.YesNo(e.EnableUntaggedSubinterface),
			StaticIps:                  util.StrToEnt(e.StaticIps),
			ManagementProfile:          e.ManagementProfile,
			NetflowProfile:             e.NetflowProfile,
		}

		if e.AdjustTcpMss || e.Ipv4MssAdjust != 0 || e.Ipv6MssAdjust != 0 {
			ans.L3.Mss = &mss{
				AdjustTcpMss:  util.YesNo(e.AdjustTcpMss),
				Ipv4MssAdjust: e.Ipv4MssAdjust,
				Ipv6MssAdjust: e.Ipv6MssAdjust,
			}
		}

		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 || e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
			ans.L3.Lacp = &lacpWithHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}

			if e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
				ans.L3.Lacp.Ha = &lacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				}

				if e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
					ans.L3.Lacp.Ha.Mac = &lacpHaMac{
						LacpHaEnableSameSystemMac:  util.YesNo(e.LacpHaEnableSameSystemMac),
						LacpHaSameSystemMacAddress: e.LacpHaSameSystemMacAddress,
					}
				}
			}
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.L3.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.L3.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6addr := e.raw["v6addr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6addr != "" || v6nd != "" {
			ans.L3.Ipv6 = &ipv6{
				Ipv6Enabled:     util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}

			if v6addr != "" {
				ans.L3.Ipv6.Address = &util.RawXml{v6addr}
			}
			if v6nd != "" {
				ans.L3.Ipv6.Neighbor = &util.RawXml{v6nd}
			}
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
			ans.L3.Dhcp = &dhcpSettings_v1{
				EnableDhcp:             util.YesNo(e.EnableDhcp),
				CreateDhcpDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				DhcpDefaultRouteMetric: e.DhcpDefaultRouteMetric,
			}
		}

		if text := e.raw["arp"]; text != "" {
			ans.L3.Arp = &util.RawXml{text}
		}
		if text := e.raw["ndp"]; text != "" {
			ans.L3.Ndp = &util.RawXml{text}
		}
		if text := e.raw["l3si"]; text != "" {
			ans.L3.Subinterfaces = &util.RawXml{text}
		}
	}

	return ans
}

type entry_v2 struct {
	XMLName       xml.Name   `xml:"entry"`
	Name          string     `xml:"name,attr"`
	Ha            *ha        `xml:"ha"`
	DecryptMirror *string    `xml:"decrypt-mirror"`
	VirtualWire   *layer2    `xml:"virtual-wire"`
	L2            *layer2    `xml:"layer2"`
	L3            *layer3_v2 `xml:"layer3"`
	Comment       string     `xml:"comment,omitempty"`
}

type layer3_v2 struct {
	Mtu                        int              `xml:"mtu,omitempty"`
	Mss                        *mss             `xml:"adjust-tcp-mss"`
	EnableUntaggedSubinterface string           `xml:"untagged-sub-interface"`
	DecryptForward             string           `xml:"decrypt-forward,omitempty"`
	Lacp                       *lacpWithHa      `xml:"lacp"`
	Lldp                       *lldp            `xml:"lldp"`
	StaticIps                  *util.EntryType  `xml:"ip"`
	Ipv6                       *ipv6            `xml:"ipv6"`
	Arp                        *util.RawXml     `xml:"arp"`
	Ndp                        *util.RawXml     `xml:"ndp-proxy"`
	ManagementProfile          string           `xml:"interface-management-profile,omitempty"`
	Dhcp                       *dhcpSettings_v1 `xml:"dhcp-client"`
	Subinterfaces              *util.RawXml     `xml:"units"`
	NetflowProfile             string           `xml:"netflow-profile,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:    e.Name,
		Comment: e.Comment,
	}

	switch e.Mode {
	case ModeHa:
		ans.Ha = &ha{}
		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 {
			ans.Ha.Lacp = &lacpNoHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}
		}
	case ModeDecryptMirror:
		s := ""
		ans.DecryptMirror = &s
	case ModeVirtualWire:
		ans.VirtualWire = &layer2{
			NetflowProfile: e.NetflowProfile,
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.VirtualWire.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.VirtualWire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		if text := e.raw["vwsi"]; text != "" {
			ans.VirtualWire.Subinterfaces = &util.RawXml{text}
		}
	case ModeLayer2:
		ans.L2 = &layer2{
			NetflowProfile: e.NetflowProfile,
		}

		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 || e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
			ans.L2.Lacp = &lacpWithHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}

			if e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
				ans.L2.Lacp.Ha = &lacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				}

				if e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
					ans.L2.Lacp.Ha.Mac = &lacpHaMac{
						LacpHaEnableSameSystemMac:  util.YesNo(e.LacpHaEnableSameSystemMac),
						LacpHaSameSystemMacAddress: e.LacpHaSameSystemMacAddress,
					}
				}
			}
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.L2.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.L2.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		if text := e.raw["l2si"]; text != "" {
			ans.L2.Subinterfaces = &util.RawXml{text}
		}
	case ModeLayer3:
		ans.L3 = &layer3_v2{
			Mtu:                        e.Mtu,
			EnableUntaggedSubinterface: util.YesNo(e.EnableUntaggedSubinterface),
			StaticIps:                  util.StrToEnt(e.StaticIps),
			ManagementProfile:          e.ManagementProfile,
			NetflowProfile:             e.NetflowProfile,
		}

		if e.DecryptForward {
			ans.L3.DecryptForward = util.YesNo(e.DecryptForward)
		}

		if e.AdjustTcpMss || e.Ipv4MssAdjust != 0 || e.Ipv6MssAdjust != 0 {
			ans.L3.Mss = &mss{
				AdjustTcpMss:  util.YesNo(e.AdjustTcpMss),
				Ipv4MssAdjust: e.Ipv4MssAdjust,
				Ipv6MssAdjust: e.Ipv6MssAdjust,
			}
		}

		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 || e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
			ans.L3.Lacp = &lacpWithHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}

			if e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
				ans.L3.Lacp.Ha = &lacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				}

				if e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
					ans.L3.Lacp.Ha.Mac = &lacpHaMac{
						LacpHaEnableSameSystemMac:  util.YesNo(e.LacpHaEnableSameSystemMac),
						LacpHaSameSystemMacAddress: e.LacpHaSameSystemMacAddress,
					}
				}
			}
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.L3.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.L3.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6addr := e.raw["v6addr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6addr != "" || v6nd != "" {
			ans.L3.Ipv6 = &ipv6{
				Ipv6Enabled:     util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}

			if v6addr != "" {
				ans.L3.Ipv6.Address = &util.RawXml{v6addr}
			}
			if v6nd != "" {
				ans.L3.Ipv6.Neighbor = &util.RawXml{v6nd}
			}
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
			ans.L3.Dhcp = &dhcpSettings_v1{
				EnableDhcp:             util.YesNo(e.EnableDhcp),
				CreateDhcpDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				DhcpDefaultRouteMetric: e.DhcpDefaultRouteMetric,
			}
		}

		if text := e.raw["arp"]; text != "" {
			ans.L3.Arp = &util.RawXml{text}
		}
		if text := e.raw["ndp"]; text != "" {
			ans.L3.Ndp = &util.RawXml{text}
		}
		if text := e.raw["l3si"]; text != "" {
			ans.L3.Subinterfaces = &util.RawXml{text}
		}
	}

	return ans
}

type entry_v3 struct {
	XMLName       xml.Name   `xml:"entry"`
	Name          string     `xml:"name,attr"`
	Ha            *ha        `xml:"ha"`
	DecryptMirror *string    `xml:"decrypt-mirror"`
	VirtualWire   *layer2    `xml:"virtual-wire"`
	L2            *layer2    `xml:"layer2"`
	L3            *layer3_v3 `xml:"layer3"`
	Comment       string     `xml:"comment,omitempty"`
}

type layer3_v3 struct {
	Mtu                        int              `xml:"mtu,omitempty"`
	Mss                        *mss             `xml:"adjust-tcp-mss"`
	EnableUntaggedSubinterface string           `xml:"untagged-sub-interface"`
	DecryptForward             string           `xml:"decrypt-forward,omitempty"`
	Lacp                       *lacpWithHa      `xml:"lacp"`
	Lldp                       *lldp            `xml:"lldp"`
	StaticIps                  *util.EntryType  `xml:"ip"`
	Ipv6                       *ipv6            `xml:"ipv6"`
	Arp                        *util.RawXml     `xml:"arp"`
	Ndp                        *util.RawXml     `xml:"ndp-proxy"`
	ManagementProfile          string           `xml:"interface-management-profile,omitempty"`
	Dhcp                       *dhcpSettings_v2 `xml:"dhcp-client"`
	Subinterfaces              *util.RawXml     `xml:"units"`
	NetflowProfile             string           `xml:"netflow-profile,omitempty"`
}

type dhcpSettings_v2 struct {
	EnableDhcp             string        `xml:"enable"`
	CreateDhcpDefaultRoute string        `xml:"create-default-route"`
	DhcpDefaultRouteMetric int           `xml:"default-route-metric,omitempty"`
	Hostname               *dhcpHostname `xml:"send-hostname"`
}

type dhcpHostname struct {
	DhcpSendHostnameEnable string `xml:"enable"`
	DhcpSendHostnameValue  string `xml:"hostname,omitempty"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:    e.Name,
		Comment: e.Comment,
	}

	switch e.Mode {
	case ModeHa:
		ans.Ha = &ha{}
		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 {
			ans.Ha.Lacp = &lacpNoHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}
		}
	case ModeDecryptMirror:
		s := ""
		ans.DecryptMirror = &s
	case ModeVirtualWire:
		ans.VirtualWire = &layer2{
			NetflowProfile: e.NetflowProfile,
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.VirtualWire.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.VirtualWire.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		if text := e.raw["vwsi"]; text != "" {
			ans.VirtualWire.Subinterfaces = &util.RawXml{text}
		}
	case ModeLayer2:
		ans.L2 = &layer2{
			NetflowProfile: e.NetflowProfile,
		}

		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 || e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
			ans.L2.Lacp = &lacpWithHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}

			if e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
				ans.L2.Lacp.Ha = &lacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				}

				if e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
					ans.L2.Lacp.Ha.Mac = &lacpHaMac{
						LacpHaEnableSameSystemMac:  util.YesNo(e.LacpHaEnableSameSystemMac),
						LacpHaSameSystemMacAddress: e.LacpHaSameSystemMacAddress,
					}
				}
			}
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.L2.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.L2.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		if text := e.raw["l2si"]; text != "" {
			ans.L2.Subinterfaces = &util.RawXml{text}
		}
	case ModeLayer3:
		ans.L3 = &layer3_v3{
			Mtu:                        e.Mtu,
			EnableUntaggedSubinterface: util.YesNo(e.EnableUntaggedSubinterface),
			StaticIps:                  util.StrToEnt(e.StaticIps),
			ManagementProfile:          e.ManagementProfile,
			NetflowProfile:             e.NetflowProfile,
		}

		if e.DecryptForward {
			ans.L3.DecryptForward = util.YesNo(e.DecryptForward)
		}

		if e.AdjustTcpMss || e.Ipv4MssAdjust != 0 || e.Ipv6MssAdjust != 0 {
			ans.L3.Mss = &mss{
				AdjustTcpMss:  util.YesNo(e.AdjustTcpMss),
				Ipv4MssAdjust: e.Ipv4MssAdjust,
				Ipv6MssAdjust: e.Ipv6MssAdjust,
			}
		}

		if e.LacpEnable || e.LacpFastFailover || e.LacpMode != "" || e.LacpTransmissionRate != "" || e.LacpSystemPriority != 0 || e.LacpMaxPorts != 0 || e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
			ans.L3.Lacp = &lacpWithHa{
				LacpEnable:           util.YesNo(e.LacpEnable),
				LacpFastFailover:     util.YesNo(e.LacpFastFailover),
				LacpMode:             e.LacpMode,
				LacpTransmissionRate: e.LacpTransmissionRate,
				LacpSystemPriority:   e.LacpSystemPriority,
				LacpMaxPorts:         e.LacpMaxPorts,
			}

			if e.LacpHaPassivePreNegotiation || e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
				ans.L3.Lacp.Ha = &lacpHa{
					LacpHaPassivePreNegotiation: util.YesNo(e.LacpHaPassivePreNegotiation),
				}

				if e.LacpHaEnableSameSystemMac || e.LacpHaSameSystemMacAddress != "" {
					ans.L3.Lacp.Ha.Mac = &lacpHaMac{
						LacpHaEnableSameSystemMac:  util.YesNo(e.LacpHaEnableSameSystemMac),
						LacpHaSameSystemMacAddress: e.LacpHaSameSystemMacAddress,
					}
				}
			}
		}

		if e.LldpEnable || e.LldpProfile != "" || e.LldpHaPassivePreNegotiation {
			ans.L3.Lldp = &lldp{
				LldpEnable:  util.YesNo(e.LldpEnable),
				LldpProfile: e.LldpProfile,
			}
			if e.LldpHaPassivePreNegotiation {
				ans.L3.Lldp.Ha = &lldpHa{
					LldpHaPassivePreNegotiation: util.YesNo(e.LldpHaPassivePreNegotiation),
				}
			}
		}

		v6addr := e.raw["v6addr"]
		v6nd := e.raw["v6nd"]
		if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6addr != "" || v6nd != "" {
			ans.L3.Ipv6 = &ipv6{
				Ipv6Enabled:     util.YesNo(e.Ipv6Enabled),
				Ipv6InterfaceId: e.Ipv6InterfaceId,
			}

			if v6addr != "" {
				ans.L3.Ipv6.Address = &util.RawXml{v6addr}
			}
			if v6nd != "" {
				ans.L3.Ipv6.Neighbor = &util.RawXml{v6nd}
			}
		}

		if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 || e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != "" {
			ans.L3.Dhcp = &dhcpSettings_v2{
				EnableDhcp:             util.YesNo(e.EnableDhcp),
				CreateDhcpDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
				DhcpDefaultRouteMetric: e.DhcpDefaultRouteMetric,
			}

			if e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != "" {
				ans.L3.Dhcp.Hostname = &dhcpHostname{
					DhcpSendHostnameEnable: util.YesNo(e.DhcpSendHostnameEnable),
					DhcpSendHostnameValue:  e.DhcpSendHostnameValue,
				}
			}
		}

		if text := e.raw["arp"]; text != "" {
			ans.L3.Arp = &util.RawXml{text}
		}
		if text := e.raw["ndp"]; text != "" {
			ans.L3.Ndp = &util.RawXml{text}
		}
		if text := e.raw["l3si"]; text != "" {
			ans.L3.Subinterfaces = &util.RawXml{text}
		}
	}

	return ans
}
