package logfwd

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a log forwarding profile.
//
// PAN-OS 8.0+.
type Entry struct {
	Name            string
	Description     string
	EnhancedLogging bool // 8.1+

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.EnhancedLogging = s.EnhancedLogging
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

	if o.MatchList != nil {
		ans.raw = map[string]string{
			"ml": util.CleanRawXml(o.MatchList.Text),
		}
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
		Name:            o.Name,
		Description:     o.Description,
		EnhancedLogging: util.AsBool(o.EnhancedLogging),
	}

	if o.MatchList != nil {
		ans.raw = map[string]string{
			"ml": util.CleanRawXml(o.MatchList.Text),
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name     `xml:"entry"`
	Name        string       `xml:"name,attr"`
	Description string       `xml:"description,omitempty"`
	MatchList   *util.RawXml `xml:"match-list"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
	}

	if text := e.raw["ml"]; text != "" {
		ans.MatchList = &util.RawXml{text}
	}

	return ans
}

type entry_v2 struct {
	XMLName         xml.Name     `xml:"entry"`
	Name            string       `xml:"name,attr"`
	Description     string       `xml:"description,omitempty"`
	EnhancedLogging string       `xml:"enhanced-application-logging"`
	MatchList       *util.RawXml `xml:"match-list"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:            e.Name,
		Description:     e.Description,
		EnhancedLogging: util.YesNo(e.EnhancedLogging),
	}

	if text := e.raw["ml"]; text != "" {
		ans.MatchList = &util.RawXml{text}
	}

	return ans
}
