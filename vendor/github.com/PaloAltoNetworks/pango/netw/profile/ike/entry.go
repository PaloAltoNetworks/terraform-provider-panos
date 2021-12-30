package ike

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of
// IKE crypto profile.
type Entry struct {
	Name                   string
	DhGroup                []string
	Authentication         []string
	Encryption             []string
	LifetimeType           string
	LifetimeValue          int
	AuthenticationMultiple int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.DhGroup = s.DhGroup
	if s.DhGroup == nil {
		o.DhGroup = nil
	} else {
		o.DhGroup = make([]string, len(s.DhGroup))
		copy(o.DhGroup, s.DhGroup)
	}
	if s.Authentication == nil {
		o.Authentication = nil
	} else {
		o.Authentication = make([]string, len(s.Authentication))
		copy(o.Authentication, s.Authentication)
	}
	if s.Encryption == nil {
		o.Encryption = nil
	} else {
		o.Encryption = make([]string, len(s.Encryption))
		copy(o.Encryption, s.Encryption)
	}
	o.LifetimeType = s.LifetimeType
	o.LifetimeValue = s.LifetimeValue
	o.AuthenticationMultiple = s.AuthenticationMultiple
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

func specifyEncryption(vals []string, v int) []string {
	if vals == nil {
		return nil
	}

	nv := make([]string, len(vals))
	switch v {
	case 2:
		for i := range vals {
			switch vals[i] {
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
			case EncryptionAes128Gcm:
				nv[i] = "aes-128-gcm"
			case EncryptionAes256Gcm:
				nv[i] = "aes-256-gcm"
			default:
				nv[i] = vals[i]
			}
		}
	case 1:
		for i := range vals {
			switch vals[i] {
			case Encryption3des:
				nv[i] = "3des"
			case EncryptionAes128:
				nv[i] = "aes128"
			case EncryptionAes192:
				nv[i] = "aes192"
			case EncryptionAes256:
				nv[i] = "aes256"
			default:
				nv[i] = vals[i]
			}
		}
	default:
		copy(nv, vals)
	}

	return nv
}

func normalizeEncryption(v []string) []string {
	if v == nil {
		return nil
	}

	ans := make([]string, len(v))
	for i := range v {
		switch v[i] {
		case "des":
			ans[i] = EncryptionDes
		case "3des":
			ans[i] = Encryption3des
		case "aes-128-cbc", "aes128":
			ans[i] = EncryptionAes128
		case "aes-192-cbc", "aes192":
			ans[i] = EncryptionAes192
		case "aes-256-cbc", "aes256":
			ans[i] = EncryptionAes256
		case "aes-128-gcm":
			ans[i] = EncryptionAes128Gcm
		case "aes-256-gcm":
			ans[i] = EncryptionAes256Gcm
		default:
			ans[i] = v[i]
		}
	}

	return ans
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
		Name:           o.Name,
		DhGroup:        util.MemToStr(o.DhGroup),
		Authentication: util.MemToStr(o.Authentication),
		Encryption:     normalizeEncryption(util.MemToStr(o.Encryption)),
	}

	if o.Lifetime != nil {
		if o.Lifetime.Seconds != 0 {
			ans.LifetimeType = TimeSeconds
			ans.LifetimeValue = o.Lifetime.Seconds
		} else if o.Lifetime.Minutes != 0 {
			ans.LifetimeType = TimeMinutes
			ans.LifetimeValue = o.Lifetime.Minutes
		} else if o.Lifetime.Hours != 0 {
			ans.LifetimeType = TimeHours
			ans.LifetimeValue = o.Lifetime.Hours
		} else if o.Lifetime.Days != 0 {
			ans.LifetimeType = TimeDays
			ans.LifetimeValue = o.Lifetime.Days
		}
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
		Name:                   o.Name,
		DhGroup:                util.MemToStr(o.DhGroup),
		Authentication:         util.MemToStr(o.Authentication),
		Encryption:             normalizeEncryption(util.MemToStr(o.Encryption)),
		AuthenticationMultiple: o.AuthenticationMultiple,
	}

	if o.Lifetime != nil {
		if o.Lifetime.Seconds != 0 {
			ans.LifetimeType = TimeSeconds
			ans.LifetimeValue = o.Lifetime.Seconds
		} else if o.Lifetime.Minutes != 0 {
			ans.LifetimeType = TimeMinutes
			ans.LifetimeValue = o.Lifetime.Minutes
		} else if o.Lifetime.Hours != 0 {
			ans.LifetimeType = TimeHours
			ans.LifetimeValue = o.Lifetime.Hours
		} else if o.Lifetime.Days != 0 {
			ans.LifetimeType = TimeDays
			ans.LifetimeValue = o.Lifetime.Days
		}
	}

	// Normalize the encryption values.
	for i := range ans.Encryption {
		switch ans.Encryption[i] {
		case "des":
			ans.Encryption[i] = EncryptionDes
		case "3des":
			ans.Encryption[i] = Encryption3des
		case "aes-128-cbc", "aes128":
			ans.Encryption[i] = EncryptionAes128
		case "aes-192-cbc", "aes192":
			ans.Encryption[i] = EncryptionAes192
		case "aes-256-cbc", "aes256":
			ans.Encryption[i] = EncryptionAes256
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName        xml.Name         `xml:"entry"`
	Name           string           `xml:"name,attr"`
	DhGroup        *util.MemberType `xml:"dh-group"`
	Authentication *util.MemberType `xml:"hash"`
	Encryption     *util.MemberType `xml:"encryption"`
	Lifetime       *lifetimeType    `xml:"lifetime"`
}

type lifetimeType struct {
	Seconds int `xml:"seconds,omitempty"`
	Minutes int `xml:"minutes,omitempty"`
	Hours   int `xml:"hours,omitempty"`
	Days    int `xml:"days,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:           e.Name,
		DhGroup:        util.StrToMem(e.DhGroup),
		Authentication: util.StrToMem(e.Authentication),
		Encryption:     util.StrToMem(specifyEncryption(e.Encryption, 1)),
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
	XMLName                xml.Name         `xml:"entry"`
	Name                   string           `xml:"name,attr"`
	DhGroup                *util.MemberType `xml:"dh-group"`
	Authentication         *util.MemberType `xml:"hash"`
	Encryption             *util.MemberType `xml:"encryption"`
	AuthenticationMultiple int              `xml:"authentication-multiple,omitempty"`
	Lifetime               *lifetimeType    `xml:"lifetime"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                   e.Name,
		DhGroup:                util.StrToMem(e.DhGroup),
		Authentication:         util.StrToMem(e.Authentication),
		Encryption:             util.StrToMem(specifyEncryption(e.Encryption, 2)),
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
