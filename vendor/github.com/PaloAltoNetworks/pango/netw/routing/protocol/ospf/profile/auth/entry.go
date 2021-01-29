package auth

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an OSPF
// auth profile.
type Entry struct {
	Name     string
	AuthType string
	Password string
	Md5Keys  []Md5Key
}

type Md5Key struct {
	KeyId     int
	Key       string
	Preferred bool
}

func (o *Entry) Copy(s Entry) {
	o.AuthType = s.AuthType
	o.Password = s.Password
	if s.Md5Keys == nil {
		o.Md5Keys = nil
	} else {
		o.Md5Keys = make([]Md5Key, len(s.Md5Keys))
		copy(o.Md5Keys, s.Md5Keys)
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
		Name: o.Name,
	}

	if o.Password != "" {
		ans.AuthType = AuthTypePassword
		ans.Password = o.Password
	} else if len(o.Md5.Md5Keys) > 0 {
		ans.AuthType = AuthTypeMd5
		ans.Md5Keys = make([]Md5Key, 0, len(o.Md5.Md5Keys))
		for i := range o.Md5.Md5Keys {
			key := Md5Key{
				KeyId:     o.Md5.Md5Keys[i].Name,
				Key:       o.Md5.Md5Keys[i].Key,
				Preferred: util.AsBool(o.Md5.Md5Keys[i].Preferred),
			}
			ans.Md5Keys = append(ans.Md5Keys, key)
		}
	}

	return ans
}

type entry_v1 struct {
	XMLName  xml.Name `xml:"entry"`
	Name     string   `xml:"name,attr"`
	Password string   `xml:"password,omitempty"`
	Md5      *md5     `xml:"md5"`
}

type md5 struct {
	Md5Keys []md5Key `xml:"entry"`
}

type md5Key struct {
	Name      int    `xml:"name,attr"`
	Key       string `xml:"key"`
	Preferred string `xml:"preferred"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name: e.Name,
	}

	switch e.AuthType {
	case AuthTypePassword:
		ans.Password = e.Password
	case AuthTypeMd5:
		ans.Md5 = &md5{}
		ans.Md5.Md5Keys = make([]md5Key, 0, len(e.Md5Keys))
		for i := range e.Md5Keys {
			key := md5Key{
				Name:      e.Md5Keys[i].KeyId,
				Key:       e.Md5Keys[i].Key,
				Preferred: util.YesNo(e.Md5Keys[i].Preferred),
			}
			ans.Md5.Md5Keys = append(ans.Md5.Md5Keys, key)
		}
	}

	return ans
}
