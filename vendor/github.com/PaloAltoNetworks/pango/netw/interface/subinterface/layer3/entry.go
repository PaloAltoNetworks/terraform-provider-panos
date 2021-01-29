package layer3

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a layer3
// subinterface.
type Entry struct {
	Name                   string
	Tag                    int
	StaticIps              []string // ordered
	Ipv6Enabled            bool
	Ipv6InterfaceId        string
	ManagementProfile      string
	Mtu                    int
	AdjustTcpMss           bool
	Ipv4MssAdjust          int
	Ipv6MssAdjust          int
	NetflowProfile         string
	Comment                string
	EnableDhcp             bool
	CreateDhcpDefaultRoute bool
	DhcpDefaultRouteMetric int
	DhcpSendHostnameEnable bool   // 9.0
	DhcpSendHostnameValue  string // 9.0
	DecryptForward         bool   // 8.1

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Tag = s.Tag
	o.StaticIps = s.StaticIps
	o.Ipv6Enabled = s.Ipv6Enabled
	o.ManagementProfile = s.ManagementProfile
	o.Mtu = s.Mtu
	o.AdjustTcpMss = s.AdjustTcpMss
	o.NetflowProfile = s.NetflowProfile
	o.Comment = s.Comment
	o.Ipv4MssAdjust = s.Ipv4MssAdjust
	o.Ipv6MssAdjust = s.Ipv6MssAdjust
	o.EnableDhcp = s.EnableDhcp
	o.CreateDhcpDefaultRoute = s.CreateDhcpDefaultRoute
	o.DhcpDefaultRouteMetric = s.DhcpDefaultRouteMetric
	o.DhcpSendHostnameEnable = s.DhcpSendHostnameEnable
	o.DhcpSendHostnameValue = s.DhcpSendHostnameValue
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, string, interface{}) {
	_, fn := versioning(v)
	return o.Name, o.Name, fn(o)
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
		Name:              o.Name,
		Tag:               o.Tag,
		StaticIps:         util.EntToStr(o.StaticIps),
		ManagementProfile: o.ManagementProfile,
		Mtu:               o.Mtu,
		NetflowProfile:    o.NetflowProfile,
		Comment:           o.Comment,
	}
	ans.raw = make(map[string]string)

	if o.Ipv6 != nil {
		ans.Ipv6Enabled = util.AsBool(o.Ipv6.Ipv6Enabled)
		ans.Ipv6InterfaceId = o.Ipv6.Ipv6InterfaceId
		if o.Ipv6.Addresses != nil {
			ans.raw["v6adr"] = util.CleanRawXml(o.Ipv6.Addresses.Text)
		}
		if o.Ipv6.Neighbor != nil {
			ans.raw["v6nbr"] = util.CleanRawXml(o.Ipv6.Neighbor.Text)
		}
	}

	if o.Mss != nil {
		ans.AdjustTcpMss = util.AsBool(o.Mss.AdjustTcpMss)
		ans.Ipv4MssAdjust = o.Mss.Ipv4MssAdjust
		ans.Ipv6MssAdjust = o.Mss.Ipv6MssAdjust
	}

	if o.Dhcp != nil {
		ans.EnableDhcp = util.AsBool(o.Dhcp.EnableDhcp)
		ans.CreateDhcpDefaultRoute = util.AsBool(o.Dhcp.CreateDhcpDefaultRoute)
		ans.DhcpDefaultRouteMetric = o.Dhcp.DhcpDefaultRouteMetric
	}

	if o.Arp != nil {
		ans.raw["arp"] = util.CleanRawXml(o.Arp.Text)
	}
	if o.NdpProxy != nil {
		ans.raw["ndp"] = util.CleanRawXml(o.NdpProxy.Text)
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}
	return ans
}

type entry_v1 struct {
	XMLName           xml.Name        `xml:"entry"`
	Name              string          `xml:"name,attr"`
	Tag               int             `xml:"tag,omitempty"`
	StaticIps         *util.EntryType `xml:"ip"`
	Ipv6              *ipv6           `xml:"ipv6"`
	ManagementProfile string          `xml:"interface-management-profile,omitempty"`
	Mtu               int             `xml:"mtu,omitempty"`
	Mss               *mss            `xml:"adjust-tcp-mss"`
	NetflowProfile    string          `xml:"netflow-profile,omitempty"`
	Comment           string          `xml:"comment,omitempty"`
	Dhcp              *dhcp_v1        `xml:"dhcp-client"`
	Arp               *util.RawXml    `xml:"arp"`
	NdpProxy          *util.RawXml    `xml:"ndp-proxy"`
}

