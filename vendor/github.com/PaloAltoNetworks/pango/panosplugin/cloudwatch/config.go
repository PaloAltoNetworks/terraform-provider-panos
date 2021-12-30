package cloudwatch

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/plugin"
	"github.com/PaloAltoNetworks/pango/util"
)

// Config is a normalized, version independent representation of an
// AWS CloudWatch config.
//
// PAN-OS 9.0+
type Config struct {
	Enabled        bool
	Namespace      string
	UpdateInterval int
}

// Copy copies the information from source Config `s` to this object.
func (o *Config) Copy(s Config) {
	o.Enabled = s.Enabled
	o.Namespace = s.Namespace
	o.UpdateInterval = s.UpdateInterval
}

/** Structs / functions for this namespace. **/

func (o Config) Specify(list []plugin.Info) (string, interface{}, error) {
	_, fn, err := versioning(list)
	if err != nil {
		return "", nil, err
	}

	return "", fn(o), nil
}

type normalizer interface {
	Names() []string
	Normalize() []Config
}

type container_v1 struct {
	Answer []config_v1 `xml:"aws-cloudwatch"`
}

func (o *container_v1) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	return nil
}

type config_v1 struct {
	XMLName        xml.Name `xml:"aws-cloudwatch"`
	Enabled        string   `xml:"enabled"`
	Namespace      string   `xml:"name,omitempty"`
	UpdateInterval int      `xml:"timeout,omitempty"`
}

func (e *config_v1) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type localConfig_v1 config_v1
	ans := localConfig_v1{
		Namespace:      "VMseries",
		UpdateInterval: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = config_v1(ans)
	return nil
}

func (o *config_v1) normalize() Config {
	ans := Config{
		Enabled:        util.AsBool(o.Enabled),
		Namespace:      o.Namespace,
		UpdateInterval: o.UpdateInterval,
	}

	return ans
}

func specify_v1(e Config) interface{} {
	ans := config_v1{
		Enabled:        util.YesNo(e.Enabled),
		Namespace:      e.Namespace,
		UpdateInterval: e.UpdateInterval,
	}

	return ans
}
