package panos

import (
	"fmt"
	"testing"

	"github.com/PaloAltoNetworks/pango"
	"github.com/PaloAltoNetworks/pango/poli/pbf"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccPanosPanoramaPbfRuleGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o1, o2, o3 pbf.Entry
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	n1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	n2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	n3 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaPbfRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaPbfRuleGroupConfig(dg, n1, n2, n3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaPbfRuleGroupExists("panos_panorama_pbf_rule_group.top", "panos_panorama_pbf_rule_group.bot", &o1, &o2, &o3),
					testAccCheckPanosPanoramaPbfRuleGroupAttributes(&o1, &o2, &o3, n1, n2, n3),
					testAccCheckPanosPanoramaPbfRuleGroupOrdering(dg, n1, n2, n3),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaPbfRuleGroupExists(top, bot string, o1, o2, o3 *pbf.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var dg, base string
		var err error
		pano := testAccProvider.Meta().(*pango.Panorama)

		// Top one.
		rTop, ok := s.RootModule().Resources[top]
		if !ok {
			return fmt.Errorf("Resource not found: %s", top)
		}
		if rTop.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}
		dg, base, _, _, topList := parsePanoramaPbfRuleGroupId(rTop.Primary.ID)
		if len(topList) != 1 {
			return fmt.Errorf("top is not len 1")
		}
		v1, err := pano.Policies.PolicyBasedForwarding.Get(dg, base, topList[0])
		if err != nil {
			return fmt.Errorf("Failed to get top: %s", err)
		}
		*o1 = v1

		// Bottom two.
		rBot, ok := s.RootModule().Resources[bot]
		if !ok {
			return fmt.Errorf("Resource not found: %s", bot)
		}
		if rBot.Primary.ID == "" {
			return fmt.Errorf("Object label ID is not set")
		}
		dg, base, _, _, botList := parsePanoramaPbfRuleGroupId(rBot.Primary.ID)
		if len(botList) != 2 {
			return fmt.Errorf("bot is not len 2")
		}
		v2, err := pano.Policies.PolicyBasedForwarding.Get(dg, base, botList[0])
		if err != nil {
			return fmt.Errorf("Failed to get bot: %s", err)
		}
		*o2 = v2
		v3, err := pano.Policies.PolicyBasedForwarding.Get(dg, base, botList[1])
		if err != nil {
			return fmt.Errorf("Failed to get bot1: %s", err)
		}
		*o3 = v3

		return nil
	}
}

