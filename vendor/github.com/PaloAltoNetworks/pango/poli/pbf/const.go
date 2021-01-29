package pbf

// Valid FromType values.
const (
	FromTypeZone      = "zone"
	FromTypeInterface = "interface"
)

// Valid ForwardNextHopType values.
const (
	ForwardNextHopTypeIpAddress = "ip-address"
	ForwardNextHopTypeFqdn      = "fqdn"
)

// Valid Action values.
const (
	ActionForward     = "forward"
	ActionVsysForward = "forward-to-vsys"
	ActionDiscard     = "discard"
	ActionNoPbf       = "no-pbf"
)

const (
	singular = "policy based forwarding rule"
	plural   = "policy based forwarding rules"
)
