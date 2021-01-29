package nat

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a NAT
// policy.  The prefix "Sat" stands for "Source Address Translation" while
// the prefix "Dat" stands for "Destination Address Translation".
//
// Targets is a map where the key is the serial number of the target device and
// the value is a list of specific vsys on that device.  The list of vsys is
// nil if all vsys on that device should be included or if the device is a
// virtual firewall (and thus only has vsys1).
//
// The following Sat params are linked:
//
// SatType = nat.DynamicIpAndPort && SatAddressType = nat.TranslatedAddress:
//
//      * SatTranslatedAddresses
//
// SatType = nat.DynamicIpAndPort && SatAddressType = nat.InterfaceAddress:
//
//      * SatInterface
//      * SatIpAddress
//
// For ALL SatType = nat.DynamicIp:
//
//      * SatTranslatedAddresses
//
// For ALL SatType = nat.DynamicIp and SatFallbackType = nat.InterfaceAddress:
//
//      * SatFallbackInterface
//
// SatType = nat.DynamicIp && SatFallbackType = nat.InterfaceAddress && SatFallbackIpType = nat.Ip:
//
//      * SatFallbackIpAddress
//
// SatType = nat.DynamicIp && SatFallbackType = nat.InterfaceAddress && SatFallbackIpType = nat.FloatingIp:
//
//      * SatFallbackIpAddress
//
// SatType = nat.DynamicIp and SatFallbackType = nat.TranslatedAddress:
//
//      * SatFallbackTranslatedAddresses
//
// SatType = nat.StaticIp:
//
//      * SatStaticTranslatedAddress
//      * SatStaticBiDirectional
//
// If both DatAddress and DatPort are unintialized, then no destination
// address translation will be enabled; setting DatType by itself is not
// good enough.
type Entry struct {
	Name                           string
	Description                    string
	Type                           string
	SourceZones                    []string // unordered
	DestinationZone                string
	ToInterface                    string
	Service                        string
	SourceAddresses                []string // unordered
	DestinationAddresses           []string // unordered
	SatType                        string
	SatAddressType                 string
	SatTranslatedAddresses         []string // unordered
	SatInterface                   string
	SatIpAddress                   string
	SatFallbackType                string
	SatFallbackTranslatedAddresses []string // unordered
	SatFallbackInterface           string
	SatFallbackIpType              string
	SatFallbackIpAddress           string
	SatStaticTranslatedAddress     string
	SatStaticBiDirectional         bool
	DatType                        string
	DatAddress                     string
	DatPort                        int
	DatDynamicDistribution         string // 8.1+
	Disabled                       bool
	Targets                        map[string][]string
	NegateTarget                   bool
	Tags                           []string // ordered
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * Type: "ipv4"
//      * ToInterface: "any"
//      * Service: "any"
//      * SourceAddresses: ["any"]
//      * DestinationAddresses: ["any"]
//      * SatType: None
//      * DatType: DatTypeStatic
func (o *Entry) Defaults() {
	if o.Type == "" {
		o.Type = "ipv4"
	}

	if o.ToInterface == "" {
		o.ToInterface = "any"
	}

	if o.Service == "" {
		o.Service = "any"
	}

	if len(o.SourceAddresses) == 0 {
		o.SourceAddresses = []string{"any"}
	}

	if len(o.DestinationAddresses) == 0 {
		o.DestinationAddresses = []string{"any"}
	}

	if o.SatType == "" {
		o.SatType = None
	}

	if o.DatType == "" {
		o.DatType = DatTypeStatic
	}
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.Type = s.Type
	o.SourceZones = s.SourceZones
	o.DestinationZone = s.DestinationZone
	o.ToInterface = s.ToInterface
	o.Service = s.Service
	o.SourceAddresses = s.SourceAddresses
	o.DestinationAddresses = s.DestinationAddresses
	o.SatType = s.SatType
	o.SatAddressType = s.SatAddressType
	o.SatTranslatedAddresses = s.SatTranslatedAddresses
	o.SatInterface = s.SatInterface
	o.SatIpAddress = s.SatIpAddress
	o.SatFallbackType = s.SatFallbackType
	o.SatFallbackTranslatedAddresses = s.SatFallbackTranslatedAddresses
	o.SatFallbackInterface = s.SatFallbackInterface
	o.SatFallbackIpType = s.SatFallbackIpType
	o.SatFallbackIpAddress = s.SatFallbackIpAddress
	o.SatStaticTranslatedAddress = s.SatStaticTranslatedAddress
	o.SatStaticBiDirectional = s.SatStaticBiDirectional
	o.DatAddress = s.DatAddress
	o.DatPort = s.DatPort
	o.Disabled = s.Disabled
	o.Targets = s.Targets
	o.NegateTarget = s.NegateTarget
	o.Tags = s.Tags
	o.DatType = s.DatType
	o.DatDynamicDistribution = s.DatDynamicDistribution
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

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                 o.Name,
		Description:          o.Description,
		Type:                 o.Type,
		SourceZones:          util.MemToStr(o.SourceZones),
		DestinationZone:      o.DestinationZone,
		ToInterface:          o.ToInterface,
		Service:              o.Service,
		SourceAddresses:      util.MemToStr(o.SourceAddresses),
		DestinationAddresses: util.MemToStr(o.DestinationAddresses),
		Disabled:             util.AsBool(o.Disabled),
		Tags:                 util.MemToStr(o.Tags),
	}

	if o.Sat == nil {
		ans.SatType = None
	} else {
		switch {
		case o.Sat.Diap != nil:
			ans.SatType = DynamicIpAndPort
			if o.Sat.Diap.InterfaceAddress != nil {
				ans.SatAddressType = InterfaceAddress
				ans.SatInterface = o.Sat.Diap.InterfaceAddress.Interface
				ans.SatIpAddress = o.Sat.Diap.InterfaceAddress.Ip
			} else {
				ans.SatAddressType = TranslatedAddress
				ans.SatTranslatedAddresses = util.MemToStr(o.Sat.Diap.TranslatedAddress)
			}
		case o.Sat.Di != nil:
			ans.SatType = DynamicIp
			ans.SatTranslatedAddresses = util.MemToStr(o.Sat.Di.TranslatedAddress)
			if o.Sat.Di.Fallback == nil {
				ans.SatFallbackType = None
			} else if o.Sat.Di.Fallback.TranslatedAddress != nil {
				ans.SatFallbackType = TranslatedAddress
				ans.SatFallbackTranslatedAddresses = util.MemToStr(o.Sat.Di.Fallback.TranslatedAddress)
			} else if o.Sat.Di.Fallback.InterfaceAddress != nil {
				ans.SatFallbackType = InterfaceAddress
				ans.SatFallbackInterface = o.Sat.Di.Fallback.InterfaceAddress.Interface
				if o.Sat.Di.Fallback.InterfaceAddress.Ip != "" {
					ans.SatFallbackIpType = Ip
					ans.SatFallbackIpAddress = o.Sat.Di.Fallback.InterfaceAddress.Ip
				} else if o.Sat.Di.Fallback.InterfaceAddress.FloatingIp != "" {
					ans.SatFallbackIpType = FloatingIp
					ans.SatFallbackIpAddress = o.Sat.Di.Fallback.InterfaceAddress.FloatingIp
				}
			}
		case o.Sat.Static != nil:
			ans.SatType = StaticIp
			ans.SatStaticTranslatedAddress = o.Sat.Static.Address
			ans.SatStaticBiDirectional = util.AsBool(o.Sat.Static.BiDirectional)
		}
	}

	if o.Dat != nil {
		ans.DatType = DatTypeStatic
		ans.DatAddress = o.Dat.Address
		ans.DatPort = o.Dat.Port
	}

	if o.Target != nil {
		ans.Targets = util.VsysEntToMap(o.Target.Targets)
		ans.NegateTarget = util.AsBool(o.Target.NegateTarget)
	}

	return ans
}

