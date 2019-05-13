package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/profile/auth"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosBgpAuthProfile_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o auth.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpAuthProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpAuthProfileConfig(vr, name, "sec1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpAuthProfileExists("panos_bgp_auth_profile.test", &o),
					testAccCheckPanosBgpAuthProfileAttributes(&o, name),
				),
			},
			{
				Config: testAccBgpAuthProfileConfig(vr, name, "sec2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpAuthProfileExists("panos_bgp_auth_profile.test", &o),
					testAccCheckPanosBgpAuthProfileAttributes(&o, name),
				),
			},
		},
	})
}

func testAccCheckPanosBgpAuthProfileExists(n string, o *auth.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, name := parseBgpAuthProfileId(rs.Primary.ID)
		v, err := fw.Network.BgpAuthProfile.Get(vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpAuthProfileAttributes(o *auth.Entry, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		return nil
	}
}

func testAccPanosBgpAuthProfileDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_auth_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, name := parseBgpAuthProfileId(rs.Primary.ID)
			_, err := fw.Network.BgpAuthProfile.Get(vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpAuthProfileConfig(vr, name, sec string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "e1" {
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_virtual_router" "vr" {
    name = %q
    interfaces = [panos_ethernet_interface.e1.name]
}

resource "panos_bgp" "conf" {
    virtual_router = panos_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_bgp_auth_profile" "test" {
    virtual_router = panos_bgp.conf.virtual_router
    name = %q
    secret = %q
}
`, vr, name, sec)
}
