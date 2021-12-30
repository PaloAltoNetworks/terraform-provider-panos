package certificate

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of
// a certificate profile.
//
// Leave UsernameField as an empty string to specify a username field
// of `None`.
//
// Note:  Entry.CertificateStatusTimeout=0 is a valid setting, so make
// sure to have the desired value configured before doing Set() / Edit().
//
// Note:
type Entry struct {
	Name                            string
	UsernameField                   string
	UsernameFieldValue              string
	Domain                          string
	Certificates                    []Certificate
	UseCrl                          bool
	UseOcsp                         bool
	CrlReceiveTimeout               int
	OcspReceiveTimeout              int
	CertificateStatusTimeout        int
	BlockUnknownCertificate         bool
	BlockCertificateTimeout         bool
	BlockUnauthenticatedCertificate bool // 7.1+
	BlockExpiredCertificate         bool // 8.1+
	OcspExcludeNonce                bool // 9.0+
}

type Certificate struct {
	Name                  string
	DefaultOcspUrl        string
	OcspVerifyCertificate string
	TemplateName          string // 9.0+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.UsernameField = s.UsernameField
	o.UsernameFieldValue = s.UsernameFieldValue
	o.Domain = s.Domain
	if s.Certificates == nil {
		o.Certificates = nil
	} else {
		o.Certificates = make([]Certificate, 0, len(s.Certificates))
		for _, x := range s.Certificates {
			o.Certificates = append(o.Certificates, Certificate{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
				TemplateName:          x.TemplateName,
			})
		}
	}
	o.UseCrl = s.UseCrl
	o.UseOcsp = s.UseOcsp
	o.CrlReceiveTimeout = s.CrlReceiveTimeout
	o.OcspReceiveTimeout = s.OcspReceiveTimeout
	o.CertificateStatusTimeout = s.CertificateStatusTimeout
	o.BlockUnknownCertificate = s.BlockUnknownCertificate
	o.BlockCertificateTimeout = s.BlockCertificateTimeout
	o.BlockUnauthenticatedCertificate = s.BlockUnauthenticatedCertificate
	o.BlockExpiredCertificate = s.BlockExpiredCertificate
	o.OcspExcludeNonce = s.OcspExcludeNonce
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

func (o container_v1) Names() []string {
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

type entry_v1 struct {
	XMLName                  xml.Name  `xml:"entry"`
	Name                     string    `xml:"name,attr"`
	UsernameField            *uf       `xml:"username-field"`
	Domain                   string    `xml:"domain,omitempty"`
	Certificates             *certs_v1 `xml:"CA"`
	UseCrl                   string    `xml:"use-crl"`
	UseOcsp                  string    `xml:"use-ocsp"`
	CrlReceiveTimeout        int       `xml:"crl-receive-timeout,omitempty"`
	OcspReceiveTimeout       int       `xml:"ocsp-receive-timeout,omitempty"`
	CertificateStatusTimeout int       `xml:"cert-status-timeout"`
	BlockUnknownCertificate  string    `xml:"block-unknown-cert"`
	BlockCertificateTimeout  string    `xml:"block-timeout-cert"`
}

type uf struct {
	Subject    string `xml:"subject,omitempty"`
	SubjectAlt string `xml:"subject-alt,omitempty"`
}

type certs_v1 struct {
	Entries []certEntry_v1 `xml:"entry"`
}

type certEntry_v1 struct {
	Name                  string `xml:"name,attr"`
	DefaultOcspUrl        string `xml:"default-ocsp-url,omitempty"`
	OcspVerifyCertificate string `xml:"ocsp-verify-cert,omitempty"`
}

func (e *entry_v1) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local entry_v1
	ans := local{
		CrlReceiveTimeout:        5,
		OcspReceiveTimeout:       5,
		CertificateStatusTimeout: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v1(ans)
	return nil
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                     o.Name,
		Domain:                   o.Domain,
		UseCrl:                   util.AsBool(o.UseCrl),
		UseOcsp:                  util.AsBool(o.UseOcsp),
		CrlReceiveTimeout:        o.CrlReceiveTimeout,
		OcspReceiveTimeout:       o.OcspReceiveTimeout,
		CertificateStatusTimeout: o.CertificateStatusTimeout,
		BlockUnknownCertificate:  util.AsBool(o.BlockUnknownCertificate),
		BlockCertificateTimeout:  util.AsBool(o.BlockCertificateTimeout),
	}

	if o.UsernameField != nil {
		if o.UsernameField.Subject != "" {
			ans.UsernameField = UsernameFieldSubject
			ans.UsernameFieldValue = o.UsernameField.Subject
		} else if o.UsernameField.SubjectAlt != "" {
			ans.UsernameField = UsernameFieldSubjectAlt
			ans.UsernameFieldValue = o.UsernameField.SubjectAlt
		}
	}

	if o.Certificates != nil {
		ans.Certificates = make([]Certificate, 0, len(o.Certificates.Entries))
		for _, x := range o.Certificates.Entries {
			ans.Certificates = append(ans.Certificates, Certificate{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
			})
		}
	}

	return ans
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                     e.Name,
		Domain:                   e.Domain,
		UseCrl:                   util.YesNo(e.UseCrl),
		UseOcsp:                  util.YesNo(e.UseOcsp),
		CrlReceiveTimeout:        e.CrlReceiveTimeout,
		OcspReceiveTimeout:       e.OcspReceiveTimeout,
		CertificateStatusTimeout: e.CertificateStatusTimeout,
		BlockUnknownCertificate:  util.YesNo(e.BlockUnknownCertificate),
		BlockCertificateTimeout:  util.YesNo(e.BlockCertificateTimeout),
	}

	switch e.UsernameField {
	case UsernameFieldSubject:
		ans.UsernameField = &uf{
			Subject: e.UsernameFieldValue,
		}
	case UsernameFieldSubjectAlt:
		ans.UsernameField = &uf{
			SubjectAlt: e.UsernameFieldValue,
		}
	}

	if len(e.Certificates) > 0 {
		list := make([]certEntry_v1, 0, len(e.Certificates))
		for _, x := range e.Certificates {
			list = append(list, certEntry_v1{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
			})
		}
		ans.Certificates = &certs_v1{Entries: list}
	}

	return ans
}

