package zone

// These are valid values for the Mode.
const (
	ModeL2          = "layer2"
	ModeL3          = "layer3"
	ModeVirtualWire = "virtual-wire"
	ModeTap         = "tap"
	ModeExternal    = "external"
	ModeTunnel      = "tunnel" // 8.0+
)

const (
	singular = "zone"
	plural   = "zones"
)
