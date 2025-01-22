package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/policies/rules/security"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type expectServerSecurityRulesOrder struct {
	Location  security.Location
	Prefix    string
	RuleNames []string
}

func ExpectServerSecurityRulesOrder(prefix string, location security.Location, ruleNames []string) *expectServerSecurityRulesOrder {
	return &expectServerSecurityRulesOrder{
		Location:  location,
		Prefix:    prefix,
		RuleNames: ruleNames,
	}
}

func (o *expectServerSecurityRulesOrder) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := security.NewService(sdkClient)

	objects, err := service.List(ctx, o.Location, "get", "", "")
	if err != nil {
		resp.Error = fmt.Errorf("failed to query server for rules: %w", err)
		return
	}

	type ruleWithState struct {
		Idx   int
		State int
	}

	rulesWithIdx := make(map[string]ruleWithState)
	for idx, elt := range o.RuleNames {
		rulesWithIdx[fmt.Sprintf("%s-%s", o.Prefix, elt)] = ruleWithState{
			Idx:   idx,
			State: 0,
		}
	}

	var prevActualIdx = -1
	for actualIdx, elt := range objects {
		if state, ok := rulesWithIdx[elt.Name]; !ok {
			continue
		} else {
			state.State = 1
			rulesWithIdx[elt.Name] = state

			if state.Idx == 0 {
				prevActualIdx = actualIdx
				continue
			} else if prevActualIdx == -1 {
				resp.Error = fmt.Errorf("rules missing from the server")
				return
			} else if actualIdx-prevActualIdx > 1 {
				resp.Error = fmt.Errorf("invalid rules order on the server")
				return
			}
			prevActualIdx = actualIdx
		}
	}

	var missing []string
	for name, elt := range rulesWithIdx {
		if elt.State != 1 {
			missing = append(missing, name)
		}
	}

	if len(missing) > 0 {
		resp.Error = fmt.Errorf("not all rules are present on the server: %s", strings.Join(missing, ", "))
		return
	}
}

type expectServerSecurityRulesCount struct {
	Prefix   string
	Location security.Location
	Count    int
}

func ExpectServerSecurityRulesCount(prefix string, location security.Location, count int) *expectServerSecurityRulesCount {
	return &expectServerSecurityRulesCount{
		Prefix:   prefix,
		Location: location,
		Count:    count,
	}
}

func (o *expectServerSecurityRulesCount) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := security.NewService(sdkClient)

	objects, err := service.List(ctx, o.Location, "get", "", "")
	if err != nil {
		resp.Error = fmt.Errorf("failed to query server for rules: %w", err)
		return
	}

	var count int
	for _, elt := range objects {
		if strings.HasPrefix(elt.Name, o.Prefix) {
			count += 1
		}
	}

	if count != o.Count {
		resp.Error = UnexpectedRulesError
		return
	}
}

const securityPolicyExtendedResource1Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-secgroup-tmpl1", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-secgroup-dg1", var.prefix)
  templates = [ resource.panos_template.template.name ]
}


resource "panos_security_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule1", var.prefix)

    source_zones     = ["any"]
    source_addresses = ["any"]
    source_users     = ["some-user"]
    # source_hips      = ["hip-profile"]
    negate_source    = false

    destination_zone      = "any"
    destination_addresses = ["any"]
    # destination_hips = ["hip-device"]

    services = ["any"]
    applications = ["any"]

    action = "drop"

    profile_setting = {
      profiles = {
        url_filtering      = ["default"]
        # data_filtering     = ["default"]
        file_blocking      = ["basic file blocking"]
        virus              = ["default"]
        spyware            = ["strict"]
        vulnerability      = ["strict"]
        wildfire_analysis  = ["default"]
      }
    }

    qos = {
      marking = {
        ip_dscp = "af11"
      }
    }

    disable_server_response_inspection = true

    icmp_unreachable = true
    # schedule         = "schedule"
    # log_setting = "log-forwarding-test-1"
  }]
}
`

func TestAccSecurityPolicyExtended(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyExtendedResource1Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_zones"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_addresses"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_users"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("some-user"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("negate_source"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_addresses"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_zone"),
						knownvalue.StringExact("any"),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("services"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("applications"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("action"),
						knownvalue.StringExact("drop"),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("url_filtering"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("url_filtering"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("file_blocking"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("basic file blocking"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("virus"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("spyware"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("strict"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("vulnerability"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("strict"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("profile_setting").
							AtMapKey("profiles").
							AtMapKey("wildfire_analysis"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_security_policy.policy",
					// 	tfjsonpath.New("rules").
					// 		AtSliceIndex(0).
					// 		AtMapKey("qos").
					// 		AtMapKey("marking").
					// 		AtMapKey("ip_dscp"),
					// 	knownvalue.StringExact("af11"),
					// ),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("disable_server_response_inspection"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("icmp_unreachable"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const securityPolicyOrderingTmpl = `
variable "rule_names" { type = list(string) }
variable "location" { type = map }