type ipv6 struct {
	Ipv6Enabled     string       `xml:"enabled"`
	Ipv6InterfaceId string       `xml:"interface-id,omitempty"`
	Addresses       *util.RawXml `xml:"address"`
	Neighbor        *util.RawXml `xml:"neighbor-discovery"`
}

type mss struct {
	AdjustTcpMss  string `xml:"enable"`
	Ipv4MssAdjust int    `xml:"ipv4-mss-adjustment,omitempty"`
	Ipv6MssAdjust int    `xml:"ipv6-mss-adjustment,omitempty"`
}

type dhcp_v1 struct {
	EnableDhcp             string `xml:"enable"`
	CreateDhcpDefaultRoute string `xml:"create-default-route"`
	DhcpDefaultRouteMetric int    `xml:"default-route-metric,omitempty"`
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
		Name:              o.Name,
		Tag:               o.Tag,
		StaticIps:         util.EntToStr(o.StaticIps),
		ManagementProfile: o.ManagementProfile,
		Mtu:               o.Mtu,
		NetflowProfile:    o.NetflowProfile,
		Comment:           o.Comment,
		DecryptForward:    util.AsBool(o.DecryptForward),
	}
	ans.raw = make(map[string]string)

	if o.Ipv6 != nil {
		ans.Ipv6Enabled = util.AsBool(o.Ipv6.Ipv6Enabled)
		ans.Ipv6InterfaceId = o.Ipv6.Ipv6InterfaceId
		if o.Ipv6.Addresses != nil {
			ans.raw["v6adr"] = util.CleanRawXml(o.Ipv6.Addresses.Text)
		}
		if o.Ipv6.Neighbor != nil {
			ans.raw["v6nbr"] = util.CleanRawXml(o.Ipv6.Neighbor.Text)
		}
	}

	if o.Mss != nil {
		ans.AdjustTcpMss = util.AsBool(o.Mss.AdjustTcpMss)
		ans.Ipv4MssAdjust = o.Mss.Ipv4MssAdjust
		ans.Ipv6MssAdjust = o.Mss.Ipv6MssAdjust
	}

	if o.Dhcp != nil {
		ans.EnableDhcp = util.AsBool(o.Dhcp.EnableDhcp)
		ans.CreateDhcpDefaultRoute = util.AsBool(o.Dhcp.CreateDhcpDefaultRoute)
		ans.DhcpDefaultRouteMetric = o.Dhcp.DhcpDefaultRouteMetric
	}

	if o.Arp != nil {
		ans.raw["arp"] = util.CleanRawXml(o.Arp.Text)
	}
	if o.NdpProxy != nil {
		ans.raw["ndp"] = util.CleanRawXml(o.NdpProxy.Text)
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
		Name:              o.Name,
		Tag:               o.Tag,
		StaticIps:         util.EntToStr(o.StaticIps),
		ManagementProfile: o.ManagementProfile,
		Mtu:               o.Mtu,
		NetflowProfile:    o.NetflowProfile,
		Comment:           o.Comment,
		DecryptForward:    util.AsBool(o.DecryptForward),
	}
	ans.raw = make(map[string]string)

	if o.Ipv6 != nil {
		ans.Ipv6Enabled = util.AsBool(o.Ipv6.Ipv6Enabled)
		ans.Ipv6InterfaceId = o.Ipv6.Ipv6InterfaceId
		if o.Ipv6.Addresses != nil {
			ans.raw["v6adr"] = util.CleanRawXml(o.Ipv6.Addresses.Text)
		}
		if o.Ipv6.Neighbor != nil {
			ans.raw["v6nbr"] = util.CleanRawXml(o.Ipv6.Neighbor.Text)
		}
	}

	if o.Mss != nil {
		ans.AdjustTcpMss = util.AsBool(o.Mss.AdjustTcpMss)
		ans.Ipv4MssAdjust = o.Mss.Ipv4MssAdjust
		ans.Ipv6MssAdjust = o.Mss.Ipv6MssAdjust
	}

	if o.Dhcp != nil {
		ans.EnableDhcp = util.AsBool(o.Dhcp.EnableDhcp)
		ans.CreateDhcpDefaultRoute = util.AsBool(o.Dhcp.CreateDhcpDefaultRoute)
		ans.DhcpDefaultRouteMetric = o.Dhcp.DhcpDefaultRouteMetric
		if o.Dhcp.Hostname != nil {
			ans.DhcpSendHostnameEnable = util.AsBool(o.Dhcp.Hostname.DhcpSendHostnameEnable)
			ans.DhcpSendHostnameValue = o.Dhcp.Hostname.DhcpSendHostnameValue
		}
	}

	if o.Arp != nil {
		ans.raw["arp"] = util.CleanRawXml(o.Arp.Text)
	}
	if o.NdpProxy != nil {
		ans.raw["ndp"] = util.CleanRawXml(o.NdpProxy.Text)
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}
	return ans
}

