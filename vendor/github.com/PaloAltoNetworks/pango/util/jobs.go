package util

import (
	"encoding/xml"
	"strings"
)

// JobResponse parses a XML response that includes a job ID.
type JobResponse struct {
	XMLName xml.Name `xml:"response"`
	Id      uint     `xml:"result>job"`
}

// BasicJob is a struct for parsing minimal information about a submitted
// job to PANOS.
type BasicJob struct {
	XMLName  xml.Name        `xml:"response"`
	Result   string          `xml:"result>job>result"`
	Progress uint            `xml:"result>job>progress"`
	Details  BasicJobDetails `xml:"result>job>details"`
	Devices  []devJob        `xml:"result>job>devices>entry"`
}

type BasicJobDetails struct {
	Lines []LineOrCdata `xml:"line"`
}

func (o *BasicJobDetails) String() string {
	ans := make([]string, 0, len(o.Lines))

	for _, line := range o.Lines {
		if line.Cdata != nil {
			ans = append(ans, strings.TrimSpace(*line.Cdata))
		} else if line.Text != nil {
			ans = append(ans, *line.Text)
		} else {
			ans = append(ans, "huh")
		}
	}

	return strings.Join(ans, " | ")
}

type LineOrCdata struct {
	Cdata *string `xml:",cdata"`
	Text  *string `xml:",chardata"`
}

type devJob struct {
	Serial string `xml:"serial-no"`
	Result string `xml:"result"`
}
