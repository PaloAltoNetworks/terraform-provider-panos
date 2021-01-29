package neighbor

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Config is a normalized, version independent representation of an IPv6
// neighbor discovery configuration.
//
// Due to the fact that RaLifetime is a rangedint field where not only is 0
// a valid value, but the default is 1800, this field will always be present
// in the marshalled XML sent to PAN-OS if either of the following is true:
//
// * another Ra* field is non-zero
// * this field is a value other than 1800
//
// PAN-OS 8.0+
type Config struct {
	EnableRa                          bool
	RaMaxInterval                     int
	RaMinInterval                     int
	RaManagedFlag                     bool
	RaOtherFlag                       bool
	RaLinkMtu                         string
	RaReachableTime                   string
	RaRetransmissionTimer             string
	RaHopLimit                        string
	RaLifetime                        int
	RaRouterPreference                string
	RaEnableConsistencyCheck          bool
	RaEnableDnsSupport                bool
	RaDnsServers                      []RaDnsServer
	RaDnsSuffixes                     []RaDnsSuffix
	EnableNdpMonitor                  bool
	EnableDuplicateAddressDetection   bool
	DuplicateAddressDetectionAttempts int
	NeighborSolicitationInterval      int
	ReachableTime                     int
	Neighbors                         []Neighbor
}

type RaDnsServer struct {
	Name     string
	Lifetime int
}

type RaDnsSuffix struct {
	Name     string
	Lifetime int
}

type Neighbor struct {
	Name       string
	MacAddress string
}

// Copy copies the information from source Config `s` to this object.
func (o *Config) Copy(s Config) {
	o.EnableRa = s.EnableRa
	o.RaMaxInterval = s.RaMaxInterval
	o.RaMinInterval = s.RaMinInterval
	o.RaManagedFlag = s.RaManagedFlag
	o.RaOtherFlag = s.RaOtherFlag
	o.RaLinkMtu = s.RaLinkMtu
	o.RaReachableTime = s.RaReachableTime
	o.RaRetransmissionTimer = s.RaRetransmissionTimer
	o.RaHopLimit = s.RaHopLimit
	o.RaLifetime = s.RaLifetime
	o.RaRouterPreference = s.RaRouterPreference
	o.RaEnableConsistencyCheck = s.RaEnableConsistencyCheck
	o.RaEnableDnsSupport = s.RaEnableDnsSupport
	o.RaDnsServers = s.RaDnsServers
	o.RaDnsSuffixes = s.RaDnsSuffixes
	o.EnableNdpMonitor = s.EnableNdpMonitor
	o.EnableDuplicateAddressDetection = s.EnableDuplicateAddressDetection
	o.DuplicateAddressDetectionAttempts = s.DuplicateAddressDetectionAttempts
	o.NeighborSolicitationInterval = s.NeighborSolicitationInterval
	o.ReachableTime = s.ReachableTime
	o.Neighbors = s.Neighbors
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
	Answer []entry_v1 `xml:"neighbor-discovery"`
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
		EnableNdpMonitor:                  util.AsBool(o.EnableNdpMonitor),
		EnableDuplicateAddressDetection:   util.AsBool(o.EnableDuplicateAddressDetection),
		DuplicateAddressDetectionAttempts: o.DuplicateAddressDetectionAttempts,
		NeighborSolicitationInterval:      o.NeighborSolicitationInterval,
		ReachableTime:                     o.ReachableTime,
	}

	if o.Ra != nil {
		ans.EnableRa = util.AsBool(o.Ra.EnableRa)
		ans.RaMaxInterval = o.Ra.RaMaxInterval
		ans.RaMinInterval = o.Ra.RaMinInterval
		ans.RaManagedFlag = util.AsBool(o.Ra.RaManagedFlag)
		ans.RaOtherFlag = util.AsBool(o.Ra.RaOtherFlag)
		ans.RaLinkMtu = o.Ra.RaLinkMtu
		ans.RaReachableTime = o.Ra.RaReachableTime
		ans.RaRetransmissionTimer = o.Ra.RaRetransmissionTimer
		ans.RaHopLimit = o.Ra.RaHopLimit
		ans.RaLifetime = o.Ra.RaLifetime
		ans.RaRouterPreference = o.Ra.RaRouterPreference
		ans.RaEnableConsistencyCheck = util.AsBool(o.Ra.RaEnableConsistencyCheck)

		if o.Ra.Dns != nil {
			ans.RaEnableDnsSupport = util.AsBool(o.Ra.Dns.RaEnableDnsSupport)

			if o.Ra.Dns.Servers != nil {
				list := make([]RaDnsServer, 0, len(o.Ra.Dns.Servers.Entries))
				for _, v := range o.Ra.Dns.Servers.Entries {
					list = append(list, RaDnsServer{
						Name:     v.Name,
						Lifetime: v.Lifetime,
					})
				}

				ans.RaDnsServers = list
			}

			if o.Ra.Dns.Suffixes != nil {
				list := make([]RaDnsSuffix, 0, len(o.Ra.Dns.Suffixes.Entries))
				for _, v := range o.Ra.Dns.Suffixes.Entries {
					list = append(list, RaDnsSuffix{
						Name:     v.Name,
						Lifetime: v.Lifetime,
					})
				}

				ans.RaDnsSuffixes = list
			}
		}
	}

	if o.Neighbors != nil {
		list := make([]Neighbor, 0, len(o.Neighbors.Entries))
		for _, v := range o.Neighbors.Entries {
			list = append(list, Neighbor{
				Name:       v.Name,
				MacAddress: v.MacAddress,
			})
		}

		ans.Neighbors = list
	}

	return ans
}

