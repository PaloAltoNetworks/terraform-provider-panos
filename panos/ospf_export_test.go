package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/ospf/exp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Data source test (listing).
func TestAccPanosDsOspfExportList(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("10.%d.%d.0/24", acctest.RandInt()%50+50, acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfExportConfig(tmpl, vr, name),
				Check:  checkDataSourceListing("panos_ospf_exports"),
			},
		},
	})
}

// Data source test.
func TestAccPanosDsOspfExport(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("10.%d.%d.0/24", acctest.RandInt()%50+50, acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfExportConfig(tmpl, vr, name),
				Check: checkDataSource("panos_ospf_export", []string{
					"name", "path_type", "tag", "metric",
				}),
			},
		},
	})
}

func testAccDsOspfExportConfig(tmpl, vr, name string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_ospf_exports" "test" {
    template = panos_ospf_export.x.template
    virtual_router = panos_ospf_export.x.virtual_router
}

data "panos_ospf_export" "test" {
    template = panos_ospf_export.x.template
    virtual_router = panos_ospf_export.x.virtual_router
    name = panos_ospf_export.x.name
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

resource "panos_ospf_export" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    tag = "10.5.15.151"
    metric = 42
}
`, tmpl, vr, name)
	}

	return fmt.Sprintf(`
data "panos_ospf_exports" "test" {
    virtual_router = panos_ospf_export.x.virtual_router
}

data "panos_ospf_export" "test" {
    virtual_router = panos_ospf_export.x.virtual_router
    name = panos_ospf_export.x.name
}

resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_export" "x" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    tag = "10.5.15.151"
    metric = 42
}
`, vr, name)
}

// Resource tests.
func TestAccPanosOspfExport(t *testing.T) {
	var o exp.Entry
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("10.%d.%d.0/24", acctest.RandInt()%50+50, acctest.RandInt()%50+50)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosOspfExportDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOspfExportConfig(tmpl, vr, name, exp.PathTypeExt1, "10.5.2.7", 50),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfExportExists("panos_ospf_export.test", &o),
					testAccCheckPanosOspfExportAttributes(&o, name, exp.PathTypeExt1, "10.5.2.7", 50),
				),
			},
			{
				Config: testAccOspfExportConfig(tmpl, vr, name, exp.PathTypeExt2, "10.6.3.8", 42),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfExportExists("panos_ospf_export.test", &o),
					testAccCheckPanosOspfExportAttributes(&o, name, exp.PathTypeExt2, "10.6.3.8", 42),
				),
			},
		},
	})
}

func testAccCheckPanosOspfExportExists(n string, o *exp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v exp.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vr, name := parseFirewallOspfExportId(rs.Primary.ID)
			v, err = con.Network.OspfExport.Get(vr, name)
		case *pango.Panorama:
			tmpl, ts, vr, name := parsePanoramaOspfExportId(rs.Primary.ID)
			v, err = con.Network.OspfExport.Get(tmpl, ts, vr, name)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosOspfExportAttributes(o *exp.Entry, name, pt, tag string, metric int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.PathType != pt {
			return fmt.Errorf("Path type is %q, not %q", o.PathType, pt)
		}

		if o.Tag != tag {
			return fmt.Errorf("Tag is %q, not %q", o.Tag, tag)
		}

		if o.Metric != metric {
			return fmt.Errorf("Metric is %d, not %d", o.Metric, metric)
		}

		return nil
	}
}

func testAccPanosOspfExportDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ospf_export" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vr, name := parseFirewallOspfExportId(rs.Primary.ID)
				_, err = con.Network.OspfExport.Get(vr, name)
			case *pango.Panorama:
				tmpl, ts, vr, name := parsePanoramaOspfExportId(rs.Primary.ID)
				_, err = con.Network.OspfExport.Get(tmpl, ts, vr, name)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccOspfExportConfig(tmpl, vr, name, pt, tag string, metric int) string {
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
    template = panos_panorama_template.x.name
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false
}

resource "panos_ospf_export" "test" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    path_type = %q
    tag = %q
    metric = %d
}
`, tmpl, vr, name, pt, tag, metric)
	}

	return fmt.Sprintf(`
resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = false
}

resource "panos_ospf_export" "test" {
    virtual_router = panos_ospf.x.virtual_router
    name = %q
    path_type = %q
    tag = %q
    metric = %d
}
`, vr, name, pt, tag, metric)
}
