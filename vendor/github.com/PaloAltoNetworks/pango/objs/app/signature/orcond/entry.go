package orcond

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of an application signature and-condition.
type Entry struct {
	Name       string
	Operator   string
	Context    string
	Pattern    string
	Value      string
	Position   string
	Mask       string
	Qualifiers map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Operator = s.Operator
	o.Context = s.Context
	o.Pattern = s.Pattern
	o.Value = s.Value
	o.Position = s.Position
	o.Mask = s.Mask
	o.Qualifiers = s.Qualifiers
}

/** Structs / functions for this namespace. **/

type normalizer interface {
	Normalize() Entry
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>entry"`
}

func (o *container_v1) Normalize() Entry {
	ans := Entry{
		Name: o.Answer.Name,
	}

	var q *qual
	if o.Answer.Operator.Pm != nil {
		ans.Operator = OperatorPatternMatch
		ans.Context = o.Answer.Operator.Pm.Context
		ans.Pattern = o.Answer.Operator.Pm.Pattern
		q = o.Answer.Operator.Pm.Qualifier
	} else if o.Answer.Operator.Gt != nil {
		ans.Operator = OperatorGreaterThan
		ans.Context = o.Answer.Operator.Gt.Context
		ans.Value = o.Answer.Operator.Gt.Value
		q = o.Answer.Operator.Gt.Qualifier
	} else if o.Answer.Operator.Lt != nil {
		ans.Operator = OperatorLessThan
		ans.Context = o.Answer.Operator.Lt.Context
		ans.Value = o.Answer.Operator.Lt.Value
		q = o.Answer.Operator.Lt.Qualifier
	} else if o.Answer.Operator.Eq != nil {
		ans.Operator = OperatorEqualTo
		ans.Context = o.Answer.Operator.Eq.Context
		ans.Position = o.Answer.Operator.Eq.Position
		ans.Mask = o.Answer.Operator.Eq.Mask
		ans.Value = o.Answer.Operator.Eq.Value
	}

	if q != nil {
		ans.Qualifiers = make(map[string]string)
		for i := range q.Entry {
			ans.Qualifiers[q.Entry[i].Name] = q.Entry[i].Value
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName  xml.Name `xml:"entry"`
	Name     string   `xml:"name,attr"`
	Operator op       `xml:"operator"`
}

type op struct {
	Pm *pm `xml:"pattern-match"`
	Gt *gt `xml:"greater-than"`
	Lt *lt `xml:"less-than"`
	Eq *eq `xml:"equal-to"`
}

type pm struct {
	Context   string `xml:"context"`
	Pattern   string `xml:"pattern"`
	Qualifier *qual  `xml:"qualifier"`
}

type qual struct {
	Entry []qe `xml:"entry"`
}

type qe struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value"`
}

type gt struct {
	Context   string `xml:"context"`
	Value     string `xml:"value"`
	Qualifier *qual  `xml:"qualifier"`
}

type lt struct {
	Context   string `xml:"context"`
	Value     string `xml:"value"`
	Qualifier *qual  `xml:"qualifier"`
}

type eq struct {
	Context  string `xml:"context"`
	Position string `xml:"position,omitempty"`
	Mask     string `xml:"mask,omitempty"`
	Value    string `xml:"value"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	var q *qual
	if len(e.Qualifiers) > 0 {
		q = &qual{Entry: make([]qe, 0, len(e.Qualifiers))}
		for k, v := range e.Qualifiers {
			q.Entry = append(q.Entry, qe{k, v})
		}
	}

	switch e.Operator {
	case OperatorPatternMatch:
		ans.Operator.Pm = &pm{
			Context:   e.Context,
			Pattern:   e.Pattern,
			Qualifier: q,
		}
	case OperatorGreaterThan:
		ans.Operator.Gt = &gt{
			Context:   e.Context,
			Value:     e.Value,
			Qualifier: q,
		}
	case OperatorLessThan:
		ans.Operator.Lt = &lt{
			Context:   e.Context,
			Value:     e.Value,
			Qualifier: q,
		}
	case OperatorEqualTo:
		ans.Operator.Eq = &eq{
			Context:  e.Context,
			Position: e.Position,
			Mask:     e.Mask,
			Value:    e.Value,
		}
	}

	return ans
}
