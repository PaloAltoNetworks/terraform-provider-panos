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

func TestAccPanosVirtualRouterEntry_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o router.Entry
	eth_name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%7+1)
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosVirtualRouterEntryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccVirtualRouterEntryConfig(eth_name, vr),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosVirtualRouterEntryExists("panos_virtual_router_entry.test", &o),
					testAccCheckPanosVirtualRouterEntryAttributes(&o, eth_name),
				),
			},
		},
	})
}

func testAccCheckPanosVirtualRouterEntryExists(n string, o *router.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, _ := parseVirtualRouterEntryId(rs.Primary.ID)
		v, err := fw.Network.VirtualRouter.Get(vr)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosVirtualRouterEntryAttributes(o *router.Entry, eth_name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(o.Interfaces) != 1 || o.Interfaces[0] != eth_name {
			return fmt.Errorf("Virtual router interfaces is %#v, not [%s]", o.Interfaces, eth_name)
		}

		return nil
	}
}

func testAccPanosVirtualRouterEntryDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_virtual_router_entry" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, _ := parseVirtualRouterEntryId(rs.Primary.ID)
			_, err := fw.Network.VirtualRouter.Get(vr)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccVirtualRouterEntryConfig(eth_name, vr string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "eth" {
    name = %q
    mode = "layer3"
}

resource "panos_virtual_router" "vr" {
    name = %q
}

resource "panos_virtual_router_entry" "test" {
    virtual_router = "${panos_virtual_router.vr.name}"
    interface = "${panos_ethernet_interface.eth.name}"
}
`, eth_name, vr)
}
