package provider_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/PaloAltoNetworks/pango/policies/rules/security"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

const securityPolicyRulesImportInitial = `
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_security_policy_rules" "rules" {
  location = {
    device_group = {
      name = panos_device_group.example.name
      ruleset = "pre-ruleset"
    }
  }

  position = { where = "last" }

  rules = [
    for idx in range(2, 5) : {
        name = format("rule-%s", idx)

        source_addresses = ["any"]
        source_zones = ["any"]

        destination_addresses = ["any"]
        destination_zones = ["any"]

        services = ["any"]
        applications = ["any"]
    }
  ]
}
`

const securityPolicyRulesImportStep = `
resource "panos_security_policy_rules" "imported" {}
`

func TestAccSecurityPolicyRulesImport(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	configStep2 := mergeConfigs(
		securityPolicyRulesImportInitial,
		securityPolicyRulesImportStep,
	)

	importStateGenerateIDWithPrefixAndRules := func(rules []string) func(state *terraform.State) (string, error) {
		return func(state *terraform.State) (string, error) {
			return securityPolicyRulesGenerateImportID(state, prefix, rules)
		}
	}

	importStateGenerateIDInvalid := importStateGenerateIDWithPrefixAndRules([]string{"rule-2", "rule-3", "rule-4", "rule-5"})
	importStateGenerateIDValid := importStateGenerateIDWithPrefixAndRules([]string{"rule-2", "rule-3", "rule-4"})

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyRulesImportInitial,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
			{
				Config: configStep2,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ResourceName:      "panos_security_policy_rules.imported",
				ImportState:       true,
				ImportStateIdFunc: importStateGenerateIDInvalid,
				ExpectError:       regexp.MustCompile("Not all entries found on the server"),
			},
			{
				Config: configStep2,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ResourceName:      "panos_security_policy_rules.imported",
				ImportState:       true,
				ImportStateIdFunc: importStateGenerateIDValid,
			},
		},
	})
}

func securityPolicyRulesGenerateImportID(_ *terraform.State, prefix string, names []string) (string, error) {
	locationData := map[string]any{
		"device_group": map[string]any{
			"panorama_device": "localhost.localdomain",
			"name":            fmt.Sprintf("%s-dg", prefix),
			"rulebase":        "pre-rulebase",
		},
	}

	positionData := map[string]any{
		"where": "last",
	}

	importState := map[string]any{
		"position": positionData,
		"location": locationData,
		"names":    names,
	}

	marshalled, err := json.Marshal(importState)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal import state into JSON: %w", err)
	}

	return base64.StdEncoding.EncodeToString(marshalled), nil
}

const securityPolicyRulesPositionFirst = `
variable "position" { type = any }
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" } }

  position = var.position

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
variable "position" { type = any }
variable "rule_names" { type = list(string) }
variable "prefix" { type = string }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" }}

  position = var.position

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

func prefixed(prefix string, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}

func withPrefix(prefix string, rules []string) []config.Variable {
	var result []config.Variable
	for _, elt := range rules {
		result = append(result, config.StringVariable(prefixed(prefix, elt)))
	}

	return result
}

func TestAccSecurityPolicyRulesPositioning(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	ruleNames := []string{"rule-2", "rule-3", "rule-4", "rule-5", "rule-6"}

	stateExpectedRuleName := func(idx int, value string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(
			"panos_security_policy_rules.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(prefix, value)),
		)
	}

	// planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
	// 	return plancheck.ExpectKnownValue(
	// 		"panos_security_policy_rules.policy",
	// 		tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
	// 		knownvalue.StringExact(prefixed(prefix, value)),
	// 	)
	// }

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			securityPolicyRulesPreCheck(prefix)

		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyRulesPositionFirst,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"prefix":     config.StringVariable(prefix),
					"position": config.ObjectVariable(map[string]config.Variable{
						"where": config.StringVariable("first"),
					}),
				},
			},
			{
				Config: securityPolicyRulesPositionFirst,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable([]config.Variable{}...),
					"prefix":     config.StringVariable(prefix),
					"position": config.ObjectVariable(map[string]config.Variable{
						"where": config.StringVariable("first"),
					}),
				},
				PlanOnly:           true,
				ExpectNonEmptyPlan: false,
			},
			{
				Config: securityPolicyRulesPositionFirst,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
					"prefix":     config.StringVariable(prefix),
					"position": config.ObjectVariable(map[string]config.Variable{
						"where": config.StringVariable("first"),
					}),
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
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
					"prefix":     config.StringVariable(prefix),
					"position": config.ObjectVariable(map[string]config.Variable{
						"where":    config.StringVariable("before"),
						"directly": config.BoolVariable(false),
						"pivot":    config.StringVariable(fmt.Sprintf("%s-rule-99", prefix)),
					}),
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
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
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
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
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
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-2"),
					stateExpectedRuleName(1, "rule-3"),
					stateExpectedRuleName(2, "rule-4"),
					stateExpectedRuleName(3, "rule-5"),
					stateExpectedRuleName(4, "rule-6"),
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-0", "rule-1", "rule-99", "rule-2", "rule-3", "rule-4", "rule-5", "rule-6"}),
				},
			},
		},
	})
}

const securityPolicyRulesOrderingDependantInitial = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_security_policy_rules" "rule-1" {
  location = { device_group = { name = panos_device_group.dg.name } }

  position = { where = "first" }

  rules = [{
    name = format("%s-rule-1", var.prefix)

    source_zones     = ["any"]
    source_addresses = ["1.1.1.1"]

    destination_zones     = ["any"]
    destination_addresses = ["172.0.0.1/8"]

    services = ["any"]
    applications = ["any"]
  }]
}
`

