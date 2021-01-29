package spyware

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// anti-spyware security profile.
//
// PAN-OS 8.0+
type Entry struct {
	Name                string
	Description         string
	PacketCapture       string // 8.x only.
	BotnetLists         []BotnetList
	DnsCategories       []DnsCategory // 10.0
	WhiteLists          []WhiteList   // 10.0
	SinkholeIpv4Address string
	SinkholeIpv6Address string
	ThreatExceptions    []string
	Rules               []Rule
	Exceptions          []Exception
}

type BotnetList struct {
	Name          string
	Action        string
	PacketCapture string // 9.0+
}

// DnsCategory is present in PAN-OS 10.0+.
type DnsCategory struct {
	Name          string
	Action        string
	LogLevel      string
	PacketCapture string
}

// WhiteList is present in PAN-OS 10.0+.
type WhiteList struct {
	Name        string
	Description string
}

type Rule struct {
	Name            string
	ThreatName      string
	Category        string
	Severities      []string
	PacketCapture   string
	Action          string
	BlockIpTrackBy  string
	BlockIpDuration int
}

type Exception struct {
	Name            string
	PacketCapture   string
	Action          string
	BlockIpTrackBy  string
	BlockIpDuration int
	ExemptIps       []string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.PacketCapture = s.PacketCapture
	if len(s.BotnetLists) == 0 {
		o.BotnetLists = nil
	} else {
		o.BotnetLists = make([]BotnetList, 0, len(s.BotnetLists))
		for _, x := range s.BotnetLists {
			o.BotnetLists = append(o.BotnetLists, BotnetList{
				Name:          x.Name,
				Action:        x.Action,
				PacketCapture: x.PacketCapture,
			})
		}
	}
	if len(s.DnsCategories) == 0 {
		o.DnsCategories = nil
	} else {
		o.DnsCategories = make([]DnsCategory, 0, len(s.DnsCategories))
		for _, x := range s.DnsCategories {
			o.DnsCategories = append(o.DnsCategories, DnsCategory{
				Name:          x.Name,
				Action:        x.Action,
				LogLevel:      x.LogLevel,
				PacketCapture: x.PacketCapture,
			})
		}
	}
	if len(s.WhiteLists) == 0 {
		o.WhiteLists = nil
	} else {
		o.WhiteLists = make([]WhiteList, 0, len(s.WhiteLists))
		for _, x := range s.WhiteLists {
			o.WhiteLists = append(o.WhiteLists, WhiteList{
				Name:        x.Name,
				Description: x.Description,
			})
		}
	}
	o.SinkholeIpv4Address = s.SinkholeIpv4Address
	o.SinkholeIpv6Address = s.SinkholeIpv6Address
	if len(s.ThreatExceptions) == 0 {
		o.ThreatExceptions = nil
	} else {
		o.ThreatExceptions = make([]string, len(s.ThreatExceptions))
		copy(o.ThreatExceptions, s.ThreatExceptions)
	}
	if len(s.Rules) == 0 {
		o.Rules = nil
	} else {
		o.Rules = make([]Rule, 0, len(s.Rules))
		for _, x := range s.Rules {
			var sevs []string
			if len(x.Severities) > 0 {
				sevs = make([]string, len(x.Severities))
				copy(sevs, x.Severities)
			}
			o.Rules = append(o.Rules, Rule{
				Name:            x.Name,
				ThreatName:      x.ThreatName,
				Category:        x.Category,
				Severities:      sevs,
				PacketCapture:   x.PacketCapture,
				Action:          x.Action,
				BlockIpTrackBy:  x.BlockIpTrackBy,
				BlockIpDuration: x.BlockIpDuration,
			})
		}
	}
	if len(s.Exceptions) == 0 {
		o.Exceptions = nil
	} else {
		o.Exceptions = make([]Exception, 0, len(s.Exceptions))
		for _, x := range s.Exceptions {
			var eis []string
			if len(x.ExemptIps) > 0 {
				eis = make([]string, len(x.ExemptIps))
				copy(eis, x.ExemptIps)
			}
			o.Exceptions = append(o.Exceptions, Exception{
				Name:            x.Name,
				PacketCapture:   x.PacketCapture,
				Action:          x.Action,
				BlockIpTrackBy:  x.BlockIpTrackBy,
				BlockIpDuration: x.BlockIpDuration,
				ExemptIps:       eis,
			})
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

// 8.0
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
	}

	if o.Botnet != nil {
		ans.PacketCapture = o.Botnet.PacketCapture
		ans.ThreatExceptions = util.EntToStr(o.Botnet.ThreatExceptions)

		if o.Botnet.Lists != nil {
			lists := make([]BotnetList, 0, len(o.Botnet.Lists.Entries))
			for _, x := range o.Botnet.Lists.Entries {
				val := BotnetList{
					Name: x.Name,
				}

				switch {
				case x.Action.Alert != nil:
					val.Action = ActionAlert
				case x.Action.Allow != nil:
					val.Action = ActionAllow
				case x.Action.Block != nil:
					val.Action = ActionBlock
				case x.Action.Sinkhole != nil:
					val.Action = ActionSinkhole
				}

				lists = append(lists, val)
			}

			ans.BotnetLists = lists
		}

		if o.Botnet.Sinkhole != nil {
			ans.SinkholeIpv4Address = o.Botnet.Sinkhole.SinkholeIpv4Address
			ans.SinkholeIpv6Address = o.Botnet.Sinkhole.SinkholeIpv6Address
		}
	}

	if o.Rules != nil {
		list := make([]Rule, 0, len(o.Rules.Entries))
		for _, x := range o.Rules.Entries {
			item := Rule{
				Name:          x.Name,
				ThreatName:    x.ThreatName,
				Category:      x.Category,
				Severities:    util.MemToStr(x.Severities),
				PacketCapture: x.PacketCapture,
			}

			if x.Action != nil {
				switch {
				case x.Action.Default != nil:
					item.Action = ActionDefault
				case x.Action.Allow != nil:
					item.Action = ActionAllow
				case x.Action.Alert != nil:
					item.Action = ActionAlert
				case x.Action.Drop != nil:
					item.Action = ActionDrop
				case x.Action.ResetClient != nil:
					item.Action = ActionResetClient
				case x.Action.ResetServer != nil:
					item.Action = ActionResetServer
				case x.Action.ResetBoth != nil:
					item.Action = ActionResetBoth
				case x.Action.BlockIp != nil:
					item.Action = ActionBlockIp
					item.BlockIpTrackBy = x.Action.BlockIp.TrackBy
					item.BlockIpDuration = x.Action.BlockIp.Duration
				}
			}

			list = append(list, item)
		}
		ans.Rules = list
	}

	if o.Exceptions != nil {
		list := make([]Exception, 0, len(o.Exceptions.Entries))
		for _, x := range o.Exceptions.Entries {
			item := Exception{
				Name:          x.Name,
				PacketCapture: x.PacketCapture,
				ExemptIps:     util.EntToStr(x.ExemptIps),
			}

			if x.Action != nil {
				switch {
				case x.Action.Default != nil:
					item.Action = ActionDefault
				case x.Action.Allow != nil:
					item.Action = ActionAllow
				case x.Action.Alert != nil:
					item.Action = ActionAlert
				case x.Action.Drop != nil:
					item.Action = ActionDrop
				case x.Action.ResetClient != nil:
					item.Action = ActionResetClient
				case x.Action.ResetServer != nil:
					item.Action = ActionResetServer
				case x.Action.ResetBoth != nil:
					item.Action = ActionResetBoth
				case x.Action.BlockIp != nil:
					item.Action = ActionBlockIp
					item.BlockIpTrackBy = x.Action.BlockIp.TrackBy
					item.BlockIpDuration = x.Action.BlockIp.Duration
				}
			}
			list = append(list, item)
		}
		ans.Exceptions = list
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name    `xml:"entry"`
	Name        string      `xml:"name,attr"`
	Description string      `xml:"description,omitempty"`
	Botnet      *botnet_v1  `xml:"botnet-domains"`
	Rules       *rules      `xml:"rules"`
	Exceptions  *exceptions `xml:"threat-exception"`
}

type botnet_v1 struct {
	Lists            *bnLists_v1     `xml:"lists"`
	Sinkhole         *sinkhole       `xml:"sinkhole"`
	PacketCapture    string          `xml:"packet-capture,omitempty"`
	ThreatExceptions *util.EntryType `xml:"threat-exception"`
}

type bnLists_v1 struct {
	Entries []bnlEntry_v1 `xml:"entry"`
}

type bnlEntry_v1 struct {
	Name   string `xml:"name,attr"`
	Action action `xml:"action"`
}

type action struct {
	Alert    *string `xml:"alert"`
	Allow    *string `xml:"allow"`
	Block    *string `xml:"block"`
	Sinkhole *string `xml:"sinkhole"`
}

type sinkhole struct {
	SinkholeIpv4Address string `xml:"ipv4-address,omitempty"`
	SinkholeIpv6Address string `xml:"ipv6-address,omitempty"`
}

type rules struct {
	Entries []rule `xml:"entry"`
}

type rule struct {
	Name          string           `xml:"name,attr"`
	ThreatName    string           `xml:"threat-name,omitempty"`
	Category      string           `xml:"category"`
	Severities    *util.MemberType `xml:"severity"`
	PacketCapture string           `xml:"packet-capture,omitempty"`
	Action        *ruleAction      `xml:"action"`
}

type ruleAction struct {
	Default     *string  `xml:"default"`
	Allow       *string  `xml:"allow"`
	Alert       *string  `xml:"alert"`
	Drop        *string  `xml:"drop"`
	ResetClient *string  `xml:"reset-client"`
	ResetServer *string  `xml:"reset-server"`
	ResetBoth   *string  `xml:"reset-both"`
	BlockIp     *blockIp `xml:"block-ip"`
}

type blockIp struct {
	TrackBy  string `xml:"track-by"`
	Duration int    `xml:"duration"`
}

type exceptions struct {
	Entries []exception `xml:"entry"`
}

type exception struct {
	Name          string           `xml:"name,attr"`
	PacketCapture string           `xml:"packet-capture,omitempty"`
	Action        *exceptionAction `xml:"action"`
	ExemptIps     *util.EntryType  `xml:"exempt-ip"`
}

type exceptionAction struct {
	Default     *string  `xml:"default"`
	Allow       *string  `xml:"allow"`
	Alert       *string  `xml:"alert"`
	Drop        *string  `xml:"drop"`
	ResetClient *string  `xml:"reset-client"`
	ResetServer *string  `xml:"reset-server"`
	ResetBoth   *string  `xml:"reset-both"`
	BlockIp     *blockIp `xml:"block-ip"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
	}

	s := ""

	if e.PacketCapture != "" || len(e.ThreatExceptions) > 0 || len(e.BotnetLists) > 0 || e.SinkholeIpv4Address != "" || e.SinkholeIpv6Address != "" {
		spec := botnet_v1{
			PacketCapture:    e.PacketCapture,
			ThreatExceptions: util.StrToEnt(e.ThreatExceptions),
		}

		if len(e.BotnetLists) > 0 {
			list := make([]bnlEntry_v1, 0, len(e.BotnetLists))
			for _, x := range e.BotnetLists {
				val := bnlEntry_v1{
					Name: x.Name,
				}

				switch x.Action {
				case ActionAlert:
					val.Action.Alert = &s
				case ActionAllow:
					val.Action.Allow = &s
				case ActionBlock:
					val.Action.Block = &s
				case ActionSinkhole:
					val.Action.Sinkhole = &s
				}

				list = append(list, val)
			}

			spec.Lists = &bnLists_v1{Entries: list}
		}

		if e.SinkholeIpv4Address != "" || e.SinkholeIpv6Address != "" {
			spec.Sinkhole = &sinkhole{
				SinkholeIpv4Address: e.SinkholeIpv4Address,
				SinkholeIpv6Address: e.SinkholeIpv6Address,
			}
		}

		ans.Botnet = &spec
	}

	if len(e.Rules) > 0 {
		list := make([]rule, 0, len(e.Rules))
		for _, x := range e.Rules {
			item := rule{
				Name:          x.Name,
				ThreatName:    x.ThreatName,
				Category:      x.Category,
				Severities:    util.StrToMem(x.Severities),
				PacketCapture: x.PacketCapture,
			}

			switch x.Action {
			case ActionDefault:
				item.Action = &ruleAction{
					Default: &s,
				}
			case ActionAllow:
				item.Action = &ruleAction{
					Allow: &s,
				}
			case ActionAlert:
				item.Action = &ruleAction{
					Alert: &s,
				}
			case ActionDrop:
				item.Action = &ruleAction{
					Drop: &s,
				}
			case ActionResetClient:
				item.Action = &ruleAction{
					ResetClient: &s,
				}
			case ActionResetServer:
				item.Action = &ruleAction{
					ResetServer: &s,
				}
			case ActionResetBoth:
				item.Action = &ruleAction{
					ResetBoth: &s,
				}
			case ActionBlockIp:
				item.Action = &ruleAction{
					BlockIp: &blockIp{
						TrackBy:  x.BlockIpTrackBy,
						Duration: x.BlockIpDuration,
					},
				}
			}

			list = append(list, item)
		}
		ans.Rules = &rules{Entries: list}
	}

	if len(e.Exceptions) > 0 {
		list := make([]exception, 0, len(e.Exceptions))
		for _, x := range e.Exceptions {
			item := exception{
				Name:          x.Name,
				PacketCapture: x.PacketCapture,
				ExemptIps:     util.StrToEnt(x.ExemptIps),
			}

			switch x.Action {
			case ActionDefault:
				item.Action = &exceptionAction{
					Default: &s,
				}
			case ActionAllow:
				item.Action = &exceptionAction{
					Allow: &s,
				}
			case ActionAlert:
				item.Action = &exceptionAction{
					Alert: &s,
				}
			case ActionDrop:
				item.Action = &exceptionAction{
					Drop: &s,
				}
			case ActionResetClient:
				item.Action = &exceptionAction{
					ResetClient: &s,
				}
			case ActionResetServer:
				item.Action = &exceptionAction{
					ResetServer: &s,
				}
			case ActionResetBoth:
				item.Action = &exceptionAction{
					ResetBoth: &s,
				}
			case ActionBlockIp:
				item.Action = &exceptionAction{
					BlockIp: &blockIp{
						TrackBy:  x.BlockIpTrackBy,
						Duration: x.BlockIpDuration,
					},
				}
			}
			list = append(list, item)
		}
		ans.Exceptions = &exceptions{Entries: list}
	}

	return ans
}

// 9.0
type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v2) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:        o.Name,
		Description: o.Description,
	}

	if o.Botnet != nil {
		ans.ThreatExceptions = util.EntToStr(o.Botnet.ThreatExceptions)

		if o.Botnet.Lists != nil {
			lists := make([]BotnetList, 0, len(o.Botnet.Lists.Entries))
			for _, x := range o.Botnet.Lists.Entries {
				val := BotnetList{
					Name:          x.Name,
					PacketCapture: x.PacketCapture,
				}

				switch {
				case x.Action.Alert != nil:
					val.Action = ActionAlert
				case x.Action.Allow != nil:
					val.Action = ActionAllow
				case x.Action.Block != nil:
					val.Action = ActionBlock
				case x.Action.Sinkhole != nil:
					val.Action = ActionSinkhole
				}

				lists = append(lists, val)
			}

			ans.BotnetLists = lists
		}

		if o.Botnet.Sinkhole != nil {
			ans.SinkholeIpv4Address = o.Botnet.Sinkhole.SinkholeIpv4Address
			ans.SinkholeIpv6Address = o.Botnet.Sinkhole.SinkholeIpv6Address
		}
	}

	if o.Rules != nil {
		list := make([]Rule, 0, len(o.Rules.Entries))
		for _, x := range o.Rules.Entries {
			item := Rule{
				Name:          x.Name,
				ThreatName:    x.ThreatName,
				Category:      x.Category,
				Severities:    util.MemToStr(x.Severities),
				PacketCapture: x.PacketCapture,
			}

			if x.Action != nil {
				switch {
				case x.Action.Default != nil:
					item.Action = ActionDefault
				case x.Action.Allow != nil:
					item.Action = ActionAllow
				case x.Action.Alert != nil:
					item.Action = ActionAlert
				case x.Action.Drop != nil:
					item.Action = ActionDrop
				case x.Action.ResetClient != nil:
					item.Action = ActionResetClient
				case x.Action.ResetServer != nil:
					item.Action = ActionResetServer
				case x.Action.ResetBoth != nil:
					item.Action = ActionResetBoth
				case x.Action.BlockIp != nil:
					item.Action = ActionBlockIp
					item.BlockIpTrackBy = x.Action.BlockIp.TrackBy
					item.BlockIpDuration = x.Action.BlockIp.Duration
				}
			}

			list = append(list, item)
		}
		ans.Rules = list
	}

	if o.Exceptions != nil {
		list := make([]Exception, 0, len(o.Exceptions.Entries))
		for _, x := range o.Exceptions.Entries {
			item := Exception{
				Name:          x.Name,
				PacketCapture: x.PacketCapture,
				ExemptIps:     util.EntToStr(x.ExemptIps),
			}

			if x.Action != nil {
				switch {
				case x.Action.Default != nil:
					item.Action = ActionDefault
				case x.Action.Allow != nil:
					item.Action = ActionAllow
				case x.Action.Alert != nil:
					item.Action = ActionAlert
				case x.Action.Drop != nil:
					item.Action = ActionDrop
				case x.Action.ResetClient != nil:
					item.Action = ActionResetClient
				case x.Action.ResetServer != nil:
					item.Action = ActionResetServer
				case x.Action.ResetBoth != nil:
					item.Action = ActionResetBoth
				case x.Action.BlockIp != nil:
					item.Action = ActionBlockIp
					item.BlockIpTrackBy = x.Action.BlockIp.TrackBy
					item.BlockIpDuration = x.Action.BlockIp.Duration
				}
			}
			list = append(list, item)
		}
		ans.Exceptions = list
	}

	return ans
}

type entry_v2 struct {
	XMLName     xml.Name    `xml:"entry"`
	Name        string      `xml:"name,attr"`
	Description string      `xml:"description,omitempty"`
	Botnet      *botnet_v2  `xml:"botnet-domains"`
	Rules       *rules      `xml:"rules"`
	Exceptions  *exceptions `xml:"threat-exception"`
}

type botnet_v2 struct {
	Lists            *bnLists_v2     `xml:"lists"`
	Sinkhole         *sinkhole       `xml:"sinkhole"`
	ThreatExceptions *util.EntryType `xml:"threat-exception"`
}

type bnLists_v2 struct {
	Entries []bnlEntry_v2 `xml:"entry"`
}

type bnlEntry_v2 struct {
	Name          string `xml:"name,attr"`
	Action        action `xml:"action"`
	PacketCapture string `xml:"packet-capture,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:        e.Name,
		Description: e.Description,
	}

	if len(e.ThreatExceptions) > 0 || len(e.BotnetLists) > 0 || e.SinkholeIpv4Address != "" || e.SinkholeIpv6Address != "" {
		spec := botnet_v2{
			ThreatExceptions: util.StrToEnt(e.ThreatExceptions),
		}

		if len(e.BotnetLists) > 0 {
			list := make([]bnlEntry_v2, 0, len(e.BotnetLists))
			for _, x := range e.BotnetLists {
				val := bnlEntry_v2{
					Name:          x.Name,
					PacketCapture: x.PacketCapture,
				}

				s := ""
				switch x.Action {
				case ActionAlert:
					val.Action.Alert = &s
				case ActionAllow:
					val.Action.Allow = &s
				case ActionBlock:
					val.Action.Block = &s
				case ActionSinkhole:
					val.Action.Sinkhole = &s
				}

				list = append(list, val)
			}

			spec.Lists = &bnLists_v2{Entries: list}
		}

		if e.SinkholeIpv4Address != "" || e.SinkholeIpv6Address != "" {
			spec.Sinkhole = &sinkhole{
				SinkholeIpv4Address: e.SinkholeIpv4Address,
				SinkholeIpv6Address: e.SinkholeIpv6Address,
			}
		}

		ans.Botnet = &spec
	}

	s := ""

	if len(e.Rules) > 0 {
		list := make([]rule, 0, len(e.Rules))
		for _, x := range e.Rules {
			item := rule{
				Name:          x.Name,
				ThreatName:    x.ThreatName,
				Category:      x.Category,
				Severities:    util.StrToMem(x.Severities),
				PacketCapture: x.PacketCapture,
			}

			switch x.Action {
			case ActionDefault:
				item.Action = &ruleAction{
					Default: &s,
				}
			case ActionAllow:
				item.Action = &ruleAction{
					Allow: &s,
				}
			case ActionAlert:
				item.Action = &ruleAction{
					Alert: &s,
				}
			case ActionDrop:
				item.Action = &ruleAction{
					Drop: &s,
				}
			case ActionResetClient:
				item.Action = &ruleAction{
					ResetClient: &s,
				}
			case ActionResetServer:
				item.Action = &ruleAction{
					ResetServer: &s,
				}
			case ActionResetBoth:
				item.Action = &ruleAction{
					ResetBoth: &s,
				}
			case ActionBlockIp:
				item.Action = &ruleAction{
					BlockIp: &blockIp{
						TrackBy:  x.BlockIpTrackBy,
						Duration: x.BlockIpDuration,
					},
				}
			}

			list = append(list, item)
		}
		ans.Rules = &rules{Entries: list}
	}

	if len(e.Exceptions) > 0 {
		list := make([]exception, 0, len(e.Exceptions))
		for _, x := range e.Exceptions {
			item := exception{
				Name:          x.Name,
				PacketCapture: x.PacketCapture,
				ExemptIps:     util.StrToEnt(x.ExemptIps),
			}

			switch x.Action {
			case ActionDefault:
				item.Action = &exceptionAction{
					Default: &s,
				}
			case ActionAllow:
				item.Action = &exceptionAction{
					Allow: &s,
				}
			case ActionAlert:
				item.Action = &exceptionAction{
					Alert: &s,
				}
			case ActionDrop:
				item.Action = &exceptionAction{
					Drop: &s,
				}
			case ActionResetClient:
				item.Action = &exceptionAction{
					ResetClient: &s,
				}
			case ActionResetServer:
				item.Action = &exceptionAction{
					ResetServer: &s,
				}
			case ActionResetBoth:
				item.Action = &exceptionAction{
					ResetBoth: &s,
				}
			case ActionBlockIp:
				item.Action = &exceptionAction{
					BlockIp: &blockIp{
						TrackBy:  x.BlockIpTrackBy,
						Duration: x.BlockIpDuration,
					},
				}
			}
			list = append(list, item)
		}
		ans.Exceptions = &exceptions{Entries: list}
	}

	return ans
}

// 10.0
type container_v3 struct {
	Answer []entry_v3 `xml:"entry"`
}

func (o *container_v3) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v3) Normalize() []Entry {
	arr := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		arr = append(arr, o.Answer[i].normalize())
	}
	return arr
}

