package provider_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango/policies/rules/security"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type expectServerSecurityRulesOrder struct {
	Location  security.Location
	Prefix    string
	RuleNames []string
}

func ExpectServerSecurityRulesOrder(prefix string, ruleNames []string) *expectServerSecurityRulesOrder {
	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	return &expectServerSecurityRulesOrder{
		Location:  *location,
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
		Name        string
		ExpectedIdx int
		ActualIdx   int
	}

	rulesWithIdx := make(map[string]ruleWithState)
	for idx, elt := range o.RuleNames {
		name := fmt.Sprintf("%s-%s", o.Prefix, elt)
		rulesWithIdx[name] = ruleWithState{
			Name:        name,
			ExpectedIdx: idx,
			ActualIdx:   -1,
		}
	}

	var serverRules []string
	for _, elt := range objects {
		serverRules = append(serverRules, elt.EntryName())
	}

	for actualIdx, elt := range objects {
		if state, ok := rulesWithIdx[elt.Name]; !ok {
			rulesWithIdx[elt.Name] = ruleWithState{
				Name:        elt.Name,
				ExpectedIdx: -1,
				ActualIdx:   actualIdx,
			}
		} else {
			state.ActualIdx = actualIdx
			rulesWithIdx[elt.Name] = state
		}
	}

	var missing []ruleWithState
	final := make([]*ruleWithState, len(rulesWithIdx))
	var rulesError bool
	for _, state := range rulesWithIdx {
		if state.ActualIdx == -1 || state.ActualIdx != state.ExpectedIdx {
			rulesError = true
		}

		if state.ActualIdx >= 0 {
			final[state.ActualIdx] = &state
		} else {
			missing = append(missing, state)
		}
	}

	for idx, elt := range final {
		if elt == nil {
			final[idx] = &missing[0]
			missing = missing[1:]
		}
	}

	if rulesError {
		var rulesView []string
		for _, elt := range final {
			finalElt := fmt.Sprintf("{%s: %d->%d}", elt.Name, elt.ExpectedIdx, elt.ActualIdx)
			rulesView = append(rulesView, finalElt)
		}

		resp.Error = fmt.Errorf("Unexpected server state: %s", strings.Join(rulesView, " "))
		return
	}
}

type expectServerSecurityRulesCount struct {
	Prefix   string
	Location security.Location
	Count    int
}

func ExpectServerSecurityRulesCount(prefix string, count int) *expectServerSecurityRulesCount {
	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)
	return &expectServerSecurityRulesCount{
		Prefix:   prefix,
		Location: *location,
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

const securityPolicyDuplicatedTmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
  templates = [ resource.panos_template.template.name ]
}


resource "panos_security_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [
    {
      name = format("%s-rule", var.prefix)
      source_zones     = ["any"]
      source_addresses = ["any"]

      destination_zones     = ["any"]
      destination_addresses = ["any"]
    },
    {
      name = format("%s-rule", var.prefix)
      source_zones     = ["any"]
      source_addresses = ["any"]

      destination_zones     = ["any"]
      destination_addresses = ["any"]
    }
  ]
}
`

const securityPolicyExtendedResource1Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
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

    destination_zones     = ["any"]
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

func TestAccSecurityPolicyDuplicatedPlan(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyDuplicatedTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile("Non-unique entry names in the list"),
			},
		},
	})
}

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
							AtMapKey("destination_zones"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
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
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_security_policy" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix) }}

  rules = [
    for index, name in var.rule_names: {
      name = name

      source_zones     = ["any"]
      source_addresses = ["any"]

      destination_zones     = ["any"]
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
			securityPolicyPreCheck(prefix)

		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable([]config.Variable{}...),
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable([]config.Variable{}...),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-1"),
					stateExpectedRuleName(1, "rule-2"),
					stateExpectedRuleName(2, "rule-3"),
					stateExpectedRuleName(3, "rule-4"),
					stateExpectedRuleName(4, "rule-5"),
					ExpectServerSecurityRulesCount(prefix, len(rulesInitial)),
					ExpectServerSecurityRulesOrder(prefix, rulesInitial),
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
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
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(rulesReordered)...),
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
					ExpectServerSecurityRulesOrder(prefix, rulesReordered),
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable([]config.Variable{}...),
				},
			},
			{
				Config: securityPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable([]config.Variable{}...),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func securityPolicyPreCheck(prefix string) {
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

	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	for _, elt := range rules {
		_, err := service.Create(ctx, *location, &elt)
		if err != nil {
			panic(fmt.Sprintf("natPolicyPreCheck failed: %s", err))
		}

	}
}
