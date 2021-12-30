package zone

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a zone.
type Entry struct {
	Name                         string
	Mode                         string
	Interfaces                   []string // unordered
	ZoneProfile                  string
	LogSetting                   string
	EnableUserId                 bool
	IncludeAcls                  []string // unordered
	ExcludeAcls                  []string // unordered
	EnablePacketBufferProtection bool     // 8.0+
	EnableDeviceIdentification   bool     // 10.0+
	DeviceIncludeAcls            []string // 10.0+, unordered?
	DeviceExcludeAcls            []string // 10.0+, unordered?
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Mode = s.Mode
	if s.Interfaces == nil {
		o.Interfaces = nil
	} else {
		o.Interfaces = make([]string, len(s.Interfaces))
		copy(o.Interfaces, s.Interfaces)
	}
	o.ZoneProfile = s.ZoneProfile
	o.LogSetting = s.LogSetting
	o.EnableUserId = s.EnableUserId
	if s.IncludeAcls == nil {
		o.IncludeAcls = nil
	} else {
		o.IncludeAcls = make([]string, len(s.IncludeAcls))
		copy(o.IncludeAcls, s.IncludeAcls)
	}
	if s.ExcludeAcls == nil {
		o.ExcludeAcls = nil
	} else {
		o.ExcludeAcls = make([]string, len(s.ExcludeAcls))
		copy(o.ExcludeAcls, s.ExcludeAcls)
	}
	o.EnablePacketBufferProtection = s.EnablePacketBufferProtection
	o.EnableDeviceIdentification = s.EnableDeviceIdentification
	if s.DeviceIncludeAcls == nil {
		o.DeviceIncludeAcls = nil
	} else {
		o.DeviceIncludeAcls = make([]string, len(s.DeviceIncludeAcls))
		copy(o.DeviceIncludeAcls, s.DeviceIncludeAcls)
	}
	if s.DeviceExcludeAcls == nil {
		o.DeviceExcludeAcls = nil
	} else {
		o.DeviceExcludeAcls = make([]string, len(s.DeviceExcludeAcls))
		copy(o.DeviceExcludeAcls, s.DeviceExcludeAcls)
	}
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

// 6.1
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
		Name:         o.Name,
		EnableUserId: util.AsBool(o.EnableUserId),
	}

	if o.Network != nil {
		ans.ZoneProfile = o.Network.ZoneProfile
		ans.LogSetting = o.Network.LogSetting

		var ilist []string
		if o.Network.Tap != nil {
			ans.Mode = ModeTap
			ilist = util.MemToStr(o.Network.Tap)
		} else if o.Network.VWire != nil {
			ans.Mode = ModeVirtualWire
			ilist = util.MemToStr(o.Network.VWire)
		} else if o.Network.L2 != nil {
			ans.Mode = ModeL2
			ilist = util.MemToStr(o.Network.L2)
		} else if o.Network.L3 != nil {
			ans.Mode = ModeL3
			ilist = util.MemToStr(o.Network.L3)
		} else if o.Network.External != nil {
			ans.Mode = ModeExternal
			ilist = util.MemToStr(o.Network.External)
		}

		if len(ilist) > 0 {
			ans.Interfaces = ilist
		}
	}

	if o.UserAcls != nil {
		ans.IncludeAcls = util.MemToStr(o.UserAcls.IncludeAcls)
		ans.ExcludeAcls = util.MemToStr(o.UserAcls.ExcludeAcls)
	}

	return ans
}

type entry_v1 struct {
	XMLName      xml.Name    `xml:"entry"`
	Name         string      `xml:"name,attr"`
	EnableUserId string      `xml:"enable-user-identification"`
	Network      *network_v1 `xml:"network"`
	UserAcls     *acls       `xml:"user-acl"`
}

type network_v1 struct {
	ZoneProfile string           `xml:"zone-protection-profile,omitempty"`
	LogSetting  string           `xml:"log-setting,omitempty"`
	Tap         *util.MemberType `xml:"tap"`
	VWire       *util.MemberType `xml:"virtual-wire"`
	L2          *util.MemberType `xml:"layer2"`
	L3          *util.MemberType `xml:"layer3"`
	External    *util.MemberType `xml:"external"`
}

