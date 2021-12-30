package app

import (
	"encoding/xml"

	"github.com/PaloAltoNetworks/pango/util"
	"github.com/PaloAltoNetworks/pango/version"
)

// Entry is a normalized, version independent representation of an application.
type Entry struct {
	Name                                 string
	DefaultType                          string
	DefaultPorts                         []string // ordered
	DefaultIpProtocol                    string
	DefaultIcmpType                      int
	DefaultIcmpCode                      int
	Category                             string
	Subcategory                          string
	Technology                           string
	Description                          string
	Timeout                              int
	TcpTimeout                           int
	UdpTimeout                           int
	TcpHalfClosedTimeout                 int
	TcpTimeWaitTimeout                   int
	Risk                                 int
	AbleToFileTransfer                   bool
	ExcessiveBandwidth                   bool
	TunnelsOtherApplications             bool
	HasKnownVulnerability                bool
	UsedByMalware                        bool
	EvasiveBehavior                      bool
	PervasiveUse                         bool
	ProneToMisuse                        bool
	ContinueScanningForOtherApplications bool
	FileTypeIdent                        bool
	VirusIdent                           bool
	DataIdent                            bool
	AlgDisableCapability                 string
	ParentApp                            string
	NoAppIdCaching                       bool // 8.1+

	raw map[string]string
}

// Copy copies the information from source Entry `s` to this object.  As the
// Name field relates to the XPATH of this object, this field is not copied.
func (o *Entry) Copy(s Entry) {
	o.DefaultType = s.DefaultType
	if s.DefaultPorts == nil {
		o.DefaultPorts = nil
	} else {
		o.DefaultPorts = make([]string, len(s.DefaultPorts))
		copy(o.DefaultPorts, s.DefaultPorts)
	}
	o.DefaultIpProtocol = s.DefaultIpProtocol
	o.DefaultIcmpType = s.DefaultIcmpType
	o.DefaultIcmpCode = s.DefaultIcmpCode
	o.Category = s.Category
	o.Subcategory = s.Subcategory
	o.Technology = s.Technology
	o.Description = s.Description
	o.Timeout = s.Timeout
	o.TcpTimeout = s.TcpTimeout
	o.UdpTimeout = s.UdpTimeout
	o.TcpHalfClosedTimeout = s.TcpHalfClosedTimeout
	o.TcpTimeWaitTimeout = s.TcpTimeWaitTimeout
	o.Risk = s.Risk
	o.AbleToFileTransfer = s.AbleToFileTransfer
	o.ExcessiveBandwidth = s.ExcessiveBandwidth
	o.TunnelsOtherApplications = s.TunnelsOtherApplications
	o.HasKnownVulnerability = s.HasKnownVulnerability
	o.UsedByMalware = s.UsedByMalware
	o.EvasiveBehavior = s.EvasiveBehavior
	o.PervasiveUse = s.PervasiveUse
	o.ProneToMisuse = s.ProneToMisuse
	o.ContinueScanningForOtherApplications = s.ContinueScanningForOtherApplications
	o.FileTypeIdent = s.FileTypeIdent
	o.VirusIdent = s.VirusIdent
	o.DataIdent = s.DataIdent
	o.AlgDisableCapability = s.AlgDisableCapability
	o.ParentApp = s.ParentApp
	o.NoAppIdCaching = s.NoAppIdCaching
}

/** Structs / functions for this namespace. **/

func (o Entry) Specify(v version.Number) (string, interface{}) {
	_, fn := versioning(v)
	return o.Name, fn(o)
}

type normalizer interface {
	Normalize() []Entry
	Names() []string
}

type container_v1 struct {
	Answer []entry_v1 `xml:"entry"`
}

