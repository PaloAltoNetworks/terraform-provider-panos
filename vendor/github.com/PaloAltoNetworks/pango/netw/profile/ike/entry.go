package ike

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


const (
    EncryptionDes = "des"
    Encryption3des = "3des"
    EncryptionAes128 = "aes-128-cbc"
    EncryptionAes192 = "aes-192-cbc"
    EncryptionAes256 = "aes-256-cbc"
)

const (
    TimeSeconds = "seconds"
    TimeMinutes = "minutes"
    TimeHours = "hours"
    TimeDays = "days"
)

// Entry is a normalized, version independent representation of an interface
// management profile.
type Entry struct {
    Name string
    DhGroup []string
    Authentication []string
    Encryption []string
    LifetimeType string
    LifetimeValue int
    AuthenticationMultiple int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.DhGroup = s.DhGroup
    o.Authentication = s.Authentication
    o.Encryption = s.Encryption
    o.LifetimeType = s.LifetimeType
    o.LifetimeValue = s.LifetimeValue
    o.AuthenticationMultiple = s.AuthenticationMultiple
}

// SpecifyEncryption takes normalizes encryption values and changes them to the
// version specific values PAN-OS will be expecting.
//
// Param v should be 1 if you're running against PAN-OS 6.1; 2 if you're
// running against 7.0 or later.
func (o *Entry) SpecifyEncryption(v int) {
    if len(o.Encryption) == 0 {
        return
    }

    nv := make([]string, len(o.Encryption))
    switch v {
    case 2:
        for i := range o.Encryption {
            switch o.Encryption[i] {
            case EncryptionDes:
                nv[i] = "des"
            case Encryption3des:
                nv[i] = "3des"
            case EncryptionAes128:
                nv[i] = "aes-128-cbc"
            case EncryptionAes192:
                nv[i] = "aes-192-cbc"
            case EncryptionAes256:
                nv[i] = "aes-256-cbc"
            default:
                nv[i] = o.Encryption[i]
            }
        }
    case 1:
        for i := range o.Encryption {
            switch o.Encryption[i] {
            case Encryption3des:
                nv[i] = "3des"
            case EncryptionAes128:
                nv[i] = "aes128"
            case EncryptionAes192:
                nv[i] = "aes192"
            case EncryptionAes256:
                nv[i] = "aes256"
            default:
                nv[i] = o.Encryption[i]
            }
        }
    default:
        copy(nv, o.Encryption)
    }

    o.Encryption = nv
}

// NormalizeEncryption normalizes the fields in o.Encryption.
func (o *Entry) NormalizeEncryption() {
    if len(o.Encryption) == 0 {
        return
    }

    nv := make([]string, len(o.Encryption))
    for i := range o.Encryption {
        switch o.Encryption[i] {
        case "des":
            nv[i] = EncryptionDes
        case "3des":
            nv[i] = Encryption3des
        case "aes-128-cbc", "aes128":
            nv[i] = EncryptionAes128
        case "aes-192-cbc", "aes192":
            nv[i] = EncryptionAes192
        case "aes-256-cbc", "aes256":
            nv[i] = EncryptionAes256
        default:
            nv[i] = o.Encryption[i]
        }
    }

    o.Encryption = nv
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
        DhGroup: util.MemToStr(o.Answer.DhGroup),
        Authentication: util.MemToStr(o.Answer.Authentication),
        Encryption: util.MemToStr(o.Answer.Encryption),
    }

    if o.Answer.Lifetime != nil {
        if o.Answer.Lifetime.Seconds != 0 {
            ans.LifetimeType = TimeSeconds
            ans.LifetimeValue = o.Answer.Lifetime.Seconds
        } else if o.Answer.Lifetime.Minutes != 0 {
            ans.LifetimeType = TimeMinutes
            ans.LifetimeValue = o.Answer.Lifetime.Minutes
        } else if o.Answer.Lifetime.Hours != 0 {
            ans.LifetimeType = TimeHours
            ans.LifetimeValue = o.Answer.Lifetime.Hours
        } else if o.Answer.Lifetime.Days != 0 {
            ans.LifetimeType = TimeDays
            ans.LifetimeValue = o.Answer.Lifetime.Days
        }
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        DhGroup: util.MemToStr(o.Answer.DhGroup),
        Authentication: util.MemToStr(o.Answer.Authentication),
        Encryption: util.MemToStr(o.Answer.Encryption),
        AuthenticationMultiple: o.Answer.AuthenticationMultiple,
    }

    if o.Answer.Lifetime != nil {
        if o.Answer.Lifetime.Seconds != 0 {
            ans.LifetimeType = TimeSeconds
            ans.LifetimeValue = o.Answer.Lifetime.Seconds
        } else if o.Answer.Lifetime.Minutes != 0 {
            ans.LifetimeType = TimeMinutes
            ans.LifetimeValue = o.Answer.Lifetime.Minutes
        } else if o.Answer.Lifetime.Hours != 0 {
            ans.LifetimeType = TimeHours
            ans.LifetimeValue = o.Answer.Lifetime.Hours
        } else if o.Answer.Lifetime.Days != 0 {
            ans.LifetimeType = TimeDays
            ans.LifetimeValue = o.Answer.Lifetime.Days
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    DhGroup *util.MemberType `xml:"dh-group"`
    Authentication *util.MemberType `xml:"hash"`
    Encryption *util.MemberType `xml:"encryption"`
    Lifetime *lifetimeType `xml:"lifetime"`
}

type lifetimeType struct {
    Seconds int `xml:"seconds,omitempty"`
    Minutes int `xml:"minutes,omitempty"`
    Hours int `xml:"hours,omitempty"`
    Days int `xml:"days,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        DhGroup: util.StrToMem(e.DhGroup),
        Authentication: util.StrToMem(e.Authentication),
        Encryption: util.StrToMem(e.Encryption),
    }

    switch e.LifetimeType {
    case TimeSeconds:
        ans.Lifetime = &lifetimeType{Seconds: e.LifetimeValue}
    case TimeMinutes:
        ans.Lifetime = &lifetimeType{Minutes: e.LifetimeValue}
    case TimeHours:
        ans.Lifetime = &lifetimeType{Hours: e.LifetimeValue}
    case TimeDays:
        ans.Lifetime = &lifetimeType{Days: e.LifetimeValue}
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    DhGroup *util.MemberType `xml:"dh-group"`
    Authentication *util.MemberType `xml:"hash"`
    Encryption *util.MemberType `xml:"encryption"`
    AuthenticationMultiple int `xml:"authentication-multiple,omitempty"`
    Lifetime *lifetimeType `xml:"lifetime"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        DhGroup: util.StrToMem(e.DhGroup),
        Authentication: util.StrToMem(e.Authentication),
        Encryption: util.StrToMem(e.Encryption),
        AuthenticationMultiple: e.AuthenticationMultiple,
    }

    switch e.LifetimeType {
    case TimeSeconds:
        ans.Lifetime = &lifetimeType{Seconds: e.LifetimeValue}
    case TimeMinutes:
        ans.Lifetime = &lifetimeType{Minutes: e.LifetimeValue}
    case TimeHours:
        ans.Lifetime = &lifetimeType{Hours: e.LifetimeValue}
    case TimeDays:
        ans.Lifetime = &lifetimeType{Days: e.LifetimeValue}
    }

    return ans
}