// 8.1
type entry_v2 struct {
	XMLName           xml.Name        `xml:"entry"`
	Name              string          `xml:"name,attr"`
	Tag               int             `xml:"tag,omitempty"`
	StaticIps         *util.EntryType `xml:"ip"`
	Ipv6              *ipv6           `xml:"ipv6"`
	ManagementProfile string          `xml:"interface-management-profile,omitempty"`
	Mtu               int             `xml:"mtu,omitempty"`
	Mss               *mss            `xml:"adjust-tcp-mss"`
	NetflowProfile    string          `xml:"netflow-profile,omitempty"`
	Comment           string          `xml:"comment,omitempty"`
	Dhcp              *dhcp_v1        `xml:"dhcp-client"`
	Arp               *util.RawXml    `xml:"arp"`
	NdpProxy          *util.RawXml    `xml:"ndp-proxy"`
	DecryptForward    string          `xml:"decrypt-forward,omitempty"`
}

// 9.0
type entry_v3 struct {
	XMLName           xml.Name        `xml:"entry"`
	Name              string          `xml:"name,attr"`
	Tag               int             `xml:"tag,omitempty"`
	StaticIps         *util.EntryType `xml:"ip"`
	Ipv6              *ipv6           `xml:"ipv6"`
	ManagementProfile string          `xml:"interface-management-profile,omitempty"`
	Mtu               int             `xml:"mtu,omitempty"`
	Mss               *mss            `xml:"adjust-tcp-mss"`
	NetflowProfile    string          `xml:"netflow-profile,omitempty"`
	Comment           string          `xml:"comment,omitempty"`
	Dhcp              *dhcp_v2        `xml:"dhcp-client"`
	Arp               *util.RawXml    `xml:"arp"`
	NdpProxy          *util.RawXml    `xml:"ndp-proxy"`
	DecryptForward    string          `xml:"decrypt-forward,omitempty"`
	DdnsConfig        *util.RawXml    `xml:"ddns-config"`
}

type dhcp_v2 struct {
	EnableDhcp             string        `xml:"enable"`
	CreateDhcpDefaultRoute string        `xml:"create-default-route"`
	DhcpDefaultRouteMetric int           `xml:"default-route-metric,omitempty"`
	Hostname               *dhcpHostname `xml:"send-hostname"`
}

