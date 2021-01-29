package matchlist

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a log forwarding profile match list.
//
// PAN-OS 8.0+.
type Entry struct {
	Name           string
	Description    string
	LogType        string
	Filter         string
	SendToPanorama bool
	SnmpProfiles   []string // unordered
	EmailProfiles  []string // unordered
	SyslogProfiles []string // unordered
	HttpProfiles   []string // unordered

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.LogType = s.LogType
	o.Filter = s.Filter
	o.SendToPanorama = s.SendToPanorama
	o.SnmpProfiles = s.SnmpProfiles
	o.EmailProfiles = s.EmailProfiles
	o.SyslogProfiles = s.SyslogProfiles
	o.HttpProfiles = s.HttpProfiles
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
		Name:           o.Answer.Name,
		Description:    o.Answer.Description,
		LogType:        o.Answer.LogType,
		Filter:         o.Answer.Filter,
		SendToPanorama: util.AsBool(o.Answer.SendToPanorama),
		SnmpProfiles:   util.MemToStr(o.Answer.SnmpProfiles),
		EmailProfiles:  util.MemToStr(o.Answer.EmailProfiles),
		SyslogProfiles: util.MemToStr(o.Answer.SyslogProfiles),
		HttpProfiles:   util.MemToStr(o.Answer.HttpProfiles),
	}

	if o.Answer.Actions != nil {
		ans.raw = map[string]string{
			"act": util.CleanRawXml(o.Answer.Actions.Text),
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName        xml.Name         `xml:"entry"`
	Name           string           `xml:"name,attr"`
	Description    string           `xml:"action-desc,omitempty"`
	LogType        string           `xml:"log-type"`
	Filter         string           `xml:"filter"`
	SendToPanorama string           `xml:"send-to-panorama"`
	SnmpProfiles   *util.MemberType `xml:"send-snmptrap"`
	EmailProfiles  *util.MemberType `xml:"send-email"`
	SyslogProfiles *util.MemberType `xml:"send-syslog"`
	HttpProfiles   *util.MemberType `xml:"send-http"`
	Actions        *util.RawXml     `xml:"actions"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:           e.Name,
		Description:    e.Description,
		LogType:        e.LogType,
		Filter:         e.Filter,
		SendToPanorama: util.YesNo(e.SendToPanorama),
		SnmpProfiles:   util.StrToMem(e.SnmpProfiles),
		EmailProfiles:  util.StrToMem(e.EmailProfiles),
		SyslogProfiles: util.StrToMem(e.SyslogProfiles),
		HttpProfiles:   util.StrToMem(e.HttpProfiles),
	}

	if text := e.raw["act"]; text != "" {
		ans.Actions = &util.RawXml{text}
	}

	return ans
}