resource "panos_security_policy" "policy" {
  location = var.location

  rules = [
    for index, name in var.rule_names: {
      name = name

      source_zones     = ["any"]
      source_addresses = ["any"]

      destination_zone      = "any"
      destination_addresses = ["any"]

      services = ["any"]
      applications = ["any"]
    }
  ]
}
`

func TestAccSecurityPolicyOrdering(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	rulesInitial := []string{"rule-1", "rule-2", "rule-3", "rule-4", "rule-5"}
	rulesReordered := []string{"rule-2", "rule-1", "rule-3", "rule-4", "rule-5"}

	prefixed := func(name string) string {
		return fmt.Sprintf("%s-%s", prefix, name)
	}

	withPrefix := func(rules []string) []config.Variable {
		var result []config.Variable
		for _, elt := range rules {
			result = append(result, config.StringVariable(prefixed(elt)))
		}

		return result
	}

	device := devicePanorama

	sdkLocation, cfgLocation := securityPolicyLocationByDeviceType(device, "pre-rulebase")

	stateExpectedRuleName := func(idx int, value string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(
			"panos_security_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
		return plancheck.ExpectKnownValue(
			"panos_security_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			securityPolicyPreCheck(prefix, sdkLocation)

		},
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             securityPolicyCheckDestroy(prefix, sdkLocation),
		Steps: []resource.TestStep{
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"location":   cfgLocation,
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"location":   cfgLocation,
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
					"location":   cfgLocation,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-1"),
					stateExpectedRuleName(1, "rule-2"),
					stateExpectedRuleName(2, "rule-3"),
					stateExpectedRuleName(3, "rule-4"),
					stateExpectedRuleName(4, "rule-5"),
					ExpectServerSecurityRulesCount(prefix, sdkLocation, len(rulesInitial)),
					ExpectServerSecurityRulesOrder(prefix, sdkLocation, rulesInitial),
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
					"location":   cfgLocation,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(rulesReordered)...),
					"location":   cfgLocation,
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planExpectedRuleName(0, "rule-2"),
						planExpectedRuleName(1, "rule-1"),
						planExpectedRuleName(2, "rule-3"),
						planExpectedRuleName(3, "rule-4"),
						planExpectedRuleName(4, "rule-5"),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-1"),
					stateExpectedRuleName(2, "rule-3"),
					stateExpectedRuleName(3, "rule-4"),
					stateExpectedRuleName(4, "rule-5"),
					ExpectServerSecurityRulesOrder(prefix, sdkLocation, rulesReordered),
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"location":   cfgLocation,
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"location":   cfgLocation,
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func securityPolicyLocationByDeviceType(typ deviceType, rulebase string) (security.Location, config.Variable) {
	var sdkLocation security.Location
	var cfgLocation config.Variable
	switch typ {
	case devicePanorama:
		sdkLocation = security.Location{
			Shared: &security.SharedLocation{
				Rulebase: rulebase,
			},
		}
		cfgLocation = config.ObjectVariable(map[string]config.Variable{
			"shared": config.ObjectVariable(map[string]config.Variable{
				"rulebase": config.StringVariable(rulebase),
			}),
		})
	case deviceFirewall:
		sdkLocation = security.Location{
			Vsys: &security.VsysLocation{
				NgfwDevice: "localhost.localdomain",
				Vsys:       "vsys1",
			},
		}
		cfgLocation = config.ObjectVariable(map[string]config.Variable{
			"vsys": config.ObjectVariable(map[string]config.Variable{
				"name": config.StringVariable("vsys1"),
			}),
		})
	}

	return sdkLocation, cfgLocation
}

func securityPolicyPreCheck(prefix string, location security.Location) {
	service := security.NewService(sdkClient)
	ctx := context.TODO()

	stringPointer := func(value string) *string { return &value }

	rules := []security.Entry{
		{
			Name:        fmt.Sprintf("%s-rule0", prefix),
			Description: stringPointer("Rule 0"),
			Source:      []string{"any"},
			Destination: []string{"any"},
			From:        []string{"any"},
			To:          []string{"any"},
			Service:     []string{"any"},
		},
		{
			Name:        fmt.Sprintf("%s-rule99", prefix),
			Description: stringPointer("Rule 99"),
			Source:      []string{"any"},
			Destination: []string{"any"},
			From:        []string{"any"},
			To:          []string{"any"},
			Service:     []string{"any"},
		},
	}

	for _, elt := range rules {
		_, err := service.Create(ctx, location, &elt)
		if err != nil {
			panic(fmt.Sprintf("natPolicyPreCheck failed: %s", err))
		}

	}
}

func securityPolicyCheckDestroy(prefix string, location security.Location) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		service := security.NewService(sdkClient)
		ctx := context.TODO()

		rules, err := service.List(ctx, location, "get", "", "")
		if err != nil && !sdkerrors.IsObjectNotFound(err) {
			return err
		}

		var danglingNames []string
		for _, elt := range rules {
			if strings.HasPrefix(elt.Name, prefix) {
				danglingNames = append(danglingNames, elt.Name)
			}
		}

		if len(danglingNames) > 0 {
			err := DanglingObjectsError
			delErr := service.Delete(ctx, location, danglingNames...)
			if delErr != nil {
				err = errors.Join(err, delErr)
			}

			return err
		}

		return nil
	}
}

func init() {
	resource.AddTestSweepers("pango_security_policy", &resource.Sweeper{
		Name: "pango_security_policy",
		F: func(typ string) error {
			service := security.NewService(sdkClient)

			var deviceTyp deviceType
			switch typ {
			case "panorama":
				deviceTyp = devicePanorama
			case "firewall":
				deviceTyp = deviceFirewall
			default:
				panic("invalid device type")
			}

			for _, rulebase := range []string{"pre-rulebase", "post-rulebase"} {
				location, _ := securityPolicyLocationByDeviceType(deviceTyp, rulebase)
				ctx := context.TODO()
				objects, err := service.List(ctx, location, "get", "", "")
				if err != nil && !sdkerrors.IsObjectNotFound(err) {
					return fmt.Errorf("Failed to list Security Rules during sweep: %w", err)
				}

				var names []string
				for _, elt := range objects {
					if strings.HasPrefix(elt.Name, "test-acc") {
						names = append(names, elt.Name)
					}
				}

				if len(names) > 0 {
					err = service.Delete(ctx, location, names...)
					if err != nil {
						return fmt.Errorf("Failed to delete Security Rules during sweep: %w", err)
					}
				}
			}

			return nil
		},
	})
}
