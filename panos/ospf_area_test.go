package panos

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/area"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source test (listing).
func TestAccPanosDsOspfAreaList(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAreaConfig(tmpl, vr, name),
				Check:  checkDataSourceListing("panos_ospf_areas"),
			},
		},
	})
}

// Data source test.
func TestAccPanosDsOspfArea(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfAreaConfig(tmpl, vr, name),
				Check: checkDataSource("panos_ospf_area", []string{
					"name",
					"type",
					"accept_summary",
					"default_route_advertise",
					"advertise_metric",
					"advertise_type",
				}),
			},
		},
	})
}

func testAccDsOspfAreaConfig(tmpl, vr, name string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_ospf_areas" "test" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
}

data "panos_ospf_area" "test" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    name = panos_ospf_area.x.name
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
    type = %q
    accept_summary = true
    default_route_advertise = true
    advertise_metric = 50
    advertise_type = %q
}
`, tmpl, vr, name, area.TypeNssa, area.AdvertiseTypeExt2)
	}

	return fmt.Sprintf(`
data "panos_ospf_areas" "test" {
    virtual_router = panos_ospf_area.x.virtual_router
}

data "panos_ospf_area" "test" {
    virtual_router = panos_ospf_area.x.virtual_router
    name = panos_ospf_area.x.name
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
    type = %q
    accept_summary = true
    default_route_advertise = true
    advertise_metric = 50
    advertise_type = %q
}
`, vr, name, area.TypeNssa, area.AdvertiseTypeExt2)
}

// Resource tests.
func TestAccPanosOspfArea(t *testing.T) {
	var o area.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("10.%d.%d.%d",
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
		acctest.RandInt()%50+50,
	)

	ranges := []area.Range{
		{"10.1.1.0/24", area.ActionSuppress},
		{"10.2.3.0/24", area.ActionAdvertise},
		{"10.3.3.0/24", area.ActionAdvertise},
		{"10.4.1.0/24", area.ActionSuppress},
		{"10.5.30.0/24", area.ActionAdvertise},
		{"10.6.2.0/24", area.ActionSuppress},
		{"10.7.77.0/24", area.ActionAdvertise},
		{"10.8.88.0/24", area.ActionSuppress},
		{"10.9.191.0/24", area.ActionAdvertise},
		{"10.10.10.0/24", area.ActionSuppress},
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosOspfAreaDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeNormal, "", false, false, 0, nil, []area.Range{ranges[0], ranges[1]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeNormal, "", false, false, 0, nil, []area.Range{ranges[0], ranges[1]}),
				),
			},
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeNormal, "", false, false, 0, nil, []area.Range{ranges[0]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeNormal, "", false, false, 0, nil, []area.Range{ranges[0]}),
				),
			},
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeStub, "", false, true, 42, nil, []area.Range{ranges[7], ranges[6], ranges[5]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeStub, "", false, true, 42, nil, []area.Range{ranges[7], ranges[6], ranges[5]}),
				),
			},
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeStub, "", true, false, 0, nil, []area.Range{ranges[2], ranges[3], ranges[4]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeStub, "", true, false, 0, nil, []area.Range{ranges[2], ranges[3], ranges[4]}),
				),
			},
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeNssa, "", true, false, 0, []area.Range{ranges[0], ranges[1]}, []area.Range{ranges[2], ranges[3]}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeNssa, "", true, false, 0, []area.Range{ranges[0], ranges[1]}, []area.Range{ranges[2], ranges[3]}),
				),
			},
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeNssa, area.AdvertiseTypeExt2, true, true, 101, ranges[0:3], ranges[3:6]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeNssa, area.AdvertiseTypeExt2, true, true, 101, ranges[0:3], ranges[3:6]),
				),
			},
			{
				Config: testAccOspfAreaConfig(tmpl, vr, name, area.TypeNssa, area.AdvertiseTypeExt1, true, true, 64, ranges[0:2], ranges[2:10]),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfAreaExists("panos_ospf_area.test", &o),
					testAccCheckPanosOspfAreaAttributes(&o, name, area.TypeNssa, area.AdvertiseTypeExt1, true, true, 64, ranges[0:2], ranges[2:10]),
				),
			},
		},
	})
}

func testAccCheckPanosOspfAreaExists(n string, o *area.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v area.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vr, name := parseFirewallOspfAreaId(rs.Primary.ID)
			v, err = con.Network.OspfArea.Get(vr, name)
		case *pango.Panorama:
			tmpl, ts, vr, name := parsePanoramaOspfAreaId(rs.Primary.ID)
			v, err = con.Network.OspfArea.Get(tmpl, ts, vr, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosOspfAreaAttributes(o *area.Entry, name, t1, t2 string, as, dra bool, am int, erList, rList []area.Range) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.Type != t1 {
			return fmt.Errorf("Type is %q, not %q", o.Type, t1)
		}

		if o.AcceptSummary != as {
			return fmt.Errorf("Accept summary is %t, not %t", o.AcceptSummary, as)
		}

		if o.DefaultRouteAdvertise != dra {
			return fmt.Errorf("Default route advertise is %t, not %t", o.DefaultRouteAdvertise, dra)
		}

		if o.AdvertiseMetric != am {
			return fmt.Errorf("Advertise metric is %d, not %d", o.AdvertiseMetric, am)
		}

		if o.AdvertiseType != t2 {
			return fmt.Errorf("Advertise type is %q, not %q", o.AdvertiseType, t2)
		}

		if len(o.ExtRanges) != len(erList) {
			return fmt.Errorf("Ext ranges len %d is not %d", len(o.ExtRanges), len(erList))
		}

		for i := range o.ExtRanges {
			if o.ExtRanges[i].Network != erList[i].Network {
				return fmt.Errorf("ext range network %d is %q, not %q", i, o.ExtRanges[i].Network, erList[i].Network)
			}

			if o.ExtRanges[i].Action != erList[i].Action {
				return fmt.Errorf("ext range action %d is %q, not %q", i, o.ExtRanges[i].Action, erList[i].Action)
			}
		}

		if len(o.Ranges) != len(rList) {
			return fmt.Errorf("Ranges len %d is not %d", len(o.Ranges), len(rList))
		}

		for i := range o.Ranges {
			if o.Ranges[i].Network != rList[i].Network {
				return fmt.Errorf("range network %d is %q, not %q", i, o.Ranges[i].Network, rList[i].Network)
			}

			if o.Ranges[i].Action != rList[i].Action {
				return fmt.Errorf("range action %d is %q, not %q", i, o.Ranges[i].Action, rList[i].Action)
			}
		}

		return nil
	}
}

func testAccPanosOspfAreaDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ospf_area" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vr, name := parseFirewallOspfAreaId(rs.Primary.ID)
				_, err = con.Network.OspfArea.Get(vr, name)
			case *pango.Panorama:
				tmpl, ts, vr, name := parsePanoramaOspfAreaId(rs.Primary.ID)
				_, err = con.Network.OspfArea.Get(tmpl, ts, vr, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccOspfAreaConfig(tmpl, vr, name, t1, t2 string, as, dra bool, am int, erList, rList []area.Range) string {
	var b strings.Builder
	b.WriteString(`
resource "panos_ospf_area" "test" {`)
	if testAccIsPanorama {
		b.WriteString(`
    template = panos_ospf.x.template`)
	}
	fmt.Fprintf(&b, `
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    type = %q
    accept_summary = %t
    default_route_advertise = %t
    advertise_metric = %d
    advertise_type = %q`, name, t1, as, dra, am, t2)
	if len(erList) > 0 {
		for _, x := range erList {
			fmt.Fprintf(&b, `
    ext_range {
        network = %q
        action = %q
    }`, x.Network, x.Action)
		}
	}
	if len(rList) > 0 {
		for _, x := range rList {
			fmt.Fprintf(&b, `
    range {
        network = %q
        action = %q
    }`, x.Network, x.Action)
		}
	}
	b.WriteString(`
}`)

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

%s
`, tmpl, vr, b.String())
	}

	return fmt.Sprintf(`
resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

%s
`, vr, b.String())
}
