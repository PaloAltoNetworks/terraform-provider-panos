package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog/server"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaSyslogServerProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var (
		o          syslog.Entry
		serverList []server.Entry
	)

	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	sn := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaSyslogServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaSyslogServerProfileConfig(tmpl, name, "$serial $severity", "$from $rule", sn, "server1.example.com", server.TransportUdp, server.SyslogFormatIetf, server.FacilityLocal3, 515),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaSyslogServerProfileExists("panos_panorama_syslog_server_profile.test", &o, &serverList),
					testAccCheckPanosPanoramaSyslogServerProfileAttributes(&o, &serverList, name, "$serial $severity", "$from $rule", sn, "server1.example.com", server.TransportUdp, server.SyslogFormatIetf, server.FacilityLocal3, 515),
				),
			},
			{
				Config: testAccPanoramaSyslogServerProfileConfig(tmpl, name, "$from $severity", "$serial $rule", sn, "server2.example.com", server.TransportTcp, server.SyslogFormatBsd, server.FacilityUser, 514),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaSyslogServerProfileExists("panos_panorama_syslog_server_profile.test", &o, &serverList),
					testAccCheckPanosPanoramaSyslogServerProfileAttributes(&o, &serverList, name, "$from $severity", "$serial $rule", sn, "server2.example.com", server.TransportTcp, server.SyslogFormatBsd, server.FacilityUser, 514),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaSyslogServerProfileExists(n string, o *syslog.Entry, serverList *[]server.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vsys, dg, name := parsePanoramaSyslogServerProfileId(rs.Primary.ID)
		v, err := pano.Device.SyslogServerProfile.Get(tmpl, ts, vsys, dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		list, err := pano.Device.SyslogServer.GetList(tmpl, ts, vsys, dg, name)
		if err != nil {
			return err
		}
		entries := make([]server.Entry, 0, len(list))
		for i := range list {
			entry, err := pano.Device.SyslogServer.Get(tmpl, ts, vsys, dg, name, list[i])
			if err != nil {
				return err
			}
			entries = append(entries, entry)
		}

		*serverList = entries

		return nil
	}
}

func testAccCheckPanosPanoramaSyslogServerProfileAttributes(o *syslog.Entry, serverList *[]server.Entry, name, system, url, sn, ss, trans, sf, fac string, port int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.System != system {
			return fmt.Errorf("System format is %q, not %q", o.System, system)
		}

		if o.Url != url {
			return fmt.Errorf("URL format is %q, not %q", o.Url, url)
		}

		if len(*serverList) != 1 {
			return fmt.Errorf("Server list is len %d, not 1", len(*serverList))
		}

		e := (*serverList)[0]

		if e.Name != sn {
			return fmt.Errorf("Server name is %s, not %s", e.Name, sn)
		}

		if e.Server != ss {
			return fmt.Errorf("Server syslog server is %q, not %q", e.Server, ss)
		}

		if e.Transport != trans {
			return fmt.Errorf("Server transport is %s, not %s", e.Transport, trans)
		}

		if e.SyslogFormat != sf {
			return fmt.Errorf("Server syslog format is %s, not %s", e.SyslogFormat, sf)
		}

		if e.Facility != fac {
			return fmt.Errorf("Server facility is %s, not %s", e.Facility, fac)
		}

		if e.Port != port {
			return fmt.Errorf("Server port is %d, not %d", e.Port, port)
		}

		return nil
	}
}

func testAccPanosPanoramaSyslogServerProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_syslog_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vsys, dg, name := parsePanoramaSyslogServerProfileId(rs.Primary.ID)
			_, err := pano.Device.SyslogServerProfile.Get(tmpl, ts, vsys, dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaSyslogServerProfileConfig(tmpl, name, system, url, sn, ss, trans, sf, fac string, port int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "for syslog server profile test"
}

resource "panos_panorama_syslog_server_profile" "test" {
    template = panos_panorama_template.x.name
    vsys = "shared"
    name = %q
    system_format = %q
    url_format = %q
    syslog_server {
        name = %q
        server = %q
        transport = %q
        syslog_format = %q
        facility = %q
        port = %d
    }
}
`, tmpl, name, system, url, sn, ss, trans, sf, fac, port)
}
