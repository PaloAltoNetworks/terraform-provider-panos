package edl

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an
// external dynamic list.
type Entry struct {
	Name               string
	Type               string
	Description        string
	Source             string
	CertificateProfile string // PAN-OS 8.0+
	Username           string // PAN-OS 8.0+
	Password           string // PAN-OS 8.0+
	ExpandDomain       bool   // PAN-OS 9.0+
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
	o.ExpandDomain = s.ExpandDomain
	o.Repeat = s.Repeat
	o.RepeatAt = s.RepeatAt
	o.RepeatDayOfWeek = s.RepeatDayOfWeek
	o.RepeatDayOfMonth = s.RepeatDayOfMonth
	if s.Exceptions == nil {
		o.Exceptions = nil
	} else {
		o.Exceptions = make([]string, len(s.Exceptions))
		copy(o.Exceptions, s.Exceptions)
	}
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
		Name:        o.Name,
		Type:        o.Type,
		Description: o.Description,
		Source:      o.Source,
	}

	if o.Repeat.FiveMinute != nil {
		ans.Repeat = RepeatEveryFiveMinutes
	} else if o.Repeat.Hourly != nil {
		ans.Repeat = RepeatHourly
	} else if o.Repeat.Daily != nil {
		ans.Repeat = RepeatDaily
		ans.RepeatAt = o.Repeat.Daily.At
	} else if o.Repeat.Weekly != nil {
		ans.Repeat = RepeatWeekly
		ans.RepeatAt = o.Repeat.Weekly.At
		ans.RepeatDayOfWeek = o.Repeat.Weekly.DayOfWeek
	} else if o.Repeat.Monthly != nil {
		ans.Repeat = RepeatMonthly
		ans.RepeatAt = o.Repeat.Monthly.At
		ans.RepeatDayOfMonth = o.Repeat.Monthly.DayOfMonth
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name `xml:"entry"`
	Name        string   `xml:"name,attr"`
	Type        string   `xml:"type"`
	Description string   `xml:"description,omitempty"`
	Source      string   `xml:"url"`
	Repeat      repeat   `xml:"recurring"`
}

type repeat struct {
	FiveMinute *string    `xml:"five-minute"`
	Hourly     *string    `xml:"hourly"`
	Daily      *timeDay   `xml:"daily"`
	Weekly     *timeWeek  `xml:"weekly"`
	Monthly    *timeMonth `xml:"monthly"`
}

type timeDay struct {
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

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Type:        e.Type,
		Description: e.Description,
		Source:      e.Source,
	}

	switch e.Repeat {
	case RepeatEveryFiveMinutes:
		s := ""
		ans.Repeat.FiveMinute = &s
	case RepeatHourly:
		s := ""
		ans.Repeat.Hourly = &s
	case RepeatDaily:
		ans.Repeat.Daily = &timeDay{e.RepeatAt}
	case RepeatWeekly:
		ans.Repeat.Weekly = &timeWeek{e.RepeatAt, e.RepeatDayOfWeek}
	case RepeatMonthly:
		ans.Repeat.Monthly = &timeMonth{e.RepeatAt, e.RepeatDayOfMonth}
	}

	return ans
}

// PAN-OS 8.0.
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
		Name: o.Name,
	}

	var sp *typeSpec

	if o.PredefinedIp != nil {
		ans.Type = TypePredefinedIp
		ans.Description = o.PredefinedIp.Description
		ans.Source = o.PredefinedIp.Source
		ans.Exceptions = util.MemToStr(o.PredefinedIp.Exceptions)
	} else if o.Ip != nil {
		ans.Type = TypeIp
		sp = o.Ip
	} else if o.Domain != nil {
		ans.Type = TypeDomain
		sp = o.Domain
	} else if o.Url != nil {
		ans.Type = TypeUrl
		sp = o.Url
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
	Repeat             repeat           `xml:"recurring"`
	Exceptions         *util.MemberType `xml:"exception-list"`
}

