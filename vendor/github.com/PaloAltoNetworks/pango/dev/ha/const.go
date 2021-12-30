package ha

const (
	// Mode Values
	ModeActivePassive = "active-passive"
	ModeActiveActive  = "active-active"

	// ApPassiveLinkState Values
	ApPassiveLinkStateAuto     = "auto"
	ApPassiveLinkStateShutdown = "shutdown"

	// AaSessionOwnerSelection Values
	AaSessionOwnerSelectionPrimaryDevice = "primary-device"
	AaSessionOwnerSelectionFirstPacket   = "first-packet"

	// AaFpSessionSetup Values
	AaFpSessionSetupPrimaryDevice = "primary-device"
	AaFpSessionSetupFirstPacket   = "first-packet"
	AaFpSessionSetupIpModulo      = "ip-modulo"
	AaFpSessionSetupIpHash        = "ip-hash"

	// AaFpSessionSetupIpHashKey Values
	AaFpSessionSetupIpHashKeySource     = "source"
	AaFpSessionSetupIpHashKeySourceDest = "source-and-destination"

	// ElectionTimersMode Values
	ElectionTimersModeRecommended = "recommended"
	ElectionTimersModeAggressive  = "aggressive"
	ElectionTimersModeAdvanced    = "advanced"

	// ElectionTimersAdvFlapMax Values
	ElectionTimersAdvFlapMaxInfinite = "infinite"
	ElectionTimersAdvFlapMaxDisable  = "disable"

	// Ha2StateSyncTransport Values
	Ha2StateSyncTransportEthernet = "ethernet"
	Ha2StateSyncTransportIp       = "ip"
	Ha2StateSyncTransportUdp      = "udp"

	// Ha2StateSyncKeepAliveAction Values
	Ha2StateSyncKeepAliveActionLogOnly       = "log-only"
	Ha2StateSyncKeepAliveActionSplitDatapath = "split-datapath"

	// LinkMonitorFailureCondition Values
	LinkMonitorFailureConditionAny = "any"
	LinkMonitorFailureConditionAll = "all"
)

const (
	singular = "HA config"
)