// PAN-OS 7.1
type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v2 struct {
	XMLName                         xml.Name  `xml:"entry"`
	Name                            string    `xml:"name,attr"`
	UsernameField                   *uf       `xml:"username-field"`
	Domain                          string    `xml:"domain,omitempty"`
	Certificates                    *certs_v1 `xml:"CA"`
	UseCrl                          string    `xml:"use-crl"`
	UseOcsp                         string    `xml:"use-ocsp"`
	CrlReceiveTimeout               int       `xml:"crl-receive-timeout,omitempty"`
	OcspReceiveTimeout              int       `xml:"ocsp-receive-timeout,omitempty"`
	CertificateStatusTimeout        int       `xml:"cert-status-timeout"`
	BlockUnknownCertificate         string    `xml:"block-unknown-cert"`
	BlockCertificateTimeout         string    `xml:"block-timeout-cert"`
	BlockUnauthenticatedCertificate string    `xml:"block-unauthenticated-cert"`
}

func (e *entry_v2) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local entry_v2
	ans := local{
		CrlReceiveTimeout:        5,
		OcspReceiveTimeout:       5,
		CertificateStatusTimeout: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v2(ans)
	return nil
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:                            o.Name,
		Domain:                          o.Domain,
		UseCrl:                          util.AsBool(o.UseCrl),
		UseOcsp:                         util.AsBool(o.UseOcsp),
		CrlReceiveTimeout:               o.CrlReceiveTimeout,
		OcspReceiveTimeout:              o.OcspReceiveTimeout,
		CertificateStatusTimeout:        o.CertificateStatusTimeout,
		BlockUnknownCertificate:         util.AsBool(o.BlockUnknownCertificate),
		BlockCertificateTimeout:         util.AsBool(o.BlockCertificateTimeout),
		BlockUnauthenticatedCertificate: util.AsBool(o.BlockUnauthenticatedCertificate),
	}

	if o.UsernameField != nil {
		if o.UsernameField.Subject != "" {
			ans.UsernameField = UsernameFieldSubject
			ans.UsernameFieldValue = o.UsernameField.Subject
		} else if o.UsernameField.SubjectAlt != "" {
			ans.UsernameField = UsernameFieldSubjectAlt
			ans.UsernameFieldValue = o.UsernameField.SubjectAlt
		}
	}

	if o.Certificates != nil {
		ans.Certificates = make([]Certificate, 0, len(o.Certificates.Entries))
		for _, x := range o.Certificates.Entries {
			ans.Certificates = append(ans.Certificates, Certificate{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
			})
		}
	}

	return ans
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                            e.Name,
		Domain:                          e.Domain,
		UseCrl:                          util.YesNo(e.UseCrl),
		UseOcsp:                         util.YesNo(e.UseOcsp),
		CrlReceiveTimeout:               e.CrlReceiveTimeout,
		OcspReceiveTimeout:              e.OcspReceiveTimeout,
		CertificateStatusTimeout:        e.CertificateStatusTimeout,
		BlockUnknownCertificate:         util.YesNo(e.BlockUnknownCertificate),
		BlockCertificateTimeout:         util.YesNo(e.BlockCertificateTimeout),
		BlockUnauthenticatedCertificate: util.YesNo(e.BlockUnauthenticatedCertificate),
	}

	switch e.UsernameField {
	case UsernameFieldSubject:
		ans.UsernameField = &uf{
			Subject: e.UsernameFieldValue,
		}
	case UsernameFieldSubjectAlt:
		ans.UsernameField = &uf{
			SubjectAlt: e.UsernameFieldValue,
		}
	}

	if len(e.Certificates) > 0 {
		list := make([]certEntry_v1, 0, len(e.Certificates))
		for _, x := range e.Certificates {
			list = append(list, certEntry_v1{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
			})
		}
		ans.Certificates = &certs_v1{Entries: list}
	}

	return ans
}

