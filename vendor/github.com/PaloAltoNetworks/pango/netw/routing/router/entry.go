package router

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a virtual
// router.
type Entry struct {
	Name                             string
	Interfaces                       []string
	StaticDist                       int
	StaticIpv6Dist                   int
	OspfIntDist                      int
	OspfExtDist                      int
	Ospfv3IntDist                    int
	Ospfv3ExtDist                    int
	IbgpDist                         int
	EbgpDist                         int
	RipDist                          int
	EnableEcmp                       bool
	EcmpSymmetricReturn              bool
	EcmpStrictSourcePath             bool
	EcmpMaxPath                      int
	EcmpLoadBalanceMethod            string
	EcmpHashSourceOnly               bool
	EcmpHashUsePort                  bool
	EcmpHashSeed                     int
	EcmpWeightedRoundRobinInterfaces map[string]int

	raw map[string]string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * StaticDist: 10
//      * StaticIpv6Dist: 10
//      * OspfIntDist: 30
//      * OspfExtDist: 110
//      * Ospfv3IntDist: 30
//      * Ospfv3ExtDist: 110
//      * IbgpDist: 200
//      * EbgpDist: 20
//      * RipDist: 120
func (o *Entry) Defaults() {
	if o.StaticDist == 0 {
		o.StaticDist = 10
	}

	if o.StaticIpv6Dist == 0 {
		o.StaticIpv6Dist = 10
	}

	if o.OspfIntDist == 0 {
		o.OspfIntDist = 30
	}

	if o.OspfExtDist == 0 {
		o.OspfExtDist = 110
	}

	if o.Ospfv3IntDist == 0 {
		o.Ospfv3IntDist = 30
	}

	if o.Ospfv3ExtDist == 0 {
		o.Ospfv3ExtDist = 110
	}

	if o.IbgpDist == 0 {
		o.IbgpDist = 200
	}

	if o.EbgpDist == 0 {
		o.EbgpDist = 20
	}

	if o.RipDist == 0 {
		o.RipDist = 120
	}
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Interfaces = make([]string, len(s.Interfaces))
	copy(o.Interfaces, s.Interfaces)
	o.StaticDist = s.StaticDist
	o.StaticIpv6Dist = s.StaticIpv6Dist
	o.OspfIntDist = s.OspfIntDist
	o.OspfExtDist = s.OspfExtDist
	o.Ospfv3IntDist = s.Ospfv3IntDist
	o.Ospfv3ExtDist = s.Ospfv3ExtDist
	o.IbgpDist = s.IbgpDist
	o.EbgpDist = s.EbgpDist
	o.RipDist = s.RipDist
	o.EnableEcmp = s.EnableEcmp
	o.EcmpSymmetricReturn = s.EcmpSymmetricReturn
	o.EcmpStrictSourcePath = s.EcmpStrictSourcePath
	o.EcmpMaxPath = s.EcmpMaxPath
	o.EcmpLoadBalanceMethod = s.EcmpLoadBalanceMethod
	o.EcmpHashSourceOnly = s.EcmpHashSourceOnly
	o.EcmpHashUsePort = s.EcmpHashUsePort
	o.EcmpHashSeed = s.EcmpHashSeed
	o.EcmpWeightedRoundRobinInterfaces = s.EcmpWeightedRoundRobinInterfaces
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, string, interface{}) {
	_, fn := versioning(v)

	return o.Name, o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:       o.Name,
		Interfaces: util.MemToStr(o.Interfaces),
	}

	if o.Dist != nil {
		ans.StaticDist = o.Dist.StaticDist
		ans.StaticIpv6Dist = o.Dist.StaticIpv6Dist
		ans.OspfIntDist = o.Dist.OspfIntDist
		ans.OspfExtDist = o.Dist.OspfExtDist
		ans.Ospfv3IntDist = o.Dist.Ospfv3IntDist
		ans.Ospfv3ExtDist = o.Dist.Ospfv3ExtDist
		ans.IbgpDist = o.Dist.IbgpDist
		ans.EbgpDist = o.Dist.EbgpDist
		ans.RipDist = o.Dist.RipDist
	}

	if o.Ecmp != nil {
		ans.EnableEcmp = util.AsBool(o.Ecmp.EnableEcmp)
		ans.EcmpSymmetricReturn = util.AsBool(o.Ecmp.EcmpSymmetricReturn)
		ans.EcmpStrictSourcePath = util.AsBool(o.Ecmp.EcmpStrictSourcePath)
		ans.EcmpMaxPath = o.Ecmp.EcmpMaxPath

		if o.Ecmp.Algorithm != nil {
			if o.Ecmp.Algorithm.IpModulo != nil {
				ans.EcmpLoadBalanceMethod = EcmpLoadBalanceMethodIpModulo
			} else if o.Ecmp.Algorithm.IpHash != nil {
				ans.EcmpLoadBalanceMethod = EcmpLoadBalanceMethodIpHash
				ans.EcmpHashSourceOnly = util.AsBool(o.Ecmp.Algorithm.IpHash.EcmpHashSourceOnly)
				ans.EcmpHashUsePort = util.AsBool(o.Ecmp.Algorithm.IpHash.EcmpHashUsePort)
				ans.EcmpHashSeed = o.Ecmp.Algorithm.IpHash.EcmpHashSeed
			} else if o.Ecmp.Algorithm.Wrr != nil {
				ans.EcmpLoadBalanceMethod = EcmpLoadBalanceMethodWeightedRoundRobin
				if o.Ecmp.Algorithm.Wrr.Interfaces != nil {
					ans.EcmpWeightedRoundRobinInterfaces = make(map[string]int)
					for _, v := range o.Ecmp.Algorithm.Wrr.Interfaces.Entries {
						ans.EcmpWeightedRoundRobinInterfaces[v.Interface] = v.Weight
					}
				}
			} else if o.Ecmp.Algorithm.Brr != nil {
				ans.EcmpLoadBalanceMethod = EcmpLoadBalanceMethodBalancedRoundRobin
			}
		}
	}

	ans.raw = make(map[string]string)
	if o.Multicast != nil {
		ans.raw["multicast"] = util.CleanRawXml(o.Multicast.Text)
	}
	if o.Protocol != nil {
		ans.raw["protocol"] = util.CleanRawXml(o.Protocol.Text)
	}
	if o.Routing != nil {
		ans.raw["routing"] = util.CleanRawXml(o.Routing.Text)
	}

	if len(ans.raw) == 0 {
		ans.raw = nil
	}

	return ans
}

