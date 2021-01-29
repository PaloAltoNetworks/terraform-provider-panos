package bgp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Config is a normalized, version independent representation of a virtual
// router's BGP configuration.
type Config struct {
	Enable                        bool
	RouterId                      string
	AsNumber                      string // XML: local-as
	BfdProfile                    string // 7.1+ ; XML: global-bfd/profile or the word "None"
	RejectDefaultRoute            bool
	InstallRoute                  bool
	AggregateMed                  bool
	DefaultLocalPreference        string
	AsFormat                      string
	AlwaysCompareMed              bool
	DeterministicMedComparison    bool
	EcmpMultiAs                   bool // 7.0+
	EnforceFirstAs                bool // 8.0+
	EnableGracefulRestart         bool
	StaleRouteTime                int
	LocalRestartTime              int
	MaxPeerRestartTime            int
	ReflectorClusterId            string
	ConfederationMemberAs         string
	AllowRedistributeDefaultRoute bool

	raw map[string]string
}

// Copy copies the information from source Config `s` to this object.
func (o *Config) Copy(s Config) {
	o.Enable = s.Enable
	o.RouterId = s.RouterId
	o.AsNumber = s.AsNumber
	o.BfdProfile = s.BfdProfile
	o.RejectDefaultRoute = s.RejectDefaultRoute
	o.InstallRoute = s.InstallRoute
	o.AggregateMed = s.AggregateMed
	o.DefaultLocalPreference = s.DefaultLocalPreference
	o.AsFormat = s.AsFormat
	o.AlwaysCompareMed = s.AlwaysCompareMed
	o.DeterministicMedComparison = s.DeterministicMedComparison
	o.EcmpMultiAs = s.EcmpMultiAs
	o.EnforceFirstAs = s.EnforceFirstAs
	o.EnableGracefulRestart = s.EnableGracefulRestart
	o.StaleRouteTime = s.StaleRouteTime
	o.LocalRestartTime = s.LocalRestartTime
	o.MaxPeerRestartTime = s.MaxPeerRestartTime
	o.ReflectorClusterId = s.ReflectorClusterId
	o.ConfederationMemberAs = s.ConfederationMemberAs
	o.AllowRedistributeDefaultRoute = s.AllowRedistributeDefaultRoute
}

/** Structs / functions for this namespace. **/

func (o Config) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return "", fn(o)
}

