package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v2c"
	"github.com/PaloAltoNetworks/pango/dev/profile/snmp/v3"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaSnmptrapServerProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var (
		o       snmp.Entry
		v2cList []v2c.Entry
		v3List  []v3.Entry
	)

	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vsys := "shared"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaSnmptrapServerProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaSnmptrapServerProfileConfig(tmpl, vsys, name, snmp.SnmpVersionV2c),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaSnmptrapServerProfileExists("panos_panorama_snmptrap_server_profile.test", &o, &v2cList, &v3List),
					testAccCheckPanosPanoramaSnmptrapServerProfileAttributes(&o, name, snmp.SnmpVersionV2c, &v2cList, &v3List),
				),
			},
			{
				Config: testAccPanoramaSnmptrapServerProfileConfig(tmpl, vsys, name, snmp.SnmpVersionV3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaSnmptrapServerProfileExists("panos_panorama_snmptrap_server_profile.test", &o, &v2cList, &v3List),
					testAccCheckPanosPanoramaSnmptrapServerProfileAttributes(&o, name, snmp.SnmpVersionV3, &v2cList, &v3List),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaSnmptrapServerProfileExists(n string, o *snmp.Entry, v2cList *[]v2c.Entry, v3List *[]v3.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vsys, dg, name := parsePanoramaSnmptrapServerProfileId(rs.Primary.ID)
		v, err := pano.Device.SnmpServerProfile.Get(tmpl, ts, vsys, dg, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		switch v.SnmpVersion {
		case snmp.SnmpVersionV2c:
			list, err := pano.Device.SnmpV2cServer.GetList(tmpl, ts, vsys, dg, name)
			if err != nil {
				return err
			}
			entries := make([]v2c.Entry, 0, len(list))
			for i := range list {
				entry, err := pano.Device.SnmpV2cServer.Get(tmpl, ts, vsys, dg, name, list[i])
				if err != nil {
					return err
				}
				entries = append(entries, entry)
			}
			*v2cList = entries
			*v3List = nil
		case snmp.SnmpVersionV3:
			list, err := pano.Device.SnmpV3Server.GetList(tmpl, ts, vsys, dg, name)
			if err != nil {
				return err
			}
			entries := make([]v3.Entry, 0, len(list))
			for i := range list {
				entry, err := pano.Device.SnmpV3Server.Get(tmpl, ts, vsys, dg, name, list[i])
				if err != nil {
					return err
				}
				entries = append(entries, entry)
			}
			*v2cList = nil
			*v3List = entries
		}

		return nil
	}
}

func testAccCheckPanosPanoramaSnmptrapServerProfileAttributes(o *snmp.Entry, name, snmp_version string, v2cList *[]v2c.Entry, v3List *[]v3.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.SnmpVersion != snmp_version {
			return fmt.Errorf("Snmp version is %s, not %s", o.SnmpVersion, snmp_version)
		}

		switch snmp_version {
		case snmp.SnmpVersionV2c:
			if len(*v2cList) != 1 {
				return fmt.Errorf("SNMP v2c servers are len %d, not 1", len(*v2cList))
			}

			entry := (*v2cList)[0]

			if entry.Name != "my-v2c-server" {
				return fmt.Errorf("v2c name is %s, not %s", entry.Name, "my-v2c-server")
			}

			if entry.Manager != "snmp2.example.com" {
				return fmt.Errorf("v2c manager is %s, not 'snmp2.example.com'", entry.Manager)
			}

			if entry.Community != "public" {
				return fmt.Errorf("v2c community is %s, not 'public'", entry.Community)
			}
		case snmp.SnmpVersionV3:
			if len(*v3List) != 1 {
				return fmt.Errorf("SNMP v3 servers are len %d, not 1", len(*v3List))
			}

			entry := (*v3List)[0]

			if entry.Name != "some-v3-server" {
				return fmt.Errorf("v3 name is %s, not %s", entry.Name, "some-v3-server")
			}

			if entry.Manager != "snmp3.foobar.com" {
				return fmt.Errorf("v3 manager is %s, not 'snmp3.foobar.com'", entry.Manager)
			}

			if entry.EngineId != "0x0123456789" {
				return fmt.Errorf("v3 engine id is %s, not '0x0123456789'", entry.EngineId)
			}

			if entry.User != "makoto" {
				return fmt.Errorf("v3 user is %s, not 'makoto'", entry.User)
			}
		}

		return nil
	}
}

func testAccPanosPanoramaSnmptrapServerProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_snmptrap_server_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vsys, dg, name := parsePanoramaSnmptrapServerProfileId(rs.Primary.ID)
			_, err := pano.Device.SnmpServerProfile.Get(tmpl, ts, vsys, dg, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaSnmptrapServerProfileConfig(tmpl, vsys, name, snmp_version string) string {
	config := map[string]string{
		snmp.SnmpVersionV2c: `
    v2c_server {
        name = "my-v2c-server"
        manager = "snmp2.example.com"
        community = "public"
    }`,
		snmp.SnmpVersionV3: `
    v3_server {
        name = "some-v3-server"
        manager = "snmp3.foobar.com"
        engine_id = "0x0123456789"
        user = "makoto"
        auth_password = "password"
        priv_password = "drowssap"
    }`,
	}

	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
    description = "for snmptrap server profile"
}

resource "panos_panorama_snmptrap_server_profile" "test" {
    template = panos_panorama_template.x.name
    vsys = %q
    name = %q
%s
}
`, tmpl, vsys, name, config[snmp_version])
}
