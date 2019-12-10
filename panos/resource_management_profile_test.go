package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/profile/mngtprof"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosManagementProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var mp mngtprof.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosManagementProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccManagementProfileConfig(name, true, false, true, "10.1.1.1", "192.168.1.1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosManagementProfileExists("panos_management_profile.test", &mp),
					testAccCheckPanosManagementProfileAttributes(&mp, name, true, false, true, "10.1.1.1", "192.168.1.1"),
				),
			},
			{
				Config: testAccManagementProfileConfig(name, false, true, false, "10.1.1.2", "192.168.1.2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosManagementProfileExists("panos_management_profile.test", &mp),
					testAccCheckPanosManagementProfileAttributes(&mp, name, false, true, false, "10.1.1.2", "192.168.1.2"),
				),
			},
		},
	})
}

func testAccCheckPanosManagementProfileExists(n string, o *mngtprof.Entry) resource.TestCheckFunc {
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
		v, err := fw.Network.ManagementProfile.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosManagementProfileAttributes(o *mngtprof.Entry, n string, h, p, ssh bool, pi1, pi2 string) resource.TestCheckFunc {
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

func testAccPanosManagementProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_management_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			_, err := fw.Network.ManagementProfile.Get(rs.Primary.ID)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccManagementProfileConfig(n string, h, p, s bool, pi1, pi2 string) string {
	return fmt.Sprintf(`
resource "panos_management_profile" "test" {
    name = "%s"
    https = %t
    ping = %t
    ssh = %t
    permitted_ips = ["%s", "%s"]
}
`, n, h, p, s, pi1, pi2)
}
