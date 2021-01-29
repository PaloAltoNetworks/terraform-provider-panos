package email

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of an email profile.
//
// PAN-OS 7.1+.
type Entry struct {
	Name              string
	Config            string
	System            string
	Threat            string
	Traffic           string
	HipMatch          string
	Url               string // 8.0+
	Data              string // 8.0+
	Wildfire          string // 8.0+
	Tunnel            string // 8.0+
	UserId            string // 8.0+
	Gtp               string // 8.0+
	Auth              string // 8.0+
	Sctp              string // 8.1+
	Iptag             string // 9.0+
	EscapedCharacters string
	EscapeCharacter   string

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Config = s.Config
	o.System = s.System
	o.Threat = s.Threat
	o.Traffic = s.Traffic
	o.HipMatch = s.HipMatch
	o.Url = s.Url
	o.Data = s.Data
	o.Wildfire = s.Wildfire
	o.Tunnel = s.Tunnel
	o.UserId = s.UserId
	o.Gtp = s.Gtp
	o.Auth = s.Auth
	o.Sctp = s.Sctp
	o.Iptag = s.Iptag
	o.EscapedCharacters = s.EscapedCharacters
	o.EscapeCharacter = s.EscapeCharacter
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
	}

	if o.Answer.Server != nil {
		ans.raw = map[string]string{
			"srv": util.CleanRawXml(o.Answer.Server.Text),
		}
	}

	if o.Answer.Format != nil {
		ans.Config = o.Answer.Format.Config
		ans.System = o.Answer.Format.System
		ans.Threat = o.Answer.Format.Threat
		ans.Traffic = o.Answer.Format.Traffic
		ans.HipMatch = o.Answer.Format.HipMatch

		if o.Answer.Format.Esc != nil {
			ans.EscapedCharacters = o.Answer.Format.Esc.EscapedCharacters
			ans.EscapeCharacter = o.Answer.Format.Esc.EscapeCharacter
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
	}

	if o.Answer.Server != nil {
		ans.raw = map[string]string{
			"srv": util.CleanRawXml(o.Answer.Server.Text),
		}
	}

	if o.Answer.Format != nil {
		ans.Config = o.Answer.Format.Config
		ans.System = o.Answer.Format.System
		ans.Threat = o.Answer.Format.Threat
		ans.Traffic = o.Answer.Format.Traffic
		ans.HipMatch = o.Answer.Format.HipMatch
		ans.Url = o.Answer.Format.Url
		ans.Data = o.Answer.Format.Data
		ans.Wildfire = o.Answer.Format.Wildfire
		ans.Tunnel = o.Answer.Format.Tunnel
		ans.UserId = o.Answer.Format.UserId
		ans.Gtp = o.Answer.Format.Gtp
		ans.Auth = o.Answer.Format.Auth

		if o.Answer.Format.Esc != nil {
			ans.EscapedCharacters = o.Answer.Format.Esc.EscapedCharacters
			ans.EscapeCharacter = o.Answer.Format.Esc.EscapeCharacter
		}
	}

	return ans
}

type container_v3 struct {
	Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
	ans := Entry{
		Name: o.Answer.Name,
	}

	if o.Answer.Server != nil {
		ans.raw = map[string]string{
			"srv": util.CleanRawXml(o.Answer.Server.Text),
		}
	}

	if o.Answer.Format != nil {
		ans.Config = o.Answer.Format.Config
		ans.System = o.Answer.Format.System
		ans.Threat = o.Answer.Format.Threat
		ans.Traffic = o.Answer.Format.Traffic
		ans.HipMatch = o.Answer.Format.HipMatch
		ans.Url = o.Answer.Format.Url
		ans.Data = o.Answer.Format.Data
		ans.Wildfire = o.Answer.Format.Wildfire
		ans.Tunnel = o.Answer.Format.Tunnel
		ans.UserId = o.Answer.Format.UserId
		ans.Gtp = o.Answer.Format.Gtp
		ans.Auth = o.Answer.Format.Auth
		ans.Sctp = o.Answer.Format.Sctp

		if o.Answer.Format.Esc != nil {
			ans.EscapedCharacters = o.Answer.Format.Esc.EscapedCharacters
			ans.EscapeCharacter = o.Answer.Format.Esc.EscapeCharacter
		}
	}

	return ans
}

type container_v4 struct {
	Answer entry_v4 `xml:"result>entry"`
}

