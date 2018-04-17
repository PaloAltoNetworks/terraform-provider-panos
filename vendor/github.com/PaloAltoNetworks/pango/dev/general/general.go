// Package general is the client.Device.GeneralSettings namespace.
//
// Normalized object: Config
package general

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


// Constants for NTP auth types.
const (
    NoAuth = "none"
    AutokeyAuth = "autokey"
    SymmetricKeyAuth = "symmetric-key"
)

// Constants for NTP algorithms.
const (
    Sha1 = "sha1"
    Md5 = "md5"
)

// Config is a normalized, version independent representation of a device's
// general settings.
type Config struct {
    Hostname string
    IpAddress string
    Netmask string
    Gateway string
    Timezone string
    Domain string
    UpdateServer string
    VerifyUpdateServer bool
    LoginBanner string
    PanoramaPrimary string
    PanoramaSecondary string
    DnsPrimary string
    DnsSecondary string
    NtpPrimaryAddress string
    NtpPrimaryAuthType string
    NtpPrimaryKeyId int
    NtpPrimaryAlgorithm string
    NtpPrimaryAuthKey string
    NtpSecondaryAddress string
    NtpSecondaryAuthType string
    NtpSecondaryKeyId int
    NtpSecondaryAlgorithm string
    NtpSecondaryAuthKey string

    raw map[string] string
}

// Defaults sets params with uninitialized values to their GUI default setting.
//
// The defaults are as follows:
//      * UpdateServer: updates.paloaltonetworks.com
func (o *Config) Defaults() {
    if o.UpdateServer == "" {
        o.UpdateServer = "updates.paloaltonetworks.com"
    }
}

// Merge copies non connectivity variables from source Config `s` to this
// object.  The fields that are not copied are as follows:
//
//      * IpAddress
//      * Netmask
//      * Gateway
func (o *Config) Merge(s Config) {
    if s.Hostname != "" {
        o.Hostname = s.Hostname
    }

    if s.Timezone != "" {
        o.Timezone = s.Timezone
    }

    if s.Domain != "" {
        o.Domain = s.Domain
    }

    if s.UpdateServer != "" {
        o.UpdateServer = s.UpdateServer
    }

    o.VerifyUpdateServer = s.VerifyUpdateServer

    if s.LoginBanner != "" {
        o.LoginBanner = s.LoginBanner
    }

    if s.PanoramaPrimary != "" {
        o.PanoramaPrimary = s.PanoramaPrimary
    }

    if s.PanoramaSecondary != "" {
        o.PanoramaSecondary = s.PanoramaSecondary
    }

    if s.DnsPrimary != "" {
        o.DnsPrimary = s.DnsPrimary
    }

    if s.DnsSecondary != "" {
        o.DnsSecondary = s.DnsSecondary
    }

    if s.NtpPrimaryAddress != "" {
        o.NtpPrimaryAddress = s.NtpPrimaryAddress
    }

    if s.NtpPrimaryAuthType != "" {
        o.NtpPrimaryAuthType = s.NtpPrimaryAuthType
    }

    if s.NtpPrimaryKeyId != 0 {
        o.NtpPrimaryKeyId = s.NtpPrimaryKeyId
    }

    if s.NtpPrimaryAlgorithm != "" {
        o.NtpPrimaryAlgorithm = s.NtpPrimaryAlgorithm
    }

    if s.NtpPrimaryAuthKey != "" {
        o.NtpPrimaryAuthKey = s.NtpPrimaryAuthKey
    }

    if s.NtpSecondaryAddress != "" {
        o.NtpSecondaryAddress = s.NtpSecondaryAddress
    }

    if s.NtpSecondaryAuthType != "" {
        o.NtpSecondaryAuthType = s.NtpSecondaryAuthType
    }

    if s.NtpSecondaryKeyId != 0 {
        o.NtpSecondaryKeyId = s.NtpSecondaryKeyId
    }

    if s.NtpSecondaryAlgorithm != "" {
        o.NtpSecondaryAlgorithm = s.NtpSecondaryAlgorithm
    }

    if s.NtpSecondaryAuthKey != "" {
        o.NtpSecondaryAuthKey = s.NtpSecondaryAuthKey
    }
}

