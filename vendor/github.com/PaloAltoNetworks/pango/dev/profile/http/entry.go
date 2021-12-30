package http

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an http profile.
//
// PAN-OS 7.1+.
type Entry struct {
	Name            string
	TagRegistration bool
	Servers         []Server
	Config          *PayloadFormat
	System          *PayloadFormat
	Threat          *PayloadFormat
	Traffic         *PayloadFormat
	HipMatch        *PayloadFormat
	Url             *PayloadFormat
	Data            *PayloadFormat
	Wildfire        *PayloadFormat
	Tunnel          *PayloadFormat
	UserId          *PayloadFormat
	Gtp             *PayloadFormat
	Auth            *PayloadFormat
	Sctp            *PayloadFormat // 8.1+
	Iptag           *PayloadFormat // 9.0+
}

// Server is an HTTP server spec.
type Server struct {
	Name               string
	Address            string
	Protocol           string
	Port               int
	HttpMethod         string
	Username           string
	Password           string // encrypted
	TlsVersion         string // 9.0+
	CertificateProfile string // 9.0+
}

// PayloadFormat is payload config for a given log type.
type PayloadFormat struct {
	Name       string
	UriFormat  string
	Payload    string
	Headers    []Header
	Parameters []Parameter
}

// Header is an HTTP header.
type Header struct {
	Name  string
	Value string
}

// Parameter is an HTTP parameter.
type Parameter struct {
	Name  string
	Value string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.TagRegistration = s.TagRegistration
	if s.Servers == nil {
		o.Servers = nil
	} else {
		o.Servers = make([]Server, 0, len(s.Servers))
		for _, x := range s.Servers {
			o.Servers = append(o.Servers, Server{
				Name:               x.Name,
				Address:            x.Address,
				Protocol:           x.Protocol,
				Port:               x.Port,
				HttpMethod:         x.HttpMethod,
				Username:           x.Username,
				Password:           x.Password,
				TlsVersion:         x.TlsVersion,
				CertificateProfile: x.CertificateProfile,
			})
		}
	}
	o.Config = copyPayloadFormat(s.Config)
	o.System = copyPayloadFormat(s.System)
	o.Threat = copyPayloadFormat(s.Threat)
	o.Traffic = copyPayloadFormat(s.Traffic)
	o.HipMatch = copyPayloadFormat(s.HipMatch)
	o.Url = copyPayloadFormat(s.Url)
	o.Data = copyPayloadFormat(s.Data)
	o.Wildfire = copyPayloadFormat(s.Wildfire)
	o.Tunnel = copyPayloadFormat(s.Tunnel)
	o.UserId = copyPayloadFormat(s.UserId)
	o.Gtp = copyPayloadFormat(s.Gtp)
	o.Auth = copyPayloadFormat(s.Auth)
	o.Sctp = copyPayloadFormat(s.Sctp)
	o.Iptag = copyPayloadFormat(s.Iptag)
}

func copyPayloadFormat(s *PayloadFormat) *PayloadFormat {
	if s == nil {
		return nil
	}

	ans := PayloadFormat{
		Name:      s.Name,
		UriFormat: s.UriFormat,
		Payload:   s.Payload,
	}

	if len(s.Headers) > 0 {
		ans.Headers = make([]Header, 0, len(s.Headers))
		for _, x := range s.Headers {
			ans.Headers = append(ans.Headers, Header{
				Name:  x.Name,
				Value: x.Value,
			})
		}
	}

	if len(s.Parameters) > 0 {
		ans.Parameters = make([]Parameter, 0, len(s.Parameters))
		for _, x := range s.Parameters {
			ans.Parameters = append(ans.Parameters, Parameter{
				Name:  x.Name,
				Value: x.Value,
			})
		}
	}

	return &ans
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

func (o container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:            o.Name,
		TagRegistration: util.AsBool(o.TagRegistration),
	}

	if o.Server != nil {
		list := make([]Server, 0, len(o.Server.Entries))
		for _, x := range o.Server.Entries {
			list = append(list, Server{
				Name:       x.Name,
				Address:    x.Address,
				Protocol:   x.Protocol,
				Port:       x.Port,
				HttpMethod: x.HttpMethod,
				Username:   x.Username,
				Password:   x.Password,
			})
		}
		ans.Servers = list
	}

	if o.Format != nil {
		ans.Config = normalizePayloadFormat(o.Format.Config)
		ans.System = normalizePayloadFormat(o.Format.System)
		ans.Threat = normalizePayloadFormat(o.Format.Threat)
		ans.Traffic = normalizePayloadFormat(o.Format.Traffic)
		ans.HipMatch = normalizePayloadFormat(o.Format.HipMatch)
		ans.Url = normalizePayloadFormat(o.Format.Url)
		ans.Data = normalizePayloadFormat(o.Format.Data)
		ans.Wildfire = normalizePayloadFormat(o.Format.Wildfire)
		ans.Tunnel = normalizePayloadFormat(o.Format.Tunnel)
		ans.UserId = normalizePayloadFormat(o.Format.UserId)
		ans.Gtp = normalizePayloadFormat(o.Format.Gtp)
		ans.Auth = normalizePayloadFormat(o.Format.Auth)
	}

	return ans
}

