package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/profile/mngtprof"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaManagementProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var mp mngtprof.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaManagementProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaManagementProfileConfig(tmpl, name, true, false, true, "10.1.1.1", "192.168.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaManagementProfileExists("panos_panorama_management_profile.test", &mp),
					testAccCheckPanosPanoramaManagementProfileAttributes(&mp, name, true, false, true, "10.1.1.1", "192.168.1.1"),
				),
			},
			{
				Config: testAccPanoramaManagementProfileConfig(tmpl, name, false, true, false, "10.1.1.2", "192.168.1.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaManagementProfileExists("panos_panorama_management_profile.test", &mp),
					testAccCheckPanosPanoramaManagementProfileAttributes(&mp, name, false, true, false, "10.1.1.2", "192.168.1.2"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaManagementProfileExists(n string, o *mngtprof.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, name := parsePanoramaManagementProfileId(rs.Primary.ID)
		v, err := pano.Network.ManagementProfile.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaManagementProfileAttributes(o *mngtprof.Entry, n string, h, p, ssh bool, pi1, pi2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %s, expected %s", o.Name, n)
		}

		if o.Https != h {
			return fmt.Errorf("HTTPS is %t, expected %t", o.Https, h)
		}

		if o.Ping != p {
			return fmt.Errorf("ping is %t, expected %t", o.Ping, p)
		}

		if o.Ssh != ssh {
			return fmt.Errorf("SSH is %t, expected %t", o.Ssh, ssh)
		}

		if len(o.PermittedIps) != 2 {
			return fmt.Errorf("len(PermittedIps) is %d, expected 2", len(o.PermittedIps))
		}

		if o.PermittedIps[0] != pi1 {
			return fmt.Errorf("Permitted IP 0 is %s, expected %s", o.PermittedIps[0], pi1)
		}

		if o.PermittedIps[1] != pi2 {
			return fmt.Errorf("Permitted IP 1 is %s, expected %s", o.PermittedIps[1], pi2)
		}

		return nil
	}
}

func testAccPanosPanoramaManagementProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_management_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, name := parsePanoramaManagementProfileId(rs.Primary.ID)
			_, err := pano.Network.ManagementProfile.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaManagementProfileConfig(tmpl, n string, h, p, s bool, pi1, pi2 string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_management_profile" "test" {
    template = panos_panorama_template.x.name
    name = %q
    https = %t
    ping = %t
    ssh = %t
    permitted_ips = [%q, %q]
}
`, tmpl, n, h, p, s, pi1, pi2)
}
