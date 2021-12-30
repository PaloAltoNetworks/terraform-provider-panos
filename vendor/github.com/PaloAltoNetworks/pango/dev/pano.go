package dev

import (
	"github.com/PaloAltoNetworks/pango/util"

	cert "github.com/PaloAltoNetworks/pango/dev/certificate"
	"github.com/PaloAltoNetworks/pango/dev/localuserdb/group"
	"github.com/PaloAltoNetworks/pango/dev/localuserdb/user"
	"github.com/PaloAltoNetworks/pango/dev/profile/certificate"
	"github.com/PaloAltoNetworks/pango/dev/profile/email"
	"github.com/PaloAltoNetworks/pango/dev/profile/http"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"
	"github.com/PaloAltoNetworks/pango/dev/ssldecrypt"
	"github.com/PaloAltoNetworks/pango/dev/vminfosource"
)

// Panorama is the client.Device namespace.
type Panorama struct {
	Certificate         *cert.Panorama
	CertificateProfile  *certificate.Panorama
	EmailServerProfile  *email.Panorama
	HttpServerProfile   *http.Panorama
	LocalUserDbGroup    *group.Panorama
	LocalUserDbUser     *user.Panorama
	SslDecrypt          *ssldecrypt.Panorama
	SnmpServerProfile   *snmp.Panorama
	SyslogServerProfile *syslog.Panorama
	VmInfoSource        *vminfosource.Panorama
}

// Initialize is invoked on client.Initialize().
func PanoramaNamespace(x util.XapiClient) *Panorama {
	return &Panorama{
		Certificate:         cert.PanoramaNamespace(x),
		CertificateProfile:  certificate.PanoramaNamespace(x),
		EmailServerProfile:  email.PanoramaNamespace(x),
		HttpServerProfile:   http.PanoramaNamespace(x),
		LocalUserDbGroup:    group.PanoramaNamespace(x),
		LocalUserDbUser:     user.PanoramaNamespace(x),
		SslDecrypt:          ssldecrypt.PanoramaNamespace(x),
		SnmpServerProfile:   snmp.PanoramaNamespace(x),
		SyslogServerProfile: syslog.PanoramaNamespace(x),
		VmInfoSource:        vminfosource.PanoramaNamespace(x),
	}
}
