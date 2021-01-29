package ipv6

// Valid NextHop values.
const (
	NextHopDiscard     = "discard"
	NextHopIpv6Address = "ipv6-address"
	NextHopNextVr      = "next-vr"
	NextHopFqdn        = "fqdn" // 9.0+
)

// Valid RouteTable values.
const (
	RouteTableNoInstall = "no install"
	RouteTableUnicast   = "unicast"
)

const (
	singular = "ipv6 static route"
	plural   = "ipv6 static routes"
)
