package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/syslog"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source listing test.
func TestAccPanosDsSyslogServerProfileList(t *testing.T) {
	var tmpl string

	if testAccIsPanorama {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
	}

	o := syslog.Entry{
		Name:   fmt.Sprintf("tf%s", acctest.RandString(6)),
		System: "$serial $severity",
		Url:    "$from $rule",
		Servers: []syslog.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				Server:       "server1.example.com",
				Transport:    syslog.TransportUdp,
				SyslogFormat: syslog.SyslogFormatIetf,
				Facility:     syslog.FacilityLocal3,
				Port:         515,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogServerProfileConfig(tmpl, o),
				Check:  checkDataSourceListing("panos_syslog_server_profiles"),
			},
		},
	})
}

// Data source test.
func TestAccPanosDsSyslogServerProfile(t *testing.T) {
	var tmpl string

	if testAccIsPanorama {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
	}

	o := syslog.Entry{
		Name:   fmt.Sprintf("tf%s", acctest.RandString(6)),
		System: "$serial $severity",
		Url:    "$from $rule",
		Servers: []syslog.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				Server:       "server1.example.com",
				Transport:    syslog.TransportUdp,
				SyslogFormat: syslog.SyslogFormatIetf,
				Facility:     syslog.FacilityLocal3,
				Port:         515,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogServerProfileConfig(tmpl, o),
				Check: checkDataSource(
					"panos_syslog_server_profile",
					[]string{"name", "system_format", "url_format"},
				),
			},
		},
	})
}

// Resource tests.
func TestAccPanosSyslogServerProfile_basic(t *testing.T) {
	var tmpl, rsName string

	if testAccIsFirewall {
		rsName = "panos_syslog_server_profile"
	} else {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
		rsName = "panos_panorama_syslog_server_profile"
	}

	var o syslog.Entry
	one := syslog.Entry{
		Name:   fmt.Sprintf("tf%s", acctest.RandString(6)),
		System: "$serial $severity",
		Url:    "$from $rule",
		Servers: []syslog.Server{
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				Server:       "server1.example.com",
				Transport:    syslog.TransportUdp,
				SyslogFormat: syslog.SyslogFormatIetf,
				Facility:     syslog.FacilityLocal3,
				Port:         515,
			},
		},
	}
	two := syslog.Entry{
		Name:   one.Name,
		System: "$from $severity",
		Url:    "$serial $rule",
		Servers: []syslog.Server{
			{
				Name:         one.Servers[0].Name,
				Server:       "server1.example.com",
				Transport:    syslog.TransportUdp,
				SyslogFormat: syslog.SyslogFormatIetf,
				Facility:     syslog.FacilityLocal3,
				Port:         515,
			},
			{
				Name:         fmt.Sprintf("tf%s", acctest.RandString(6)),
				Server:       "server2.example.com",
				Transport:    syslog.TransportTcp,
				SyslogFormat: syslog.SyslogFormatBsd,
				Facility:     syslog.FacilityLocal4,
				Port:         514,
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosSyslogServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSyslogServerProfileConfig(tmpl, one),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSyslogServerProfileExists(rsName, &o),
					testAccCheckPanosSyslogServerProfileAttributes(&o, &one),
				),
			},
			{
				Config: testAccSyslogServerProfileConfig(tmpl, two),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSyslogServerProfileExists(rsName, &o),
					testAccCheckPanosSyslogServerProfileAttributes(&o, &two),
				),
			},
		},
	})
}

func testAccCheckPanosSyslogServerProfileExists(n string, o *syslog.Entry) resource.TestCheckFunc {
	name := fmt.Sprintf("%s.test", n)
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Resource not found: %s / %#v", n, s.RootModule().Resources)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v syslog.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseSyslogServerProfileId(rs.Primary.ID)
			v, err = con.Device.SyslogServerProfile.Get(vsys, name)
		case *pango.Panorama:
			tmpl, ts, vsys, name := parsePanoramaSyslogServerProfileId(rs.Primary.ID)
			v, err = con.Device.SyslogServerProfile.Get(tmpl, ts, vsys, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosSyslogServerProfileAttributes(o, conf *syslog.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != conf.Name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, conf.Name)
		}

		if o.System != conf.System {
			return fmt.Errorf("System format is %q, not %q", o.System, conf.System)
		}

		if o.Url != conf.Url {
			return fmt.Errorf("URL format is %q, not %q", o.Url, conf.Url)
		}

		if len(o.Servers) != len(conf.Servers) {
			return fmt.Errorf("Server list len mismatch: %d vs %d", len(o.Servers), len(conf.Servers))
		}

		for i := range o.Servers {
			if o.Servers[i].Name != conf.Servers[i].Name {
				return fmt.Errorf("Server name is %s, not %s", o.Servers[i].Name, conf.Servers[i].Name)
			}

			if o.Servers[i].Server != conf.Servers[i].Server {
				return fmt.Errorf("Server syslog server is %q, not %q", o.Servers[i].Server, conf.Servers[i].Server)
			}

			if o.Servers[i].Transport != conf.Servers[i].Transport {
				return fmt.Errorf("Server transport is %s, not %s", o.Servers[i].Transport, conf.Servers[i].Transport)
			}

			if o.Servers[i].SyslogFormat != conf.Servers[i].SyslogFormat {
				return fmt.Errorf("Server syslog format is %s, not %s", o.Servers[i].SyslogFormat, conf.Servers[i].SyslogFormat)
			}

			if o.Servers[i].Facility != conf.Servers[i].Facility {
				return fmt.Errorf("Server facility is %s, not %s", o.Servers[i].Facility, conf.Servers[i].Facility)
			}

			if o.Servers[i].Port != conf.Servers[i].Port {
				return fmt.Errorf("Server port is %d, not %d", o.Servers[i].Port, conf.Servers[i].Port)
			}
		}

		return nil
	}
}

func testAccPanosSyslogServerProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_syslog_server_profile" && rs.Type != "panos_panorama_syslog_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseSyslogServerProfileId(rs.Primary.ID)
				_, err = con.Device.SyslogServerProfile.Get(vsys, name)
			case *pango.Panorama:
				tmpl, ts, vsys, name := parsePanoramaSyslogServerProfileId(rs.Primary.ID)
				_, err = con.Device.SyslogServerProfile.Get(tmpl, ts, vsys, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccSyslogServerProfileConfig(tmpl string, o syslog.Entry) string {
	var b strings.Builder
	var rsName, tmplSpec, tmplAttr string

	if tmpl == "" {
		rsName = "panos_syslog_server_profile"
	} else {
		rsName = "panos_panorama_syslog_server_profile"
		tmplSpec = fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}
`, tmpl)
		tmplAttr = `
    template = panos_panorama_template.x.name`
	}

	b.Grow(200 * len(o.Servers))
	for _, x := range o.Servers {
		b.WriteString(fmt.Sprintf(`
    syslog_server {
        name = %q
        server = %q
        transport = %q
        syslog_format = %q
        facility = %q
        port = %d
    }
`, x.Name, x.Server, x.Transport, x.SyslogFormat, x.Facility, x.Port))
	}

	return fmt.Sprintf(`
%s

resource %q "test" {
%s
    vsys = "shared"
    name = %q
    system_format = %q
    url_format = %q
%s
}

data "panos_syslog_server_profile" "test" {
%s
    vsys = "shared"
    name = %s.test.name
}

data "panos_syslog_server_profiles" "test" {
%s
    vsys = "shared"
}
`, tmplSpec, rsName, tmplAttr, o.Name, o.System, o.Url, b.String(), tmplAttr, rsName, tmplAttr)
}
