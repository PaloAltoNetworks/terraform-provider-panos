package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/security"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccPanosSecurityRuleGroup_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o1, o2, o3, o4, o5 security.Entry
	d1 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d2 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d3 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d4 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d5 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d6 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d7 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d8 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d9 := fmt.Sprintf("desc %s", acctest.RandString(6))
	d10 := fmt.Sprintf("desc %s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosSecurityRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSecurityRuleGroupConfig(d1, d2, d3, d4, d5),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSecurityRuleGroupExists("panos_security_rule_group.top", "panos_security_rule_group.mid", "panos_security_rule_group.bot", &o1, &o2, &o3, &o4, &o5),
					testAccCheckPanosSecurityRuleGroupAttributes(&o1, &o2, &o3, &o4, &o5, d1, d2, d3, d4, d5),
					testAccCheckPanosSecurityRuleGroupOrdering(),
				),
			},
			{
				Config: testAccSecurityRuleGroupConfig(d6, d7, d8, d9, d10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosSecurityRuleGroupExists("panos_security_rule_group.top", "panos_security_rule_group.mid", "panos_security_rule_group.bot", &o1, &o2, &o3, &o4, &o5),
					testAccCheckPanosSecurityRuleGroupAttributes(&o1, &o2, &o3, &o4, &o5, d6, d7, d8, d9, d10),
					testAccCheckPanosSecurityRuleGroupOrdering(),
				),
			},
		},
	})
}

func testAccCheckPanosSecurityRuleGroupExists(top, mid, bot string, o1, o2, o3, o4, o5 *security.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var vsys string
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
		vsys, _, _, topList := parseSecurityRuleGroupId(rTop.Primary.ID)
		if len(topList) != 2 {
			return fmt.Errorf("top is not len 2")
		}
		v1, err := fw.Policies.Security.Get(vsys, topList[0])
		if err != nil {
			return fmt.Errorf("Failed to get top0: %s", err)
		}
		*o1 = v1
		v2, err := fw.Policies.Security.Get(vsys, topList[1])
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
		vsys, _, _, midList := parseSecurityRuleGroupId(rMid.Primary.ID)
		if len(midList) != 1 {
			return fmt.Errorf("mid is not len 1")
		}
		v3, err := fw.Policies.Security.Get(vsys, midList[0])
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
		vsys, _, _, botList := parseSecurityRuleGroupId(rBot.Primary.ID)
		if len(botList) != 2 {
			return fmt.Errorf("bot is not len 2")
		}
		v4, err := fw.Policies.Security.Get(vsys, botList[0])
		if err != nil {
			return fmt.Errorf("Failed to get bot: %s", err)
		}
		*o4 = v4
		v5, err := fw.Policies.Security.Get(vsys, botList[1])
		if err != nil {
			return fmt.Errorf("Failed to get bot1: %s", err)
		}
		*o5 = v5

		return nil
	}
}

func testAccCheckPanosSecurityRuleGroupAttributes(o1, o2, o3, o4, o5 *security.Entry, d1, d2, d3, d4, d5 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o1.Name != "mary" {
			return fmt.Errorf("Name 1 is %q, expected \"mary\"", o1.Name)
		} else if o1.Description != d1 {
			return fmt.Errorf("Description 1 is %q, expected %q", o1.Description, d1)
		}

		if o2.Name != "had" {
			return fmt.Errorf("Name 2 is %q, expected \"had\"", o2.Name)
		} else if o2.Description != d2 {
			return fmt.Errorf("Description 2 is %q, expected %q", o2.Description, d2)
		}

		if o3.Name != "a" {
			return fmt.Errorf("Name 3 is %q, expected \"a\"", o3.Name)
		} else if o3.Description != d3 {
			return fmt.Errorf("Description 3 is %q, expected %q", o3.Description, d3)
		}

		if o4.Name != "little" {
			return fmt.Errorf("Name 4 is %q, expected \"little\"", o4.Name)
		} else if o4.Description != d4 {
			return fmt.Errorf("Description 4 is %q, expected %q", o4.Description, d4)
		}

		if o5.Name != "lamb" {
			return fmt.Errorf("Name 1 is %q, expected \"lamb\"", o5.Name)
		} else if o5.Description != d5 {
			return fmt.Errorf("Description 5 is %q, expected %q", o5.Description, d5)
		}

		return nil
	}
}

func testAccCheckPanosSecurityRuleGroupOrdering() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fw := testAccProvider.Meta().(*pango.Firewall)
		p3i := -1

		list, err := fw.Policies.Security.GetList("")
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

func testAccPanosSecurityRuleGroupDestroy(s *terraform.State) error {
	fw := testAccProvider.Meta().(*pango.Firewall)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_security_rule_group" {
			continue
		}

		if rs.Primary.ID != "" {
			vsys, _, _, list := parseSecurityRuleGroupId(rs.Primary.ID)
			for _, rule := range list {
				_, err := fw.Policies.Security.Get(vsys, rule)
				if err == nil {
					return fmt.Errorf("Security policy %q still exists", rule)
				}
			}
		}
	}

	return nil
}

func testAccSecurityRuleGroupConfig(d1, d2, d3, d4, d5 string) string {
	return fmt.Sprintf(`
resource "panos_security_rule_group" "top" {
    position_keyword = "top"
    rule {
        name = "mary"
        description = "%s"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
    rule {
        name = "had"
        description = "%s"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
}

resource "panos_security_rule_group" "mid" {
    position_keyword = "before"
    position_reference = panos_security_rule_group.bot.rule.0.name
    rule {
        name = "a"
        description = "%s"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
}

resource "panos_security_rule_group" "bot" {
    position_keyword = "bottom"
    rule {
        name = "little"
        description = "%s"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
    rule {
        name = "lamb"
        description = "%s"
        source_zones = ["any"]
        source_addresses = ["any"]
        source_users = ["any"]
        hip_profiles = ["any"]
        destination_zones = ["any"]
        destination_addresses = ["any"]
        applications = ["any"]
        services = ["application-default"]
        categories = ["any"]
        action = "allow"
    }
}
`, d1, d2, d3, d4, d5)
}
