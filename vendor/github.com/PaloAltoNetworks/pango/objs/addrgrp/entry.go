package addrgrp

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an address
// group.  The value set in DynamicMatch should be something like the following:
//
//  * 'tag1'
//  * 'tag1' or 'tag2' and 'tag3'
//
// The Tags param is for administrative tags for this address object
// group itself.
type Entry struct {
	Name            string
	Description     string
	StaticAddresses []string // unordered
	DynamicMatch    string
	Tags            []string // ordered
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	if s.StaticAddresses == nil {
		o.StaticAddresses = nil
	} else {
		o.StaticAddresses = make([]string, len(s.StaticAddresses))
		copy(o.StaticAddresses, s.StaticAddresses)
	}
	o.DynamicMatch = s.DynamicMatch
	if s.Tags == nil {
		o.Tags = nil
	} else {
		o.Tags = make([]string, len(s.Tags))
		copy(o.Tags, s.Tags)
	}
}

/** Structs / functions for normalization. **/

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
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:            o.Name,
		Description:     o.Description,
		StaticAddresses: util.MemToStr(o.StaticAddresses),
		Tags:            util.MemToStr(o.Tags),
	}
	if o.DynamicMatch != nil {
		ans.DynamicMatch = *o.DynamicMatch
	}

	return ans
}

type entry_v1 struct {
	XMLName         xml.Name         `xml:"entry"`
	Name            string           `xml:"name,attr"`
	Description     string           `xml:"description,omitempty"`
	StaticAddresses *util.MemberType `xml:"static"`
	DynamicMatch    *string          `xml:"dynamic>filter"`
	Tags            *util.MemberType `xml:"tag"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:            e.Name,
		Description:     e.Description,
		StaticAddresses: util.StrToMem(e.StaticAddresses),
		Tags:            util.StrToMem(e.Tags),
	}
	if e.DynamicMatch != "" {
		ans.DynamicMatch = &e.DynamicMatch
	}

	return ans
}