type entry_v1 struct {
	XMLName                           xml.Name   `xml:"neighbor-discovery"`
	Ra                                *ra        `xml:"router-advertisement"`
	EnableNdpMonitor                  string     `xml:"enable-ndp-monitor"`
	EnableDuplicateAddressDetection   string     `xml:"enable-dad"`
	DuplicateAddressDetectionAttempts int        `xml:"dad-attempts,omitempty"`
	NeighborSolicitationInterval      int        `xml:"ns-interval,omitempty"`
	ReachableTime                     int        `xml:"reachable-time,omitempty"`
	Neighbors                         *neighbors `xml:"neighbor"`
}

type ra struct {
	EnableRa                 string `xml:"enable"`
	RaMaxInterval            int    `xml:"max-interval,omitempty"`
	RaMinInterval            int    `xml:"min-interval,omitempty"`
	RaManagedFlag            string `xml:"managed-flag"`
	RaOtherFlag              string `xml:"other-flag"`
	RaLinkMtu                string `xml:"link-mtu,omitempty"`
	RaReachableTime          string `xml:"reachable-time,omitempty"`
	RaRetransmissionTimer    string `xml:"retransmission-timer,omitempty"`
	RaHopLimit               string `xml:"hop-limit,omitempty"`
	RaLifetime               int    `xml:"lifetime"`
	RaRouterPreference       string `xml:"router-preference,omitempty"`
	RaEnableConsistencyCheck string `xml:"enable-consistency-check"`
	Dns                      *dns   `xml:"dns"`
}

type dns struct {
	RaEnableDnsSupport string       `xml:"enable"`
	Servers            *dnsServers  `xml:"server"`
	Suffixes           *dnsSuffixes `xml:"suffix"`
}

type dnsServers struct {
	Entries []dnsServer `xml:"entry"`
}

type dnsServer struct {
	Name     string `xml:"name,attr"`
	Lifetime int    `xml:"lifetime,omitempty"`
}

type dnsSuffixes struct {
	Entries []dnsSuffix `xml:"entry"`
}

type dnsSuffix struct {
	Name     string `xml:"name,attr"`
	Lifetime int    `xml:"lifetime,omitempty"`
}

type neighbors struct {
	Entries []neighbor `xml:"entry"`
}

type neighbor struct {
	Name       string `xml:"name,attr"`
	MacAddress string `xml:"hw-address"`
}

func specify_v1(e Config) interface{} {
	ans := entry_v1{
		EnableNdpMonitor:                  util.YesNo(e.EnableNdpMonitor),
		EnableDuplicateAddressDetection:   util.YesNo(e.EnableDuplicateAddressDetection),
		DuplicateAddressDetectionAttempts: e.DuplicateAddressDetectionAttempts,
		NeighborSolicitationInterval:      e.NeighborSolicitationInterval,
		ReachableTime:                     e.ReachableTime,
	}

	incDns := e.RaEnableDnsSupport || len(e.RaDnsServers) > 0 || len(e.RaDnsSuffixes) > 0
	if e.EnableRa || e.RaMaxInterval > 0 || e.RaMinInterval > 0 || e.RaManagedFlag ||
		e.RaOtherFlag || e.RaLinkMtu != "" || e.RaReachableTime != "" ||
		e.RaRetransmissionTimer != "" || e.RaHopLimit != "" || e.RaLifetime != 1800 ||
		e.RaRouterPreference != "" || e.RaEnableConsistencyCheck ||
		incDns {
		ans.Ra = &ra{
			EnableRa:                 util.YesNo(e.EnableRa),
			RaMaxInterval:            e.RaMaxInterval,
			RaMinInterval:            e.RaMinInterval,
			RaManagedFlag:            util.YesNo(e.RaManagedFlag),
			RaOtherFlag:              util.YesNo(e.RaOtherFlag),
			RaLinkMtu:                e.RaLinkMtu,
			RaReachableTime:          e.RaReachableTime,
			RaRetransmissionTimer:    e.RaRetransmissionTimer,
			RaHopLimit:               e.RaHopLimit,
			RaLifetime:               e.RaLifetime,
			RaRouterPreference:       e.RaRouterPreference,
			RaEnableConsistencyCheck: util.YesNo(e.RaEnableConsistencyCheck),
		}

		if incDns {
			ans.Ra.Dns = &dns{
				RaEnableDnsSupport: util.YesNo(e.RaEnableDnsSupport),
			}

			if len(e.RaDnsServers) > 0 {
				list := make([]dnsServer, 0, len(e.RaDnsServers))
				for _, v := range e.RaDnsServers {
					list = append(list, dnsServer{
						Name:     v.Name,
						Lifetime: v.Lifetime,
					})
				}

				ans.Ra.Dns.Servers = &dnsServers{Entries: list}
			}

			if len(e.RaDnsSuffixes) > 0 {
				list := make([]dnsSuffix, 0, len(e.RaDnsSuffixes))
				for _, v := range e.RaDnsSuffixes {
					list = append(list, dnsSuffix{
						Name:     v.Name,
						Lifetime: v.Lifetime,
					})
				}

				ans.Ra.Dns.Suffixes = &dnsSuffixes{Entries: list}
			}
		}
	}

	if len(e.Neighbors) > 0 {
		list := make([]neighbor, 0, len(e.Neighbors))
		for _, v := range e.Neighbors {
			list = append(list, neighbor{
				Name:       v.Name,
				MacAddress: v.MacAddress,
			})
		}

		ans.Neighbors = &neighbors{Entries: list}
	}

	return ans
}