func (o *entry_v3) normalize() Entry {
	ans := Entry{
		Name:        o.Name,
		Description: o.Description,
	}

	if o.Botnet != nil {
		ans.ThreatExceptions = util.EntToStr(o.Botnet.ThreatExceptions)

		if o.Botnet.Lists != nil {
			lists := make([]BotnetList, 0, len(o.Botnet.Lists.Entries))
			for _, x := range o.Botnet.Lists.Entries {
				val := BotnetList{
					Name:          x.Name,
					PacketCapture: x.PacketCapture,
				}

				switch {
				case x.Action.Alert != nil:
					val.Action = ActionAlert
				case x.Action.Allow != nil:
					val.Action = ActionAllow
				case x.Action.Block != nil:
					val.Action = ActionBlock
				case x.Action.Sinkhole != nil:
					val.Action = ActionSinkhole
				}

				lists = append(lists, val)
			}

			ans.BotnetLists = lists
		}

		if o.Botnet.Dns != nil {
			list := make([]DnsCategory, 0, len(o.Botnet.Dns.Entries))
			for _, x := range o.Botnet.Dns.Entries {
				list = append(list, DnsCategory{
					Name:          x.Name,
					Action:        x.Action,
					LogLevel:      x.LogLevel,
					PacketCapture: x.PacketCapture,
				})
			}

			ans.DnsCategories = list
		}

		if o.Botnet.WhiteLists != nil {
			list := make([]WhiteList, 0, len(o.Botnet.WhiteLists.Entries))
			for _, x := range o.Botnet.WhiteLists.Entries {
				list = append(list, WhiteList{
					Name:        x.Name,
					Description: x.Description,
				})
			}

			ans.WhiteLists = list
		}

		if o.Botnet.Sinkhole != nil {
			ans.SinkholeIpv4Address = o.Botnet.Sinkhole.SinkholeIpv4Address
			ans.SinkholeIpv6Address = o.Botnet.Sinkhole.SinkholeIpv6Address
		}
	}

	if o.Rules != nil {
		list := make([]Rule, 0, len(o.Rules.Entries))
		for _, x := range o.Rules.Entries {
			item := Rule{
				Name:          x.Name,
				ThreatName:    x.ThreatName,
				Category:      x.Category,
				Severities:    util.MemToStr(x.Severities),
				PacketCapture: x.PacketCapture,
			}

			if x.Action != nil {
				switch {
				case x.Action.Default != nil:
					item.Action = ActionDefault
				case x.Action.Allow != nil:
					item.Action = ActionAllow
				case x.Action.Alert != nil:
					item.Action = ActionAlert
				case x.Action.Drop != nil:
					item.Action = ActionDrop
				case x.Action.ResetClient != nil:
					item.Action = ActionResetClient
				case x.Action.ResetServer != nil:
					item.Action = ActionResetServer
				case x.Action.ResetBoth != nil:
					item.Action = ActionResetBoth
				case x.Action.BlockIp != nil:
					item.Action = ActionBlockIp
					item.BlockIpTrackBy = x.Action.BlockIp.TrackBy
					item.BlockIpDuration = x.Action.BlockIp.Duration
				}
			}

			list = append(list, item)
		}
		ans.Rules = list
	}

	if o.Exceptions != nil {
		list := make([]Exception, 0, len(o.Exceptions.Entries))
		for _, x := range o.Exceptions.Entries {
			item := Exception{
				Name:          x.Name,
				PacketCapture: x.PacketCapture,
				ExemptIps:     util.EntToStr(x.ExemptIps),
			}

			if x.Action != nil {
				switch {
				case x.Action.Default != nil:
					item.Action = ActionDefault
				case x.Action.Allow != nil:
					item.Action = ActionAllow
				case x.Action.Alert != nil:
					item.Action = ActionAlert
				case x.Action.Drop != nil:
					item.Action = ActionDrop
				case x.Action.ResetClient != nil:
					item.Action = ActionResetClient
				case x.Action.ResetServer != nil:
					item.Action = ActionResetServer
				case x.Action.ResetBoth != nil:
					item.Action = ActionResetBoth
				case x.Action.BlockIp != nil:
					item.Action = ActionBlockIp
					item.BlockIpTrackBy = x.Action.BlockIp.TrackBy
					item.BlockIpDuration = x.Action.BlockIp.Duration
				}
			}
			list = append(list, item)
		}
		ans.Exceptions = list
	}

	return ans
}

