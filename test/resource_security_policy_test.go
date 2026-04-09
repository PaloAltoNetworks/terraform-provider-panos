package provider_test

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/profiles/customurlcategory"
	"github.com/PaloAltoNetworks/pango/movement"
	"github.com/PaloAltoNetworks/pango/policies/rules/security"
	"github.com/PaloAltoNetworks/pango/util"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func CreateServerSecurityRules(prefix string, ruleNames []string, beforePivot string) {
	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	service := security.NewService(sdkClient)
	ctx := context.TODO()

	var entries []*security.Entry
	for _, name := range ruleNames {
		entry := &security.Entry{
			Name:        prefixed(prefix, name),
			Source:      []string{"any"},
			Destination: []string{"any"},
			From:        []string{"any"},
			To:          []string{"any"},
			Service:     []string{"any"},
		}
		_, err := service.Create(ctx, *location, entry)
		if err != nil {
			panic(fmt.Sprintf("CreateServerSecurityRules failed: %s", err))
		}
		entries = append(entries, entry)
	}

	position := movement.PositionBefore{
		Pivot:    prefixed(prefix, beforePivot),
		Directly: true,
	}
	err := service.MoveGroup(ctx, *location, position, entries, 10)
	if err != nil {
		panic(fmt.Sprintf("CreateServerSecurityRules move failed: %s", err))
	}
}

func DeleteServerSecurityRules(prefix string, ruleNames []string) {
	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	service := security.NewService(sdkClient)

	var names []string
	for _, elt := range ruleNames {
		names = append(names, prefixed(prefix, elt))
	}

	err := service.Delete(context.TODO(), *location, names...)
	if err != nil {
		panic("failed to delete entries from the server")
	}

	return
}

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
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
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

	stateExpectedRuleName := func(idx int, value string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(
			"panos_security_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(prefix, value)),
		)
	}

	planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
		return plancheck.ExpectKnownValue(
			"panos_security_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(prefix, value)),
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
					"rule_names": config.ListVariable(withPrefix(prefix, rulesInitial)...),
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
					"rule_names": config.ListVariable(withPrefix(prefix, rulesInitial)...),
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
					"rule_names": config.ListVariable(withPrefix(prefix, rulesReordered)...),
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

const securityPolicy_UpdateMissing_Tmpl = `
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_security_policy_rules" "policy" {
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

func TestAccSecurityPolicy_UpdateMissing(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	rules := []string{"rule-1", "rule-2", "rule-3"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyRules_UpdateMissing_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(prefix, rules)...),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, rules),
				},
			},
			{
				Config: securityPolicyRules_UpdateMissing_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(prefix, rules)...),
				},
				PreConfig: func() {
					DeleteServerSecurityRules(prefix, []string{"rule-2"})
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, rules),
				},
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

const securityPolicy_DeletePartiallyMissing_Initial_Tmpl = `
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
      name = format("%s-rule-1", var.prefix)
      source_zones     = ["any"]
      source_addresses = ["any"]
      destination_zones     = ["any"]
      destination_addresses = ["any"]
      services = ["any"]
      applications = ["any"]
    },
    {
      name = format("%s-rule-2", var.prefix)
      source_zones     = ["any"]
      source_addresses = ["any"]
      destination_zones     = ["any"]
      destination_addresses = ["any"]
      services = ["any"]
      applications = ["any"]
    },
    {
      name = format("%s-rule-3", var.prefix)
      source_zones     = ["any"]
      source_addresses = ["any"]
      destination_zones     = ["any"]
      destination_addresses = ["any"]
      services = ["any"]
      applications = ["any"]
    },
    {
      name = format("%s-rule-4", var.prefix)
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

const securityPolicy_DeletePartiallyMissing_EmptyRules_Tmpl = `
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

  rules = []
}
`

const securityPolicy_DeletePartiallyMissing_Empty_Tmpl = `
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
`

func TestAccSecurityPolicy_DeletePartiallyMissing(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	rules := []string{"rule-1", "rule-2", "rule-3", "rule-4"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicy_DeletePartiallyMissing_Initial_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-2", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(2).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-3", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(3).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-4", prefix)),
					),
					ExpectServerSecurityRulesCount(prefix, 4),
					ExpectServerSecurityRulesOrder(prefix, rules),
				},
			},
			{
				Config: securityPolicy_DeletePartiallyMissing_EmptyRules_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesCount(prefix, 0),
				},
			},
			{
				Config: securityPolicy_DeletePartiallyMissing_Initial_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-2", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(2).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-3", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(3).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-4", prefix)),
					),
					ExpectServerSecurityRulesCount(prefix, 4),
					ExpectServerSecurityRulesOrder(prefix, rules),
				},
			},
			{
				PreConfig: func() {
					DeleteServerSecurityRules(prefix, []string{"rule-2", "rule-3"})
				},
				Config: securityPolicy_DeletePartiallyMissing_Empty_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesCount(prefix, 0),
				},
			},
		},
	})
}