type entry_v1 struct {
	XMLName              xml.Name         `xml:"entry"`
	Name                 string           `xml:"name,attr"`
	Description          string           `xml:"description"`
	Type                 string           `xml:"nat-type"`
	SourceZones          *util.MemberType `xml:"from"`
	DestinationZone      string           `xml:"to>member"`
	ToInterface          string           `xml:"to-interface"`
	Service              string           `xml:"service"`
	SourceAddresses      *util.MemberType `xml:"source"`
	DestinationAddresses *util.MemberType `xml:"destination"`
	Sat                  *srcXlate        `xml:"source-translation"`
	Dat                  *dstXlate        `xml:"destination-translation"`
	Disabled             string           `xml:"disabled"`
	Target               *targetInfo      `xml:"target"`
	Tags                 *util.MemberType `xml:"tag"`
}

type dstXlate struct {
	Address      string `xml:"translated-address,omitempty"`
	Port         int    `xml:"translated-port,omitempty"`
	Distribution string `xml:"distribution,omitempty"`
}

type srcXlate struct {
	Diap   *srcXlateDiap   `xml:"dynamic-ip-and-port"`
	Di     *srcXlateDi     `xml:"dynamic-ip"`
	Static *srcXlateStatic `xml:"static-ip"`
}

type srcXlateDiap struct {
	TranslatedAddress *util.MemberType `xml:"translated-address"`
	InterfaceAddress  *srcXlateDiapIa  `xml:"interface-address"`
}

