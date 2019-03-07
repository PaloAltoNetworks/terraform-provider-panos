package ikegw

import (
    "encoding/xml"

    "github.com/PaloAltoNetworks/pango/util"
)


const (
    Ikev1 = "ikev1"
    Ikev2 = "ikev2"
    Ikev2Preferred = "ikev2-preferred"
)

const (
    IdTypeIpAddress = "ipaddr"
    IdTypeFqdn = "fqdn"
    IdTypeUfqdn = "ufqdn"
    IdTypeKeyId = "keyid"
    IdTypeDn = "dn"
)

const (
    PeerTypeIp = "ip"
    PeerTypeDynamic = "dynamic"
    PeerTypeFqdn = "fqdn"
)

const (
    AuthPreSharedKey = "pre-shared-key"
    AuthCertificate = "certificate"
)

const (
    PeerIdCheckExact = "exact"
    PeerIdCheckWildcard = "wildcard"
)

// Entry is a normalized, version independent representation of an IKE gateway.
type Entry struct {
    Name string
    Version string
    EnableIpv6 bool
    Disabled bool
    PeerIpType string
    PeerIpValue string
    Interface string
    LocalIpAddressType string
    LocalIpAddressValue string
    AuthType string
    PreSharedKey string
    LocalIdType string
    LocalIdValue string
    PeerIdType string
    PeerIdValue string
    PeerIdCheck string
    LocalCert string
    CertEnableHashAndUrl bool
    CertBaseUrl string
    CertUseManagementAsSource bool
    CertPermitPayloadMismatch bool
    CertProfile string
    CertEnableStrictValidation bool
    EnablePassiveMode bool
    EnableNatTraversal bool
    NatTraversalKeepAlive int
    NatTraversalEnableUdpChecksum bool
    EnableFragmentation bool
    Ikev1ExchangeMode string
    Ikev1CryptoProfile string
    EnableDeadPeerDetection bool
    DeadPeerDetectionInterval int
    DeadPeerDetectionRetry int
    Ikev2CryptoProfile string
    Ikev2CookieValidation bool
    EnableLivenessCheck bool
    LivenessCheckInterval int
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
    o.Version = s.Version
    o.EnableIpv6 = s.EnableIpv6
    o.Disabled = s.Disabled
    o.PeerIpType = s.PeerIpType
    o.PeerIpValue = s.PeerIpValue
    o.Interface = s.Interface
    o.LocalIpAddressType = s.LocalIpAddressType
    o.LocalIpAddressValue = s.LocalIpAddressValue
    o.AuthType = s.AuthType
    o.PreSharedKey = s.PreSharedKey
    o.LocalIdType = s.LocalIdType
    o.LocalIdValue = s.LocalIdValue
    o.PeerIdType = s.PeerIdType
    o.PeerIdValue = s.PeerIdValue
    o.PeerIdCheck = s.PeerIdCheck
    o.LocalCert = s.LocalCert
    o.CertEnableHashAndUrl = s.CertEnableHashAndUrl
    o.CertBaseUrl = s.CertBaseUrl
    o.CertUseManagementAsSource = s.CertUseManagementAsSource
    o.CertPermitPayloadMismatch = s.CertPermitPayloadMismatch
    o.CertProfile = s.CertProfile
    o.CertEnableStrictValidation = s.CertEnableStrictValidation
    o.EnablePassiveMode = s.EnablePassiveMode
    o.EnableNatTraversal = s.EnableNatTraversal
    o.NatTraversalKeepAlive = s.NatTraversalKeepAlive
    o.NatTraversalEnableUdpChecksum = s.NatTraversalEnableUdpChecksum
    o.EnableFragmentation = s.EnableFragmentation
    o.Ikev1ExchangeMode = s.Ikev1ExchangeMode
    o.Ikev1CryptoProfile = s.Ikev1CryptoProfile
    o.EnableDeadPeerDetection = s.EnableDeadPeerDetection
    o.DeadPeerDetectionInterval = s.DeadPeerDetectionInterval
    o.DeadPeerDetectionRetry = s.DeadPeerDetectionRetry
    o.Ikev2CryptoProfile = s.Ikev2CryptoProfile
    o.Ikev2CookieValidation = s.Ikev2CookieValidation
    o.EnableLivenessCheck = s.EnableLivenessCheck
    o.LivenessCheckInterval = s.LivenessCheckInterval
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
        Version: Ikev1,
        Interface: o.Answer.LocalIp.Interface,
    }

    if o.Answer.PeerIp.Dynamic != nil {
        ans.PeerIpType = PeerTypeDynamic
    } else {
        ans.PeerIpType = PeerTypeIp
        ans.PeerIpValue = o.Answer.PeerIp.Static
    }

    if o.Answer.PeerId != nil {
        ans.PeerIdType = o.Answer.PeerId.PeerIdType
        ans.PeerIdValue = o.Answer.PeerId.PeerIdValue
        ans.PeerIdCheck = o.Answer.PeerId.PeerIdCheck
    }

    if o.Answer.LocalIp.StaticIp != "" {
        ans.LocalIpAddressType = PeerTypeIp
        ans.LocalIpAddressValue = o.Answer.LocalIp.StaticIp
    }

    if o.Answer.LocalId != nil {
        ans.LocalIdType = o.Answer.LocalId.LocalIdType
        ans.LocalIdValue = o.Answer.LocalId.LocalIdValue
    }

    if o.Answer.PskAuth != nil {
        ans.AuthType = AuthPreSharedKey
        ans.PreSharedKey = o.Answer.PskAuth.Key
    } else if o.Answer.CAuth != nil {
        ans.AuthType = AuthCertificate
        ans.LocalCert = o.Answer.CAuth.LocalCert
        ans.CertProfile = o.Answer.CAuth.CertProfile
        ans.CertEnableStrictValidation = util.AsBool(o.Answer.CAuth.CertEnableStrictValidation)
        ans.CertPermitPayloadMismatch = util.AsBool(o.Answer.CAuth.CertPermitPayloadMismatch)
    }

    if o.Answer.Proto != nil {
        if o.Answer.Proto.Ikev1 != nil {
            ans.Ikev1ExchangeMode = o.Answer.Proto.Ikev1.Ikev1ExchangeMode
            ans.Ikev1CryptoProfile = o.Answer.Proto.Ikev1.Ikev1CryptoProfile

            if o.Answer.Proto.Ikev1.Dpd != nil {
                ans.EnableDeadPeerDetection = util.AsBool(o.Answer.Proto.Ikev1.Dpd.EnableDeadPeerDetection)
                ans.DeadPeerDetectionInterval = o.Answer.Proto.Ikev1.Dpd.DeadPeerDetectionInterval
                ans.DeadPeerDetectionRetry = o.Answer.Proto.Ikev1.Dpd.DeadPeerDetectionRetry
            }
        }
    }

    if o.Answer.ProtoCommon != nil {
        ans.EnablePassiveMode = util.AsBool(o.Answer.ProtoCommon.EnablePassiveMode)
        if o.Answer.ProtoCommon.Nat != nil {
            ans.EnableNatTraversal = util.AsBool(o.Answer.ProtoCommon.Nat.EnableNatTraversal)
            ans.NatTraversalKeepAlive = o.Answer.ProtoCommon.Nat.NatTraversalKeepAlive
            ans.NatTraversalEnableUdpChecksum = util.AsBool(o.Answer.ProtoCommon.Nat.NatTraversalEnableUdpChecksum)
        }
        if o.Answer.ProtoCommon.Frag != nil {
            ans.EnableFragmentation = util.AsBool(o.Answer.ProtoCommon.Frag.EnableFragmentation)
        }
    }

    return ans
}