// PAN-OS 8.1
type container_v3 struct {
	Answer []entry_v3 `xml:"entry"`
}

func (o container_v3) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v3) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v3 struct {
	XMLName                         xml.Name  `xml:"entry"`
	Name                            string    `xml:"name,attr"`
	UsernameField                   *uf       `xml:"username-field"`
	Domain                          string    `xml:"domain,omitempty"`
	Certificates                    *certs_v1 `xml:"CA"`
	UseCrl                          string    `xml:"use-crl"`
	UseOcsp                         string    `xml:"use-ocsp"`
	CrlReceiveTimeout               int       `xml:"crl-receive-timeout,omitempty"`
	OcspReceiveTimeout              int       `xml:"ocsp-receive-timeout,omitempty"`
	CertificateStatusTimeout        int       `xml:"cert-status-timeout"`
	BlockUnknownCertificate         string    `xml:"block-unknown-cert"`
	BlockCertificateTimeout         string    `xml:"block-timeout-cert"`
	BlockUnauthenticatedCertificate string    `xml:"block-unauthenticated-cert"`
	BlockExpiredCertificate         string    `xml:"block-expired-cert"`
}

func (e *entry_v3) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local entry_v3
	ans := local{
		CrlReceiveTimeout:        5,
		OcspReceiveTimeout:       5,
		CertificateStatusTimeout: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v3(ans)
	return nil
}

