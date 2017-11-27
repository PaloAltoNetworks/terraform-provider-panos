package panos

import (
    "fmt"
    "testing"

    "github.com/PaloAltoNetworks/pango"
    "github.com/PaloAltoNetworks/pango/netw/mngtprof"

    "github.com/hashicorp/terraform/helper/acctest"
    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
)


func TestPanosManagementProfile_basic(t *testing.T) {
    var mp mngtprof.Entry
    name := fmt.Sprintf("tf%s", acctest.RandString(6))

    resource.Test(t, resource.TestCase{
        PreCheck: func() { testAccPreCheck(t) },
        Providers: testAccProviders,
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
            return fmt.Errorf("Management profile label ID is not set")
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

        if len(o.PermittedIp) != 2 {
            return fmt.Errorf("len(PermittedIp) is %d, expected 2", len(o.PermittedIp))
        }

        if o.PermittedIp[0] != pi1 {
            return fmt.Errorf("Permitted IP 0 is %s, expected %s", o.PermittedIp[0], pi1)
        }

        if o.PermittedIp[1] != pi2 {
            return fmt.Errorf("Permitted IP 1 is %s, expected %s", o.PermittedIp[1], pi2)
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

        name := rs.Primary.ID
        if name != "" {
            info, err := fw.Network.ManagementProfile.Get(name)
            if err == nil {
                return fmt.Errorf("Management profile %q still exists: %v %#v", name, err, info)
            }
        }
        //return nil
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
    permitted_ip = ["%s", "%s"]
}
`, n, h, p, s, pi1, pi2)
}