type container_v2 struct {
    Answer entry_v2 `xml:"result>entry"`
}

func (o *container_v2) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Interface: o.Answer.LocalIp.Interface,
        Disabled: util.AsBool(o.Answer.Disabled),
        EnableIpv6: util.AsBool(o.Answer.EnableIpv6),
    }

    if o.Answer.PeerIp.Dynamic != nil {
        ans.PeerIpType = PeerTypeDynamic
    } else {
        ans.PeerIpType = PeerTypeIp
        ans.PeerIpValue = o.Answer.PeerIp.Static
    }

    if o.Answer.PeerId != nil {
        ans.PeerIdType = o.Answer.PeerId.PeerIdType
        ans.PeerIdValue = o.Answer.PeerId.PeerIdValue
        ans.PeerIdCheck = o.Answer.PeerId.PeerIdCheck
    }

    if o.Answer.LocalIp.StaticIp != "" {
        ans.LocalIpAddressType = PeerTypeIp
        ans.LocalIpAddressValue = o.Answer.LocalIp.StaticIp
    }

    if o.Answer.LocalId != nil {
        ans.LocalIdType = o.Answer.LocalId.LocalIdType
        ans.LocalIdValue = o.Answer.LocalId.LocalIdValue
    }

    if o.Answer.PskAuth != nil {
        ans.AuthType = AuthPreSharedKey
        ans.PreSharedKey = o.Answer.PskAuth.Key
    } else if o.Answer.CAuth != nil {
        ans.AuthType = AuthCertificate
        ans.LocalCert = o.Answer.CAuth.CLocal.LocalCert
        ans.CertProfile = o.Answer.CAuth.CertProfile
        ans.CertEnableStrictValidation = util.AsBool(o.Answer.CAuth.CertEnableStrictValidation)
        ans.CertPermitPayloadMismatch = util.AsBool(o.Answer.CAuth.CertPermitPayloadMismatch)

        if o.Answer.CAuth.CLocal.Hau != nil {
            ans.CertEnableHashAndUrl = util.AsBool(o.Answer.CAuth.CLocal.Hau.CertEnableHashAndUrl)
            ans.CertBaseUrl = o.Answer.CAuth.CLocal.Hau.CertBaseUrl
        }
    }

    if o.Answer.Proto != nil {
        ans.Version = o.Answer.Proto.Version

        if o.Answer.Proto.Ikev1 != nil {
            ans.Ikev1ExchangeMode = o.Answer.Proto.Ikev1.Ikev1ExchangeMode
            ans.Ikev1CryptoProfile = o.Answer.Proto.Ikev1.Ikev1CryptoProfile

            if o.Answer.Proto.Ikev1.Dpd != nil {
                ans.EnableDeadPeerDetection = util.AsBool(o.Answer.Proto.Ikev1.Dpd.EnableDeadPeerDetection)
                ans.DeadPeerDetectionInterval = o.Answer.Proto.Ikev1.Dpd.DeadPeerDetectionInterval
                ans.DeadPeerDetectionRetry = o.Answer.Proto.Ikev1.Dpd.DeadPeerDetectionRetry
            }
        }

        if o.Answer.Proto.Ikev2 != nil {
            ans.Ikev2CryptoProfile = o.Answer.Proto.Ikev2.Ikev2CryptoProfile
            ans.Ikev2CookieValidation = util.AsBool(o.Answer.Proto.Ikev2.Ikev2CookieValidation)

            if o.Answer.Proto.Ikev2.Dpd != nil {
                ans.EnableLivenessCheck = util.AsBool(o.Answer.Proto.Ikev2.Dpd.EnableLivenessCheck)
                ans.LivenessCheckInterval = o.Answer.Proto.Ikev2.Dpd.LivenessCheckInterval
            }
        }
    }

    if o.Answer.ProtoCommon != nil {
        ans.EnablePassiveMode = util.AsBool(o.Answer.ProtoCommon.EnablePassiveMode)
        if o.Answer.ProtoCommon.Nat != nil {
            ans.EnableNatTraversal = util.AsBool(o.Answer.ProtoCommon.Nat.EnableNatTraversal)
            ans.NatTraversalKeepAlive = o.Answer.ProtoCommon.Nat.NatTraversalKeepAlive
            ans.NatTraversalEnableUdpChecksum = util.AsBool(o.Answer.ProtoCommon.Nat.NatTraversalEnableUdpChecksum)
        }
        if o.Answer.ProtoCommon.Frag != nil {
            ans.EnableFragmentation = util.AsBool(o.Answer.ProtoCommon.Frag.EnableFragmentation)
        }
    }

    return ans
}

