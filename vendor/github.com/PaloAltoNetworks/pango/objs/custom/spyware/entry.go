package spyware

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// custom data pattern object.
//
// PAN-OS 7.0+.
type Entry struct {
	Name                     string
	ThreatName               string
	Comment                  string
	Severity                 string
	Direction                string
	DefaultAction            string
	BlockIpTrackBy           string
	BlockIpDuration          int
	Cves                     []string
	Bugtraqs                 []string
	Vendors                  []string
	References               []string
	StandardSignatureType    *StandardSignatureType
	CombinationSignatureType *CombinationSignatureType
}

type StandardSignatureType struct {
	Signatures []StandardSignature
}

type StandardSignature struct {
	Name         string
	Comment      string
	Scope        string
	OrderFree    bool
	StandardAnds []StandardAnd
}

type StandardAnd struct {
	Name        string
	StandardOrs []StandardOr
}

type StandardOr struct {
	Name         string
	LessThan     *Condition
	EqualTo      *EqualTo
	GreaterThan  *Condition
	PatternMatch *Pattern
}

type Condition struct {
	Value      int
	Context    string
	Qualifiers []Qualifier
}

type EqualTo struct {
	Value      int
	Context    string
	Negate     bool // PAN-OS 10.0
	Qualifiers []Qualifier
}

type Qualifier struct {
	Qualifier string
	Value     string
}

type Pattern struct {
	Pattern    string
	Context    string
	Negate     bool
	Qualifiers []Qualifier
}

type CombinationSignatureType struct {
	OrderFree           bool
	ThresholdTime       int
	IntervalTime        int
	AggregationCriteria string
	Signatures          []CombinationSignature
}

type CombinationSignature struct {
	Name           string
	CombinationOrs []CombinationOr
}

type CombinationOr struct {
	Name     string
	ThreatId string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.ThreatName = s.ThreatName
	o.Comment = s.Comment
	o.Severity = s.Severity
	o.Direction = s.Direction
	o.DefaultAction = s.DefaultAction
	o.BlockIpTrackBy = s.BlockIpTrackBy
	o.BlockIpDuration = s.BlockIpDuration
	if s.Cves == nil {
		o.Cves = nil
	} else {
		o.Cves = make([]string, len(s.Cves))
		copy(o.Cves, s.Cves)
	}
	if s.Bugtraqs == nil {
		o.Bugtraqs = nil
	} else {
		o.Bugtraqs = make([]string, len(s.Bugtraqs))
		copy(o.Bugtraqs, s.Bugtraqs)
	}
	if s.Vendors == nil {
		o.Vendors = nil
	} else {
		o.Vendors = make([]string, len(s.Vendors))
		copy(o.Vendors, s.Vendors)
	}
	if s.References == nil {
		o.References = nil
	} else {
		o.References = make([]string, len(s.References))
		copy(o.References, s.References)
	}

	if s.StandardSignatureType == nil {
		o.StandardSignatureType = nil
	} else {
		o.StandardSignatureType = &StandardSignatureType{}
		if len(s.StandardSignatureType.Signatures) > 0 {
			sstl := make([]StandardSignature, 0, len(s.StandardSignatureType.Signatures))
			for _, x := range s.StandardSignatureType.Signatures {
				var acl []StandardAnd
				if x.StandardAnds != nil {
					acl = make([]StandardAnd, 0, len(x.StandardAnds))
					for _, aobj := range x.StandardAnds {
						var olist []StandardOr
						if aobj.StandardOrs != nil {
							olist = make([]StandardOr, 0, len(aobj.StandardOrs))
							for _, oobj := range aobj.StandardOrs {
								olist = append(olist, StandardOr{
									Name:         oobj.Name,
									LessThan:     copyCondition(oobj.LessThan),
									EqualTo:      copyEqualTo(oobj.EqualTo),
									GreaterThan:  copyCondition(oobj.GreaterThan),
									PatternMatch: copyPattern(oobj.PatternMatch),
								})
							}
						}
						acl = append(acl, StandardAnd{
							Name:        aobj.Name,
							StandardOrs: olist,
						})
					}
				}
				sstl = append(sstl, StandardSignature{
					Name:         x.Name,
					Comment:      x.Comment,
					Scope:        x.Scope,
					OrderFree:    x.OrderFree,
					StandardAnds: acl,
				})
			}
			o.StandardSignatureType.Signatures = sstl
		}
	}

	if s.CombinationSignatureType == nil {
		o.CombinationSignatureType = nil
	} else {
		var slist []CombinationSignature
		if s.CombinationSignatureType.Signatures != nil {
			slist = make([]CombinationSignature, 0, len(s.CombinationSignatureType.Signatures))
			for _, sobj := range s.CombinationSignatureType.Signatures {
				var olist []CombinationOr
				if sobj.CombinationOrs != nil {
					olist = make([]CombinationOr, 0, len(sobj.CombinationOrs))
					for _, oobj := range sobj.CombinationOrs {
						olist = append(olist, CombinationOr{
							Name:     oobj.Name,
							ThreatId: oobj.ThreatId,
						})
					}
				}
				slist = append(slist, CombinationSignature{
					Name:           sobj.Name,
					CombinationOrs: olist,
				})
			}
		}
		o.CombinationSignatureType = &CombinationSignatureType{
			OrderFree:           s.CombinationSignatureType.OrderFree,
			ThresholdTime:       s.CombinationSignatureType.ThresholdTime,
			IntervalTime:        s.CombinationSignatureType.IntervalTime,
			AggregationCriteria: s.CombinationSignatureType.AggregationCriteria,
			Signatures:          slist,
		}
	}
}

