package ipsec

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

const (
	ProtocolEsp = "esp"
	ProtocolAh  = "ah"
)

const (
	EncryptionDes       = "des"
	Encryption3des      = "3des"
	EncryptionAes128    = "aes-128-cbc"
	EncryptionAes192    = "aes-192-cbc"
	EncryptionAes256    = "aes-256-cbc"
	EncryptionAes128Gcm = "aes-128-gcm"
	EncryptionAes256Gcm = "aes-256-gcm"
	EncryptionNull      = "null"
)

const (
	TimeSeconds = "seconds"
	TimeMinutes = "minutes"
	TimeHours   = "hours"
	TimeDays    = "days"
)

const (
	SizeKb = "kb"
	SizeMb = "mb"
	SizeGb = "gb"
	SizeTb = "tb"
)

// Entry is a normalized, version independent representation of an interface
// management profile.
type Entry struct {
	Name           string
	Protocol       string
	Encryption     []string
	Authentication []string
	DhGroup        string
	LifetimeType   string
	LifetimeValue  int
	LifesizeType   string
	LifesizeValue  int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Protocol = s.Protocol
	o.Encryption = s.Encryption
	o.Authentication = s.Authentication
	o.DhGroup = s.DhGroup
	o.LifetimeType = s.LifetimeType
	o.LifetimeValue = s.LifetimeValue
	o.LifesizeType = s.LifesizeType
	o.LifesizeValue = s.LifesizeValue
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
			case EncryptionAes128Gcm:
				nv[i] = "aes-128-gcm"
			case EncryptionAes256Gcm:
				nv[i] = "aes-256-gcm"
			case EncryptionNull:
				nv[i] = "null"
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
			case EncryptionNull:
				nv[i] = "null"
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
		case "aes-128-gcm":
			nv[i] = EncryptionAes128Gcm
		case "aes-256-gcm":
			nv[i] = EncryptionAes256Gcm
		case "null":
			nv[i] = EncryptionNull
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
		Name:    o.Answer.Name,
		DhGroup: o.Answer.DhGroup,
	}

	if o.Answer.Esp != nil {
		ans.Protocol = ProtocolEsp
		ans.Encryption = util.MemToStr(o.Answer.Esp.Encryption)
		ans.Authentication = util.MemToStr(o.Answer.Esp.Authentication)
	} else if o.Answer.Ah != nil {
		ans.Protocol = ProtocolAh
		ans.Authentication = util.MemToStr(o.Answer.Ah.Authentication)
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

	if o.Answer.Lifesize != nil {
		if o.Answer.Lifesize.Kb != 0 {
			ans.LifesizeType = SizeKb
			ans.LifesizeValue = o.Answer.Lifesize.Kb
		} else if o.Answer.Lifesize.Mb != 0 {
			ans.LifesizeType = SizeMb
			ans.LifesizeValue = o.Answer.Lifesize.Mb
		} else if o.Answer.Lifesize.Gb != 0 {
			ans.LifesizeType = SizeGb
			ans.LifesizeValue = o.Answer.Lifesize.Gb
		} else if o.Answer.Lifesize.Tb != 0 {
			ans.LifesizeType = SizeTb
			ans.LifesizeValue = o.Answer.Lifesize.Tb
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName  xml.Name      `xml:"entry"`
	Name     string        `xml:"name,attr"`
	Esp      *espType      `xml:"esp"`
	Ah       *ahType       `xml:"ah"`
	DhGroup  string        `xml:"dh-group,omitempty"`
	Lifetime *lifetimeType `xml:"lifetime"`
	Lifesize *lifesizeType `xml:"lifesize"`
}

type espType struct {
	Encryption     *util.MemberType `xml:"encryption"`
	Authentication *util.MemberType `xml:"authentication"`
}

type ahType struct {
	Authentication *util.MemberType `xml:"authentication"`
}

type lifetimeType struct {
	Seconds int `xml:"seconds,omitempty"`
	Minutes int `xml:"minutes,omitempty"`
	Hours   int `xml:"hours,omitempty"`
	Days    int `xml:"days,omitempty"`
}

type lifesizeType struct {
	Kb int `xml:"kb,omitempty"`
	Mb int `xml:"mb,omitempty"`
	Gb int `xml:"gb,omitempty"`
	Tb int `xml:"tb,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:    e.Name,
		DhGroup: e.DhGroup,
	}

	switch e.Protocol {
	case ProtocolEsp:
		ans.Esp = &espType{
			Encryption:     util.StrToMem(e.Encryption),
			Authentication: util.StrToMem(e.Authentication),
		}
	case ProtocolAh:
		ans.Ah = &ahType{util.StrToMem(e.Authentication)}
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

	switch e.LifesizeType {
	case SizeKb:
		ans.Lifesize = &lifesizeType{Kb: e.LifesizeValue}
	case SizeMb:
		ans.Lifesize = &lifesizeType{Mb: e.LifesizeValue}
	case SizeGb:
		ans.Lifesize = &lifesizeType{Gb: e.LifesizeValue}
	case SizeTb:
		ans.Lifesize = &lifesizeType{Tb: e.LifesizeValue}
	}

	return ans
}