// Custom state check to verify URL categories using SDK
type expectSecurityRuleCustomUrlCategory struct {
	RuleName           string
	DeviceGroup        string
	ExpectedCategories []string
}

func ExpectSecurityRuleCustomUrlCategory(ruleName string, deviceGroup string, expectedCategories []string) *expectSecurityRuleCustomUrlCategory {
	return &expectSecurityRuleCustomUrlCategory{
		RuleName:           ruleName,
		DeviceGroup:        deviceGroup,
		ExpectedCategories: expectedCategories,
	}
}

func (o *expectSecurityRuleCustomUrlCategory) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	// Use SDK to read the security rule
	securityService := security.NewService(sdkClient)
	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = o.DeviceGroup

	// Build XPath for the rule
	vn := sdkClient.Versioning()
	path, err := location.XpathWithComponents(vn, util.AsEntryXpath(o.RuleName))
	if err != nil {
		resp.Error = fmt.Errorf("failed to build xpath for security rule %s: %w", o.RuleName, err)
		return
	}

	// Read the rule from the server using XPath
	entry, err := securityService.ReadWithXpath(ctx, util.AsXpath(path), "get")
	if err != nil {
		resp.Error = fmt.Errorf("failed to read security rule %s: %w", o.RuleName, err)
		return
	}

	// Verify categories
	if len(entry.Category) != len(o.ExpectedCategories) {
		resp.Error = fmt.Errorf("expected %d categories, got %d. Expected: %v, Got: %v",
			len(o.ExpectedCategories), len(entry.Category), o.ExpectedCategories, entry.Category)
		return
	}

	// Check each expected category is present
	categoryMap := make(map[string]bool)
	for _, cat := range entry.Category {
		categoryMap[cat] = true
	}

	for _, expectedCat := range o.ExpectedCategories {
		if !categoryMap[expectedCat] {
			resp.Error = fmt.Errorf("expected category %s not found in rule. Got categories: %v",
				expectedCat, entry.Category)
			return
		}
	}
}

// Custom state check to verify custom URL category exists and has correct properties
type expectCustomUrlCategoryOnServer struct {
	CategoryName string
	DeviceGroup  string
	ExpectedUrls []string
}

func ExpectCustomUrlCategoryOnServer(categoryName string, deviceGroup string, expectedUrls []string) *expectCustomUrlCategoryOnServer {
	return &expectCustomUrlCategoryOnServer{
		CategoryName: categoryName,
		DeviceGroup:  deviceGroup,
		ExpectedUrls: expectedUrls,
	}
}

func (o *expectCustomUrlCategoryOnServer) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	// Use SDK to read the custom URL category
	urlCategoryService := customurlcategory.NewService(sdkClient)
	location := customurlcategory.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = o.DeviceGroup

	// Build XPath for the category
	vn := sdkClient.Versioning()
	path, err := location.XpathWithComponents(vn, util.AsEntryXpath(o.CategoryName))
	if err != nil {
		resp.Error = fmt.Errorf("failed to build xpath for custom URL category %s: %w", o.CategoryName, err)
		return
	}

	// Read the category from the server using XPath
	entry, err := urlCategoryService.ReadWithXpath(ctx, util.AsXpath(path), "get")
	if err != nil {
		resp.Error = fmt.Errorf("failed to read custom URL category %s: %w", o.CategoryName, err)
		return
	}

	// Verify the URLs
	if len(entry.List) != len(o.ExpectedUrls) {
		resp.Error = fmt.Errorf("expected %d URLs, got %d. Expected: %v, Got: %v",
			len(o.ExpectedUrls), len(entry.List), o.ExpectedUrls, entry.List)
		return
	}

	urlMap := make(map[string]bool)
	for _, url := range entry.List {
		urlMap[url] = true
	}

	for _, expectedUrl := range o.ExpectedUrls {
		if !urlMap[expectedUrl] {
			resp.Error = fmt.Errorf("expected URL %s not found in category. Got URLs: %v",
				expectedUrl, entry.List)
			return
		}
	}
}

const securityPolicyWithCustomUrlCategoryTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_custom_url_category" "test_category" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-custom-url-cat", var.prefix)
  type = "URL List"
  list = ["malicious.example.com", "phishing.example.org"]
  description = "Test custom URL category for security policy"
}

resource "panos_security_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-rule-with-custom-url", var.prefix)

    source_zones     = ["any"]
    source_addresses = ["any"]

    destination_zones     = ["any"]
    destination_addresses = ["any"]

    services     = ["any"]
    applications = ["any"]

    # Reference the custom URL category
    category = [panos_custom_url_category.test_category.name]

    action = "deny"
    log_end = true
  }]
}
`

const securityPolicyWithMultipleUrlCategoriesTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_custom_url_category" "malware_category" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-malware-urls", var.prefix)
  type = "URL List"
  list = ["malware1.example.com", "malware2.example.com"]
}

resource "panos_custom_url_category" "phishing_category" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-phishing-urls", var.prefix)
  type = "URL List"
  list = ["phishing1.example.com", "phishing2.example.com"]
}

resource "panos_security_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-rule-multi-url", var.prefix)

    source_zones     = ["any"]
    source_addresses = ["any"]

    destination_zones     = ["any"]
    destination_addresses = ["any"]

    services     = ["any"]
    applications = ["any"]

    # Reference multiple custom URL categories
    category = [
      panos_custom_url_category.malware_category.name,
      panos_custom_url_category.phishing_category.name,
    ]

    action = "deny"
    log_end = true
  }]
}
`

const securityPolicyUrlCategoryUpdatedTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_custom_url_category" "test_category" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-custom-url-cat", var.prefix)
  type = "URL List"
  # Updated list
  list = ["updated1.example.com", "updated2.example.com", "updated3.example.com"]
  description = "Updated test custom URL category"
}

resource "panos_security_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-rule-with-custom-url", var.prefix)

    source_zones     = ["any"]
    source_addresses = ["any"]

    destination_zones     = ["any"]
    destination_addresses = ["any"]

    services     = ["any"]
    applications = ["any"]

    # Still references the same category (by name)
    category = [panos_custom_url_category.test_category.name]

    action = "deny"
    log_end = true
  }]
}
`

func TestAccSecurityPolicy_WithCustomUrlCategory(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	deviceGroup := fmt.Sprintf("%s-dg", prefix)
	customUrlCategoryName := fmt.Sprintf("%s-custom-url-cat", prefix)
	ruleName := fmt.Sprintf("%s-rule-with-custom-url", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// Step 1: Create custom URL category and security policy referencing it
			{
				Config: securityPolicyWithCustomUrlCategoryTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Verify Terraform state
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.test_category",
						tfjsonpath.New("name"),
						knownvalue.StringExact(customUrlCategoryName),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.test_category",
						tfjsonpath.New("list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("malicious.example.com"),
							knownvalue.StringExact("phishing.example.org"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("category"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(customUrlCategoryName),
						}),
					),
					// Verify via SDK that custom URL category exists on server
					ExpectCustomUrlCategoryOnServer(
						customUrlCategoryName,
						deviceGroup,
						[]string{"malicious.example.com", "phishing.example.org"},
					),
					// Verify via SDK that security rule has correct category reference
					ExpectSecurityRuleCustomUrlCategory(
						ruleName,
						deviceGroup,
						[]string{customUrlCategoryName},
					),
				},
			},
			// Step 2: Update custom URL category (URLs change, but name stays the same)
			{
				Config: securityPolicyUrlCategoryUpdatedTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Verify updated URLs
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.test_category",
						tfjsonpath.New("list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("updated1.example.com"),
							knownvalue.StringExact("updated2.example.com"),
							knownvalue.StringExact("updated3.example.com"),
						}),
					),
					// Verify security rule still references the category
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("category"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(customUrlCategoryName),
						}),
					),
					// Verify via SDK
					ExpectCustomUrlCategoryOnServer(
						customUrlCategoryName,
						deviceGroup,
						[]string{"updated1.example.com", "updated2.example.com", "updated3.example.com"},
					),
					ExpectSecurityRuleCustomUrlCategory(
						ruleName,
						deviceGroup,
						[]string{customUrlCategoryName},
					),
				},
			},
			// Step 3: Plan-only check to ensure no drift
			{
				Config: securityPolicyUrlCategoryUpdatedTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
		},
	})
}

func TestAccSecurityPolicy_WithMultipleCustomUrlCategories(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	deviceGroup := fmt.Sprintf("%s-dg", prefix)
	malwareCategoryName := fmt.Sprintf("%s-malware-urls", prefix)
	phishingCategoryName := fmt.Sprintf("%s-phishing-urls", prefix)
	ruleName := fmt.Sprintf("%s-rule-multi-url", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyWithMultipleUrlCategoriesTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Verify both categories exist
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.malware_category",
						tfjsonpath.New("name"),
						knownvalue.StringExact(malwareCategoryName),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.phishing_category",
						tfjsonpath.New("name"),
						knownvalue.StringExact(phishingCategoryName),
					),
					// Verify security rule references both categories
					statecheck.ExpectKnownValue(
						"panos_security_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("category"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(malwareCategoryName),
							knownvalue.StringExact(phishingCategoryName),
						}),
					),
					// Verify via SDK that both categories are on the server
					ExpectCustomUrlCategoryOnServer(
						malwareCategoryName,
						deviceGroup,
						[]string{"malware1.example.com", "malware2.example.com"},
					),
					ExpectCustomUrlCategoryOnServer(
						phishingCategoryName,
						deviceGroup,
						[]string{"phishing1.example.com", "phishing2.example.com"},
					),
					// Verify via SDK that security rule references both categories
					ExpectSecurityRuleCustomUrlCategory(
						ruleName,
						deviceGroup,
						[]string{malwareCategoryName, phishingCategoryName},
					),
				},
			},
		},
	})
}
