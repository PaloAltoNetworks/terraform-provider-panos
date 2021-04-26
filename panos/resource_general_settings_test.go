package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/general"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosGeneralSettings_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o general.Config

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccGeneralSettingsConfig("acctest", "10.15.15.15", "10.5.5.5", "10.10.10.10", "autokey"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosGeneralSettingsExists("panos_general_settings.test", &o),
					testAccCheckPanosGeneralSettingsAttributes(&o, "acctest", "10.15.15.15", "10.5.5.5", "10.10.10.10", "autokey"),
				),
			},
			{
				Config: testAccGeneralSettingsConfig("ngfw", "10.25.25.25", "10.15.15.15", "10.20.20.20", "none"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosGeneralSettingsExists("panos_general_settings.test", &o),
					testAccCheckPanosGeneralSettingsAttributes(&o, "ngfw", "10.25.25.25", "10.15.15.15", "10.20.20.20", "none"),
				),
			},
		},
	})
}

func testAccCheckPanosGeneralSettingsExists(n string, o *general.Config) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		v, err := fw.Device.GeneralSettings.Get()
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosGeneralSettingsAttributes(o *general.Config, h, ps, ds, nsa, nsat string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Hostname != h {
			return fmt.Errorf("Hostname is %s, expected %s", o.Hostname, h)
		}

		if o.PanoramaSecondary != ds {
			return fmt.Errorf("Secondary Panorama is %s, expected %s", o.PanoramaSecondary, ps)
		}

		if o.DnsSecondary != ds {
			return fmt.Errorf("Secondary DNS is %s, expected %s", o.DnsSecondary, ds)
		}

		if o.NtpSecondaryAddress != nsa {
			return fmt.Errorf("Secondary NTP Address is %s, expected %s", o.NtpSecondaryAddress, nsa)
		}

		if o.NtpSecondaryAuthType != nsat {
			return fmt.Errorf("Hostname is %s, expected %s", o.NtpSecondaryAuthType, nsat)
		}

		return nil
	}
}

func testAccGeneralSettingsConfig(h, ps, ds, nsa, nsat string) string {
	return fmt.Sprintf(`
resource "panos_general_settings" "test" {
    hostname = "%s"
    panorama_secondary = "%s"
    dns_secondary = "%s"
    ntp_secondary_address = "%s"
    ntp_secondary_auth_type = "%s"
}
`, h, ps, ds, nsa, nsat)
}
