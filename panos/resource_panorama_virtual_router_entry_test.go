package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosPanoramaVirtualRouterEntry_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o router.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	eth_name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%7+1)
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaVirtualRouterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaVirtualRouterEntryConfig(tmpl, eth_name, vr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVirtualRouterEntryExists("panos_panorama_virtual_router_entry.test", &o),
					testAccCheckPanosPanoramaVirtualRouterEntryAttributes(&o, eth_name),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaVirtualRouterEntryExists(n string, o *router.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, vr, _ := parsePanoramaVirtualRouterEntryId(rs.Primary.ID)
		v, err := pano.Network.VirtualRouter.Get(tmpl, ts, vr)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaVirtualRouterEntryAttributes(o *router.Entry, eth_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(o.Interfaces) != 1 || o.Interfaces[0] != eth_name {
			return fmt.Errorf("Virtual router interfaces is %#v, not [%s]", o.Interfaces, eth_name)
		}

		return nil
	}
}

func testAccPanosPanoramaVirtualRouterEntryDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_virtual_router_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, vr, _ := parsePanoramaVirtualRouterEntryId(rs.Primary.ID)
			_, err := pano.Network.VirtualRouter.Get(tmpl, ts, vr)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaVirtualRouterEntryConfig(tmpl, eth_name, vr string) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "tmpl" {
    name = %q
    description = "for template acc test"
}

resource "panos_panorama_ethernet_interface" "eth" {
    template = "${panos_panorama_template.tmpl.name}"
    name = %q
    mode = "layer3"
}

resource "panos_panorama_virtual_router" "vr" {
    template = "${panos_panorama_template.tmpl.name}"
    name = %q
}

resource "panos_panorama_virtual_router_entry" "test" {
    template = "${panos_panorama_template.tmpl.name}"
    virtual_router = "${panos_panorama_virtual_router.vr.name}"
    interface = "${panos_panorama_ethernet_interface.eth.name}"
}
`, tmpl, eth_name, vr)
}
