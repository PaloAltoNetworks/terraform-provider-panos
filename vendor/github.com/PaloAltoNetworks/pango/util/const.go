package util

// Rulebase constants for various policies.
const (
	Rulebase     = "rulebase"
	PreRulebase  = "pre-rulebase"
	PostRulebase = "post-rulebase"
)

// Valid values to use for VsysImport() or VsysUnimport().
const (
	InterfaceImport     = "interface"
	VirtualRouterImport = "virtual-router"
	VirtualWireImport   = "virtual-wire"
	VlanImport          = "vlan"
)

// These constants are valid move locations to pass to various movement
// functions (aka - policy management).
const (
	MoveSkip = iota
	MoveBefore
	MoveDirectlyBefore
	MoveAfter
	MoveDirectlyAfter
	MoveTop
	MoveBottom
)

// Valid values to use for any function expecting a pango query type `qt`.
const (
	Get  = "get"
	Show = "show"
)

// PanosTimeWithoutTimezoneFormat is a time (missing the timezone) that PAN-OS
// will give sometimes.  Combining this with `Clock()` to get a usable time.
// report that does not contain
const PanosTimeWithoutTimezoneFormat = "2006/01/02 15:04:05"
