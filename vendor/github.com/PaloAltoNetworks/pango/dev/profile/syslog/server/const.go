package server

// Valid values for Transport.
const (
	TransportUdp = "UDP"
	TransportTcp = "TCP"
	TransportSsl = "SSL"
)

// Valid values for SyslogFormat.
const (
	SyslogFormatBsd  = "BSD"
	SyslogFormatIetf = "IETF"
)

// Valid facilities.
const (
	FacilityUser   = "LOG_USER"
	FacilityLocal0 = "LOG_LOCAL0"
	FacilityLocal1 = "LOG_LOCAL1"
	FacilityLocal2 = "LOG_LOCAL2"
	FacilityLocal3 = "LOG_LOCAL3"
	FacilityLocal4 = "LOG_LOCAL4"
	FacilityLocal5 = "LOG_LOCAL5"
	FacilityLocal6 = "LOG_LOCAL6"
	FacilityLocal7 = "LOG_LOCAL7"
)

const (
	singular = "syslog server"
	plural   = "syslog servers"
)
