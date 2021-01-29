package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area/iface"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source test (listing).
func TestAccPanosDsOspfAreaInterfaceList(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	area := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)
	name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%3+1)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAreaInterfaceConfig(tmpl, vr, area, name),
				Check:  checkDataSourceListing("panos_ospf_area_interfaces"),
			},
		},
	})
}

// Data source test.
func TestAccPanosDsOspfAreaInterface(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	area := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)
	name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%3+1)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAreaInterfaceConfig(tmpl, vr, area, name),
				Check: checkDataSource("panos_ospf_area_interface", []string{
					"name",
					"enable",
					"passive",
					"link_type",
					"metric",
					"priority",
					"hello_interval",
					"dead_counts",
					"retransmit_interval",
					"transit_delay",
					"grace_restart_delay",
				}),
			},
		},
	})
}

func testAccDsOspfAreaInterfaceConfig(tmpl, vr, area, name string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_ospf_area_interfaces" "test" {
    template = panos_ospf_area_interface.x.template
    virtual_router = panos_ospf_area_interface.x.virtual_router
    ospf_area = panos_ospf_area_interface.x.ospf_area
}

data "panos_ospf_area_interface" "test" {
    template = panos_ospf_area_interface.x.template
    virtual_router = panos_ospf_area_interface.x.virtual_router
    ospf_area = panos_ospf_area_interface.x.ospf_area
    name = panos_ospf_area_interface.x.name
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    interfaces = [panos_panorama_ethernet_interface.x.name]
    name = %q
}

resource "panos_ospf" "x" {
    template = panos_panorama_virtual_router.x.template
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = %q
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    mode = "layer3"
    vsys = "vsys1"
}

resource "panos_ospf_area_interface" "x" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = panos_panorama_ethernet_interface.x.name
    enable = true
    passive = true
}
`, tmpl, vr, area, name)
	}

	return fmt.Sprintf(`
data "panos_ospf_area_interfaces" "test" {
    virtual_router = panos_ospf_area_interface.x.virtual_router
    ospf_area = panos_ospf_area_interface.x.ospf_area
}

data "panos_ospf_area_interface" "test" {
    virtual_router = panos_ospf_area_interface.x.virtual_router
    ospf_area = panos_ospf_area_interface.x.ospf_area
    name = panos_ospf_area_interface.x.name
}

resource "panos_virtual_router" "x" {
    name = %q
    interfaces = [panos_ethernet_interface.x.name]
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
}

resource "panos_ethernet_interface" "x" {
    name = %q
    mode = "layer3"
    vsys = "vsys1"
}