func (o *container_v1) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v1) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v1) normalize() Entry {
	ans := Entry{
		Name:                                 o.Name,
		Category:                             o.Category,
		Subcategory:                          o.Subcategory,
		Technology:                           o.Technology,
		Description:                          o.Description,
		Timeout:                              o.Timeout,
		TcpTimeout:                           o.TcpTimeout,
		UdpTimeout:                           o.UdpTimeout,
		TcpHalfClosedTimeout:                 o.TcpHalfClosedTimeout,
		TcpTimeWaitTimeout:                   o.TcpTimeWaitTimeout,
		Risk:                                 o.Risk,
		AbleToFileTransfer:                   util.AsBool(o.AbleToFileTransfer),
		ExcessiveBandwidth:                   util.AsBool(o.ExcessiveBandwidth),
		TunnelsOtherApplications:             util.AsBool(o.TunnelsOtherApplications),
		HasKnownVulnerability:                util.AsBool(o.HasKnownVulnerability),
		UsedByMalware:                        util.AsBool(o.UsedByMalware),
		EvasiveBehavior:                      util.AsBool(o.EvasiveBehavior),
		PervasiveUse:                         util.AsBool(o.PervasiveUse),
		ProneToMisuse:                        util.AsBool(o.ProneToMisuse),
		ContinueScanningForOtherApplications: util.AsBool(o.ContinueScanningForOtherApplications),
		FileTypeIdent:                        util.AsBool(o.FileTypeIdent),
		VirusIdent:                           util.AsBool(o.VirusIdent),
		DataIdent:                            util.AsBool(o.DataIdent),
		AlgDisableCapability:                 o.AlgDisableCapability,
		ParentApp:                            o.ParentApp,
	}

	raw := make(map[string]string)

	if o.Default == nil {
		ans.DefaultType = DefaultTypeNone
	} else if o.Default.DefaultPorts != nil {
		ans.DefaultType = DefaultTypePort
		ans.DefaultPorts = util.MemToStr(o.Default.DefaultPorts)
	} else if o.Default.DefaultIpProtocol != "" {
		ans.DefaultType = DefaultTypeIpProtocol
		ans.DefaultIpProtocol = o.Default.DefaultIpProtocol
	} else if o.Default.Icmp != nil {
		ans.DefaultType = DefaultTypeIcmp
		ans.DefaultIcmpType = o.Default.Icmp.Type
		ans.DefaultIcmpCode = o.Default.Icmp.Code
	} else if o.Default.Icmp6 != nil {
		ans.DefaultType = DefaultTypeIcmp6
		ans.DefaultIcmpType = o.Default.Icmp6.Type
		ans.DefaultIcmpCode = o.Default.Icmp6.Code
	}

	if o.Sigs != nil {
		raw["sigs"] = util.CleanRawXml(o.Sigs.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v1 struct {
	XMLName                              xml.Name     `xml:"entry"`
	Name                                 string       `xml:"name,attr"`
	Default                              *theDefault  `xml:"default"`
	Category                             string       `xml:"category"`
	Subcategory                          string       `xml:"subcategory"`
	Technology                           string       `xml:"technology"`
	Description                          string       `xml:"description,omitempty"`
	Timeout                              int          `xml:"timeout,omitempty"`
	TcpTimeout                           int          `xml:"tcp-timeout,omitempty"`
	UdpTimeout                           int          `xml:"udp-timeout,omitempty"`
	TcpHalfClosedTimeout                 int          `xml:"tcp-half-closed-timeout,omitempty"`
	TcpTimeWaitTimeout                   int          `xml:"tcp-time-wait-timeout,omitempty"`
	Risk                                 int          `xml:"risk"`
	AbleToFileTransfer                   string       `xml:"able-to-transfer-file"`
	ExcessiveBandwidth                   string       `xml:"consume-big-bandwidth"`
	TunnelsOtherApplications             string       `xml:"tunnel-other-application"`
	HasKnownVulnerability                string       `xml:"has-known-vulnerability"`
	UsedByMalware                        string       `xml:"used-by-malware"`
	EvasiveBehavior                      string       `xml:"evasive-behavior"`
	PervasiveUse                         string       `xml:"pervasive-use"`
	ProneToMisuse                        string       `xml:"prone-to-misuse"`
	ContinueScanningForOtherApplications string       `xml:"tunnel-applications"`
	FileTypeIdent                        string       `xml:"file-type-ident"`
	VirusIdent                           string       `xml:"virus-ident"`
	DataIdent                            string       `xml:"data-ident"`
	AlgDisableCapability                 string       `xml:"alg-disable-capability,omitempty"`
	ParentApp                            string       `xml:"parent-app,omitempty"`
	Sigs                                 *util.RawXml `xml:"signature"`
}

type theDefault struct {
	DefaultPorts      *util.MemberType `xml:"port"`
	DefaultIpProtocol string           `xml:"ident-by-ip-protocol,omitempty"`
	Icmp              *icmp            `xml:"ident-by-icmp-type"`
	Icmp6             *icmp            `xml:"ident-by-icmp6-type"`
}

type icmp struct {
	Type int `xml:"type"`
	Code int `xml:"code"`
}

func specify_v1(e Entry) interface{} {
	ans := entry_v1{
		Name:                                 e.Name,
		Category:                             e.Category,
		Subcategory:                          e.Subcategory,
		Technology:                           e.Technology,
		Description:                          e.Description,
		Timeout:                              e.Timeout,
		TcpTimeout:                           e.TcpTimeout,
		UdpTimeout:                           e.UdpTimeout,
		TcpHalfClosedTimeout:                 e.TcpHalfClosedTimeout,
		TcpTimeWaitTimeout:                   e.TcpTimeWaitTimeout,
		Risk:                                 e.Risk,
		AbleToFileTransfer:                   util.YesNo(e.AbleToFileTransfer),
		ExcessiveBandwidth:                   util.YesNo(e.ExcessiveBandwidth),
		TunnelsOtherApplications:             util.YesNo(e.TunnelsOtherApplications),
		HasKnownVulnerability:                util.YesNo(e.HasKnownVulnerability),
		UsedByMalware:                        util.YesNo(e.UsedByMalware),
		EvasiveBehavior:                      util.YesNo(e.EvasiveBehavior),
		PervasiveUse:                         util.YesNo(e.PervasiveUse),
		ProneToMisuse:                        util.YesNo(e.ProneToMisuse),
		ContinueScanningForOtherApplications: util.YesNo(e.ContinueScanningForOtherApplications),
		FileTypeIdent:                        util.YesNo(e.FileTypeIdent),
		VirusIdent:                           util.YesNo(e.VirusIdent),
		DataIdent:                            util.YesNo(e.DataIdent),
		AlgDisableCapability:                 e.AlgDisableCapability,
		ParentApp:                            e.ParentApp,
	}

	switch e.DefaultType {
	case DefaultTypePort:
		ans.Default = &theDefault{
			DefaultPorts: util.StrToMem(e.DefaultPorts),
		}
	case DefaultTypeIpProtocol:
		ans.Default = &theDefault{
			DefaultIpProtocol: e.DefaultIpProtocol,
		}
	case DefaultTypeIcmp:
		ans.Default = &theDefault{
			Icmp: &icmp{
				Type: e.DefaultIcmpType,
				Code: e.DefaultIcmpCode,
			},
		}
	case DefaultTypeIcmp6:
		ans.Default = &theDefault{
			Icmp6: &icmp{
				Type: e.DefaultIcmpType,
				Code: e.DefaultIcmpCode,
			},
		}
	}

	if text := e.raw["sigs"]; text != "" {
		ans.Sigs = &util.RawXml{text}
	}

	return ans
}

// PAN-OS 8.1
type container_v2 struct {
	Answer []entry_v2 `xml:"entry"`
}

func (o *container_v2) Normalize() []Entry {
	ans := make([]Entry, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].normalize())
	}

	return ans
}

func (o *container_v2) Names() []string {
	ans := make([]string, 0, len(o.Answer))
	for i := range o.Answer {
		ans = append(ans, o.Answer[i].Name)
	}

	return ans
}

func (o *entry_v2) normalize() Entry {
	ans := Entry{
		Name:                                 o.Name,
		Category:                             o.Category,
		Subcategory:                          o.Subcategory,
		Technology:                           o.Technology,
		Description:                          o.Description,
		Timeout:                              o.Timeout,
		TcpTimeout:                           o.TcpTimeout,
		UdpTimeout:                           o.UdpTimeout,
		TcpHalfClosedTimeout:                 o.TcpHalfClosedTimeout,
		TcpTimeWaitTimeout:                   o.TcpTimeWaitTimeout,
		Risk:                                 o.Risk,
		AbleToFileTransfer:                   util.AsBool(o.AbleToFileTransfer),
		ExcessiveBandwidth:                   util.AsBool(o.ExcessiveBandwidth),
		TunnelsOtherApplications:             util.AsBool(o.TunnelsOtherApplications),
		HasKnownVulnerability:                util.AsBool(o.HasKnownVulnerability),
		UsedByMalware:                        util.AsBool(o.UsedByMalware),
		EvasiveBehavior:                      util.AsBool(o.EvasiveBehavior),
		PervasiveUse:                         util.AsBool(o.PervasiveUse),
		ProneToMisuse:                        util.AsBool(o.ProneToMisuse),
		ContinueScanningForOtherApplications: util.AsBool(o.ContinueScanningForOtherApplications),
		FileTypeIdent:                        util.AsBool(o.FileTypeIdent),
		VirusIdent:                           util.AsBool(o.VirusIdent),
		DataIdent:                            util.AsBool(o.DataIdent),
		AlgDisableCapability:                 o.AlgDisableCapability,
		ParentApp:                            o.ParentApp,
		NoAppIdCaching:                       util.AsBool(o.NoAppIdCaching),
	}

	raw := make(map[string]string)

	if o.Default == nil {
		ans.DefaultType = DefaultTypeNone
	} else if o.Default.DefaultPorts != nil {
		ans.DefaultType = DefaultTypePort
		ans.DefaultPorts = util.MemToStr(o.Default.DefaultPorts)
	} else if o.Default.DefaultIpProtocol != "" {
		ans.DefaultType = DefaultTypeIpProtocol
		ans.DefaultIpProtocol = o.Default.DefaultIpProtocol
	} else if o.Default.Icmp != nil {
		ans.DefaultType = DefaultTypeIcmp
		ans.DefaultIcmpType = o.Default.Icmp.Type
		ans.DefaultIcmpCode = o.Default.Icmp.Code
	} else if o.Default.Icmp6 != nil {
		ans.DefaultType = DefaultTypeIcmp6
		ans.DefaultIcmpType = o.Default.Icmp6.Type
		ans.DefaultIcmpCode = o.Default.Icmp6.Code
	}

	if o.Sigs != nil {
		raw["sigs"] = util.CleanRawXml(o.Sigs.Text)
	}

	if len(raw) != 0 {
		ans.raw = raw
	}

	return ans
}

type entry_v2 struct {
	XMLName                              xml.Name     `xml:"entry"`
	Name                                 string       `xml:"name,attr"`
	Default                              *theDefault  `xml:"default"`
	Category                             string       `xml:"category"`
	Subcategory                          string       `xml:"subcategory"`
	Technology                           string       `xml:"technology"`
	Description                          string       `xml:"description,omitempty"`
	Timeout                              int          `xml:"timeout,omitempty"`
	TcpTimeout                           int          `xml:"tcp-timeout,omitempty"`
	UdpTimeout                           int          `xml:"udp-timeout,omitempty"`
	TcpHalfClosedTimeout                 int          `xml:"tcp-half-closed-timeout,omitempty"`
	TcpTimeWaitTimeout                   int          `xml:"tcp-time-wait-timeout,omitempty"`
	Risk                                 int          `xml:"risk"`
	AbleToFileTransfer                   string       `xml:"able-to-transfer-file"`
	ExcessiveBandwidth                   string       `xml:"consume-big-bandwidth"`
	TunnelsOtherApplications             string       `xml:"tunnel-other-application"`
	HasKnownVulnerability                string       `xml:"has-known-vulnerability"`
	UsedByMalware                        string       `xml:"used-by-malware"`
	EvasiveBehavior                      string       `xml:"evasive-behavior"`
	PervasiveUse                         string       `xml:"pervasive-use"`
	ProneToMisuse                        string       `xml:"prone-to-misuse"`
	ContinueScanningForOtherApplications string       `xml:"tunnel-applications"`
	FileTypeIdent                        string       `xml:"file-type-ident"`
	VirusIdent                           string       `xml:"virus-ident"`
	DataIdent                            string       `xml:"data-ident"`
	AlgDisableCapability                 string       `xml:"alg-disable-capability,omitempty"`
	ParentApp                            string       `xml:"parent-app,omitempty"`
	NoAppIdCaching                       string       `xml:"no-appid-caching"`
	Sigs                                 *util.RawXml `xml:"signature"`
}

func specify_v2(e Entry) interface{} {
	ans := entry_v2{
		Name:                                 e.Name,
		Category:                             e.Category,
		Subcategory:                          e.Subcategory,
		Technology:                           e.Technology,
		Description:                          e.Description,
		Timeout:                              e.Timeout,
		TcpTimeout:                           e.TcpTimeout,
		UdpTimeout:                           e.UdpTimeout,
		TcpHalfClosedTimeout:                 e.TcpHalfClosedTimeout,
		TcpTimeWaitTimeout:                   e.TcpTimeWaitTimeout,
		Risk:                                 e.Risk,
		AbleToFileTransfer:                   util.YesNo(e.AbleToFileTransfer),
		ExcessiveBandwidth:                   util.YesNo(e.ExcessiveBandwidth),
		TunnelsOtherApplications:             util.YesNo(e.TunnelsOtherApplications),
		HasKnownVulnerability:                util.YesNo(e.HasKnownVulnerability),
		UsedByMalware:                        util.YesNo(e.UsedByMalware),
		EvasiveBehavior:                      util.YesNo(e.EvasiveBehavior),
		PervasiveUse:                         util.YesNo(e.PervasiveUse),
		ProneToMisuse:                        util.YesNo(e.ProneToMisuse),
		ContinueScanningForOtherApplications: util.YesNo(e.ContinueScanningForOtherApplications),
		FileTypeIdent:                        util.YesNo(e.FileTypeIdent),
		VirusIdent:                           util.YesNo(e.VirusIdent),
		DataIdent:                            util.YesNo(e.DataIdent),
		AlgDisableCapability:                 e.AlgDisableCapability,
		ParentApp:                            e.ParentApp,
		NoAppIdCaching:                       util.YesNo(e.NoAppIdCaching),
	}

	switch e.DefaultType {
	case DefaultTypePort:
		ans.Default = &theDefault{
			DefaultPorts: util.StrToMem(e.DefaultPorts),
		}
	case DefaultTypeIpProtocol:
		ans.Default = &theDefault{
			DefaultIpProtocol: e.DefaultIpProtocol,
		}
	case DefaultTypeIcmp:
		ans.Default = &theDefault{
			Icmp: &icmp{
				Type: e.DefaultIcmpType,
				Code: e.DefaultIcmpCode,
			},
		}
	case DefaultTypeIcmp6:
		ans.Default = &theDefault{
			Icmp6: &icmp{
				Type: e.DefaultIcmpType,
				Code: e.DefaultIcmpCode,
			},
		}
	}

	if text := e.raw["sigs"]; text != "" {
		ans.Sigs = &util.RawXml{text}
	}

	return ans
}