type container_v3 struct {
    Answer entry_v3 `xml:"result>entry"`
}

func (o *container_v3) Normalize() Entry {
    ans := Entry{
        Name: o.Answer.Name,
        Interface: o.Answer.LocalIp.Interface,
        Disabled: util.AsBool(o.Answer.Disabled),
        EnableIpv6: util.AsBool(o.Answer.EnableIpv6),
    }

    if o.Answer.PeerIp.Dynamic != nil {
        ans.PeerIpType = PeerTypeDynamic
    } else if o.Answer.PeerIp.Fqdn != "" {
        ans.PeerIpType = PeerTypeFqdn
        ans.PeerIpValue = o.Answer.PeerIp.Fqdn
    } else {
        ans.PeerIpType = PeerTypeIp
        ans.PeerIpValue = o.Answer.PeerIp.Static
    }

    if o.Answer.PeerId != nil {
        ans.PeerIdType = o.Answer.PeerId.PeerIdType
        ans.PeerIdValue = o.Answer.PeerId.PeerIdValue
        ans.PeerIdCheck = o.Answer.PeerId.PeerIdCheck
    }

    if o.Answer.LocalIp.StaticIp != "" {
        ans.LocalIpAddressType = PeerTypeIp
        ans.LocalIpAddressValue = o.Answer.LocalIp.StaticIp
    }

    if o.Answer.LocalId != nil {
        ans.LocalIdType = o.Answer.LocalId.LocalIdType
        ans.LocalIdValue = o.Answer.LocalId.LocalIdValue
    }

    if o.Answer.PskAuth != nil {
        ans.AuthType = AuthPreSharedKey
        ans.PreSharedKey = o.Answer.PskAuth.Key
    } else if o.Answer.CAuth != nil {
        ans.AuthType = AuthCertificate
        ans.LocalCert = o.Answer.CAuth.CLocal.LocalCert
        ans.CertProfile = o.Answer.CAuth.CertProfile
        ans.CertEnableStrictValidation = util.AsBool(o.Answer.CAuth.CertEnableStrictValidation)
        ans.CertPermitPayloadMismatch = util.AsBool(o.Answer.CAuth.CertPermitPayloadMismatch)

        if o.Answer.CAuth.CLocal.Hau != nil {
            ans.CertEnableHashAndUrl = util.AsBool(o.Answer.CAuth.CLocal.Hau.CertEnableHashAndUrl)
            ans.CertBaseUrl = o.Answer.CAuth.CLocal.Hau.CertBaseUrl
        }
    }

    if o.Answer.Proto != nil {
        ans.Version = o.Answer.Proto.Version

        if o.Answer.Proto.Ikev1 != nil {
            ans.Ikev1ExchangeMode = o.Answer.Proto.Ikev1.Ikev1ExchangeMode
            ans.Ikev1CryptoProfile = o.Answer.Proto.Ikev1.Ikev1CryptoProfile

            if o.Answer.Proto.Ikev1.Dpd != nil {
                ans.EnableDeadPeerDetection = util.AsBool(o.Answer.Proto.Ikev1.Dpd.EnableDeadPeerDetection)
                ans.DeadPeerDetectionInterval = o.Answer.Proto.Ikev1.Dpd.DeadPeerDetectionInterval
                ans.DeadPeerDetectionRetry = o.Answer.Proto.Ikev1.Dpd.DeadPeerDetectionRetry
            }
        }

        if o.Answer.Proto.Ikev2 != nil {
            ans.Ikev2CryptoProfile = o.Answer.Proto.Ikev2.Ikev2CryptoProfile
            ans.Ikev2CookieValidation = util.AsBool(o.Answer.Proto.Ikev2.Ikev2CookieValidation)

            if o.Answer.Proto.Ikev2.Dpd != nil {
                ans.EnableLivenessCheck = util.AsBool(o.Answer.Proto.Ikev2.Dpd.EnableLivenessCheck)
                ans.LivenessCheckInterval = o.Answer.Proto.Ikev2.Dpd.LivenessCheckInterval
            }
        }
    }

    if o.Answer.ProtoCommon != nil {
        ans.EnablePassiveMode = util.AsBool(o.Answer.ProtoCommon.EnablePassiveMode)
        if o.Answer.ProtoCommon.Nat != nil {
            ans.EnableNatTraversal = util.AsBool(o.Answer.ProtoCommon.Nat.EnableNatTraversal)
            ans.NatTraversalKeepAlive = o.Answer.ProtoCommon.Nat.NatTraversalKeepAlive
            ans.NatTraversalEnableUdpChecksum = util.AsBool(o.Answer.ProtoCommon.Nat.NatTraversalEnableUdpChecksum)
        }
        if o.Answer.ProtoCommon.Frag != nil {
            ans.EnableFragmentation = util.AsBool(o.Answer.ProtoCommon.Frag.EnableFragmentation)
        }
    }

    return ans
}

