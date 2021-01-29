package neighbor

// Valid values for RaRouterPreference.
const (
	RaRouterPreferenceHigh   = "High"
	RaRouterPreferenceMedium = "Medium"
	RaRouterPreferenceLow    = "Low"
)

// Valid values for the iType param.
const (
	TypeEthernet  = "ethernet"
	TypeAggregate = "aggregate-ethernet"
	TypeVlan      = "vlan"
)

const (
	singular = "ipv6 neighbor"
)