func copyCondition(v *Condition) *Condition {
	if v == nil {
		return nil
	}

	return &Condition{
		Value:      v.Value,
		Context:    v.Context,
		Qualifiers: copyQualifiers(v.Qualifiers),
	}
}

func copyQualifiers(v []Qualifier) []Qualifier {
	if v == nil {
		return nil
	}

	list := make([]Qualifier, 0, len(v))
	for _, x := range v {
		list = append(list, Qualifier{
			Qualifier: x.Qualifier,
			Value:     x.Value,
		})
	}

	return list
}

func copyEqualTo(v *EqualTo) *EqualTo {
	if v == nil {
		return nil
	}

	return &EqualTo{
		Value:      v.Value,
		Context:    v.Context,
		Negate:     v.Negate,
		Qualifiers: copyQualifiers(v.Qualifiers),
	}
}

func copyPattern(v *Pattern) *Pattern {
	if v == nil {
		return nil
	}

	return &Pattern{
		Pattern:    v.Pattern,
		Context:    v.Context,
		Negate:     v.Negate,
		Qualifiers: copyQualifiers(v.Qualifiers),
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
		Name:       o.Name,
		ThreatName: o.ThreatName,
		Comment:    o.Comment,
		Severity:   o.Severity,
		Direction:  o.Direction,
		Cves:       util.MemToStr(o.Cves),
		Bugtraqs:   util.MemToStr(o.Bugtraqs),
		Vendors:    util.MemToStr(o.Vendors),
		References: util.MemToStr(o.References),
	}

	if o.DefaultAction == nil {
		ans.DefaultAction = DefaultActionAlert
	} else {
		switch {
		case o.DefaultAction.Allow != nil:
			ans.DefaultAction = DefaultActionAllow
		case o.DefaultAction.Alert != nil:
			ans.DefaultAction = DefaultActionAlert
		case o.DefaultAction.Drop != nil:
			ans.DefaultAction = DefaultActionDrop
		case o.DefaultAction.ResetClient != nil:
			ans.DefaultAction = DefaultActionResetClient
		case o.DefaultAction.ResetServer != nil:
			ans.DefaultAction = DefaultActionResetServer
		case o.DefaultAction.ResetBoth != nil:
			ans.DefaultAction = DefaultActionResetBoth
		case o.DefaultAction.BlockIp != nil:
			ans.DefaultAction = DefaultActionBlockIp
			ans.BlockIpTrackBy = o.DefaultAction.BlockIp.BlockIpTrackBy
			ans.BlockIpDuration = o.DefaultAction.BlockIp.BlockIpDuration
		}
	}

	if o.Signature.Standard != nil {
		ans.StandardSignatureType = &StandardSignatureType{}
		if len(o.Signature.Standard.Entries) > 0 {
			sigs := make([]StandardSignature, 0, len(o.Signature.Standard.Entries))
			for _, x := range o.Signature.Standard.Entries {
				var alist []StandardAnd
				if x.Ands != nil && len(x.Ands.Entries) > 0 {
					alist = make([]StandardAnd, 0, len(x.Ands.Entries))
					for _, aobj := range x.Ands.Entries {
						var olist []StandardOr
						if aobj.Ors != nil && len(aobj.Ors.Entries) > 0 {
							olist = make([]StandardOr, 0, len(aobj.Ors.Entries))
							for _, oobj := range aobj.Ors.Entries {
								olist = append(olist, StandardOr{
									Name:         oobj.Name,
									LessThan:     normalizeCondition(oobj.Op.LessThan),
									EqualTo:      normalizeEqualToFromCondition(oobj.Op.EqualTo),
									GreaterThan:  normalizeCondition(oobj.Op.GreaterThan),
									PatternMatch: normalizePatternMatch(oobj.Op.PatternMatch),
								})
							}
						}
						alist = append(alist, StandardAnd{
							Name:        aobj.Name,
							StandardOrs: olist,
						})
					}
				}
				sigs = append(sigs, StandardSignature{
					Name:         x.Name,
					Comment:      x.Comment,
					Scope:        x.Scope,
					OrderFree:    util.AsBool(x.OrderFree),
					StandardAnds: alist,
				})
			}
			ans.StandardSignatureType.Signatures = sigs
		}
	}

	if o.Signature.Combination != nil {
		ans.CombinationSignatureType = &CombinationSignatureType{
			OrderFree:           util.AsBool(o.Signature.Combination.OrderFree),
			ThresholdTime:       o.Signature.Combination.TimeAttribute.ThresholdTime,
			IntervalTime:        o.Signature.Combination.TimeAttribute.IntervalTime,
			AggregationCriteria: o.Signature.Combination.TimeAttribute.AggregationCriteria,
		}

		if o.Signature.Combination.AndConditions != nil && len(o.Signature.Combination.AndConditions.Entries) > 0 {
			list := make([]CombinationSignature, 0, len(o.Signature.Combination.AndConditions.Entries))
			for _, x := range o.Signature.Combination.AndConditions.Entries {
				var olist []CombinationOr
				if x.OrConditions != nil && len(x.OrConditions.Entries) > 0 {
					olist = make([]CombinationOr, 0, len(x.OrConditions.Entries))
					for _, oobj := range x.OrConditions.Entries {
						olist = append(olist, CombinationOr{
							Name:     oobj.Name,
							ThreatId: oobj.ThreatId,
						})
					}
				}
				list = append(list, CombinationSignature{
					Name:           x.Name,
					CombinationOrs: olist,
				})
			}
			ans.CombinationSignatureType.Signatures = list
		}
	}

	return ans
}

