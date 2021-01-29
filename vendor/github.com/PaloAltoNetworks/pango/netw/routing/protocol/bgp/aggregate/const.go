package aggregate

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

// Valid values for AsPathType.  As of PAN-OS 8.1, "prepend" and
// "remove-and-prepend" are disabled.
const (
	AsPathTypeNone             = "none"
	AsPathTypeRemove           = "remove"
	AsPathTypePrepend          = "prepend"
	AsPathTypeRemoveAndPrepend = "remove-and-prepend"
)

const (
	singular = "bgp aggregation policy"
	plural   = "bgp aggregation policies"
)