type entry_v1 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    PeerIp peerIp_v1 `xml:"peer-address"`
    PeerId *peerId `xml:"peer-id"`
    LocalIp localIp `xml:"local-address"`
    LocalId *localId `xml:"local-id"`
    PskAuth *pskAuth `xml:"authentication>pre-shared-key"`
    CAuth *cAuth_v1 `xml:"authentication>certificate"`
    Proto *proto_v1 `xml:"protocol"`
    ProtoCommon *protoCommon `xml:"protocol-common"`
}

type peerIp_v1 struct {
    Static string `xml:"ip,omitempty"`
    Dynamic *string `xml:"dynamic"`
}

type peerId struct {
    PeerIdType string `xml:"type"`
    PeerIdValue string `xml:"id,omitempty"`
    PeerIdCheck string `xml:"matching,omitempty"`
}

type localIp struct {
    Interface string `xml:"interface,omitempty"`
    StaticIp string `xml:"ip,omitempty"`
}

type localId struct {
    LocalIdType string `xml:"type"`
    LocalIdValue string `xml:"id"`
}

type pskAuth struct {
    Key string `xml:"key"`
}

type cAuth_v1 struct {
    LocalCert string `xml:"local-certificate"`
    CertProfile string `xml:"certificate-profile"`
    CertEnableStrictValidation string `xml:"strict-validation-revocation"`
    CertPermitPayloadMismatch string `xml:"allow-id-payload-mismatch"`
}