// General is a namespace struct, included as part of pango.Client.
type General struct {
    con util.XapiClient
}

// Initialize is invoked by client.Initialize().
func (c *General) Initialize(con util.XapiClient) {
    c.con = con
}

// Show performs SHOW to retrieve the device's general settings.
func (c *General) Show() (Config, error) {
    c.con.LogQuery("(show) general settings")
    return c.details(c.con.Show)
}

// Get performs GET to retrieve the device's general settings.
func (c *General) Get() (Config, error) {
    c.con.LogQuery("(get) general settings")
    return c.details(c.con.Get)
}

// Set performs SET to create / update the device's general settings.
func (c *General) Set(e Config) error {
    var err error
    _, fn := c.versioning()
    c.con.LogAction("(set) general settings")

    path := c.xpath()
    path = path[:len(path) - 1]

    _, err = c.con.Set(path, fn(e), nil, nil)
    return err
}

// Edit performs EDIT to update the device's general settings.
func (c *General) Edit(e Config) error {
    var err error
    _, fn := c.versioning()
    c.con.LogAction("(edit) general settings")

    path := c.xpath()

    _, err = c.con.Edit(path, fn(e), nil, nil)
    return err
}

/** Internal functions for the General struct **/

func (c *General) versioning() (normalizer, func(Config) (interface{})) {
    return &container_v1{}, specify_v1
}

func (c *General) details(fn util.Retriever) (Config, error) {
    path := c.xpath()
    obj, _ := c.versioning()
    if _, err := fn(path, nil, obj); err != nil {
        return Config{}, err
    }
    ans := obj.Normalize()

    return ans, nil
}

func (c *General) xpath() []string {
    return []string{
        "config",
        "devices",
        util.AsEntryXpath([]string{"localhost.localdomain"}),
        "deviceconfig",
        "system",
    }
}

/** Structs / functions for this namespace. **/

type normalizer interface {
    Normalize() Config
}

type container_v1 struct {
    Answer config_v1 `xml:"result>system"`
}