type normalizer interface {
	Normalize() []Config
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"bgp"`
}

func (o *container_v1) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	return nil
}

func (o *entry_v1) normalize() Config {
	ans := Config{
		Enable:                        util.AsBool(o.Enable),
		RouterId:                      o.RouterId,
		AsNumber:                      o.AsNumber,
		RejectDefaultRoute:            util.AsBool(o.RejectDefaultRoute),
		InstallRoute:                  util.AsBool(o.InstallRoute),
		AllowRedistributeDefaultRoute: util.AsBool(o.AllowRedistributeDefaultRoute),
	}

	raw := make(map[string]string)

	if o.Options != nil {
		ans.AsFormat = o.Options.AsFormat
		ans.DefaultLocalPreference = o.Options.DefaultLocalPreference
		ans.ReflectorClusterId = o.Options.ReflectorClusterId
		ans.ConfederationMemberAs = o.Options.ConfederationMemberAs

		if o.Options.Med != nil {
			ans.AlwaysCompareMed = util.AsBool(o.Options.Med.AlwaysCompareMed)
			ans.DeterministicMedComparison = util.AsBool(o.Options.Med.DeterministicMedComparison)
		}

		if o.Options.GracefulRestart != nil {
			ans.EnableGracefulRestart = util.AsBool(o.Options.GracefulRestart.EnableGracefulRestart)
			ans.StaleRouteTime = o.Options.GracefulRestart.StaleRouteTime
			ans.LocalRestartTime = o.Options.GracefulRestart.LocalRestartTime
			ans.MaxPeerRestartTime = o.Options.GracefulRestart.MaxPeerRestartTime
		}

		if o.Options.Aggregate != nil {
			ans.AggregateMed = util.AsBool(o.Options.Aggregate.AggregateMed)
		}

		if o.Options.OutboundRouteFilter != nil {
			raw["orf"] = util.CleanRawXml(o.Options.OutboundRouteFilter.Text)
		}
	}

	if o.AuthProfile != nil {
		raw["ap"] = util.CleanRawXml(o.AuthProfile.Text)
	}
	if o.DampeningProfile != nil {
		raw["dp"] = util.CleanRawXml(o.DampeningProfile.Text)
	}
	if o.PeerGroup != nil {
		raw["pg"] = util.CleanRawXml(o.PeerGroup.Text)
	}
	if o.Policy != nil {
		raw["poli"] = util.CleanRawXml(o.Policy.Text)
	}
	if o.RedistRules != nil {
		raw["rr"] = util.CleanRawXml(o.RedistRules.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type container_v2 struct {
	Answer []entry_v2 `xml:"bgp"`
}

func (o *container_v2) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v2) Names() []string {
	return nil
}

func (o *entry_v2) normalize() Config {
	ans := Config{
		Enable:                        util.AsBool(o.Enable),
		RouterId:                      o.RouterId,
		AsNumber:                      o.AsNumber,
		RejectDefaultRoute:            util.AsBool(o.RejectDefaultRoute),
		InstallRoute:                  util.AsBool(o.InstallRoute),
		EcmpMultiAs:                   util.AsBool(o.EcmpMultiAs),
		AllowRedistributeDefaultRoute: util.AsBool(o.AllowRedistributeDefaultRoute),
	}

	raw := make(map[string]string)

	if o.Options != nil {
		ans.AsFormat = o.Options.AsFormat
		ans.DefaultLocalPreference = o.Options.DefaultLocalPreference
		ans.ReflectorClusterId = o.Options.ReflectorClusterId
		ans.ConfederationMemberAs = o.Options.ConfederationMemberAs

		if o.Options.Med != nil {
			ans.AlwaysCompareMed = util.AsBool(o.Options.Med.AlwaysCompareMed)
			ans.DeterministicMedComparison = util.AsBool(o.Options.Med.DeterministicMedComparison)
		}

		if o.Options.GracefulRestart != nil {
			ans.EnableGracefulRestart = util.AsBool(o.Options.GracefulRestart.EnableGracefulRestart)
			ans.StaleRouteTime = o.Options.GracefulRestart.StaleRouteTime
			ans.LocalRestartTime = o.Options.GracefulRestart.LocalRestartTime
			ans.MaxPeerRestartTime = o.Options.GracefulRestart.MaxPeerRestartTime
		}

		if o.Options.Aggregate != nil {
			ans.AggregateMed = util.AsBool(o.Options.Aggregate.AggregateMed)
		}

		if o.Options.OutboundRouteFilter != nil {
			raw["orf"] = util.CleanRawXml(o.Options.OutboundRouteFilter.Text)
		}
	}

	if o.AuthProfile != nil {
		raw["ap"] = util.CleanRawXml(o.AuthProfile.Text)
	}
	if o.DampeningProfile != nil {
		raw["dp"] = util.CleanRawXml(o.DampeningProfile.Text)
	}
	if o.PeerGroup != nil {
		raw["pg"] = util.CleanRawXml(o.PeerGroup.Text)
	}
	if o.Policy != nil {
		raw["poli"] = util.CleanRawXml(o.Policy.Text)
	}
	if o.RedistRules != nil {
		raw["rr"] = util.CleanRawXml(o.RedistRules.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type container_v3 struct {
	Answer []entry_v3 `xml:"bgp"`
}

func (o *container_v3) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v3) Names() []string {
	return nil
}

func (o *entry_v3) normalize() Config {
	ans := Config{
		Enable:                        util.AsBool(o.Enable),
		RouterId:                      o.RouterId,
		AsNumber:                      o.AsNumber,
		RejectDefaultRoute:            util.AsBool(o.RejectDefaultRoute),
		InstallRoute:                  util.AsBool(o.InstallRoute),
		EcmpMultiAs:                   util.AsBool(o.EcmpMultiAs),
		AllowRedistributeDefaultRoute: util.AsBool(o.AllowRedistributeDefaultRoute),
	}

	raw := make(map[string]string)

	if o.GlobalBfd != nil {
		ans.BfdProfile = o.GlobalBfd.BfdProfile
	}

	if o.Options != nil {
		ans.AsFormat = o.Options.AsFormat
		ans.DefaultLocalPreference = o.Options.DefaultLocalPreference
		ans.ReflectorClusterId = o.Options.ReflectorClusterId
		ans.ConfederationMemberAs = o.Options.ConfederationMemberAs

		if o.Options.Med != nil {
			ans.AlwaysCompareMed = util.AsBool(o.Options.Med.AlwaysCompareMed)
			ans.DeterministicMedComparison = util.AsBool(o.Options.Med.DeterministicMedComparison)
		}

		if o.Options.GracefulRestart != nil {
			ans.EnableGracefulRestart = util.AsBool(o.Options.GracefulRestart.EnableGracefulRestart)
			ans.StaleRouteTime = o.Options.GracefulRestart.StaleRouteTime
			ans.LocalRestartTime = o.Options.GracefulRestart.LocalRestartTime
			ans.MaxPeerRestartTime = o.Options.GracefulRestart.MaxPeerRestartTime
		}

		if o.Options.Aggregate != nil {
			ans.AggregateMed = util.AsBool(o.Options.Aggregate.AggregateMed)
		}

		if o.Options.OutboundRouteFilter != nil {
			raw["orf"] = util.CleanRawXml(o.Options.OutboundRouteFilter.Text)
		}
	}

	if o.AuthProfile != nil {
		raw["ap"] = util.CleanRawXml(o.AuthProfile.Text)
	}
	if o.DampeningProfile != nil {
		raw["dp"] = util.CleanRawXml(o.DampeningProfile.Text)
	}
	if o.PeerGroup != nil {
		raw["pg"] = util.CleanRawXml(o.PeerGroup.Text)
	}
	if o.Policy != nil {
		raw["poli"] = util.CleanRawXml(o.Policy.Text)
	}
	if o.RedistRules != nil {
		raw["rr"] = util.CleanRawXml(o.RedistRules.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type container_v4 struct {
	Answer []entry_v4 `xml:"bgp"`
}

func (o *container_v4) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v4) Names() []string {
	return nil
}

func (o *entry_v4) normalize() Config {
	ans := Config{
		Enable:                        util.AsBool(o.Enable),
		RouterId:                      o.RouterId,
		AsNumber:                      o.AsNumber,
		RejectDefaultRoute:            util.AsBool(o.RejectDefaultRoute),
		InstallRoute:                  util.AsBool(o.InstallRoute),
		EcmpMultiAs:                   util.AsBool(o.EcmpMultiAs),
		EnforceFirstAs:                util.AsBool(o.EnforceFirstAs),
		AllowRedistributeDefaultRoute: util.AsBool(o.AllowRedistributeDefaultRoute),
	}

	raw := make(map[string]string)

	if o.GlobalBfd != nil {
		ans.BfdProfile = o.GlobalBfd.BfdProfile
	}

	if o.Options != nil {
		ans.AsFormat = o.Options.AsFormat
		ans.DefaultLocalPreference = o.Options.DefaultLocalPreference
		ans.ReflectorClusterId = o.Options.ReflectorClusterId
		ans.ConfederationMemberAs = o.Options.ConfederationMemberAs

		if o.Options.Med != nil {
			ans.AlwaysCompareMed = util.AsBool(o.Options.Med.AlwaysCompareMed)
			ans.DeterministicMedComparison = util.AsBool(o.Options.Med.DeterministicMedComparison)
		}

		if o.Options.GracefulRestart != nil {
			ans.EnableGracefulRestart = util.AsBool(o.Options.GracefulRestart.EnableGracefulRestart)
			ans.StaleRouteTime = o.Options.GracefulRestart.StaleRouteTime
			ans.LocalRestartTime = o.Options.GracefulRestart.LocalRestartTime
			ans.MaxPeerRestartTime = o.Options.GracefulRestart.MaxPeerRestartTime
		}

		if o.Options.Aggregate != nil {
			ans.AggregateMed = util.AsBool(o.Options.Aggregate.AggregateMed)
		}

		if o.Options.OutboundRouteFilter != nil {
			raw["orf"] = util.CleanRawXml(o.Options.OutboundRouteFilter.Text)
		}
	}

	if o.AuthProfile != nil {
		raw["ap"] = util.CleanRawXml(o.AuthProfile.Text)
	}
	if o.DampeningProfile != nil {
		raw["dp"] = util.CleanRawXml(o.DampeningProfile.Text)
	}
	if o.PeerGroup != nil {
		raw["pg"] = util.CleanRawXml(o.PeerGroup.Text)
	}
	if o.Policy != nil {
		raw["poli"] = util.CleanRawXml(o.Policy.Text)
	}
	if o.RedistRules != nil {
		raw["rr"] = util.CleanRawXml(o.RedistRules.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName                       xml.Name  `xml:"bgp"`
	Enable                        string    `xml:"enable"`
	RouterId                      string    `xml:"router-id,omitempty"`
	AsNumber                      string    `xml:"local-as,omitempty"`
	RejectDefaultRoute            string    `xml:"reject-default-route"`
	InstallRoute                  string    `xml:"install-route"`
	AllowRedistributeDefaultRoute string    `xml:"allow-redist-default-route"`
	Options                       *rOptions `xml:"routing-options"`

	AuthProfile      *util.RawXml `xml:"auth-profile"`
	DampeningProfile *util.RawXml `xml:"dampening-profile"`
	PeerGroup        *util.RawXml `xml:"peer-group"`
	Policy           *util.RawXml `xml:"policy"`
	RedistRules      *util.RawXml `xml:"redist-rules"`
}

type rOptions struct {
	AsFormat               string           `xml:"as-format,omitempty"`
	Med                    *med             `xml:"med"`
	DefaultLocalPreference string           `xml:"default-local-preference,omitempty"`
	ReflectorClusterId     string           `xml:"reflector-cluster-id,omitempty"`
	ConfederationMemberAs  string           `xml:"confederation-member-as,omitempty"`
	GracefulRestart        *gracefulRestart `xml:"graceful-restart"`
	Aggregate              *aggOptions      `xml:"aggregate"`

	OutboundRouteFilter *util.RawXml `xml:"outbound-route-filter"`
}

type med struct {
	AlwaysCompareMed           string `xml:"always-compare-med"`
	DeterministicMedComparison string `xml:"deterministic-med-comparison"`
}

type gracefulRestart struct {
	EnableGracefulRestart string `xml:"enable"`
	StaleRouteTime        int    `xml:"stale-route-time,omitempty"`
	LocalRestartTime      int    `xml:"local-restart-time,omitempty"`
	MaxPeerRestartTime    int    `xml:"max-peer-restart-time,omitempty"`
}

type aggOptions struct {
	AggregateMed string `xml:"aggregate-med"`
}

func specify_v1(e Config) interface{} {
	ans := entry_v1{
		Enable:                        util.YesNo(e.Enable),
		RouterId:                      e.RouterId,
		AsNumber:                      e.AsNumber,
		RejectDefaultRoute:            util.YesNo(e.RejectDefaultRoute),
		InstallRoute:                  util.YesNo(e.InstallRoute),
		AllowRedistributeDefaultRoute: util.YesNo(e.AllowRedistributeDefaultRoute),
	}

	hasMed := e.AlwaysCompareMed || e.DeterministicMedComparison
	hasGracefulRestart := e.EnableGracefulRestart || e.StaleRouteTime != 0 || e.LocalRestartTime != 0 || e.MaxPeerRestartTime != 0
	hasAggOptions := e.AggregateMed

	if hasMed || hasGracefulRestart || hasAggOptions || e.AsFormat != "" || e.DefaultLocalPreference != "" || e.ReflectorClusterId != "" || e.ConfederationMemberAs != "" {
		o := rOptions{
			AsFormat:               e.AsFormat,
			DefaultLocalPreference: e.DefaultLocalPreference,
			ReflectorClusterId:     e.ReflectorClusterId,
			ConfederationMemberAs:  e.ConfederationMemberAs,
		}

		if hasMed {
			o.Med = &med{
				AlwaysCompareMed:           util.YesNo(e.AlwaysCompareMed),
				DeterministicMedComparison: util.YesNo(e.DeterministicMedComparison),
			}
		}

		if hasGracefulRestart {
			o.GracefulRestart = &gracefulRestart{
				EnableGracefulRestart: util.YesNo(e.EnableGracefulRestart),
				StaleRouteTime:        e.StaleRouteTime,
				LocalRestartTime:      e.LocalRestartTime,
				MaxPeerRestartTime:    e.MaxPeerRestartTime,
			}
		}

		if hasAggOptions {
			o.Aggregate = &aggOptions{
				AggregateMed: util.YesNo(e.AggregateMed),
			}
		}

		if text, present := e.raw["orf"]; present {
			o.OutboundRouteFilter = &util.RawXml{text}
		}

		ans.Options = &o
	}

	if text, present := e.raw["ap"]; present {
		ans.AuthProfile = &util.RawXml{text}
	}
	if text, present := e.raw["dp"]; present {
		ans.DampeningProfile = &util.RawXml{text}
	}
	if text, present := e.raw["pg"]; present {
		ans.PeerGroup = &util.RawXml{text}
	}
	if text, present := e.raw["poli"]; present {
		ans.Policy = &util.RawXml{text}
	}
	if text, present := e.raw["rr"]; present {
		ans.RedistRules = &util.RawXml{text}
	}

	return ans
}

type entry_v2 struct {
	XMLName                       xml.Name  `xml:"bgp"`
	Enable                        string    `xml:"enable"`
	RouterId                      string    `xml:"router-id,omitempty"`
	AsNumber                      string    `xml:"local-as,omitempty"`
	RejectDefaultRoute            string    `xml:"reject-default-route"`
	InstallRoute                  string    `xml:"install-route"`
	EcmpMultiAs                   string    `xml:"ecmp-multi-as"`
	AllowRedistributeDefaultRoute string    `xml:"allow-redist-default-route"`
	Options                       *rOptions `xml:"routing-options"`

	AuthProfile      *util.RawXml `xml:"auth-profile"`
	DampeningProfile *util.RawXml `xml:"dampening-profile"`
	PeerGroup        *util.RawXml `xml:"peer-group"`
	Policy           *util.RawXml `xml:"policy"`
	RedistRules      *util.RawXml `xml:"redist-rules"`
}

func specify_v2(e Config) interface{} {
	ans := entry_v2{
		Enable:                        util.YesNo(e.Enable),
		RouterId:                      e.RouterId,
		AsNumber:                      e.AsNumber,
		RejectDefaultRoute:            util.YesNo(e.RejectDefaultRoute),
		InstallRoute:                  util.YesNo(e.InstallRoute),
		EcmpMultiAs:                   util.YesNo(e.EcmpMultiAs),
		AllowRedistributeDefaultRoute: util.YesNo(e.AllowRedistributeDefaultRoute),
	}

	hasMed := e.AlwaysCompareMed || e.DeterministicMedComparison
	hasGracefulRestart := e.EnableGracefulRestart || e.StaleRouteTime != 0 || e.LocalRestartTime != 0 || e.MaxPeerRestartTime != 0
	hasAggOptions := e.AggregateMed

	if hasMed || hasGracefulRestart || hasAggOptions || e.AsFormat != "" || e.DefaultLocalPreference != "" || e.ReflectorClusterId != "" || e.ConfederationMemberAs != "" {
		o := rOptions{
			AsFormat:               e.AsFormat,
			DefaultLocalPreference: e.DefaultLocalPreference,
			ReflectorClusterId:     e.ReflectorClusterId,
			ConfederationMemberAs:  e.ConfederationMemberAs,
		}

		if hasMed {
			o.Med = &med{
				AlwaysCompareMed:           util.YesNo(e.AlwaysCompareMed),
				DeterministicMedComparison: util.YesNo(e.DeterministicMedComparison),
			}
		}

		if hasGracefulRestart {
			o.GracefulRestart = &gracefulRestart{
				EnableGracefulRestart: util.YesNo(e.EnableGracefulRestart),
				StaleRouteTime:        e.StaleRouteTime,
				LocalRestartTime:      e.LocalRestartTime,
				MaxPeerRestartTime:    e.MaxPeerRestartTime,
			}
		}

		if hasAggOptions {
			o.Aggregate = &aggOptions{
				AggregateMed: util.YesNo(e.AggregateMed),
			}
		}

		if text, present := e.raw["orf"]; present {
			o.OutboundRouteFilter = &util.RawXml{text}
		}

		ans.Options = &o
	}

	if text, present := e.raw["ap"]; present {
		ans.AuthProfile = &util.RawXml{text}
	}
	if text, present := e.raw["dp"]; present {
		ans.DampeningProfile = &util.RawXml{text}
	}
	if text, present := e.raw["pg"]; present {
		ans.PeerGroup = &util.RawXml{text}
	}
	if text, present := e.raw["poli"]; present {
		ans.Policy = &util.RawXml{text}
	}
	if text, present := e.raw["rr"]; present {
		ans.RedistRules = &util.RawXml{text}
	}

	return ans
}

type entry_v3 struct {
	XMLName                       xml.Name   `xml:"bgp"`
	Enable                        string     `xml:"enable"`
	RouterId                      string     `xml:"router-id,omitempty"`
	AsNumber                      string     `xml:"local-as,omitempty"`
	GlobalBfd                     *globalBfd `xml:"global-bfd"`
	RejectDefaultRoute            string     `xml:"reject-default-route"`
	InstallRoute                  string     `xml:"install-route"`
	EcmpMultiAs                   string     `xml:"ecmp-multi-as"`
	AllowRedistributeDefaultRoute string     `xml:"allow-redist-default-route"`
	Options                       *rOptions  `xml:"routing-options"`

	AuthProfile      *util.RawXml `xml:"auth-profile"`
	DampeningProfile *util.RawXml `xml:"dampening-profile"`
	PeerGroup        *util.RawXml `xml:"peer-group"`
	Policy           *util.RawXml `xml:"policy"`
	RedistRules      *util.RawXml `xml:"redist-rules"`
}

type globalBfd struct {
	BfdProfile string `xml:"profile,omitempty"`
}

func specify_v3(e Config) interface{} {
	ans := entry_v3{
		Enable:                        util.YesNo(e.Enable),
		RouterId:                      e.RouterId,
		AsNumber:                      e.AsNumber,
		RejectDefaultRoute:            util.YesNo(e.RejectDefaultRoute),
		InstallRoute:                  util.YesNo(e.InstallRoute),
		EcmpMultiAs:                   util.YesNo(e.EcmpMultiAs),
		AllowRedistributeDefaultRoute: util.YesNo(e.AllowRedistributeDefaultRoute),
	}

	if e.BfdProfile != "" {
		ans.GlobalBfd = &globalBfd{
			BfdProfile: e.BfdProfile,
		}
	}

	hasMed := e.AlwaysCompareMed || e.DeterministicMedComparison
	hasGracefulRestart := e.EnableGracefulRestart || e.StaleRouteTime != 0 || e.LocalRestartTime != 0 || e.MaxPeerRestartTime != 0
	hasAggOptions := e.AggregateMed

	if hasMed || hasGracefulRestart || hasAggOptions || e.AsFormat != "" || e.DefaultLocalPreference != "" || e.ReflectorClusterId != "" || e.ConfederationMemberAs != "" {
		o := rOptions{
			AsFormat:               e.AsFormat,
			DefaultLocalPreference: e.DefaultLocalPreference,
			ReflectorClusterId:     e.ReflectorClusterId,
			ConfederationMemberAs:  e.ConfederationMemberAs,
		}

		if hasMed {
			o.Med = &med{
				AlwaysCompareMed:           util.YesNo(e.AlwaysCompareMed),
				DeterministicMedComparison: util.YesNo(e.DeterministicMedComparison),
			}
		}

		if hasGracefulRestart {
			o.GracefulRestart = &gracefulRestart{
				EnableGracefulRestart: util.YesNo(e.EnableGracefulRestart),
				StaleRouteTime:        e.StaleRouteTime,
				LocalRestartTime:      e.LocalRestartTime,
				MaxPeerRestartTime:    e.MaxPeerRestartTime,
			}
		}

		if hasAggOptions {
			o.Aggregate = &aggOptions{
				AggregateMed: util.YesNo(e.AggregateMed),
			}
		}

		if text, present := e.raw["orf"]; present {
			o.OutboundRouteFilter = &util.RawXml{text}
		}

		ans.Options = &o
	}

	if text, present := e.raw["ap"]; present {
		ans.AuthProfile = &util.RawXml{text}
	}
	if text, present := e.raw["dp"]; present {
		ans.DampeningProfile = &util.RawXml{text}
	}
	if text, present := e.raw["pg"]; present {
		ans.PeerGroup = &util.RawXml{text}
	}
	if text, present := e.raw["poli"]; present {
		ans.Policy = &util.RawXml{text}
	}
	if text, present := e.raw["rr"]; present {
		ans.RedistRules = &util.RawXml{text}
	}

	return ans
}

type entry_v4 struct {
	XMLName                       xml.Name   `xml:"bgp"`
	Enable                        string     `xml:"enable"`
	RouterId                      string     `xml:"router-id,omitempty"`
	AsNumber                      string     `xml:"local-as,omitempty"`
	GlobalBfd                     *globalBfd `xml:"global-bfd"`
	RejectDefaultRoute            string     `xml:"reject-default-route"`
	InstallRoute                  string     `xml:"install-route"`
	EcmpMultiAs                   string     `xml:"ecmp-multi-as"`
	EnforceFirstAs                string     `xml:"enforce-first-as"`
	AllowRedistributeDefaultRoute string     `xml:"allow-redist-default-route"`
	Options                       *rOptions  `xml:"routing-options"`

	AuthProfile      *util.RawXml `xml:"auth-profile"`
	DampeningProfile *util.RawXml `xml:"dampening-profile"`
	PeerGroup        *util.RawXml `xml:"peer-group"`
	Policy           *util.RawXml `xml:"policy"`
	RedistRules      *util.RawXml `xml:"redist-rules"`
}

func specify_v4(e Config) interface{} {
	ans := entry_v4{
		Enable:                        util.YesNo(e.Enable),
		RouterId:                      e.RouterId,
		AsNumber:                      e.AsNumber,
		RejectDefaultRoute:            util.YesNo(e.RejectDefaultRoute),
		InstallRoute:                  util.YesNo(e.InstallRoute),
		EcmpMultiAs:                   util.YesNo(e.EcmpMultiAs),
		EnforceFirstAs:                util.YesNo(e.EnforceFirstAs),
		AllowRedistributeDefaultRoute: util.YesNo(e.AllowRedistributeDefaultRoute),
	}

	if e.BfdProfile != "" {
		ans.GlobalBfd = &globalBfd{
			BfdProfile: e.BfdProfile,
		}
	}

	hasMed := e.AlwaysCompareMed || e.DeterministicMedComparison
	hasGracefulRestart := e.EnableGracefulRestart || e.StaleRouteTime != 0 || e.LocalRestartTime != 0 || e.MaxPeerRestartTime != 0
	hasAggOptions := e.AggregateMed

	if hasMed || hasGracefulRestart || hasAggOptions || e.AsFormat != "" || e.DefaultLocalPreference != "" || e.ReflectorClusterId != "" || e.ConfederationMemberAs != "" {
		o := rOptions{
			AsFormat:               e.AsFormat,
			DefaultLocalPreference: e.DefaultLocalPreference,
			ReflectorClusterId:     e.ReflectorClusterId,
			ConfederationMemberAs:  e.ConfederationMemberAs,
		}

		if hasMed {
			o.Med = &med{
				AlwaysCompareMed:           util.YesNo(e.AlwaysCompareMed),
				DeterministicMedComparison: util.YesNo(e.DeterministicMedComparison),
			}
		}

		if hasGracefulRestart {
			o.GracefulRestart = &gracefulRestart{
				EnableGracefulRestart: util.YesNo(e.EnableGracefulRestart),
				StaleRouteTime:        e.StaleRouteTime,
				LocalRestartTime:      e.LocalRestartTime,
				MaxPeerRestartTime:    e.MaxPeerRestartTime,
			}
		}

		if hasAggOptions {
			o.Aggregate = &aggOptions{
				AggregateMed: util.YesNo(e.AggregateMed),
			}
		}

		if text, present := e.raw["orf"]; present {
			o.OutboundRouteFilter = &util.RawXml{text}
		}

		ans.Options = &o
	}

	if text, present := e.raw["ap"]; present {
		ans.AuthProfile = &util.RawXml{text}
	}
	if text, present := e.raw["dp"]; present {
		ans.DampeningProfile = &util.RawXml{text}
	}
	if text, present := e.raw["pg"]; present {
		ans.PeerGroup = &util.RawXml{text}
	}
	if text, present := e.raw["poli"]; present {
		ans.Policy = &util.RawXml{text}
	}
	if text, present := e.raw["rr"]; present {
		ans.RedistRules = &util.RawXml{text}
	}

	return ans
}
