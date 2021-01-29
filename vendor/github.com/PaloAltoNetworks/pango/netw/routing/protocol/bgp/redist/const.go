package redist

// Valid values for AddressFamily (PAN-OS 8.0+).
const (
	AddressFamilyIpv4 = "ipv4"
	AddressFamilyIpv6 = "ipv6"
)

// Valid values for RouteTable (PAN-OS 8.0+).
const (
	RouteTableUnicast   = "unicast"
	RouteTableMulticast = "multicast"
	RouteTableBoth      = "both"
)

// Valid values for Origin.
const (
	SetOriginIgp        = "igp"
	SetOriginEgp        = "egp"
	SetOriginIncomplete = "incomplete"
)

const (
	singular = "bgp redistribution rule"
	plural   = "bgp redistribution rules"
)
