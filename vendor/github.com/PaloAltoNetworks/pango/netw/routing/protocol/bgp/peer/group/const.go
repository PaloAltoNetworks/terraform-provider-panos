package group

// Valid values for Type.
const (
	TypeEbgp       = "ebgp"
	TypeEbgpConfed = "ebgp-confed"
	TypeIbgp       = "ibgp"
	TypeIbgpConfed = "ibgp-confed"
)

// Valid values for ExportNextHop and ImportNextHop.
// NextHopResolve is valid only for ExportNextHop for TypeEbgp.
// NextHopUsePeer is valid only for ImportNextHop for TypeEbgp.
const (
	NextHopOriginal = "original"
	NextHopUseSelf  = "use-self"
	NextHopResolve  = "resolve"
	NextHopUsePeer  = "use-peer"
)

const (
	singular = "bgp peer group"
	plural   = "bgp peer groups"
)
