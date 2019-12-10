package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/netw/routing/protocol/bgp/imp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosBgpImportRuleGroup_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o1, o2, o3, o4, o5 imp.Entry
	vr := fmt.Sprintf("tf%s", acctest.RandString(6))
	m1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m3 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m4 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m5 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m6 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m7 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m8 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m9 := fmt.Sprintf("tf%s", acctest.RandString(6))
	m10 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosBgpImportRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBgpImportRuleGroupConfig(vr, m1, m2, m3, m4, m5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpImportRuleGroupExists("panos_bgp_import_rule_group.top", "panos_bgp_import_rule_group.mid", "panos_bgp_import_rule_group.bot", &o1, &o2, &o3, &o4, &o5),
					testAccCheckPanosBgpImportRuleGroupAttributes(&o1, &o2, &o3, &o4, &o5, m1, m2, m3, m4, m5),
					testAccCheckPanosBgpImportRuleGroupOrdering(vr),
				),
			},
			{
				Config: testAccBgpImportRuleGroupConfig(vr, m6, m7, m8, m9, m10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosBgpImportRuleGroupExists("panos_bgp_import_rule_group.top", "panos_bgp_import_rule_group.mid", "panos_bgp_import_rule_group.bot", &o1, &o2, &o3, &o4, &o5),
					testAccCheckPanosBgpImportRuleGroupAttributes(&o1, &o2, &o3, &o4, &o5, m6, m7, m8, m9, m10),
					testAccCheckPanosBgpImportRuleGroupOrdering(vr),
				),
			},
		},
	})
}

func testAccCheckPanosBgpImportRuleGroupExists(top, mid, bot string, o1, o2, o3, o4, o5 *imp.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var vr string
		var err error
		fw := testAccProvider.Meta().(*pango.Firewall)

		// Top two.
		rTop, ok := s.RootModule().Resources[top]
		if !ok {
			return fmt.Errorf("Resource not found: %s", top)
		}
		if rTop.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}
		vr, _, _, topList := parseBgpImportRuleGroupId(rTop.Primary.ID)
		if len(topList) != 2 {
			return fmt.Errorf("top is not len 2")
		}
		v1, err := fw.Network.BgpImport.Get(vr, topList[0])
		if err != nil {
			return fmt.Errorf("Failed to get top0: %s", err)
		}
		*o1 = v1
		v2, err := fw.Network.BgpImport.Get(vr, topList[1])
		if err != nil {
			return fmt.Errorf("Failed to get top1: %s", err)
		}
		*o2 = v2

		// Middle one.
		rMid, ok := s.RootModule().Resources[mid]
		if !ok {
			return fmt.Errorf("Resource not found: %s", mid)
		}
		if rMid.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}
		vr, _, _, midList := parseBgpImportRuleGroupId(rMid.Primary.ID)
		if len(midList) != 1 {
			return fmt.Errorf("mid is not len 1")
		}
		v3, err := fw.Network.BgpImport.Get(vr, midList[0])
		if err != nil {
			return fmt.Errorf("Failed to get mid: %s", err)
		}
		*o3 = v3

		// Bottom two.
		rBot, ok := s.RootModule().Resources[bot]
		if !ok {
			return fmt.Errorf("Resource not found: %s", bot)
		}
		if rBot.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}
		vr, _, _, botList := parseBgpImportRuleGroupId(rBot.Primary.ID)
		if len(botList) != 2 {
			return fmt.Errorf("bot is not len 2")
		}
		v4, err := fw.Network.BgpImport.Get(vr, botList[0])
		if err != nil {
			return fmt.Errorf("Failed to get bot: %s", err)
		}
		*o4 = v4
		v5, err := fw.Network.BgpImport.Get(vr, botList[1])
		if err != nil {
			return fmt.Errorf("Failed to get bot1: %s", err)
		}
		*o5 = v5

		return nil
	}
}

