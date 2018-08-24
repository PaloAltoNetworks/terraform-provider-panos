package ipsectunnel

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


const (
    TypeAutoKey = "auto-key"
    TypeManualKey = "manual-key"
    TypeGlobalProtectSatellite = "global-protect-satellite"
)

const (
    MkEspEncryptionDes = "des"
    MkEspEncryption3des = "3des"
    MkEspEncryptionAes128 = "aes-128-cbc"
    MkEspEncryptionAes192 = "aes-192-cbc"
    MkEspEncryptionAes256 = "aes-256-cbc"
    MkEspEncryptionNull = "null"
)

const (
    MkProtocolEsp = "esp"
    MkProtocolAh = "ah"
)

const (
    MkAuthTypeMd5 = "md5"
    MkAuthTypeSha1 = "sha1"
    MkAuthTypeSha256 = "sha256"
    MkAuthTypeSha384 = "sha384"
    MkAuthTypeSha512 = "sha512"
    MkAuthTypeNone = "none"
)

// Entry is a normalized, version independent representation of an IKE gateway.
type Entry struct {
    Name string
    TunnelInterface string
    AntiReplay bool
    EnableIpv6 bool
    Type string
    AkIkeGateway string
    AkIpsecCryptoProfile string
    MkLocalSpi string
    MkInterface string
    MkRemoteSpi string
    MkRemoteAddress string
    MkLocalAddressIp string
    MkLocalAddressFloatingIp string
    MkProtocol string
    MkAuthType string
    MkAuthKey string
    MkEspEncryptionType string
    MkEspEncryptionKey string
    GpsInterface string
    GpsPortalAddress string
    GpsPreferIpv6 bool
    GpsInterfaceIpIpv4 string
    GpsInterfaceIpIpv6 string
    GpsInterfaceFloatingIpIpv4 string
    GpsInterfaceFloatingIpIpv6 string
    GpsPublishConnectedRoutes bool
    GpsPublishRoutes []string
    GpsLocalCertificate string
    GpsCertificateProfile string
    CopyTos bool
    CopyFlowLabel bool
    EnableTunnelMonitor bool
    TunnelMonitorDestinationIp string
    TunnelMonitorSourceIp string
    TunnelMonitorProxyId string
    TunnelMonitorProfile string
    Disabled bool

    raw map[string] string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.TunnelInterface = s.TunnelInterface
    o.AntiReplay = s.AntiReplay
    o.EnableIpv6 = s.EnableIpv6
    o.Type = s.Type
    o.AkIkeGateway = s.AkIkeGateway
    o.AkIpsecCryptoProfile = s.AkIpsecCryptoProfile
    o.MkLocalSpi = s.MkLocalSpi
    o.MkInterface = s.MkInterface
    o.MkRemoteSpi = s.MkRemoteSpi
    o.MkRemoteAddress = s.MkRemoteAddress
    o.MkLocalAddressIp = s.MkLocalAddressIp
    o.MkLocalAddressFloatingIp = s.MkLocalAddressFloatingIp
    o.MkProtocol = s.MkProtocol
    o.MkAuthType = s.MkAuthType
    o.MkAuthKey = s.MkAuthKey
    o.MkEspEncryptionType = s.MkEspEncryptionType
    o.MkEspEncryptionKey = s.MkEspEncryptionKey
    o.GpsInterface = s.GpsInterface
    o.GpsPreferIpv6 = s.GpsPreferIpv6
    o.GpsInterfaceIpIpv4 = s.GpsInterfaceIpIpv4
    o.GpsInterfaceIpIpv6 = s.GpsInterfaceIpIpv6
    o.GpsInterfaceFloatingIpIpv4 = s.GpsInterfaceFloatingIpIpv4
    o.GpsInterfaceFloatingIpIpv6 = s.GpsInterfaceFloatingIpIpv6
    o.GpsPublishConnectedRoutes = s.GpsPublishConnectedRoutes
    o.GpsPublishRoutes = s.GpsPublishRoutes
    o.GpsLocalCertificate = s.GpsLocalCertificate
    o.GpsCertificateProfile = s.GpsCertificateProfile
    o.AntiReplay = s.AntiReplay
    o.CopyTos = s.CopyTos
    o.CopyFlowLabel = s.CopyFlowLabel
    o.EnableTunnelMonitor = s.EnableTunnelMonitor
    o.TunnelMonitorDestinationIp = s.TunnelMonitorDestinationIp
    o.TunnelMonitorSourceIp = s.TunnelMonitorSourceIp
    o.TunnelMonitorProxyId = s.TunnelMonitorProxyId
    o.TunnelMonitorProfile = s.TunnelMonitorProfile
    o.Disabled = s.Disabled
}

// SpecifyEncryption takes normalized encryption values and changes them to the
// version specific values PAN-OS will be expecting.
//
// Param v should be 1 if you're running against PAN-OS 6.1, 2 if you're
// running against 7.0 or later.
func (o *Entry) SpecifyEncryption(v int) {
    switch v {
    case 2:
        switch o.MkEspEncryptionType {
        case MkEspEncryptionDes:
            o.MkEspEncryptionType = "des"
        case MkEspEncryption3des:
            o.MkEspEncryptionType = "3des"
        case MkEspEncryptionAes128:
            o.MkEspEncryptionType = "aes-128-cbc"
        case MkEspEncryptionAes192:
            o.MkEspEncryptionType = "aes-192-cbc"
        case MkEspEncryptionAes256:
            o.MkEspEncryptionType = "aes-256-cbc"
        case MkEspEncryptionNull:
            o.MkEspEncryptionType = "null"
        }
    case 1:
        switch o.MkEspEncryptionType {
        case MkEspEncryptionDes:
            o.MkEspEncryptionType = "des"
        case MkEspEncryption3des:
            o.MkEspEncryptionType = "3des"
        case MkEspEncryptionAes128:
            o.MkEspEncryptionType = "aes128"
        case MkEspEncryptionAes192:
            o.MkEspEncryptionType = "aes192"
        case MkEspEncryptionAes256:
            o.MkEspEncryptionType = "aes256"
        case MkEspEncryptionNull:
            o.MkEspEncryptionType = "null"
        }
    }
}