type authType struct {
	Username string `xml:"username"`
	Password string `xml:"password"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name: e.Name,
	}

	switch e.Type {
	case TypePredefinedIp:
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
			spec.Repeat.Daily = &timeDay{e.RepeatAt}
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

// PAN-OS 9.0.
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
		Name: o.Name,
	}

	var sp *typeSpec

	if o.PredefinedIp != nil {
		ans.Type = TypePredefinedIp
		ans.Description = o.PredefinedIp.Description
		ans.Source = o.PredefinedIp.Source
		ans.Exceptions = util.MemToStr(o.PredefinedIp.Exceptions)
	} else if o.Ip != nil {
		ans.Type = TypeIp
		sp = o.Ip
	} else if o.Domain != nil {
		ans.Type = TypeDomain
		ans.ExpandDomain = util.AsBool(o.Domain.ExpandDomain)
		sp = &typeSpec{
			Description:        o.Domain.Description,
			Source:             o.Domain.Source,
			CertificateProfile: o.Domain.CertificateProfile,
			Auth:               o.Domain.Auth,
			Repeat:             o.Domain.Repeat,
			Exceptions:         o.Domain.Exceptions,
		}
	} else if o.Url != nil {
		ans.Type = TypeUrl
		sp = o.Url
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

type entry_v3 struct {
	XMLName      xml.Name        `xml:"entry"`
	Name         string          `xml:"name,attr"`
	PredefinedIp *typePredefined `xml:"type>predefined-ip"`
	Ip           *typeSpec       `xml:"type>ip"`
	Domain       *domainSpec     `xml:"type>domain"`
	Url          *typeSpec       `xml:"type>url"`
}

type domainSpec struct {
	Description        string           `xml:"description,omitempty"`
	Source             string           `xml:"url"`
	CertificateProfile string           `xml:"certificate-profile,omitempty"`
	Auth               *authType        `xml:"auth"`
	Repeat             repeat           `xml:"recurring"`
	Exceptions         *util.MemberType `xml:"exception-list"`
	ExpandDomain       string           `xml:"expand-domain"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name: e.Name,
	}

	switch e.Type {
	case TypePredefinedIp:
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
			spec.Repeat.Daily = &timeDay{e.RepeatAt}
		case RepeatWeekly:
			spec.Repeat.Weekly = &timeWeek{e.RepeatAt, e.RepeatDayOfWeek}
		case RepeatMonthly:
			spec.Repeat.Monthly = &timeMonth{e.RepeatAt, e.RepeatDayOfMonth}
		}

		switch e.Type {
		case TypeIp:
			ans.Ip = spec
		case TypeDomain:
			ans.Domain = &domainSpec{
				Description:        spec.Description,
				Source:             spec.Source,
				CertificateProfile: spec.CertificateProfile,
				Auth:               spec.Auth,
				Repeat:             spec.Repeat,
				Exceptions:         spec.Exceptions,
				ExpandDomain:       util.YesNo(e.ExpandDomain),
			}
		case TypeUrl:
			ans.Url = spec
		}
	}

	return ans
}

// PAN-OS 10.0.
type container_v4 struct {
	Answer []entry_v4 `xml:"entry"`
}

func (o *container_v4) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v4) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v4) normalize() Entry {
	ans := Entry{
		Name: o.Name,
	}

	var sp *typeSpec

	if o.PredefinedIp != nil {
		ans.Type = TypePredefinedIp
		ans.Description = o.PredefinedIp.Description
		ans.Source = o.PredefinedIp.Source
		ans.Exceptions = util.MemToStr(o.PredefinedIp.Exceptions)
	} else if o.PredefinedUrl != nil {
		ans.Type = TypePredefinedUrl
		ans.Description = o.PredefinedUrl.Description
		ans.Source = o.PredefinedUrl.Source
		ans.Exceptions = util.MemToStr(o.PredefinedUrl.Exceptions)
	} else if o.Ip != nil {
		ans.Type = TypeIp
		sp = o.Ip
	} else if o.Domain != nil {
		ans.Type = TypeDomain
		ans.ExpandDomain = util.AsBool(o.Domain.ExpandDomain)
		sp = &typeSpec{
			Description:        o.Domain.Description,
			Source:             o.Domain.Source,
			CertificateProfile: o.Domain.CertificateProfile,
			Auth:               o.Domain.Auth,
			Repeat:             o.Domain.Repeat,
			Exceptions:         o.Domain.Exceptions,
		}
	} else if o.Url != nil {
		ans.Type = TypeUrl
		sp = o.Url
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

type entry_v4 struct {
	XMLName       xml.Name        `xml:"entry"`
	Name          string          `xml:"name,attr"`
	PredefinedIp  *typePredefined `xml:"type>predefined-ip"`
	PredefinedUrl *typePredefined `xml:"type>predefined-url"`
	Ip            *typeSpec       `xml:"type>ip"`
	Domain        *domainSpec     `xml:"type>domain"`
	Url           *typeSpec       `xml:"type>url"`
}

func specify_v4(e Entry) interface{} {
	ans := entry_v4{
		Name: e.Name,
	}

	switch e.Type {
	case TypePredefinedIp:
		ans.PredefinedIp = &typePredefined{
			Description: e.Description,
			Source:      e.Source,
			Exceptions:  util.StrToMem(e.Exceptions),
		}
	case TypePredefinedUrl:
		ans.PredefinedUrl = &typePredefined{
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
			spec.Repeat.Daily = &timeDay{e.RepeatAt}
		case RepeatWeekly:
			spec.Repeat.Weekly = &timeWeek{e.RepeatAt, e.RepeatDayOfWeek}
		case RepeatMonthly:
			spec.Repeat.Monthly = &timeMonth{e.RepeatAt, e.RepeatDayOfMonth}
		}

		switch e.Type {
		case TypeIp:
			ans.Ip = spec
		case TypeDomain:
			ans.Domain = &domainSpec{
				Description:        spec.Description,
				Source:             spec.Source,
				CertificateProfile: spec.CertificateProfile,
				Auth:               spec.Auth,
				Repeat:             spec.Repeat,
				Exceptions:         spec.Exceptions,
				ExpandDomain:       util.YesNo(e.ExpandDomain),
			}
		case TypeUrl:
			ans.Url = spec
		}
	}

	return ans
}
