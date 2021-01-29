package aggregate

// Valid Mode values.
const (
	ModeHa            = "ha"
	ModeDecryptMirror = "decrypt-mirror"
	ModeVirtualWire   = "virtual-wire"
	ModeLayer2        = "layer2"
	ModeLayer3        = "layer3"
)

// Valid values for LacpMode.
const (
	LacpModePassive = "passive"
	LacpModeActive  = "active"
)

// Valid values for LacpTransmissionRate.
const (
	LacpTransmissionRateFast = "fast"
	LacpTransmissionRateSlow = "slow"
)

const (
	singular = "aggregate ethernet interface"
	plural   = "aggregate ethernet interfaces"
)