// NormalizeEncryption normalizes the fields in o.MkEspEncryption.
func (o *Entry) NormalizeEncryption() {
    switch o.MkEspEncryptionType {
    case "des":
        o.MkEspEncryptionType = MkEspEncryptionDes
    case "3des":
        o.MkEspEncryptionType = MkEspEncryption3des
    case "aes-128-cbc", "aes128":
        o.MkEspEncryptionType = MkEspEncryptionAes128
    case "aes-192-cbc", "aes192":
        o.MkEspEncryptionType = MkEspEncryptionAes192
    case "aes-256-cbc", "aes256":
        o.MkEspEncryptionType = MkEspEncryptionAes256
    case "null":
        o.MkEspEncryptionType = MkEspEncryptionNull
    }
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
        TunnelInterface: o.Answer.TunnelInterface,
        AntiReplay: util.AsBool(o.Answer.AntiReplay),
        CopyTos: util.AsBool(o.Answer.CopyTos),
    }

    ans.raw = make(map[string] string)

    if o.Answer.Ak != nil {
        ans.Type = TypeAutoKey
        ans.AkIkeGateway = util.EntToOneStr(o.Answer.Ak.AkIkeGateway)
        ans.AkIpsecCryptoProfile = o.Answer.Ak.AkIpsecCryptoProfile

        if o.Answer.Ak.ProxyIpv4 != nil {
            ans.raw["pv4"] = util.CleanRawXml(o.Answer.Ak.ProxyIpv4.Text)
        }

        if o.Answer.Ak.ProxyIpv6 != nil {
            ans.raw["pv6"] = util.CleanRawXml(o.Answer.Ak.ProxyIpv6.Text)
        }
    } else if o.Answer.Mk != nil {
        ans.Type = TypeManualKey
        ans.MkLocalSpi = o.Answer.Mk.MkLocalSpi
        ans.MkInterface = o.Answer.Mk.Local.MkInterface
        ans.MkLocalAddressIp = o.Answer.Mk.Local.MkLocalAddressIp
        ans.MkRemoteAddress = o.Answer.Mk.Peer.MkRemoteAddress
        ans.MkRemoteSpi = o.Answer.Mk.MkRemoteSpi

        if o.Answer.Mk.Esp != nil {
            ans.MkProtocol = MkProtocolEsp
            if o.Answer.Mk.Esp.AuthMd5 != nil {
                ans.MkAuthType = MkAuthTypeMd5
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthMd5.Value
            } else if o.Answer.Mk.Esp.AuthSha1 != nil {
                ans.MkAuthType = MkAuthTypeSha1
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha1.Value
            } else if o.Answer.Mk.Esp.AuthSha256 != nil {
                ans.MkAuthType = MkAuthTypeSha256
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha256.Value
            } else if o.Answer.Mk.Esp.AuthSha384 != nil {
                ans.MkAuthType = MkAuthTypeSha384
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha384.Value
            } else if o.Answer.Mk.Esp.AuthSha512 != nil {
                ans.MkAuthType = MkAuthTypeSha512
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha512.Value
            } else if o.Answer.Mk.Esp.AuthNone != nil {
                ans.MkAuthType = MkAuthTypeNone
            }

            ans.MkEspEncryptionType = o.Answer.Mk.Esp.MkEspEncryptionType
            ans.MkEspEncryptionKey = o.Answer.Mk.Esp.MkEspEncryptionKey
        } else if o.Answer.Mk.Ah != nil {
            ans.MkProtocol = MkProtocolAh
            if o.Answer.Mk.Ah.AuthMd5 != nil {
                ans.MkAuthType = MkAuthTypeMd5
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthMd5.Value
            } else if o.Answer.Mk.Ah.AuthSha1 != nil {
                ans.MkAuthType = MkAuthTypeSha1
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha1.Value
            } else if o.Answer.Mk.Ah.AuthSha256 != nil {
                ans.MkAuthType = MkAuthTypeSha256
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha256.Value
            } else if o.Answer.Mk.Ah.AuthSha384 != nil {
                ans.MkAuthType = MkAuthTypeSha384
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha384.Value
            } else if o.Answer.Mk.Ah.AuthSha512 != nil {
                ans.MkAuthType = MkAuthTypeSha512
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha512.Value
            }
        }
    } else if o.Answer.Gps != nil {
        ans.Type = TypeGlobalProtectSatellite
        ans.GpsPortalAddress = o.Answer.Gps.GpsPortalAddress
        ans.GpsPublishRoutes = util.MemToStr(o.Answer.Gps.GpsPublishRoutes)
        ans.GpsInterface = o.Answer.Gps.Local.GpsInterface
        ans.GpsInterfaceIpIpv4 = o.Answer.Gps.Local.GpsInterfaceIpIpv4

        if o.Answer.Gps.Pcr != nil {
            ans.GpsPublishConnectedRoutes = util.AsBool(o.Answer.Gps.Pcr.GpsPublishConnectedRoutes)
        }

        if o.Answer.Gps.Ca != nil {
            ans.GpsLocalCertificate = o.Answer.Gps.Ca.GpsLocalCertificate
            ans.GpsCertificateProfile = o.Answer.Gps.Ca.GpsCertificateProfile
        }
    }

    if o.Answer.TunnelMonitor != nil {
        ans.EnableTunnelMonitor = util.AsBool(o.Answer.TunnelMonitor.EnableTunnelMonitor)
        ans.TunnelMonitorDestinationIp = o.Answer.TunnelMonitor.TunnelMonitorDestinationIp
        ans.TunnelMonitorSourceIp = o.Answer.TunnelMonitor.TunnelMonitorSourceIp
        ans.TunnelMonitorProfile = o.Answer.TunnelMonitor.TunnelMonitorProfile
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        TunnelInterface: o.Answer.TunnelInterface,
        AntiReplay: util.AsBool(o.Answer.AntiReplay),
        CopyTos: util.AsBool(o.Answer.CopyTos),
        EnableIpv6: util.AsBool(o.Answer.EnableIpv6),
        Disabled: util.AsBool(o.Answer.Disabled),
        CopyFlowLabel: util.AsBool(o.Answer.CopyFlowLabel),
    }

    ans.raw = make(map[string] string)

    if o.Answer.Ak != nil {
        ans.Type = TypeAutoKey
        ans.AkIkeGateway = util.EntToOneStr(o.Answer.Ak.AkIkeGateway)
        ans.AkIpsecCryptoProfile = o.Answer.Ak.AkIpsecCryptoProfile

        if o.Answer.Ak.ProxyIpv4 != nil {
            ans.raw["pv4"] = util.CleanRawXml(o.Answer.Ak.ProxyIpv4.Text)
        }

        if o.Answer.Ak.ProxyIpv6 != nil {
            ans.raw["pv6"] = util.CleanRawXml(o.Answer.Ak.ProxyIpv6.Text)
        }
    } else if o.Answer.Mk != nil {
        ans.Type = TypeManualKey
        ans.MkLocalSpi = o.Answer.Mk.MkLocalSpi
        ans.MkInterface = o.Answer.Mk.Local.MkInterface
        ans.MkLocalAddressIp = o.Answer.Mk.Local.MkLocalAddressIp
        ans.MkLocalAddressFloatingIp = o.Answer.Mk.Local.MkLocalAddressFloatingIp
        ans.MkRemoteAddress = o.Answer.Mk.Peer.MkRemoteAddress
        ans.MkRemoteSpi = o.Answer.Mk.MkRemoteSpi

        if o.Answer.Mk.Esp != nil {
            ans.MkProtocol = MkProtocolEsp
            if o.Answer.Mk.Esp.AuthMd5 != nil {
                ans.MkAuthType = MkAuthTypeMd5
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthMd5.Value
            } else if o.Answer.Mk.Esp.AuthSha1 != nil {
                ans.MkAuthType = MkAuthTypeSha1
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha1.Value
            } else if o.Answer.Mk.Esp.AuthSha256 != nil {
                ans.MkAuthType = MkAuthTypeSha256
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha256.Value
            } else if o.Answer.Mk.Esp.AuthSha384 != nil {
                ans.MkAuthType = MkAuthTypeSha384
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha384.Value
            } else if o.Answer.Mk.Esp.AuthSha512 != nil {
                ans.MkAuthType = MkAuthTypeSha512
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha512.Value
            } else if o.Answer.Mk.Esp.AuthNone != nil {
                ans.MkAuthType = MkAuthTypeNone
            }

            ans.MkEspEncryptionType = o.Answer.Mk.Esp.MkEspEncryptionType
            ans.MkEspEncryptionKey = o.Answer.Mk.Esp.MkEspEncryptionKey
        } else if o.Answer.Mk.Ah != nil {
            ans.MkProtocol = MkProtocolAh
            if o.Answer.Mk.Ah.AuthMd5 != nil {
                ans.MkAuthType = MkAuthTypeMd5
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthMd5.Value
            } else if o.Answer.Mk.Ah.AuthSha1 != nil {
                ans.MkAuthType = MkAuthTypeSha1
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha1.Value
            } else if o.Answer.Mk.Ah.AuthSha256 != nil {
                ans.MkAuthType = MkAuthTypeSha256
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha256.Value
            } else if o.Answer.Mk.Ah.AuthSha384 != nil {
                ans.MkAuthType = MkAuthTypeSha384
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha384.Value
            } else if o.Answer.Mk.Ah.AuthSha512 != nil {
                ans.MkAuthType = MkAuthTypeSha512
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha512.Value
            }
        }
    } else if o.Answer.Gps != nil {
        ans.Type = TypeGlobalProtectSatellite
        ans.GpsPortalAddress = o.Answer.Gps.GpsPortalAddress
        ans.GpsPublishRoutes = util.MemToStr(o.Answer.Gps.GpsPublishRoutes)
        ans.GpsInterface = o.Answer.Gps.Local.GpsInterface
        ans.GpsInterfaceIpIpv4 = o.Answer.Gps.Local.GpsInterfaceIpIpv4
        ans.GpsInterfaceFloatingIpIpv4 = o.Answer.Gps.Local.GpsInterfaceFloatingIpIpv4

        if o.Answer.Gps.Pcr != nil {
            ans.GpsPublishConnectedRoutes = util.AsBool(o.Answer.Gps.Pcr.GpsPublishConnectedRoutes)
        }

        if o.Answer.Gps.Ca != nil {
            ans.GpsLocalCertificate = o.Answer.Gps.Ca.GpsLocalCertificate
            ans.GpsCertificateProfile = o.Answer.Gps.Ca.GpsCertificateProfile
        }
    }

    if o.Answer.TunnelMonitor != nil {
        ans.EnableTunnelMonitor = util.AsBool(o.Answer.TunnelMonitor.EnableTunnelMonitor)
        ans.TunnelMonitorDestinationIp = o.Answer.TunnelMonitor.TunnelMonitorDestinationIp
        ans.TunnelMonitorSourceIp = o.Answer.TunnelMonitor.TunnelMonitorSourceIp
        ans.TunnelMonitorProfile = o.Answer.TunnelMonitor.TunnelMonitorProfile
        ans.TunnelMonitorProxyId = o.Answer.TunnelMonitor.TunnelMonitorProxyId
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type container_v3 struct {
    Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        TunnelInterface: o.Answer.TunnelInterface,
        AntiReplay: util.AsBool(o.Answer.AntiReplay),
        CopyTos: util.AsBool(o.Answer.CopyTos),
        EnableIpv6: util.AsBool(o.Answer.EnableIpv6),
        Disabled: util.AsBool(o.Answer.Disabled),
        CopyFlowLabel: util.AsBool(o.Answer.CopyFlowLabel),
    }

    ans.raw = make(map[string] string)

    if o.Answer.Ak != nil {
        ans.Type = TypeAutoKey
        ans.AkIkeGateway = util.EntToOneStr(o.Answer.Ak.AkIkeGateway)
        ans.AkIpsecCryptoProfile = o.Answer.Ak.AkIpsecCryptoProfile

        if o.Answer.Ak.ProxyIpv4 != nil {
            ans.raw["pv4"] = util.CleanRawXml(o.Answer.Ak.ProxyIpv4.Text)
        }

        if o.Answer.Ak.ProxyIpv6 != nil {
            ans.raw["pv6"] = util.CleanRawXml(o.Answer.Ak.ProxyIpv6.Text)
        }
    } else if o.Answer.Mk != nil {
        ans.Type = TypeManualKey
        ans.MkLocalSpi = o.Answer.Mk.MkLocalSpi
        ans.MkInterface = o.Answer.Mk.Local.MkInterface
        ans.MkLocalAddressIp = o.Answer.Mk.Local.MkLocalAddressIp
        ans.MkLocalAddressFloatingIp = o.Answer.Mk.Local.MkLocalAddressFloatingIp
        ans.MkRemoteAddress = o.Answer.Mk.Peer.MkRemoteAddress
        ans.MkRemoteSpi = o.Answer.Mk.MkRemoteSpi

        if o.Answer.Mk.Esp != nil {
            ans.MkProtocol = MkProtocolEsp
            if o.Answer.Mk.Esp.AuthMd5 != nil {
                ans.MkAuthType = MkAuthTypeMd5
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthMd5.Value
            } else if o.Answer.Mk.Esp.AuthSha1 != nil {
                ans.MkAuthType = MkAuthTypeSha1
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha1.Value
            } else if o.Answer.Mk.Esp.AuthSha256 != nil {
                ans.MkAuthType = MkAuthTypeSha256
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha256.Value
            } else if o.Answer.Mk.Esp.AuthSha384 != nil {
                ans.MkAuthType = MkAuthTypeSha384
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha384.Value
            } else if o.Answer.Mk.Esp.AuthSha512 != nil {
                ans.MkAuthType = MkAuthTypeSha512
                ans.MkAuthKey = o.Answer.Mk.Esp.AuthSha512.Value
            } else if o.Answer.Mk.Esp.AuthNone != nil {
                ans.MkAuthType = MkAuthTypeNone
            }

            ans.MkEspEncryptionType = o.Answer.Mk.Esp.MkEspEncryptionType
            ans.MkEspEncryptionKey = o.Answer.Mk.Esp.MkEspEncryptionKey
        } else if o.Answer.Mk.Ah != nil {
            ans.MkProtocol = MkProtocolAh
            if o.Answer.Mk.Ah.AuthMd5 != nil {
                ans.MkAuthType = MkAuthTypeMd5
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthMd5.Value
            } else if o.Answer.Mk.Ah.AuthSha1 != nil {
                ans.MkAuthType = MkAuthTypeSha1
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha1.Value
            } else if o.Answer.Mk.Ah.AuthSha256 != nil {
                ans.MkAuthType = MkAuthTypeSha256
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha256.Value
            } else if o.Answer.Mk.Ah.AuthSha384 != nil {
                ans.MkAuthType = MkAuthTypeSha384
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha384.Value
            } else if o.Answer.Mk.Ah.AuthSha512 != nil {
                ans.MkAuthType = MkAuthTypeSha512
                ans.MkAuthKey = o.Answer.Mk.Ah.AuthSha512.Value
            }
        }
    } else if o.Answer.Gps != nil {
        ans.Type = TypeGlobalProtectSatellite
        ans.GpsPortalAddress = o.Answer.Gps.GpsPortalAddress
        ans.GpsPreferIpv6 = util.AsBool(o.Answer.Gps.GpsPreferIpv6)
        ans.GpsPublishRoutes = util.MemToStr(o.Answer.Gps.GpsPublishRoutes)
        ans.GpsInterface = o.Answer.Gps.Local.GpsInterface

        if o.Answer.Gps.Pcr != nil {
            ans.GpsPublishConnectedRoutes = util.AsBool(o.Answer.Gps.Pcr.GpsPublishConnectedRoutes)
        }

        if o.Answer.Gps.Local.Ip != nil {
            ans.GpsInterfaceIpIpv4 = o.Answer.Gps.Local.Ip.Ipv4
            ans.GpsInterfaceIpIpv6 = o.Answer.Gps.Local.Ip.Ipv6
        }

        if o.Answer.Gps.Local.Floating != nil {
            ans.GpsInterfaceFloatingIpIpv4 = o.Answer.Gps.Local.Floating.Ipv4
            ans.GpsInterfaceFloatingIpIpv6 = o.Answer.Gps.Local.Floating.Ipv6
        }

        if o.Answer.Gps.Ca != nil {
            ans.GpsLocalCertificate = o.Answer.Gps.Ca.GpsLocalCertificate
            ans.GpsCertificateProfile = o.Answer.Gps.Ca.GpsCertificateProfile
        }
    }

    if o.Answer.TunnelMonitor != nil {
        ans.EnableTunnelMonitor = util.AsBool(o.Answer.TunnelMonitor.EnableTunnelMonitor)
        ans.TunnelMonitorDestinationIp = o.Answer.TunnelMonitor.TunnelMonitorDestinationIp
        ans.TunnelMonitorSourceIp = o.Answer.TunnelMonitor.TunnelMonitorSourceIp
        ans.TunnelMonitorProfile = o.Answer.TunnelMonitor.TunnelMonitorProfile
        ans.TunnelMonitorProxyId = o.Answer.TunnelMonitor.TunnelMonitorProxyId
    }

    if len(ans.raw) == 0 {
        ans.raw = nil
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TunnelInterface string `xml:"tunnel-interface"`
    AntiReplay string `xml:"anti-replay,omitempty"`
    Ak *ak `xml:"auto-key"`
    Mk *mk_v1 `xml:"manual-key"`
    Gps *gps_v1 `xml:"global-protect-satellite"`
    TunnelMonitor *tunMon_v1 `xml:"tunnel-monitor"`
    CopyTos string `xml:"copy-tos"`
}

type ak struct {
    AkIkeGateway *util.EntryType `xml:"ike-gateway"`
    AkIpsecCryptoProfile string `xml:"ipsec-crypto-profile,omitempty"`
    ProxyIpv4 *util.RawXml `xml:"proxy-id"`
    ProxyIpv6 *util.RawXml `xml:"proxy-id-v6"`
}

type mk_v1 struct {
    MkLocalSpi string `xml:"local-spi"`
    Local mkLocal_v1 `xml:"local-address"`
    Peer mkPeer `xml:"peer-address"`
    MkRemoteSpi string `xml:"remote-spi"`
    Esp *mkEsp `xml:"esp"`
    Ah *mkAh `xml:"ah"`
}

type mkLocal_v1 struct {
    MkInterface string `xml:"interface"`
    MkLocalAddressIp string `xml:"ip,omitempty"`
}

type mkPeer struct {
    MkRemoteAddress string `xml:"ip"`
}

type mkEsp struct {
    AuthMd5 *authKey `xml:"authentication>md5"`
    AuthSha1 *authKey `xml:"authentication>sha1"`
    AuthSha256 *authKey `xml:"authentication>sha256"`
    AuthSha384 *authKey `xml:"authentication>sha384"`
    AuthSha512 *authKey `xml:"authentication>sha512"`
    AuthNone *string `xml:"authentication>none"`
    MkEspEncryptionType string `xml:"encryption>algorithm"`
    MkEspEncryptionKey string `xml:"encryption>key,omitempty"`
}

type authKey struct {
    Value string `xml:"key"`
}

type mkAh struct {
    AuthMd5 *authKey `xml:"md5"`
    AuthSha1 *authKey `xml:"sha1"`
    AuthSha256 *authKey `xml:"sha256"`
    AuthSha384 *authKey `xml:"sha384"`
    AuthSha512 *authKey `xml:"sha512"`
}

type gps_v1 struct {
    GpsPortalAddress string `xml:"portal-address"`
    Local gpsLocal_v1 `xml:"local-address"`
    Pcr *gpsPcr `xml:"publish-connected-routes"`
    GpsPublishRoutes *util.MemberType `xml:"publish-routes"`
    Ca *gpsCa `xml:"external-ca"`
}

type gpsLocal_v1 struct {
    GpsInterface string `xml:"interface"`
    GpsInterfaceIpIpv4 string `xml:"ip,omitempty"`
}

type gpsPcr struct {
    GpsPublishConnectedRoutes string `xml:"enable"`
}

type gpsCa struct {
    GpsLocalCertificate string `xml:"local-certificate,omitempty"`
    GpsCertificateProfile string `xml:"certificate-profile,omitempty"`
}

type tunMon_v1 struct {
    EnableTunnelMonitor string `xml:"enable"`
    TunnelMonitorDestinationIp string `xml:"destination-ip,omitempty"`
    TunnelMonitorSourceIp string `xml:"source-ip,omitempty"`
    TunnelMonitorProfile string `xml:"tunnel-monitor-profile,omitempty"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        TunnelInterface: e.TunnelInterface,
        CopyTos: util.YesNo(e.CopyTos),
    }
    if e.AntiReplay {
        // NOTE(gfreeman) PAN-OS errors if you send this as false...???
        ans.AntiReplay = util.YesNo(e.AntiReplay)
    }

    switch e.Type {
    case TypeAutoKey:
        ans.Ak = &ak{
            AkIkeGateway: util.OneStrToEnt(e.AkIkeGateway),
            AkIpsecCryptoProfile: e.AkIpsecCryptoProfile,
        }

        if text, present := e.raw["pv4"]; present {
            ans.Ak.ProxyIpv4 = &util.RawXml{text}
        }

        if text, present := e.raw["pv6"]; present {
            ans.Ak.ProxyIpv6 = &util.RawXml{text}
        }
    case TypeManualKey:
        ans.Mk = &mk_v1{
            MkLocalSpi: e.MkLocalSpi,
            MkRemoteSpi: e.MkRemoteSpi,
            Local: mkLocal_v1{
                MkInterface: e.MkInterface,
                MkLocalAddressIp: e.MkLocalAddressIp,
            },
            Peer: mkPeer{
                MkRemoteAddress: e.MkRemoteAddress,
            },
        }

        switch e.MkProtocol {
        case MkProtocolEsp:
            ans.Mk.Esp = &mkEsp{
                MkEspEncryptionType: e.MkEspEncryptionType,
                MkEspEncryptionKey: e.MkEspEncryptionKey,
            }

            switch e.MkAuthType {
            case MkAuthTypeMd5:
                ans.Mk.Esp.AuthMd5 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha1:
                ans.Mk.Esp.AuthSha1 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha256:
                ans.Mk.Esp.AuthSha256 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha384:
                ans.Mk.Esp.AuthSha384 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha512:
                ans.Mk.Esp.AuthSha512 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeNone:
                s := ""
                ans.Mk.Esp.AuthNone = &s
            }
        case MkProtocolAh:
            switch e.MkAuthType {
            case MkAuthTypeMd5:
                ans.Mk.Ah = &mkAh{AuthMd5: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha1:
                ans.Mk.Ah = &mkAh{AuthSha1: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha256:
                ans.Mk.Ah = &mkAh{AuthSha256: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha384:
                ans.Mk.Ah = &mkAh{AuthSha384: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha512:
                ans.Mk.Ah = &mkAh{AuthSha512: &authKey{Value: e.MkAuthKey}}
            }
        }
    case TypeGlobalProtectSatellite:
        ans.Gps = &gps_v1{
            GpsPortalAddress: e.GpsPortalAddress,
            Local: gpsLocal_v1{
                GpsInterface: e.GpsInterface,
                GpsInterfaceIpIpv4: e.GpsInterfaceIpIpv4,
            },
            GpsPublishRoutes: util.StrToMem(e.GpsPublishRoutes),
        }

        if e.GpsPublishConnectedRoutes {
            ans.Gps.Pcr = &gpsPcr{util.YesNo(e.GpsPublishConnectedRoutes)}
        }

        if e.GpsLocalCertificate != "" || e.GpsCertificateProfile != "" {
            ans.Gps.Ca = &gpsCa{
                GpsLocalCertificate: e.GpsLocalCertificate,
                GpsCertificateProfile: e.GpsCertificateProfile,
            }
        }
    }

    if e.EnableTunnelMonitor || e.TunnelMonitorDestinationIp != "" || e.TunnelMonitorSourceIp != "" || e.TunnelMonitorProfile != "" {
        ans.TunnelMonitor = &tunMon_v1{
            EnableTunnelMonitor: util.YesNo(e.EnableTunnelMonitor),
            TunnelMonitorDestinationIp: e.TunnelMonitorDestinationIp,
            TunnelMonitorSourceIp: e.TunnelMonitorSourceIp,
            TunnelMonitorProfile: e.TunnelMonitorProfile,
        }
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TunnelInterface string `xml:"tunnel-interface"`
    AntiReplay string `xml:"anti-replay,omitempty"`
    Ak *ak `xml:"auto-key"`
    Mk *mk_v2 `xml:"manual-key"`
    Gps *gps_v2 `xml:"global-protect-satellite"`
    TunnelMonitor *tunMon_v2 `xml:"tunnel-monitor"`
    CopyTos string `xml:"copy-tos"`
    EnableIpv6 string `xml:"ipv6"`
    Disabled string `xml:"disabled"`
    CopyFlowLabel string `xml:"copy-flow-label"`
}

type mk_v2 struct {
    MkLocalSpi string `xml:"local-spi"`
    Local mkLocal_v2 `xml:"local-address"`
    Peer mkPeer `xml:"peer-address"`
    MkRemoteSpi string `xml:"remote-spi"`
    Esp *mkEsp `xml:"esp"`
    Ah *mkAh `xml:"ah"`
}

type mkLocal_v2 struct {
    MkInterface string `xml:"interface"`
    MkLocalAddressIp string `xml:"ip,omitempty"`
    MkLocalAddressFloatingIp string `xml:"floating-ip,omitempty"`
}

type gps_v2 struct {
    GpsPortalAddress string `xml:"portal-address"`
    Local gpsLocal_v2 `xml:"local-address"`
    Pcr *gpsPcr `xml:"publish-connected-routes"`
    GpsPublishRoutes *util.MemberType `xml:"publish-routes"`
    Ca *gpsCa `xml:"external-ca"`
}

type gpsLocal_v2 struct {
    GpsInterface string `xml:"interface"`
    GpsInterfaceIpIpv4 string `xml:"ip,omitempty"`
    GpsInterfaceFloatingIpIpv4 string `xml:"floating-ip,omitempty"`
}

type tunMon_v2 struct {
    EnableTunnelMonitor string `xml:"enable"`
    TunnelMonitorDestinationIp string `xml:"destination-ip,omitempty"`
    TunnelMonitorSourceIp string `xml:"source-ip,omitempty"`
    TunnelMonitorProfile string `xml:"tunnel-monitor-profile,omitempty"`
    TunnelMonitorProxyId string `xml:"proxy-id,omitempty"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        TunnelInterface: e.TunnelInterface,
        CopyTos: util.YesNo(e.CopyTos),
        EnableIpv6: util.YesNo(e.EnableIpv6),
        Disabled: util.YesNo(e.Disabled),
        CopyFlowLabel: util.YesNo(e.CopyFlowLabel),
    }
    if e.AntiReplay {
        // NOTE(gfreeman) PAN-OS errors if you send this as false...???
        ans.AntiReplay = util.YesNo(e.AntiReplay)
    }

    switch e.Type {
    case TypeAutoKey:
        ans.Ak = &ak{
            AkIkeGateway: util.OneStrToEnt(e.AkIkeGateway),
            AkIpsecCryptoProfile: e.AkIpsecCryptoProfile,
        }

        if text, present := e.raw["pv4"]; present {
            ans.Ak.ProxyIpv4 = &util.RawXml{text}
        }

        if text, present := e.raw["pv6"]; present {
            ans.Ak.ProxyIpv6 = &util.RawXml{text}
        }
    case TypeManualKey:
        ans.Mk = &mk_v2{
            MkLocalSpi: e.MkLocalSpi,
            MkRemoteSpi: e.MkRemoteSpi,
            Local: mkLocal_v2{
                MkInterface: e.MkInterface,
                MkLocalAddressIp: e.MkLocalAddressIp,
                MkLocalAddressFloatingIp: e.MkLocalAddressFloatingIp,
            },
            Peer: mkPeer{
                MkRemoteAddress: e.MkRemoteAddress,
            },
        }

        switch e.MkProtocol {
        case MkProtocolEsp:
            ans.Mk.Esp = &mkEsp{
                MkEspEncryptionType: e.MkEspEncryptionType,
                MkEspEncryptionKey: e.MkEspEncryptionKey,
            }

            switch e.MkAuthType {
            case MkAuthTypeMd5:
                ans.Mk.Esp.AuthMd5 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha1:
                ans.Mk.Esp.AuthSha1 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha256:
                ans.Mk.Esp.AuthSha256 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha384:
                ans.Mk.Esp.AuthSha384 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha512:
                ans.Mk.Esp.AuthSha512 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeNone:
                s := ""
                ans.Mk.Esp.AuthNone = &s
            }
        case MkProtocolAh:
            switch e.MkAuthType {
            case MkAuthTypeMd5:
                ans.Mk.Ah = &mkAh{AuthMd5: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha1:
                ans.Mk.Ah = &mkAh{AuthSha1: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha256:
                ans.Mk.Ah = &mkAh{AuthSha256: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha384:
                ans.Mk.Ah = &mkAh{AuthSha384: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha512:
                ans.Mk.Ah = &mkAh{AuthSha512: &authKey{Value: e.MkAuthKey}}
            }
        }
    case TypeGlobalProtectSatellite:
        ans.Gps = &gps_v2{
            GpsPortalAddress: e.GpsPortalAddress,
            Local: gpsLocal_v2{
                GpsInterface: e.GpsInterface,
                GpsInterfaceIpIpv4: e.GpsInterfaceIpIpv4,
                GpsInterfaceFloatingIpIpv4: e.GpsInterfaceFloatingIpIpv4,
            },
            GpsPublishRoutes: util.StrToMem(e.GpsPublishRoutes),
        }

        if e.GpsPublishConnectedRoutes {
            ans.Gps.Pcr = &gpsPcr{util.YesNo(e.GpsPublishConnectedRoutes)}
        }

        if e.GpsLocalCertificate != "" || e.GpsCertificateProfile != "" {
            ans.Gps.Ca = &gpsCa{
                GpsLocalCertificate: e.GpsLocalCertificate,
                GpsCertificateProfile: e.GpsCertificateProfile,
            }
        }
    }

    if e.EnableTunnelMonitor || e.TunnelMonitorDestinationIp != "" || e.TunnelMonitorSourceIp != "" || e.TunnelMonitorProfile != "" || e.TunnelMonitorProxyId != "" {
        ans.TunnelMonitor = &tunMon_v2{
            EnableTunnelMonitor: util.YesNo(e.EnableTunnelMonitor),
            TunnelMonitorDestinationIp: e.TunnelMonitorDestinationIp,
            TunnelMonitorSourceIp: e.TunnelMonitorSourceIp,
            TunnelMonitorProfile: e.TunnelMonitorProfile,
            TunnelMonitorProxyId: e.TunnelMonitorProxyId,
        }
    }

    return ans
}

type entry_v3 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    TunnelInterface string `xml:"tunnel-interface"`
    AntiReplay string `xml:"anti-replay,omitempty"`
    Ak *ak `xml:"auto-key"`
    Mk *mk_v2 `xml:"manual-key"`
    Gps *gps_v3 `xml:"global-protect-satellite"`
    TunnelMonitor *tunMon_v2 `xml:"tunnel-monitor"`
    CopyTos string `xml:"copy-tos"`
    EnableIpv6 string `xml:"ipv6"`
    Disabled string `xml:"disabled"`
    CopyFlowLabel string `xml:"copy-flow-label"`
}

type gps_v3 struct {
    GpsPreferIpv6 string `xml:"ipv6-preferred"`
    GpsPortalAddress string `xml:"portal-address"`
    Local gpsLocal_v3 `xml:"local-address"`
    Pcr *gpsPcr `xml:"publish-connected-routes"`
    GpsPublishRoutes *util.MemberType `xml:"publish-routes"`
    Ca *gpsCa `xml:"external-ca"`
}

type gpsLocal_v3 struct {
    GpsInterface string `xml:"interface"`
    Ip *gpsLocalIp `xml:"ip"`
    Floating *gpsLocalIp `xml:"floating-ip"`
}

type gpsLocalIp struct {
    Ipv4 string `xml:"ipv4,omitempty"`
    Ipv6 string `xml:"ipv6,omitempty"`
}

func specify_v3(e Entry) interface{} {
    ans := entry_v3{
        Name: e.Name,
        TunnelInterface: e.TunnelInterface,
        CopyTos: util.YesNo(e.CopyTos),
        EnableIpv6: util.YesNo(e.EnableIpv6),
        Disabled: util.YesNo(e.Disabled),
        CopyFlowLabel: util.YesNo(e.CopyFlowLabel),
    }
    if e.AntiReplay {
        // NOTE(gfreeman) PAN-OS errors if you send this as false...???
        ans.AntiReplay = util.YesNo(e.AntiReplay)
    }

    switch e.Type {
    case TypeAutoKey:
        ans.Ak = &ak{
            AkIkeGateway: util.OneStrToEnt(e.AkIkeGateway),
            AkIpsecCryptoProfile: e.AkIpsecCryptoProfile,
        }

        if text, present := e.raw["pv4"]; present {
            ans.Ak.ProxyIpv4 = &util.RawXml{text}
        }

        if text, present := e.raw["pv6"]; present {
            ans.Ak.ProxyIpv6 = &util.RawXml{text}
        }
    case TypeManualKey:
        ans.Mk = &mk_v2{
            MkLocalSpi: e.MkLocalSpi,
            MkRemoteSpi: e.MkRemoteSpi,
            Local: mkLocal_v2{
                MkInterface: e.MkInterface,
                MkLocalAddressIp: e.MkLocalAddressIp,
                MkLocalAddressFloatingIp: e.MkLocalAddressFloatingIp,
            },
            Peer: mkPeer{
                MkRemoteAddress: e.MkRemoteAddress,
            },
        }

        switch e.MkProtocol {
        case MkProtocolEsp:
            ans.Mk.Esp = &mkEsp{
                MkEspEncryptionType: e.MkEspEncryptionType,
                MkEspEncryptionKey: e.MkEspEncryptionKey,
            }

            switch e.MkAuthType {
            case MkAuthTypeMd5:
                ans.Mk.Esp.AuthMd5 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha1:
                ans.Mk.Esp.AuthSha1 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha256:
                ans.Mk.Esp.AuthSha256 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha384:
                ans.Mk.Esp.AuthSha384 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeSha512:
                ans.Mk.Esp.AuthSha512 = &authKey{Value: e.MkAuthKey}
            case MkAuthTypeNone:
                s := ""
                ans.Mk.Esp.AuthNone = &s
            }
        case MkProtocolAh:
            switch e.MkAuthType {
            case MkAuthTypeMd5:
                ans.Mk.Ah = &mkAh{AuthMd5: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha1:
                ans.Mk.Ah = &mkAh{AuthSha1: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha256:
                ans.Mk.Ah = &mkAh{AuthSha256: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha384:
                ans.Mk.Ah = &mkAh{AuthSha384: &authKey{Value: e.MkAuthKey}}
            case MkAuthTypeSha512:
                ans.Mk.Ah = &mkAh{AuthSha512: &authKey{Value: e.MkAuthKey}}
            }
        }
    case TypeGlobalProtectSatellite:
        ans.Gps = &gps_v3{
            GpsPortalAddress: e.GpsPortalAddress,
            GpsPreferIpv6: util.YesNo(e.GpsPreferIpv6),
            Local: gpsLocal_v3{
                GpsInterface: e.GpsInterface,
            },
            GpsPublishRoutes: util.StrToMem(e.GpsPublishRoutes),
        }

        if e.GpsPublishConnectedRoutes {
            ans.Gps.Pcr = &gpsPcr{util.YesNo(e.GpsPublishConnectedRoutes)}
        }

        if e.GpsInterfaceIpIpv4 != "" || e.GpsInterfaceIpIpv6 != "" {
            ans.Gps.Local.Ip = &gpsLocalIp{
                Ipv4: e.GpsInterfaceIpIpv4,
                Ipv6: e.GpsInterfaceIpIpv6,
            }
        }

        if e.GpsInterfaceFloatingIpIpv4 != "" || e.GpsInterfaceFloatingIpIpv6 != "" {
            ans.Gps.Local.Floating = &gpsLocalIp{
                Ipv4: e.GpsInterfaceFloatingIpIpv4,
                Ipv6: e.GpsInterfaceFloatingIpIpv6,
            }
        }

        if e.GpsLocalCertificate != "" || e.GpsCertificateProfile != "" {
            ans.Gps.Ca = &gpsCa{
                GpsLocalCertificate: e.GpsLocalCertificate,
                GpsCertificateProfile: e.GpsCertificateProfile,
            }
        }
    }

    if e.EnableTunnelMonitor || e.TunnelMonitorDestinationIp != "" || e.TunnelMonitorSourceIp != "" || e.TunnelMonitorProfile != "" || e.TunnelMonitorProxyId != "" {
        ans.TunnelMonitor = &tunMon_v2{
            EnableTunnelMonitor: util.YesNo(e.EnableTunnelMonitor),
            TunnelMonitorDestinationIp: e.TunnelMonitorDestinationIp,
            TunnelMonitorSourceIp: e.TunnelMonitorSourceIp,
            TunnelMonitorProfile: e.TunnelMonitorProfile,
            TunnelMonitorProxyId: e.TunnelMonitorProxyId,
        }
    }

    return ans
}
