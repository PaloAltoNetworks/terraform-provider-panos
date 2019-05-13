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

func TestAccPanosPanoramaBgpAuthProfile_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o auth.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaBgpAuthProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaBgpAuthProfileConfig(tmpl, vr, name, "sec1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpAuthProfileExists("panos_panorama_bgp_auth_profile.test", &o),
					testAccCheckPanosPanoramaBgpAuthProfileAttributes(&o, name),
				),
			},
			{
				Config: testAccPanoramaBgpAuthProfileConfig(tmpl, vr, name, "sec2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaBgpAuthProfileExists("panos_panorama_bgp_auth_profile.test", &o),
					testAccCheckPanosPanoramaBgpAuthProfileAttributes(&o, name),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaBgpAuthProfileExists(n string, o *auth.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, name := parsePanoramaBgpAuthProfileId(rs.Primary.ID)
		v, err := pano.Network.BgpAuthProfile.Get(tmpl, ts, vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaBgpAuthProfileAttributes(o *auth.Entry, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		return nil
	}
}

func testAccPanosPanoramaBgpAuthProfileDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_bgp_auth_profile" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, name := parsePanoramaBgpAuthProfileId(rs.Primary.ID)
			_, err := pano.Network.BgpAuthProfile.Get(tmpl, ts, vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaBgpAuthProfileConfig(tmpl, vr, name, sec string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "t" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "e1" {
    template = panos_panorama_template.t.name
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_panorama_virtual_router" "vr" {
    template = panos_panorama_template.t.name
    name = %q
    interfaces = [panos_panorama_ethernet_interface.e1.name]
}

resource "panos_panorama_bgp" "conf" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "42"
    enable = false
}

resource "panos_panorama_bgp_auth_profile" "test" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.conf.virtual_router
    name = %q
    secret = %q
}
`, tmpl, vr, name, sec)
}
