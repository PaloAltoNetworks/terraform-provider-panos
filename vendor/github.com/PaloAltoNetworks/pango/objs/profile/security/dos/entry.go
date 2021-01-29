package dos

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// DOS protection security profile.
type Entry struct {
	Name                      string
	Description               string
	Type                      string
	EnableSessionsProtections bool
	MaxConcurrentSessions     int
	Syn                       *SynProtection
	Udp                       *Protection
	Icmp                      *Protection
	Icmpv6                    *Protection
	Other                     *Protection
}

type SynProtection struct {
	Enable        bool
	Action        string
	AlarmRate     int
	ActivateRate  int
	MaxRate       int
	BlockDuration int
}

type Protection struct {
	Enable        bool
	AlarmRate     int
	ActivateRate  int
	MaxRate       int
	BlockDuration int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.Type = s.Type
	o.EnableSessionsProtections = s.EnableSessionsProtections
	o.MaxConcurrentSessions = s.MaxConcurrentSessions
	if s.Syn == nil {
		o.Syn = nil
	} else {
		o.Syn = &SynProtection{
			Enable:        s.Syn.Enable,
			Action:        s.Syn.Action,
			AlarmRate:     s.Syn.AlarmRate,
			ActivateRate:  s.Syn.ActivateRate,
			MaxRate:       s.Syn.MaxRate,
			BlockDuration: s.Syn.BlockDuration,
		}
	}
	if s.Udp == nil {
		o.Udp = nil
	} else {
		o.Udp = &Protection{
			Enable:        s.Udp.Enable,
			AlarmRate:     s.Udp.AlarmRate,
			ActivateRate:  s.Udp.ActivateRate,
			MaxRate:       s.Udp.MaxRate,
			BlockDuration: s.Udp.BlockDuration,
		}
	}
	if s.Icmp == nil {
		o.Icmp = nil
	} else {
		o.Icmp = &Protection{
			Enable:        s.Icmp.Enable,
			AlarmRate:     s.Icmp.AlarmRate,
			ActivateRate:  s.Icmp.ActivateRate,
			MaxRate:       s.Icmp.MaxRate,
			BlockDuration: s.Icmp.BlockDuration,
		}
	}
	if s.Icmpv6 == nil {
		o.Icmpv6 = nil
	} else {
		o.Icmpv6 = &Protection{
			Enable:        s.Icmpv6.Enable,
			AlarmRate:     s.Icmpv6.AlarmRate,
			ActivateRate:  s.Icmpv6.ActivateRate,
			MaxRate:       s.Icmpv6.MaxRate,
			BlockDuration: s.Icmpv6.BlockDuration,
		}
	}
	if s.Other == nil {
		o.Other = nil
	} else {
		o.Other = &Protection{
			Enable:        s.Other.Enable,
			AlarmRate:     s.Other.AlarmRate,
			ActivateRate:  s.Other.ActivateRate,
			MaxRate:       s.Other.MaxRate,
			BlockDuration: s.Other.BlockDuration,
		}
	}
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

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v1) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:        o.Name,
		Description: o.Description,
		Type:        o.Type,
	}

	if o.Flood != nil {
		if o.Flood.Syn != nil {
			if o.Flood.Syn.Red != nil {
				ans.Syn = &SynProtection{
					Action:       SynActionRed,
					AlarmRate:    o.Flood.Syn.Red.AlarmRate,
					ActivateRate: o.Flood.Syn.Red.ActivateRate,
					MaxRate:      o.Flood.Syn.Red.MaxRate,
				}
				if o.Flood.Syn.Red.Block != nil {
					ans.Syn.BlockDuration = o.Flood.Syn.Red.Block.BlockDuration
				}
			} else if o.Flood.Syn.Cookies != nil {
				ans.Syn = &SynProtection{
					Action:       SynActionCookies,
					AlarmRate:    o.Flood.Syn.Cookies.AlarmRate,
					ActivateRate: o.Flood.Syn.Cookies.ActivateRate,
					MaxRate:      o.Flood.Syn.Cookies.MaxRate,
				}
				if o.Flood.Syn.Cookies.Block != nil {
					ans.Syn.BlockDuration = o.Flood.Syn.Cookies.Block.BlockDuration
				}
			} else {
				ans.Syn = &SynProtection{}
			}
			ans.Syn.Enable = util.AsBool(o.Flood.Syn.Enable)
		}

		if o.Flood.Udp != nil {
			if o.Flood.Udp.Red != nil {
				ans.Udp = &Protection{
					AlarmRate:    o.Flood.Udp.Red.AlarmRate,
					ActivateRate: o.Flood.Udp.Red.ActivateRate,
					MaxRate:      o.Flood.Udp.Red.MaxRate,
				}
				if o.Flood.Udp.Red.Block != nil {
					ans.Udp.BlockDuration = o.Flood.Udp.Red.Block.BlockDuration
				}
			} else {
				ans.Udp = &Protection{}
			}
			ans.Udp.Enable = util.AsBool(o.Flood.Udp.Enable)
		}

		if o.Flood.Icmp != nil {
			if o.Flood.Icmp.Red != nil {
				ans.Icmp = &Protection{
					AlarmRate:    o.Flood.Icmp.Red.AlarmRate,
					ActivateRate: o.Flood.Icmp.Red.ActivateRate,
					MaxRate:      o.Flood.Icmp.Red.MaxRate,
				}
				if o.Flood.Icmp.Red.Block != nil {
					ans.Icmp.BlockDuration = o.Flood.Icmp.Red.Block.BlockDuration
				}
			} else {
				ans.Icmp = &Protection{}
			}
			ans.Icmp.Enable = util.AsBool(o.Flood.Icmp.Enable)
		}

		if o.Flood.Icmpv6 != nil {
			if o.Flood.Icmpv6.Red != nil {
				ans.Icmpv6 = &Protection{
					AlarmRate:    o.Flood.Icmpv6.Red.AlarmRate,
					ActivateRate: o.Flood.Icmpv6.Red.ActivateRate,
					MaxRate:      o.Flood.Icmpv6.Red.MaxRate,
				}
				if o.Flood.Icmpv6.Red.Block != nil {
					ans.Icmpv6.BlockDuration = o.Flood.Icmpv6.Red.Block.BlockDuration
				}
			} else {
				ans.Icmpv6 = &Protection{}
			}
			ans.Icmpv6.Enable = util.AsBool(o.Flood.Icmpv6.Enable)
		}

		if o.Flood.Other != nil {
			if o.Flood.Other.Red != nil {
				ans.Other = &Protection{
					AlarmRate:    o.Flood.Other.Red.AlarmRate,
					ActivateRate: o.Flood.Other.Red.ActivateRate,
					MaxRate:      o.Flood.Other.Red.MaxRate,
				}
				if o.Flood.Other.Red.Block != nil {
					ans.Other.BlockDuration = o.Flood.Other.Red.Block.BlockDuration
				}
			} else {
				ans.Other = &Protection{}
			}
			ans.Other.Enable = util.AsBool(o.Flood.Other.Enable)
		}
	}

	if o.Resource != nil {
		if o.Resource.Sess != nil {
			ans.EnableSessionsProtections = util.AsBool(o.Resource.Sess.EnableSessionsProtections)
			ans.MaxConcurrentSessions = o.Resource.Sess.MaxConcurrentSessions
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name  `xml:"entry"`
	Name        string    `xml:"name,attr"`
	Description string    `xml:"description,omitempty"`
	Type        string    `xml:"type"`
	Flood       *flood    `xml:"flood"`
	Resource    *resource `xml:"resource"`
}

type flood struct {
	Syn    *syn    `xml:"tcp-syn"`
	Udp    *common `xml:"udp"`
	Icmp   *common `xml:"icmp"`
	Icmpv6 *common `xml:"icmpv6"`
	Other  *common `xml:"other-ip"`
}

type syn struct {
	Enable  string   `xml:"enable"`
	Red     *details `xml:"red"`
	Cookies *details `xml:"syn-cookies"`
}

type details struct {
	AlarmRate    int    `xml:"alarm-rate"`
	ActivateRate int    `xml:"activate-rate"`
	MaxRate      int    `xml:"maximal-rate"`
	Block        *block `xml:"block"`
}

type block struct {
	BlockDuration int `xml:"duration,omitempty"`
}

type common struct {
	Enable string   `xml:"enable"`
	Red    *details `xml:"red"`
}

type resource struct {
	Sess *sess `xml:"sessions"`
}

type sess struct {
	EnableSessionsProtections string `xml:"enabled"`
	MaxConcurrentSessions     int    `xml:"max-concurrent-limit,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		Type:        e.Type,
	}

	if e.Syn != nil || e.Udp != nil || e.Icmp != nil || e.Icmpv6 != nil || e.Other != nil {
		ans.Flood = &flood{}
		if e.Syn != nil {
			ans.Flood.Syn = &syn{
				Enable: util.YesNo(e.Syn.Enable),
			}

			switch e.Syn.Action {
			case SynActionRed:
				ans.Flood.Syn.Red = &details{
					AlarmRate:    e.Syn.AlarmRate,
					ActivateRate: e.Syn.ActivateRate,
					MaxRate:      e.Syn.MaxRate,
				}
				if e.Syn.BlockDuration != 0 {
					ans.Flood.Syn.Red.Block = &block{
						BlockDuration: e.Syn.BlockDuration,
					}
				}
			case SynActionCookies:
				ans.Flood.Syn.Cookies = &details{
					AlarmRate:    e.Syn.AlarmRate,
					ActivateRate: e.Syn.ActivateRate,
					MaxRate:      e.Syn.MaxRate,
				}
				if e.Syn.BlockDuration != 0 {
					ans.Flood.Syn.Cookies.Block = &block{
						BlockDuration: e.Syn.BlockDuration,
					}
				}
			}
		}

		if e.Udp != nil {
			ans.Flood.Udp = &common{
				Enable: util.YesNo(e.Udp.Enable),
			}

			if e.Udp.AlarmRate != 0 || e.Udp.ActivateRate != 0 || e.Udp.MaxRate != 0 || e.Udp.BlockDuration != 0 {
				ans.Flood.Udp.Red = &details{
					AlarmRate:    e.Udp.AlarmRate,
					ActivateRate: e.Udp.ActivateRate,
					MaxRate:      e.Udp.MaxRate,
				}
				if e.Udp.BlockDuration != 0 {
					ans.Flood.Udp.Red.Block = &block{
						BlockDuration: e.Udp.BlockDuration,
					}
				}
			}
		}

		if e.Icmp != nil {
			ans.Flood.Icmp = &common{
				Enable: util.YesNo(e.Icmp.Enable),
			}

			if e.Icmp.AlarmRate != 0 || e.Icmp.ActivateRate != 0 || e.Icmp.MaxRate != 0 || e.Icmp.BlockDuration != 0 {
				ans.Flood.Icmp.Red = &details{
					AlarmRate:    e.Icmp.AlarmRate,
					ActivateRate: e.Icmp.ActivateRate,
					MaxRate:      e.Icmp.MaxRate,
				}
				if e.Icmp.BlockDuration != 0 {
					ans.Flood.Icmp.Red.Block = &block{
						BlockDuration: e.Icmp.BlockDuration,
					}
				}
			}
		}

		if e.Icmpv6 != nil {
			ans.Flood.Icmpv6 = &common{
				Enable: util.YesNo(e.Icmpv6.Enable),
			}

			if e.Icmpv6.AlarmRate != 0 || e.Icmpv6.ActivateRate != 0 || e.Icmpv6.MaxRate != 0 || e.Icmpv6.BlockDuration != 0 {
				ans.Flood.Icmpv6.Red = &details{
					AlarmRate:    e.Icmpv6.AlarmRate,
					ActivateRate: e.Icmpv6.ActivateRate,
					MaxRate:      e.Icmpv6.MaxRate,
				}
				if e.Icmpv6.BlockDuration != 0 {
					ans.Flood.Icmpv6.Red.Block = &block{
						BlockDuration: e.Icmpv6.BlockDuration,
					}
				}
			}
		}

		if e.Other != nil {
			ans.Flood.Other = &common{
				Enable: util.YesNo(e.Other.Enable),
			}

			if e.Other.AlarmRate != 0 || e.Other.ActivateRate != 0 || e.Other.MaxRate != 0 || e.Other.BlockDuration != 0 {
				ans.Flood.Other.Red = &details{
					AlarmRate:    e.Other.AlarmRate,
					ActivateRate: e.Other.ActivateRate,
					MaxRate:      e.Other.MaxRate,
				}
				if e.Other.BlockDuration != 0 {
					ans.Flood.Other.Red.Block = &block{
						BlockDuration: e.Other.BlockDuration,
					}
				}
			}
		}
	}

	if e.EnableSessionsProtections || e.MaxConcurrentSessions != 0 {
		ans.Resource = &resource{
			Sess: &sess{
				EnableSessionsProtections: util.YesNo(e.EnableSessionsProtections),
				MaxConcurrentSessions:     e.MaxConcurrentSessions,
			},
		}
	}

	return ans
}
