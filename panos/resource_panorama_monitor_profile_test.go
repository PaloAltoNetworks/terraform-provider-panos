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

func TestAccPanosPanoramaMonitorProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o monitor.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaMonitorProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaMonitorProfileConfig(tmpl, name, monitor.ActionFailOver, 11, 7),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaMonitorProfileExists("panos_panorama_monitor_profile.test", &o),
					testAccCheckPanosPanoramaMonitorProfileAttributes(&o, name, monitor.ActionFailOver, 11, 7),
				),
			},
			{
				Config: testAccPanoramaMonitorProfileConfig(tmpl, name, monitor.ActionWaitRecover, 10, 4),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaMonitorProfileExists("panos_panorama_monitor_profile.test", &o),
					testAccCheckPanosPanoramaMonitorProfileAttributes(&o, name, monitor.ActionWaitRecover, 10, 4),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaMonitorProfileExists(n string, o *monitor.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaMonitorProfileId(rs.Primary.ID)
		v, err := pano.Network.MonitorProfile.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaMonitorProfileAttributes(o *monitor.Entry, name, action string, inter, thresh int) resource.TestCheckFunc {
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

func testAccPanosPanoramaMonitorProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_monitor_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaMonitorProfileId(rs.Primary.ID)
			_, err := pano.Network.MonitorProfile.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaMonitorProfileConfig(tmpl, name, action string, inter, thresh int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_monitor_profile" "test" {
    template = panos_panorama_template.x.name
    name = %q
    action = %q
    interval = %d
    threshold = %d
}
`, tmpl, name, action, inter, thresh)
}
