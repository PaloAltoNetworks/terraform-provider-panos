package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/dev/profile/snmp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Resource tests.
func TestAccPanosSnmptrapServerProfile_basic(t *testing.T) {
	var tmpl, rsName string

	if testAccIsFirewall {
		rsName = "panos_snmptrap_server_profile"
	} else {
		tmpl = fmt.Sprintf("tf%s", acctest.RandString(6))
		rsName = "panos_panorama_snmptrap_server_profile"
	}

	var o snmp.Entry
	one := snmp.Entry{
		Name: fmt.Sprintf("tf%s", acctest.RandString(6)),
		V2cServers: []snmp.V2cServer{
			{
				Name:      "my-v2c-server",
				Manager:   "snmp2.example.com",
				Community: "public",
			},
		},
	}
	two := snmp.Entry{
		Name: one.Name,
		V3Servers: []snmp.V3Server{
			{
				Name:         "some-v3-server",
				Manager:      "snmp3.foobar.com",
				User:         "makoto",
				EngineId:     "0x0123456789",
				AuthPassword: "password",
				PrivPassword: "drowssap",
			},
		},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosSnmptrapServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSnmptrapServerProfileConfig(tmpl, one),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSnmptrapServerProfileExists(rsName, &o),
					testAccCheckPanosSnmptrapServerProfileAttributes(&o, &one),
				),
			},
			{
				Config: testAccSnmptrapServerProfileConfig(tmpl, two),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSnmptrapServerProfileExists(rsName, &o),
					testAccCheckPanosSnmptrapServerProfileAttributes(&o, &two),
				),
			},
		},
	})
}

func testAccCheckPanosSnmptrapServerProfileExists(n string, o *snmp.Entry) resource.TestCheckFunc {
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
		var v snmp.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, name := parseSnmptrapServerProfileId(rs.Primary.ID)
			v, err = con.Device.SnmpServerProfile.Get(vsys, name)
		case *pango.Panorama:
			tmpl, ts, vsys, name := parsePanoramaSnmptrapServerProfileId(rs.Primary.ID)
			v, err = con.Device.SnmpServerProfile.Get(tmpl, ts, vsys, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosSnmptrapServerProfileAttributes(o, conf *snmp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != conf.Name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, conf.Name)
		}

		if len(o.V2cServers) != len(conf.V2cServers) {
			return fmt.Errorf("v2 Server list len mismatch: %d vs %d", len(o.V2cServers), len(conf.V2cServers))
		}

		for i := range o.V2cServers {
			if o.V2cServers[i].Name != conf.V2cServers[i].Name {
				return fmt.Errorf("Server name is %s, not %s", o.V2cServers[i].Name, conf.V2cServers[i].Name)
			}

			if o.V2cServers[i].Manager != conf.V2cServers[i].Manager {
				return fmt.Errorf("Server Manager is %q, not %q", o.V2cServers[i].Manager, conf.V2cServers[i].Manager)
			}

			if o.V2cServers[i].Community != conf.V2cServers[i].Community {
				return fmt.Errorf("Server community is %s, not %s", o.V2cServers[i].Community, conf.V2cServers[i].Community)
			}
		}

		if len(o.V3Servers) != len(conf.V3Servers) {
			return fmt.Errorf("v3 Server list len mismatch: %d vs %d", len(o.V3Servers), len(conf.V3Servers))
		}

		for i := range o.V3Servers {
			if o.V3Servers[i].Name != conf.V3Servers[i].Name {
				return fmt.Errorf("Server name is %s, not %s", o.V3Servers[i].Name, conf.V3Servers[i].Name)
			}

			if o.V3Servers[i].Manager != conf.V3Servers[i].Manager {
				return fmt.Errorf("Server Manager is %q, not %q", o.V3Servers[i].Manager, conf.V3Servers[i].Manager)
			}

			if o.V3Servers[i].EngineId != conf.V3Servers[i].EngineId {
				return fmt.Errorf("Server engine id is %s, not %s", o.V3Servers[i].EngineId, conf.V3Servers[i].EngineId)
			}

			if o.V3Servers[i].User != conf.V3Servers[i].User {
				return fmt.Errorf("Server user is %s, not %s", o.V3Servers[i].User, conf.V3Servers[i].User)
			}
		}

		return nil
	}
}

func testAccPanosSnmptrapServerProfileDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_snmptrap_server_profile" && rs.Type != "panos_panorama_snmptrap_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, name := parseSnmptrapServerProfileId(rs.Primary.ID)
				_, err = con.Device.SnmpServerProfile.Get(vsys, name)
			case *pango.Panorama:
				tmpl, ts, vsys, name := parsePanoramaSnmptrapServerProfileId(rs.Primary.ID)
				_, err = con.Device.SnmpServerProfile.Get(tmpl, ts, vsys, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccSnmptrapServerProfileConfig(tmpl string, o snmp.Entry) string {
	var v2, v3 strings.Builder
	var rsName, tmplSpec, tmplAttr string

	if tmpl == "" {
		rsName = "panos_snmptrap_server_profile"
	} else {
		rsName = "panos_panorama_snmptrap_server_profile"
		tmplSpec = fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}
`, tmpl)
		tmplAttr = `
    template = panos_panorama_template.x.name`
	}

	v2.Grow(100 * len(o.V2cServers))
	for _, x := range o.V2cServers {
		v2.WriteString(fmt.Sprintf(`
    v2c_server {
        name = %q
        manager = %q
        community = %q
    }
`, x.Name, x.Manager, x.Community))
	}

	v3.Grow(200 * len(o.V3Servers))
	for _, x := range o.V3Servers {
		v3.WriteString(fmt.Sprintf(`
    v3_server {
        name = %q
        manager = %q
        engine_id = %q
        user = %q
        auth_password = %q
        priv_password = %q
    }
`, x.Name, x.Manager, x.EngineId, x.User, x.AuthPassword, x.PrivPassword))
	}

	return fmt.Sprintf(`
%s

resource %q "test" {
%s
    vsys = "shared"
    name = %q
%s
%s
}
`, tmplSpec, rsName, tmplAttr, o.Name, v2.String(), v3.String())
}