type dhcpHostname struct {
	DhcpSendHostnameEnable string `xml:"enable"`
	DhcpSendHostnameValue  string `xml:"hostname,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:              e.Name,
		Tag:               e.Tag,
		StaticIps:         util.StrToEnt(e.StaticIps),
		ManagementProfile: e.ManagementProfile,
		Mtu:               e.Mtu,
		NetflowProfile:    e.NetflowProfile,
		Comment:           e.Comment,
	}

	v6adr := e.raw["v6adr"]
	v6nbr := e.raw["v6nbr"]
	if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nbr != "" {
		i6 := ipv6{
			Ipv6Enabled:     util.YesNo(e.Ipv6Enabled),
			Ipv6InterfaceId: e.Ipv6InterfaceId,
		}

		if v6adr != "" {
			i6.Addresses = &util.RawXml{v6adr}
		}
		if v6nbr != "" {
			i6.Neighbor = &util.RawXml{v6nbr}
		}
		ans.Ipv6 = &i6
	}

	if e.AdjustTcpMss || e.Ipv4MssAdjust != 0 || e.Ipv6MssAdjust != 0 {
		ans.Mss = &mss{
			AdjustTcpMss:  util.YesNo(e.AdjustTcpMss),
			Ipv4MssAdjust: e.Ipv4MssAdjust,
			Ipv6MssAdjust: e.Ipv6MssAdjust,
		}
	}

	if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
		ans.Dhcp = &dhcp_v1{
			EnableDhcp:             util.YesNo(e.EnableDhcp),
			CreateDhcpDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
			DhcpDefaultRouteMetric: e.DhcpDefaultRouteMetric,
		}
	}

	if text, present := e.raw["arp"]; present {
		ans.Arp = &util.RawXml{text}
	}
	if text, present := e.raw["ndp"]; present {
		ans.NdpProxy = &util.RawXml{text}
	}

	return ans
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:              e.Name,
		Tag:               e.Tag,
		StaticIps:         util.StrToEnt(e.StaticIps),
		ManagementProfile: e.ManagementProfile,
		Mtu:               e.Mtu,
		NetflowProfile:    e.NetflowProfile,
		Comment:           e.Comment,
	}

	if e.DecryptForward {
		ans.DecryptForward = util.YesNo(e.DecryptForward)
	}

	v6adr := e.raw["v6adr"]
	v6nbr := e.raw["v6nbr"]
	if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nbr != "" {
		i6 := ipv6{
			Ipv6Enabled:     util.YesNo(e.Ipv6Enabled),
			Ipv6InterfaceId: e.Ipv6InterfaceId,
		}

		if v6adr != "" {
			i6.Addresses = &util.RawXml{v6adr}
		}
		if v6nbr != "" {
			i6.Neighbor = &util.RawXml{v6nbr}
		}
		ans.Ipv6 = &i6
	}

	if e.AdjustTcpMss || e.Ipv4MssAdjust != 0 || e.Ipv6MssAdjust != 0 {
		ans.Mss = &mss{
			AdjustTcpMss:  util.YesNo(e.AdjustTcpMss),
			Ipv4MssAdjust: e.Ipv4MssAdjust,
			Ipv6MssAdjust: e.Ipv6MssAdjust,
		}
	}

	if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 {
		ans.Dhcp = &dhcp_v1{
			EnableDhcp:             util.YesNo(e.EnableDhcp),
			CreateDhcpDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
			DhcpDefaultRouteMetric: e.DhcpDefaultRouteMetric,
		}
	}

	if text, present := e.raw["arp"]; present {
		ans.Arp = &util.RawXml{text}
	}
	if text, present := e.raw["ndp"]; present {
		ans.NdpProxy = &util.RawXml{text}
	}

	return ans
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:              e.Name,
		Tag:               e.Tag,
		StaticIps:         util.StrToEnt(e.StaticIps),
		ManagementProfile: e.ManagementProfile,
		Mtu:               e.Mtu,
		NetflowProfile:    e.NetflowProfile,
		Comment:           e.Comment,
	}

	if e.DecryptForward {
		ans.DecryptForward = util.YesNo(e.DecryptForward)
	}

	v6adr := e.raw["v6adr"]
	v6nbr := e.raw["v6nbr"]
	if e.Ipv6Enabled || e.Ipv6InterfaceId != "" || v6adr != "" || v6nbr != "" {
		i6 := ipv6{
			Ipv6Enabled:     util.YesNo(e.Ipv6Enabled),
			Ipv6InterfaceId: e.Ipv6InterfaceId,
		}

		if v6adr != "" {
			i6.Addresses = &util.RawXml{v6adr}
		}
		if v6nbr != "" {
			i6.Neighbor = &util.RawXml{v6nbr}
		}
		ans.Ipv6 = &i6
	}

	if e.AdjustTcpMss || e.Ipv4MssAdjust != 0 || e.Ipv6MssAdjust != 0 {
		ans.Mss = &mss{
			AdjustTcpMss:  util.YesNo(e.AdjustTcpMss),
			Ipv4MssAdjust: e.Ipv4MssAdjust,
			Ipv6MssAdjust: e.Ipv6MssAdjust,
		}
	}

	dhn := e.DhcpSendHostnameEnable || e.DhcpSendHostnameValue != ""
	if e.EnableDhcp || e.CreateDhcpDefaultRoute || e.DhcpDefaultRouteMetric != 0 || dhn {
		ans.Dhcp = &dhcp_v2{
			EnableDhcp:             util.YesNo(e.EnableDhcp),
			CreateDhcpDefaultRoute: util.YesNo(e.CreateDhcpDefaultRoute),
			DhcpDefaultRouteMetric: e.DhcpDefaultRouteMetric,
		}

		if dhn {
			ans.Dhcp.Hostname = &dhcpHostname{
				DhcpSendHostnameEnable: util.YesNo(e.DhcpSendHostnameEnable),
				DhcpSendHostnameValue:  e.DhcpSendHostnameValue,
			}
		}
	}

	if text, present := e.raw["arp"]; present {
		ans.Arp = &util.RawXml{text}
	}
	if text, present := e.raw["ndp"]; present {
		ans.NdpProxy = &util.RawXml{text}
	}

	return ans
}
