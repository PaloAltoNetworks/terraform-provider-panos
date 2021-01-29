package virus

/*
Valid values for Decoder.Action and ApplicationException.Action.

Block:          PAN-OS 6.1 only.
Drop:           PAN-OS 7.0+
ResetClient:    PAN-OS 7.0+
ResetServer:    PAN-OS 7.0+
ResetBoth:      PAN-OS 7.0+
*/
const (
	Default     = "default"
	Allow       = "allow"
	Alert       = "alert"
	Block       = "block"
	Drop        = "drop"
	ResetClient = "reset-client"
	ResetServer = "reset-server"
	ResetBoth   = "reset-both"
)

const (
	singular = "antivirus security profile"
	plural   = "antivirus security profiles"
)
