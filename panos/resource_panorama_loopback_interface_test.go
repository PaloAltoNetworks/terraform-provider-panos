package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/interface/loopback"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaLoopbackInterface_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o loopback.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("loopback.%d", num)
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaLoopbackInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaLoopbackInterfaceConfig(tmpl, name, "first comment", "10.8.9.1", 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLoopbackInterfaceExists("panos_panorama_loopback_interface.test", &o),
					testAccCheckPanosPanoramaLoopbackInterfaceAttributes(&o, name, "first comment", "10.8.9.1", 600),
				),
			},
			{
				Config: testAccPanoramaLoopbackInterfaceConfig(tmpl, name, "second comment", "10.9.10.1", 700),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaLoopbackInterfaceExists("panos_panorama_loopback_interface.test", &o),
					testAccCheckPanosPanoramaLoopbackInterfaceAttributes(&o, name, "second comment", "10.9.10.1", 700),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaLoopbackInterfaceExists(n string, o *loopback.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, _, name := parsePanoramaLoopbackInterfaceId(rs.Primary.ID)
		v, err := pano.Network.LoopbackInterface.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaLoopbackInterfaceAttributes(o *loopback.Entry, name, cmt, ip string, mtu int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Comment != cmt {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, cmt)
		}

		if len(o.StaticIps) != 1 || o.StaticIps[0] != ip {
			return fmt.Errorf("Static IPs is %#v, expected [%q]", o.StaticIps, ip)
		}

		if o.Mtu != mtu {
			return fmt.Errorf("MTU is %d, expected %d", o.Mtu, mtu)
		}

		return nil
	}
}

func testAccPanosPanoramaLoopbackInterfaceDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_loopback_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, _, name := parsePanoramaLoopbackInterfaceId(rs.Primary.ID)
			_, err := pano.Network.LoopbackInterface.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaLoopbackInterfaceConfig(tmpl, name, cmt, ip string, mtu int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_loopback_interface" "test" {
    template = panos_panorama_template.x.name
    name = %q
    comment = %q
    static_ips = [%q]
    mtu = %d
}
`, tmpl, name, cmt, ip, mtu)
}
