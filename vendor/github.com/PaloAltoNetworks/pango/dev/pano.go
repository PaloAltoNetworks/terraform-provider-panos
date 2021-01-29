package dev

import (
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/PaloAltoNetworks/pango/dev/profile/email"
	emailsrv "github.com/PaloAltoNetworks/pango/dev/profile/email/server"
	"github.com/PaloAltoNetworks/pango/dev/profile/http"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/header"
	"github.com/PaloAltoNetworks/pango/dev/profile/http/param"
	httpsrv "github.com/PaloAltoNetworks/pango/dev/profile/http/server"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v2c"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v3"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"
	syslogsrv "github.com/PaloAltoNetworks/pango/dev/profile/syslog/server"
)

// PanoDev is the client.Device namespace.
type PanoDev struct {
	EmailServer         *emailsrv.PanoServer
	EmailServerProfile  *email.PanoEmail
	HttpHeader          *header.PanoHeader
	HttpParam           *param.PanoParam
	HttpServer          *httpsrv.PanoServer
	HttpServerProfile   *http.PanoHttp
	SnmpServerProfile   *snmp.PanoSnmp
	SnmpV2cServer       *v2c.PanoV2c
	SnmpV3Server        *v3.PanoV3
	SyslogServer        *syslogsrv.PanoServer
	SyslogServerProfile *syslog.PanoSyslog
}

// Initialize is invoked on client.Initialize().
func (c *PanoDev) Initialize(i util.XapiClient) {
	c.EmailServer = &emailsrv.PanoServer{}
	c.EmailServer.Initialize(i)

	c.EmailServerProfile = &email.PanoEmail{}
	c.EmailServerProfile.Initialize(i)

	c.HttpHeader = &header.PanoHeader{}
	c.HttpHeader.Initialize(i)

	c.HttpParam = &param.PanoParam{}
	c.HttpParam.Initialize(i)

	c.HttpServer = &httpsrv.PanoServer{}
	c.HttpServer.Initialize(i)

	c.HttpServerProfile = &http.PanoHttp{}
	c.HttpServerProfile.Initialize(i)

	c.SnmpServerProfile = &snmp.PanoSnmp{}
	c.SnmpServerProfile.Initialize(i)

	c.SnmpV2cServer = &v2c.PanoV2c{}
	c.SnmpV2cServer.Initialize(i)

	c.SnmpV3Server = &v3.PanoV3{}
	c.SnmpV3Server.Initialize(i)

	c.SyslogServer = &syslogsrv.PanoServer{}
	c.SyslogServer.Initialize(i)

	c.SyslogServerProfile = &syslog.PanoSyslog{}
	c.SyslogServerProfile.Initialize(i)
}
