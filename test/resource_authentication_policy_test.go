package provider_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango/policies/rules/authentication"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type expectServerAuthenticationRulesOrder struct {
	Location  authentication.Location
	Prefix    string
	RuleNames []string
}

func ExpectServerAuthenticationRulesOrder(prefix string, ruleNames []string) *expectServerAuthenticationRulesOrder {
	location := authentication.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	return &expectServerAuthenticationRulesOrder{
		Location:  *location,
		Prefix:    prefix,
		RuleNames: ruleNames,
	}
}

func (o *expectServerAuthenticationRulesOrder) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := authentication.NewService(sdkClient)

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

type expectServerAuthenticationRulesCount struct {
	Prefix   string
	Location authentication.Location
	Count    int
}

func ExpectServerAuthenticationRulesCount(prefix string, count int) *expectServerAuthenticationRulesCount {
	location := authentication.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)
	return &expectServerAuthenticationRulesCount{
		Location: *location,
		Prefix:   prefix,
		Count:    count,
	}
}

func (o *expectServerAuthenticationRulesCount) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := authentication.NewService(sdkClient)

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

// TestAccAuthenticationPolicy_Basic - Comprehensive test covering 19 standard parameters
const authenticationPolicyBasicTmpl = `
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

resource "panos_administrative_tag" "tag" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-tag1", var.prefix)
}

resource "panos_authentication_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-rule1", var.prefix)
    description = "Authentication policy rule for testing"

    category = ["any"]

    source_zones = ["any"]
    destination_zones = ["any"]

    source_addresses = ["any"]
    destination_addresses = ["any"]

    source_hip = ["any"]
    destination_hip = ["any"]

    source_users = ["any"]

    services = ["any"]

    negate_source = true
    negate_destination = false

    disabled = false

    log_authentication_timeout = true

    timeout = 120

    tags = [panos_administrative_tag.tag.name]
  }]
}
`

func TestAccAuthenticationPolicy_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: authenticationPolicyBasicTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("description"),
						knownvalue.StringExact("Authentication policy rule for testing"),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("category"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_zones"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_zones"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_addresses"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_addresses"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_hip"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_hip"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_users"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("services"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("negate_source"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("negate_destination"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("log_authentication_timeout"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("timeout"),
						knownvalue.Int64Exact(120),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tag1", prefix)),
						}),
					),
				},
			},
		},
	})
}

// TestAccAuthenticationPolicy_Target - Test the nested target object parameter
const authenticationPolicyTargetTmpl = `
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

resource "panos_administrative_tag" "tag" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-tag1", var.prefix)
}

resource "panos_authentication_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-rule-target", var.prefix)

    source_zones = ["any"]
    destination_zones = ["any"]
    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]

    target = {
      negate = true
      tags = [panos_administrative_tag.tag.name]
    }
  }]
}
`

func TestAccAuthenticationPolicy_Target(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: authenticationPolicyTargetTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-target", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("target").
							AtMapKey("negate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("target").
							AtMapKey("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tag1", prefix)),
						}),
					),
				},
			},
		},
	})
}

// TestAccAuthenticationPolicy_Ordering - Test rule ordering functionality
const authenticationPolicyOrderingTmpl = `
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
  templates = [ resource.panos_template.template.name ]
}

resource "panos_authentication_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [
    for index, name in var.rule_names: {
      name = name

      source_zones = ["any"]
      destination_zones = ["any"]
      source_addresses = ["any"]
      destination_addresses = ["any"]
      services = ["any"]
    }
  ]
}
`

func TestAccAuthenticationPolicy_Ordering(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	rulesInitial := []string{"rule-1", "rule-2", "rule-3"}
	rulesReordered := []string{"rule-2", "rule-1", "rule-3"}

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
			"panos_authentication_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
		return plancheck.ExpectKnownValue(
			"panos_authentication_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: authenticationPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-1"),
					stateExpectedRuleName(1, "rule-2"),
					stateExpectedRuleName(2, "rule-3"),
					ExpectServerAuthenticationRulesCount(prefix, len(rulesInitial)),
					ExpectServerAuthenticationRulesOrder(prefix, rulesInitial),
				},
			},
			{
				Config: authenticationPolicyOrderingTmpl,
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
				Config: authenticationPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(rulesReordered)...),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						planExpectedRuleName(0, "rule-2"),
						planExpectedRuleName(1, "rule-1"),
						planExpectedRuleName(2, "rule-3"),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-1"),
					stateExpectedRuleName(2, "rule-3"),
					ExpectServerAuthenticationRulesOrder(prefix, rulesReordered),
				},
			},
		},
	})
}

// TestAccAuthenticationPolicy_Minimal - Test minimal configuration with defaults
const authenticationPolicyMinimalTmpl = `
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

resource "panos_authentication_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-minimal", var.prefix)

    source_zones = ["any"]
    destination_zones = ["any"]
    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]
  }]
}
`

func TestAccAuthenticationPolicy_Minimal(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: authenticationPolicyMinimalTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-minimal", prefix)),
					),
					// Verify default timeout value (60 minutes according to spec)
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("timeout"),
						knownvalue.Int64Exact(60),
					),
					// Verify optional fields are null/not set
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("description"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("authentication_enforcement"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("disabled"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("negate_source"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("negate_destination"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

// TestAccAuthenticationPolicy_Validation - Test edge cases and validation constraints
const authenticationPolicyValidationDescriptionTmpl = `
variable "prefix" { type = string }
variable "description" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
  templates = [ resource.panos_template.template.name ]
}

resource "panos_authentication_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-validation", var.prefix)
    description = var.description

    source_zones = ["any"]
    destination_zones = ["any"]
    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]
  }]
}
`

const authenticationPolicyValidationTimeoutTmpl = `
variable "prefix" { type = string }
variable "timeout" { type = number }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
  templates = [ resource.panos_template.template.name ]
}

resource "panos_authentication_policy" "policy" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [{
    name = format("%s-validation", var.prefix)
    timeout = var.timeout

    source_zones = ["any"]
    destination_zones = ["any"]
    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]
  }]
}
`


func TestAccAuthenticationPolicy_Validation(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate strings for length validation
	maxDescription := strings.Repeat("a", 1024)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// Test maximum description length (1024 chars)
			{
				Config: authenticationPolicyValidationDescriptionTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"description": config.StringVariable(maxDescription),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("description"),
						knownvalue.StringExact(maxDescription),
					),
				},
			},
			// Test timeout minimum value (1)
			{
				Config: authenticationPolicyValidationTimeoutTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":  config.StringVariable(prefix),
					"timeout": config.IntegerVariable(1),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("timeout"),
						knownvalue.Int64Exact(1),
					),
				},
			},
			// Test timeout maximum value (1440)
			{
				Config: authenticationPolicyValidationTimeoutTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":  config.StringVariable(prefix),
					"timeout": config.IntegerVariable(1440),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_authentication_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("timeout"),
						knownvalue.Int64Exact(1440),
					),
				},
			},
		},
	})
}
