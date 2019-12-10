package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog/server"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosSyslogServerProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var (
		o          syslog.Entry
		serverList []server.Entry
	)

	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	sn := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosSyslogServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogServerProfileConfig(name, "$serial $severity", "$from $rule", sn, "server1.example.com", server.TransportUdp, server.SyslogFormatIetf, server.FacilityLocal3, 515),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSyslogServerProfileExists("panos_syslog_server_profile.test", &o, &serverList),
					testAccCheckPanosSyslogServerProfileAttributes(&o, &serverList, name, "$serial $severity", "$from $rule", sn, "server1.example.com", server.TransportUdp, server.SyslogFormatIetf, server.FacilityLocal3, 515),
				),
			},
			{
				Config: testAccSyslogServerProfileConfig(name, "$from $severity", "$serial $rule", sn, "server2.example.com", server.TransportTcp, server.SyslogFormatBsd, server.FacilityUser, 514),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSyslogServerProfileExists("panos_syslog_server_profile.test", &o, &serverList),
					testAccCheckPanosSyslogServerProfileAttributes(&o, &serverList, name, "$from $severity", "$serial $rule", sn, "server2.example.com", server.TransportTcp, server.SyslogFormatBsd, server.FacilityUser, 514),
				),
			},
		},
	})
}

func testAccCheckPanosSyslogServerProfileExists(n string, o *syslog.Entry, serverList *[]server.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vsys, name := parseSyslogServerProfileId(rs.Primary.ID)
		v, err := fw.Device.SyslogServerProfile.Get(vsys, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		list, err := fw.Device.SyslogServer.GetList(vsys, name)
		if err != nil {
			return err
		}
		entries := make([]server.Entry, 0, len(list))
		for i := range list {
			entry, err := fw.Device.SyslogServer.Get(vsys, name, list[i])
			if err != nil {
				return err
			}
			entries = append(entries, entry)
		}

		*serverList = entries

		return nil
	}
}

func testAccCheckPanosSyslogServerProfileAttributes(o *syslog.Entry, serverList *[]server.Entry, name, system, url, sn, ss, trans, sf, fac string, port int) resource.TestCheckFunc {
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

func testAccPanosSyslogServerProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_syslog_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, name := parseSyslogServerProfileId(rs.Primary.ID)
			_, err := fw.Device.SyslogServerProfile.Get(vsys, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccSyslogServerProfileConfig(name, system, url, sn, ss, trans, sf, fac string, port int) string {
	return fmt.Sprintf(`
resource "panos_syslog_server_profile" "test" {
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
`, name, system, url, sn, ss, trans, sf, fac, port)
}
