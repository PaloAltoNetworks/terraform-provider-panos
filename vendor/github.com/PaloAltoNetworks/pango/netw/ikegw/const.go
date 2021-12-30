package ikegw

const (
	Ikev1          = "ikev1"
	Ikev2          = "ikev2"
	Ikev2Preferred = "ikev2-preferred"
)

const (
	IdTypeIpAddress = "ipaddr"
	IdTypeFqdn      = "fqdn"
	IdTypeUfqdn     = "ufqdn"
	IdTypeKeyId     = "keyid"
	IdTypeDn        = "dn"
)

const (
	PeerTypeIp      = "ip"
	PeerTypeDynamic = "dynamic"
	PeerTypeFqdn    = "fqdn"
)

const (
	LocalTypeIp         = "ip"
	LocalTypeFloatingIp = "floating-ip"
)

const (
	AuthPreSharedKey = "pre-shared-key"
	AuthCertificate  = "certificate"
)

const (
	PeerIdCheckExact    = "exact"
	PeerIdCheckWildcard = "wildcard"
)

const (
	singular = "ike gateways"
	plural   = "ike gateways"
)
