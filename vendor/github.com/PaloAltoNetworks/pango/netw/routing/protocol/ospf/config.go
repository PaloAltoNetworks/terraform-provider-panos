package ospf

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Config is a normalized, version independent representation of a virtual
// router's OSPF configuration.
type Config struct {
	Enable                        bool
	RouterId                      string
	RejectDefaultRoute            bool
	AllowRedistributeDefaultRoute bool
	Rfc1583                       bool
	SpfCalculationDelay           float64
	LsaInterval                   float64
	EnableGracefulRestart         bool
	GracePeriod                   int
	HelperEnable                  bool
	StrictLsaChecking             bool
	MaxNeighborRestartTime        int
	BfdProfile                    string // BFD profile or "None" to disable BFD

	raw map[string]string
}

// Copy copies the information from source Config `s` to this object.
func (o *Config) Copy(s Config) {
	o.Enable = s.Enable
	o.RouterId = s.RouterId
	o.RejectDefaultRoute = s.RejectDefaultRoute
	o.AllowRedistributeDefaultRoute = s.AllowRedistributeDefaultRoute
	o.Rfc1583 = s.Rfc1583
	o.SpfCalculationDelay = s.SpfCalculationDelay
	o.LsaInterval = s.LsaInterval
	o.EnableGracefulRestart = s.EnableGracefulRestart
	o.GracePeriod = s.GracePeriod
	o.HelperEnable = s.HelperEnable
	o.StrictLsaChecking = s.StrictLsaChecking
	o.MaxNeighborRestartTime = s.MaxNeighborRestartTime
	o.BfdProfile = s.BfdProfile
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
	Answer []entry_v1 `xml:"ospf"`
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
		RejectDefaultRoute:            util.AsBool(o.RejectDefaultRoute),
		AllowRedistributeDefaultRoute: util.AsBool(o.AllowRedistributeDefaultRoute),
		Rfc1583:                       util.AsBool(o.Rfc1583),
	}

	raw := make(map[string]string)

	if o.Timers != nil {
		ans.SpfCalculationDelay = o.Timers.SpfCalculationDelay
		ans.LsaInterval = o.Timers.LsaInterval
	}
	if o.GracefulRestart != nil {
		ans.EnableGracefulRestart = util.AsBool(o.GracefulRestart.EnableGracefulRestart)
		ans.GracePeriod = o.GracefulRestart.GracePeriod
		ans.HelperEnable = util.AsBool(o.GracefulRestart.HelperEnable)
		ans.StrictLsaChecking = util.AsBool(o.GracefulRestart.StrictLsaChecking)
		ans.MaxNeighborRestartTime = o.GracefulRestart.MaxNeighborRestartTime
	}
	if o.GlobalBfd != nil {
		ans.BfdProfile = o.GlobalBfd.BfdProfile
	}

	if o.AuthProfile != nil {
		raw["ap"] = util.CleanRawXml(o.AuthProfile.Text)
	}
	if o.Area != nil {
		raw["area"] = util.CleanRawXml(o.Area.Text)
	}
	if o.ExportRules != nil {
		raw["exp"] = util.CleanRawXml(o.ExportRules.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName                       xml.Name         `xml:"ospf"`
	Enable                        string           `xml:"enable"`
	RouterId                      string           `xml:"router-id,omitempty"`
	RejectDefaultRoute            string           `xml:"reject-default-route"`
	AllowRedistributeDefaultRoute string           `xml:"allow-redist-default-route"`
	Rfc1583                       string           `xml:"rfc1583"`
	Timers                        *timers          `xml:"timers"`
	GracefulRestart               *gracefulRestart `xml:"graceful-restart"`
	GlobalBfd                     *globalBfd       `xml:"global-bfd"`

	AuthProfile *util.RawXml `xml:"auth-profile"`
	Area        *util.RawXml `xml:"area"`
	ExportRules *util.RawXml `xml:"export-rules"`
}

type timers struct {
	SpfCalculationDelay float64 `xml:"spf-calculation-delay,omitempty"`
	LsaInterval         float64 `xml:"lsa-interval,omitempty"`
}

type gracefulRestart struct {
	EnableGracefulRestart  string `xml:"enable"`
	GracePeriod            int    `xml:"grace-period,omitempty"`
	HelperEnable           string `xml:"helper-enable"`
	StrictLsaChecking      string `xml:"strict-LSA-checking"`
	MaxNeighborRestartTime int    `xml:"max-neighbor-restart-time,omitempty""`
}

type globalBfd struct {
	BfdProfile string `xml:"profile,omitempty"`
}

func specify_v1(e Config) interface{} {
	ans := entry_v1{
		Enable:                        util.YesNo(e.Enable),
		RouterId:                      e.RouterId,
		RejectDefaultRoute:            util.YesNo(e.RejectDefaultRoute),
		AllowRedistributeDefaultRoute: util.YesNo(e.AllowRedistributeDefaultRoute),
		Rfc1583:                       util.YesNo(e.Rfc1583),
	}

	if e.SpfCalculationDelay != 0 || e.LsaInterval != 0 {
		ans.Timers = &timers{
			SpfCalculationDelay: e.SpfCalculationDelay,
			LsaInterval:         e.LsaInterval,
		}
	}

	// EnableGracefulRestart, HelperEnable, StrictLsaChecking schema elements
	// are default="yes"
	if !e.EnableGracefulRestart || !e.HelperEnable || !e.StrictLsaChecking ||
		e.GracePeriod != 0 || e.MaxNeighborRestartTime != 0 {
		ans.GracefulRestart = &gracefulRestart{
			EnableGracefulRestart:  util.YesNo(e.EnableGracefulRestart),
			GracePeriod:            e.GracePeriod,
			HelperEnable:           util.YesNo(e.HelperEnable),
			StrictLsaChecking:      util.YesNo(e.StrictLsaChecking),
			MaxNeighborRestartTime: e.MaxNeighborRestartTime,
		}
	}

	if e.BfdProfile != "" {
		ans.GlobalBfd = &globalBfd{BfdProfile: e.BfdProfile}
	}

	if text, present := e.raw["ap"]; present {
		ans.AuthProfile = &util.RawXml{text}
	}
	if text, present := e.raw["area"]; present {
		ans.Area = &util.RawXml{text}
	}
	if text, present := e.raw["exp"]; present {
		ans.ExportRules = &util.RawXml{text}
	}

	return ans
}