func (o *container_v4) Normalize() Entry {
	ans := Entry{
		Name: o.Answer.Name,
	}

	if o.Answer.Server != nil {
		ans.raw = map[string]string{
			"srv": util.CleanRawXml(o.Answer.Server.Text),
		}
	}

	if o.Answer.Format != nil {
		ans.Config = o.Answer.Format.Config
		ans.System = o.Answer.Format.System
		ans.Threat = o.Answer.Format.Threat
		ans.Traffic = o.Answer.Format.Traffic
		ans.HipMatch = o.Answer.Format.HipMatch
		ans.Url = o.Answer.Format.Url
		ans.Data = o.Answer.Format.Data
		ans.Wildfire = o.Answer.Format.Wildfire
		ans.Tunnel = o.Answer.Format.Tunnel
		ans.UserId = o.Answer.Format.UserId
		ans.Gtp = o.Answer.Format.Gtp
		ans.Auth = o.Answer.Format.Auth
		ans.Sctp = o.Answer.Format.Sctp
		ans.Iptag = o.Answer.Format.Iptag

		if o.Answer.Format.Esc != nil {
			ans.EscapedCharacters = o.Answer.Format.Esc.EscapedCharacters
			ans.EscapeCharacter = o.Answer.Format.Esc.EscapeCharacter
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name     `xml:"entry"`
	Name    string       `xml:"name,attr"`
	Server  *util.RawXml `xml:"server"`
	Format  *format_v1   `xml:"format"`
}

type format_v1 struct {
	Config   string `xml:"config,omitempty"`
	System   string `xml:"system,omitempty"`
	Threat   string `xml:"thread,omitempty"`
	Traffic  string `xml:"traffic,omitempty"`
	HipMatch string `xml:"hip-match,omitempty"`
	Esc      *esc   `xml:"escaping"`
}

type esc struct {
	EscapedCharacters string `xml:"escaped-characters"`
	EscapeCharacter   string `xml:"escape-character"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	hasEsc := e.EscapedCharacters != "" || e.EscapeCharacter != ""
	if e.Config != "" || e.System != "" || e.Threat != "" || e.Traffic != "" || e.HipMatch != "" || hasEsc {
		ans.Format = &format_v1{
			Config:   e.Config,
			System:   e.System,
			Threat:   e.Threat,
			Traffic:  e.Traffic,
			HipMatch: e.HipMatch,
		}

		if hasEsc {
			ans.Format.Esc = &esc{
				EscapedCharacters: e.EscapedCharacters,
				EscapeCharacter:   e.EscapeCharacter,
			}
		}
	}

	return ans
}

type entry_v2 struct {
	XMLName xml.Name     `xml:"entry"`
	Name    string       `xml:"name,attr"`
	Server  *util.RawXml `xml:"server"`
	Format  *format_v2   `xml:"format"`
}

type format_v2 struct {
	Config   string `xml:"config,omitempty"`
	System   string `xml:"system,omitempty"`
	Threat   string `xml:"thread,omitempty"`
	Traffic  string `xml:"traffic,omitempty"`
	HipMatch string `xml:"hip-match,omitempty"`
	Url      string `xml:"url,omitempty"`
	Data     string `xml:"data,omitempty"`
	Wildfire string `xml:"wildfire,omitempty"`
	Tunnel   string `xml:"tunnel,omitempty"`
	UserId   string `xml:"userid,omitempty"`
	Gtp      string `xml:"gtp,omitempty"`
	Auth     string `xml:"auth,omitempty"`
	Esc      *esc   `xml:"escaping"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name: e.Name,
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	hasEsc := e.EscapedCharacters != "" || e.EscapeCharacter != ""
	if e.Config != "" || e.System != "" || e.Threat != "" || e.Traffic != "" || e.HipMatch != "" || e.Url != "" || e.Data != "" || e.Wildfire != "" || e.Tunnel != "" || e.UserId != "" || e.Gtp != "" || e.Auth != "" || hasEsc {
		ans.Format = &format_v2{
			Config:   e.Config,
			System:   e.System,
			Threat:   e.Threat,
			Traffic:  e.Traffic,
			HipMatch: e.HipMatch,
			Url:      e.Url,
			Data:     e.Data,
			Wildfire: e.Wildfire,
			Tunnel:   e.Tunnel,
			UserId:   e.UserId,
			Gtp:      e.Gtp,
			Auth:     e.Auth,
		}

		if hasEsc {
			ans.Format.Esc = &esc{
				EscapedCharacters: e.EscapedCharacters,
				EscapeCharacter:   e.EscapeCharacter,
			}
		}
	}

	return ans
}

type entry_v3 struct {
	XMLName xml.Name     `xml:"entry"`
	Name    string       `xml:"name,attr"`
	Server  *util.RawXml `xml:"server"`
	Format  *format_v3   `xml:"format"`
}

type format_v3 struct {
	Config   string `xml:"config,omitempty"`
	System   string `xml:"system,omitempty"`
	Threat   string `xml:"thread,omitempty"`
	Traffic  string `xml:"traffic,omitempty"`
	HipMatch string `xml:"hip-match,omitempty"`
	Url      string `xml:"url,omitempty"`
	Data     string `xml:"data,omitempty"`
	Wildfire string `xml:"wildfire,omitempty"`
	Tunnel   string `xml:"tunnel,omitempty"`
	UserId   string `xml:"userid,omitempty"`
	Gtp      string `xml:"gtp,omitempty"`
	Auth     string `xml:"auth,omitempty"`
	Sctp     string `xml:"sctp,omitempty"`
	Esc      *esc   `xml:"escaping"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name: e.Name,
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	hasEsc := e.EscapedCharacters != "" || e.EscapeCharacter != ""
	if e.Config != "" || e.System != "" || e.Threat != "" || e.Traffic != "" || e.HipMatch != "" || e.Url != "" || e.Data != "" || e.Wildfire != "" || e.Tunnel != "" || e.UserId != "" || e.Gtp != "" || e.Auth != "" || e.Sctp != "" || hasEsc {
		ans.Format = &format_v3{
			Config:   e.Config,
			System:   e.System,
			Threat:   e.Threat,
			Traffic:  e.Traffic,
			HipMatch: e.HipMatch,
			Url:      e.Url,
			Data:     e.Data,
			Wildfire: e.Wildfire,
			Tunnel:   e.Tunnel,
			UserId:   e.UserId,
			Gtp:      e.Gtp,
			Auth:     e.Auth,
			Sctp:     e.Sctp,
		}

		if hasEsc {
			ans.Format.Esc = &esc{
				EscapedCharacters: e.EscapedCharacters,
				EscapeCharacter:   e.EscapeCharacter,
			}
		}
	}

	return ans
}

type entry_v4 struct {
	XMLName xml.Name     `xml:"entry"`
	Name    string       `xml:"name,attr"`
	Server  *util.RawXml `xml:"server"`
	Format  *format_v4   `xml:"format"`
}

type format_v4 struct {
	Config   string `xml:"config,omitempty"`
	System   string `xml:"system,omitempty"`
	Threat   string `xml:"thread,omitempty"`
	Traffic  string `xml:"traffic,omitempty"`
	HipMatch string `xml:"hip-match,omitempty"`
	Url      string `xml:"url,omitempty"`
	Data     string `xml:"data,omitempty"`
	Wildfire string `xml:"wildfire,omitempty"`
	Tunnel   string `xml:"tunnel,omitempty"`
	UserId   string `xml:"userid,omitempty"`
	Gtp      string `xml:"gtp,omitempty"`
	Auth     string `xml:"auth,omitempty"`
	Sctp     string `xml:"sctp,omitempty"`
	Iptag    string `xml:"iptag,omitempty"`
	Esc      *esc   `xml:"escaping"`
}

func specify_v4(e Entry) interface{} {
	ans := entry_v4{
		Name: e.Name,
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	hasEsc := e.EscapedCharacters != "" || e.EscapeCharacter != ""
	if e.Config != "" || e.System != "" || e.Threat != "" || e.Traffic != "" || e.HipMatch != "" || e.Url != "" || e.Data != "" || e.Wildfire != "" || e.Tunnel != "" || e.UserId != "" || e.Gtp != "" || e.Auth != "" || e.Sctp != "" || e.Iptag != "" || hasEsc {
		ans.Format = &format_v4{
			Config:   e.Config,
			System:   e.System,
			Threat:   e.Threat,
			Traffic:  e.Traffic,
			HipMatch: e.HipMatch,
			Url:      e.Url,
			Data:     e.Data,
			Wildfire: e.Wildfire,
			Tunnel:   e.Tunnel,
			UserId:   e.UserId,
			Gtp:      e.Gtp,
			Auth:     e.Auth,
			Sctp:     e.Sctp,
			Iptag:    e.Iptag,
		}

		if hasEsc {
			ans.Format.Esc = &esc{
				EscapedCharacters: e.EscapedCharacters,
				EscapeCharacter:   e.EscapeCharacter,
			}
		}
	}

	return ans
}