type entry_v1 struct {
	XMLName    xml.Name         `xml:"entry"`
	Name       string           `xml:"name,attr"`
	Interfaces *util.MemberType `xml:"interface"`
	Dist       *dist            `xml:"admin-dists"`
	Ecmp       *ecmp            `xml:"ecmp"`
	Multicast  *util.RawXml     `xml:"multicast"`
	Protocol   *util.RawXml     `xml:"protocol"`
	Routing    *util.RawXml     `xml:"routing-table"`
}

type dist struct {
	StaticDist     int `xml:"static,omitempty"`
	StaticIpv6Dist int `xml:"static-ipv6,omitempty"`
	OspfIntDist    int `xml:"ospf-int,omitempty"`
	OspfExtDist    int `xml:"ospf-ext,omitempty"`
	Ospfv3IntDist  int `xml:"ospfv3-int,omitempty"`
	Ospfv3ExtDist  int `xml:"ospfv3-ext,omitempty"`
	IbgpDist       int `xml:"ibgp,omitempty"`
	EbgpDist       int `xml:"ebgp,omitempty"`
	RipDist        int `xml:"rip,omitempty"`
}

type ecmp struct {
	EnableEcmp           string     `xml:"enable"`
	EcmpSymmetricReturn  string     `xml:"symmetric-return"`
	EcmpStrictSourcePath string     `xml:"strict-source-path"`
	EcmpMaxPath          int        `xml:"max-path,omitempty"`
	Algorithm            *algorithm `xml:"algorithm"`
}

