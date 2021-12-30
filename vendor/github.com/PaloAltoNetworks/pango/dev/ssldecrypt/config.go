package ssldecrypt

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Config is a normalized, version independent representation of
// SSL decrypt settings associated with certificates.
//
// Note: PAN-OS 8.0+
type Config struct {
	ForwardTrustCertificateRsa            string
	ForwardTrustCertificateEcdsa          string
	ForwardUntrustCertificateRsa          string
	ForwardUntrustCertificateEcdsa        string
	RootCaExcludes                        []string
	TrustedRootCas                        []string
	DisabledPredefinedExcludeCertificates []string
	SslDecryptExcludeCertificates         []SslDecryptExcludeCertificate
}

type SslDecryptExcludeCertificate struct {
	Name        string
	Description string
	Exclude     bool
}

// Copy copies the information from source Config `s` to this object.
func (o *Config) Copy(s Config) {
	o.ForwardTrustCertificateRsa = s.ForwardTrustCertificateRsa
	o.ForwardTrustCertificateEcdsa = s.ForwardTrustCertificateEcdsa
	o.ForwardUntrustCertificateRsa = s.ForwardUntrustCertificateRsa
	o.ForwardUntrustCertificateEcdsa = s.ForwardUntrustCertificateEcdsa
	o.RootCaExcludes = util.CopyStringSlice(s.RootCaExcludes)
	o.TrustedRootCas = util.CopyStringSlice(s.TrustedRootCas)
	o.DisabledPredefinedExcludeCertificates = util.CopyStringSlice(s.DisabledPredefinedExcludeCertificates)
	if s.SslDecryptExcludeCertificates == nil {
		o.SslDecryptExcludeCertificates = nil
	} else {
		o.SslDecryptExcludeCertificates = make([]SslDecryptExcludeCertificate, 0, len(s.SslDecryptExcludeCertificates))
		for _, x := range s.SslDecryptExcludeCertificates {
			o.SslDecryptExcludeCertificates = append(o.SslDecryptExcludeCertificates, SslDecryptExcludeCertificate{
				Name:        x.Name,
				Description: x.Description,
				Exclude:     x.Exclude,
			})
		}
	}
}

/** Structs / functions for this namespace. **/

func (o Config) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return "", fn(o)
}

type normalizer interface {
	Normalize() []Config
	Names() []string
}

type container_v1 struct {
	Answer []config_v1 `xml:"ssl-decrypt"`
}

func (o *container_v1) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for _ = range o.Answer {
		ans = append(ans, "")
	}

	return ans
}

func (o *config_v1) normalize() Config {
	ans := Config{
		RootCaExcludes:                        util.MemToStr(o.RootCaExcludes),
		TrustedRootCas:                        util.MemToStr(o.TrustedRootCas),
		DisabledPredefinedExcludeCertificates: util.MemToStr(o.DisabledPredefinedExcludeCertificates),
	}

	if o.ForwardTrustCerts != nil {
		ans.ForwardTrustCertificateRsa = o.ForwardTrustCerts.Rsa
		ans.ForwardTrustCertificateEcdsa = o.ForwardTrustCerts.Ecdsa
	}

	if o.ForwardUntrustCerts != nil {
		ans.ForwardUntrustCertificateRsa = o.ForwardUntrustCerts.Rsa
		ans.ForwardUntrustCertificateEcdsa = o.ForwardUntrustCerts.Ecdsa
	}

	if o.SslDecryptExcludeCertificates != nil {
		list := make([]SslDecryptExcludeCertificate, 0, len(o.SslDecryptExcludeCertificates.Entries))
		for _, x := range o.SslDecryptExcludeCertificates.Entries {
			list = append(list, SslDecryptExcludeCertificate{
				Name:        x.Name,
				Description: x.Description,
				Exclude:     util.AsBool(x.Exclude),
			})
		}
		ans.SslDecryptExcludeCertificates = list
	}

	return ans
}

type config_v1 struct {
	XMLName                               xml.Name         `xml:"ssl-decrypt"`
	ForwardTrustCerts                     *trustUntrust    `xml:"forward-trust-certificate"`
	ForwardUntrustCerts                   *trustUntrust    `xml:"forward-untrust-certificate"`
	RootCaExcludes                        *util.MemberType `xml:"root-ca-exclude-list"`
	TrustedRootCas                        *util.MemberType `xml:"trusted-root-CA"`
	DisabledPredefinedExcludeCertificates *util.MemberType `xml:"disabled-ssl-exclude-cert-from-predefined"`
	SslDecryptExcludeCertificates         *sdec            `xml:"ssl-exclude-cert"`
}

type trustUntrust struct {
	Rsa   string `xml:"rsa,omitempty"`
	Ecdsa string `xml:"ecdsa,omitempty"`
}

type sdec struct {
	Entries []sdecEntry `xml:"entry"`
}

type sdecEntry struct {
	Name        string `xml:"name,attr"`
	Description string `xml:"description,omitempty"`
	Exclude     string `xml:"exclude"`
}

func specify_v1(e Config) interface{} {
	ans := config_v1{
		RootCaExcludes:                        util.StrToMem(e.RootCaExcludes),
		TrustedRootCas:                        util.StrToMem(e.TrustedRootCas),
		DisabledPredefinedExcludeCertificates: util.StrToMem(e.DisabledPredefinedExcludeCertificates),
	}

	if e.ForwardTrustCertificateRsa != "" || e.ForwardTrustCertificateEcdsa != "" {
		ans.ForwardTrustCerts = &trustUntrust{
			Rsa:   e.ForwardTrustCertificateRsa,
			Ecdsa: e.ForwardTrustCertificateEcdsa,
		}
	}

	if e.ForwardUntrustCertificateRsa != "" || e.ForwardUntrustCertificateEcdsa != "" {
		ans.ForwardUntrustCerts = &trustUntrust{
			Rsa:   e.ForwardUntrustCertificateRsa,
			Ecdsa: e.ForwardUntrustCertificateEcdsa,
		}
	}

	if len(e.SslDecryptExcludeCertificates) > 0 {
		list := make([]sdecEntry, 0, len(e.SslDecryptExcludeCertificates))
		for _, x := range e.SslDecryptExcludeCertificates {
			list = append(list, sdecEntry{
				Name:        x.Name,
				Description: x.Description,
				Exclude:     util.YesNo(x.Exclude),
			})
		}

		ans.SslDecryptExcludeCertificates = &sdec{Entries: list}
	}

	return ans
}