func testAccCheckPanosPanoramaPbfRuleGroupAttributes(o1, o2, o3 *pbf.Entry, n1, n2, n3 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o1.Name != n1 {
			return fmt.Errorf("1. Name is %q, not %q", o1.Name, n1)
		} else if o1.Description != "wu" {
			return fmt.Errorf("1. Description is %q, not 'wu'q", o1.Description)
		} else if o1.FromType != pbf.FromTypeZone {
			return fmt.Errorf("1. FromType is %q, not %q", o1.FromType, pbf.FromTypeZone)
		} else if len(o1.FromValues) != 2 {
			return fmt.Errorf("1. FromValues is len %d, not 2", len(o1.FromValues))
		} else if len(o1.SourceAddresses) != 1 || o1.SourceAddresses[0] != "10.50.50.50" {
			return fmt.Errorf("1. SourceAddresses is %#v, not [10.50.50.50]", o1.SourceAddresses)
		} else if len(o1.SourceUsers) != 1 || o1.SourceUsers[0] != "any" {
			return fmt.Errorf("1. SourceUsers is %#v, not [any]", o1.SourceUsers)
		} else if !o1.NegateSource {
			return fmt.Errorf("1. NegateSource is not true")
		} else if len(o1.DestinationAddresses) != 1 || o1.DestinationAddresses[0] != "10.80.80.80" {
			return fmt.Errorf("1. DestinationAddresses is %#v, not [10.80.80.80]", o1.DestinationAddresses)
		} else if len(o1.Applications) != 1 || o1.Applications[0] != "any" {
			return fmt.Errorf("1. Applications is %#v, not [any]", o1.Applications)
		} else if len(o1.Services) != 1 || o1.Services[0] != "application-default" {
			return fmt.Errorf("1. Services is %#v, not [application-default]", o1.Services)
		} else if o1.Action != "discard" {
			return fmt.Errorf("1. Action is %s, not 'discard'", o1.Action)
		}

		if o2.Name != n2 {
			return fmt.Errorf("2. Name is %q, not %q", o2.Name, n2)
		} else if o2.Description != "tang" {
			return fmt.Errorf("2. Description 2 is %q, not 'tang'", o2.Description)
			/*
			   } else if len(o2.Tags) != 2 || o2.Tags[0] != "tagx" || o2.Tags[1] != "tagy" {
			       return fmt.Errorf("2. Tags is %#v, not [tagx, tagy]", o2.Tags)
			*/
		} else if o2.FromType != pbf.FromTypeZone {
			return fmt.Errorf("2. FromType is %q, not %q", o2.FromType, pbf.FromTypeZone)
		} else if len(o2.FromValues) != 1 {
			return fmt.Errorf("2. FromValues is len %d, not 1", len(o2.FromValues))
		} else if len(o2.SourceAddresses) != 1 || o2.SourceAddresses[0] != "10.60.60.60" {
			return fmt.Errorf("2. SourceAddresses is %#v, not [10.60.60.60]", o2.SourceAddresses)
		} else if len(o2.SourceUsers) != 1 || o2.SourceUsers[0] != "any" {
			return fmt.Errorf("2. SourceUsers is %#v, not [any]", o2.SourceUsers)
		} else if len(o2.DestinationAddresses) != 1 || o2.DestinationAddresses[0] != "10.90.90.90" {
			return fmt.Errorf("2. DestinationAddresses is %#v, not [10.90.90.90]", o2.DestinationAddresses)
		} else if !o2.NegateDestination {
			return fmt.Errorf("2. NegateDestination is not true")
		} else if len(o2.Applications) != 1 || o2.Applications[0] != "any" {
			return fmt.Errorf("2. Applications is %#v, not [any]", o2.Applications)
		} else if len(o2.Services) != 1 || o2.Services[0] != "service-http" {
			return fmt.Errorf("2. Services is %#v, not [service-http]", o2.Services)
		} else if o2.Action != "no-pbf" {
			return fmt.Errorf("2. Action is %s, not 'no-pbf'", o2.Action)
		}

		if o3.Name != n3 {
			return fmt.Errorf("3. Name is %q, not %q", o3.Name, n3)
		} else if o3.Description != "clan" {
			return fmt.Errorf("3. Description is %q, not 'clan'", o3.Description)
			/*
			   } else if len(o3.Tags) != 2 || o3.Tags[0] != "tagy" || o3.Tags[1] != "tagx" {
			       return fmt.Errorf("3. Tags is %#v, not [tagy, tagx]", o3.Tags)
			*/
		} else if o3.FromType != pbf.FromTypeInterface {
			return fmt.Errorf("3. FromType is %s, not %s", o3.FromType, pbf.FromTypeInterface)
		} else if len(o3.FromValues) != 1 || o3.FromValues[0] != "ethernet1/2" {
			return fmt.Errorf("3. FromValues is %#v, not [ethernet1/2]", o3.FromValues)
		} else if len(o3.SourceAddresses) != 1 || o3.SourceAddresses[0] != "10.70.70.70" {
			return fmt.Errorf("3. SourceAddresses is %#v, not [10.70.70.70]", o3.SourceAddresses)
		} else if len(o3.SourceUsers) != 1 || o3.SourceUsers[0] != "any" {
			return fmt.Errorf("3. SourceUsers is %#v, not [any]", o3.SourceUsers)
		} else if len(o3.DestinationAddresses) != 1 || o3.DestinationAddresses[0] != "10.100.100.100" {
			return fmt.Errorf("3. DestinationAddresses is %#v, not [10.100.100.100]", o3.DestinationAddresses)
		} else if len(o3.Applications) != 1 || o3.Applications[0] != "any" {
			return fmt.Errorf("3. Applications is %#v, not [any]", o3.Applications)
		} else if len(o3.Services) != 1 || o3.Services[0] != "service-https" {
			return fmt.Errorf("3. Services is %#v, not [service-https]", o3.Services)
		} else if o3.Action != pbf.ActionForward {
			return fmt.Errorf("3. Action is %s, not %s", o3.Action, pbf.ActionForward)
		} else if o3.ForwardEgressInterface != "ethernet1/2" {
			return fmt.Errorf("3. ForwardEgressInterface is %s, not ethernet1/2", o3.ForwardEgressInterface)
		} else if o3.ForwardMonitorProfile != "my-monitor-profile" {
			return fmt.Errorf("3. ForwardMonitorProfile is %s, not my-monitor-profile", o3.ForwardMonitorProfile)
		} else if !o3.EnableEnforceSymmetricReturn {
			return fmt.Errorf("3. EnableEnforceSymmetricReturn is not set")
		} else if len(o3.SymmetricReturnAddresses) != 2 || o3.SymmetricReturnAddresses[0] != "10.20.50.90" || o3.SymmetricReturnAddresses[1] != "5.4.3.2" {
			return fmt.Errorf("3. SymmetricReturnAddresses is %#v, not [10.20.50.90, 5.4.3.2]", o3.SymmetricReturnAddresses)
		}

		return nil
	}
}

