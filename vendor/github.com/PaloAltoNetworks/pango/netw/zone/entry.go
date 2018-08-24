package zone

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Entry is a normalized, version independent representation of a zone.
type Entry struct {
    Name string
    Mode string
    Interfaces []string
    ZoneProfile string
    LogSetting string
    EnableUserId bool
    IncludeAcls []string
    ExcludeAcls []string
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

type normalizer interface {
    Normalize() Entry
}

type container_v1 struct {
    Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        ZoneProfile: o.Answer.Profile,
        LogSetting: o.Answer.LogSetting,
        EnableUserId: util.AsBool(o.Answer.EnableUserId),
    }
    if o.Answer.L3 != nil {
        ans.Mode = "layer3"
        ans.Interfaces = o.Answer.L3.Interfaces
    } else if o.Answer.L2 != nil {
        ans.Mode = "layer2"
        ans.Interfaces = o.Answer.L2.Interfaces
    } else if o.Answer.VWire != nil {
        ans.Mode = "virtual-wire"
        ans.Interfaces = o.Answer.VWire.Interfaces
    } else if o.Answer.Tap != nil {
        ans.Mode = "tap"
        ans.Interfaces = o.Answer.Tap.Interfaces
    } else if o.Answer.External != nil {
        ans.Mode = "external"
        ans.Interfaces = o.Answer.External.Interfaces
    }
    if o.Answer.IncludeAcls != nil {
        ans.IncludeAcls = o.Answer.IncludeAcls.Acls
    }
    if o.Answer.ExcludeAcls != nil {
        ans.ExcludeAcls = o.Answer.ExcludeAcls.Acls
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    L3 *zoneInterfaceList `xml:"network>layer3"`
    L2 *zoneInterfaceList `xml:"network>layer2"`
    VWire *zoneInterfaceList `xml:"network>virtual-wire"`
    Tap *zoneInterfaceList `xml:"network>tap"`
    External *zoneInterfaceList `xml:"network>external"`
    Profile string `xml:"network>zone-protection-profile,omitempty"`
    LogSetting string `xml:"network>log-setting,omitempty"`
    EnableUserId string `xml:"enable-user-identification"`
    IncludeAcls *aclList `xml:"user-acl>include-list"`
    ExcludeAcls *aclList `xml:"user-acl>exclude-list"`
}

type zoneInterfaceList struct {
    Interfaces []string `xml:"member"`
}

type aclList struct {
    Acls []string `xml:"member"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        Profile: e.ZoneProfile,
        LogSetting: e.LogSetting,
        EnableUserId: util.YesNo(e.EnableUserId),
    }
    il := &zoneInterfaceList{e.Interfaces}
    switch e.Mode {
    case "layer2":
        ans.L2 = il
    case "layer3":
        ans.L3 = il
    case "virtual-wire":
        ans.VWire = il
    case "tap":
        ans.Tap = il
    case "external":
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
