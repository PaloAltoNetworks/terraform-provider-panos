package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/monitor"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosMonitorProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o monitor.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosMonitorProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccMonitorProfileConfig(name, monitor.ActionFailOver, 11, 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosMonitorProfileExists("panos_monitor_profile.test", &o),
					testAccCheckPanosMonitorProfileAttributes(&o, name, monitor.ActionFailOver, 11, 7),
				),
			},
			{
				Config: testAccMonitorProfileConfig(name, monitor.ActionWaitRecover, 10, 4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosMonitorProfileExists("panos_monitor_profile.test", &o),
					testAccCheckPanosMonitorProfileAttributes(&o, name, monitor.ActionWaitRecover, 10, 4),
				),
			},
		},
	})
}

func testAccCheckPanosMonitorProfileExists(n string, o *monitor.Entry) resource.TestCheckFunc {
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
		v, err := fw.Network.MonitorProfile.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosMonitorProfileAttributes(o *monitor.Entry, name, action string, inter, thresh int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Action != action {
			return fmt.Errorf("Action is %q, expected %q", o.Action, action)
		}

		if o.Interval != inter {
			return fmt.Errorf("Interval is %d, expected %d", o.Interval, inter)
		}

		if o.Threshold != thresh {
			return fmt.Errorf("Threshold is %d, expected %d", o.Threshold, thresh)
		}

		return nil
	}
}

func testAccPanosMonitorProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_monitor_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			_, err := fw.Network.MonitorProfile.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccMonitorProfileConfig(name, action string, inter, thresh int) string {
	return fmt.Sprintf(`
resource "panos_monitor_profile" "test" {
    name = %q
    action = %q
    interval = %d
    threshold = %d
}
`, name, action, inter, thresh)
}
