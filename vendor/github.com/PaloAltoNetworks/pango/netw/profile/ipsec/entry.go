package ipsec

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
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
	if s.Encryption == nil {
		o.Encryption = nil
	} else {
		o.Encryption = make([]string, len(s.Encryption))
		copy(o.Encryption, s.Encryption)
	}
	if s.Authentication == nil {
		o.Authentication = nil
	} else {
		o.Authentication = make([]string, len(s.Authentication))
		copy(o.Authentication, s.Authentication)
	}
	o.DhGroup = s.DhGroup
	o.LifetimeType = s.LifetimeType
	o.LifetimeValue = s.LifetimeValue
	o.LifesizeType = s.LifesizeType
	o.LifesizeValue = s.LifesizeValue
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
			case EncryptionNull:
				nv[i] = "null"
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
			case EncryptionNull:
				nv[i] = "null"
			default:
				nv[i] = vals[i]
			}
		}
	default:
		copy(nv, vals)
	}

	return nv
}

func normalizeEncryption(vals []string) []string {
	if vals == nil {
		return nil
	}

	nv := make([]string, len(vals))
	for i := range vals {
		switch vals[i] {
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
			nv[i] = vals[i]
		}
	}

	return nv
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
		Name:    o.Name,
		DhGroup: o.DhGroup,
	}

	if o.Esp != nil {
		ans.Protocol = ProtocolEsp
		ans.Encryption = normalizeEncryption(util.MemToStr(o.Esp.Encryption))
		ans.Authentication = util.MemToStr(o.Esp.Authentication)
	} else if o.Ah != nil {
		ans.Protocol = ProtocolAh
		ans.Authentication = util.MemToStr(o.Ah.Authentication)
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

	if o.Lifesize != nil {
		if o.Lifesize.Kb != 0 {
			ans.LifesizeType = SizeKb
			ans.LifesizeValue = o.Lifesize.Kb
		} else if o.Lifesize.Mb != 0 {
			ans.LifesizeType = SizeMb
			ans.LifesizeValue = o.Lifesize.Mb
		} else if o.Lifesize.Gb != 0 {
			ans.LifesizeType = SizeGb
			ans.LifesizeValue = o.Lifesize.Gb
		} else if o.Lifesize.Tb != 0 {
			ans.LifesizeType = SizeTb
			ans.LifesizeValue = o.Lifesize.Tb
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
			Encryption:     util.StrToMem(specifyEncryption(e.Encryption, 1)),
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
		DhGroup: o.DhGroup,
	}

	if o.Esp != nil {
		ans.Protocol = ProtocolEsp
		ans.Encryption = normalizeEncryption(util.MemToStr(o.Esp.Encryption))
		ans.Authentication = util.MemToStr(o.Esp.Authentication)
	} else if o.Ah != nil {
		ans.Protocol = ProtocolAh
		ans.Authentication = util.MemToStr(o.Ah.Authentication)
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

	if o.Lifesize != nil {
		if o.Lifesize.Kb != 0 {
			ans.LifesizeType = SizeKb
			ans.LifesizeValue = o.Lifesize.Kb
		} else if o.Lifesize.Mb != 0 {
			ans.LifesizeType = SizeMb
			ans.LifesizeValue = o.Lifesize.Mb
		} else if o.Lifesize.Gb != 0 {
			ans.LifesizeType = SizeGb
			ans.LifesizeValue = o.Lifesize.Gb
		} else if o.Lifesize.Tb != 0 {
			ans.LifesizeType = SizeTb
			ans.LifesizeValue = o.Lifesize.Tb
		}
	}

	return ans
}

type entry_v2 struct {
	XMLName  xml.Name      `xml:"entry"`
	Name     string        `xml:"name,attr"`
	Esp      *espType      `xml:"esp"`
	Ah       *ahType       `xml:"ah"`
	DhGroup  string        `xml:"dh-group,omitempty"`
	Lifetime *lifetimeType `xml:"lifetime"`
	Lifesize *lifesizeType `xml:"lifesize"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:    e.Name,
		DhGroup: e.DhGroup,
	}

	switch e.Protocol {
	case ProtocolEsp:
		ans.Esp = &espType{
			Encryption:     util.StrToMem(specifyEncryption(e.Encryption, 2)),
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