func normalizeCondition(v *condition) *Condition {
	if v == nil {
		return nil
	}

	return &Condition{
		Value:      v.Value,
		Context:    v.Context,
		Qualifiers: normalizeQualifiers(v.Qualifiers),
	}
}

func normalizeQualifiers(v *qualifiers) []Qualifier {
	if v == nil || len(v.Entries) == 0 {
		return nil
	}

	list := make([]Qualifier, 0, len(v.Entries))
	for _, x := range v.Entries {
		list = append(list, Qualifier{
			Qualifier: x.Qualifier,
			Value:     x.Value,
		})
	}

	return list
}

func normalizeEqualToFromCondition(v *condition) *EqualTo {
	if v == nil {
		return nil
	}

	return &EqualTo{
		Value:      v.Value,
		Context:    v.Context,
		Qualifiers: normalizeQualifiers(v.Qualifiers),
	}
}

func normalizePatternMatch(v *patternMatch) *Pattern {
	if v == nil {
		return nil
	}

	return &Pattern{
		Pattern:    v.Pattern,
		Context:    v.Context,
		Negate:     util.AsBool(v.Negate),
		Qualifiers: normalizeQualifiers(v.Qualifiers),
	}
}

type entry_v1 struct {
	XMLName       xml.Name         `xml:"entry"`
	Name          string           `xml:"name,attr"`
	ThreatName    string           `xml:"threatname"`
	Comment       string           `xml:"comment,omitempty"`
	Severity      string           `xml:"severity"`
	Direction     string           `xml:"direction"`
	DefaultAction *action          `xml:"default-action"`
	Cves          *util.MemberType `xml:"cve"`
	Bugtraqs      *util.MemberType `xml:"bugtraq"`
	Vendors       *util.MemberType `xml:"vendor"`
	References    *util.MemberType `xml:"reference"`
	Signature     sig_v1           `xml:"signature"`
}

