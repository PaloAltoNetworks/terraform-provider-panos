package header

// Valid logtype values for logtype.
const (
	Config   = "config"
	System   = "system"
	Threat   = "threat"
	Traffic  = "traffic"
	HipMatch = "hip-match"
	Url      = "url"
	Data     = "data"
	Wildfire = "wildfire"
	Tunnel   = "tunnel"
	UserId   = "userid"
	Gtp      = "gtp"
	Auth     = "auth"
	Sctp     = "sctp"
	Iptag    = "iptag"
)

const (
	singular = "http header"
	plural   = "http headers"
)
