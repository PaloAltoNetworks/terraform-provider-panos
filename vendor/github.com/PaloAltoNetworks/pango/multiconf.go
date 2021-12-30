package pango

import (
	"encoding/xml"
	"fmt"
	"strings"

	"github.com/PaloAltoNetworks/pango/util"
)

// MultiConfigure is a container object for making a type=multi-config call.
type MultiConfigure struct {
	XMLName xml.Name `xml:"multi-configure-request"`
	Reqs    []MultiConfigureRequest
}

// IncrementalIds assigns incremental ID numbers to all requests.
//
// Any request that already has an ID is skipped, and the number is discarded.
func (m *MultiConfigure) IncrementalIds() {
	for i := range m.Reqs {
		if m.Reqs[i].Id == "" {
			m.Reqs[i].Id = fmt.Sprintf("%d", i+1)
		}
	}
}

// MultiConfigureRequest is an individual request in a MultiConfigure instance.
//
// These are built up automatically when invoking Client.Set / Client.Edit after
// Client.PrepareMultiConfigure is invoked.
type MultiConfigureRequest struct {
	XMLName xml.Name
	Id      string `xml:"id,attr,omitempty"`
	Xpath   string `xml:"xpath,attr"`
	Data    interface{}
}

// MultiConfigureResponse is a struct to handle the response from multi-config
// commands.
type MultiConfigureResponse struct {
	XMLName xml.Name                     `xml:"response"`
	Status  string                       `xml:"status,attr"`
	Code    int                          `xml:"code,attr"`
	Results []MultiConfigResponseElement `xml:"response"`
}

// Ok returns if there was an error or not.
func (m *MultiConfigureResponse) Ok() bool {
	return m.Status == "success"
}

// Error returns the error if there was one.
func (m *MultiConfigureResponse) Error() string {
	if len(m.Results) == 0 {
		return ""
	}

	r := m.Results[len(m.Results)-1]
	if r.Ok() {
		return ""
	}

	return r.Message()
}

// MultiConfigResponseElement is a single response from a multi-config request.
type MultiConfigResponseElement struct {
	XMLName xml.Name `xml:"response"`
	Status  string   `xml:"status,attr"`
	Code    int      `xml:"code,attr"`
	Id      string   `xml:"id,attr,omitempty"`
	Msg     McreMsg  `xml:"msg"`
}

type McreMsg struct {
	Line    []util.CdataText `xml:"line"`
	Message string           `xml:",chardata"`
}

func (m *MultiConfigResponseElement) Ok() bool {
	return m.Status == "success"
}

func (m *MultiConfigResponseElement) Message() string {
	if len(m.Msg.Line) > 0 {
		var b strings.Builder
		for i := range m.Msg.Line {
			if i != 0 {
				b.WriteString(" | ")
			}
			b.WriteString(strings.TrimSpace(m.Msg.Line[i].Text))
		}
		return b.String()
	}

	return m.Msg.Message
}
