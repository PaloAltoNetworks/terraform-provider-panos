package nat

// Values for Entry.SatType.
const (
	DynamicIpAndPort = "dynamic-ip-and-port"
	DynamicIp        = "dynamic-ip"
	StaticIp         = "static-ip"
)

// Values for Entry.SatAddressType.
const (
	InterfaceAddress  = "interface-address"
	TranslatedAddress = "translated-address"
)

// None is a valid value for both Entry.SatType and Entry.SatAddressType.
const None = "none"

// These are the valid settings for Entry.SatFallbackIpType.
const (
	Ip         = "ip"
	FloatingIp = "floating"
)

// These are valid settings for DatType.
const (
	DatTypeStatic  = "destination-translation"
	DatTypeDynamic = "dynamic-destination-translation"
)

// Valid values for the Type value.
const (
	TypeIpv4  = "ipv4"
	TypeNat64 = "nat64"
	TypeNptv6 = "nptv6"
)

const (
	singular = "nat rule"
	plural   = "nat rules"
)