type proto_v1 struct {
    Ikev1 *ikev1_v1 `xml:"ikev1"`
}

type ikev1_v1 struct {
    Ikev1ExchangeMode string `xml:"exchange-mode,omitempty"`
    Ikev1CryptoProfile string `xml:"ike-crypto-profile,omitempty"`
    Dpd *ikev1Dpd `xml:"dpd"`
}

type ikev1Dpd struct {
    EnableDeadPeerDetection string `xml:"enable"`
    DeadPeerDetectionInterval int `xml:"interval,omitempty"`
    DeadPeerDetectionRetry int `xml:"retry,omitempty"`
}

type protoCommon struct {
    EnablePassiveMode string `xml:"passive-mode"`
    Nat *protoNat `xml:"nat-traversal"`
    Frag *protoFrag `xml:"fragmentation"`
}

type protoNat struct {
    EnableNatTraversal string `xml:"enable"`
    NatTraversalKeepAlive int `xml:"keep-alive-interval"`
    NatTraversalEnableUdpChecksum string `xml:"udp-checksum-enable"`
}

type protoFrag struct {
    EnableFragmentation string `xml:"enable"`
}

func specify_v1(e Entry) interface{} {
    ans := entry_v1{
        Name: e.Name,
        LocalIp: localIp{
            Interface: e.Interface,
            StaticIp: e.LocalIpAddressValue,
        },
    }

    switch e.PeerIpType {
    case PeerTypeIp:
        ans.PeerIp.Static = e.PeerIpValue
    case PeerTypeDynamic:
        s := ""
        ans.PeerIp.Dynamic = &s
    }

    if e.PeerIdType != "" || e.PeerIdValue != "" || e.PeerIdCheck != "" {
        ans.PeerId = &peerId{
            PeerIdType: e.PeerIdType,
            PeerIdValue: e.PeerIdValue,
            PeerIdCheck: e.PeerIdCheck,
        }
    }

    if e.LocalIdType != "" || e.LocalIdValue != "" {
        ans.LocalId = &localId{
            LocalIdType: e.LocalIdType,
            LocalIdValue: e.LocalIdValue,
        }
    }

    switch e.AuthType {
    case AuthPreSharedKey:
        ans.PskAuth = &pskAuth{
            Key: e.PreSharedKey,
        }
    case AuthCertificate:
        ans.CAuth = &cAuth_v1{
            LocalCert: e.LocalCert,
            CertProfile: e.CertProfile,
            CertEnableStrictValidation: util.YesNo(e.CertEnableStrictValidation),
            CertPermitPayloadMismatch: util.YesNo(e.CertPermitPayloadMismatch),
        }
    }

    if e.Ikev1ExchangeMode != "" || e.Ikev1CryptoProfile != "" || e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 {
        ans.Proto = &proto_v1{
            Ikev1: &ikev1_v1{
                Ikev1ExchangeMode: e.Ikev1ExchangeMode,
                Ikev1CryptoProfile: e.Ikev1CryptoProfile,
            },
        }

        if e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 {
            ans.Proto.Ikev1.Dpd = &ikev1Dpd{
                EnableDeadPeerDetection: util.YesNo(e.EnableDeadPeerDetection),
                DeadPeerDetectionInterval: e.DeadPeerDetectionInterval,
                DeadPeerDetectionRetry: e.DeadPeerDetectionRetry,
            }
        }
    }

    if e.EnablePassiveMode || e.EnableNatTraversal || e.NatTraversalKeepAlive != 0 || e.NatTraversalEnableUdpChecksum || e.EnableFragmentation {
        s := protoCommon{
            EnablePassiveMode: util.YesNo(e.EnablePassiveMode),
        }

        if e.EnableNatTraversal || e.NatTraversalKeepAlive != 0 || e.NatTraversalEnableUdpChecksum {
            s.Nat = &protoNat{
                EnableNatTraversal: util.YesNo(e.EnableNatTraversal),
                NatTraversalKeepAlive: e.NatTraversalKeepAlive,
                NatTraversalEnableUdpChecksum: util.YesNo(e.NatTraversalEnableUdpChecksum),
            }
        }

        if e.EnableFragmentation {
            s.Frag = &protoFrag{
                EnableFragmentation: util.YesNo(e.EnableFragmentation),
            }
        }

        ans.ProtoCommon = &s
    }

    return ans
}