type acls struct {
	IncludeAcls *util.MemberType `xml:"include-list"`
	ExcludeAcls *util.MemberType `xml:"exclude-list"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:         e.Name,
		EnableUserId: util.YesNo(e.EnableUserId),
	}

	if e.Mode != "" || e.ZoneProfile != "" || e.LogSetting != "" {
		ans.Network = &network_v1{
			ZoneProfile: e.ZoneProfile,
			LogSetting:  e.LogSetting,
		}

		ilist := &util.MemberType{}
		if len(e.Interfaces) > 0 {
			ilist = util.StrToMem(e.Interfaces)
		}

		switch e.Mode {
		case ModeTap:
			ans.Network.Tap = ilist
		case ModeVirtualWire:
			ans.Network.VWire = ilist
		case ModeL2:
			ans.Network.L2 = ilist
		case ModeL3:
			ans.Network.L3 = ilist
		case ModeExternal:
			ans.Network.External = ilist
		}
	}

	if len(e.IncludeAcls) > 0 || len(e.ExcludeAcls) > 0 {
		ans.UserAcls = &acls{
			IncludeAcls: util.StrToMem(e.IncludeAcls),
			ExcludeAcls: util.StrToMem(e.ExcludeAcls),
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
		Name:         o.Name,
		EnableUserId: util.AsBool(o.EnableUserId),
	}

	if o.Network != nil {
		ans.ZoneProfile = o.Network.ZoneProfile
		ans.EnablePacketBufferProtection = util.AsBool(o.Network.EnablePacketBufferProtection)
		ans.LogSetting = o.Network.LogSetting

		var ilist []string
		if o.Network.Tap != nil {
			ans.Mode = ModeTap
			ilist = util.MemToStr(o.Network.Tap)
		} else if o.Network.VWire != nil {
			ans.Mode = ModeVirtualWire
			ilist = util.MemToStr(o.Network.VWire)
		} else if o.Network.L2 != nil {
			ans.Mode = ModeL2
			ilist = util.MemToStr(o.Network.L2)
		} else if o.Network.L3 != nil {
			ans.Mode = ModeL3
			ilist = util.MemToStr(o.Network.L3)
		} else if o.Network.External != nil {
			ans.Mode = ModeExternal
			ilist = util.MemToStr(o.Network.External)
		} else if o.Network.Tunnel != nil {
			ans.Mode = ModeTunnel
		}

		if len(ilist) > 0 {
			ans.Interfaces = ilist
		}
	}

	if o.UserAcls != nil {
		ans.IncludeAcls = util.MemToStr(o.UserAcls.IncludeAcls)
		ans.ExcludeAcls = util.MemToStr(o.UserAcls.ExcludeAcls)
	}

	return ans
}

type entry_v2 struct {
	XMLName      xml.Name    `xml:"entry"`
	Name         string      `xml:"name,attr"`
	EnableUserId string      `xml:"enable-user-identification"`
	Network      *network_v2 `xml:"network"`
	UserAcls     *acls       `xml:"user-acl"`
}

type network_v2 struct {
	ZoneProfile                  string           `xml:"zone-protection-profile,omitempty"`
	EnablePacketBufferProtection string           `xml:"enable-packet-buffer-protection"`
	LogSetting                   string           `xml:"log-setting,omitempty"`
	Tap                          *util.MemberType `xml:"tap"`
	VWire                        *util.MemberType `xml:"virtual-wire"`
	L2                           *util.MemberType `xml:"layer2"`
	L3                           *util.MemberType `xml:"layer3"`
	External                     *util.MemberType `xml:"external"`
	Tunnel                       *string          `xml:"tunnel"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:         e.Name,
		EnableUserId: util.YesNo(e.EnableUserId),
	}

	if e.Mode != "" || e.ZoneProfile != "" || e.EnablePacketBufferProtection || e.LogSetting != "" {
		ans.Network = &network_v2{
			ZoneProfile:                  e.ZoneProfile,
			EnablePacketBufferProtection: util.YesNo(e.EnablePacketBufferProtection),
			LogSetting:                   e.LogSetting,
		}

		ilist := &util.MemberType{}
		if len(e.Interfaces) > 0 {
			ilist = util.StrToMem(e.Interfaces)
		}

		switch e.Mode {
		case ModeTap:
			ans.Network.Tap = ilist
		case ModeVirtualWire:
			ans.Network.VWire = ilist
		case ModeL2:
			ans.Network.L2 = ilist
		case ModeL3:
			ans.Network.L3 = ilist
		case ModeExternal:
			ans.Network.External = ilist
		case ModeTunnel:
			s := ""
			ans.Network.Tunnel = &s
		}
	}

	if len(e.IncludeAcls) > 0 || len(e.ExcludeAcls) > 0 {
		ans.UserAcls = &acls{
			IncludeAcls: util.StrToMem(e.IncludeAcls),
			ExcludeAcls: util.StrToMem(e.ExcludeAcls),
		}
	}

	return ans
}

