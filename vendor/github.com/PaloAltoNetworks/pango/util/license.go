package util

import (
	"encoding/xml"
)

// License defines a license entry.
type License struct {
	XMLName     xml.Name `xml:"entry"`
	Feature     string   `xml:"feature"`
	Description string   `xml:"description"`
	Serial      string   `xml:"serial"`
	Issued      string   `xml:"issued"`
	Expires     string   `xml:"expires"`
	Expired     string   `xml:"expired"`
	AuthCode    string   `xml:"authcode"`
}