func normalizePayloadFormat(val *formatSpec) *PayloadFormat {
	if val == nil {
		return nil
	}

	ans := PayloadFormat{
		Name:      val.Name,
		UriFormat: val.UriFormat,
		Payload:   val.Payload,
	}

	if val.Header != nil {
		list := make([]Header, 0, len(val.Header.Entries))
		for _, x := range val.Header.Entries {
			list = append(list, Header{
				Name:  x.Name,
				Value: x.Value,
			})
		}
		ans.Headers = list
	}

	if val.Parameter != nil {
		list := make([]Parameter, 0, len(val.Parameter.Entries))
		for _, x := range val.Parameter.Entries {
			list = append(list, Parameter{
				Name:  x.Name,
				Value: x.Value,
			})
		}
		ans.Parameters = list
	}

	return &ans
}

type entry_v1 struct {
	XMLName         xml.Name    `xml:"entry"`
	Name            string      `xml:"name,attr"`
	TagRegistration string      `xml:"tag-registration"`
	Server          *servers_v1 `xml:"server"`
	Format          *format_v1  `xml:"format"`
}

type servers_v1 struct {
	Entries []serverEntry_v1 `xml:"entry"`
}

type serverEntry_v1 struct {
	Name       string `xml:"name,attr"`
	Address    string `xml:"address"`
	Protocol   string `xml:"protocol,omitempty"`
	Port       int    `xml:"port,omitempty"`
	HttpMethod string `xml:"http-method"`
	Username   string `xml:"username,omitempty"`
	Password   string `xml:"password,omitempty"`
}

type format_v1 struct {
	Config   *formatSpec `xml:"config"`
	System   *formatSpec `xml:"system"`
	Threat   *formatSpec `xml:"threat"`
	Traffic  *formatSpec `xml:"traffic"`
	HipMatch *formatSpec `xml:"hip-match"`
	Url      *formatSpec `xml:"url"`
	Data     *formatSpec `xml:"data"`
	Wildfire *formatSpec `xml:"wildfire"`
	Tunnel   *formatSpec `xml:"tunnel"`
	UserId   *formatSpec `xml:"userid"`
	Gtp      *formatSpec `xml:"gtp"`
	Auth     *formatSpec `xml:"auth"`
}

type formatSpec struct {
	Name      string   `xml:"name,omitempty"`
	UriFormat string   `xml:"url-format,omitempty"`
	Payload   string   `xml:"payload,omitempty"`
	Header    *headers `xml:"headers"`
	Parameter *params  `xml:"params"`
}

type headers struct {
	Entries []headerEntry `xml:"entry"`
}

type headerEntry struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

type params struct {
	Entries []paramEntry `xml:"entry"`
}