type action struct {
	Allow       *string  `xml:"allow"`
	Alert       *string  `xml:"alert"`
	Drop        *string  `xml:"drop"`
	ResetClient *string  `xml:"reset-client"`
	ResetServer *string  `xml:"reset-server"`
	ResetBoth   *string  `xml:"reset-both"`
	BlockIp     *blockIp `xml:"block-ip"`
}

type blockIp struct {
	BlockIpTrackBy  string `xml:"track-by"`
	BlockIpDuration int    `xml:"duration"`
}

type sig_v1 struct {
	Standard    *standard_v1 `xml:"standard"`
	Combination *combination `xml:"combination"`
}

type standard_v1 struct {
	Entries []standardEntry_v1 `xml:"entry"`
}

type standardEntry_v1 struct {
	Name      string   `xml:"name,attr"`
	Comment   string   `xml:"comment,omitempty"`
	Scope     string   `xml:"scope,omitempty"`
	OrderFree string   `xml:"order-free"`
	Ands      *ands_v1 `xml:"and-condition"`
}

type ands_v1 struct {
	Entries []standardAnd_v1 `xml:"entry"`
}

type standardAnd_v1 struct {
	Name string  `xml:"name,attr"`
	Ors  *ors_v1 `xml:"or-condition"`
}

type ors_v1 struct {
	Entries []standardOr_v1 `xml:"entry"`
}

type standardOr_v1 struct {
	Name string        `xml:"name,attr"`
	Op   standardOp_v1 `xml:"operator"`
}

type standardOp_v1 struct {
	LessThan     *condition    `xml:"less-than"`
	EqualTo      *condition    `xml:"equal-to"`
	GreaterThan  *condition    `xml:"greater-than"`
	PatternMatch *patternMatch `xml:"pattern-match"`
}

type condition struct {
	Value      int         `xml:"value"`
	Context    string      `xml:"context"`
	Qualifiers *qualifiers `xml:"qualifier"`
}

type qualifiers struct {
	Entries []qualifier `xml:"entry"`
}

type qualifier struct {
	Qualifier string `xml:"name,attr"`
	Value     string `xml:"value"`
}

type combination struct {
	OrderFree     string    `xml:"order-free"`
	TimeAttribute comboTime `xml:"time-attribute"`
	AndConditions *andCond  `xml:"and-condition"`
}

type patternMatch struct {
	Pattern    string      `xml:"pattern"`
	Context    string      `xml:"context"`
	Negate     string      `xml:"negate"`
	Qualifiers *qualifiers `xml:"qualifier"`
}

type comboTime struct {
	ThresholdTime       int    `xml:"threshold"`
	IntervalTime        int    `xml:"interval"`
	AggregationCriteria string `xml:"track-by,omitempty"`
}

type andCond struct {
	Entries []andCondEntry `xml:"entry"`
}

type andCondEntry struct {
	Name         string  `xml:"name,attr"`
	OrConditions *orCond `xml:"or-condition"`
}

type orCond struct {
	Entries []orCondEntry `xml:"entry"`
}

