package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosBgp_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o bgp.Config
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpConfig(name, "5.5.5.5", "42", true, false, true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpExists("panos_bgp.test", &o),
					testAccCheckPanosBgpAttributes(&o, "5.5.5.5", "42", true, false, true, false),
				),
			},
			{
				Config: testAccBgpConfig(name, "6.6.6.6", "420", false, true, false, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpExists("panos_bgp.test", &o),
					testAccCheckPanosBgpAttributes(&o, "6.6.6.6", "420", false, true, false, true),
				),
			},
		},
	})
}

func testAccCheckPanosBgpExists(n string, o *bgp.Config) resource.TestCheckFunc {
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
		v, err := fw.Network.BgpConfig.Get(name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpAttributes(o *bgp.Config, ri, an string, en, rdr, ir, am bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.RouterId != ri {
			return fmt.Errorf("Router Id is %q, expected %q", o.RouterId, ri)
		}

		if o.AsNumber != an {
			return fmt.Errorf("AS number is %q, expected %q", o.AsNumber, an)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.RejectDefaultRoute != rdr {
			return fmt.Errorf("Reject default route is %t, not %t", o.RejectDefaultRoute, rdr)
		}

		if o.InstallRoute != ir {
			return fmt.Errorf("Install route is %t, not %t", o.InstallRoute, ir)
		}

		if o.AggregateMed != am {
			return fmt.Errorf("Aggregate MED is %t, not %t", o.AggregateMed, am)
		}

		return nil
	}
}

func testAccPanosBgpDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp" {
			continue
		}

		if rs.Primary.ID != "" {
			name := rs.Primary.ID
			_, err := fw.Network.BgpConfig.Get(name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpConfig(name, ri, an string, en, rdr, ir, am bool) string {
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

resource "panos_bgp" "test" {
    virtual_router = panos_virtual_router.vr.name
    router_id = %q
    as_number = %q
    enable = %t
    reject_default_route = %t
    install_route = %t
    aggregate_med = %t
}
`, name, ri, an, en, rdr, ir, am)
}
