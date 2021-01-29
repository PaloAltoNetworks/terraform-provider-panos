package arp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an arp entry.
type Entry struct {
	Ip         string
	MacAddress string
	Interface  string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Ip field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.MacAddress = s.MacAddress
	o.Interface = s.Interface
}

/** Structs / functions for normalization. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Ip, fn(o)
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
		ans = append(ans, o.Answer[i].Ip)
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
		Ip:         o.Ip,
		MacAddress: o.MacAddress,
		Interface:  o.Interface,
	}

	return ans
}

type entry_v1 struct {
	XMLName    xml.Name `xml:"entry"`
	Ip         string   `xml:"name,attr"`
	MacAddress string   `xml:"hw-address"`
	Interface  string   `xml:"interface,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Ip:         e.Ip,
		MacAddress: e.MacAddress,
		Interface:  e.Interface,
	}

	return ans
}