type entry_v2 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Disabled string `xml:"disabled"`
    EnableIpv6 string `xml:"ipv6"`
    PeerIp peerIp_v1 `xml:"peer-address"`
    PeerId *peerId `xml:"peer-id"`
    LocalIp localIp `xml:"local-address"`
    LocalId *localId `xml:"local-id"`
    PskAuth *pskAuth `xml:"authentication>pre-shared-key"`
    CAuth *cAuth_v2 `xml:"authentication>certificate"`
    Proto *proto_v2 `xml:"protocol"`
    ProtoCommon *protoCommon `xml:"protocol-common"`
}

type cAuth_v2 struct {
    CLocal cLocal `xml:"local-certificate"`
    CertProfile string `xml:"certificate-profile"`
    CertUseManagementAsSource string `xml:"use-management-as-source"`
    CertEnableStrictValidation string `xml:"strict-validation-revocation"`
    CertPermitPayloadMismatch string `xml:"allow-id-payload-mismatch"`
}

type cLocal struct {
    LocalCert string `xml:"name"`
    Hau *hau `xml:"hash-and-url"`
}

type hau struct {
    CertEnableHashAndUrl string `xml:"enable"`
    CertBaseUrl string `xml:"base-url"`
}

type proto_v2 struct {
    Version string `xml:"version"`
    Ikev1 *ikev1_v1 `xml:"ikev1"`
    Ikev2 *ikev2_v1 `xml:"ikev2"`
}

type ikev2_v1 struct {
    Ikev2CryptoProfile string `xml:"ike-crypto-profile,omitempty"`
    Ikev2CookieValidation string `xml:"require-cookie"`
    Dpd *ikev2Dpd `xml:"dpd"`
}

type ikev2Dpd struct {
    EnableLivenessCheck string `xml:"enable"`
    LivenessCheckInterval int `xml:"interval,omitempty"`
}

