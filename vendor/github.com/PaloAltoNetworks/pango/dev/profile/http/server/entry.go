package server

import (
	"encoding/xml"
)

// Entry is a normalized, version independent representation of an http server.
//
// PAN-OS 7.1+.
type Entry struct {
	Name               string
	Address            string
	Protocol           string
	Port               int
	HttpMethod         string
	Username           string
	Password           string // encrypted
	TlsVersion         string // 9.0+
	CertificateProfile string // 9.0+
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.Address = s.Address
	o.Protocol = s.Protocol
	o.Port = s.Port
	o.HttpMethod = s.HttpMethod
	o.Username = s.Username
	o.Password = s.Password
	o.TlsVersion = s.TlsVersion
	o.CertificateProfile = s.CertificateProfile
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
		Name:       o.Answer.Name,
		Address:    o.Answer.Address,
		Protocol:   o.Answer.Protocol,
		Port:       o.Answer.Port,
		HttpMethod: o.Answer.HttpMethod,
		Username:   o.Answer.Username,
		Password:   o.Answer.Password,
	}

	return ans
}

type container_v2 struct {
	Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
	ans := Entry{
		Name:               o.Answer.Name,
		Address:            o.Answer.Address,
		Protocol:           o.Answer.Protocol,
		Port:               o.Answer.Port,
		HttpMethod:         o.Answer.HttpMethod,
		Username:           o.Answer.Username,
		Password:           o.Answer.Password,
		TlsVersion:         o.Answer.TlsVersion,
		CertificateProfile: o.Answer.CertificateProfile,
	}

	return ans
}

type entry_v1 struct {
	XMLName    xml.Name `xml:"entry"`
	Name       string   `xml:"name,attr"`
	Address    string   `xml:"address"`
	Protocol   string   `xml:"protocol,omitempty"`
	Port       int      `xml:"port,omitempty"`
	HttpMethod string   `xml:"http-method"`
	Username   string   `xml:"username,omitempty"`
	Password   string   `xml:"password,omitempty"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:       e.Name,
		Address:    e.Address,
		Protocol:   e.Protocol,
		Port:       e.Port,
		HttpMethod: e.HttpMethod,
		Username:   e.Username,
		Password:   e.Password,
	}

	return ans
}

type entry_v2 struct {
	XMLName            xml.Name `xml:"entry"`
	Name               string   `xml:"name,attr"`
	Address            string   `xml:"address"`
	Protocol           string   `xml:"protocol,omitempty"`
	Port               int      `xml:"port,omitempty"`
	HttpMethod         string   `xml:"http-method"`
	Username           string   `xml:"username,omitempty"`
	Password           string   `xml:"password,omitempty"`
	TlsVersion         string   `xml:"tls-version,omitempty"`
	CertificateProfile string   `xml:"certificate-profile,omitempty"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:               e.Name,
		Address:            e.Address,
		Protocol:           e.Protocol,
		Port:               e.Port,
		HttpMethod:         e.HttpMethod,
		Username:           e.Username,
		Password:           e.Password,
		TlsVersion:         e.TlsVersion,
		CertificateProfile: e.CertificateProfile,
	}

	return ans
}
