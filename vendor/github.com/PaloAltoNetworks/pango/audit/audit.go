package audit

import (
	"encoding/xml"
	"time"

	"github.com/PaloAltoNetworks/pango/util"
)

// CommentHistory is a container for audit comment history.
type CommentHistory struct {
	Comments []Comment `xml:"result>log>logs>entry"`
}

// Comment is a historical audit comment.
//
// The time generated returned from PAN-OS doesn't contain timezone
// information, so the `TimeGenerated` is being left as a string for now.
type Comment struct {
	Admin         string    `xml:"admin"`
	Comment       string    `xml:"comment"`
	ConfigVersion int       `xml:"config_ver"`
	TimeGenerated string    `xml:"time_generated"`
	Time          time.Time `xml:"-"`
}

// SetTime attempts to parse the `TimeGenerated` into a usable time.Time while
// pulling the timezone information from the system clock.
func (c *Comment) SetTime(clock time.Time) {
	if t, err := time.ParseInLocation(util.PanosTimeWithoutTimezoneFormat, c.TimeGenerated, clock.Location()); err == nil {
		c.Time = t
	}
}

// SetComment is for configuring an audit comment for the given XPATH.
type SetComment struct {
	XMLName xml.Name `xml:"set"`
	Xpath   string   `xml:"audit-comment>xpath"`
	Comment string   `xml:"audit-comment>comment"`
}

// GetComment is a query to get the current audit comment.
type GetComment struct {
	XMLName xml.Name `xml:"show"`
	Xpath   string   `xml:"config>list>audit-comments>xpath"`
}

// UncommittedComment is returned when getting the current audit comment for a rule.
type UncommittedComment struct {
	Comment string `xml:"result>entry>comment"`
}
