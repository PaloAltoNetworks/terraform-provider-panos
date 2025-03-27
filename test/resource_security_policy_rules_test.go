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
	//"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const securityPolicyRulesPositionFirst = `
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" } }

  position = { where = "first" }

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

const securityPolicyRulesPositionIndirectlyBefore = `
variable "rule_names" { type = list(string) }
variable "prefix" { type = string }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" }}

  position = { where = "before", directly = false, pivot = format("%s-rule-99", var.prefix) }

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

const securityPolicyRulesPositionDirectlyBefore = `
variable "rule_names" { type = list(string) }
variable "prefix" { type = string }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" }}

  position = {
    where = "before"
    directly = true
    pivot = format("%s-rule-99", var.prefix)
  }

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

const securityPolicyRulesPositionDirectlyAfter = `
variable "rule_names" { type = list(string) }
variable "prefix" { type = string }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" }}

  position = {
    where = "after"
    directly = true
    pivot = format("%s-rule-0", var.prefix)
  }

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

const securityPolicyRulesPositionLast = `
variable "rule_names" { type = list(string) }
variable "prefix" { type = string }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" }}

  position = {
    where = "last"
  }

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

func TestAccSecurityPolicyRulesPositioning(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	ruleNames := []string{"rule-2", "rule-3", "rule-4", "rule-5", "rule-6"}

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
			"panos_security_policy_rules.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	// planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
	// 	return plancheck.ExpectKnownValue(
	// 		"panos_security_policy_rules.policy",
	// 		tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
	// 		knownvalue.StringExact(prefixed(value)),
	// 	)
	// }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			securityPolicyRulesPreCheck(prefix)

		},
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             securityPolicyRulesCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: securityPolicyRulesPositionFirst,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"prefix":     config.StringVariable(prefix),
				},
			},
			{
				Config: securityPolicyRulesPositionFirst,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"prefix":     config.StringVariable(prefix),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: securityPolicyRulesPositionFirst,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(ruleNames)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-3"),
					stateExpectedRuleName(2, "rule-4"),
					stateExpectedRuleName(3, "rule-5"),
					stateExpectedRuleName(4, "rule-6"),
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-2", "rule-3", "rule-4", "rule-5", "rule-6", "rule-0", "rule-1", "rule-99"}),
				},
			},
			{
				Config: securityPolicyRulesPositionIndirectlyBefore,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(ruleNames)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-3"),
					stateExpectedRuleName(2, "rule-4"),
					stateExpectedRuleName(3, "rule-5"),
					stateExpectedRuleName(4, "rule-6"),
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-2", "rule-3", "rule-4", "rule-5", "rule-6", "rule-0", "rule-1", "rule-99"}),
				},
			},
			{
				Config: securityPolicyRulesPositionDirectlyBefore,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(ruleNames)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-3"),
					stateExpectedRuleName(2, "rule-4"),
					stateExpectedRuleName(3, "rule-5"),
					stateExpectedRuleName(4, "rule-6"),
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-0", "rule-1", "rule-2", "rule-3", "rule-4", "rule-5", "rule-6", "rule-99"}),
				},
			},
			{
				Config: securityPolicyRulesPositionDirectlyAfter,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(ruleNames)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-3"),
					stateExpectedRuleName(2, "rule-4"),
					stateExpectedRuleName(3, "rule-5"),
					stateExpectedRuleName(4, "rule-6"),
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-0", "rule-2", "rule-3", "rule-4", "rule-5", "rule-6", "rule-1", "rule-99"}),
				},
			},
			{
				Config: securityPolicyRulesPositionLast,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(ruleNames)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-3"),
					stateExpectedRuleName(2, "rule-4"),
					stateExpectedRuleName(3, "rule-5"),
					stateExpectedRuleName(4, "rule-6"),
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-0", "rule-1", "rule-99", "rule-1", "rule-2", "rule-3", "rule-4", "rule-5"}),
				},
			},
		},
	})
}

func securityPolicyRulesPreCheck(prefix string) {
	service := security.NewService(sdkClient)
	ctx := context.TODO()

	stringPointer := func(value string) *string { return &value }

	location := security.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	rules := []security.Entry{
		{
			Name:        fmt.Sprintf("%s-rule-0", prefix),
			Description: stringPointer("Rule 0"),
			Source:      []string{"any"},
			Destination: []string{"any"},
			From:        []string{"any"},
			To:          []string{"any"},
			Service:     []string{"any"},
		},
		{
			Name:        fmt.Sprintf("%s-rule-1", prefix),
			Description: stringPointer("Rule 0"),
			Source:      []string{"any"},
			Destination: []string{"any"},
			From:        []string{"any"},
			To:          []string{"any"},
			Service:     []string{"any"},
		},
		{
			Name:        fmt.Sprintf("%s-rule-99", prefix),
			Description: stringPointer("Rule 99"),
			Source:      []string{"any"},
			Destination: []string{"any"},
			From:        []string{"any"},
			To:          []string{"any"},
			Service:     []string{"any"},
		},
	}

	for _, elt := range rules {
		_, err := service.Create(ctx, *location, &elt)
		if err != nil {
			panic(fmt.Sprintf("natPolicyPreCheck failed: %s", err))
		}

	}
}

func securityPolicyRulesCheckDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {

		location := security.NewDeviceGroupLocation()
		location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

		service := security.NewService(sdkClient)
		ctx := context.TODO()

		rules, err := service.List(ctx, *location, "get", "", "")
		if err != nil && !sdkerrors.IsObjectNotFound(err) {
			return err
		}

		var danglingNames []string

		seededRule := func(name string) bool {
			seeded := []string{"rule-0", "rule-1", "rule-99"}
			for _, elt := range seeded {
				if strings.HasSuffix(name, elt) {
					return true
				}
			}

			return false
		}

		for _, elt := range rules {
			if strings.HasPrefix(elt.Name, prefix) && !seededRule(elt.Name) {
				danglingNames = append(danglingNames, elt.Name)
			}
		}

		if len(danglingNames) > 0 {
			err := fmt.Errorf("%w: %s", DanglingObjectsError, strings.Join(danglingNames, ", "))
			delErr := service.Delete(ctx, *location, danglingNames...)
			if delErr != nil {
				err = errors.Join(err, delErr)
			}

			return err
		}

		return nil
	}
}