const securityPolicyRulesOrderingDependant2 = `
resource "panos_security_policy_rules" "rule-255" {
  location = { device_group = { name = panos_device_group.dg.name } }

  position = { where = "last" }

  rules = [{
    name = format("%s-rule-255", var.prefix)

    source_zones     = ["any"]
    source_addresses = ["1.1.1.255"]

    destination_zones     = ["any"]
    destination_addresses = ["172.0.0.255/8"]

    services     = ["any"]
    applications = ["any"]
  }]
}
`

const securityPolicyRulesOrderingDependant3 = `
resource "panos_security_policy_rules" "example-directly-after" {
  location = { device_group = { name = panos_device_group.dg.name } }

  position = { where = "after", directly = true, pivot = "${var.prefix}-rule-1" }

  rules = [for k in [2, 3, ] :
    {
      name                  = "${var.prefix}-rule-${k}"
      source_zones          = ["any"]
      source_addresses      = ["1.1.1.${k}"]
      destination_zones     = ["any"]
      destination_addresses = ["172.0.0.${k}/8"]
      services              = ["any"]
      applications          = ["any"]
    }
  ]
}
`

const securityPolicyRulesOrderingDependant4 = `
resource "panos_security_policy_rules" "rules-after-rule-1" {
  location = { device_group = { name = panos_device_group.dg.name } }

  position = { where = "after", directly = false, pivot = format("%s-rule-1", var.prefix) }

  rules = [for k in [4, 5] :
    {
      name = format("%s-rule-%s", var.prefix, k)

      source_zones          = ["any"],
      source_addresses      = ["1.1.1.${k}"],
      destination_zones     = ["any"],
      destination_addresses = ["172.0.0.${k}/8"],
      services              = ["any"],
      applications          = ["any"],
    }
  ]
}
`

const securityPolicyRulesOrderingDependant5 = `
resource "panos_security_policy_rules" "rules-directly-before-rule-255" {
  location = { device_group = { name = panos_device_group.dg.name } }

  position = { where = "before", directly = true, pivot = "${var.prefix}-rule-255" }

  rules = [for k in [6, 7] :
    {
      name                  = "${var.prefix}-rule-${k}",
      source_zones          = ["any"],
      source_addresses      = ["1.1.1.${k}"],
      destination_zones     = ["any"],
      destination_addresses = ["172.0.0.${k}/8"],
      services              = ["any"],
      applications          = ["any"],
    }
  ]
}
`

const securityPolicyRulesOrderingDependant6 = `
resource "panos_security_policy_rules" "rules-before-rule-255" {
  location = { device_group = { name = panos_device_group.dg.name } }

  position = { where = "before", directly = false, pivot = "${var.prefix}-rule-255" }

  rules = [for k in [8, 9] :
    {
      name                  = "${var.prefix}-rule-${k}",
      source_zones          = ["any"],
      source_addresses      = ["1.1.1.${k}"],
      destination_zones     = ["any"],
      destination_addresses = ["172.0.0.${k}/8"],
      services              = ["any"],
      applications          = ["any"],
    }
  ]
}
`

func TestAccSecurityPolicyRulesOrderingDependant(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	configStep1 := securityPolicyRulesOrderingDependantInitial
	configStep2 := mergeConfigs(
		securityPolicyRulesOrderingDependantInitial,
		securityPolicyRulesOrderingDependant2,
	)
	configStep3 := mergeConfigs(
		configStep2,
		securityPolicyRulesOrderingDependant3,
	)
	configStep4 := mergeConfigs(
		configStep3,
		securityPolicyRulesOrderingDependant4,
	)
	configStep5 := mergeConfigs(
		configStep4,
		securityPolicyRulesOrderingDependant5,
	)
	configStep6 := mergeConfigs(
		configStep5,
		securityPolicyRulesOrderingDependant6,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)

		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configStep1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1"}),
				},
			},
			{
				Config: configStep2,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-255"}),
				},
			},
			{
				Config: configStep3,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-2", "rule-3", "rule-255"}),
				},
			},
			{
				Config: configStep4,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectNonEmptyPlan: true,
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-2", "rule-3", "rule-255", "rule-4", "rule-5"}),
				},
			},
			{
				Config: configStep4,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-2", "rule-3", "rule-4", "rule-5", "rule-255"}),
				},
			},
			{
				Config: configStep5,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-2", "rule-3", "rule-4", "rule-5", "rule-6", "rule-7", "rule-255"}),
				},
			},
			{
				Config: configStep6,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectNonEmptyPlan: true,
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-2", "rule-3", "rule-4", "rule-5", "rule-6", "rule-7", "rule-8", "rule-9", "rule-255"}),
				},
			},
			{
				Config: configStep6,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-1", "rule-2", "rule-3", "rule-4", "rule-5", "rule-8", "rule-9", "rule-6", "rule-7", "rule-255"}),
				},
			},
		},
	})
}

