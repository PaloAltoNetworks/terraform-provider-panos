package address

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an
// IPv6 address.
//
// Note that loopback and tunnel interfaces do not have neighbor discovery config,
// so all router advertisement params should be left as empty for those types.
type Entry struct {
	Name                string
	Enabled             bool
	InterfaceIdAsHost   bool
	Anycast             bool
	EnableRa            bool
	RaValidLifetime     string
	RaPreferredLifetime string
	RaOnLink            bool
	RaAutonomous        bool
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Enabled = s.Enabled
	o.InterfaceIdAsHost = s.InterfaceIdAsHost
	o.Anycast = s.Anycast
	o.EnableRa = s.EnableRa
	o.RaValidLifetime = s.RaValidLifetime
	o.RaPreferredLifetime = s.RaPreferredLifetime
	o.RaOnLink = s.RaOnLink
	o.RaAutonomous = s.RaAutonomous
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
		Name:    o.Name,
		Enabled: util.AsBool(o.Enabled),
	}

	if o.InterfaceIdAsHost != nil {
		ans.InterfaceIdAsHost = true
	}

	if o.Anycast != nil {
		ans.Anycast = true
	}

	if o.Ra != nil {
		ans.EnableRa = util.AsBool(o.Ra.EnableRa)
		ans.RaValidLifetime = o.Ra.RaValidLifetime
		ans.RaPreferredLifetime = o.Ra.RaPreferredLifetime
		ans.RaOnLink = util.AsBool(o.Ra.RaOnLink)
		ans.RaAutonomous = util.AsBool(o.Ra.RaAutonomous)
	}

	return ans
}

type entry_v1 struct {
	XMLName           xml.Name `xml:"entry"`
	Name              string   `xml:"name,attr"`
	Enabled           string   `xml:"enable-on-interface"`
	InterfaceIdAsHost *string  `xml:"prefix"`
	Anycast           *string  `xml:"anycast"`
	Ra                *ra      `xml:"advertise"`
}

type ra struct {
	EnableRa            string `xml:"enable"`
	RaValidLifetime     string `xml:"valid-lifetime,omitempty"`
	RaPreferredLifetime string `xml:"preferred-lifetime,omitempty"`
	RaOnLink            string `xml:"onlink-flag"`
	RaAutonomous        string `xml:"auto-config-flag"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:    e.Name,
		Enabled: util.YesNo(e.Enabled),
	}

	s := ""

	if e.InterfaceIdAsHost {
		ans.InterfaceIdAsHost = &s
	}

	if e.Anycast {
		ans.Anycast = &s
	}

	if e.EnableRa || e.RaValidLifetime != "" || e.RaPreferredLifetime != "" || e.RaOnLink || e.RaAutonomous {
		ans.Ra = &ra{
			EnableRa:            util.YesNo(e.EnableRa),
			RaValidLifetime:     e.RaValidLifetime,
			RaPreferredLifetime: e.RaPreferredLifetime,
			RaOnLink:            util.YesNo(e.RaOnLink),
			RaAutonomous:        util.YesNo(e.RaAutonomous),
		}
	}

	return ans
}
