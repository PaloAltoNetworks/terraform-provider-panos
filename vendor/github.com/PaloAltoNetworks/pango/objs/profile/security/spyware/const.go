package spyware

// Valid values for WhiteList.LogLevel.
const (
	LogLevelDefault       = "default"
	LogLevelNone          = "none"
	LogLevelLow           = "low"
	LogLevelInformational = "informational"
	LogLevelMedium        = "medium"
	LogLevelHigh          = "high"
	LogLevelCritical      = "critical"
)

// Valid values for PacketCapture params.
const (
	Disable         = "disable"
	SinglePacket    = "single-packet"
	ExtendedCapture = "extended-capture"
)

// Valid values for Action params.
const (
	ActionAlert       = "alert"        // BotnetList, BlockList, Rule, Exception
	ActionAllow       = "allow"        // BotnetList, DnsCategory, Rule, Exception
	ActionBlock       = "block"        // BotnetList, DnsCategory
	ActionBlockIp     = "block-ip"     // Rule, Exception
	ActionDefault     = "default"      // DnsCategory, Rule, Exception
	ActionDrop        = "drop"         // Rule, Exception
	ActionResetBoth   = "reset-both"   // Rule, Exception
	ActionResetClient = "reset-client" // Rule, Exception
	ActionResetServer = "reset-server" // Rule, Exception
	ActionSinkhole    = "sinkhole"     // BotnetList, DnsCategory
)

// Valid values for BlockIpTrackBy.
const (
	TrackBySource               = "source"
	TrackBySourceAndDestination = "source-and-destination"
)

const (
	singular = "anti-spyware security profile"
	plural   = "anti-spyware security profiles"
)