func TestAccSecurityPolicyRulesPositionAsVariable(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	ruleNames := []string{"rule-2", "rule-3"}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			securityPolicyRulesPreCheck(prefix)

		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: securityPolicyRulesPositionAsVariableTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
					"prefix":     config.StringVariable(prefix),
					"position": config.ObjectVariable(map[string]config.Variable{
						"where":    config.StringVariable("before"),
						"directly": config.BoolVariable(true),
						"pivot":    config.StringVariable(fmt.Sprintf("%s-rule-1", prefix)),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-0", "rule-2", "rule-3", "rule-1", "rule-99"}),
				},
			},
			{
				Config: securityPolicyRulesPositionAsVariableTmpl,
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(prefix, ruleNames)...),
					"prefix":     config.StringVariable(prefix),
					"position": config.ObjectVariable(map[string]config.Variable{
						"where":    config.StringVariable("before"),
						"directly": config.BoolVariable(true),
						"pivot":    config.StringVariable(fmt.Sprintf("%s-rule-99", prefix)),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ExpectServerSecurityRulesOrder(prefix, []string{"rule-0", "rule-1", "rule-2", "rule-3", "rule-99"}),
				},
			},
		},
	})
}

const securityPolicyRulesPositionAsVariableTmpl = `
variable "position" { type = any }
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_security_policy_rules" "policy" {
  location = { device_group = { name = format("%s-dg", var.prefix), rulebase = "pre-rulebase" } }

  position = var.position

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

const securityPolicyRules_Hierarchy_Initial_Tmpl = `
variable prefix { type = string }

resource "panos_device_group" "parent" {
  location = { panorama = {} }

  name = format("%s-parent", var.prefix)
}

resource "panos_device_group" "child" {
  location = { panorama = {} }

  name = format("%s-child", var.prefix)
}

resource "panos_device_group_parent" "relation" {
  location = { panorama = {} }

  device_group = panos_device_group.child.name
  parent       = panos_device_group.parent.name
}
`

const securityPolicyRules_Hierarchy_Parent_Entries_Tmpl = `
variable "parent_rule_names" {
  type = list(string)
}

resource "panos_security_policy_rules" "parent" {
  location = { device_group = { name = panos_device_group.parent.name } }

  position = { where = "first" }

  rules = [
    for index, name in var.parent_rule_names: {
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

const securityPolicyRules_Hierarchy_Child_Entries_Tmpl = `
variable "child_rule_names" {
  type = list(string)
}

resource "panos_security_policy_rules" "child" {
  location = { device_group = { name = panos_device_group.child.name} }

  position = { where = "first" }

  rules = [
    for index, name in var.child_rule_names: {
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

func testAccSecurityPolicyRules_Hierarchy(t *testing.T, parent config.Variable, child config.Variable) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	configStep1 := securityPolicyRules_Hierarchy_Initial_Tmpl
	configStep2 := mergeConfigs(
		securityPolicyRules_Hierarchy_Initial_Tmpl,
		securityPolicyRules_Hierarchy_Parent_Entries_Tmpl,
	)
	configStep3 := mergeConfigs(
		securityPolicyRules_Hierarchy_Initial_Tmpl,
		securityPolicyRules_Hierarchy_Parent_Entries_Tmpl,
		securityPolicyRules_Hierarchy_Child_Entries_Tmpl,
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			securityPolicyPreCheck(prefix)

		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configStep1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
			},
			{
				Config: configStep2,
				ConfigVariables: map[string]config.Variable{
					"prefix":            config.StringVariable(prefix),
					"parent_rule_names": parent,
				},
			},
			{
				Config: configStep3,
				ConfigVariables: map[string]config.Variable{
					"prefix":            config.StringVariable(prefix),
					"parent_rule_names": parent,
					"child_rule_names":  child,
				},
			},
		},
	})
}

func TestAccSecurityPolicyRules_Hierarchy_UniqueNames(t *testing.T) {
	parentRules := config.ListVariable(config.StringVariable("rule-1"), config.StringVariable("rule-2"))
	childRules := config.ListVariable(config.StringVariable("rule-3"), config.StringVariable("rule-4"))
	testAccSecurityPolicyRules_Hierarchy(t, parentRules, childRules)
}

func mergeConfigs(configs ...string) string {
	return strings.Join(configs, "\n")
}
