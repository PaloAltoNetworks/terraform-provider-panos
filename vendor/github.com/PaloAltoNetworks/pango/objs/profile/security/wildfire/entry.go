package wildfire

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// Wildfire analysis security profile.
//
// PAN-OS 7.0+
type Entry struct {
	Name        string
	Description string
	Rules       []Rule
}

type Rule struct {
	Name         string
	Applications []string // ordered
	FileTypes    []string // ordered
	Direction    string
	Analysis     string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	if s.Rules == nil {
		o.Rules = nil
	} else {
		o.Rules = make([]Rule, 0, len(s.Rules))
		for _, x := range s.Rules {
			item := Rule{
				Name:      x.Name,
				Direction: x.Direction,
				Analysis:  x.Analysis,
			}
			if x.Applications != nil {
				item.Applications = make([]string, len(x.Applications))
				copy(item.Applications, x.Applications)
			}
			if x.FileTypes != nil {
				item.FileTypes = make([]string, len(x.FileTypes))
				copy(item.FileTypes, x.FileTypes)
			}
			o.Rules = append(o.Rules, item)
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
	}

	if o.Rules != nil {
		rules := make([]Rule, 0, len(o.Rules.Entries))
		for _, v := range o.Rules.Entries {
			r := Rule{
				Name:         v.Name,
				Applications: util.MemToStr(v.Applications),
				FileTypes:    util.MemToStr(v.FileTypes),
				Direction:    v.Direction,
				Analysis:     v.Analysis,
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
	Rules       *rules_v1 `xml:"rules"`
}

type rules_v1 struct {
	Entries []rule_v1 `xml:"entry"`
}

type rule_v1 struct {
	Name         string           `xml:"name,attr"`
	Applications *util.MemberType `xml:"application"`
	FileTypes    *util.MemberType `xml:"file-type"`
	Direction    string           `xml:"direction,omitempty"`
	Analysis     string           `xml:"analysis,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
	}

	if len(e.Rules) > 0 {
		rules := make([]rule_v1, 0, len(e.Rules))
		for _, er := range e.Rules {
			r := rule_v1{
				Name:         er.Name,
				Applications: util.StrToMem(er.Applications),
				FileTypes:    util.StrToMem(er.FileTypes),
				Direction:    er.Direction,
				Analysis:     er.Analysis,
			}
			rules = append(rules, r)
		}
		ans.Rules = &rules_v1{Entries: rules}
	}

	return ans
}
