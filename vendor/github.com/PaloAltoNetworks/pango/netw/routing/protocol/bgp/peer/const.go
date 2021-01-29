package peer

// Valid values for ReflectorClient.
const (
	ReflectorClientNonClient    = "non-client"
	ReflectorClientClient       = "client"
	ReflectorClientMeshedClient = "meshed-client"
)

// Valid values for PeeringType.
const (
	PeeringTypeBilateral   = "bilateral"
	PeeringTypeUnspecified = "unspecified"
)

// Valid values for BfdProfile, besides an actual BFD profile's name.
const (
	BfdProfileInherit = "Inherit-vr-global-setting"
	BfdProfileNone    = "None"
)

// Valid values for AddressFamilyType.
const (
	AddressFamilyTypeIpv4 = "ipv4"
	AddressFamilyTypeIpv6 = "ipv6"
)

// Valid non-int value for MaxPrefixes.
const MaxPrefixesUnlimited = "unlimited"

const (
	singular = "bgp peer group peer"
	plural   = "bgp peer group peers"
)
