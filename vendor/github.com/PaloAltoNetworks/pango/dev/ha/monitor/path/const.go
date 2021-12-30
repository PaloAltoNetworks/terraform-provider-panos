package path

const (
	// FailureCondition Values
	FailureConditionAny = "any"
	FailureConditionAll = "all"

	// gType Values
	VirtualWire   = "virtual-wire"
	Vlan          = "vlan"
	VirtualRouter = "virtual-router"
	LogicalRouter = "logical-router"
)

const (
	singular = "HA path monitor group"
	plural   = singular + "s"
)
