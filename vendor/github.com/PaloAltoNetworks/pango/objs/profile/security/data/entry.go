package data

import (
	"encoding/xml"
	"fmt"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// data filtering profile.
type Entry struct {
	Name        string
	Description string
	DataCapture bool
	Rules       []Rule
}

// Rule is a data filtering rule.  In PAN-OS, the Name param is
// auto generated and look like "rule1", "rule2"...  If the Name param
// is an empty string, the name will be auto populated as appropriate.
type Rule struct {
	Name           string
	DataPattern    string
	Applications   []string // ordered
	FileTypes      []string // ordered
	Direction      string
	AlertThreshold int
	BlockThreshold int
	LogSeverity    string // 8.0+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.DataCapture = s.DataCapture
	if s.Rules == nil {
		o.Rules = nil
	} else {
		o.Rules = make([]Rule, 0, len(s.Rules))
		for _, x := range s.Rules {
			r := Rule{
				Name:           x.Name,
				DataPattern:    x.DataPattern,
				Direction:      x.Direction,
				AlertThreshold: x.AlertThreshold,
				BlockThreshold: x.BlockThreshold,
				LogSeverity:    x.LogSeverity,
			}
			if len(x.Applications) != 0 {
				r.Applications = make([]string, len(x.Applications))
				copy(r.Applications, x.Applications)
			}
			if len(x.FileTypes) != 0 {
				r.FileTypes = make([]string, len(x.FileTypes))
				copy(r.FileTypes, x.FileTypes)
			}
			o.Rules = append(o.Rules, r)
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
		DataCapture: util.AsBool(o.DataCapture),
	}

	if o.Rules != nil {
		rules := make([]Rule, 0, len(o.Rules.Entries))
		for _, v := range o.Rules.Entries {
			r := Rule{
				Name:           v.Name,
				DataPattern:    v.DataPattern,
				Applications:   util.MemToStr(v.Applications),
				FileTypes:      util.MemToStr(v.FileTypes),
				Direction:      v.Direction,
				AlertThreshold: v.AlertThreshold,
				BlockThreshold: v.BlockThreshold,
			}
			rules = append(rules, r)
		}
		ans.Rules = rules
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name  `xml:"entry"`
	Name        string    `xml:"name,attr"`
	Description string    `xml:"description,omitempty"`
	DataCapture string    `xml:"data-capture,omitempty"`
	Rules       *rules_v1 `xml:"rules"`
}

type rules_v1 struct {
	Entries []rule_v1 `xml:"entry"`
}

type rule_v1 struct {
	Name           string           `xml:"name,attr"`
	DataPattern    string           `xml:"data-object"`
	Applications   *util.MemberType `xml:"application"`
	FileTypes      *util.MemberType `xml:"file-type"`
	Direction      string           `xml:"direction,omitempty"`
	AlertThreshold int              `xml:"alert-threshold,omitempty"`
	BlockThreshold int              `xml:"block-threshold,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		DataCapture: util.YesNo(e.DataCapture),
	}

	if len(e.Rules) > 0 {
		rules := make([]rule_v1, 0, len(e.Rules))
		for num, er := range e.Rules {
			r := rule_v1{
				Name:           er.Name,
				DataPattern:    er.DataPattern,
				Applications:   util.StrToMem(er.Applications),
				FileTypes:      util.StrToMem(er.FileTypes),
				Direction:      er.Direction,
				AlertThreshold: er.AlertThreshold,
				BlockThreshold: er.BlockThreshold,
			}
			if er.Name == "" {
				r.Name = fmt.Sprintf("rule%d", num)
			}
			rules = append(rules, r)
		}
		ans.Rules = &rules_v1{Entries: rules}
	}

	return ans
}

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
		DataCapture: util.AsBool(o.DataCapture),
	}

	if o.Rules != nil {
		rules := make([]Rule, 0, len(o.Rules.Entries))
		for _, v := range o.Rules.Entries {
			r := Rule{
				Name:           v.Name,
				DataPattern:    v.DataPattern,
				Applications:   util.MemToStr(v.Applications),
				FileTypes:      util.MemToStr(v.FileTypes),
				Direction:      v.Direction,
				AlertThreshold: v.AlertThreshold,
				BlockThreshold: v.BlockThreshold,
				LogSeverity:    v.LogSeverity,
			}
			rules = append(rules, r)
		}
		ans.Rules = rules
	}

	return ans
}

type entry_v2 struct {
	XMLName     xml.Name  `xml:"entry"`
	Name        string    `xml:"name,attr"`
	Description string    `xml:"description,omitempty"`
	DataCapture string    `xml:"data-capture,omitempty"`
	Rules       *rules_v2 `xml:"rules"`
}

type rules_v2 struct {
	Entries []rule_v2 `xml:"entry"`
}

type rule_v2 struct {
	Name           string           `xml:"name,attr"`
	DataPattern    string           `xml:"data-object"`
	Applications   *util.MemberType `xml:"application"`
	FileTypes      *util.MemberType `xml:"file-type"`
	Direction      string           `xml:"direction,omitempty"`
	AlertThreshold int              `xml:"alert-threshold,omitempty"`
	BlockThreshold int              `xml:"block-threshold,omitempty"`
	LogSeverity    string           `xml:"log-severity,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:        e.Name,
		Description: e.Description,
		DataCapture: util.YesNo(e.DataCapture),
	}

	if len(e.Rules) > 0 {
		rules := make([]rule_v2, 0, len(e.Rules))
		for num, er := range e.Rules {
			r := rule_v2{
				Name:           er.Name,
				DataPattern:    er.DataPattern,
				Applications:   util.StrToMem(er.Applications),
				FileTypes:      util.StrToMem(er.FileTypes),
				Direction:      er.Direction,
				AlertThreshold: er.AlertThreshold,
				BlockThreshold: er.BlockThreshold,
				LogSeverity:    er.LogSeverity,
			}
			if er.Name == "" {
				r.Name = fmt.Sprintf("rule%d", num)
			}
			rules = append(rules, r)
		}
		ans.Rules = &rules_v2{Entries: rules}
	}

	return ans
}