func (o *entry_v3) normalize() Entry {
	ans := Entry{
		Name:                            o.Name,
		Domain:                          o.Domain,
		UseCrl:                          util.AsBool(o.UseCrl),
		UseOcsp:                         util.AsBool(o.UseOcsp),
		CrlReceiveTimeout:               o.CrlReceiveTimeout,
		OcspReceiveTimeout:              o.OcspReceiveTimeout,
		CertificateStatusTimeout:        o.CertificateStatusTimeout,
		BlockUnknownCertificate:         util.AsBool(o.BlockUnknownCertificate),
		BlockCertificateTimeout:         util.AsBool(o.BlockCertificateTimeout),
		BlockUnauthenticatedCertificate: util.AsBool(o.BlockUnauthenticatedCertificate),
		BlockExpiredCertificate:         util.AsBool(o.BlockExpiredCertificate),
	}

	if o.UsernameField != nil {
		if o.UsernameField.Subject != "" {
			ans.UsernameField = UsernameFieldSubject
			ans.UsernameFieldValue = o.UsernameField.Subject
		} else if o.UsernameField.SubjectAlt != "" {
			ans.UsernameField = UsernameFieldSubjectAlt
			ans.UsernameFieldValue = o.UsernameField.SubjectAlt
		}
	}

	if o.Certificates != nil {
		ans.Certificates = make([]Certificate, 0, len(o.Certificates.Entries))
		for _, x := range o.Certificates.Entries {
			ans.Certificates = append(ans.Certificates, Certificate{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
			})
		}
	}

	return ans
}

func specify_v3(e Entry) interface{} {
	ans := entry_v3{
		Name:                            e.Name,
		Domain:                          e.Domain,
		UseCrl:                          util.YesNo(e.UseCrl),
		UseOcsp:                         util.YesNo(e.UseOcsp),
		CrlReceiveTimeout:               e.CrlReceiveTimeout,
		OcspReceiveTimeout:              e.OcspReceiveTimeout,
		CertificateStatusTimeout:        e.CertificateStatusTimeout,
		BlockUnknownCertificate:         util.YesNo(e.BlockUnknownCertificate),
		BlockCertificateTimeout:         util.YesNo(e.BlockCertificateTimeout),
		BlockUnauthenticatedCertificate: util.YesNo(e.BlockUnauthenticatedCertificate),
		BlockExpiredCertificate:         util.YesNo(e.BlockExpiredCertificate),
	}

	switch e.UsernameField {
	case UsernameFieldSubject:
		ans.UsernameField = &uf{
			Subject: e.UsernameFieldValue,
		}
	case UsernameFieldSubjectAlt:
		ans.UsernameField = &uf{
			SubjectAlt: e.UsernameFieldValue,
		}
	}

	if len(e.Certificates) > 0 {
		list := make([]certEntry_v1, 0, len(e.Certificates))
		for _, x := range e.Certificates {
			list = append(list, certEntry_v1{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
			})
		}
		ans.Certificates = &certs_v1{Entries: list}
	}

	return ans
}

// PAN-OS 9.0
type container_v4 struct {
	Answer []entry_v4 `xml:"entry"`
}

func (o container_v4) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *container_v4) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

type entry_v4 struct {
	XMLName                         xml.Name  `xml:"entry"`
	Name                            string    `xml:"name,attr"`
	UsernameField                   *uf       `xml:"username-field"`
	Domain                          string    `xml:"domain,omitempty"`
	Certificates                    *certs_v2 `xml:"CA"`
	UseCrl                          string    `xml:"use-crl"`
	UseOcsp                         string    `xml:"use-ocsp"`
	CrlReceiveTimeout               int       `xml:"crl-receive-timeout,omitempty"`
	OcspReceiveTimeout              int       `xml:"ocsp-receive-timeout,omitempty"`
	CertificateStatusTimeout        int       `xml:"cert-status-timeout"`
	BlockUnknownCertificate         string    `xml:"block-unknown-cert"`
	BlockCertificateTimeout         string    `xml:"block-timeout-cert"`
	BlockUnauthenticatedCertificate string    `xml:"block-unauthenticated-cert"`
	BlockExpiredCertificate         string    `xml:"block-expired-cert"`
	OcspExcludeNonce                string    `xml:"ocsp-exclude-nonce"`
}

type certs_v2 struct {
	Entries []certEntry_v2 `xml:"entry"`
}

