package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/interface/tunnel"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaTunnelInterface_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o tunnel.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("tunnel.%d", num)
	tmpl := fmt.Sprintf("tfTmpl%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaTunnelInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaTunnelInterfaceConfig(tmpl, name, "first comment", "10.8.9.1/24", 600),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTunnelInterfaceExists("panos_panorama_tunnel_interface.test", &o),
					testAccCheckPanosPanoramaTunnelInterfaceAttributes(&o, name, "first comment", "10.8.9.1/24", 600),
				),
			},
			{
				Config: testAccPanoramaTunnelInterfaceConfig(tmpl, name, "second comment", "10.9.10.1/24", 700),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaTunnelInterfaceExists("panos_panorama_tunnel_interface.test", &o),
					testAccCheckPanosPanoramaTunnelInterfaceAttributes(&o, name, "second comment", "10.9.10.1/24", 700),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaTunnelInterfaceExists(n string, o *tunnel.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, _, name := parsePanoramaTunnelInterfaceId(rs.Primary.ID)
		v, err := pano.Network.TunnelInterface.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaTunnelInterfaceAttributes(o *tunnel.Entry, name, cmt, ip string, mtu int) resource.TestCheckFunc {
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

func testAccPanosPanoramaTunnelInterfaceDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_tunnel_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, _, name := parsePanoramaTunnelInterfaceId(rs.Primary.ID)
			_, err := pano.Network.TunnelInterface.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaTunnelInterfaceConfig(tmpl, name, cmt, ip string, mtu int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_tunnel_interface" "test" {
    name = %q
    template = panos_panorama_template.x.name
    comment = %q
    static_ips = [%q]
    mtu = %d
}
`, tmpl, name, cmt, ip, mtu)
}
