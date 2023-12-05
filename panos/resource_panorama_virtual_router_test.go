package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/router"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaVirtualRouter_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o router.Entry
	name := fmt.Sprintf("tf%s", acctest.RandString(6))
	tmpl := fmt.Sprintf("tfTemplate%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaVirtualRouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaVirtualRouterConfig(tmpl, name, "ethernet1/1", 10, 10, 30, 110, 30, 110, 200, 20, 120),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVirtualRouterExists("panos_panorama_virtual_router.test", &o),
					testAccCheckPanosPanoramaVirtualRouterAttributes(&o, name, "ethernet1/1", 10, 10, 30, 110, 30, 110, 200, 20, 120),
				),
			},
			{
				Config: testAccPanoramaVirtualRouterConfig(tmpl, name, "ethernet1/2", 11, 12, 33, 114, 35, 116, 207, 28, 129),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaVirtualRouterExists("panos_panorama_virtual_router.test", &o),
					testAccCheckPanosPanoramaVirtualRouterAttributes(&o, name, "ethernet1/2", 11, 12, 33, 114, 35, 116, 207, 28, 129),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaVirtualRouterExists(n string, o *router.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		pano := testAccProvider.Meta().(*pango.Panorama)
		tmpl, ts, _, name := parseVirtualRouterId(rs.Primary.ID)
		v, err := pano.Network.VirtualRouter.Get(tmpl, ts, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosPanoramaVirtualRouterAttributes(o *router.Entry, name, eth string, sd, sid, oid, oed, ov3id, ov3ed, id, ed, rd int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if len(o.Interfaces) != 1 || o.Interfaces[0] != eth {
			return fmt.Errorf("Interfaces is %#v, expected [%s]", o.Interfaces, eth)
		}

		if o.StaticDist != sd {
			return fmt.Errorf("Static is %d, expected %d", o.StaticDist, sd)
		}

		if o.StaticIpv6Dist != sid {
			return fmt.Errorf("Static IPV6 is %d, expected %d", o.StaticIpv6Dist, sid)
		}

		if o.OspfIntDist != oid {
			return fmt.Errorf("OSPF Int is %d, expected %d", o.OspfIntDist, oid)
		}

		if o.OspfExtDist != oed {
			return fmt.Errorf("OSPF Ext is %d, expected %d", o.OspfExtDist, oed)
		}

		if o.Ospfv3IntDist != ov3id {
			return fmt.Errorf("OSPFv3 Int is %d, expected %d", o.Ospfv3IntDist, ov3id)
		}

		if o.Ospfv3ExtDist != ov3ed {
			return fmt.Errorf("OSPFv3 Ext is %d, expected %d", o.Ospfv3ExtDist, ov3ed)
		}

		if o.IbgpDist != id {
			return fmt.Errorf("IBGP is %d, expected %d", o.IbgpDist, id)
		}

		if o.EbgpDist != ed {
			return fmt.Errorf("EBGP is %d, expected %d", o.EbgpDist, ed)
		}

		if o.RipDist != rd {
			return fmt.Errorf("RIP is %d, expected %d", o.RipDist, rd)
		}

		return nil
	}
}

func testAccPanosPanoramaVirtualRouterDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_virtual_router" {
			continue
		}

		if rs.Primary.ID != "" {
			tmpl, ts, _, name := parseVirtualRouterId(rs.Primary.ID)
			_, err := pano.Network.VirtualRouter.Get(tmpl, ts, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccPanoramaVirtualRouterConfig(tmpl, name, eth string, sd, sid, oid, oed, ov3id, ov3ed, id, ed, rd int) string {
	return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_ethernet_interface" "eth1_1" {
    template = panos_panorama_template.x.name
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_panorama_ethernet_interface" "eth1_2" {
    template = panos_panorama_template.x.name
    name = "ethernet1/2"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_panorama_virtual_router" "test" {
    template = panos_panorama_template.x.name
    name = "%s"
    interfaces = ["%s"]
    static_dist = %d
    static_ipv6_dist = %d
    ospf_int_dist = %d
    ospf_ext_dist = %d
    ospfv3_int_dist = %d
    ospfv3_ext_dist = %d
    ibgp_dist = %d
    ebgp_dist = %d
    rip_dist = %d
    depends_on = ["panos_panorama_ethernet_interface.eth1_1", "panos_panorama_ethernet_interface.eth1_2"]
}
`, tmpl, name, eth, sd, sid, oid, oed, ov3id, ov3ed, id, ed, rd)
}
