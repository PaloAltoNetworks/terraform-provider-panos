package snmp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Entry is a normalized, version independent representation of a snmptrap profile.
//
// PAN-OS 7.1+.
type Entry struct {
	Name        string
	SnmpVersion string

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.SnmpVersion = s.SnmpVersion
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

	if o.Answer.V2c != nil {
		ans.SnmpVersion = SnmpVersionV2c
		if o.Answer.V2c.Servers != nil {
			ans.raw = map[string]string{
				"v2c": util.CleanRawXml(o.Answer.V2c.Servers.Text),
			}
		}
	} else if o.Answer.V3 != nil {
		ans.SnmpVersion = SnmpVersionV3
		if o.Answer.V3.Servers != nil {
			ans.raw = map[string]string{
				"v3": util.CleanRawXml(o.Answer.V3.Servers.Text),
			}
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName xml.Name `xml:"entry"`
	Name    string   `xml:"name,attr"`
	V2c     *details `xml:"version>v2c"`
	V3      *details `xml:"version>v3"`
}

type details struct {
	Servers *util.RawXml `xml:"server"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	switch e.SnmpVersion {
	case SnmpVersionV2c:
		ans.V2c = &details{}
		if text := e.raw["v2c"]; text != "" {
			ans.V2c.Servers = &util.RawXml{text}
		}
	case SnmpVersionV3:
		ans.V3 = &details{}
		if text := e.raw["v3"]; text != "" {
			ans.V3.Servers = &util.RawXml{text}
		}
	}

	return ans
}
