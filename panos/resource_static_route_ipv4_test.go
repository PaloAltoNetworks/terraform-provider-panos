package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/route/static/ipv4"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosStaticRouteIpv4(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o ipv4.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosStaticRouteIpv4Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStaticRouteIpv4Config(vr, name, "10.1.7.0/32", ipv4.NextHopIpAddress, "10.1.7.4", 42, 21, ipv4.RouteTableUnicast),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosStaticRouteIpv4Exists("panos_static_route_ipv4.test", &o),
					testAccCheckPanosStaticRouteIpv4Attributes(&o, name, "10.1.7.0/32", ipv4.NextHopIpAddress, "10.1.7.4", 42, 21, ipv4.RouteTableUnicast),
				),
			},
			{
				Config: testAccStaticRouteIpv4Config(vr, name, "10.1.9.0/32", "", "", 46, 23, ipv4.RouteTableBoth),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosStaticRouteIpv4Exists("panos_static_route_ipv4.test", &o),
					testAccCheckPanosStaticRouteIpv4Attributes(&o, name, "10.1.9.0/32", "", "", 46, 23, ipv4.RouteTableBoth),
				),
			},
		},
	})
}

func testAccCheckPanosStaticRouteIpv4Exists(n string, o *ipv4.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, name := parseStaticRouteIpv4Id(rs.Primary.ID)
		v, err := fw.Network.StaticRoute.Get(vr, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosStaticRouteIpv4Attributes(o *ipv4.Entry, name, dest, ty, nh string, ad, metric int, rt string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.Destination != dest {
			return fmt.Errorf("Destination is %q, expected %q", o.Destination, dest)
		}

		if o.Interface != "ethernet1/6" {
			return fmt.Errorf("Interface is %q, expected \"ethernet1/6\"", o.Interface)
		}

		if o.Type != ty {
			return fmt.Errorf("Type is %q, expected %q", o.Type, ty)
		}

		if o.NextHop != nh {
			return fmt.Errorf("Next hop is %q, expected %q", o.NextHop, nh)
		}

		if o.AdminDistance != ad {
			return fmt.Errorf("Admin dist is %d, expected %d", o.AdminDistance, ad)
		}

		if o.Metric != metric {
			return fmt.Errorf("Metric is %d, expected %d", o.Metric, metric)
		}

		if o.RouteTable != rt {
			return fmt.Errorf("Route table is %q, expected %q", o.RouteTable, rt)
		}

		return nil
	}
}

func testAccPanosStaticRouteIpv4Destroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_static_route_ipv4" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, name := parseStaticRouteIpv4Id(rs.Primary.ID)
			_, err := fw.Network.StaticRoute.Get(vr, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccStaticRouteIpv4Config(vr, name, dest, ty, nh string, ad, metric int, rt string) string {
	return fmt.Sprintf(`
resource "panos_ethernet_interface" "eth6" {
    name = "ethernet1/6"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
}

resource "panos_virtual_router" "vr" {
    name = %q
    interfaces = [panos_ethernet_interface.eth6.name]
}

resource "panos_static_route_ipv4" "test" {
    name = %q
    virtual_router = panos_virtual_router.vr.name
    destination = %q
    interface = panos_ethernet_interface.eth6.name
    type = %q
    next_hop = %q
    admin_distance = %d
    metric = %d
    route_table = %q
}
`, vr, name, dest, ty, nh, ad, metric, rt)
}
