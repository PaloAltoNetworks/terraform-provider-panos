package util

import (
	"encoding/xml"
)

// NewHitCountRequest returns a new hit count request struct.
//
// If the rules param is nil, then the hit count for all rules is returned.
func NewHitCountRequest(rulebase, vsys string, rules []string) interface{} {
	req := hcReq{
		Vsys: hcReqVsys{
			Name: vsys,
			Rulebase: hcReqRulebase{
				Name: rulebase,
				Rules: hcReqRules{
					List: StrToMem(rules),
				},
			},
		},
	}

	if req.Vsys.Rulebase.Rules.List == nil {
		s := ""
		req.Vsys.Rulebase.Rules.All = &s
	}

	return req
}

// HitCountResponse is the hit count response struct.
type HitCountResponse struct {
	XMLName xml.Name   `xml:"response"`
	Results []HitCount `xml:"result>rule-hit-count>vsys>entry>rule-base>entry>rules>entry"`
}

// HitCount is the hit count data for a specific rule.
type HitCount struct {
	Name                      string `xml:"name,attr"`
	Latest                    string `xml:"latest"`
	HitCount                  uint   `xml:"hit-count"`
	LastHitTimestamp          int    `xml:"last-hit-timestamp"`
	LastResetTimestamp        int    `xml:"last-reset-timestamp"`
	FirstHitTimestamp         int    `xml:"first-hit-timestamp"`
	RuleCreationTimestamp     int    `xml:"rule-creation-timestamp"`
	RuleModificationTimestamp int    `xml:"rule-modification-timestamp"`
}

type hcReq struct {
	XMLName xml.Name  `xml:"show"`
	Vsys    hcReqVsys `xml:"rule-hit-count>vsys>vsys-name>entry"`
}

type hcReqVsys struct {
	Name     string        `xml:"name,attr"`
	Rulebase hcReqRulebase `xml:"rule-base>entry"`
}

type hcReqRulebase struct {
	Name  string     `xml:"name,attr"`
	Rules hcReqRules `xml:"rules"`
}

type hcReqRules struct {
	All  *string     `xml:"all"`
	List *MemberType `xml:"list"`
}
