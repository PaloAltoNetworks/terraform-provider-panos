package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/telemetry"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosTelemetry_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o telemetry.Config

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosTelemetryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccTelemetryConfig(true, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTelemetryExists("panos_telemetry.test", &o),
					testAccCheckPanosTelemetryAttributes(&o, true, false, true, false),
				),
			},
			{
				Config: testAccTelemetryConfig(false, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosTelemetryExists("panos_telemetry.test", &o),
					testAccCheckPanosTelemetryAttributes(&o, false, true, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosTelemetryExists(n string, o *telemetry.Config) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		v, err := fw.Device.Telemetry.Get()
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosTelemetryAttributes(o *telemetry.Config, ar, tpr, ur, pdm bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.ApplicationReports != ar {
			return fmt.Errorf("Application reports is %t, expected %t", o.ApplicationReports, ar)
		}

		if o.ThreatPreventionReports != tpr {
			return fmt.Errorf("Threat prevention reports is %t, expected %t", o.ThreatPreventionReports, tpr)
		}

		if o.UrlReports != ur {
			return fmt.Errorf("URL reports is %t, expected %t", o.UrlReports, ur)
		}

		if o.PassiveDnsMonitoring != pdm {
			return fmt.Errorf("Passive DNS monitoring is %t, expected %t", o.PassiveDnsMonitoring, pdm)
		}

		return nil
	}
}

func testAccPanosTelemetryDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_telemetry" {
			continue
		}

		if rs.Primary.ID != "" {
			if _, err := fw.Device.Telemetry.Get(); err == nil {
				return fmt.Errorf("Telemetry still exists")
			}
		}
	}

	return nil
}

func testAccTelemetryConfig(ar, tpr, ur, pdm bool) string {
	return fmt.Sprintf(`
resource "panos_telemetry" "test" {
    application_reports = %t
    threat_prevention_reports = %t
    url_reports = %t
    passive_dns_monitoring = %t
}
`, ar, tpr, ur, pdm)
}
