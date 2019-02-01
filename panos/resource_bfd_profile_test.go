package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/bfd"
	"github.com/PaloAltoNetworks/pango/version"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosBdfProfile_basic(t *testing.T) {
	versionAdded := version.Number{7, 1, 0, ""}

	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	} else if !testAccPanosVersion.Gte(versionAdded) {
		t.Skip("This test is only valid for PAN-OS 7.1+")
	}

	var o bfd.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBfdProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBfdProfileConfig(name, bfd.ModeActive, 201, 202, 3, 4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBfdProfileExists("panos_bfd_profile.test", &o),
					testAccCheckPanosBfdProfileAttributes(&o, name, bfd.ModeActive, 201, 202, 3, 4),
				),
			},
			{
				Config: testAccBfdProfileConfig(name, bfd.ModePassive, 301, 302, 33, 44),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBfdProfileExists("panos_bfd_profile.test", &o),
					testAccCheckPanosBfdProfileAttributes(&o, name, bfd.ModePassive, 301, 302, 33, 44),
				),
			},
		},
	})
}

func testAccCheckPanosBfdProfileExists(n string, o *bfd.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		name := rs.Primary.ID
		v, err := fw.Network.BfdProfile.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBfdProfileAttributes(o *bfd.Entry, name, mode string, txi, rxi, dm, ht int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Mode != mode {
			return fmt.Errorf("Mode is %q, expected %q", o.Mode, mode)
		}

		if o.MinimumTxInterval != txi {
			return fmt.Errorf("Min TX Interval is %d, expected %d", o.MinimumTxInterval, txi)
		}

		if o.MinimumRxInterval != rxi {
			return fmt.Errorf("Min RX Interval is %d, expected %d", o.MinimumRxInterval, txi)
		}

		if o.DetectionMultiplier != dm {
			return fmt.Errorf("Detection multiplier is %d, expected %d", o.DetectionMultiplier, dm)
		}

		if o.HoldTime != ht {
			return fmt.Errorf("Hold time is %d, expected %d", o.HoldTime, ht)
		}

		return nil
	}
}

func testAccPanosBfdProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bfd_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			_, err := fw.Network.BfdProfile.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBfdProfileConfig(name, mode string, txi, rxi, dm, ht int) string {
	return fmt.Sprintf(`
resource "panos_bfd_profile" "test" {
    name = %q
    mode = %q
    minimum_tx_interval = %d
    minimum_rx_interval = %d
    detection_multiplier = %d
    hold_time = %d
}
`, name, mode, txi, rxi, dm, ht)
}