// 10.0
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
		Name:                       o.Name,
		EnableUserId:               util.AsBool(o.EnableUserId),
		EnableDeviceIdentification: util.AsBool(o.EnableDeviceIdentification),
	}

	if o.Network != nil {
		ans.ZoneProfile = o.Network.ZoneProfile
		ans.EnablePacketBufferProtection = util.AsBool(o.Network.EnablePacketBufferProtection)
		ans.LogSetting = o.Network.LogSetting

		var ilist []string
		if o.Network.Tap != nil {
			ans.Mode = ModeTap
			ilist = util.MemToStr(o.Network.Tap)
		} else if o.Network.VWire != nil {
			ans.Mode = ModeVirtualWire
			ilist = util.MemToStr(o.Network.VWire)
		} else if o.Network.L2 != nil {
			ans.Mode = ModeL2
			ilist = util.MemToStr(o.Network.L2)
		} else if o.Network.L3 != nil {
			ans.Mode = ModeL3
			ilist = util.MemToStr(o.Network.L3)
		} else if o.Network.External != nil {
			ans.Mode = ModeExternal
			ilist = util.MemToStr(o.Network.External)
		} else if o.Network.Tunnel != nil {
			ans.Mode = ModeTunnel
		}

		if len(ilist) > 0 {
			ans.Interfaces = ilist
		}
	}

	if o.UserAcls != nil {
		ans.IncludeAcls = util.MemToStr(o.UserAcls.IncludeAcls)
		ans.ExcludeAcls = util.MemToStr(o.UserAcls.ExcludeAcls)
	}

	if o.DeviceAcls != nil {
		ans.DeviceIncludeAcls = util.MemToStr(o.DeviceAcls.IncludeAcls)
		ans.DeviceExcludeAcls = util.MemToStr(o.DeviceAcls.ExcludeAcls)
	}

	return ans
}

type entry_v3 struct {
	XMLName                    xml.Name    `xml:"entry"`
	Name                       string      `xml:"name,attr"`
	EnableUserId               string      `xml:"enable-user-identification"`
	EnableDeviceIdentification string      `xml:"enable-device-identification"`
	Network                    *network_v2 `xml:"network"`
	UserAcls                   *acls       `xml:"user-acl"`
	DeviceAcls                 *acls       `xml:"device-acl"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:                       e.Name,
		EnableUserId:               util.YesNo(e.EnableUserId),
		EnableDeviceIdentification: util.YesNo(e.EnableDeviceIdentification),
	}

	if e.Mode != "" || e.ZoneProfile != "" || e.EnablePacketBufferProtection || e.LogSetting != "" {
		ans.Network = &network_v2{
			ZoneProfile:                  e.ZoneProfile,
			EnablePacketBufferProtection: util.YesNo(e.EnablePacketBufferProtection),
			LogSetting:                   e.LogSetting,
		}

		ilist := &util.MemberType{}
		if len(e.Interfaces) > 0 {
			ilist = util.StrToMem(e.Interfaces)
		}

		switch e.Mode {
		case ModeTap:
			ans.Network.Tap = ilist
		case ModeVirtualWire:
			ans.Network.VWire = ilist
		case ModeL2:
			ans.Network.L2 = ilist
		case ModeL3:
			ans.Network.L3 = ilist
		case ModeExternal:
			ans.Network.External = ilist
		case ModeTunnel:
			s := ""
			ans.Network.Tunnel = &s
		}
	}

	if len(e.IncludeAcls) > 0 || len(e.ExcludeAcls) > 0 {
		ans.UserAcls = &acls{
			IncludeAcls: util.StrToMem(e.IncludeAcls),
			ExcludeAcls: util.StrToMem(e.ExcludeAcls),
		}
	}

	if len(e.DeviceIncludeAcls) > 0 || len(e.DeviceExcludeAcls) > 0 {
		ans.DeviceAcls = &acls{
			IncludeAcls: util.StrToMem(e.DeviceIncludeAcls),
			ExcludeAcls: util.StrToMem(e.DeviceExcludeAcls),
		}
	}

	return ans
}
