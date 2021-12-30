package telemetry

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Config is a normalized, version independent representation of telemetry
// sharing configuration.
type Config struct {
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
func (o *Config) Copy(s Config) {
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

func (o Config) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return "", fn(o)
}

type normalizer interface {
	Normalize() []Config
	Names() []string
}

type container_v1 struct {
	Answer []config_v1 `xml:"statistics-service"`
}

func (o *container_v1) Normalize() []Config {
	ans := make([]Config, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for _ = range o.Answer {
		ans = append(ans, "")
	}

	return ans
}

func (o *config_v1) normalize() Config {
	ans := Config{
		ApplicationReports:             util.AsBool(o.ApplicationReports),
		ThreatPreventionReports:        util.AsBool(o.ThreatPreventionReports),
		UrlReports:                     util.AsBool(o.UrlReports),
		FileTypeIdentificationReports:  util.AsBool(o.FileTypeIdentificationReports),
		ThreatPreventionData:           util.AsBool(o.ThreatPreventionData),
		ThreatPreventionPacketCaptures: util.AsBool(o.ThreatPreventionPacketCaptures),
		ProductUsageStats:              util.AsBool(o.ProductUsageStats),
		PassiveDnsMonitoring:           util.AsBool(o.PassiveDnsMonitoring),
	}

	return ans
}

type config_v1 struct {
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

func specify_v1(e Config) interface{} {
	ans := config_v1{
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
