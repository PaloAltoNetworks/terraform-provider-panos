package router

// Valid values for EcmpLoadBalanceMethod.
const (
	EcmpLoadBalanceMethodIpModulo           = "ip-modulo"
	EcmpLoadBalanceMethodIpHash             = "ip-hash"
	EcmpLoadBalanceMethodWeightedRoundRobin = "weighted-round-robin"
	EcmpLoadBalanceMethodBalancedRoundRobin = "balanced-round-robin"
)

const (
	singular = "virtual router"
	plural   = "virtual routers"
)
