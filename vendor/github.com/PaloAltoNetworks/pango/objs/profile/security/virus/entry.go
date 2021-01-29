package virus

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// antivirus security profile.
type Entry struct {
	Name                      string
	Description               string
	PacketCapture             bool
	Decoders                  []Decoder
	ApplicationExceptions     []ApplicationException
	ThreatExceptions          []string
	MachineLearningModels     []MachineLearningModel     // 10.0
	MachineLearningExceptions []MachineLearningException // 10.0
}

type Decoder struct {
	Name                  string
	Action                string
	WildfireAction        string
	MachineLearningAction string // 10.0
}

type ApplicationException struct {
	Application string
	Action      string
}

type MachineLearningModel struct {
	Model  string
	Action string
}

type MachineLearningException struct {
	Name        string
	Description string
	Filename    string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.PacketCapture = s.PacketCapture
	if s.Decoders == nil {
		o.Decoders = nil
	} else {
		o.Decoders = make([]Decoder, 0, len(s.Decoders))
		for _, x := range s.Decoders {
			o.Decoders = append(o.Decoders, Decoder{
				Name:                  x.Name,
				Action:                x.Action,
				WildfireAction:        x.WildfireAction,
				MachineLearningAction: x.MachineLearningAction,
			})
		}
	}
	if s.ApplicationExceptions == nil {
		o.ApplicationExceptions = nil
	} else {
		o.ApplicationExceptions = make([]ApplicationException, 0, len(s.ApplicationExceptions))
		for _, x := range s.ApplicationExceptions {
			o.ApplicationExceptions = append(o.ApplicationExceptions, ApplicationException{
				Application: x.Application,
				Action:      x.Action,
			})
		}
	}
	o.ThreatExceptions = s.ThreatExceptions
	if s.MachineLearningModels == nil {
		o.MachineLearningModels = nil
	} else {
		o.MachineLearningModels = make([]MachineLearningModel, 0, len(s.MachineLearningModels))
		for _, x := range s.MachineLearningModels {
			o.MachineLearningModels = append(o.MachineLearningModels, MachineLearningModel{
				Model:  x.Model,
				Action: x.Action,
			})
		}
	}
	if s.MachineLearningExceptions == nil {
		o.MachineLearningExceptions = nil
	} else {
		o.MachineLearningExceptions = make([]MachineLearningException, 0, len(s.MachineLearningExceptions))
		for _, x := range s.MachineLearningExceptions {
			o.MachineLearningExceptions = append(o.MachineLearningExceptions, MachineLearningException{
				Name:        x.Name,
				Description: x.Description,
				Filename:    x.Filename,
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
		Name:             o.Name,
		Description:      o.Description,
		PacketCapture:    util.AsBool(o.PacketCapture),
		ThreatExceptions: util.EntToStr(o.ThreatExceptions),
	}

	if o.Decoder != nil {
		data := make([]Decoder, 0, len(o.Decoder.Entries))
		for _, d := range o.Decoder.Entries {
			data = append(data, Decoder{
				Name:           d.Name,
				Action:         d.Action,
				WildfireAction: d.WildfireAction,
			})
		}

		ans.Decoders = data
	}

	if o.Application != nil {
		data := make([]ApplicationException, 0, len(o.Application.Entries))
		for _, d := range o.Application.Entries {
			data = append(data, ApplicationException{
				Application: d.Application,
				Action:      d.Action,
			})
		}

		ans.ApplicationExceptions = data
	}

	return ans
}

type entry_v1 struct {
	XMLName          xml.Name        `xml:"entry"`
	Name             string          `xml:"name,attr"`
	Description      string          `xml:"description,omitempty"`
	PacketCapture    string          `xml:"packet-capture"`
	Decoder          *decoder_v1     `xml:"decoder"`
	Application      *application    `xml:"application"`
	ThreatExceptions *util.EntryType `xml:"threat-exception"`
}

type decoder_v1 struct {
	Entries []decoderEntry_v1 `xml:"entry"`
}

type decoderEntry_v1 struct {
	Name           string `xml:"name,attr"`
	Action         string `xml:"action,omitempty"`
	WildfireAction string `xml:"wildfire-action,omitempty"`
}

type application struct {
	Entries []applicationEntry `xml:"entry"`
}

type applicationEntry struct {
	Application string `xml:"name,attr"`
	Action      string `xml:"action,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:             e.Name,
		Description:      e.Description,
		PacketCapture:    util.YesNo(e.PacketCapture),
		ThreatExceptions: util.StrToEnt(e.ThreatExceptions),
	}

	if len(e.Decoders) > 0 {
		data := make([]decoderEntry_v1, 0, len(e.Decoders))
		for _, d := range e.Decoders {
			data = append(data, decoderEntry_v1{
				Name:           d.Name,
				Action:         d.Action,
				WildfireAction: d.WildfireAction,
			})
		}

		ans.Decoder = &decoder_v1{Entries: data}
	}

	if len(e.ApplicationExceptions) > 0 {
		data := make([]applicationEntry, 0, len(e.ApplicationExceptions))
		for _, d := range e.ApplicationExceptions {
			data = append(data, applicationEntry{
				Application: d.Application,
				Action:      d.Action,
			})
		}

		ans.Application = &application{Entries: data}
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
		Name:             o.Name,
		Description:      o.Description,
		PacketCapture:    util.AsBool(o.PacketCapture),
		ThreatExceptions: util.EntToStr(o.ThreatExceptions),
	}

	if o.Decoder != nil {
		data := make([]Decoder, 0, len(o.Decoder.Entries))
		for _, d := range o.Decoder.Entries {
			data = append(data, Decoder{
				Name:                  d.Name,
				Action:                d.Action,
				WildfireAction:        d.WildfireAction,
				MachineLearningAction: d.MachineLearningAction,
			})
		}

		ans.Decoders = data
	}

	if o.Application != nil {
		data := make([]ApplicationException, 0, len(o.Application.Entries))
		for _, d := range o.Application.Entries {
			data = append(data, ApplicationException{
				Application: d.Application,
				Action:      d.Action,
			})
		}

		ans.ApplicationExceptions = data
	}

	if o.Ml != nil {
		data := make([]MachineLearningModel, 0, len(o.Ml.Entries))
		for _, d := range o.Ml.Entries {
			data = append(data, MachineLearningModel{
				Model:  d.Model,
				Action: d.Action,
			})
		}

		ans.MachineLearningModels = data
	}

	if o.MlException != nil {
		data := make([]MachineLearningException, 0, len(o.MlException.Entries))
		for _, d := range o.MlException.Entries {
			data = append(data, MachineLearningException{
				Name:        d.Name,
				Description: d.Description,
				Filename:    d.Filename,
			})
		}

		ans.MachineLearningExceptions = data
	}

	return ans
}

type entry_v2 struct {
	XMLName          xml.Name        `xml:"entry"`
	Name             string          `xml:"name,attr"`
	Description      string          `xml:"description,omitempty"`
	PacketCapture    string          `xml:"packet-capture"`
	Ml               *mlConfig       `xml:"mlav-engine-filebased-enabled"`
	Decoder          *decoder_v2     `xml:"decoder"`
	Application      *application    `xml:"application"`
	ThreatExceptions *util.EntryType `xml:"threat-exception"`
	MlException      *mleConfig      `xml:"mlav-exception"`
}

type mlConfig struct {
	Entries []mlEntry `xml:"entry"`
}

type mlEntry struct {
	Model  string `xml:"name,attr"`
	Action string `xml:"mlav-policy-action"`
}

type decoder_v2 struct {
	Entries []decoderEntry_v2 `xml:"entry"`
}

type decoderEntry_v2 struct {
	Name                  string `xml:"name,attr"`
	Action                string `xml:"action,omitempty"`
	WildfireAction        string `xml:"wildfire-action,omitempty"`
	MachineLearningAction string `xml:"mlav-action,omitempty"`
}

type mleConfig struct {
	Entries []mlExceptionEntry `xml:"entry"`
}

type mlExceptionEntry struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,omitempty"`
	Filename    string `xml:"filename,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:             e.Name,
		Description:      e.Description,
		PacketCapture:    util.YesNo(e.PacketCapture),
		ThreatExceptions: util.StrToEnt(e.ThreatExceptions),
	}

	if len(e.Decoders) > 0 {
		data := make([]decoderEntry_v2, 0, len(e.Decoders))
		for _, d := range e.Decoders {
			data = append(data, decoderEntry_v2{
				Name:                  d.Name,
				Action:                d.Action,
				WildfireAction:        d.WildfireAction,
				MachineLearningAction: d.MachineLearningAction,
			})
		}

		ans.Decoder = &decoder_v2{Entries: data}
	}

	if len(e.ApplicationExceptions) > 0 {
		data := make([]applicationEntry, 0, len(e.ApplicationExceptions))
		for _, d := range e.ApplicationExceptions {
			data = append(data, applicationEntry{
				Application: d.Application,
				Action:      d.Action,
			})
		}

		ans.Application = &application{Entries: data}
	}

	if len(e.MachineLearningModels) > 0 {
		data := make([]mlEntry, 0, len(e.MachineLearningModels))
		for _, d := range e.MachineLearningModels {
			data = append(data, mlEntry{
				Model:  d.Model,
				Action: d.Action,
			})
		}

		ans.Ml = &mlConfig{Entries: data}
	}

	if len(e.MachineLearningExceptions) > 0 {
		data := make([]mlExceptionEntry, 0, len(e.MachineLearningExceptions))
		for _, d := range e.MachineLearningExceptions {
			data = append(data, mlExceptionEntry{
				Name:        d.Name,
				Description: d.Description,
				Filename:    d.Filename,
			})
		}

		ans.MlException = &mleConfig{Entries: data}
	}

	return ans
}