func (o *container_v1) Normalize() Config {
    ans := Config{
        Hostname: o.Answer.Hostname,
        IpAddress: o.Answer.IpAddress,
        Netmask: o.Answer.Netmask,
        Gateway: o.Answer.Gateway,
        Timezone: o.Answer.Timezone,
        Domain: o.Answer.Domain,
        UpdateServer: o.Answer.UpdateServer,
        VerifyUpdateServer: util.AsBool(o.Answer.VerifyUpdateServer),
        LoginBanner: o.Answer.LoginBanner,
        PanoramaPrimary: o.Answer.PanoramaPrimary,
        PanoramaSecondary: o.Answer.PanoramaSecondary,
    }
    if o.Answer.Dns != nil {
        ans.DnsPrimary = o.Answer.Dns.Primary
        ans.DnsSecondary = o.Answer.Dns.Secondary
    }
    if o.Answer.Ntp != nil {
        if o.Answer.Ntp.Primary != nil {
            ans.NtpPrimaryAddress = o.Answer.Ntp.Primary.IpAddress
            switch {
            case o.Answer.Ntp.Primary.Auth.None != nil:
                ans.NtpPrimaryAuthType = NoAuth
            case o.Answer.Ntp.Primary.Auth.Autokey != nil:
                ans.NtpPrimaryAuthType = AutokeyAuth
            case o.Answer.Ntp.Primary.Auth.SymmetricKey != nil:
                ans.NtpPrimaryAuthType = SymmetricKeyAuth
                ans.NtpPrimaryKeyId = o.Answer.Ntp.Primary.Auth.SymmetricKey.KeyId
                switch {
                case o.Answer.Ntp.Primary.Auth.SymmetricKey.Algorithm.Sha1 != nil:
                    ans.NtpPrimaryAlgorithm = Sha1
                    ans.NtpPrimaryAuthKey = o.Answer.Ntp.Primary.Auth.SymmetricKey.Algorithm.Sha1.AuthenticationKey
                case o.Answer.Ntp.Primary.Auth.SymmetricKey.Algorithm.Md5 != nil:
                    ans.NtpPrimaryAlgorithm = Md5
                    ans.NtpPrimaryAuthKey = o.Answer.Ntp.Primary.Auth.SymmetricKey.Algorithm.Md5.AuthenticationKey
                }
            }
        }
        if o.Answer.Ntp.Secondary != nil {
            ans.NtpSecondaryAddress = o.Answer.Ntp.Secondary.IpAddress
            switch {
            case o.Answer.Ntp.Secondary.Auth.None != nil:
                ans.NtpSecondaryAuthType = NoAuth
            case o.Answer.Ntp.Secondary.Auth.Autokey != nil:
                ans.NtpSecondaryAuthType = AutokeyAuth
            case o.Answer.Ntp.Secondary.Auth.SymmetricKey != nil:
                ans.NtpSecondaryAuthType = SymmetricKeyAuth
                ans.NtpSecondaryKeyId = o.Answer.Ntp.Secondary.Auth.SymmetricKey.KeyId
                switch {
                case o.Answer.Ntp.Secondary.Auth.SymmetricKey.Algorithm.Sha1 != nil:
                    ans.NtpSecondaryAlgorithm = Sha1
                    ans.NtpSecondaryAuthKey = o.Answer.Ntp.Secondary.Auth.SymmetricKey.Algorithm.Sha1.AuthenticationKey
                case o.Answer.Ntp.Secondary.Auth.SymmetricKey.Algorithm.Md5 != nil:
                    ans.NtpSecondaryAlgorithm = Md5
                    ans.NtpSecondaryAuthKey = o.Answer.Ntp.Secondary.Auth.SymmetricKey.Algorithm.Md5.AuthenticationKey
                }
            }
        }
    }

    ans.raw = make(map[string] string)
    if o.Answer.AckLoginBanner != nil {
        ans.raw["alb"] = util.CleanRawXml(o.Answer.AckLoginBanner.Text)
    }
    if o.Answer.AuthenticationProfile != nil {
        ans.raw["ap"] = util.CleanRawXml(o.Answer.AuthenticationProfile.Text)
    }
    if o.Answer.CertificateProfile != nil {
        ans.raw["cp"] = util.CleanRawXml(o.Answer.CertificateProfile.Text)
    }
    if o.Answer.DomainLookupUrl != nil {
        ans.raw["dlu"] = util.CleanRawXml(o.Answer.DomainLookupUrl.Text)
    }
    if o.Answer.FqdnForceRefreshTime != nil {
        ans.raw["ffrt"] = util.CleanRawXml(o.Answer.FqdnForceRefreshTime.Text)
    }
    if o.Answer.FqdnRefreshTime != nil {
        ans.raw["frt"] = util.CleanRawXml(o.Answer.FqdnRefreshTime.Text)
    }
    if o.Answer.GeoLocation != nil {
        ans.raw["gl"] = util.CleanRawXml(o.Answer.GeoLocation.Text)
    }
    if o.Answer.HsmSettings != nil {
        ans.raw["hs"] = util.CleanRawXml(o.Answer.HsmSettings.Text)
    }
    if o.Answer.IpAddressLookupUrl != nil {
        ans.raw["ialu"] = util.CleanRawXml(o.Answer.IpAddressLookupUrl.Text)
    }
    if o.Answer.Ipv6Address != nil {
        ans.raw["i6a"] = util.CleanRawXml(o.Answer.Ipv6Address.Text)
    }
    if o.Answer.Ipv6DefaultGateway != nil {
        ans.raw["i6dg"] = util.CleanRawXml(o.Answer.Ipv6DefaultGateway.Text)
    }
    if o.Answer.Locale != nil {
        ans.raw["locale"] = util.CleanRawXml(o.Answer.Locale.Text)
    }
    if o.Answer.LogExportSchedule != nil {
        ans.raw["les"] = util.CleanRawXml(o.Answer.LogExportSchedule.Text)
    }
    if o.Answer.LogLink != nil {
        ans.raw["ll"] = util.CleanRawXml(o.Answer.LogLink.Text)
    }
    if o.Answer.MotdAndBanner != nil {
        ans.raw["mab"] = util.CleanRawXml(o.Answer.MotdAndBanner.Text)
    }
    if o.Answer.Mtu != nil {
        ans.raw["mtu"] = util.CleanRawXml(o.Answer.Mtu.Text)
    }
    if o.Answer.PermittedIp != nil {
        ans.raw["pi"] = util.CleanRawXml(o.Answer.PermittedIp.Text)
    }
    if o.Answer.Route != nil {
        ans.raw["route"] = util.CleanRawXml(o.Answer.Route.Text)
    }
    if o.Answer.SecureProxyPassword != nil {
        ans.raw["sppassword"] = util.CleanRawXml(o.Answer.SecureProxyPassword.Text)
    }
    if o.Answer.SecureProxyPort != nil {
        ans.raw["spport"] = util.CleanRawXml(o.Answer.SecureProxyPort.Text)
    }
    if o.Answer.SecureProxyServer != nil {
        ans.raw["sps"] = util.CleanRawXml(o.Answer.SecureProxyServer.Text)
    }
    if o.Answer.SecureProxyUser != nil {
        ans.raw["spu"] = util.CleanRawXml(o.Answer.SecureProxyUser.Text)
    }
    if o.Answer.Service != nil {
        ans.raw["service"] = util.CleanRawXml(o.Answer.Service.Text)
    }
    if o.Answer.SnmpSetting != nil {
        ans.raw["ss"] = util.CleanRawXml(o.Answer.SnmpSetting.Text)
    }
    if o.Answer.SpeedDuplex != nil {
        ans.raw["sd"] = util.CleanRawXml(o.Answer.SpeedDuplex.Text)
    }
    if o.Answer.SslTlsServiceProfile != nil {
        ans.raw["stsp"] = util.CleanRawXml(o.Answer.SslTlsServiceProfile.Text)
    }
    if o.Answer.SyslogCertificate != nil {
        ans.raw["sc"] = util.CleanRawXml(o.Answer.SyslogCertificate.Text)
    }
    if o.Answer.Type != nil {
        ans.raw["type"] = util.CleanRawXml(o.Answer.Type.Text)
    }
    if o.Answer.UpdateSchedule != nil {
        ans.raw["us"] = util.CleanRawXml(o.Answer.UpdateSchedule.Text)
    }
    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type config_v1 struct {
    XMLName xml.Name `xml:"system"`
    Hostname string `xml:"hostname"`
    IpAddress string `xml:"ip-address,omitempty"`
    Netmask string `xml:"netmask,omitempty"`
    Gateway string `xml:"default-gateway,omitempty"`
    Timezone string `xml:"timezone"`
    Domain string `xml:"domain,omitempty"`
    UpdateServer string `xml:"update-server,omitempty"`
    VerifyUpdateServer string `xml:"server-verification"`
    LoginBanner string `xml:"login-banner,omitempty"`
    PanoramaPrimary string `xml:"panorama-server,omitempty"`
    PanoramaSecondary string `xml:"panorama-server-2,omitempty"`
    Dns *deviceDns `xml:"dns-setting"`
    Ntp *deviceNtp `xml:"ntp-servers"`
    AckLoginBanner *util.RawXml `xml:"ack-login-banner"`
    AuthenticationProfile *util.RawXml `xml:"authentication-profile"`
    CertificateProfile *util.RawXml `xml:"certificate-profile"`
    DomainLookupUrl *util.RawXml `xml:"domain-lookup-url"`
    FqdnForceRefreshTime *util.RawXml `xml:"fqdn-forcerefresh-time"`
    FqdnRefreshTime *util.RawXml `xml:"fqdn-refresh-time"`
    GeoLocation *util.RawXml `xml:"geo-location"`
    HsmSettings *util.RawXml `xml:"hsm-settings"`
    IpAddressLookupUrl *util.RawXml `xml:"ip-address-lookup-url"`
    Ipv6Address *util.RawXml `xml:"ipv6-address"`
    Ipv6DefaultGateway *util.RawXml `xml:"ipv6-default-gateway"`
    Locale *util.RawXml `xml:"locale"`
    LogExportSchedule *util.RawXml `xml:"log-export-schedule"`
    LogLink *util.RawXml `xml:"log-link"`
    MotdAndBanner *util.RawXml `xml:"motd-and-banner"`
    Mtu *util.RawXml `xml:"mtu"`
    PermittedIp *util.RawXml `xml:"permitted-ip"`
    Route *util.RawXml `xml:"route"`
    SecureProxyPassword *util.RawXml `xml:"secure-proxy-password"`
    SecureProxyPort *util.RawXml `xml:"secure-proxy-port"`
    SecureProxyServer *util.RawXml `xml:"secure-proxy-server"`
    SecureProxyUser *util.RawXml `xml:"secure-proxy-user"`
    Service *util.RawXml `xml:"service"`
    SnmpSetting *util.RawXml `xml:"snmp-setting"`
    SpeedDuplex *util.RawXml `xml:"speed-duplex"`
    SslTlsServiceProfile *util.RawXml `xml:"ssl-tls-service-profile"`
    SyslogCertificate *util.RawXml `xml:"syslog-certificate"`
    Type *util.RawXml `xml:"type"`
    UpdateSchedule *util.RawXml `xml:"update-schedule"`
}

type deviceDns struct {
    Primary string `xml:"servers>primary,omitempty"`
    Secondary string `xml:"servers>secondary,omitempty"`
}

type deviceNtp struct {
    Primary *ntpConfig `xml:"primary-ntp-server"`
    Secondary *ntpConfig `xml:"secondary-ntp-server"`
}

type ntpConfig struct {
    IpAddress string `xml:"ntp-server-address"`
    Auth ntpAuth `xml:"authentication-type"`
}

type ntpAuth struct {
    None *string `xml:"none"`
    Autokey *string `xml:"autokey"`
    SymmetricKey *symKey `xml:"symmetric-key"`
}

type symKey struct {
    KeyId int `xml:"key-id"`
    Algorithm symKeyAlgorithm `xml:"algorithm"`
}

type symKeyAlgorithm struct {
    Sha1 *algorithmAuthKey `xml:"sha1"`
    Md5 *algorithmAuthKey `xml:"md5"`
}

type algorithmAuthKey struct {
    AuthenticationKey string `xml:"authentication-key"`
}

func specify_v1(c Config) interface{} {
    ans := config_v1{
        Hostname: c.Hostname,
        IpAddress: c.IpAddress,
        Netmask: c.Netmask,
        Gateway: c.Gateway,
        Timezone: c.Timezone,
        Domain: c.Domain,
        UpdateServer: c.UpdateServer,
        VerifyUpdateServer: util.YesNo(c.VerifyUpdateServer),
        LoginBanner: c.LoginBanner,
        PanoramaPrimary: c.PanoramaPrimary,
        PanoramaSecondary: c.PanoramaSecondary,
    }
    if c.DnsPrimary != "" || c.DnsSecondary != "" {
        ans.Dns = &deviceDns{
            c.DnsPrimary,
            c.DnsSecondary,
        }
    }
    // NTP
    ntp_config := &deviceNtp{}
    if c.NtpPrimaryAddress != "" || c.NtpPrimaryAuthType != "" {
        ntp_config.Primary = &ntpConfig{
            IpAddress: c.NtpPrimaryAddress,
        }
        var es string
        switch c.NtpPrimaryAuthType {
        case NoAuth:
            ntp_config.Primary.Auth.None = &es
        case AutokeyAuth:
            ntp_config.Primary.Auth.Autokey = &es
        case SymmetricKeyAuth:
            ntp_config.Primary.Auth.SymmetricKey = &symKey{
                KeyId: c.NtpPrimaryKeyId,
            }
            switch c.NtpPrimaryAlgorithm {
            case Sha1:
                ntp_config.Primary.Auth.SymmetricKey.Algorithm.Sha1 = &algorithmAuthKey{c.NtpPrimaryAuthKey}
            case Md5:
                ntp_config.Primary.Auth.SymmetricKey.Algorithm.Md5 = &algorithmAuthKey{c.NtpPrimaryAuthKey}
            }
        }
    }
    if c.NtpSecondaryAddress != "" || c.NtpSecondaryAuthType != "" {
        ntp_config.Secondary = &ntpConfig{
            IpAddress: c.NtpSecondaryAddress,
        }
        var es string
        switch c.NtpSecondaryAuthType {
        case NoAuth:
            ntp_config.Secondary.Auth.None = &es
        case AutokeyAuth:
            ntp_config.Secondary.Auth.Autokey = &es
        case SymmetricKeyAuth:
            ntp_config.Secondary.Auth.SymmetricKey = &symKey{
                KeyId: c.NtpSecondaryKeyId,
            }
            switch c.NtpSecondaryAlgorithm {
            case Sha1:
                ntp_config.Secondary.Auth.SymmetricKey.Algorithm.Sha1 = &algorithmAuthKey{c.NtpSecondaryAuthKey}
            case Md5:
                ntp_config.Secondary.Auth.SymmetricKey.Algorithm.Md5 = &algorithmAuthKey{c.NtpSecondaryAuthKey}
            }
        }
    }
    if ntp_config.Primary != nil || ntp_config.Secondary != nil {
        ans.Ntp = ntp_config
    }

    if text, present := c.raw["alb"]; present {
        ans.AckLoginBanner = &util.RawXml{text}
    }
    if text, present := c.raw["ap"]; present {
        ans.AuthenticationProfile = &util.RawXml{text}
    }
    if text, present := c.raw["cp"]; present {
        ans.CertificateProfile = &util.RawXml{text}
    }
    if text, present := c.raw["dlu"]; present {
        ans.DomainLookupUrl = &util.RawXml{text}
    }
    if text, present := c.raw["ffrt"]; present {
        ans.FqdnForceRefreshTime = &util.RawXml{text}
    }
    if text, present := c.raw["frt"]; present {
        ans.FqdnRefreshTime = &util.RawXml{text}
    }
    if text, present := c.raw["gl"]; present {
        ans.GeoLocation = &util.RawXml{text}
    }
    if text, present := c.raw["hs"]; present {
        ans.HsmSettings = &util.RawXml{text}
    }
    if text, present := c.raw["ialu"]; present {
        ans.IpAddressLookupUrl = &util.RawXml{text}
    }
    if text, present := c.raw["i6a"]; present {
        ans.Ipv6Address = &util.RawXml{text}
    }
    if text, present := c.raw["i6dg"]; present {
        ans.Ipv6DefaultGateway = &util.RawXml{text}
    }
    if text, present := c.raw["locale"]; present {
        ans.Locale = &util.RawXml{text}
    }
    if text, present := c.raw["les"]; present {
        ans.LogExportSchedule = &util.RawXml{text}
    }
    if text, present := c.raw["ll"]; present {
        ans.LogLink = &util.RawXml{text}
    }
    if text, present := c.raw["mab"]; present {
        ans.MotdAndBanner = &util.RawXml{text}
    }
    if text, present := c.raw["mtu"]; present {
        ans.Mtu = &util.RawXml{text}
    }
    if text, present := c.raw["pi"]; present {
        ans.PermittedIp = &util.RawXml{text}
    }
    if text, present := c.raw["route"]; present {
        ans.Route = &util.RawXml{text}
    }
    if text, present := c.raw["sppassword"]; present {
        ans.SecureProxyPassword = &util.RawXml{text}
    }
    if text, present := c.raw["spport"]; present {
        ans.SecureProxyPort = &util.RawXml{text}
    }
    if text, present := c.raw["sps"]; present {
        ans.SecureProxyServer = &util.RawXml{text}
    }
    if text, present := c.raw["spu"]; present {
        ans.SecureProxyUser = &util.RawXml{text}
    }
    if text, present := c.raw["service"]; present {
        ans.Service = &util.RawXml{text}
    }
    if text, present := c.raw["ss"]; present {
        ans.SnmpSetting = &util.RawXml{text}
    }
    if text, present := c.raw["sd"]; present {
        ans.SpeedDuplex = &util.RawXml{text}
    }
    if text, present := c.raw["stsp"]; present {
        ans.SslTlsServiceProfile = &util.RawXml{text}
    }
    if text, present := c.raw["sc"]; present {
        ans.SyslogCertificate = &util.RawXml{text}
    }
    if text, present := c.raw["type"]; present {
        ans.Type = &util.RawXml{text}
    }
    if text, present := c.raw["us"]; present {
        ans.UpdateSchedule = &util.RawXml{text}
    }

    return ans
}
