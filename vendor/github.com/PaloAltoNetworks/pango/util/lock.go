package util

import (
	"encoding/xml"
)

// Lock represents either a config lock or a commit lock.
type Lock struct {
	XMLName  xml.Name  `xml:"entry"`
	Owner    string    `xml:"name,attr"`
	Name     string    `xml:"name"`
	Type     string    `xml:"type"`
	LoggedIn string    `xml:"loggedin"`
	Comment  CdataText `xml:"comment"`
}