resource "panos_ospf_area_interface" "x" {
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = panos_ethernet_interface.x.name
}
`, vr, area, name)
}

// Resource tests.
func TestAccPanosOspfAreaInterface(t *testing.T) {
	var o iface.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	area := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)
	name := fmt.Sprintf("ethernet1/%d", acctest.RandInt()%3+1)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosOspfAreaInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOspfAreaInterfaceConfig(tmpl, vr, area, name, iface.LinkTypeBroadcast, true, false, 32768, 128, 1800, 10, 1800, 1800, 5, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaInterfaceExists("panos_ospf_area_interface.test", &o),
					testAccCheckPanosOspfAreaInterfaceAttributes(&o, name, iface.LinkTypeBroadcast, true, false, 32768, 128, 1800, 10, 1800, 1800, 5, nil),
				),
			},
			{
				Config: testAccOspfAreaInterfaceConfig(tmpl, vr, area, name, iface.LinkTypePointToPoint, false, true, 20000, 200, 300, 7, 400, 500, 6, nil),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaInterfaceExists("panos_ospf_area_interface.test", &o),
					testAccCheckPanosOspfAreaInterfaceAttributes(&o, name, iface.LinkTypePointToPoint, false, true, 20000, 200, 300, 7, 400, 500, 6, nil),
				),
			},
			{
				Config: testAccOspfAreaInterfaceConfig(tmpl, vr, area, name, iface.LinkTypePointToMultiPoint, false, true, 20000, 200, 300, 7, 400, 500, 6, []string{"10.1.5.151", "10.2.3.4"}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaInterfaceExists("panos_ospf_area_interface.test", &o),
					testAccCheckPanosOspfAreaInterfaceAttributes(&o, name, iface.LinkTypePointToMultiPoint, false, true, 20000, 200, 300, 7, 400, 500, 6, []string{"10.1.5.151", "10.2.3.4"}),
				),
			},
		},
	})
}

func testAccCheckPanosOspfAreaInterfaceExists(n string, o *iface.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v iface.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vr, area, name := parseFirewallOspfAreaInterfaceId(rs.Primary.ID)
			v, err = con.Network.OspfAreaInterface.Get(vr, area, name)
		case *pango.Panorama:
			tmpl, ts, vr, area, name := parsePanoramaOspfAreaInterfaceId(rs.Primary.ID)
			v, err = con.Network.OspfAreaInterface.Get(tmpl, ts, vr, area, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosOspfAreaInterfaceAttributes(o *iface.Entry, name, lt string, en, passive bool, metric, pri, hi, dc, ri, td, grd int, neighbors []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.LinkType != lt {
			return fmt.Errorf("Link type is %q, not %q", o.LinkType, lt)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.Passive != passive {
			return fmt.Errorf("Passive is %t, not %t", o.Passive, passive)
		}

		if o.Metric != metric {
			return fmt.Errorf("Metric is %d, not %d", o.Metric, metric)
		}

		if o.Priority != pri {
			return fmt.Errorf("Priority is %d, not %d", o.Priority, pri)
		}

		if o.HelloInterval != hi {
			return fmt.Errorf("Hellow interval is %d, not %d", o.HelloInterval, hi)
		}

		if o.DeadCounts != dc {
			return fmt.Errorf("Dead counts is %d, not %d", o.DeadCounts, dc)
		}

		if o.RetransmitInterval != ri {
			return fmt.Errorf("Retransmit interval is %d, not %d", o.RetransmitInterval, ri)
		}

		if o.TransitDelay != td {
			return fmt.Errorf("Transit delay is %d, not %d", o.TransitDelay, td)
		}

		if o.GraceRestartDelay != grd {
			return fmt.Errorf("Grace restart delay is %d, not %d", o.GraceRestartDelay, grd)
		}

		if len(o.Neighbors) != len(neighbors) {
			return fmt.Errorf("Neighbors len %d is not %d", len(o.Neighbors), len(neighbors))
		}

		for i := range o.Neighbors {
			if o.Neighbors[i] != neighbors[i] {
				return fmt.Errorf("Neighbor.%d is %q, not %q", i, o.Neighbors[i], neighbors[i])
			}
		}

		return nil
	}
}

func testAccPanosOspfAreaInterfaceDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ospf_area_interface" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vr, area, name := parseFirewallOspfAreaInterfaceId(rs.Primary.ID)
				_, err = con.Network.OspfAreaInterface.Get(vr, area, name)
			case *pango.Panorama:
				tmpl, ts, vr, area, name := parsePanoramaOspfAreaInterfaceId(rs.Primary.ID)
				_, err = con.Network.OspfAreaInterface.Get(tmpl, ts, vr, area, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccOspfAreaInterfaceConfig(tmpl, vr, area, name, lt string, en, passive bool, metric, pri, hi, dc, ri, td, grd int, neighbors []string) string {
	var b strings.Builder
	for num, x := range neighbors {
		if num != 0 {
			b.WriteString(", ")
		}
		fmt.Fprintf(&b, "%q", x)
	}

	if testAccIsPanorama {
		return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    interfaces = [panos_panorama_ethernet_interface.x.name]
    name = %q
}

resource "panos_ospf" "x" {
    template = panos_panorama_virtual_router.x.template
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = %q
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = %q
    mode = "layer3"
    vsys = "vsys1"
}

resource "panos_ospf_area_interface" "test" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = panos_panorama_ethernet_interface.x.name
    link_type = %q
    enable = %t
    passive = %t
    metric = %d
    priority = %d
    hello_interval = %d
    dead_counts = %d
    retransmit_interval = %d
    transit_delay = %d
    grace_restart_delay = %d
    neighbors = [%s]
}
`, tmpl, vr, area, name, lt, en, passive, metric, pri, hi, dc, ri, td, grd, b.String())
	}

	return fmt.Sprintf(`
resource "panos_virtual_router" "x" {
    name = %q
    interfaces = [panos_ethernet_interface.x.name]
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
}

resource "panos_ethernet_interface" "x" {
    name = %q
    mode = "layer3"
    vsys = "vsys1"
}

resource "panos_ospf_area_interface" "test" {
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = panos_ethernet_interface.x.name
    link_type = %q
    enable = %t
    passive = %t
    metric = %d
    priority = %d
    hello_interval = %d
    dead_counts = %d
    retransmit_interval = %d
    transit_delay = %d
    grace_restart_delay = %d
    neighbors = [%s]
}
`, vr, area, name, lt, en, passive, metric, pri, hi, dc, ri, td, grd, b.String())
}
