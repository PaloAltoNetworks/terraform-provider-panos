package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/ospf"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source tests.
func TestAccPanosDsOspf(t *testing.T) {
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsOspfConfig(tmpl, vr),
				Check: checkDataSource("panos_ospf", []string{
					"enable",
					"router_id",
					"enable_graceful_restart",
					"grace_period",
					"helper_enable",
					"lsa_interval",
					"max_neighbor_restart_time",
					"spf_calculation_delay",
				}),
			},
		},
	})
}

func testAccDsOspfConfig(tmpl, vr string) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
data "panos_ospf" "test" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
}

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
    enable = true
    router_id = "10.5.7.9"
    enable_graceful_restart = true
    grace_period = 121
    helper_enable = true
    lsa_interval = 3
    max_neighbor_restart_time = 141
    spf_calculation_delay = 4
}
`, tmpl, vr)
	}

	return fmt.Sprintf(`
data "panos_ospf" "test" {
    virtual_router = panos_ospf.x.virtual_router
}

resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "x" {
    virtual_router = panos_virtual_router.x.name
    enable = true
    router_id = "10.5.7.9"
    enable_graceful_restart = true
    grace_period = 121
    helper_enable = true
    lsa_interval = 3
    max_neighbor_restart_time = 141
    spf_calculation_delay = 4
}
`, vr)
}

// Resource tests.
func TestAccPanosOspf(t *testing.T) {
	var o ospf.Config
	tmpl := fmt.Sprintf("tf%s", acctest.RandString(6))
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosOspfDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOspfConfig(tmpl, vr, "10.2.3.5", false, true, false, true, false, true, false, 4, 3, 121, 141),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfExists("panos_ospf.test", &o),
					testAccCheckPanosOspfAttributes(&o, "10.2.3.5", false, true, false, true, false, true, false, 4, 3, 121, 141),
				),
			},
			{
				Config: testAccOspfConfig(tmpl, vr, "10.5.8.13", true, false, true, false, true, false, true, 5, 4, 120, 140),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosOspfExists("panos_ospf.test", &o),
					testAccCheckPanosOspfAttributes(&o, "10.5.8.13", true, false, true, false, true, false, true, 5, 4, 120, 140),
				),
			},
		},
	})
}

func testAccCheckPanosOspfExists(n string, o *ospf.Config) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v ospf.Config

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vr := rs.Primary.ID
			v, err = con.Network.OspfConfig.Get(vr)
		case *pango.Panorama:
			tmpl, ts, vr := parsePanoramaOspfId(rs.Primary.ID)
			v, err = con.Network.OspfConfig.Get(tmpl, ts, vr)
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosOspfAttributes(o *ospf.Config, rid string, enable, rdr, ardr, rfc1583, egr, he, slc bool, scd, li float64, gp, mnrt int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Enable != enable {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, enable)
		}

		if o.RouterId != rid {
			return fmt.Errorf("Router ID is %q, not %q", o.RouterId, rid)
		}

		if o.RejectDefaultRoute != rdr {
			return fmt.Errorf("Reject default route is %t, not %t", o.RejectDefaultRoute, rdr)
		}

		if o.AllowRedistributeDefaultRoute != ardr {
			return fmt.Errorf("Allow redistribute default route is %t, not %t", o.AllowRedistributeDefaultRoute, ardr)
		}

		if o.Rfc1583 != rfc1583 {
			return fmt.Errorf("rfc1583 is %t, not %t", o.Rfc1583, rfc1583)
		}

		if o.SpfCalculationDelay != scd {
			return fmt.Errorf("SPF calcuation delay is %f, not %f", o.SpfCalculationDelay, scd)
		}

		if o.LsaInterval != li {
			return fmt.Errorf("LSA interval is %f, not %f", o.LsaInterval, li)
		}

		if o.EnableGracefulRestart != egr {
			return fmt.Errorf("Enable graceful restart is %t, not %t", o.EnableGracefulRestart, egr)
		}

		if o.GracePeriod != gp {
			return fmt.Errorf("Grace period is %d, not %d", o.GracePeriod, gp)
		}

		if o.HelperEnable != he {
			return fmt.Errorf("Helper enable is %t, not %t", o.HelperEnable, he)
		}

		if o.StrictLsaChecking != slc {
			return fmt.Errorf("Strict LSA checking is %t, not %t", o.StrictLsaChecking, slc)
		}

		if o.MaxNeighborRestartTime != mnrt {
			return fmt.Errorf("Max neighbor restart time is %d, not %d", o.MaxNeighborRestartTime, mnrt)
		}

		return nil
	}
}

func testAccPanosOspfDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_ospf" {
			continue
		}

		if rs.Primary.ID != "" {
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vr := rs.Primary.ID
				_, err = con.Network.OspfConfig.Get(vr)
			case *pango.Panorama:
				tmpl, ts, vr := parsePanoramaOspfId(rs.Primary.ID)
				_, err = con.Network.OspfConfig.Get(tmpl, ts, vr)
			}
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccOspfConfig(tmpl, vr, rid string, enable, rdr, ardr, rfc1583, egr, he, slc bool, scd, li float64, gp, mnrt int) string {
	if testAccIsPanorama {
		return fmt.Sprintf(`
resource "panos_panorama_template" "x" {
    name = %q
}

resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    name = %q
}

resource "panos_ospf" "test" {
    template = panos_panorama_template.x.name
    virtual_router = panos_panorama_virtual_router.x.name
    enable = %t
    router_id = %q
    reject_default_route = %t
    allow_redistribute_default_route = %t
    rfc_1583 = %t
    spf_calculation_delay = %f
    lsa_interval = %f
    enable_graceful_restart = %t
    grace_period = %d
    helper_enable = %t
    strict_lsa_checking = %t
    max_neighbor_restart_time = %d
}
`, tmpl, vr, enable, rid, rdr, ardr, rfc1583, scd, li, egr, gp, he, slc, mnrt)
	}

	return fmt.Sprintf(`
resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_ospf" "test" {
    virtual_router = panos_virtual_router.x.name
    enable = %t
    router_id = %q
    reject_default_route = %t
    allow_redistribute_default_route = %t
    rfc_1583 = %t
    spf_calculation_delay = %f
    lsa_interval = %f
    enable_graceful_restart = %t
    grace_period = %d
    helper_enable = %t
    strict_lsa_checking = %t
    max_neighbor_restart_time = %d
}
`, vr, enable, rid, rdr, ardr, rfc1583, scd, li, egr, gp, he, slc, mnrt)
}