func specify_v2(e Entry) interface{} {
    ans := entry_v2{
        Name: e.Name,
        Disabled: util.YesNo(e.Disabled),
        EnableIpv6: util.YesNo(e.EnableIpv6),
        LocalIp: localIp{
            Interface: e.Interface,
            StaticIp: e.LocalIpAddressValue,
        },
    }

    switch e.PeerIpType {
    case PeerTypeIp:
        ans.PeerIp.Static = e.PeerIpValue
    case PeerTypeDynamic:
        s := ""
        ans.PeerIp.Dynamic = &s
    }

    if e.PeerIdType != "" || e.PeerIdValue != "" || e.PeerIdCheck != "" {
        ans.PeerId = &peerId{
            PeerIdType: e.PeerIdType,
            PeerIdValue: e.PeerIdValue,
            PeerIdCheck: e.PeerIdCheck,
        }
    }

    if e.LocalIdType != "" || e.LocalIdValue != "" {
        ans.LocalId = &localId{
            LocalIdType: e.LocalIdType,
            LocalIdValue: e.LocalIdValue,
        }
    }

    switch e.AuthType {
    case AuthPreSharedKey:
        ans.PskAuth = &pskAuth{
            Key: e.PreSharedKey,
        }
    case AuthCertificate:
        ans.CAuth = &cAuth_v2{
            CLocal: cLocal{
                LocalCert: e.LocalCert,
            },
            CertProfile: e.CertProfile,
            CertUseManagementAsSource: util.YesNo(e.CertUseManagementAsSource),
            CertEnableStrictValidation: util.YesNo(e.CertEnableStrictValidation),
            CertPermitPayloadMismatch: util.YesNo(e.CertPermitPayloadMismatch),
        }

        if e.CertEnableHashAndUrl || e.CertBaseUrl != "" {
            ans.CAuth.CLocal.Hau = &hau{
                CertEnableHashAndUrl: util.YesNo(e.CertEnableHashAndUrl),
                CertBaseUrl: e.CertBaseUrl,
            }
        }
    }

    if e.Version != "" || e.Ikev1ExchangeMode != "" || e.Ikev1CryptoProfile != "" || e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 || e.Ikev2CryptoProfile != "" || e.Ikev2CookieValidation || e.EnableLivenessCheck || e.LivenessCheckInterval != 0 {
        ans.Proto = &proto_v2{
            Version: e.Version,
        }

        if e.Ikev1ExchangeMode != "" || e.Ikev1CryptoProfile != "" || e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 {
            ans.Proto.Ikev1 = &ikev1_v1{
                Ikev1ExchangeMode: e.Ikev1ExchangeMode,
                Ikev1CryptoProfile: e.Ikev1CryptoProfile,
            }

            if e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 {
                ans.Proto.Ikev1.Dpd = &ikev1Dpd{
                    EnableDeadPeerDetection: util.YesNo(e.EnableDeadPeerDetection),
                    DeadPeerDetectionInterval: e.DeadPeerDetectionInterval,
                    DeadPeerDetectionRetry: e.DeadPeerDetectionRetry,
                }
            }
        }

        if e.Ikev2CryptoProfile != "" || e.Ikev2CookieValidation || e.EnableLivenessCheck || e.LivenessCheckInterval != 0 {
            ans.Proto.Ikev2 = &ikev2_v1{
                Ikev2CryptoProfile: e.Ikev2CryptoProfile,
                Ikev2CookieValidation: util.YesNo(e.Ikev2CookieValidation),
            }

            if e.EnableLivenessCheck || e.LivenessCheckInterval != 0 {
                ans.Proto.Ikev2.Dpd = &ikev2Dpd{
                    EnableLivenessCheck: util.YesNo(e.EnableLivenessCheck),
                    LivenessCheckInterval: e.LivenessCheckInterval,
                }
            }
        }
    }

    if e.EnablePassiveMode || e.EnableNatTraversal || e.NatTraversalKeepAlive != 0 || e.NatTraversalEnableUdpChecksum || e.EnableFragmentation {
        s := protoCommon{
            EnablePassiveMode: util.YesNo(e.EnablePassiveMode),
        }

        if e.EnableNatTraversal || e.NatTraversalKeepAlive != 0 || e.NatTraversalEnableUdpChecksum {
            s.Nat = &protoNat{
                EnableNatTraversal: util.YesNo(e.EnableNatTraversal),
                NatTraversalKeepAlive: e.NatTraversalKeepAlive,
                NatTraversalEnableUdpChecksum: util.YesNo(e.NatTraversalEnableUdpChecksum),
            }
        }

        if e.EnableFragmentation {
            s.Frag = &protoFrag{
                EnableFragmentation: util.YesNo(e.EnableFragmentation),
            }
        }

        ans.ProtoCommon = &s
    }

    return ans
}

type entry_v3 struct {
    XMLName xml.Name `xml:"entry"`
    Name string `xml:"name,attr"`
    Disabled string `xml:"disabled"`
    EnableIpv6 string `xml:"ipv6"`
    PeerIp peerIp_v2 `xml:"peer-address"`
    PeerId *peerId `xml:"peer-id"`
    LocalIp localIp `xml:"local-address"`
    LocalId *localId `xml:"local-id"`
    PskAuth *pskAuth `xml:"authentication>pre-shared-key"`
    CAuth *cAuth_v2 `xml:"authentication>certificate"`
    Proto *proto_v2 `xml:"protocol"`
    ProtoCommon *protoCommon `xml:"protocol-common"`
}

type peerIp_v2 struct {
    Static string `xml:"ip,omitempty"`
    Dynamic *string `xml:"dynamic"`
    Fqdn string `xml:"fqdn,omitempty"`
}

