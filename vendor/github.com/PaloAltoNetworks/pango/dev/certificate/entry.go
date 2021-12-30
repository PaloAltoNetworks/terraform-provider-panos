package certificate

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of a
// certificate.
//
// PAN-OS 7.1+.
type Entry struct {
	Name            string
	CommonName      string
	Algorithm       string
	Ca              bool
	NotValidAfter   string
	NotValidBefore  string
	ExpiryEpoch     string
	Subject         string
	SubjectHash     string
	Issuer          string
	IssuerHash      string
	Csr             string
	PublicKey       string
	PrivateKey      string
	PrivateKeyOnHsm bool
	Status          string
	RevokeDateEpoch string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.CommonName = s.CommonName
	o.Algorithm = s.Algorithm
	o.Ca = s.Ca
	o.NotValidAfter = s.NotValidAfter
	o.NotValidBefore = s.NotValidBefore
	o.ExpiryEpoch = s.ExpiryEpoch
	o.Subject = s.Subject
	o.SubjectHash = s.SubjectHash
	o.Issuer = s.Issuer
	o.IssuerHash = s.IssuerHash
	o.Csr = s.Csr
	o.PublicKey = s.PublicKey
	o.PrivateKey = s.PrivateKey
	o.PrivateKeyOnHsm = s.PrivateKeyOnHsm
	o.Status = s.Status
	o.RevokeDateEpoch = s.RevokeDateEpoch
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
		Name:            o.Name,
		CommonName:      o.CommonName,
		Algorithm:       o.Algorithm,
		Ca:              util.AsBool(o.Ca),
		NotValidAfter:   o.NotValidAfter,
		NotValidBefore:  o.NotValidBefore,
		ExpiryEpoch:     o.ExpiryEpoch,
		Subject:         o.Subject,
		SubjectHash:     o.SubjectHash,
		Issuer:          o.Issuer,
		IssuerHash:      o.IssuerHash,
		Csr:             o.Csr,
		PublicKey:       o.PublicKey,
		PrivateKey:      o.PrivateKey,
		PrivateKeyOnHsm: util.AsBool(o.PrivateKeyOnHsm),
		Status:          o.Status,
		RevokeDateEpoch: o.RevokeDateEpoch,
	}

	return ans
}

type entry_v1 struct {
	XMLName         xml.Name `xml:"entry"`
	Name            string   `xml:"name,attr"`
	CommonName      string   `xml:"common-name"`
	Algorithm       string   `xml:"algorithm,omitempty"`
	Ca              string   `xml:"ca,omitempty"`
	NotValidAfter   string   `xml:"not-valid-after,omitempty"`
	NotValidBefore  string   `xml:"not-valid-before,omitempty"`
	ExpiryEpoch     string   `xml:"expiry-epoch,omitempty"`
	Subject         string   `xml:"subject,omitempty"`
	SubjectHash     string   `xml:"subject-hash,omitempty"`
	Issuer          string   `xml:"issuer,omitempty"`
	IssuerHash      string   `xml:"issuer-hash,omitempty"`
	Csr             string   `xml:"csr,omitempty"`
	PublicKey       string   `xml:"public-key,omitempty"`
	PrivateKey      string   `xml:"private-key,omitempty"`
	PrivateKeyOnHsm string   `xml:"private-key-on-hsm,omitempty"`
	Status          string   `xml:"status,omitempty"`
	RevokeDateEpoch string   `xml:"revoke-date-epoch,omitempty"`
}

func (e *entry_v1) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local entry_v1
	ans := local{
		Status: StatusValid,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v1(ans)
	return nil
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:            e.Name,
		CommonName:      e.CommonName,
		Algorithm:       e.Algorithm,
		NotValidAfter:   e.NotValidAfter,
		NotValidBefore:  e.NotValidBefore,
		ExpiryEpoch:     e.ExpiryEpoch,
		Subject:         e.Subject,
		SubjectHash:     e.SubjectHash,
		Issuer:          e.Issuer,
		IssuerHash:      e.IssuerHash,
		Csr:             e.Csr,
		PublicKey:       e.PublicKey,
		PrivateKey:      e.PrivateKey,
		Status:          e.Status,
		RevokeDateEpoch: e.RevokeDateEpoch,
	}

	if e.Ca {
		ans.Ca = util.YesNo(e.Ca)
	}

	if e.PrivateKeyOnHsm {
		ans.PrivateKeyOnHsm = util.YesNo(e.PrivateKeyOnHsm)
	}

	return ans
}
