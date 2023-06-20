package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/netw/routing/protocol/bgp/conadv/filter/advertise"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosBgpConditionalAdvAdvertiseFilter_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o advertise.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	ca := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpConditionalAdvAdvertiseFilterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpConditionalAdvAdvertiseFilterConfig(vr, ca, name, "path1", "com1", "ext1", "21", "5.5.5.0/24", "5.5.6.0/24", false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpConditionalAdvAdvertiseFilterExists("panos_bgp_conditional_adv_advertise_filter.test", &o),
					testAccCheckPanosBgpConditionalAdvAdvertiseFilterAttributes(&o, name, "path1", "com1", "ext1", "21", "5.5.5.0/24", "5.5.6.0/24", false),
				),
			},
			{
				Config: testAccBgpConditionalAdvAdvertiseFilterConfig(vr, ca, name, "path2", "com2", "ext2", "7", "6.6.6.0/24", "6.6.7.0/24", true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpConditionalAdvAdvertiseFilterExists("panos_bgp_conditional_adv_advertise_filter.test", &o),
					testAccCheckPanosBgpConditionalAdvAdvertiseFilterAttributes(&o, name, "path2", "com2", "ext2", "7", "6.6.6.0/24", "6.6.7.0/24", true),
				),
			},
		},
	})
}

func testAccCheckPanosBgpConditionalAdvAdvertiseFilterExists(n string, o *advertise.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		fw := testAccProvider.Meta().(*pango.Firewall)
		vr, ca, name := parseBgpConditionalAdvAdvertiseFilterId(rs.Primary.ID)
		v, err := fw.Network.BgpConAdvAdvertiseFilter.Get(vr, ca, name)
		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosBgpConditionalAdvAdvertiseFilterAttributes(o *advertise.Entry, name, apr, cr, ecr, med, ap, nh string, en bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %q, not %q", o.Name, name)
		}

		if o.AsPathRegex != apr {
			return fmt.Errorf("AS path regex is %q, not %q", o.AsPathRegex, apr)
		}

		if o.CommunityRegex != cr {
			return fmt.Errorf("Community regex is %q, not %q", o.CommunityRegex, cr)
		}

		if o.ExtendedCommunityRegex != ecr {
			return fmt.Errorf("Extended community regex is %q, not %q", o.ExtendedCommunityRegex, ecr)
		}

		if o.Med != med {
			return fmt.Errorf("MED is %s, not %s", o.Med, med)
		}

		if len(o.AddressPrefix) != 1 || o.AddressPrefix[0] != ap {
			return fmt.Errorf("Address prefixes is %#v, not [%s]", o.AddressPrefix, ap)
		}

		if len(o.NextHop) != 1 || o.NextHop[0] != nh {
			return fmt.Errorf("NextHop is %#v, not [%s]", o.NextHop, nh)
		}

		if o.Enable != en {
			return fmt.Errorf("Enable is %t, not %t", o.Enable, en)
		}

		return nil
	}
}

func testAccPanosBgpConditionalAdvAdvertiseFilterDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_conditional_adv_advertise_filter" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, ca, name := parseBgpConditionalAdvAdvertiseFilterId(rs.Primary.ID)
			_, err := fw.Network.BgpConAdvAdvertiseFilter.Get(vr, ca, name)
			if err == nil {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccBgpConditionalAdvAdvertiseFilterConfig(vr, ca, name, apr, cr, ecr, med, ap, nh string, en bool) string {
	return fmt.Sprintf(`
resource "panos_virtual_router" "vr" {
    name = %q
}

resource "panos_bgp" "x" {
    virtual_router = panos_virtual_router.vr.name
    router_id = "5.5.5.5"
    as_number = "55"
    enable = false
}

resource "panos_bgp_conditional_adv" "ca" {
    virtual_router = panos_bgp.x.virtual_router
    name = %q
    enable = false
}

resource "panos_bgp_conditional_adv_advertise_filter" "test" {
    virtual_router = panos_bgp.x.virtual_router
    bgp_conditional_adv = panos_bgp_conditional_adv.ca.name
    name = %q
    as_path_regex = %q
    community_regex = %q
    extended_community_regex = %q
    med = %q
    address_prefixes = [%q]
    next_hops = [%q]
    enable = %t
}

resource "panos_bgp_conditional_adv_non_exist_filter" "x" {
    virtual_router = panos_bgp.x.virtual_router
    bgp_conditional_adv = panos_bgp_conditional_adv.ca.name
    name = "nef"
    community_regex = "*bar*"
    address_prefixes = ["7.8.9.0/24"]
}
`, vr, ca, name, apr, cr, ecr, med, ap, nh, en)
}
