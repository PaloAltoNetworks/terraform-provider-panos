package edl

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Constants for Entry.Type field.  Only TypeIp is valid for PAN-OS 7.0 and
// earlier.  TypePredefined is valid for PAN-OS 8.0 and later.
const (
	TypeIp         string = "ip"
	TypeDomain     string = "domain"
	TypeUrl        string = "url"
	TypePredefined string = "predefined"
)

// Constants for the Repeat field.  Option "RepeatEveryFiveMinutes" is valid
// for PAN-OS 8.0 and higher.
const (
	RepeatEveryFiveMinutes = "every five minutes"
	RepeatHourly           = "hourly"
	RepeatDaily            = "daily"
	RepeatWeekly           = "weekly"
	RepeatMonthly          = "monthly"
)

// Entry is a normalized, version independent representation of an
// external dynamic list.
type Entry struct {
	Name               string
	Type               string
	Description        string
	Source             string
	CertificateProfile string
	Username           string
	Password           string
	Repeat             string
	RepeatAt           string
	RepeatDayOfWeek    string
	RepeatDayOfMonth   int
	Exceptions         []string // ordered
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Type = s.Type
	o.Description = s.Description
	o.Source = s.Source
	o.CertificateProfile = s.CertificateProfile
	o.Username = s.Username
	o.Password = s.Password
	o.Repeat = s.Repeat
	o.RepeatAt = s.RepeatAt
	o.RepeatDayOfWeek = s.RepeatDayOfWeek
	o.RepeatDayOfMonth = s.RepeatDayOfMonth
	o.Exceptions = s.Exceptions
}

/** Structs / functions for normalization. **/

type normalizer interface {
	Normalize() Entry
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
	ans := Entry{
		Name:        o.Answer.Name,
		Type:        o.Answer.Type,
		Description: o.Answer.Description,
		Source:      o.Answer.Source,
	}

	if o.Answer.Repeat.FiveMinute != nil {
		ans.Repeat = RepeatEveryFiveMinutes
	} else if o.Answer.Repeat.Hourly != nil {
		ans.Repeat = RepeatHourly
		ans.RepeatAt = o.Answer.Repeat.Hourly.At
	} else if o.Answer.Repeat.Daily != nil {
		ans.Repeat = RepeatDaily
		ans.RepeatAt = o.Answer.Repeat.Daily.At
	} else if o.Answer.Repeat.Weekly != nil {
		ans.Repeat = RepeatWeekly
		ans.RepeatAt = o.Answer.Repeat.Weekly.At
		ans.RepeatDayOfWeek = o.Answer.Repeat.Weekly.DayOfWeek
	} else if o.Answer.Repeat.Monthly != nil {
		ans.Repeat = RepeatMonthly
		ans.RepeatAt = o.Answer.Repeat.Monthly.At
		ans.RepeatDayOfMonth = o.Answer.Repeat.Monthly.DayOfMonth
	}

	return ans
}

type container_v2 struct {
	Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
	ans := Entry{
		Name: o.Answer.Name,
	}

	var sp *typeSpec

	if o.Answer.PredefinedIp != nil {
		ans.Type = TypePredefined
		ans.Description = o.Answer.PredefinedIp.Description
		ans.Source = o.Answer.PredefinedIp.Source
		ans.Exceptions = util.MemToStr(o.Answer.PredefinedIp.Exceptions)
	} else if o.Answer.Ip != nil {
		ans.Type = TypeIp
		sp = o.Answer.Ip
	} else if o.Answer.Domain != nil {
		ans.Type = TypeDomain
		sp = o.Answer.Domain
	} else if o.Answer.Url != nil {
		ans.Type = TypeUrl
		sp = o.Answer.Url
	}

	if sp != nil {
		ans.Description = sp.Description
		ans.Source = sp.Source
		ans.CertificateProfile = sp.CertificateProfile
		ans.Exceptions = util.MemToStr(sp.Exceptions)
		if sp.Auth != nil {
			ans.Username = sp.Auth.Username
			ans.Password = sp.Auth.Password
		}
		if sp.Repeat.FiveMinute != nil {
			ans.Repeat = RepeatEveryFiveMinutes
		} else if sp.Repeat.Hourly != nil {
			ans.Repeat = RepeatHourly
		} else if sp.Repeat.Daily != nil {
			ans.Repeat = RepeatDaily
			ans.RepeatAt = sp.Repeat.Daily.At
		} else if sp.Repeat.Weekly != nil {
			ans.Repeat = RepeatWeekly
			ans.RepeatAt = sp.Repeat.Weekly.At
			ans.RepeatDayOfWeek = sp.Repeat.Weekly.DayOfWeek
		} else if sp.Repeat.Monthly != nil {
			ans.Repeat = RepeatMonthly
			ans.RepeatAt = sp.Repeat.Monthly.At
			ans.RepeatDayOfMonth = sp.Repeat.Monthly.DayOfMonth
		}
	}

	return ans
}

