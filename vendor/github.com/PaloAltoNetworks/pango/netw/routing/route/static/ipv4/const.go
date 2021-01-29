package ipv4

// Valid NextHop values.
const (
	NextHopDiscard   = "discard"
	NextHopIpAddress = "ip-address"
	NextHopNextVr    = "next-vr"
)

// Valid RouteTable values.
const (
	RouteTableNoInstall = "no install"
	RouteTableUnicast   = "unicast"
	RouteTableMulticast = "multicast"
	RouteTableBoth      = "both"
)

const (
	singular = "ipv4 static route"
	plural   = "ipv4 static routes"
)
