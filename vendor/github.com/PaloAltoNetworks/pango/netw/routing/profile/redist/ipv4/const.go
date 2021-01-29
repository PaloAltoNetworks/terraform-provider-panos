package ipv4

// These are valid values for the Action param.
const (
	ActionRedist   = "redist"
	ActionNoRedist = "no-redist"
)

// These are valid values for Types.
const (
	TypeBgp     = "bgp"
	TypeConnect = "connect"
	TypeOspf    = "ospf"
	TypeRip     = "rip"
	TypeStatic  = "static"
)

// These are valid values for OspfPathTypes.
const (
	OspfPathTypeIntraArea = "intra-area"
	OspfPathTypeInterArea = "inter-area"
	OspfPathTypeExt1      = "ext-1"
	OspfPathTypeExt2      = "ext-2"
)

const (
	singular = "redistribution profile"
	plural   = "redistribution profiles"
)
