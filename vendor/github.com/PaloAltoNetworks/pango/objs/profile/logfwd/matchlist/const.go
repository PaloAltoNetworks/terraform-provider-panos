package matchlist

// These are valid values for LogType.
// The value "sctp" is valid for PAN-OS 8.1+.
// The value "decryption" is valid for PAN-OS 10.0+.
const (
	LogTypeTraffic    = "traffic"
	LogTypeThreat     = "threat"
	LogTypeWildfire   = "wildfire"
	LogTypeUrl        = "url"
	LogTypeData       = "data"
	LogTypeGtp        = "gtp"
	LogTypeTunnel     = "tunnel"
	LogTypeAuth       = "auth"
	LogTypeSctp       = "sctp"
	LogTypeDecryption = "decryption"
)

const (
	singular = "log forwarding profile match list"
	plural   = "log forwarding profile match lists"
)