type paramEntry struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:            e.Name,
		TagRegistration: util.YesNo(e.TagRegistration),
	}

	if len(e.Servers) > 0 {
		list := make([]serverEntry_v1, 0, len(e.Servers))
		for _, x := range e.Servers {
			list = append(list, serverEntry_v1{
				Name:       x.Name,
				Address:    x.Address,
				Protocol:   x.Protocol,
				Port:       x.Port,
				HttpMethod: x.HttpMethod,
				Username:   x.Username,
				Password:   x.Password,
			})
		}
		ans.Server = &servers_v1{Entries: list}
	}

	if e.Config != nil || e.System != nil || e.Threat != nil || e.Traffic != nil || e.HipMatch != nil || e.Url != nil || e.Data != nil || e.Wildfire != nil || e.Tunnel != nil || e.UserId != nil || e.Gtp != nil || e.Auth != nil {
		ans.Format = &format_v1{
			Config:   specifyPayloadFormat(e.Config),
			System:   specifyPayloadFormat(e.System),
			Threat:   specifyPayloadFormat(e.Threat),
			Traffic:  specifyPayloadFormat(e.Traffic),
			HipMatch: specifyPayloadFormat(e.HipMatch),
			Url:      specifyPayloadFormat(e.Url),
			Data:     specifyPayloadFormat(e.Data),
			Wildfire: specifyPayloadFormat(e.Wildfire),
			Tunnel:   specifyPayloadFormat(e.Tunnel),
			UserId:   specifyPayloadFormat(e.UserId),
			Gtp:      specifyPayloadFormat(e.Gtp),
			Auth:     specifyPayloadFormat(e.Auth),
		}
	}

	return ans
}

func specifyPayloadFormat(val *PayloadFormat) *formatSpec {
	if val == nil {
		return nil
	}

	ans := formatSpec{
		Name:      val.Name,
		UriFormat: val.UriFormat,
		Payload:   val.Payload,
	}

	if len(val.Headers) > 0 {
		list := make([]headerEntry, 0, len(val.Headers))
		for _, x := range val.Headers {
			list = append(list, headerEntry{
				Name:  x.Name,
				Value: x.Value,
			})
		}
		ans.Header = &headers{Entries: list}
	}

	if len(val.Parameters) > 0 {
		list := make([]paramEntry, 0, len(val.Parameters))
		for _, x := range val.Parameters {
			list = append(list, paramEntry{
				Name:  x.Name,
				Value: x.Value,
			})
		}
		ans.Parameter = &params{Entries: list}
	}

	return &ans
}

// PAN-OS 8.1
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

func (o container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:            o.Name,
		TagRegistration: util.AsBool(o.TagRegistration),
	}

	if o.Server != nil {
		list := make([]Server, 0, len(o.Server.Entries))
		for _, x := range o.Server.Entries {
			list = append(list, Server{
				Name:       x.Name,
				Address:    x.Address,
				Protocol:   x.Protocol,
				Port:       x.Port,
				HttpMethod: x.HttpMethod,
				Username:   x.Username,
				Password:   x.Password,
			})
		}
		ans.Servers = list
	}

	if o.Format != nil {
		ans.Config = normalizePayloadFormat(o.Format.Config)
		ans.System = normalizePayloadFormat(o.Format.System)
		ans.Threat = normalizePayloadFormat(o.Format.Threat)
		ans.Traffic = normalizePayloadFormat(o.Format.Traffic)
		ans.HipMatch = normalizePayloadFormat(o.Format.HipMatch)
		ans.Url = normalizePayloadFormat(o.Format.Url)
		ans.Data = normalizePayloadFormat(o.Format.Data)
		ans.Wildfire = normalizePayloadFormat(o.Format.Wildfire)
		ans.Tunnel = normalizePayloadFormat(o.Format.Tunnel)
		ans.UserId = normalizePayloadFormat(o.Format.UserId)
		ans.Gtp = normalizePayloadFormat(o.Format.Gtp)
		ans.Auth = normalizePayloadFormat(o.Format.Auth)
		ans.Sctp = normalizePayloadFormat(o.Format.Sctp)
	}

	return ans
}

type entry_v2 struct {
	XMLName         xml.Name    `xml:"entry"`
	Name            string      `xml:"name,attr"`
	TagRegistration string      `xml:"tag-registration"`
	Format          *format_v2  `xml:"format"`
	Server          *servers_v1 `xml:"server"`
}