type certEntry_v2 struct {
	Name                  string `xml:"name,attr"`
	DefaultOcspUrl        string `xml:"default-ocsp-url,omitempty"`
	OcspVerifyCertificate string `xml:"ocsp-verify-cert,omitempty"`
	TemplateName          string `xml:"template-name,omitempty"`
}

func (e *entry_v4) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type local entry_v4
	ans := local{
		CrlReceiveTimeout:        5,
		OcspReceiveTimeout:       5,
		CertificateStatusTimeout: 5,
	}
	if err := d.DecodeElement(&ans, &start); err != nil {
		return err
	}
	*e = entry_v4(ans)
	return nil
}

func (o *entry_v4) normalize() Entry {
	ans := Entry{
		Name:                            o.Name,
		Domain:                          o.Domain,
		UseCrl:                          util.AsBool(o.UseCrl),
		UseOcsp:                         util.AsBool(o.UseOcsp),
		CrlReceiveTimeout:               o.CrlReceiveTimeout,
		OcspReceiveTimeout:              o.OcspReceiveTimeout,
		CertificateStatusTimeout:        o.CertificateStatusTimeout,
		BlockUnknownCertificate:         util.AsBool(o.BlockUnknownCertificate),
		BlockCertificateTimeout:         util.AsBool(o.BlockCertificateTimeout),
		BlockUnauthenticatedCertificate: util.AsBool(o.BlockUnauthenticatedCertificate),
		BlockExpiredCertificate:         util.AsBool(o.BlockExpiredCertificate),
		OcspExcludeNonce:                util.AsBool(o.OcspExcludeNonce),
	}

	if o.UsernameField != nil {
		if o.UsernameField.Subject != "" {
			ans.UsernameField = UsernameFieldSubject
			ans.UsernameFieldValue = o.UsernameField.Subject
		} else if o.UsernameField.SubjectAlt != "" {
			ans.UsernameField = UsernameFieldSubjectAlt
			ans.UsernameFieldValue = o.UsernameField.SubjectAlt
		}
	}

	if o.Certificates != nil {
		ans.Certificates = make([]Certificate, 0, len(o.Certificates.Entries))
		for _, x := range o.Certificates.Entries {
			ans.Certificates = append(ans.Certificates, Certificate{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
				TemplateName:          x.TemplateName,
			})
		}
	}

	return ans
}

func specify_v4(e Entry) interface{} {
	ans := entry_v4{
		Name:                            e.Name,
		Domain:                          e.Domain,
		UseCrl:                          util.YesNo(e.UseCrl),
		UseOcsp:                         util.YesNo(e.UseOcsp),
		CrlReceiveTimeout:               e.CrlReceiveTimeout,
		OcspReceiveTimeout:              e.OcspReceiveTimeout,
		CertificateStatusTimeout:        e.CertificateStatusTimeout,
		BlockUnknownCertificate:         util.YesNo(e.BlockUnknownCertificate),
		BlockCertificateTimeout:         util.YesNo(e.BlockCertificateTimeout),
		BlockUnauthenticatedCertificate: util.YesNo(e.BlockUnauthenticatedCertificate),
		BlockExpiredCertificate:         util.YesNo(e.BlockExpiredCertificate),
		OcspExcludeNonce:                util.YesNo(e.OcspExcludeNonce),
	}

	switch e.UsernameField {
	case UsernameFieldSubject:
		ans.UsernameField = &uf{
			Subject: e.UsernameFieldValue,
		}
	case UsernameFieldSubjectAlt:
		ans.UsernameField = &uf{
			SubjectAlt: e.UsernameFieldValue,
		}
	}

	if len(e.Certificates) > 0 {
		list := make([]certEntry_v2, 0, len(e.Certificates))
		for _, x := range e.Certificates {
			list = append(list, certEntry_v2{
				Name:                  x.Name,
				DefaultOcspUrl:        x.DefaultOcspUrl,
				OcspVerifyCertificate: x.OcspVerifyCertificate,
				TemplateName:          x.TemplateName,
			})
		}
		ans.Certificates = &certs_v2{Entries: list}
	}

	return ans
}
