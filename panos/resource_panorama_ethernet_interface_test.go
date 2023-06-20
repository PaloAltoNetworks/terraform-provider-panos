package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/interface/eth"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaEthernetInterface_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o eth.Entry
	num := (acctest.RandInt() % 9) + 1
	name := fmt.Sprintf("ethernet1/%d", num)
	tmpl := fmt.Sprintf("tfTmpl%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaEthernetInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaEthernetInterfaceConfig(tmpl, name, "down", "first comment", "10.1.1.1/24", "192.168.1.1/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaEthernetInterfaceExists("panos_panorama_ethernet_interface.test", &o),
					testAccCheckPanosPanoramaEthernetInterfaceAttributes(&o, name, "down", "first comment", "10.1.1.1/24", "192.168.1.1/24"),
				),
			},
			{
				Config: testAccPanoramaEthernetInterfaceConfig(tmpl, name, "up", "second comment", "10.1.2.1/24", "192.168.2.1/24"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaEthernetInterfaceExists("panos_panorama_ethernet_interface.test", &o),
					testAccCheckPanosPanoramaEthernetInterfaceAttributes(&o, name, "up", "second comment", "10.1.2.1/24", "192.168.2.1/24"),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaEthernetInterfaceExists(n string, o *eth.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, _, name := parsePanoramaEthernetInterfaceId(rs.Primary.ID)
		v, err := pano.Network.EthernetInterface.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaEthernetInterfaceAttributes(o *eth.Entry, n, ls, c, i1, i2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != n {
			return fmt.Errorf("Name is %s, expected %s", o.Name, n)
		}

		if o.LinkState != ls {
			return fmt.Errorf("Link state is %s, expected %s", o.LinkState, ls)
		}

		if o.Comment != c {
			return fmt.Errorf("Comment is %q, expected %q", o.Comment, c)
		}

		if len(o.StaticIps) != 2 {
			return fmt.Errorf("len(StaticIps) is %d, expected 2", len(o.StaticIps))
		}

		if o.StaticIps[0] != i1 {
			return fmt.Errorf("StaticIps[0] is %s, expected %s", o.StaticIps[0], i1)
		}

		if o.StaticIps[1] != i2 {
			return fmt.Errorf("StaticIps[1] is %s, expected %s", o.StaticIps[1], i2)
		}

		return nil
	}
}

func testAccPanosPanoramaEthernetInterfaceDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_ethernet_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, _, name := parsePanoramaEthernetInterfaceId(rs.Primary.ID)
			_, err := pano.Network.EthernetInterface.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaEthernetInterfaceConfig(tmpl, n, ls, c, i1, i2 string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "test" {
    name = %q
    template = panos_panorama_template.x.name
    mode = "layer3"
    link_state = %q
    comment = %q
    static_ips = [%q, %q]
}
`, tmpl, n, ls, c, i1, i2)
}
