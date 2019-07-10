package dev

import (
    "github.com/PaloAltoNetworks/pango/util"

    "github.com/PaloAltoNetworks/pango/dev/general"
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
    "github.com/PaloAltoNetworks/pango/dev/telemetry"
)


// FwDev is the client.Device namespace.
type FwDev struct {
    EmailServer *emailsrv.FwServer
    EmailServerProfile *email.FwEmail
    GeneralSettings *general.FwGeneral
    HttpHeader *header.FwHeader
    HttpParam *param.FwParam
    HttpServer *httpsrv.FwServer
    HttpServerProfile *http.FwHttp
    SnmpServerProfile *snmp.FwSnmp
    SnmpV2cServer *v2c.FwV2c
    SnmpV3Server *v3.FwV3
    SyslogServer *syslogsrv.FwServer
    SyslogServerProfile *syslog.FwSyslog
    Telemetry *telemetry.FwTelemetry
}

// Initialize is invoked on client.Initialize().
func (c *FwDev) Initialize(i util.XapiClient) {
    c.EmailServer = &emailsrv.FwServer{}
    c.EmailServer.Initialize(i)

    c.EmailServerProfile = &email.FwEmail{}
    c.EmailServerProfile.Initialize(i)

    c.GeneralSettings = &general.FwGeneral{}
    c.GeneralSettings.Initialize(i)

    c.HttpHeader = &header.FwHeader{}
    c.HttpHeader.Initialize(i)

    c.HttpParam = &param.FwParam{}
    c.HttpParam.Initialize(i)

    c.HttpServer = &httpsrv.FwServer{}
    c.HttpServer.Initialize(i)

    c.HttpServerProfile = &http.FwHttp{}
    c.HttpServerProfile.Initialize(i)

    c.SnmpServerProfile = &snmp.FwSnmp{}
    c.SnmpServerProfile.Initialize(i)

    c.SnmpV2cServer = &v2c.FwV2c{}
    c.SnmpV2cServer.Initialize(i)

    c.SnmpV3Server = &v3.FwV3{}
    c.SnmpV3Server.Initialize(i)

    c.SyslogServer = &syslogsrv.FwServer{}
    c.SyslogServer.Initialize(i)

    c.SyslogServerProfile = &syslog.FwSyslog{}
    c.SyslogServerProfile.Initialize(i)

    c.Telemetry = &telemetry.FwTelemetry{}
    c.Telemetry.Initialize(i)
}
