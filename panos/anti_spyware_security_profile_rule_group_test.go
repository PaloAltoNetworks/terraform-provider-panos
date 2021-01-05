package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/objs/profile/security/spyware/rule"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Data source listing tests.
func TestAccPanosDsFirewallAntiSpywareSecurityProfileRulesList(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileRuleGroupConfig(prof, name),
				Check:  checkDataSourceListing("panos_anti_spyware_security_profile_rules"),
			},
		},
	})
}

func TestAccPanosDsPanoramaAntiSpywareSecurityProfileRulesList(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileRuleGroupConfig(prof, name),
				Check:  checkDataSourceListing("panos_anti_spyware_security_profile_rules"),
			},
		},
	})
}

// Data source tests.
func TestAccPanosDsFirewallAntiSpywareSecurityProfileRule_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileRuleGroupConfig(prof, name),
				Check: checkDataSource("panos_anti_spyware_security_profile_rule", []string{
					"rule.0.name", "rule.0.threat_name", "rule.0.category",
					"rule.0.packet_capture", "rule.0.action",
				}),
			},
		},
	})
}

func TestAccPanosDsPanoramaAntiSpywareSecurityProfileRule_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDsAntiSpywareSecurityProfileRuleGroupConfig(prof, name),
				Check: checkDataSource("panos_anti_spyware_security_profile_rule", []string{
					"rule.0.name", "rule.0.threat_name", "rule.0.category",
					"rule.0.packet_capture", "rule.0.action",
				}),
			},
		},
	})
}

func testAccDsAntiSpywareSecurityProfileRuleGroupConfig(prof, name string) string {
	return fmt.Sprintf(`
data "panos_anti_spyware_security_profile_rules" "test" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
}

data "panos_anti_spyware_security_profile_rule" "test" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
    name = panos_anti_spyware_security_profile_rule_group.x.rule.0.name
}

resource "panos_anti_spyware_security_profile" "x" {
    name = %q
    description = "anti_spyware sec prof acctest"
    sinkhole_ipv4_address = "pan-sinkhole-default-ip"
    sinkhole_ipv6_address = "::1"
    botnet_list {
        name = "default-paloalto-dns"
        action = "sinkhole"
        packet_capture = "disable"
    }
    botnet_list {
        name = "default-paloalto-cloud"
        action = "allow"
        packet_capture = "disable"
    }
}

resource "panos_anti_spyware_security_profile_rule_group" "x" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
    rule {
        name = %q
        threat_name = "any"
        category = "any"
        action = "default"
        packet_capture = "disable"
        severities = ["any"]
    }
}
`, prof, name)
}

// Resource tests.
func TestAccPanosFirewallAntiSpywareSecurityProfileRuleGroup_basic(t *testing.T) {
	if !testAccIsFirewall {
		t.Skip(SkipFirewallAccTest)
	}

	var o rule.Entry
	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntiSpywareSecurityProfileRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiSpywareSecurityProfileRuleGroupConfig(prof, name, "any", "any", "reset-client", "single-packet", "critical", "low"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupExists("panos_anti_spyware_security_profile_rule_group.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupAttributes(&o, name, "any", "any", "reset-client", "single-packet", "critical", "low"),
				),
			},
			{
				Config: testAccAntiSpywareSecurityProfileRuleGroupConfig(prof, name, "foo", "adware", "allow", "disable", "medium", "high"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupExists("panos_anti_spyware_security_profile_rule_group.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupAttributes(&o, name, "foo", "adware", "allow", "disable", "medium", "high"),
				),
			},
		},
	})
}

func TestAccPanosPanoramaAntiSpywareSecurityProfileRuleGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o rule.Entry
	prof := fmt.Sprintf("tf%s", acctest.RandString(6))
	name := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosAntiSpywareSecurityProfileRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAntiSpywareSecurityProfileRuleGroupConfig(prof, name, "any", "any", "reset-client", "single-packet", "critical", "low"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupExists("panos_anti_spyware_security_profile_rule_group.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupAttributes(&o, name, "any", "any", "reset-client", "single-packet", "critical", "low"),
				),
			},
			{
				Config: testAccAntiSpywareSecurityProfileRuleGroupConfig(prof, name, "foo", "adware", "allow", "disable", "medium", "high"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupExists("panos_anti_spyware_security_profile_rule_group.test", &o),
					testAccCheckPanosAntiSpywareSecurityProfileRuleGroupAttributes(&o, name, "foo", "adware", "allow", "disable", "medium", "high"),
				),
			},
		},
	})
}