type algorithm struct {
	IpModulo *string `xml:"ip-modulo"`
	IpHash   *ipHash `xml:"ip-hash"`
	Wrr      *wrr    `xml:"weighted-round-robin"`
	Brr      *string `xml:"balanced-round-robin"`
}

type ipHash struct {
	EcmpHashSourceOnly string `xml:"src-only"`
	EcmpHashUsePort    string `xml:"use-port"`
	EcmpHashSeed       int    `xml:"hash-seed,omitempty"`
}

type wrr struct {
	Interfaces *wrrInterfaces `xml:"interface"`
}

type wrrInterfaces struct {
	Entries []wrrInterface `xml:"entry"`
}

type wrrInterface struct {
	XMLName   xml.Name `xml:"entry"`
	Interface string   `xml:"name,attr"`
	Weight    int      `xml:"weight,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:       e.Name,
		Interfaces: util.StrToMem(e.Interfaces),
	}

	if e.StaticDist != 0 || e.StaticIpv6Dist != 0 || e.OspfIntDist != 0 || e.OspfExtDist != 0 || e.Ospfv3IntDist != 0 || e.Ospfv3ExtDist != 0 || e.IbgpDist != 0 || e.EbgpDist != 0 || e.RipDist != 0 {
		ans.Dist = &dist{
			StaticDist:     e.StaticDist,
			StaticIpv6Dist: e.StaticIpv6Dist,
			OspfIntDist:    e.OspfIntDist,
			OspfExtDist:    e.OspfExtDist,
			Ospfv3IntDist:  e.Ospfv3IntDist,
			Ospfv3ExtDist:  e.Ospfv3ExtDist,
			IbgpDist:       e.IbgpDist,
			EbgpDist:       e.EbgpDist,
			RipDist:        e.RipDist,
		}
	}

	if e.EnableEcmp || e.EcmpSymmetricReturn || e.EcmpStrictSourcePath || e.EcmpMaxPath != 0 || e.EcmpLoadBalanceMethod != "" {
		s := ""
		ans.Ecmp = &ecmp{
			EnableEcmp:           util.YesNo(e.EnableEcmp),
			EcmpSymmetricReturn:  util.YesNo(e.EcmpSymmetricReturn),
			EcmpStrictSourcePath: util.YesNo(e.EcmpStrictSourcePath),
			EcmpMaxPath:          e.EcmpMaxPath,
		}

		switch e.EcmpLoadBalanceMethod {
		case EcmpLoadBalanceMethodIpModulo:
			ans.Ecmp.Algorithm = &algorithm{
				IpModulo: &s,
			}
		case EcmpLoadBalanceMethodIpHash:
			ans.Ecmp.Algorithm = &algorithm{
				IpHash: &ipHash{
					EcmpHashSourceOnly: util.YesNo(e.EcmpHashSourceOnly),
					EcmpHashUsePort:    util.YesNo(e.EcmpHashUsePort),
					EcmpHashSeed:       e.EcmpHashSeed,
				},
			}
		case EcmpLoadBalanceMethodWeightedRoundRobin:
			ans.Ecmp.Algorithm = &algorithm{
				Wrr: &wrr{},
			}
			if len(e.EcmpWeightedRoundRobinInterfaces) > 0 {
				list := make([]wrrInterface, 0, len(e.EcmpWeightedRoundRobinInterfaces))
				for name, weight := range e.EcmpWeightedRoundRobinInterfaces {
					list = append(list, wrrInterface{
						Interface: name,
						Weight:    weight,
					})
				}
				ans.Ecmp.Algorithm.Wrr.Interfaces = &wrrInterfaces{
					Entries: list,
				}
			}
		case EcmpLoadBalanceMethodBalancedRoundRobin:
			ans.Ecmp.Algorithm = &algorithm{
				Brr: &s,
			}
		}
	}

	if text, present := e.raw["multicast"]; present {
		ans.Multicast = &util.RawXml{text}
	}
	if text, present := e.raw["protocol"]; present {
		ans.Protocol = &util.RawXml{text}
	}
	if text, present := e.raw["routing"]; present {
		ans.Routing = &util.RawXml{text}
	}

	return ans
}