type srcXlateDiapIa struct {
	Interface string `xml:"interface"`
	Ip        string `xml:"ip,omitempty"`
}

type srcXlateDi struct {
	TranslatedAddress *util.MemberType `xml:"translated-address"`
	Fallback          *fallback        `xml:"fallback"`
}

type fallback struct {
	TranslatedAddress *util.MemberType `xml:"translated-address"`
	InterfaceAddress  *fallbackIface   `xml:"interface-address"`
}

type fallbackIface struct {
	Ip         string `xml:"ip,omitempty"`
	Interface  string `xml:"interface,omitempty"`
	FloatingIp string `xml:"floating-ip,omitempty"`
}

type srcXlateStatic struct {
	Address       string `xml:"translated-address"`
	BiDirectional string `xml:"bi-directional"`
}

type targetInfo struct {
	Targets      *util.VsysEntryType `xml:"devices"`
	NegateTarget string              `xml:"negate,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                 e.Name,
		Description:          e.Description,
		Type:                 e.Type,
		SourceZones:          util.StrToMem(e.SourceZones),
		DestinationZone:      e.DestinationZone,
		ToInterface:          e.ToInterface,
		Service:              e.Service,
		SourceAddresses:      util.StrToMem(e.SourceAddresses),
		DestinationAddresses: util.StrToMem(e.DestinationAddresses),
		Disabled:             util.YesNo(e.Disabled),
		Tags:                 util.StrToMem(e.Tags),
	}

	var sv *srcXlate
	switch e.SatType {
	case DynamicIpAndPort:
		sv = &srcXlate{
			Diap: &srcXlateDiap{},
		}
		switch e.SatAddressType {
		case TranslatedAddress:
			sv.Diap.TranslatedAddress = util.StrToMem(e.SatTranslatedAddresses)
		case InterfaceAddress:
			sv.Diap.InterfaceAddress = &srcXlateDiapIa{
				Interface: e.SatInterface,
				Ip:        e.SatIpAddress,
			}
		}
	case DynamicIp:
		sv = &srcXlate{
			Di: &srcXlateDi{
				TranslatedAddress: util.StrToMem(e.SatTranslatedAddresses),
			},
		}
		switch e.SatFallbackType {
		case InterfaceAddress:
			sv.Di.Fallback = &fallback{
				InterfaceAddress: &fallbackIface{
					Interface: e.SatFallbackInterface,
				},
			}
			switch e.SatFallbackIpType {
			case Ip:
				sv.Di.Fallback.InterfaceAddress.Ip = e.SatFallbackIpAddress
			case FloatingIp:
				sv.Di.Fallback.InterfaceAddress.FloatingIp = e.SatFallbackIpAddress
			}
		case TranslatedAddress:
			sv.Di.Fallback = &fallback{TranslatedAddress: util.StrToMem(e.SatFallbackTranslatedAddresses)}
		}
	case StaticIp:
		sv = &srcXlate{
			Static: &srcXlateStatic{
				e.SatStaticTranslatedAddress,
				util.YesNo(e.SatStaticBiDirectional),
			},
		}
	}
	ans.Sat = sv

	if e.DatType == DatTypeStatic {
		if e.DatAddress != "" || e.DatPort != 0 {
			ans.Dat = &dstXlate{
				e.DatAddress,
				e.DatPort,
				"",
			}
		}
	}

	if len(e.Targets) != 0 || e.NegateTarget {
		ans.Target = &targetInfo{
			Targets:      util.MapToVsysEnt(e.Targets),
			NegateTarget: util.YesNo(e.NegateTarget),
		}
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

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:                 o.Name,
		Description:          o.Description,
		Type:                 o.Type,
		SourceZones:          util.MemToStr(o.SourceZones),
		DestinationZone:      o.DestinationZone,
		ToInterface:          o.ToInterface,
		Service:              o.Service,
		SourceAddresses:      util.MemToStr(o.SourceAddresses),
		DestinationAddresses: util.MemToStr(o.DestinationAddresses),
		Disabled:             util.AsBool(o.Disabled),
		Tags:                 util.MemToStr(o.Tags),
	}

	if o.Sat == nil {
		ans.SatType = None
	} else {
		switch {
		case o.Sat.Diap != nil:
			ans.SatType = DynamicIpAndPort
			if o.Sat.Diap.InterfaceAddress != nil {
				ans.SatAddressType = InterfaceAddress
				ans.SatInterface = o.Sat.Diap.InterfaceAddress.Interface
				ans.SatIpAddress = o.Sat.Diap.InterfaceAddress.Ip
			} else {
				ans.SatAddressType = TranslatedAddress
				ans.SatTranslatedAddresses = util.MemToStr(o.Sat.Diap.TranslatedAddress)
			}
		case o.Sat.Di != nil:
			ans.SatType = DynamicIp
			ans.SatTranslatedAddresses = util.MemToStr(o.Sat.Di.TranslatedAddress)
			if o.Sat.Di.Fallback == nil {
				ans.SatFallbackType = None
			} else if o.Sat.Di.Fallback.TranslatedAddress != nil {
				ans.SatFallbackType = TranslatedAddress
				ans.SatFallbackTranslatedAddresses = util.MemToStr(o.Sat.Di.Fallback.TranslatedAddress)
			} else if o.Sat.Di.Fallback.InterfaceAddress != nil {
				ans.SatFallbackType = InterfaceAddress
				ans.SatFallbackInterface = o.Sat.Di.Fallback.InterfaceAddress.Interface
				if o.Sat.Di.Fallback.InterfaceAddress.Ip != "" {
					ans.SatFallbackIpType = Ip
					ans.SatFallbackIpAddress = o.Sat.Di.Fallback.InterfaceAddress.Ip
				} else if o.Sat.Di.Fallback.InterfaceAddress.FloatingIp != "" {
					ans.SatFallbackIpType = FloatingIp
					ans.SatFallbackIpAddress = o.Sat.Di.Fallback.InterfaceAddress.FloatingIp
				}
			}
		case o.Sat.Static != nil:
			ans.SatType = StaticIp
			ans.SatStaticTranslatedAddress = o.Sat.Static.Address
			ans.SatStaticBiDirectional = util.AsBool(o.Sat.Static.BiDirectional)
		}
	}

	if o.Dat != nil {
		ans.DatType = DatTypeStatic
		ans.DatAddress = o.Dat.Address
		ans.DatPort = o.Dat.Port
	}

	if o.DatDynamic != nil {
		ans.DatType = DatTypeDynamic
		ans.DatAddress = o.DatDynamic.Address
		ans.DatPort = o.DatDynamic.Port
		ans.DatDynamicDistribution = o.DatDynamic.Distribution
	}

	if o.Target != nil {
		ans.Targets = util.VsysEntToMap(o.Target.Targets)
		ans.NegateTarget = util.AsBool(o.Target.NegateTarget)
	}

	return ans
}

type entry_v2 struct {
	XMLName              xml.Name         `xml:"entry"`
	Name                 string           `xml:"name,attr"`
	Description          string           `xml:"description"`
	Type                 string           `xml:"nat-type"`
	SourceZones          *util.MemberType `xml:"from"`
	DestinationZone      string           `xml:"to>member"`
	ToInterface          string           `xml:"to-interface"`
	Service              string           `xml:"service"`
	SourceAddresses      *util.MemberType `xml:"source"`
	DestinationAddresses *util.MemberType `xml:"destination"`
	Sat                  *srcXlate        `xml:"source-translation"`
	Dat                  *dstXlate        `xml:"destination-translation"`
	DatDynamic           *dstXlate        `xml:"dynamic-destination-translation"`
	Disabled             string           `xml:"disabled"`
	Target               *targetInfo      `xml:"target"`
	Tags                 *util.MemberType `xml:"tag"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                 e.Name,
		Description:          e.Description,
		Type:                 e.Type,
		SourceZones:          util.StrToMem(e.SourceZones),
		DestinationZone:      e.DestinationZone,
		ToInterface:          e.ToInterface,
		Service:              e.Service,
		SourceAddresses:      util.StrToMem(e.SourceAddresses),
		DestinationAddresses: util.StrToMem(e.DestinationAddresses),
		Disabled:             util.YesNo(e.Disabled),
		Tags:                 util.StrToMem(e.Tags),
	}

	var sv *srcXlate
	switch e.SatType {
	case DynamicIpAndPort:
		sv = &srcXlate{
			Diap: &srcXlateDiap{},
		}
		switch e.SatAddressType {
		case TranslatedAddress:
			sv.Diap.TranslatedAddress = util.StrToMem(e.SatTranslatedAddresses)
		case InterfaceAddress:
			sv.Diap.InterfaceAddress = &srcXlateDiapIa{
				Interface: e.SatInterface,
				Ip:        e.SatIpAddress,
			}
		}
	case DynamicIp:
		sv = &srcXlate{
			Di: &srcXlateDi{
				TranslatedAddress: util.StrToMem(e.SatTranslatedAddresses),
			},
		}
		switch e.SatFallbackType {
		case InterfaceAddress:
			sv.Di.Fallback = &fallback{
				InterfaceAddress: &fallbackIface{
					Interface: e.SatFallbackInterface,
				},
			}
			switch e.SatFallbackIpType {
			case Ip:
				sv.Di.Fallback.InterfaceAddress.Ip = e.SatFallbackIpAddress
			case FloatingIp:
				sv.Di.Fallback.InterfaceAddress.FloatingIp = e.SatFallbackIpAddress
			}
		case TranslatedAddress:
			sv.Di.Fallback = &fallback{TranslatedAddress: util.StrToMem(e.SatFallbackTranslatedAddresses)}
		}
	case StaticIp:
		sv = &srcXlate{
			Static: &srcXlateStatic{
				e.SatStaticTranslatedAddress,
				util.YesNo(e.SatStaticBiDirectional),
			},
		}
	}
	ans.Sat = sv

	if e.DatType == DatTypeStatic {
		if e.DatAddress != "" || e.DatPort != 0 {
			ans.Dat = &dstXlate{
				e.DatAddress,
				e.DatPort,
				"",
			}
		}
	} else if e.DatType == DatTypeDynamic {
		if e.DatAddress != "" || e.DatPort != 0 || e.DatDynamicDistribution != "" {
			ans.DatDynamic = &dstXlate{
				e.DatAddress,
				e.DatPort,
				e.DatDynamicDistribution,
			}
		}
	}

	if len(e.Targets) != 0 || e.NegateTarget {
		ans.Target = &targetInfo{
			Targets:      util.MapToVsysEnt(e.Targets),
			NegateTarget: util.YesNo(e.NegateTarget),
		}
	}

	return ans
}
