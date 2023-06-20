package panos

import (
	"github.com/fpluchorg/pango/dev/telemetry"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTelemetry() *schema.Resource {
	return &schema.Resource{
		Create: createTelemetry,
		Read:   readTelemetry,
		Update: updateTelemetry,
		Delete: deleteTelemetry,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"application_reports": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"threat_prevention_reports": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"url_reports": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"file_type_identification_reports": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"threat_prevention_data": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"threat_prevention_packet_captures": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"product_usage_stats": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"passive_dns_monitoring": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func parseTelemetry(d *schema.ResourceData) telemetry.Config {
	o := telemetry.Config{
		ApplicationReports:             d.Get("application_reports").(bool),
		ThreatPreventionReports:        d.Get("threat_prevention_reports").(bool),
		UrlReports:                     d.Get("url_reports").(bool),
		FileTypeIdentificationReports:  d.Get("file_type_identification_reports").(bool),
		ThreatPreventionData:           d.Get("threat_prevention_data").(bool),
		ThreatPreventionPacketCaptures: d.Get("threat_prevention_packet_captures").(bool),
		ProductUsageStats:              d.Get("product_usage_stats").(bool),
		PassiveDnsMonitoring:           d.Get("passive_dns_monitoring").(bool),
	}

	return o
}

func createTelemetry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	o := parseTelemetry(d)

	if err = fw.Device.Telemetry.Set(o); err != nil {
		return err
	}

	d.SetId(fw.Hostname)
	return readTelemetry(d, meta)
}

func readTelemetry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}

	o, err := fw.Device.Telemetry.Get()
	if err != nil {
		if isObjectNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("application_reports", o.ApplicationReports)
	d.Set("threat_prevention_reports", o.ThreatPreventionReports)
	d.Set("url_reports", o.UrlReports)
	d.Set("file_type_identification_reports", o.FileTypeIdentificationReports)
	d.Set("threat_prevention_data", o.ThreatPreventionData)
	d.Set("threat_prevention_packet_captures", o.ThreatPreventionPacketCaptures)
	d.Set("product_usage_stats", o.ProductUsageStats)
	d.Set("passive_dns_monitoring", o.PassiveDnsMonitoring)

	return nil
}

func updateTelemetry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}
	o := parseTelemetry(d)

	lo, err := fw.Device.Telemetry.Get()
	if err != nil {
		return err
	}
	lo.Copy(o)
	if err = fw.Device.Telemetry.Edit(lo); err != nil {
		return err
	}

	return readTelemetry(d, meta)
}

func deleteTelemetry(d *schema.ResourceData, meta interface{}) error {
	fw, err := firewall(meta, "")
	if err != nil {
		return err
	}

	if err = fw.Device.Telemetry.Delete(); err != nil && !isObjectNotFound(err) {
		return err
	}

	d.SetId("")
	return nil
}