// Ideally there would be one struct for PAN-OS 6.1 & 7.0 and another for
// PAN-OS 7.1, but since the difference is minimal, I'm using the same struct.
//
// Probably revisit this at a later time..?
type entry_v1 struct {
	XMLName     xml.Name `xml:"entry"`
	Name        string   `xml:"name,attr"`
	Type        string   `xml:"type"`
	Description string   `xml:"description,omitempty"`
	Source      string   `xml:"url"`
	Repeat      rep_v1   `xml:"recurring"`
}

type rep_v1 struct {
	FiveMinute *string    `xml:"five-minute"`
	Hourly     *timeAt    `xml:"hourly"`
	Daily      *timeAt    `xml:"daily"`
	Weekly     *timeWeek  `xml:"weekly"`
	Monthly    *timeMonth `xml:"monthly"`
}

type timeAt struct {
	At string `xml:"at,omitempty"`
}

type timeWeek struct {
	At        string `xml:"at"`
	DayOfWeek string `xml:"day-of-week"`
}

type timeMonth struct {
	At         string `xml:"at"`
	DayOfMonth int    `xml:"day-of-month"`
}

type entry_v2 struct {
	XMLName      xml.Name        `xml:"entry"`
	Name         string          `xml:"name,attr"`
	PredefinedIp *typePredefined `xml:"type>predefined-ip"`
	Ip           *typeSpec       `xml:"type>ip"`
	Domain       *typeSpec       `xml:"type>domain"`
	Url          *typeSpec       `xml:"type>url"`
}

type typePredefined struct {
	Description string           `xml:"description"`
	Source      string           `xml:"url"`
	Exceptions  *util.MemberType `xml:"exception-list"`
}

type typeSpec struct {
	Description        string           `xml:"description,omitempty"`
	Source             string           `xml:"url"`
	CertificateProfile string           `xml:"certificate-profile,omitempty"`
	Auth               *authType        `xml:"auth"`
	Repeat             rep_v2           `xml:"recurring"`
	Exceptions         *util.MemberType `xml:"exception-list"`
}

type authType struct {
	Username string `xml:"username"`
	Password string `xml:"password"`
}

type rep_v2 struct {
	FiveMinute *string    `xml:"five-minute"`
	Hourly     *string    `xml:"hourly"`
	Daily      *timeAt    `xml:"daily"`
	Weekly     *timeWeek  `xml:"weekly"`
	Monthly    *timeMonth `xml:"monthly"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Type:        e.Type,
		Description: e.Description,
		Source:      e.Source,
	}

	switch e.Repeat {
	case RepeatEveryFiveMinutes:
		sp := ""
		ans.Repeat.FiveMinute = &sp
	case RepeatHourly:
		ans.Repeat.Hourly = &timeAt{e.RepeatAt}
	case RepeatDaily:
		ans.Repeat.Daily = &timeAt{e.RepeatAt}
	case RepeatWeekly:
		ans.Repeat.Weekly = &timeWeek{e.RepeatAt, e.RepeatDayOfWeek}
	case RepeatMonthly:
		ans.Repeat.Monthly = &timeMonth{e.RepeatAt, e.RepeatDayOfMonth}
	}

	return ans
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name: e.Name,
	}

	switch e.Type {
	case TypePredefined:
		ans.PredefinedIp = &typePredefined{
			Description: e.Description,
			Source:      e.Source,
			Exceptions:  util.StrToMem(e.Exceptions),
		}
	default:
		spec := &typeSpec{
			Description:        e.Description,
			Source:             e.Source,
			CertificateProfile: e.CertificateProfile,
			Exceptions:         util.StrToMem(e.Exceptions),
		}

		if e.Username != "" || e.Password != "" {
			spec.Auth = &authType{e.Username, e.Password}
		}

		sp := ""
		switch e.Repeat {
		case RepeatEveryFiveMinutes:
			spec.Repeat.FiveMinute = &sp
		case RepeatHourly:
			spec.Repeat.Hourly = &sp
		case RepeatDaily:
			spec.Repeat.Daily = &timeAt{e.RepeatAt}
		case RepeatWeekly:
			spec.Repeat.Weekly = &timeWeek{e.RepeatAt, e.RepeatDayOfWeek}
		case RepeatMonthly:
			spec.Repeat.Monthly = &timeMonth{e.RepeatAt, e.RepeatDayOfMonth}
		}

		switch e.Type {
		case TypeIp:
			ans.Ip = spec
		case TypeDomain:
			ans.Domain = spec
		case TypeUrl:
			ans.Url = spec
		}
	}

	return ans
}
