package zone

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a zone.
type Entry struct {
	Name         string
	Mode         string
	Interfaces   []string // unordered
	ZoneProfile  string
	LogSetting   string
	EnableUserId bool
	IncludeAcls  []string // unordered
	ExcludeAcls  []string // unordered
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Mode = s.Mode
	o.Interfaces = s.Interfaces
	o.ZoneProfile = s.ZoneProfile
	o.LogSetting = s.LogSetting
	o.EnableUserId = s.EnableUserId
	o.IncludeAcls = s.IncludeAcls
	o.ExcludeAcls = s.ExcludeAcls
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
		Name:         o.Name,
		ZoneProfile:  o.Profile,
		LogSetting:   o.LogSetting,
		EnableUserId: util.AsBool(o.EnableUserId),
	}
	if o.L3 != nil {
		ans.Mode = ModeL3
		ans.Interfaces = o.L3.Interfaces
	} else if o.L2 != nil {
		ans.Mode = ModeL2
		ans.Interfaces = o.L2.Interfaces
	} else if o.VWire != nil {
		ans.Mode = ModeVirtualWire
		ans.Interfaces = o.VWire.Interfaces
	} else if o.Tap != nil {
		ans.Mode = ModeTap
		ans.Interfaces = o.Tap.Interfaces
	} else if o.External != nil {
		ans.Mode = ModeExternal
		ans.Interfaces = o.External.Interfaces
	}
	if o.IncludeAcls != nil {
		ans.IncludeAcls = o.IncludeAcls.Acls
	}
	if o.ExcludeAcls != nil {
		ans.ExcludeAcls = o.ExcludeAcls.Acls
	}

	return ans
}

type entry_v1 struct {
	XMLName      xml.Name           `xml:"entry"`
	Name         string             `xml:"name,attr"`
	L3           *zoneInterfaceList `xml:"network>layer3"`
	L2           *zoneInterfaceList `xml:"network>layer2"`
	VWire        *zoneInterfaceList `xml:"network>virtual-wire"`
	Tap          *zoneInterfaceList `xml:"network>tap"`
	External     *zoneInterfaceList `xml:"network>external"`
	Profile      string             `xml:"network>zone-protection-profile,omitempty"`
	LogSetting   string             `xml:"network>log-setting,omitempty"`
	EnableUserId string             `xml:"enable-user-identification"`
	IncludeAcls  *aclList           `xml:"user-acl>include-list"`
	ExcludeAcls  *aclList           `xml:"user-acl>exclude-list"`
}

type zoneInterfaceList struct {
	Interfaces []string `xml:"member"`
}

type aclList struct {
	Acls []string `xml:"member"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:         e.Name,
		Profile:      e.ZoneProfile,
		LogSetting:   e.LogSetting,
		EnableUserId: util.YesNo(e.EnableUserId),
	}
	il := &zoneInterfaceList{e.Interfaces}
	switch e.Mode {
	case ModeL2:
		ans.L2 = il
	case ModeL3:
		ans.L3 = il
	case ModeVirtualWire:
		ans.VWire = il
	case ModeTap:
		ans.Tap = il
	case ModeExternal:
		ans.External = il
	}
	if len(e.IncludeAcls) > 0 {
		inu := &aclList{e.IncludeAcls}
		ans.IncludeAcls = inu
	}
	if len(e.ExcludeAcls) > 0 {
		exu := &aclList{e.ExcludeAcls}
		ans.ExcludeAcls = exu
	}

	return ans
}
