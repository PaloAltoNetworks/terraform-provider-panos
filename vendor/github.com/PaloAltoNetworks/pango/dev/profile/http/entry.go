package http

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of an http profile.
//
// PAN-OS 7.1+.
type Entry struct {
	Name              string
	TagRegistration   bool
	ConfigName        string
	ConfigUriFormat   string
	ConfigPayload     string
	SystemName        string
	SystemUriFormat   string
	SystemPayload     string
	ThreatName        string
	ThreatUriFormat   string
	ThreatPayload     string
	TrafficName       string
	TrafficUriFormat  string
	TrafficPayload    string
	HipMatchName      string
	HipMatchUriFormat string
	HipMatchPayload   string
	UrlName           string
	UrlUriFormat      string
	UrlPayload        string
	DataName          string
	DataUriFormat     string
	DataPayload       string
	WildfireName      string
	WildfireUriFormat string
	WildfirePayload   string
	TunnelName        string
	TunnelUriFormat   string
	TunnelPayload     string
	UserIdName        string
	UserIdUriFormat   string
	UserIdPayload     string
	GtpName           string
	GtpUriFormat      string
	GtpPayload        string
	AuthName          string
	AuthUriFormat     string
	AuthPayload       string
	SctpName          string // 8.1+
	SctpUriFormat     string // 8.1+
	SctpPayload       string // 8.1+
	IptagName         string // 9.0+
	IptagUriFormat    string // 9.0+
	IptagPayload      string // 9.0+

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.TagRegistration = s.TagRegistration
	o.ConfigName = s.ConfigName
	o.ConfigUriFormat = s.ConfigUriFormat
	o.ConfigPayload = s.ConfigPayload
	o.SystemName = s.SystemName
	o.SystemUriFormat = s.SystemUriFormat
	o.SystemPayload = s.SystemPayload
	o.ThreatName = s.ThreatName
	o.ThreatUriFormat = s.ThreatUriFormat
	o.ThreatPayload = s.ThreatPayload
	o.TrafficName = s.TrafficName
	o.TrafficUriFormat = s.TrafficUriFormat
	o.TrafficPayload = s.TrafficPayload
	o.HipMatchName = s.HipMatchName
	o.HipMatchUriFormat = s.HipMatchUriFormat
	o.HipMatchPayload = s.HipMatchPayload
	o.UrlName = s.UrlName
	o.UrlUriFormat = s.UrlUriFormat
	o.UrlPayload = s.UrlPayload
	o.DataName = s.DataName
	o.DataUriFormat = s.DataUriFormat
	o.DataPayload = s.DataPayload
	o.WildfireName = s.WildfireName
	o.WildfireUriFormat = s.WildfireUriFormat
	o.WildfirePayload = s.WildfirePayload
	o.TunnelName = s.TunnelName
	o.TunnelUriFormat = s.TunnelUriFormat
	o.TunnelPayload = s.TunnelPayload
	o.UserIdName = s.UserIdName
	o.UserIdUriFormat = s.UserIdUriFormat
	o.UserIdPayload = s.UserIdPayload
	o.GtpName = s.GtpName
	o.GtpUriFormat = s.GtpUriFormat
	o.GtpPayload = s.GtpPayload
	o.AuthName = s.AuthName
	o.AuthUriFormat = s.AuthUriFormat
	o.AuthPayload = s.AuthPayload
	o.SctpName = s.SctpName
	o.SctpUriFormat = s.SctpUriFormat
	o.SctpPayload = s.SctpPayload
	o.IptagName = s.IptagName
	o.IptagUriFormat = s.IptagUriFormat
	o.IptagPayload = s.IptagPayload
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
		Name:            o.Answer.Name,
		TagRegistration: util.AsBool(o.Answer.TagRegistration),
	}

	ans.raw = make(map[string]string)

	if o.Answer.Server != nil {
		ans.raw["srv"] = util.CleanRawXml(o.Answer.Server.Text)
	}

	if o.Answer.Format != nil {
		if o.Answer.Format.Config != nil {
			ans.ConfigName = o.Answer.Format.Config.Name
			ans.ConfigUriFormat = o.Answer.Format.Config.UriFormat
			ans.ConfigPayload = o.Answer.Format.Config.Payload
			if o.Answer.Format.Config.Headers != nil {
				ans.raw["configh"] = util.CleanRawXml(o.Answer.Format.Config.Headers.Text)
			}
			if o.Answer.Format.Config.Params != nil {
				ans.raw["configp"] = util.CleanRawXml(o.Answer.Format.Config.Params.Text)
			}
		}
		if o.Answer.Format.System != nil {
			ans.SystemName = o.Answer.Format.System.Name
			ans.SystemUriFormat = o.Answer.Format.System.UriFormat
			ans.SystemPayload = o.Answer.Format.System.Payload
			if o.Answer.Format.System.Headers != nil {
				ans.raw["systemh"] = util.CleanRawXml(o.Answer.Format.System.Headers.Text)
			}
			if o.Answer.Format.System.Params != nil {
				ans.raw["systemp"] = util.CleanRawXml(o.Answer.Format.System.Params.Text)
			}
		}
		if o.Answer.Format.Threat != nil {
			ans.ThreatName = o.Answer.Format.Threat.Name
			ans.ThreatUriFormat = o.Answer.Format.Threat.UriFormat
			ans.ThreatPayload = o.Answer.Format.Threat.Payload
			if o.Answer.Format.Threat.Headers != nil {
				ans.raw["threath"] = util.CleanRawXml(o.Answer.Format.Threat.Headers.Text)
			}
			if o.Answer.Format.Threat.Params != nil {
				ans.raw["threatp"] = util.CleanRawXml(o.Answer.Format.Threat.Params.Text)
			}
		}
		if o.Answer.Format.Traffic != nil {
			ans.TrafficName = o.Answer.Format.Traffic.Name
			ans.TrafficUriFormat = o.Answer.Format.Traffic.UriFormat
			ans.TrafficPayload = o.Answer.Format.Traffic.Payload
			if o.Answer.Format.Traffic.Headers != nil {
				ans.raw["traffich"] = util.CleanRawXml(o.Answer.Format.Traffic.Headers.Text)
			}
			if o.Answer.Format.Traffic.Params != nil {
				ans.raw["trafficp"] = util.CleanRawXml(o.Answer.Format.Traffic.Params.Text)
			}
		}
		if o.Answer.Format.HipMatch != nil {
			ans.HipMatchName = o.Answer.Format.HipMatch.Name
			ans.HipMatchUriFormat = o.Answer.Format.HipMatch.UriFormat
			ans.HipMatchPayload = o.Answer.Format.HipMatch.Payload
			if o.Answer.Format.HipMatch.Headers != nil {
				ans.raw["hipmatchh"] = util.CleanRawXml(o.Answer.Format.HipMatch.Headers.Text)
			}
			if o.Answer.Format.HipMatch.Params != nil {
				ans.raw["hipmatchp"] = util.CleanRawXml(o.Answer.Format.HipMatch.Params.Text)
			}
		}
		if o.Answer.Format.Url != nil {
			ans.UrlName = o.Answer.Format.Url.Name
			ans.UrlUriFormat = o.Answer.Format.Url.UriFormat
			ans.UrlPayload = o.Answer.Format.Url.Payload
			if o.Answer.Format.Url.Headers != nil {
				ans.raw["urlh"] = util.CleanRawXml(o.Answer.Format.Url.Headers.Text)
			}
			if o.Answer.Format.Url.Params != nil {
				ans.raw["urlp"] = util.CleanRawXml(o.Answer.Format.Url.Params.Text)
			}
		}
		if o.Answer.Format.Data != nil {
			ans.DataName = o.Answer.Format.Data.Name
			ans.DataUriFormat = o.Answer.Format.Data.UriFormat
			ans.DataPayload = o.Answer.Format.Data.Payload
			if o.Answer.Format.Data.Headers != nil {
				ans.raw["datah"] = util.CleanRawXml(o.Answer.Format.Data.Headers.Text)
			}
			if o.Answer.Format.Data.Params != nil {
				ans.raw["datap"] = util.CleanRawXml(o.Answer.Format.Data.Params.Text)
			}
		}
		if o.Answer.Format.Wildfire != nil {
			ans.WildfireName = o.Answer.Format.Wildfire.Name
			ans.WildfireUriFormat = o.Answer.Format.Wildfire.UriFormat
			ans.WildfirePayload = o.Answer.Format.Wildfire.Payload
			if o.Answer.Format.Wildfire.Headers != nil {
				ans.raw["wildfireh"] = util.CleanRawXml(o.Answer.Format.Wildfire.Headers.Text)
			}
			if o.Answer.Format.Wildfire.Params != nil {
				ans.raw["wildfirep"] = util.CleanRawXml(o.Answer.Format.Wildfire.Params.Text)
			}
		}
		if o.Answer.Format.Tunnel != nil {
			ans.TunnelName = o.Answer.Format.Tunnel.Name
			ans.TunnelUriFormat = o.Answer.Format.Tunnel.UriFormat
			ans.TunnelPayload = o.Answer.Format.Tunnel.Payload
			if o.Answer.Format.Tunnel.Headers != nil {
				ans.raw["tunnelh"] = util.CleanRawXml(o.Answer.Format.Tunnel.Headers.Text)
			}
			if o.Answer.Format.Tunnel.Params != nil {
				ans.raw["tunnelp"] = util.CleanRawXml(o.Answer.Format.Tunnel.Params.Text)
			}
		}
		if o.Answer.Format.UserId != nil {
			ans.UserIdName = o.Answer.Format.UserId.Name
			ans.UserIdUriFormat = o.Answer.Format.UserId.UriFormat
			ans.UserIdPayload = o.Answer.Format.UserId.Payload
			if o.Answer.Format.UserId.Headers != nil {
				ans.raw["useridh"] = util.CleanRawXml(o.Answer.Format.UserId.Headers.Text)
			}
			if o.Answer.Format.UserId.Params != nil {
				ans.raw["useridp"] = util.CleanRawXml(o.Answer.Format.UserId.Params.Text)
			}
		}
		if o.Answer.Format.Gtp != nil {
			ans.GtpName = o.Answer.Format.Gtp.Name
			ans.GtpUriFormat = o.Answer.Format.Gtp.UriFormat
			ans.GtpPayload = o.Answer.Format.Gtp.Payload
			if o.Answer.Format.Gtp.Headers != nil {
				ans.raw["gtph"] = util.CleanRawXml(o.Answer.Format.Gtp.Headers.Text)
			}
			if o.Answer.Format.Gtp.Params != nil {
				ans.raw["gtpp"] = util.CleanRawXml(o.Answer.Format.Gtp.Params.Text)
			}
		}
		if o.Answer.Format.Auth != nil {
			ans.AuthName = o.Answer.Format.Auth.Name
			ans.AuthUriFormat = o.Answer.Format.Auth.UriFormat
			ans.AuthPayload = o.Answer.Format.Auth.Payload
			if o.Answer.Format.Auth.Headers != nil {
				ans.raw["authh"] = util.CleanRawXml(o.Answer.Format.Auth.Headers.Text)
			}
			if o.Answer.Format.Auth.Params != nil {
				ans.raw["authp"] = util.CleanRawXml(o.Answer.Format.Auth.Params.Text)
			}
		}
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type container_v2 struct {
	Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
	ans := Entry{
		Name:            o.Answer.Name,
		TagRegistration: util.AsBool(o.Answer.TagRegistration),
	}

	ans.raw = make(map[string]string)

	if o.Answer.Server != nil {
		ans.raw["srv"] = util.CleanRawXml(o.Answer.Server.Text)
	}

	if o.Answer.Format != nil {
		if o.Answer.Format.Config != nil {
			ans.ConfigName = o.Answer.Format.Config.Name
			ans.ConfigUriFormat = o.Answer.Format.Config.UriFormat
			ans.ConfigPayload = o.Answer.Format.Config.Payload
			if o.Answer.Format.Config.Headers != nil {
				ans.raw["configh"] = util.CleanRawXml(o.Answer.Format.Config.Headers.Text)
			}
			if o.Answer.Format.Config.Params != nil {
				ans.raw["configp"] = util.CleanRawXml(o.Answer.Format.Config.Params.Text)
			}
		}
		if o.Answer.Format.System != nil {
			ans.SystemName = o.Answer.Format.System.Name
			ans.SystemUriFormat = o.Answer.Format.System.UriFormat
			ans.SystemPayload = o.Answer.Format.System.Payload
			if o.Answer.Format.System.Headers != nil {
				ans.raw["systemh"] = util.CleanRawXml(o.Answer.Format.System.Headers.Text)
			}
			if o.Answer.Format.System.Params != nil {
				ans.raw["systemp"] = util.CleanRawXml(o.Answer.Format.System.Params.Text)
			}
		}
		if o.Answer.Format.Threat != nil {
			ans.ThreatName = o.Answer.Format.Threat.Name
			ans.ThreatUriFormat = o.Answer.Format.Threat.UriFormat
			ans.ThreatPayload = o.Answer.Format.Threat.Payload
			if o.Answer.Format.Threat.Headers != nil {
				ans.raw["threath"] = util.CleanRawXml(o.Answer.Format.Threat.Headers.Text)
			}
			if o.Answer.Format.Threat.Params != nil {
				ans.raw["threatp"] = util.CleanRawXml(o.Answer.Format.Threat.Params.Text)
			}
		}
		if o.Answer.Format.Traffic != nil {
			ans.TrafficName = o.Answer.Format.Traffic.Name
			ans.TrafficUriFormat = o.Answer.Format.Traffic.UriFormat
			ans.TrafficPayload = o.Answer.Format.Traffic.Payload
			if o.Answer.Format.Traffic.Headers != nil {
				ans.raw["traffich"] = util.CleanRawXml(o.Answer.Format.Traffic.Headers.Text)
			}
			if o.Answer.Format.Traffic.Params != nil {
				ans.raw["trafficp"] = util.CleanRawXml(o.Answer.Format.Traffic.Params.Text)
			}
		}
		if o.Answer.Format.HipMatch != nil {
			ans.HipMatchName = o.Answer.Format.HipMatch.Name
			ans.HipMatchUriFormat = o.Answer.Format.HipMatch.UriFormat
			ans.HipMatchPayload = o.Answer.Format.HipMatch.Payload
			if o.Answer.Format.HipMatch.Headers != nil {
				ans.raw["hipmatchh"] = util.CleanRawXml(o.Answer.Format.HipMatch.Headers.Text)
			}
			if o.Answer.Format.HipMatch.Params != nil {
				ans.raw["hipmatchp"] = util.CleanRawXml(o.Answer.Format.HipMatch.Params.Text)
			}
		}
		if o.Answer.Format.Url != nil {
			ans.UrlName = o.Answer.Format.Url.Name
			ans.UrlUriFormat = o.Answer.Format.Url.UriFormat
			ans.UrlPayload = o.Answer.Format.Url.Payload
			if o.Answer.Format.Url.Headers != nil {
				ans.raw["urlh"] = util.CleanRawXml(o.Answer.Format.Url.Headers.Text)
			}
			if o.Answer.Format.Url.Params != nil {
				ans.raw["urlp"] = util.CleanRawXml(o.Answer.Format.Url.Params.Text)
			}
		}
		if o.Answer.Format.Data != nil {
			ans.DataName = o.Answer.Format.Data.Name
			ans.DataUriFormat = o.Answer.Format.Data.UriFormat
			ans.DataPayload = o.Answer.Format.Data.Payload
			if o.Answer.Format.Data.Headers != nil {
				ans.raw["datah"] = util.CleanRawXml(o.Answer.Format.Data.Headers.Text)
			}
			if o.Answer.Format.Data.Params != nil {
				ans.raw["datap"] = util.CleanRawXml(o.Answer.Format.Data.Params.Text)
			}
		}
		if o.Answer.Format.Wildfire != nil {
			ans.WildfireName = o.Answer.Format.Wildfire.Name
			ans.WildfireUriFormat = o.Answer.Format.Wildfire.UriFormat
			ans.WildfirePayload = o.Answer.Format.Wildfire.Payload
			if o.Answer.Format.Wildfire.Headers != nil {
				ans.raw["wildfireh"] = util.CleanRawXml(o.Answer.Format.Wildfire.Headers.Text)
			}
			if o.Answer.Format.Wildfire.Params != nil {
				ans.raw["wildfirep"] = util.CleanRawXml(o.Answer.Format.Wildfire.Params.Text)
			}
		}
		if o.Answer.Format.Tunnel != nil {
			ans.TunnelName = o.Answer.Format.Tunnel.Name
			ans.TunnelUriFormat = o.Answer.Format.Tunnel.UriFormat
			ans.TunnelPayload = o.Answer.Format.Tunnel.Payload
			if o.Answer.Format.Tunnel.Headers != nil {
				ans.raw["tunnelh"] = util.CleanRawXml(o.Answer.Format.Tunnel.Headers.Text)
			}
			if o.Answer.Format.Tunnel.Params != nil {
				ans.raw["tunnelp"] = util.CleanRawXml(o.Answer.Format.Tunnel.Params.Text)
			}
		}
		if o.Answer.Format.UserId != nil {
			ans.UserIdName = o.Answer.Format.UserId.Name
			ans.UserIdUriFormat = o.Answer.Format.UserId.UriFormat
			ans.UserIdPayload = o.Answer.Format.UserId.Payload
			if o.Answer.Format.UserId.Headers != nil {
				ans.raw["useridh"] = util.CleanRawXml(o.Answer.Format.UserId.Headers.Text)
			}
			if o.Answer.Format.UserId.Params != nil {
				ans.raw["useridp"] = util.CleanRawXml(o.Answer.Format.UserId.Params.Text)
			}
		}
		if o.Answer.Format.Gtp != nil {
			ans.GtpName = o.Answer.Format.Gtp.Name
			ans.GtpUriFormat = o.Answer.Format.Gtp.UriFormat
			ans.GtpPayload = o.Answer.Format.Gtp.Payload
			if o.Answer.Format.Gtp.Headers != nil {
				ans.raw["gtph"] = util.CleanRawXml(o.Answer.Format.Gtp.Headers.Text)
			}
			if o.Answer.Format.Gtp.Params != nil {
				ans.raw["gtpp"] = util.CleanRawXml(o.Answer.Format.Gtp.Params.Text)
			}
		}
		if o.Answer.Format.Auth != nil {
			ans.AuthName = o.Answer.Format.Auth.Name
			ans.AuthUriFormat = o.Answer.Format.Auth.UriFormat
			ans.AuthPayload = o.Answer.Format.Auth.Payload
			if o.Answer.Format.Auth.Headers != nil {
				ans.raw["authh"] = util.CleanRawXml(o.Answer.Format.Auth.Headers.Text)
			}
			if o.Answer.Format.Auth.Params != nil {
				ans.raw["authp"] = util.CleanRawXml(o.Answer.Format.Auth.Params.Text)
			}
		}
		if o.Answer.Format.Sctp != nil {
			ans.SctpName = o.Answer.Format.Sctp.Name
			ans.SctpUriFormat = o.Answer.Format.Sctp.UriFormat
			ans.SctpPayload = o.Answer.Format.Sctp.Payload
			if o.Answer.Format.Sctp.Headers != nil {
				ans.raw["sctph"] = util.CleanRawXml(o.Answer.Format.Sctp.Headers.Text)
			}
			if o.Answer.Format.Sctp.Params != nil {
				ans.raw["sctpp"] = util.CleanRawXml(o.Answer.Format.Sctp.Params.Text)
			}
		}
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type container_v3 struct {
	Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
	ans := Entry{
		Name:            o.Answer.Name,
		TagRegistration: util.AsBool(o.Answer.TagRegistration),
	}

	ans.raw = make(map[string]string)

	if o.Answer.Server != nil {
		ans.raw["srv"] = util.CleanRawXml(o.Answer.Server.Text)
	}

	if o.Answer.Format != nil {
		if o.Answer.Format.Config != nil {
			ans.ConfigName = o.Answer.Format.Config.Name
			ans.ConfigUriFormat = o.Answer.Format.Config.UriFormat
			ans.ConfigPayload = o.Answer.Format.Config.Payload
			if o.Answer.Format.Config.Headers != nil {
				ans.raw["configh"] = util.CleanRawXml(o.Answer.Format.Config.Headers.Text)
			}
			if o.Answer.Format.Config.Params != nil {
				ans.raw["configp"] = util.CleanRawXml(o.Answer.Format.Config.Params.Text)
			}
		}
		if o.Answer.Format.System != nil {
			ans.SystemName = o.Answer.Format.System.Name
			ans.SystemUriFormat = o.Answer.Format.System.UriFormat
			ans.SystemPayload = o.Answer.Format.System.Payload
			if o.Answer.Format.System.Headers != nil {
				ans.raw["systemh"] = util.CleanRawXml(o.Answer.Format.System.Headers.Text)
			}
			if o.Answer.Format.System.Params != nil {
				ans.raw["systemp"] = util.CleanRawXml(o.Answer.Format.System.Params.Text)
			}
		}
		if o.Answer.Format.Threat != nil {
			ans.ThreatName = o.Answer.Format.Threat.Name
			ans.ThreatUriFormat = o.Answer.Format.Threat.UriFormat
			ans.ThreatPayload = o.Answer.Format.Threat.Payload
			if o.Answer.Format.Threat.Headers != nil {
				ans.raw["threath"] = util.CleanRawXml(o.Answer.Format.Threat.Headers.Text)
			}
			if o.Answer.Format.Threat.Params != nil {
				ans.raw["threatp"] = util.CleanRawXml(o.Answer.Format.Threat.Params.Text)
			}
		}
		if o.Answer.Format.Traffic != nil {
			ans.TrafficName = o.Answer.Format.Traffic.Name
			ans.TrafficUriFormat = o.Answer.Format.Traffic.UriFormat
			ans.TrafficPayload = o.Answer.Format.Traffic.Payload
			if o.Answer.Format.Traffic.Headers != nil {
				ans.raw["traffich"] = util.CleanRawXml(o.Answer.Format.Traffic.Headers.Text)
			}
			if o.Answer.Format.Traffic.Params != nil {
				ans.raw["trafficp"] = util.CleanRawXml(o.Answer.Format.Traffic.Params.Text)
			}
		}
		if o.Answer.Format.HipMatch != nil {
			ans.HipMatchName = o.Answer.Format.HipMatch.Name
			ans.HipMatchUriFormat = o.Answer.Format.HipMatch.UriFormat
			ans.HipMatchPayload = o.Answer.Format.HipMatch.Payload
			if o.Answer.Format.HipMatch.Headers != nil {
				ans.raw["hipmatchh"] = util.CleanRawXml(o.Answer.Format.HipMatch.Headers.Text)
			}
			if o.Answer.Format.HipMatch.Params != nil {
				ans.raw["hipmatchp"] = util.CleanRawXml(o.Answer.Format.HipMatch.Params.Text)
			}
		}
		if o.Answer.Format.Url != nil {
			ans.UrlName = o.Answer.Format.Url.Name
			ans.UrlUriFormat = o.Answer.Format.Url.UriFormat
			ans.UrlPayload = o.Answer.Format.Url.Payload
			if o.Answer.Format.Url.Headers != nil {
				ans.raw["urlh"] = util.CleanRawXml(o.Answer.Format.Url.Headers.Text)
			}
			if o.Answer.Format.Url.Params != nil {
				ans.raw["urlp"] = util.CleanRawXml(o.Answer.Format.Url.Params.Text)
			}
		}
		if o.Answer.Format.Data != nil {
			ans.DataName = o.Answer.Format.Data.Name
			ans.DataUriFormat = o.Answer.Format.Data.UriFormat
			ans.DataPayload = o.Answer.Format.Data.Payload
			if o.Answer.Format.Data.Headers != nil {
				ans.raw["datah"] = util.CleanRawXml(o.Answer.Format.Data.Headers.Text)
			}
			if o.Answer.Format.Data.Params != nil {
				ans.raw["datap"] = util.CleanRawXml(o.Answer.Format.Data.Params.Text)
			}
		}
		if o.Answer.Format.Wildfire != nil {
			ans.WildfireName = o.Answer.Format.Wildfire.Name
			ans.WildfireUriFormat = o.Answer.Format.Wildfire.UriFormat
			ans.WildfirePayload = o.Answer.Format.Wildfire.Payload
			if o.Answer.Format.Wildfire.Headers != nil {
				ans.raw["wildfireh"] = util.CleanRawXml(o.Answer.Format.Wildfire.Headers.Text)
			}
			if o.Answer.Format.Wildfire.Params != nil {
				ans.raw["wildfirep"] = util.CleanRawXml(o.Answer.Format.Wildfire.Params.Text)
			}
		}
		if o.Answer.Format.Tunnel != nil {
			ans.TunnelName = o.Answer.Format.Tunnel.Name
			ans.TunnelUriFormat = o.Answer.Format.Tunnel.UriFormat
			ans.TunnelPayload = o.Answer.Format.Tunnel.Payload
			if o.Answer.Format.Tunnel.Headers != nil {
				ans.raw["tunnelh"] = util.CleanRawXml(o.Answer.Format.Tunnel.Headers.Text)
			}
			if o.Answer.Format.Tunnel.Params != nil {
				ans.raw["tunnelp"] = util.CleanRawXml(o.Answer.Format.Tunnel.Params.Text)
			}
		}
		if o.Answer.Format.UserId != nil {
			ans.UserIdName = o.Answer.Format.UserId.Name
			ans.UserIdUriFormat = o.Answer.Format.UserId.UriFormat
			ans.UserIdPayload = o.Answer.Format.UserId.Payload
			if o.Answer.Format.UserId.Headers != nil {
				ans.raw["useridh"] = util.CleanRawXml(o.Answer.Format.UserId.Headers.Text)
			}
			if o.Answer.Format.UserId.Params != nil {
				ans.raw["useridp"] = util.CleanRawXml(o.Answer.Format.UserId.Params.Text)
			}
		}
		if o.Answer.Format.Gtp != nil {
			ans.GtpName = o.Answer.Format.Gtp.Name
			ans.GtpUriFormat = o.Answer.Format.Gtp.UriFormat
			ans.GtpPayload = o.Answer.Format.Gtp.Payload
			if o.Answer.Format.Gtp.Headers != nil {
				ans.raw["gtph"] = util.CleanRawXml(o.Answer.Format.Gtp.Headers.Text)
			}
			if o.Answer.Format.Gtp.Params != nil {
				ans.raw["gtpp"] = util.CleanRawXml(o.Answer.Format.Gtp.Params.Text)
			}
		}
		if o.Answer.Format.Auth != nil {
			ans.AuthName = o.Answer.Format.Auth.Name
			ans.AuthUriFormat = o.Answer.Format.Auth.UriFormat
			ans.AuthPayload = o.Answer.Format.Auth.Payload
			if o.Answer.Format.Auth.Headers != nil {
				ans.raw["authh"] = util.CleanRawXml(o.Answer.Format.Auth.Headers.Text)
			}
			if o.Answer.Format.Auth.Params != nil {
				ans.raw["authp"] = util.CleanRawXml(o.Answer.Format.Auth.Params.Text)
			}
		}
		if o.Answer.Format.Sctp != nil {
			ans.SctpName = o.Answer.Format.Sctp.Name
			ans.SctpUriFormat = o.Answer.Format.Sctp.UriFormat
			ans.SctpPayload = o.Answer.Format.Sctp.Payload
			if o.Answer.Format.Sctp.Headers != nil {
				ans.raw["sctph"] = util.CleanRawXml(o.Answer.Format.Sctp.Headers.Text)
			}
			if o.Answer.Format.Sctp.Params != nil {
				ans.raw["sctpp"] = util.CleanRawXml(o.Answer.Format.Sctp.Params.Text)
			}
		}
		if o.Answer.Format.Iptag != nil {
			ans.IptagName = o.Answer.Format.Iptag.Name
			ans.IptagUriFormat = o.Answer.Format.Iptag.UriFormat
			ans.IptagPayload = o.Answer.Format.Iptag.Payload
			if o.Answer.Format.Iptag.Headers != nil {
				ans.raw["iptagh"] = util.CleanRawXml(o.Answer.Format.Iptag.Headers.Text)
			}
			if o.Answer.Format.Iptag.Params != nil {
				ans.raw["iptagp"] = util.CleanRawXml(o.Answer.Format.Iptag.Params.Text)
			}
		}
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName         xml.Name     `xml:"entry"`
	Name            string       `xml:"name,attr"`
	TagRegistration string       `xml:"tag-registration"`
	Format          *format_v1   `xml:"format"`
	Server          *util.RawXml `xml:"server"`
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
	Name      string       `xml:"name,omitempty"`
	UriFormat string       `xml:"url-format,omitempty"`
	Payload   string       `xml:"payload,omitempty"`
	Headers   *util.RawXml `xml:"headers"`
	Params    *util.RawXml `xml:"params"`
}

func specify_v1(e Entry) interface{} {
	var hdr, prm string

	ans := entry_v1{
		Name:            e.Name,
		TagRegistration: util.YesNo(e.TagRegistration),
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	f := format_v1{}
	hasData := false

	hdr = e.raw["configh"]
	prm = e.raw["configp"]
	if e.ConfigName != "" || e.ConfigUriFormat != "" || e.ConfigPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Config = &formatSpec{
			Name:      e.ConfigName,
			UriFormat: e.ConfigUriFormat,
			Payload:   e.ConfigPayload,
		}
		if hdr != "" {
			f.Config.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Config.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["systemh"]
	prm = e.raw["systemp"]
	if e.SystemName != "" || e.SystemUriFormat != "" || e.SystemPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.System = &formatSpec{
			Name:      e.SystemName,
			UriFormat: e.SystemUriFormat,
			Payload:   e.SystemPayload,
		}
		if hdr != "" {
			f.System.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.System.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["threath"]
	prm = e.raw["threatp"]
	if e.ThreatName != "" || e.ThreatUriFormat != "" || e.ThreatPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Threat = &formatSpec{
			Name:      e.ThreatName,
			UriFormat: e.ThreatUriFormat,
			Payload:   e.ThreatPayload,
		}
		if hdr != "" {
			f.Threat.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Threat.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["traffich"]
	prm = e.raw["trafficp"]
	if e.TrafficName != "" || e.TrafficUriFormat != "" || e.TrafficPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Traffic = &formatSpec{
			Name:      e.TrafficName,
			UriFormat: e.TrafficUriFormat,
			Payload:   e.TrafficPayload,
		}
		if hdr != "" {
			f.Traffic.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Traffic.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["hipmatchh"]
	prm = e.raw["hipmatchp"]
	if e.HipMatchName != "" || e.HipMatchUriFormat != "" || e.HipMatchPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.HipMatch = &formatSpec{
			Name:      e.HipMatchName,
			UriFormat: e.HipMatchUriFormat,
			Payload:   e.HipMatchPayload,
		}
		if hdr != "" {
			f.HipMatch.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.HipMatch.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["urlh"]
	prm = e.raw["urlp"]
	if e.UrlName != "" || e.UrlUriFormat != "" || e.UrlPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Url = &formatSpec{
			Name:      e.UrlName,
			UriFormat: e.UrlUriFormat,
			Payload:   e.UrlPayload,
		}
		if hdr != "" {
			f.Url.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Url.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["datah"]
	prm = e.raw["datap"]
	if e.DataName != "" || e.DataUriFormat != "" || e.DataPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Data = &formatSpec{
			Name:      e.DataName,
			UriFormat: e.DataUriFormat,
			Payload:   e.DataPayload,
		}
		if hdr != "" {
			f.Data.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Data.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["wildfireh"]
	prm = e.raw["wildfirep"]
	if e.WildfireName != "" || e.WildfireUriFormat != "" || e.WildfirePayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Wildfire = &formatSpec{
			Name:      e.WildfireName,
			UriFormat: e.WildfireUriFormat,
			Payload:   e.WildfirePayload,
		}
		if hdr != "" {
			f.Wildfire.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Wildfire.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["tunnelh"]
	prm = e.raw["tunnelp"]
	if e.TunnelName != "" || e.TunnelUriFormat != "" || e.TunnelPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Tunnel = &formatSpec{
			Name:      e.TunnelName,
			UriFormat: e.TunnelUriFormat,
			Payload:   e.TunnelPayload,
		}
		if hdr != "" {
			f.Tunnel.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Tunnel.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["useridh"]
	prm = e.raw["useridp"]
	if e.UserIdName != "" || e.UserIdUriFormat != "" || e.UserIdPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.UserId = &formatSpec{
			Name:      e.UserIdName,
			UriFormat: e.UserIdUriFormat,
			Payload:   e.UserIdPayload,
		}
		if hdr != "" {
			f.UserId.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.UserId.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["gtph"]
	prm = e.raw["gtpp"]
	if e.GtpName != "" || e.GtpUriFormat != "" || e.GtpPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Gtp = &formatSpec{
			Name:      e.GtpName,
			UriFormat: e.GtpUriFormat,
			Payload:   e.GtpPayload,
		}
		if hdr != "" {
			f.Gtp.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Gtp.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["authh"]
	prm = e.raw["authp"]
	if e.AuthName != "" || e.AuthUriFormat != "" || e.AuthPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Auth = &formatSpec{
			Name:      e.AuthName,
			UriFormat: e.AuthUriFormat,
			Payload:   e.AuthPayload,
		}
		if hdr != "" {
			f.Auth.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Auth.Params = &util.RawXml{prm}
		}
	}

	if hasData {
		ans.Format = &f
	}

	return ans
}

type entry_v2 struct {
	XMLName         xml.Name     `xml:"entry"`
	Name            string       `xml:"name,attr"`
	TagRegistration string       `xml:"tag-registration"`
	Format          *format_v2   `xml:"format"`
	Server          *util.RawXml `xml:"server"`
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
	var hdr, prm string

	ans := entry_v2{
		Name:            e.Name,
		TagRegistration: util.YesNo(e.TagRegistration),
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	f := format_v2{}
	hasData := false

	hdr = e.raw["configh"]
	prm = e.raw["configp"]
	if e.ConfigName != "" || e.ConfigUriFormat != "" || e.ConfigPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Config = &formatSpec{
			Name:      e.ConfigName,
			UriFormat: e.ConfigUriFormat,
			Payload:   e.ConfigPayload,
		}
		if hdr != "" {
			f.Config.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Config.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["systemh"]
	prm = e.raw["systemp"]
	if e.SystemName != "" || e.SystemUriFormat != "" || e.SystemPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.System = &formatSpec{
			Name:      e.SystemName,
			UriFormat: e.SystemUriFormat,
			Payload:   e.SystemPayload,
		}
		if hdr != "" {
			f.System.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.System.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["threath"]
	prm = e.raw["threatp"]
	if e.ThreatName != "" || e.ThreatUriFormat != "" || e.ThreatPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Threat = &formatSpec{
			Name:      e.ThreatName,
			UriFormat: e.ThreatUriFormat,
			Payload:   e.ThreatPayload,
		}
		if hdr != "" {
			f.Threat.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Threat.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["traffich"]
	prm = e.raw["trafficp"]
	if e.TrafficName != "" || e.TrafficUriFormat != "" || e.TrafficPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Traffic = &formatSpec{
			Name:      e.TrafficName,
			UriFormat: e.TrafficUriFormat,
			Payload:   e.TrafficPayload,
		}
		if hdr != "" {
			f.Traffic.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Traffic.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["hipmatchh"]
	prm = e.raw["hipmatchp"]
	if e.HipMatchName != "" || e.HipMatchUriFormat != "" || e.HipMatchPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.HipMatch = &formatSpec{
			Name:      e.HipMatchName,
			UriFormat: e.HipMatchUriFormat,
			Payload:   e.HipMatchPayload,
		}
		if hdr != "" {
			f.HipMatch.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.HipMatch.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["urlh"]
	prm = e.raw["urlp"]
	if e.UrlName != "" || e.UrlUriFormat != "" || e.UrlPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Url = &formatSpec{
			Name:      e.UrlName,
			UriFormat: e.UrlUriFormat,
			Payload:   e.UrlPayload,
		}
		if hdr != "" {
			f.Url.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Url.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["datah"]
	prm = e.raw["datap"]
	if e.DataName != "" || e.DataUriFormat != "" || e.DataPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Data = &formatSpec{
			Name:      e.DataName,
			UriFormat: e.DataUriFormat,
			Payload:   e.DataPayload,
		}
		if hdr != "" {
			f.Data.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Data.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["wildfireh"]
	prm = e.raw["wildfirep"]
	if e.WildfireName != "" || e.WildfireUriFormat != "" || e.WildfirePayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Wildfire = &formatSpec{
			Name:      e.WildfireName,
			UriFormat: e.WildfireUriFormat,
			Payload:   e.WildfirePayload,
		}
		if hdr != "" {
			f.Wildfire.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Wildfire.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["tunnelh"]
	prm = e.raw["tunnelp"]
	if e.TunnelName != "" || e.TunnelUriFormat != "" || e.TunnelPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Tunnel = &formatSpec{
			Name:      e.TunnelName,
			UriFormat: e.TunnelUriFormat,
			Payload:   e.TunnelPayload,
		}
		if hdr != "" {
			f.Tunnel.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Tunnel.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["useridh"]
	prm = e.raw["useridp"]
	if e.UserIdName != "" || e.UserIdUriFormat != "" || e.UserIdPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.UserId = &formatSpec{
			Name:      e.UserIdName,
			UriFormat: e.UserIdUriFormat,
			Payload:   e.UserIdPayload,
		}
		if hdr != "" {
			f.UserId.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.UserId.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["gtph"]
	prm = e.raw["gtpp"]
	if e.GtpName != "" || e.GtpUriFormat != "" || e.GtpPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Gtp = &formatSpec{
			Name:      e.GtpName,
			UriFormat: e.GtpUriFormat,
			Payload:   e.GtpPayload,
		}
		if hdr != "" {
			f.Gtp.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Gtp.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["authh"]
	prm = e.raw["authp"]
	if e.AuthName != "" || e.AuthUriFormat != "" || e.AuthPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Auth = &formatSpec{
			Name:      e.AuthName,
			UriFormat: e.AuthUriFormat,
			Payload:   e.AuthPayload,
		}
		if hdr != "" {
			f.Auth.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Auth.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["sctph"]
	prm = e.raw["sctpp"]
	if e.SctpName != "" || e.SctpUriFormat != "" || e.SctpPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Sctp = &formatSpec{
			Name:      e.SctpName,
			UriFormat: e.SctpUriFormat,
			Payload:   e.SctpPayload,
		}
		if hdr != "" {
			f.Sctp.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Sctp.Params = &util.RawXml{prm}
		}
	}

	if hasData {
		ans.Format = &f
	}

	return ans
}

type entry_v3 struct {
	XMLName         xml.Name     `xml:"entry"`
	Name            string       `xml:"name,attr"`
	TagRegistration string       `xml:"tag-registration"`
	Format          *format_v3   `xml:"format"`
	Server          *util.RawXml `xml:"server"`
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
	var hdr, prm string

	ans := entry_v3{
		Name:            e.Name,
		TagRegistration: util.YesNo(e.TagRegistration),
	}

	if text := e.raw["srv"]; text != "" {
		ans.Server = &util.RawXml{text}
	}

	f := format_v3{}
	hasData := false

	hdr = e.raw["configh"]
	prm = e.raw["configp"]
	if e.ConfigName != "" || e.ConfigUriFormat != "" || e.ConfigPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Config = &formatSpec{
			Name:      e.ConfigName,
			UriFormat: e.ConfigUriFormat,
			Payload:   e.ConfigPayload,
		}
		if hdr != "" {
			f.Config.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Config.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["systemh"]
	prm = e.raw["systemp"]
	if e.SystemName != "" || e.SystemUriFormat != "" || e.SystemPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.System = &formatSpec{
			Name:      e.SystemName,
			UriFormat: e.SystemUriFormat,
			Payload:   e.SystemPayload,
		}
		if hdr != "" {
			f.System.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.System.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["threath"]
	prm = e.raw["threatp"]
	if e.ThreatName != "" || e.ThreatUriFormat != "" || e.ThreatPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Threat = &formatSpec{
			Name:      e.ThreatName,
			UriFormat: e.ThreatUriFormat,
			Payload:   e.ThreatPayload,
		}
		if hdr != "" {
			f.Threat.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Threat.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["traffich"]
	prm = e.raw["trafficp"]
	if e.TrafficName != "" || e.TrafficUriFormat != "" || e.TrafficPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Traffic = &formatSpec{
			Name:      e.TrafficName,
			UriFormat: e.TrafficUriFormat,
			Payload:   e.TrafficPayload,
		}
		if hdr != "" {
			f.Traffic.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Traffic.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["hipmatchh"]
	prm = e.raw["hipmatchp"]
	if e.HipMatchName != "" || e.HipMatchUriFormat != "" || e.HipMatchPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.HipMatch = &formatSpec{
			Name:      e.HipMatchName,
			UriFormat: e.HipMatchUriFormat,
			Payload:   e.HipMatchPayload,
		}
		if hdr != "" {
			f.HipMatch.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.HipMatch.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["urlh"]
	prm = e.raw["urlp"]
	if e.UrlName != "" || e.UrlUriFormat != "" || e.UrlPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Url = &formatSpec{
			Name:      e.UrlName,
			UriFormat: e.UrlUriFormat,
			Payload:   e.UrlPayload,
		}
		if hdr != "" {
			f.Url.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Url.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["datah"]
	prm = e.raw["datap"]
	if e.DataName != "" || e.DataUriFormat != "" || e.DataPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Data = &formatSpec{
			Name:      e.DataName,
			UriFormat: e.DataUriFormat,
			Payload:   e.DataPayload,
		}
		if hdr != "" {
			f.Data.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Data.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["wildfireh"]
	prm = e.raw["wildfirep"]
	if e.WildfireName != "" || e.WildfireUriFormat != "" || e.WildfirePayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Wildfire = &formatSpec{
			Name:      e.WildfireName,
			UriFormat: e.WildfireUriFormat,
			Payload:   e.WildfirePayload,
		}
		if hdr != "" {
			f.Wildfire.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Wildfire.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["tunnelh"]
	prm = e.raw["tunnelp"]
	if e.TunnelName != "" || e.TunnelUriFormat != "" || e.TunnelPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Tunnel = &formatSpec{
			Name:      e.TunnelName,
			UriFormat: e.TunnelUriFormat,
			Payload:   e.TunnelPayload,
		}
		if hdr != "" {
			f.Tunnel.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Tunnel.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["useridh"]
	prm = e.raw["useridp"]
	if e.UserIdName != "" || e.UserIdUriFormat != "" || e.UserIdPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.UserId = &formatSpec{
			Name:      e.UserIdName,
			UriFormat: e.UserIdUriFormat,
			Payload:   e.UserIdPayload,
		}
		if hdr != "" {
			f.UserId.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.UserId.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["gtph"]
	prm = e.raw["gtpp"]
	if e.GtpName != "" || e.GtpUriFormat != "" || e.GtpPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Gtp = &formatSpec{
			Name:      e.GtpName,
			UriFormat: e.GtpUriFormat,
			Payload:   e.GtpPayload,
		}
		if hdr != "" {
			f.Gtp.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Gtp.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["authh"]
	prm = e.raw["authp"]
	if e.AuthName != "" || e.AuthUriFormat != "" || e.AuthPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Auth = &formatSpec{
			Name:      e.AuthName,
			UriFormat: e.AuthUriFormat,
			Payload:   e.AuthPayload,
		}
		if hdr != "" {
			f.Auth.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Auth.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["sctph"]
	prm = e.raw["sctpp"]
	if e.SctpName != "" || e.SctpUriFormat != "" || e.SctpPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Sctp = &formatSpec{
			Name:      e.SctpName,
			UriFormat: e.SctpUriFormat,
			Payload:   e.SctpPayload,
		}
		if hdr != "" {
			f.Sctp.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Sctp.Params = &util.RawXml{prm}
		}
	}

	hdr = e.raw["iptagh"]
	prm = e.raw["iptagp"]
	if e.IptagName != "" || e.IptagUriFormat != "" || e.IptagPayload != "" || hdr != "" || prm != "" {
		hasData = true
		f.Iptag = &formatSpec{
			Name:      e.IptagName,
			UriFormat: e.IptagUriFormat,
			Payload:   e.IptagPayload,
		}
		if hdr != "" {
			f.Iptag.Headers = &util.RawXml{hdr}
		}
		if prm != "" {
			f.Iptag.Params = &util.RawXml{prm}
		}
	}

	if hasData {
		ans.Format = &f
	}

	return ans
}