type format_v2 struct {
	Config   *formatSpec `xml:"config"`
	System   *formatSpec `xml:"system"`
	Threat   *formatSpec `xml:"threat"`
	Traffic  *formatSpec `xml:"traffic"`
	HipMatch *formatSpec `xml:"hip-match"`
	Url      *formatSpec `xml:"url"`
	Data     *formatSpec `xml:"data"`
	Wildfire *formatSpec `xml:"wildfire"`
	Tunnel   *formatSpec `xml:"tunnel"`
	UserId   *formatSpec `xml:"userid"`
	Gtp      *formatSpec `xml:"gtp"`
	Auth     *formatSpec `xml:"auth"`
	Sctp     *formatSpec `xml:"sctp"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:            e.Name,
		TagRegistration: util.YesNo(e.TagRegistration),
	}

	if len(e.Servers) > 0 {
		list := make([]serverEntry_v1, 0, len(e.Servers))
		for _, x := range e.Servers {
			list = append(list, serverEntry_v1{
				Name:       x.Name,
				Address:    x.Address,
				Protocol:   x.Protocol,
				Port:       x.Port,
				HttpMethod: x.HttpMethod,
				Username:   x.Username,
				Password:   x.Password,
			})
		}
		ans.Server = &servers_v1{Entries: list}
	}

	if e.Config != nil || e.System != nil || e.Threat != nil || e.Traffic != nil || e.HipMatch != nil || e.Url != nil || e.Data != nil || e.Wildfire != nil || e.Tunnel != nil || e.UserId != nil || e.Gtp != nil || e.Auth != nil || e.Sctp != nil {
		ans.Format = &format_v2{
			Config:   specifyPayloadFormat(e.Config),
			System:   specifyPayloadFormat(e.System),
			Threat:   specifyPayloadFormat(e.Threat),
			Traffic:  specifyPayloadFormat(e.Traffic),
			HipMatch: specifyPayloadFormat(e.HipMatch),
			Url:      specifyPayloadFormat(e.Url),
			Data:     specifyPayloadFormat(e.Data),
			Wildfire: specifyPayloadFormat(e.Wildfire),
			Tunnel:   specifyPayloadFormat(e.Tunnel),
			UserId:   specifyPayloadFormat(e.UserId),
			Gtp:      specifyPayloadFormat(e.Gtp),
			Auth:     specifyPayloadFormat(e.Auth),
			Sctp:     specifyPayloadFormat(e.Sctp),
		}
	}

	return ans
}

// PAN-OS 9.0
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

func (o container_v3) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v3) normalize() Entry {
	ans := Entry{
		Name:            o.Name,
		TagRegistration: util.AsBool(o.TagRegistration),
	}

	if o.Server != nil {
		list := make([]Server, 0, len(o.Server.Entries))
		for _, x := range o.Server.Entries {
			list = append(list, Server{
				Name:               x.Name,
				Address:            x.Address,
				Protocol:           x.Protocol,
				Port:               x.Port,
				HttpMethod:         x.HttpMethod,
				Username:           x.Username,
				Password:           x.Password,
				TlsVersion:         x.TlsVersion,
				CertificateProfile: x.CertificateProfile,
			})
		}
		ans.Servers = list
	}

	if o.Format != nil {
		ans.Config = normalizePayloadFormat(o.Format.Config)
		ans.System = normalizePayloadFormat(o.Format.System)
		ans.Threat = normalizePayloadFormat(o.Format.Threat)
		ans.Traffic = normalizePayloadFormat(o.Format.Traffic)
		ans.HipMatch = normalizePayloadFormat(o.Format.HipMatch)
		ans.Url = normalizePayloadFormat(o.Format.Url)
		ans.Data = normalizePayloadFormat(o.Format.Data)
		ans.Wildfire = normalizePayloadFormat(o.Format.Wildfire)
		ans.Tunnel = normalizePayloadFormat(o.Format.Tunnel)
		ans.UserId = normalizePayloadFormat(o.Format.UserId)
		ans.Gtp = normalizePayloadFormat(o.Format.Gtp)
		ans.Auth = normalizePayloadFormat(o.Format.Auth)
		ans.Sctp = normalizePayloadFormat(o.Format.Sctp)
		ans.Iptag = normalizePayloadFormat(o.Format.Iptag)
	}

	return ans
}

type entry_v3 struct {
	XMLName         xml.Name    `xml:"entry"`
	Name            string      `xml:"name,attr"`
	TagRegistration string      `xml:"tag-registration"`
	Format          *format_v3  `xml:"format"`
	Server          *servers_v2 `xml:"server"`
}

type servers_v2 struct {
	Entries []serverEntry_v2 `xml:"entry"`
}

type serverEntry_v2 struct {
	XMLName            xml.Name `xml:"entry"`
	Name               string   `xml:"name,attr"`
	Address            string   `xml:"address"`
	Protocol           string   `xml:"protocol,omitempty"`
	Port               int      `xml:"port,omitempty"`
	HttpMethod         string   `xml:"http-method"`
	Username           string   `xml:"username,omitempty"`
	Password           string   `xml:"password,omitempty"`
	TlsVersion         string   `xml:"tls-version,omitempty"`
	CertificateProfile string   `xml:"certificate-profile,omitempty"`
}

func (e *serverEntry_v2) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local serverEntry_v2
	ans := local{
		TlsVersion:         "1.2",
		CertificateProfile: "None",
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = serverEntry_v2(ans)
	return nil
}

type format_v3 struct {
	Config   *formatSpec `xml:"config"`
	System   *formatSpec `xml:"system"`
	Threat   *formatSpec `xml:"threat"`
	Traffic  *formatSpec `xml:"traffic"`
	HipMatch *formatSpec `xml:"hip-match"`
	Url      *formatSpec `xml:"url"`
	Data     *formatSpec `xml:"data"`
	Wildfire *formatSpec `xml:"wildfire"`
	Tunnel   *formatSpec `xml:"tunnel"`
	UserId   *formatSpec `xml:"userid"`
	Gtp      *formatSpec `xml:"gtp"`
	Auth     *formatSpec `xml:"auth"`
	Sctp     *formatSpec `xml:"sctp"`
	Iptag    *formatSpec `xml:"iptag"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:            e.Name,
		TagRegistration: util.YesNo(e.TagRegistration),
	}

	if len(e.Servers) > 0 {
		list := make([]serverEntry_v2, 0, len(e.Servers))
		for _, x := range e.Servers {
			list = append(list, serverEntry_v2{
				Name:               x.Name,
				Address:            x.Address,
				Protocol:           x.Protocol,
				Port:               x.Port,
				HttpMethod:         x.HttpMethod,
				Username:           x.Username,
				Password:           x.Password,
				TlsVersion:         x.TlsVersion,
				CertificateProfile: x.CertificateProfile,
			})
		}
		ans.Server = &servers_v2{Entries: list}
	}

	if e.Config != nil || e.System != nil || e.Threat != nil || e.Traffic != nil || e.HipMatch != nil || e.Url != nil || e.Data != nil || e.Wildfire != nil || e.Tunnel != nil || e.UserId != nil || e.Gtp != nil || e.Auth != nil || e.Sctp != nil || e.Iptag != nil {
		ans.Format = &format_v3{
			Config:   specifyPayloadFormat(e.Config),
			System:   specifyPayloadFormat(e.System),
			Threat:   specifyPayloadFormat(e.Threat),
			Traffic:  specifyPayloadFormat(e.Traffic),
			HipMatch: specifyPayloadFormat(e.HipMatch),
			Url:      specifyPayloadFormat(e.Url),
			Data:     specifyPayloadFormat(e.Data),
			Wildfire: specifyPayloadFormat(e.Wildfire),
			Tunnel:   specifyPayloadFormat(e.Tunnel),
			UserId:   specifyPayloadFormat(e.UserId),
			Gtp:      specifyPayloadFormat(e.Gtp),
			Auth:     specifyPayloadFormat(e.Auth),
			Sctp:     specifyPayloadFormat(e.Sctp),
			Iptag:    specifyPayloadFormat(e.Iptag),
		}
	}

	return ans
}