func testAccCheckPanosPanoramaPbfRuleGroupOrdering(dg, n1, n2, n3 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pano := testAccProvider.Meta().(*pango.Panorama)

		list, err := pano.Policies.PolicyBasedForwarding.GetList(dg, util.PreRulebase)
		if err != nil {
			return fmt.Errorf("Failed GetList in ordering check: %s", err)
		}

		for i, v := range list {
			if v == n1 {
				if i+1 >= len(list) {
					return fmt.Errorf("No rules after n1 %q", n1)
				}
				if list[i+1] != n2 {
					return fmt.Errorf("Rule after n1 (%s) is %q, not %q", n1, list[i+1], n2)
				}
				if i+2 >= len(list) {
					return fmt.Errorf("No rules after n2 %q", n2)
				}
				if list[i+2] != n3 {
					return fmt.Errorf("Rule after n2 (%s) is %q, not %q", n2, list[i+2], n3)
				}
				return nil
			}
		}

		return fmt.Errorf("Rule n1 (%s) not found", n1)
	}
}

func testAccPanosPanoramaPbfRuleGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_pbf_rule_group" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, base, _, _, list := parsePanoramaPbfRuleGroupId(rs.Primary.ID)
			for _, rule := range list {
				_, err := pano.Policies.PolicyBasedForwarding.Get(dg, base, rule)
				if err == nil {
					return fmt.Errorf("PBF rule %q still exists", rule)
				}
			}
		}
	}

	return nil
}

func testAccPanoramaPbfRuleGroupConfig(dg, n1, n2, n3 string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
    description = "pbf rule group test"
}

/*
resource "panos_panorama_administrative_tag" "x" {
    device_group = panos_panorama_device_group.x.name
    name = "tagx"
    color = "color5"
}

resource "panos_panorama_administrative_tag" "y" {
    device_group = panos_panorama_device_group.x.name
    name = "tagy"
    color = "color6"
}
*/

resource "panos_panorama_pbf_rule_group" "top" {
    device_group = panos_panorama_device_group.x.name
    position_keyword = "directly before"
    position_reference = panos_panorama_pbf_rule_group.bot.rule.0.name
    rule {
        name = %q
        description = "wu"
        source {
            zones = ["zone1", "zone2"]
            addresses = ["10.50.50.50"]
            users = ["any"]
            negate = true
        }
        destination {
            addresses = ["10.80.80.80"]
            applications = ["any"]
            services = ["application-default"]
        }
        forwarding {
            action = "discard"
        }
    }
}

resource "panos_panorama_pbf_rule_group" "bot" {
    device_group = panos_panorama_device_group.x.name
    rule {
        name = %q
        description = "tang"
        //tags = [panos_panorama_administrative_tag.x.name, panos_panorama_administrative_tag.y.name]
        source {
            zones = ["zone1"]
            addresses = ["10.60.60.60"]
            users = ["any"]
        }
        destination {
            addresses = ["10.90.90.90"]
            applications = ["any"]
            services = ["service-http"]
            negate = true
        }
        forwarding {
            action = "no-pbf"
        }
    }
    rule {
        name = %q
        description = "clan"
        //tags = [panos_panorama_administrative_tag.y.name, panos_panorama_administrative_tag.x.name]
        source {
            interfaces = ["ethernet1/2"]
            addresses = ["10.70.70.70"]
            users = ["any"]
        }
        destination {
            addresses = ["10.100.100.100"]
            applications = ["any"]
            services = ["service-https"]
        }
        forwarding {
            egress_interface = "ethernet1/2"
            monitor {
                profile = "my-monitor-profile"
            }
            symmetric_return {
                enable = true
                addresses = ["10.20.50.90", "5.4.3.2"]
            }
        }
    }
}
`, dg, n1, n2, n3)
}
