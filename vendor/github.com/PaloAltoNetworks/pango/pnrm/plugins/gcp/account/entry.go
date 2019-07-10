package account

import (
    "encoding/xml"
)


// Entry is a normalized, version independent representation of GCP account credentials.
type Entry struct {
    Name string
    Description string
    ProjectId string
    ServiceAccountCredentialType string
    CredentialFile string // encrypted
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Description = s.Description
    o.ServiceAccountCredentialType = s.ServiceAccountCredentialType
    o.ProjectId = s.ProjectId
    o.CredentialFile = s.CredentialFile
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
        Description: o.Answer.Description,
        ProjectId: o.Answer.ProjectId,
    }

    if o.Answer.Type.Gcp != nil {
        ans.ServiceAccountCredentialType = Project
        ans.CredentialFile = o.Answer.Type.Gcp.CredentialFile
    } else if o.Answer.Type.Gke != nil {
        ans.ServiceAccountCredentialType = Gke
        ans.CredentialFile = o.Answer.Type.Gke.CredentialFile
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Type actType `xml:"type"`
    ProjectId string `xml:"project-id"`
    Description string `xml:"description,omitempty"`
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
        Name: e.Name,
        Description: e.Description,
        ProjectId: e.ProjectId,
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
