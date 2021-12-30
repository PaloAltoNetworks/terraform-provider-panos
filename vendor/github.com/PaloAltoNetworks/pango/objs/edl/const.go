package edl

// Constants for Entry.Type field.  Only TypeIp is valid for PAN-OS 7.0 and
// earlier.
const (
	TypeIp            string = "ip"
	TypeDomain        string = "domain"
	TypeUrl           string = "url"
	TypePredefinedIp  string = "predefined-ip"  // PAN-OS 8.0+
	TypePredefinedUrl string = "predefined-url" // PAN-OS 10.0+
)

// Constants for the Repeat field.
const (
	RepeatEveryFiveMinutes = "every five minutes" // PAN-OS 8.0+
	RepeatHourly           = "hourly"
	RepeatDaily            = "daily"
	RepeatWeekly           = "weekly"
	RepeatMonthly          = "monthly"
)

const (
	singular = "edl"
	plural   = "edls"
)
