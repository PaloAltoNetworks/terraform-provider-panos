package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/ospf/area/vlink"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source test (listing).
func TestAccPanosDsOspfAreaVirtualLinkList(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	area := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAreaVirtualLinkConfig(tmpl, vr, area, name),
				Check:  checkDataSourceListing("panos_ospf_area_virtual_links"),
			},
		},
	})
}

// Data source test.
func TestAccPanosDsOspfAreaVirtualLink(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	area := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAreaVirtualLinkConfig(tmpl, vr, area, name),
				Check: checkDataSource("panos_ospf_area_virtual_link", []string{
					"name",
					"enable",
					"neighbor_id",
					"transit_area_id",
					"hello_interval",
					"dead_counts",
					"retransmit_interval",
					"transit_delay",
				}),
			},
		},
	})
}

func testAccDsOspfAreaVirtualLinkConfig(tmpl, vr, area, name string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_ospf_area_virtual_links" "test" {
    template = panos_ospf_area_virtual_link.x.template
    virtual_router = panos_ospf_area_virtual_link.x.virtual_router
    ospf_area = panos_ospf_area_virtual_link.x.ospf_area
}

data "panos_ospf_area_virtual_link" "test" {
    template = panos_ospf_area_virtual_link.x.template
    virtual_router = panos_ospf_area_virtual_link.x.virtual_router
    ospf_area = panos_ospf_area_virtual_link.x.ospf_area
    name = panos_ospf_area_virtual_link.x.name
}

resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
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

resource "panos_ospf_area_virtual_link" "x" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = %q
    neighbor_id = "10.20.40.80"
    transit_area_id = panos_ospf_area.x.name
}
`, tmpl, vr, area, name)
	}

	return fmt.Sprintf(`
data "panos_ospf_area_virtual_links" "test" {
    virtual_router = panos_ospf_area_virtual_link.x.virtual_router
    ospf_area = panos_ospf_area_virtual_link.x.ospf_area
}

data "panos_ospf_area_virtual_link" "test" {
    virtual_router = panos_ospf_area_virtual_link.x.virtual_router
    ospf_area = panos_ospf_area_virtual_link.x.ospf_area
    name = panos_ospf_area_virtual_link.x.name
}

resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
}

resource "panos_ospf_area_virtual_link" "x" {
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = %q
    neighbor_id = "10.20.40.80"
    transit_area_id = panos_ospf_area.x.name
}
`, vr, area, name)
}

// Resource tests.
func TestAccPanosOspfAreaVirtualLink(t *testing.T) {
	var o vlink.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	area := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosOspfAreaVirtualLinkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOspfAreaVirtualLinkConfig(tmpl, vr, area, name, "10.5.7.2", area, true, 2000, 18, 1000, 500),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaVirtualLinkExists("panos_ospf_area_virtual_link.test", &o),
					testAccCheckPanosOspfAreaVirtualLinkAttributes(&o, name, "10.5.7.2", area, true, 2000, 18, 1000, 500),
				),
			},
			{
				Config: testAccOspfAreaVirtualLinkConfig(tmpl, vr, area, name, "10.6.8.3", area, false, 1999, 17, 999, 499),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaVirtualLinkExists("panos_ospf_area_virtual_link.test", &o),
					testAccCheckPanosOspfAreaVirtualLinkAttributes(&o, name, "10.6.8.3", area, false, 1999, 17, 999, 499),
				),
			},
		},
	})
}

func testAccCheckPanosOspfAreaVirtualLinkExists(n string, o *vlink.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v vlink.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vr, area, name := parseFirewallOspfAreaVirtualLinkId(rs.Primary.ID)
			v, err = con.Network.OspfAreaVirtualLink.Get(vr, area, name)
		case *pango.Panorama:
			tmpl, ts, vr, area, name := parsePanoramaOspfAreaVirtualLinkId(rs.Primary.ID)
			v, err = con.Network.OspfAreaVirtualLink.Get(tmpl, ts, vr, area, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosOspfAreaVirtualLinkAttributes(o *vlink.Entry, name, nid, taid string, en bool, hi, dc, ri, td int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		if o.NeighborId != nid {
			return fmt.Errorf("Neighbor ID is %q, not %q", o.NeighborId, nid)
		}

		if o.TransitAreaId != taid {
			return fmt.Errorf("Transit area ID is %q, not %q", o.TransitAreaId, taid)
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

		return nil
	}
}

func testAccPanosOspfAreaVirtualLinkDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ospf_area_virtual_link" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vr, area, name := parseFirewallOspfAreaVirtualLinkId(rs.Primary.ID)
				_, err = con.Network.OspfAreaVirtualLink.Get(vr, area, name)
			case *pango.Panorama:
				tmpl, ts, vr, area, name := parsePanoramaOspfAreaVirtualLinkId(rs.Primary.ID)
				_, err = con.Network.OspfAreaVirtualLink.Get(tmpl, ts, vr, area, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccOspfAreaVirtualLinkConfig(tmpl, vr, area, name, nid, taid string, en bool, hi, dc, ri, td int) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
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

resource "panos_ospf_area_virtual_link" "test" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = %q
    enable = %t
    neighbor_id = %q
    transit_area_id = %q
    hello_interval = %d
    dead_counts = %d
    retransmit_interval = %d
    transit_delay = %d
}
`, tmpl, vr, area, name, en, nid, taid, hi, dc, ri, td)
	}

	return fmt.Sprintf(`
resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
}

resource "panos_ospf_area_virtual_link" "test" {
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = %q
    enable = %t
    neighbor_id = %q
    transit_area_id = %q
    hello_interval = %d
    dead_counts = %d
    retransmit_interval = %d
    transit_delay = %d
}
`, vr, area, name, en, nid, taid, hi, dc, ri, td)
}
