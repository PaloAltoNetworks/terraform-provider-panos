package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/email"
	"github.com/PaloAltoNetworks/pango/dev/profile/email/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaEmailServerProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var (
		o          email.Entry
		serverList []server.Entry
	)

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	sn := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaEmailServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaEmailServerProfileConfig(tmpl, name, "$serial $severity", "$from $rule", sn, "foobar", "from@example.com", "to@example.com", "archive@example.com", "smtp.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaEmailServerProfileExists("panos_panorama_email_server_profile.test", &o, &serverList),
					testAccCheckPanosPanoramaEmailServerProfileAttributes(&o, &serverList, name, "$serial $severity", "$from $rule", sn, "foobar", "from@example.com", "to@example.com", "archive@example.com", "smtp.example.com"),
				),
			},
			{
				Config: testAccPanoramaEmailServerProfileConfig(tmpl, name, "$from $severity", "$serial $rule", sn, "barbaz", "john@example.com", "jacob@example.com", "jingleheimer@example.com", "smith.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaEmailServerProfileExists("panos_panorama_email_server_profile.test", &o, &serverList),
					testAccCheckPanosPanoramaEmailServerProfileAttributes(&o, &serverList, name, "$from $severity", "$serial $rule", sn, "barbaz", "john@example.com", "jacob@example.com", "jingleheimer@example.com", "smith.example.com"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaEmailServerProfileExists(n string, o *email.Entry, serverList *[]server.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vsys, dg, name := parsePanoramaEmailServerProfileId(rs.Primary.ID)
		v, err := pano.Device.EmailServerProfile.Get(tmpl, ts, vsys, dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		list, err := pano.Device.EmailServer.GetList(tmpl, ts, vsys, dg, name)
		if err != nil {
			return err
		}
		entries := make([]server.Entry, 0, len(list))
		for i := range list {
			entry, err := pano.Device.EmailServer.Get(tmpl, ts, vsys, dg, name, list[i])
			if err != nil {
				return err
			}
			entries = append(entries, entry)
		}

		*serverList = entries

		return nil
	}
}

func testAccCheckPanosPanoramaEmailServerProfileAttributes(o *email.Entry, serverList *[]server.Entry, name, hip, cf, sn, display, from, to, also, gw string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.HipMatch != hip {
			return fmt.Errorf("HIP match format is %q, not %q", o.HipMatch, hip)
		}

		if o.Config != cf {
			return fmt.Errorf("Config format is %q, not %q", o.Config, cf)
		}

		if len(*serverList) != 1 {
			return fmt.Errorf("Server list is len %d, not 1", len(*serverList))
		}

		e := (*serverList)[0]

		if e.Name != sn {
			return fmt.Errorf("Server name is %s, not %s", e.Name, sn)
		}

		if e.DisplayName != display {
			return fmt.Errorf("Display name is %s, not %s", e.DisplayName, display)
		}

		if e.From != from {
			return fmt.Errorf("From is %s, not %s", e.From, from)
		}

		if e.To != to {
			return fmt.Errorf("To is %s, not %s", e.To, to)
		}

		if e.AlsoTo != also {
			return fmt.Errorf("Also to is %s, not %s", e.AlsoTo, also)
		}

		if e.EmailGateway != gw {
			return fmt.Errorf("Email gateway is %s, not %s", e.EmailGateway, gw)
		}

		return nil
	}
}

func testAccPanosPanoramaEmailServerProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_email_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vsys, dg, name := parsePanoramaEmailServerProfileId(rs.Primary.ID)
			_, err := pano.Device.EmailServerProfile.Get(tmpl, ts, vsys, dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaEmailServerProfileConfig(tmpl, name, hip, cf, sn, display, from, to, also, gw string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "for email server profile test"
}

resource "panos_panorama_email_server_profile" "test" {
    template = panos_panorama_template.x.name
    vsys = "shared"
    name = %q
    hip_match_format = %q
    config_format = %q
    email_server {
        name = %q
        display_name = %q
        from_email = %q
        to_email = %q
        also_to_email = %q
        email_gateway= %q
    }
}
`, tmpl, name, hip, cf, sn, display, from, to, also, gw)
}