func specify_v3(e Entry) interface{} {
    ans := entry_v3{
        Name: e.Name,
        Disabled: util.YesNo(e.Disabled),
        EnableIpv6: util.YesNo(e.EnableIpv6),
        LocalIp: localIp{
            Interface: e.Interface,
            StaticIp: e.LocalIpAddressValue,
        },
    }

    switch e.PeerIpType {
    case PeerTypeIp:
        ans.PeerIp.Static = e.PeerIpValue
    case PeerTypeFqdn:
        ans.PeerIp.Fqdn = e.PeerIpValue
    case PeerTypeDynamic:
        s := ""
        ans.PeerIp.Dynamic = &s
    }

    if e.PeerIdType != "" || e.PeerIdValue != "" || e.PeerIdCheck != "" {
        ans.PeerId = &peerId{
            PeerIdType: e.PeerIdType,
            PeerIdValue: e.PeerIdValue,
            PeerIdCheck: e.PeerIdCheck,
        }
    }

    if e.LocalIdType != "" || e.LocalIdValue != "" {
        ans.LocalId = &localId{
            LocalIdType: e.LocalIdType,
            LocalIdValue: e.LocalIdValue,
        }
    }

    switch e.AuthType {
    case AuthPreSharedKey:
        ans.PskAuth = &pskAuth{
            Key: e.PreSharedKey,
        }
    case AuthCertificate:
        ans.CAuth = &cAuth_v2{
            CLocal: cLocal{
                LocalCert: e.LocalCert,
            },
            CertProfile: e.CertProfile,
            CertUseManagementAsSource: util.YesNo(e.CertUseManagementAsSource),
            CertEnableStrictValidation: util.YesNo(e.CertEnableStrictValidation),
            CertPermitPayloadMismatch: util.YesNo(e.CertPermitPayloadMismatch),
        }

        if e.CertEnableHashAndUrl || e.CertBaseUrl != "" {
            ans.CAuth.CLocal.Hau = &hau{
                CertEnableHashAndUrl: util.YesNo(e.CertEnableHashAndUrl),
                CertBaseUrl: e.CertBaseUrl,
            }
        }
    }

    if e.Version != "" || e.Ikev1ExchangeMode != "" || e.Ikev1CryptoProfile != "" || e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 || e.Ikev2CryptoProfile != "" || e.Ikev2CookieValidation || e.EnableLivenessCheck || e.LivenessCheckInterval != 0 {
        ans.Proto = &proto_v2{
            Version: e.Version,
        }

        if e.Ikev1ExchangeMode != "" || e.Ikev1CryptoProfile != "" || e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 {
            ans.Proto.Ikev1 = &ikev1_v1{
                Ikev1ExchangeMode: e.Ikev1ExchangeMode,
                Ikev1CryptoProfile: e.Ikev1CryptoProfile,
            }

            if e.EnableDeadPeerDetection || e.DeadPeerDetectionInterval != 0 || e.DeadPeerDetectionRetry != 0 {
                ans.Proto.Ikev1.Dpd = &ikev1Dpd{
                    EnableDeadPeerDetection: util.YesNo(e.EnableDeadPeerDetection),
                    DeadPeerDetectionInterval: e.DeadPeerDetectionInterval,
                    DeadPeerDetectionRetry: e.DeadPeerDetectionRetry,
                }
            }
        }

        if e.Ikev2CryptoProfile != "" || e.Ikev2CookieValidation || e.EnableLivenessCheck || e.LivenessCheckInterval != 0 {
            ans.Proto.Ikev2 = &ikev2_v1{
                Ikev2CryptoProfile: e.Ikev2CryptoProfile,
                Ikev2CookieValidation: util.YesNo(e.Ikev2CookieValidation),
            }

            if e.EnableLivenessCheck || e.LivenessCheckInterval != 0 {
                ans.Proto.Ikev2.Dpd = &ikev2Dpd{
                    EnableLivenessCheck: util.YesNo(e.EnableLivenessCheck),
                    LivenessCheckInterval: e.LivenessCheckInterval,
                }
            }
        }
    }

    if e.EnablePassiveMode || e.EnableNatTraversal || e.NatTraversalKeepAlive != 0 || e.NatTraversalEnableUdpChecksum || e.EnableFragmentation {
        s := protoCommon{
            EnablePassiveMode: util.YesNo(e.EnablePassiveMode),
        }

        if e.EnableNatTraversal || e.NatTraversalKeepAlive != 0 || e.NatTraversalEnableUdpChecksum {
            s.Nat = &protoNat{
                EnableNatTraversal: util.YesNo(e.EnableNatTraversal),
                NatTraversalKeepAlive: e.NatTraversalKeepAlive,
                NatTraversalEnableUdpChecksum: util.YesNo(e.NatTraversalEnableUdpChecksum),
            }
        }

        if e.EnableFragmentation {
            s.Frag = &protoFrag{
                EnableFragmentation: util.YesNo(e.EnableFragmentation),
            }
        }

        ans.ProtoCommon = &s
    }

    return ans
}
