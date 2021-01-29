package imp

// Valid values for MatchRouteTable.
const (
	MatchRouteTableUnicast   = "unicast"
	MatchRouteTableMulticast = "multicast"
	MatchRouteTableBoth      = "both"
)

// Valid values for Action.
const (
	ActionAllow = "allow"
	ActionDeny  = "deny"
)

// Valid values for Origin.
const (
	OriginIgp        = "igp"
	OriginEgp        = "egp"
	OriginIncomplete = "incomplete"
)

// Valid values for CommunityType.
const (
	CommunityTypeNone        = "none"
	CommunityTypeRemoveAll   = "remove-all"
	CommunityTypeRemoveRegex = "remove-regex"
	CommunityTypeAppend      = "append"
	CommunityTypeOverwrite   = "overwrite"
)

// Valid values for CommunityValue when CommunityType is "append" or
// "overwrite".
const (
	AppendNoExport    = "no-export"
	AppendNoAdvertise = "no-advertise"
	AppendLocalAs     = "local-as"
	AppendNoPeer      = "nopeer"
)

// Valid values for AsPathType.
const (
	AsPathTypeNone   = "none"
	AsPathTypeRemove = "remove"
)

const (
	singular = "bgp import rule"
	plural   = "bgp import rules"
)
