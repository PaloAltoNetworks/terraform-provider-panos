package cluster

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of a GKE cluster.
type Entry struct {
	Name              string
	GcpZone           string
	ClusterCredential string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.GcpZone = s.GcpZone
	o.ClusterCredential = s.ClusterCredential
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
		Name:              o.Answer.Name,
		GcpZone:           o.Answer.GcpZone,
		ClusterCredential: o.Answer.ClusterCredential,
	}

	return ans
}

type entry_v1 struct {
	XMLName           xml.Name `xml:"entry"`
	Name              string   `xml:"name,attr"`
	GcpZone           string   `xml:"gke-zone"`
	ClusterCredential string   `xml:"gke-creds"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:              e.Name,
		GcpZone:           e.GcpZone,
		ClusterCredential: e.ClusterCredential,
	}

	return ans
}