func testAccCheckPanosBgpImportRuleGroupAttributes(o1, o2, o3, o4, o5 *imp.Entry, m1, m2, m3, m4, m5 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o1.Name != "mary" {
			return fmt.Errorf("Name 1 is %q, expected \"mary\"", o1.Name)
		} else if o1.MatchAsPathRegex != m1 {
			return fmt.Errorf("AS path regex 1 is %q, expected %q", o1.MatchAsPathRegex, m1)
		}

		if o2.Name != "had" {
			return fmt.Errorf("Name 2 is %q, expected \"had\"", o2.Name)
		} else if o2.MatchAsPathRegex != m2 {
			return fmt.Errorf("AS path regex 2 is %q, expected %q", o2.MatchAsPathRegex, m2)
		}

		if o3.Name != "a" {
			return fmt.Errorf("Name 3 is %q, expected \"a\"", o3.Name)
		} else if o3.MatchAsPathRegex != m3 {
			return fmt.Errorf("AS path regex 3 is %q, expected %q", o3.MatchAsPathRegex, m3)
		}

		if o4.Name != "little" {
			return fmt.Errorf("Name 4 is %q, expected \"little\"", o4.Name)
		} else if o4.MatchAsPathRegex != m4 {
			return fmt.Errorf("AS path regex 4 is %q, expected %q", o4.MatchAsPathRegex, m4)
		}

		if o5.Name != "lamb" {
			return fmt.Errorf("Name 1 is %q, expected \"lamb\"", o5.Name)
		} else if o5.MatchAsPathRegex != m5 {
			return fmt.Errorf("AS path regex 5 is %q, expected %q", o5.MatchAsPathRegex, m5)
		}

		return nil
	}
}

func testAccCheckPanosBgpImportRuleGroupOrdering(vr string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fw := testAccProvider.Meta().(*pango.Firewall)
		p3i := -1

		list, err := fw.Network.BgpImport.GetList(vr)
		if err != nil {
			return fmt.Errorf("Failed GetList in ordering check: %s", err)
		}

		for i, v := range list {
			if v == "a" {
				p3i = i
				break
			}
		}

		stl := len(list) - 2
		if len(list) < 5 {
			return fmt.Errorf("Ordering expected at least 5 policies, not %d", len(list))
		} else if list[0] != "mary" {
			return fmt.Errorf("First policy is %q not \"mary\"", list[0])
		} else if list[1] != "had" {
			return fmt.Errorf("Second policy is %q not \"mary\"", list[1])
		} else if p3i == -1 || p3i >= stl {
			return fmt.Errorf("Middle policy is improperly placed: %d vs %d (stl)", p3i, stl)
		} else if list[stl] != "little" {
			return fmt.Errorf("Second to last policy is %q not \"mary\"", list[stl])
		} else if list[stl+1] != "lamb" {
			return fmt.Errorf("Last policy is %q not \"mary\"", list[stl+1])
		}

		return nil
	}
}

func testAccPanosBgpImportRuleGroupDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_bgp_import_rule_group" {
			continue
		}

		if rs.Primary.ID != "" {
			vr, _, _, list := parseBgpImportRuleGroupId(rs.Primary.ID)
			for _, rule := range list {
				_, err := fw.Network.BgpImport.Get(vr, rule)
				if err == nil {
					return fmt.Errorf("BGP import rule %q still exists", rule)
				}
			}
		}
	}

	return nil
}

func testAccBgpImportRuleGroupConfig(vr, m1, m2, m3, m4, m5 string) string {
	return fmt.Sprintf(`
data "panos_system_info" "x" {}

resource "panos_virtual_router" "x" {
    name = %q
}

resource "panos_bgp" "x" {
    virtual_router = panos_virtual_router.x.name
    router_id = "1.2.3.4"
    as_number = "42"
    enable = false
}

resource "panos_bgp_import_rule_group" "top" {
    virtual_router = panos_bgp.x.virtual_router
    position_keyword = "top"
    rule {
        name = "mary"
        match_as_path_regex = %q
        action = "deny"
        match_route_table = data.panos_system_info.x.version_major >= 8 ? "unicast" : ""
    }
    rule {
        name = "had"
        match_as_path_regex = %q
        action = "deny"
        match_route_table = data.panos_system_info.x.version_major >= 8 ? "unicast" : ""
    }
}

resource "panos_bgp_import_rule_group" "mid" {
    virtual_router = panos_bgp.x.virtual_router
    position_keyword = "before"
    position_reference = panos_bgp_import_rule_group.bot.rule.0.name
    rule {
        name = "a"
        match_as_path_regex = %q
        action = "deny"
        match_route_table = data.panos_system_info.x.version_major >= 8 ? "unicast" : ""
    }
}

resource "panos_bgp_import_rule_group" "bot" {
    virtual_router = panos_bgp.x.virtual_router
    position_keyword = "bottom"
    rule {
        name = "little"
        match_as_path_regex = %q
        action = "deny"
        match_route_table = data.panos_system_info.x.version_major >= 8 ? "unicast" : ""
    }
    rule {
        name = "lamb"
        match_as_path_regex = %q
        action = "deny"
        match_route_table = data.panos_system_info.x.version_major >= 8 ? "unicast" : ""
    }
}
`, vr, m1, m2, m3, m4, m5)
}
