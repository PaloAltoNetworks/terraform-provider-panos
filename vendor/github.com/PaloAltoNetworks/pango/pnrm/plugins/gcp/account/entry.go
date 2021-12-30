package account

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/plugin"
)

// Entry is a normalized, version independent representation of GCP account credentials.
//
// Note:  GCP Plugin v1.0
type Entry struct {
	Name                         string
	Description                  string
	ProjectId                    string
	ServiceAccountCredentialType string
	CredentialFile               string // encrypted
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Description = s.Description
	o.ProjectId = s.ProjectId
	o.ServiceAccountCredentialType = s.ServiceAccountCredentialType
	o.CredentialFile = s.CredentialFile
}

/** Structs / functions for this namespace. **/

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

func (o Entry) Specify(list []plugin.Info) (string, interface{}, error) {
	_, fn, err := versioning(list)
	if err != nil {
		return o.Name, nil, err
	}

	return o.Name, fn(o), nil
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
		Name:        o.Name,
		Description: o.Description,
		ProjectId:   o.ProjectId,
	}

	if o.Type.Gcp != nil {
		ans.ServiceAccountCredentialType = Project
		ans.CredentialFile = o.Type.Gcp.CredentialFile
	} else if o.Type.Gke != nil {
		ans.ServiceAccountCredentialType = Gke
		ans.CredentialFile = o.Type.Gke.CredentialFile
	}

	return ans
}

type entry_v1 struct {
	XMLName     xml.Name `xml:"entry"`
	Name        string   `xml:"name,attr"`
	Type        actType  `xml:"type"`
	ProjectId   string   `xml:"project-id"`
	Description string   `xml:"description,omitempty"`
}

type actType struct {
	Gke *creds `xml:"gke"`
	Gcp *creds `xml:"gcp"`
}

type creds struct {
	CredentialFile string `xml:"file"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:        e.Name,
		Description: e.Description,
		ProjectId:   e.ProjectId,
	}

	switch e.ServiceAccountCredentialType {
	case Project:
		ans.Type.Gcp = &creds{
			CredentialFile: e.CredentialFile,
		}
	case Gke:
		ans.Type.Gke = &creds{
			CredentialFile: e.CredentialFile,
		}
	}

	return ans
}