func testAccCheckPanosAntiSpywareSecurityProfileRuleGroupExists(n string, o *rule.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}

		var err error
		var v rule.Entry

		switch con := testAccProvider.Meta().(type) {
		case *pango.Firewall:
			vsys, prof, _, _, names := parseAntiSpywareSecurityProfileRuleGroupId(rs.Primary.ID)
			if len(names) != 1 {
				return fmt.Errorf("names is not len 1: %#v", names)
			}
			v, err = con.Objects.AntiSpywareRule.Get(vsys, prof, names[0])
		case *pango.Panorama:
			dg, prof, _, _, names := parseAntiSpywareSecurityProfileRuleGroupId(rs.Primary.ID)
			if len(names) != 1 {
				return fmt.Errorf("names is not len 1: %#v", names)
			}
			v, err = con.Objects.AntiSpywareRule.Get(dg, prof, names[0])
		}

		if err != nil {
			return fmt.Errorf("Error in get: %s", err)
		}

		*o = v

		return nil
	}
}

func testAccCheckPanosAntiSpywareSecurityProfileRuleGroupAttributes(o *rule.Entry, name, tn, cat, action, pc, sev1, sev2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o.Name != name {
			return fmt.Errorf("Name is %s, expected %s", o.Name, name)
		}

		if o.ThreatName != tn {
			return fmt.Errorf("Threat name is %q, expected %s", o.ThreatName, tn)
		}

		if o.Category != cat {
			return fmt.Errorf("Category is %q, expected %q", o.Category, cat)
		}

		if o.Action != action {
			return fmt.Errorf("Action is %q, expected %q", o.Action, action)
		}

		if o.PacketCapture != pc {
			return fmt.Errorf("Packet capture is %q, expected %q", o.PacketCapture, pc)
		}

		if len(o.Severities) != 2 {
			return fmt.Errorf("Severities is not 2 long: %#v", o.Severities)
		}

		if o.Severities[0] != sev1 {
			return fmt.Errorf("sev1 is %q, not %q", o.Severities[0], sev1)
		}

		if o.Severities[1] != sev2 {
			return fmt.Errorf("sev2 is %q, not %q", o.Severities[1], sev2)
		}

		return nil
	}
}

func testAccPanosAntiSpywareSecurityProfileRuleGroupDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_anti_spyware_security_profile_rule_group" {
			continue
		}

		if rs.Primary.ID != "" {
			var list []string
			var err error

			switch con := testAccProvider.Meta().(type) {
			case *pango.Firewall:
				vsys, prof, _, _, _ := parseAntiSpywareSecurityProfileRuleGroupId(rs.Primary.ID)
				list, err = con.Objects.AntiSpywareRule.GetList(vsys, prof)
			case *pango.Panorama:
				dg, prof, _, _, _ := parseAntiSpywareSecurityProfileRuleGroupId(rs.Primary.ID)
				list, err = con.Objects.AntiSpywareRule.GetList(dg, prof)
			}
			if err == nil && len(list) > 0 {
				return fmt.Errorf("Object %q still exists", rs.Primary.ID)
			}
		}
		return nil
	}

	return nil
}

func testAccAntiSpywareSecurityProfileRuleGroupConfig(prof, name, tn, cat, action, pc, sev1, sev2 string) string {
	return fmt.Sprintf(`
resource "panos_anti_spyware_security_profile" "x" {
    name = %q
    description = "anti spyware sec prof rule group acctest"
    sinkhole_ipv4_address = "pan-sinkhole-default-ip"
    sinkhole_ipv6_address = "::1"
    botnet_list {
        name = "default-paloalto-dns"
        action = "sinkhole"
        packet_capture = "disable"
    }
    botnet_list {
        name = "default-paloalto-cloud"
        action = "allow"
        packet_capture = "disable"
    }
}

resource "panos_anti_spyware_security_profile_rule_group" "test" {
    anti_spyware_security_profile = panos_anti_spyware_security_profile.x.name
    rule {
        name = %q
        threat_name = %q
        category = %q
        action = %q
        packet_capture = %q
        severities = [%q, %q]
    }
}
`, prof, name, tn, cat, action, pc, sev1, sev2)
}
