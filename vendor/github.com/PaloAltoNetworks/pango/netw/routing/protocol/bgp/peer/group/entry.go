package group

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a BGP
// peer group.
type Entry struct {
	Name                    string
	Enable                  bool
	AggregatedConfedAsPath  bool
	SoftResetWithStoredInfo bool
	Type                    string
	ExportNextHop           string
	ImportNextHop           string
	RemovePrivateAs         bool

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Enable = s.Enable
	o.AggregatedConfedAsPath = s.AggregatedConfedAsPath
	o.SoftResetWithStoredInfo = s.SoftResetWithStoredInfo
	o.Type = s.Type
	o.ExportNextHop = s.ExportNextHop
	o.ImportNextHop = s.ImportNextHop
	o.RemovePrivateAs = s.RemovePrivateAs
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

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                    o.Name,
		Enable:                  util.AsBool(o.Enable),
		AggregatedConfedAsPath:  util.AsBool(o.AggregatedConfedAsPath),
		SoftResetWithStoredInfo: util.AsBool(o.SoftResetWithStoredInfo),
	}

	if o.Type == nil {
		ans.Type = TypeEbgp
	} else if o.Type.Ebgp != nil {
		ans.Type = TypeEbgp
		ans.ExportNextHop = o.Type.Ebgp.ExportNextHop
		ans.ImportNextHop = o.Type.Ebgp.ImportNextHop
		ans.RemovePrivateAs = util.AsBool(o.Type.Ebgp.RemovePrivateAs)
	} else if o.Type.EbgpConfed != nil {
		ans.Type = TypeEbgpConfed
		ans.ExportNextHop = o.Type.EbgpConfed.ExportNextHop
	} else if o.Type.Ibgp != nil {
		ans.Type = TypeIbgp
		ans.ExportNextHop = o.Type.Ibgp.ExportNextHop
	} else if o.Type.IbgpConfed != nil {
		ans.Type = TypeIbgpConfed
		ans.ExportNextHop = o.Type.IbgpConfed.ExportNextHop
	}

	if o.Peer != nil {
		ans.raw = map[string]string{
			"peer": util.CleanRawXml(o.Peer.Text),
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName                 xml.Name     `xml:"entry"`
	Name                    string       `xml:"name,attr"`
	Enable                  string       `xml:"enable"`
	AggregatedConfedAsPath  string       `xml:"aggregated-confed-as-path"`
	SoftResetWithStoredInfo string       `xml:"soft-reset-with-stored-info"`
	Type                    *gType       `xml:"type"`
	Peer                    *util.RawXml `xml:"peer"`
}

type gType struct {
	Ebgp       *ebgpOptions  `xml:"ebgp"`
	EbgpConfed *basicOptions `xml:"ebgp-confed"`
	Ibgp       *basicOptions `xml:"ibgp"`
	IbgpConfed *basicOptions `xml:"ibgp-confed"`
}

type ebgpOptions struct {
	ExportNextHop   string `xml:"export-nexthop,omitempty"`
	ImportNextHop   string `xml:"import-nexthop,omitempty"`
	RemovePrivateAs string `xml:"remove-private-as"`
}

type basicOptions struct {
	ExportNextHop string `xml:"export-nexthop,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                    e.Name,
		Enable:                  util.YesNo(e.Enable),
		AggregatedConfedAsPath:  util.YesNo(e.AggregatedConfedAsPath),
		SoftResetWithStoredInfo: util.YesNo(e.SoftResetWithStoredInfo),
	}

	switch e.Type {
	case TypeEbgp:
		ans.Type = &gType{
			Ebgp: &ebgpOptions{
				ExportNextHop:   e.ExportNextHop,
				ImportNextHop:   e.ImportNextHop,
				RemovePrivateAs: util.YesNo(e.RemovePrivateAs),
			},
		}
	case TypeEbgpConfed:
		ans.Type = &gType{
			EbgpConfed: &basicOptions{
				ExportNextHop: e.ExportNextHop,
			},
		}
	case TypeIbgp:
		ans.Type = &gType{
			Ibgp: &basicOptions{
				ExportNextHop: e.ExportNextHop,
			},
		}
	case TypeIbgpConfed:
		ans.Type = &gType{
			IbgpConfed: &basicOptions{
				ExportNextHop: e.ExportNextHop,
			},
		}
	}

	if text, present := e.raw["peer"]; present {
		ans.Peer = &util.RawXml{text}
	}

	return ans
}