type orCondEntry struct {
	Name     string `xml:"name,attr"`
	ThreatId string `xml:"threat-id"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:          e.Name,
		ThreatName:    e.ThreatName,
		Comment:       e.Comment,
		Severity:      e.Severity,
		Direction:     e.Direction,
		DefaultAction: &action{},
		Cves:          util.StrToMem(e.Cves),
		Bugtraqs:      util.StrToMem(e.Bugtraqs),
		Vendors:       util.StrToMem(e.Vendors),
		References:    util.StrToMem(e.References),
	}

	s := ""
	switch e.DefaultAction {
	case DefaultActionAllow:
		ans.DefaultAction.Allow = &s
	case DefaultActionAlert:
		ans.DefaultAction.Alert = &s
	case DefaultActionDrop:
		ans.DefaultAction.Drop = &s
	case DefaultActionResetClient:
		ans.DefaultAction.ResetClient = &s
	case DefaultActionResetServer:
		ans.DefaultAction.ResetServer = &s
	case DefaultActionResetBoth:
		ans.DefaultAction.ResetBoth = &s
	case DefaultActionBlockIp:
		ans.DefaultAction.BlockIp = &blockIp{
			BlockIpTrackBy:  e.BlockIpTrackBy,
			BlockIpDuration: e.BlockIpDuration,
		}
	default:
		ans.DefaultAction = nil
	}

	if e.StandardSignatureType != nil {
		ans.Signature.Standard = &standard_v1{}
		if len(e.StandardSignatureType.Signatures) > 0 {
			sigList := make([]standardEntry_v1, 0, len(e.StandardSignatureType.Signatures))
			for _, ss := range e.StandardSignatureType.Signatures {
				var av *ands_v1
				if len(ss.StandardAnds) > 0 {
					alist := make([]standardAnd_v1, 0, len(ss.StandardAnds))
					for _, aobj := range ss.StandardAnds {
						var ov *ors_v1
						if len(aobj.StandardOrs) > 0 {
							olist := make([]standardOr_v1, 0, len(aobj.StandardOrs))
							for _, oobj := range aobj.StandardOrs {
								olist = append(olist, standardOr_v1{
									Name: oobj.Name,
									Op: standardOp_v1{
										LessThan:     specifyCondition(oobj.LessThan),
										EqualTo:      specifyEqualToAsCondition(oobj.EqualTo),
										GreaterThan:  specifyCondition(oobj.GreaterThan),
										PatternMatch: specifyPatternMatch(oobj.PatternMatch),
									},
								})
							}
							ov = &ors_v1{Entries: olist}
						}
						alist = append(alist, standardAnd_v1{
							Name: aobj.Name,
							Ors:  ov,
						})
					}
					av = &ands_v1{Entries: alist}
				}
				sigList = append(sigList, standardEntry_v1{
					Name:      ss.Name,
					Comment:   ss.Comment,
					Scope:     ss.Scope,
					OrderFree: util.YesNo(ss.OrderFree),
					Ands:      av,
				})
			}
			ans.Signature.Standard.Entries = sigList
		}
	}

	if e.CombinationSignatureType != nil {
		ans.Signature.Combination = &combination{
			OrderFree: util.YesNo(e.CombinationSignatureType.OrderFree),
			TimeAttribute: comboTime{
				ThresholdTime:       e.CombinationSignatureType.ThresholdTime,
				IntervalTime:        e.CombinationSignatureType.IntervalTime,
				AggregationCriteria: e.CombinationSignatureType.AggregationCriteria,
			},
		}

		if len(e.CombinationSignatureType.Signatures) > 0 {
			list := make([]andCondEntry, 0, len(e.CombinationSignatureType.Signatures))
			for _, x := range e.CombinationSignatureType.Signatures {
				ors := make([]orCondEntry, 0, len(x.CombinationOrs))
				for _, y := range x.CombinationOrs {
					ors = append(ors, orCondEntry{
						Name:     y.Name,
						ThreatId: y.ThreatId,
					})
				}

				list = append(list, andCondEntry{
					Name:         x.Name,
					OrConditions: &orCond{Entries: ors},
				})
			}

			ans.Signature.Combination.AndConditions = &andCond{
				Entries: list,
			}
		}
	}

	return ans
}

func specifyCondition(v *Condition) *condition {
	if v == nil {
		return nil
	}

	return &condition{
		Value:      v.Value,
		Context:    v.Context,
		Qualifiers: specifyQualifiers(v.Qualifiers),
	}
}

func specifyEqualToAsCondition(v *EqualTo) *condition {
	if v == nil {
		return nil
	}

	return &condition{
		Value:      v.Value,
		Context:    v.Context,
		Qualifiers: specifyQualifiers(v.Qualifiers),
	}
}

func specifyQualifiers(v []Qualifier) *qualifiers {
	if len(v) == 0 {
		return nil
	}

	list := make([]qualifier, 0, len(v))
	for _, x := range v {
		list = append(list, qualifier{
			Qualifier: x.Qualifier,
			Value:     x.Value,
		})
	}

	return &qualifiers{Entries: list}
}

func specifyPatternMatch(v *Pattern) *patternMatch {
	if v == nil {
		return nil
	}

	return &patternMatch{
		Pattern:    v.Pattern,
		Context:    v.Context,
		Negate:     util.YesNo(v.Negate),
		Qualifiers: specifyQualifiers(v.Qualifiers),
	}
}

// PAN-OS 10.0
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
		Name:       o.Name,
		ThreatName: o.ThreatName,
		Comment:    o.Comment,
		Severity:   o.Severity,
		Direction:  o.Direction,
		Cves:       util.MemToStr(o.Cves),
		Bugtraqs:   util.MemToStr(o.Bugtraqs),
		Vendors:    util.MemToStr(o.Vendors),
		References: util.MemToStr(o.References),
	}

	if o.DefaultAction == nil {
		ans.DefaultAction = DefaultActionAlert
	} else {
		switch {
		case o.DefaultAction.Allow != nil:
			ans.DefaultAction = DefaultActionAllow
		case o.DefaultAction.Alert != nil:
			ans.DefaultAction = DefaultActionAlert
		case o.DefaultAction.Drop != nil:
			ans.DefaultAction = DefaultActionDrop
		case o.DefaultAction.ResetClient != nil:
			ans.DefaultAction = DefaultActionResetClient
		case o.DefaultAction.ResetServer != nil:
			ans.DefaultAction = DefaultActionResetServer
		case o.DefaultAction.ResetBoth != nil:
			ans.DefaultAction = DefaultActionResetBoth
		case o.DefaultAction.BlockIp != nil:
			ans.DefaultAction = DefaultActionBlockIp
			ans.BlockIpTrackBy = o.DefaultAction.BlockIp.BlockIpTrackBy
			ans.BlockIpDuration = o.DefaultAction.BlockIp.BlockIpDuration
		}
	}

	if o.Signature.Standard != nil {
		ans.StandardSignatureType = &StandardSignatureType{}
		if len(o.Signature.Standard.Entries) > 0 {
			sigs := make([]StandardSignature, 0, len(o.Signature.Standard.Entries))
			for _, x := range o.Signature.Standard.Entries {
				var alist []StandardAnd
				if x.Ands != nil && len(x.Ands.Entries) > 0 {
					alist = make([]StandardAnd, 0, len(x.Ands.Entries))
					for _, aobj := range x.Ands.Entries {
						var olist []StandardOr
						if aobj.Ors != nil && len(aobj.Ors.Entries) > 0 {
							olist = make([]StandardOr, 0, len(aobj.Ors.Entries))
							for _, oobj := range aobj.Ors.Entries {
								olist = append(olist, StandardOr{
									Name:         oobj.Name,
									LessThan:     normalizeCondition(oobj.Op.LessThan),
									EqualTo:      normalizeEqualToFromEqualTo(oobj.Op.EqualTo),
									GreaterThan:  normalizeCondition(oobj.Op.GreaterThan),
									PatternMatch: normalizePatternMatch(oobj.Op.PatternMatch),
								})
							}
						}
						alist = append(alist, StandardAnd{
							Name:        aobj.Name,
							StandardOrs: olist,
						})
					}
				}
				sigs = append(sigs, StandardSignature{
					Name:         x.Name,
					Comment:      x.Comment,
					Scope:        x.Scope,
					OrderFree:    util.AsBool(x.OrderFree),
					StandardAnds: alist,
				})
			}
			ans.StandardSignatureType.Signatures = sigs
		}
	}

	if o.Signature.Combination != nil {
		ans.CombinationSignatureType = &CombinationSignatureType{
			OrderFree:           util.AsBool(o.Signature.Combination.OrderFree),
			ThresholdTime:       o.Signature.Combination.TimeAttribute.ThresholdTime,
			IntervalTime:        o.Signature.Combination.TimeAttribute.IntervalTime,
			AggregationCriteria: o.Signature.Combination.TimeAttribute.AggregationCriteria,
		}

		if o.Signature.Combination.AndConditions != nil && len(o.Signature.Combination.AndConditions.Entries) > 0 {
			list := make([]CombinationSignature, 0, len(o.Signature.Combination.AndConditions.Entries))
			for _, x := range o.Signature.Combination.AndConditions.Entries {
				var olist []CombinationOr
				if x.OrConditions != nil && len(x.OrConditions.Entries) > 0 {
					olist = make([]CombinationOr, 0, len(x.OrConditions.Entries))
					for _, oobj := range x.OrConditions.Entries {
						olist = append(olist, CombinationOr{
							Name:     oobj.Name,
							ThreatId: oobj.ThreatId,
						})
					}
				}
				list = append(list, CombinationSignature{
					Name:           x.Name,
					CombinationOrs: olist,
				})
			}
			ans.CombinationSignatureType.Signatures = list
		}
	}

	return ans
}

func normalizeEqualToFromEqualTo(v *equalTo) *EqualTo {
	if v == nil {
		return nil
	}

	return &EqualTo{
		Value:      v.Value,
		Context:    v.Context,
		Negate:     util.AsBool(v.Negate),
		Qualifiers: normalizeQualifiers(v.Qualifiers),
	}
}

type entry_v2 struct {
	XMLName       xml.Name         `xml:"entry"`
	Name          string           `xml:"name,attr"`
	ThreatName    string           `xml:"threatname"`
	Comment       string           `xml:"comment,omitempty"`
	Severity      string           `xml:"severity"`
	Direction     string           `xml:"direction"`
	DefaultAction *action          `xml:"default-action"`
	Cves          *util.MemberType `xml:"cve"`
	Bugtraqs      *util.MemberType `xml:"bugtraq"`
	Vendors       *util.MemberType `xml:"vendor"`
	References    *util.MemberType `xml:"reference"`
	Signature     sig_v2           `xml:"signature"`
}

type sig_v2 struct {
	Standard    *standard_v2 `xml:"standard"`
	Combination *combination `xml:"combination"`
}

type standard_v2 struct {
	Entries []standardEntry_v2 `xml:"entry"`
}

type standardEntry_v2 struct {
	Name      string   `xml:"name,attr"`
	Comment   string   `xml:"comment,omitempty"`
	Scope     string   `xml:"scope,omitempty"`
	OrderFree string   `xml:"order-free"`
	Ands      *ands_v2 `xml:"and-condition"`
}

type ands_v2 struct {
	Entries []standardAnd_v2 `xml:"entry"`
}

type standardAnd_v2 struct {
	Name string  `xml:"name,attr"`
	Ors  *ors_v2 `xml:"or-condition"`
}

type ors_v2 struct {
	Entries []standardOr_v2 `xml:"entry"`
}

type standardOr_v2 struct {
	Name string        `xml:"name,attr"`
	Op   standardOp_v2 `xml:"operator"`
}

type standardOp_v2 struct {
	LessThan     *condition    `xml:"less-than"`
	EqualTo      *equalTo      `xml:"equal-to"`
	GreaterThan  *condition    `xml:"greater-than"`
	PatternMatch *patternMatch `xml:"pattern-match"`
}

type equalTo struct {
	Value      int         `xml:"value"`
	Context    string      `xml:"context"`
	Negate     string      `xml:"negate"`
	Qualifiers *qualifiers `xml:"qualifier"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:          e.Name,
		ThreatName:    e.ThreatName,
		Comment:       e.Comment,
		Severity:      e.Severity,
		Direction:     e.Direction,
		DefaultAction: &action{},
		Cves:          util.StrToMem(e.Cves),
		Bugtraqs:      util.StrToMem(e.Bugtraqs),
		Vendors:       util.StrToMem(e.Vendors),
		References:    util.StrToMem(e.References),
	}

	s := ""
	switch e.DefaultAction {
	case DefaultActionAllow:
		ans.DefaultAction.Allow = &s
	case DefaultActionAlert:
		ans.DefaultAction.Alert = &s
	case DefaultActionDrop:
		ans.DefaultAction.Drop = &s
	case DefaultActionResetClient:
		ans.DefaultAction.ResetClient = &s
	case DefaultActionResetServer:
		ans.DefaultAction.ResetServer = &s
	case DefaultActionResetBoth:
		ans.DefaultAction.ResetBoth = &s
	case DefaultActionBlockIp:
		ans.DefaultAction.BlockIp = &blockIp{
			BlockIpTrackBy:  e.BlockIpTrackBy,
			BlockIpDuration: e.BlockIpDuration,
		}
	default:
		ans.DefaultAction = nil
	}

	if e.StandardSignatureType != nil {
		ans.Signature.Standard = &standard_v2{}
		if len(e.StandardSignatureType.Signatures) > 0 {
			sigList := make([]standardEntry_v2, 0, len(e.StandardSignatureType.Signatures))
			for _, ss := range e.StandardSignatureType.Signatures {
				var av *ands_v2
				if len(ss.StandardAnds) > 0 {
					alist := make([]standardAnd_v2, 0, len(ss.StandardAnds))
					for _, aobj := range ss.StandardAnds {
						var ov *ors_v2
						if len(aobj.StandardOrs) > 0 {
							olist := make([]standardOr_v2, 0, len(aobj.StandardOrs))
							for _, oobj := range aobj.StandardOrs {
								olist = append(olist, standardOr_v2{
									Name: oobj.Name,
									Op: standardOp_v2{
										LessThan:     specifyCondition(oobj.LessThan),
										EqualTo:      specifyEqualToFromEqualTo(oobj.EqualTo),
										GreaterThan:  specifyCondition(oobj.GreaterThan),
										PatternMatch: specifyPatternMatch(oobj.PatternMatch),
									},
								})
							}
							ov = &ors_v2{Entries: olist}
						}
						alist = append(alist, standardAnd_v2{
							Name: aobj.Name,
							Ors:  ov,
						})
					}
					av = &ands_v2{Entries: alist}
				}
				sigList = append(sigList, standardEntry_v2{
					Name:      ss.Name,
					Comment:   ss.Comment,
					Scope:     ss.Scope,
					OrderFree: util.YesNo(ss.OrderFree),
					Ands:      av,
				})
			}
			ans.Signature.Standard.Entries = sigList
		}
	}

	if e.CombinationSignatureType != nil {
		ans.Signature.Combination = &combination{
			OrderFree: util.YesNo(e.CombinationSignatureType.OrderFree),
			TimeAttribute: comboTime{
				ThresholdTime:       e.CombinationSignatureType.ThresholdTime,
				IntervalTime:        e.CombinationSignatureType.IntervalTime,
				AggregationCriteria: e.CombinationSignatureType.AggregationCriteria,
			},
		}

		if len(e.CombinationSignatureType.Signatures) > 0 {
			list := make([]andCondEntry, 0, len(e.CombinationSignatureType.Signatures))
			for _, x := range e.CombinationSignatureType.Signatures {
				ors := make([]orCondEntry, 0, len(x.CombinationOrs))
				for _, y := range x.CombinationOrs {
					ors = append(ors, orCondEntry{
						Name:     y.Name,
						ThreatId: y.ThreatId,
					})
				}

				list = append(list, andCondEntry{
					Name:         x.Name,
					OrConditions: &orCond{Entries: ors},
				})
			}

			ans.Signature.Combination.AndConditions = &andCond{
				Entries: list,
			}
		}
	}

	return ans
}

func specifyEqualToFromEqualTo(v *EqualTo) *equalTo {
	if v == nil {
		return nil
	}

	return &equalTo{
		Value:      v.Value,
		Context:    v.Context,
		Negate:     util.YesNo(v.Negate),
		Qualifiers: specifyQualifiers(v.Qualifiers),
	}
}