type entry_v3 struct {
	XMLName     xml.Name    `xml:"entry"`
	Name        string      `xml:"name,attr"`
	Description string      `xml:"description,omitempty"`
	Botnet      *botnet_v3  `xml:"botnet-domains"`
	Rules       *rules      `xml:"rules"`
	Exceptions  *exceptions `xml:"threat-exception"`
}

type botnet_v3 struct {
	Lists            *bnLists_v2     `xml:"lists"`
	Dns              *dns            `xml:"dns-security-categories"`
	WhiteLists       *whiteList      `xml:"whitelist"`
	Sinkhole         *sinkhole       `xml:"sinkhole"`
	ThreatExceptions *util.EntryType `xml:"threat-exception"`
}

type dns struct {
	Entries []dnsEntry `xml:"entry"`
}

type dnsEntry struct {
	Name          string `xml:"name,attr"`
	Action        string `xml:"action,omitempty"`
	LogLevel      string `xml:"log-level,omitempty"`
	PacketCapture string `xml:"packet-capture,omitempty"`
}

type whiteList struct {
	Entries []wlEntry `xml:"entry"`
}

type wlEntry struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,omitempty"`
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:        e.Name,
		Description: e.Description,
	}

	if len(e.ThreatExceptions) > 0 || len(e.BotnetLists) > 0 || e.SinkholeIpv4Address != "" || e.SinkholeIpv6Address != "" || len(e.DnsCategories) > 0 || len(e.WhiteLists) > 0 {
		spec := botnet_v3{
			ThreatExceptions: util.StrToEnt(e.ThreatExceptions),
		}

		if len(e.DnsCategories) > 0 {
			list := make([]dnsEntry, 0, len(e.DnsCategories))
			for _, x := range e.DnsCategories {
				list = append(list, dnsEntry{
					Name:          x.Name,
					Action:        x.Action,
					LogLevel:      x.LogLevel,
					PacketCapture: x.PacketCapture,
				})
			}

			spec.Dns = &dns{Entries: list}
		}

		if len(e.WhiteLists) > 0 {
			list := make([]wlEntry, 0, len(e.WhiteLists))
			for _, x := range e.WhiteLists {
				list = append(list, wlEntry{
					Name:        x.Name,
					Description: x.Description,
				})
			}

			spec.WhiteLists = &whiteList{Entries: list}
		}

		if len(e.BotnetLists) > 0 {
			list := make([]bnlEntry_v2, 0, len(e.BotnetLists))
			for _, x := range e.BotnetLists {
				val := bnlEntry_v2{
					Name:          x.Name,
					PacketCapture: x.PacketCapture,
				}

				s := ""
				switch x.Action {
				case ActionAlert:
					val.Action.Alert = &s
				case ActionAllow:
					val.Action.Allow = &s
				case ActionBlock:
					val.Action.Block = &s
				case ActionSinkhole:
					val.Action.Sinkhole = &s
				}

				list = append(list, val)
			}

			spec.Lists = &bnLists_v2{Entries: list}
		}

		if e.SinkholeIpv4Address != "" || e.SinkholeIpv6Address != "" {
			spec.Sinkhole = &sinkhole{
				SinkholeIpv4Address: e.SinkholeIpv4Address,
				SinkholeIpv6Address: e.SinkholeIpv6Address,
			}
		}

		ans.Botnet = &spec
	}

	s := ""

	if len(e.Rules) > 0 {
		list := make([]rule, 0, len(e.Rules))
		for _, x := range e.Rules {
			item := rule{
				Name:          x.Name,
				ThreatName:    x.ThreatName,
				Category:      x.Category,
				Severities:    util.StrToMem(x.Severities),
				PacketCapture: x.PacketCapture,
			}

			switch x.Action {
			case ActionDefault:
				item.Action = &ruleAction{
					Default: &s,
				}
			case ActionAllow:
				item.Action = &ruleAction{
					Allow: &s,
				}
			case ActionAlert:
				item.Action = &ruleAction{
					Alert: &s,
				}
			case ActionDrop:
				item.Action = &ruleAction{
					Drop: &s,
				}
			case ActionResetClient:
				item.Action = &ruleAction{
					ResetClient: &s,
				}
			case ActionResetServer:
				item.Action = &ruleAction{
					ResetServer: &s,
				}
			case ActionResetBoth:
				item.Action = &ruleAction{
					ResetBoth: &s,
				}
			case ActionBlockIp:
				item.Action = &ruleAction{
					BlockIp: &blockIp{
						TrackBy:  x.BlockIpTrackBy,
						Duration: x.BlockIpDuration,
					},
				}
			}

			list = append(list, item)
		}
		ans.Rules = &rules{Entries: list}
	}

	if len(e.Exceptions) > 0 {
		list := make([]exception, 0, len(e.Exceptions))
		for _, x := range e.Exceptions {
			item := exception{
				Name:          x.Name,
				PacketCapture: x.PacketCapture,
				ExemptIps:     util.StrToEnt(x.ExemptIps),
			}

			switch x.Action {
			case ActionDefault:
				item.Action = &exceptionAction{
					Default: &s,
				}
			case ActionAllow:
				item.Action = &exceptionAction{
					Allow: &s,
				}
			case ActionAlert:
				item.Action = &exceptionAction{
					Alert: &s,
				}
			case ActionDrop:
				item.Action = &exceptionAction{
					Drop: &s,
				}
			case ActionResetClient:
				item.Action = &exceptionAction{
					ResetClient: &s,
				}
			case ActionResetServer:
				item.Action = &exceptionAction{
					ResetServer: &s,
				}
			case ActionResetBoth:
				item.Action = &exceptionAction{
					ResetBoth: &s,
				}
			case ActionBlockIp:
				item.Action = &exceptionAction{
					BlockIp: &blockIp{
						TrackBy:  x.BlockIpTrackBy,
						Duration: x.BlockIpDuration,
					},
				}
			}
			list = append(list, item)
		}
		ans.Exceptions = &exceptions{Entries: list}
	}

	return ans
}
