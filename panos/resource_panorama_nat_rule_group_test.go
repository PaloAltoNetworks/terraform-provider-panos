package panos

import (
	"fmt"
	"testing"

	"github.com/fpluchorg/pango"
	"github.com/fpluchorg/pango/poli/nat"
	"github.com/fpluchorg/pango/util"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccPanosPanoramaNatRuleGroup_basic(t *testing.T) {
	if !testAccIsPanorama {
		t.Skip(SkipPanoramaAccTest)
	}

	var o1, o2, o3 nat.Entry
	dg := fmt.Sprintf("tf%s", acctest.RandString(6))
	n1 := fmt.Sprintf("tf%s", acctest.RandString(6))
	n2 := fmt.Sprintf("tf%s", acctest.RandString(6))
	n3 := fmt.Sprintf("tf%s", acctest.RandString(6))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccPanosPanoramaNatRuleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPanoramaNatRuleGroupConfig(dg, n1, n2, n3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPanosPanoramaNatRuleGroupExists("panos_panorama_nat_rule_group.top", "panos_panorama_nat_rule_group.bot", &o1, &o2, &o3),
					testAccCheckPanosPanoramaNatRuleGroupAttributes(&o1, &o2, &o3, n1, n2, n3),
					testAccCheckPanosPanoramaNatRuleGroupOrdering(dg, n1, n2, n3),
				),
			},
		},
	})
}

func testAccCheckPanosPanoramaNatRuleGroupExists(top, bot string, o1, o2, o3 *nat.Entry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		var dg, rb string
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
		dg, rb, _, _, topList := parsePanoramaNatRuleGroupId(rTop.Primary.ID)
		if len(topList) != 1 {
			return fmt.Errorf("top is not len 1")
		}
		v1, err := pano.Policies.Nat.Get(dg, rb, topList[0])
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
		dg, rb, _, _, botList := parsePanoramaNatRuleGroupId(rBot.Primary.ID)
		if len(botList) != 2 {
			return fmt.Errorf("bot is not len 2")
		}
		v2, err := pano.Policies.Nat.Get(dg, rb, botList[0])
		if err != nil {
			return fmt.Errorf("Failed to get bot: %s", err)
		}
		*o2 = v2
		v3, err := pano.Policies.Nat.Get(dg, rb, botList[1])
		if err != nil {
			return fmt.Errorf("Failed to get bot1: %s", err)
		}
		*o3 = v3

		return nil
	}
}

func testAccCheckPanosPanoramaNatRuleGroupAttributes(o1, o2, o3 *nat.Entry, n1, n2, n3 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if o1.Name != n1 {
			return fmt.Errorf("1. Name is %q, not %q", o1.Name, n1)
		} else if o1.Description != "wu" {
			return fmt.Errorf("1. Description is %q, not 'wu'q", o1.Description)
		} else if len(o1.SourceZones) != 1 || o1.SourceZones[0] != "z1" {
			return fmt.Errorf("1. Source zones is %#v, not ['z1']", o1.SourceZones)
		} else if o1.DestinationZone != "z1" {
			return fmt.Errorf("1. Destination zone is %q, not 'z1'", o1.DestinationZone)
		} else if o1.ToInterface != "any" {
			return fmt.Errorf("1. Destination interface is %q, not 'any'", o1.ToInterface)
		} else if len(o1.SourceAddresses) != 1 || o1.SourceAddresses[0] != "any" {
			return fmt.Errorf("1. Source addresses is %#v, not ['any']", o1.SourceAddresses)
		} else if len(o1.DestinationAddresses) != 1 || o1.DestinationAddresses[0] != "any" {
			return fmt.Errorf("1. Destination addresses is %#v, not ['any']", o1.DestinationAddresses)
		} else if o1.SatType != nat.DynamicIpAndPort {
			return fmt.Errorf("1. SatType is %s, not %s", o1.SatType, nat.DynamicIpAndPort)
		} else if o1.SatAddressType != nat.InterfaceAddress {
			return fmt.Errorf("1. SatAddressType is %s, not %s", o1.SatAddressType, nat.InterfaceAddress)
		} else if o1.SatInterface != "ethernet1/6" {
			return fmt.Errorf("1. SatInterface is %s, not 'ethernet1/6'", o1.SatInterface)
		} else if o1.DatType != nat.DatTypeStatic {
			return fmt.Errorf("1. DatType is %s, not %s", o1.DatType, nat.DatTypeStatic)
		} else if o1.DatAddress != "10.1.1.1" {
			return fmt.Errorf("1. DatAddress is %s, not '10.1.1.1'", o1.DatAddress)
		} else if o1.DatPort != 1234 {
			return fmt.Errorf("1. DatPort is %d, not 1234", o1.DatPort)
		}

		if o2.Name != n2 {
			return fmt.Errorf("2. Name is %q, not %q", o2.Name, n2)
		} else if o2.Description != "tang" {
			return fmt.Errorf("2. Description 2 is %q, not 'tang'", o2.Description)
		} else if len(o2.SourceZones) != 1 || o2.SourceZones[0] != "z2" {
			return fmt.Errorf("2. Source zones is %#v, not ['z2']", o2.SourceZones)
		} else if o2.DestinationZone != "z3" {
			return fmt.Errorf("2. Destination zone is %q, not 'z3'", o2.DestinationZone)
		} else if o2.ToInterface != "any" {
			return fmt.Errorf("2. Destination interface is %q, not 'any'", o2.ToInterface)
		} else if len(o2.SourceAddresses) != 1 || o2.SourceAddresses[0] != "any" {
			return fmt.Errorf("2. Source addresses is %#v, not ['any']", o2.SourceAddresses)
		} else if len(o2.DestinationAddresses) != 1 || o2.DestinationAddresses[0] != "any" {
			return fmt.Errorf("2. Destination addresses is %#v, not ['any']", o2.DestinationAddresses)
		} else if o2.SatType != nat.None {
			return fmt.Errorf("2. SatType is %s, not %s", o2.SatType, nat.None)
		} else if o2.DatType != nat.DatTypeStatic {
			return fmt.Errorf("2. DatType is %s, not %s", o2.DatType, nat.DatTypeStatic)
		} else if o2.DatAddress != "10.2.3.1" {
			return fmt.Errorf("2. DatAddress is %s, not '10.2.3.1'", o2.DatAddress)
		} else if o2.DatPort != 5678 {
			return fmt.Errorf("2. DatPort is %d, not 5678", o2.DatPort)
		}

		if o3.Name != n3 {
			return fmt.Errorf("3. Name is %q, not %q", o3.Name, n3)
		} else if o3.Description != "clan" {
			return fmt.Errorf("3. Description is %q, not 'clan'", o3.Description)
		} else if len(o3.SourceZones) != 1 || o3.SourceZones[0] != "z3" {
			return fmt.Errorf("3. Source zones is %#v, not ['z3']", o3.SourceZones)
		} else if o3.DestinationZone != "z2" {
			return fmt.Errorf("3. Destination zone is %q, not 'z2'", o3.DestinationZone)
		} else if o3.ToInterface != "any" {
			return fmt.Errorf("3. Destination interface is %q, not 'any'", o3.ToInterface)
		} else if len(o3.SourceAddresses) != 1 || o3.SourceAddresses[0] != "any" {
			return fmt.Errorf("3. Source addresses is %#v, not ['any']", o3.SourceAddresses)
		} else if len(o3.DestinationAddresses) != 1 || o3.DestinationAddresses[0] != "any" {
			return fmt.Errorf("3. Destination addresses is %#v, not ['any']", o3.DestinationAddresses)
		} else if o3.SatType != nat.StaticIp {
			return fmt.Errorf("3. SatType is %s, not %s", o3.SatType, nat.StaticIp)
		} else if o3.SatStaticTranslatedAddress != "192.168.1.5" {
			return fmt.Errorf("3. SatStaticTranslatedAddress is %s, not '192.168.1.5'", o3.SatStaticTranslatedAddress)
		} else if o3.SatStaticBiDirectional != true {
			return fmt.Errorf("3. SatStaticBiDirectional is %t, not true", o3.SatStaticBiDirectional)
		} else if o3.DatType != "" {
			return fmt.Errorf("3. DatType is %s, not ''", o3.DatType)
		}

		return nil
	}
}

