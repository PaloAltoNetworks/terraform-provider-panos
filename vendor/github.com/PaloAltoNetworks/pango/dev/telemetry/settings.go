package telemetry

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
)

// Settings is a normalized, version independent representation of telemetry
// sharing configuration.
type Settings struct {
	ApplicationReports             bool
	ThreatPreventionReports        bool
	UrlReports                     bool
	FileTypeIdentificationReports  bool
	ThreatPreventionData           bool
	ThreatPreventionPacketCaptures bool
	ProductUsageStats              bool
	PassiveDnsMonitoring           bool
}

// Copy copies the information from source Settings `s` to this object.
func (o *Settings) Copy(s Settings) {
	o.ApplicationReports = s.ApplicationReports
	o.ThreatPreventionReports = s.ThreatPreventionReports
	o.UrlReports = s.UrlReports
	o.FileTypeIdentificationReports = s.FileTypeIdentificationReports
	o.ThreatPreventionData = s.ThreatPreventionData
	o.ThreatPreventionPacketCaptures = s.ThreatPreventionPacketCaptures
	o.ProductUsageStats = s.ProductUsageStats
	o.PassiveDnsMonitoring = s.PassiveDnsMonitoring
}

/** Structs / functions for normalization. **/

type normalizer interface {
	Normalize() Settings
}

type container_v1 struct {
	Answer entry_v1 `xml:"result>statistics-service"`
}

func (o *container_v1) Normalize() Settings {
	ans := Settings{
		ApplicationReports:             util.AsBool(o.Answer.ApplicationReports),
		ThreatPreventionReports:        util.AsBool(o.Answer.ThreatPreventionReports),
		UrlReports:                     util.AsBool(o.Answer.UrlReports),
		FileTypeIdentificationReports:  util.AsBool(o.Answer.FileTypeIdentificationReports),
		ThreatPreventionData:           util.AsBool(o.Answer.ThreatPreventionData),
		ThreatPreventionPacketCaptures: util.AsBool(o.Answer.ThreatPreventionPacketCaptures),
		ProductUsageStats:              util.AsBool(o.Answer.ProductUsageStats),
		PassiveDnsMonitoring:           util.AsBool(o.Answer.PassiveDnsMonitoring),
	}

	return ans
}

type entry_v1 struct {
	XMLName                        xml.Name `xml:"statistics-service"`
	ApplicationReports             string   `xml:"application-reports"`
	ThreatPreventionReports        string   `xml:"threat-prevention-reports"`
	UrlReports                     string   `xml:"url-reports"`
	FileTypeIdentificationReports  string   `xml:"file-identification-reports"`
	ThreatPreventionData           string   `xml:"threat-prevention-information"`
	ThreatPreventionPacketCaptures string   `xml:"threat-prevention-pcap"`
	ProductUsageStats              string   `xml:"health-performance-reports"`
	PassiveDnsMonitoring           string   `xml:"passive-dns-monitoring"`
}

func specify_v1(e Settings) interface{} {
	ans := entry_v1{
		ApplicationReports:             util.YesNo(e.ApplicationReports),
		ThreatPreventionReports:        util.YesNo(e.ThreatPreventionReports),
		UrlReports:                     util.YesNo(e.UrlReports),
		FileTypeIdentificationReports:  util.YesNo(e.FileTypeIdentificationReports),
		ThreatPreventionData:           util.YesNo(e.ThreatPreventionData),
		ThreatPreventionPacketCaptures: util.YesNo(e.ThreatPreventionPacketCaptures),
		ProductUsageStats:              util.YesNo(e.ProductUsageStats),
		PassiveDnsMonitoring:           util.YesNo(e.PassiveDnsMonitoring),
	}

	return ans
}