func testAccCheckPanosPanoramaNatRuleGroupOrdering(dg, n1, n2, n3 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		pano := testAccProvider.Meta().(*pango.Panorama)

		list, err := pano.Policies.Nat.GetList(dg, util.PreRulebase)
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

func testAccPanosPanoramaNatRuleGroupDestroy(s *terraform.State) error {
	pano := testAccProvider.Meta().(*pango.Panorama)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "panos_panorama_nat_rule_group" {
			continue
		}

		if rs.Primary.ID != "" {
			dg, rb, _, _, list := parsePanoramaNatRuleGroupId(rs.Primary.ID)
			for _, rule := range list {
				_, err := pano.Policies.Nat.Get(dg, rb, rule)
				if err == nil {
					return fmt.Errorf("Nat policy %q still exists", rule)
				}
			}
		}
	}

	return nil
}

func testAccPanoramaNatRuleGroupConfig(dg, n1, n2, n3 string) string {
	return fmt.Sprintf(`
resource "panos_panorama_device_group" "x" {
    name = %q
}

resource "panos_panorama_nat_rule_group" "top" {
    device_group = panos_panorama_device_group.x.name
    position_keyword = "directly before"
    position_reference = panos_panorama_nat_rule_group.bot.rule.0.name
    rule {
        name = %q
        description = "wu"
        original_packet {
            source_zones = ["z1"]
            destination_zone = "z1"
            destination_interface = "any"
            source_addresses = ["any"]
            destination_addresses = ["any"]
        }
        translated_packet {
            source {
                dynamic_ip_and_port {
                    interface_address {
                        interface = "ethernet1/6"
                    }
                }
            }
            destination {
                static_translation {
                    address = "10.1.1.1"
                    port = 1234
                }
            }
        }
    }
}

resource "panos_panorama_nat_rule_group" "bot" {
    device_group = panos_panorama_device_group.x.name
    rule {
        name = %q
        description = "tang"
        original_packet {
            source_zones = ["z2"]
            destination_zone = "z3"
            destination_interface = "any"
            source_addresses = ["any"]
            destination_addresses = ["any"]
        }
        translated_packet {
            source {}
            destination {
                static_translation {
                    address = "10.2.3.1"
                    port = 5678
                }
            }
        }
    }
    rule {
        name = %q
        description = "clan"
        original_packet {
            source_zones = ["z3"]
            destination_zone = "z2"
            destination_interface = "any"
            source_addresses = ["any"]
            destination_addresses = ["any"]
        }
        translated_packet {
            source {
                static_ip {
                    translated_address = "192.168.1.5"
                    bi_directional = true
                }
            }
            destination {}
        }
    }
}
`, dg, n1, n2, n3)
}
