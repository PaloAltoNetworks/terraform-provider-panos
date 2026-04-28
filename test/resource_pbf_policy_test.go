package provider_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/policies/rules/pbf"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// Helper functions
func prefixedPbf(prefix string, name string) string {
	return fmt.Sprintf("%s-%s", prefix, name)
}

func withPrefixPbf(prefix string, rules []string) []config.Variable {
	var result []config.Variable
	for _, elt := range rules {
		result = append(result, config.StringVariable(prefixedPbf(prefix, elt)))
	}
	return result
}

// Server-side validation helpers
type expectServerPbfRulesOrder struct {
	Location  pbf.Location
	Prefix    string
	RuleNames []string
}

func ExpectServerPbfRulesOrder(prefix string, ruleNames []string) *expectServerPbfRulesOrder {
	location := pbf.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	return &expectServerPbfRulesOrder{
		Location:  *location,
		Prefix:    prefix,
		RuleNames: ruleNames,
	}
}

func (o *expectServerPbfRulesOrder) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := pbf.NewService(sdkClient)

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

	for actualIdx, elt := range objects {
		if state, ok := rulesWithIdx[elt.Name]; ok {
			state.ActualIdx = actualIdx
			rulesWithIdx[elt.Name] = state
		}
	}

	var rulesError bool
	for _, state := range rulesWithIdx {
		if state.ActualIdx == -1 {
			resp.Error = fmt.Errorf("rule %s not found on server", state.Name)
			return
		}
		if state.ActualIdx != state.ExpectedIdx {
			rulesError = true
		}
	}

	if rulesError {
		var rulesView []string
		for _, elt := range objects {
			if state, ok := rulesWithIdx[elt.Name]; ok {
				finalElt := fmt.Sprintf("{%s: %d->%d}", state.Name, state.ExpectedIdx, state.ActualIdx)
				rulesView = append(rulesView, finalElt)
			}
		}
		resp.Error = fmt.Errorf("Unexpected server state: %s", strings.Join(rulesView, " "))
		return
	}
}

type expectServerPbfRulesCount struct {
	Prefix   string
	Location pbf.Location
	Count    int
}

func ExpectServerPbfRulesCount(prefix string, count int) *expectServerPbfRulesCount {
	location := pbf.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)
	return &expectServerPbfRulesCount{
		Prefix:   prefix,
		Location: *location,
		Count:    count,
	}
}

func (o *expectServerPbfRulesCount) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := pbf.NewService(sdkClient)

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
		resp.Error = fmt.Errorf("expected %d PBF rules with prefix %s, got %d", o.Count, o.Prefix, count)
		return
	}
}

// Test 1: Comprehensive test with forward action (nexthop IP)
const pbfPolicyBasicTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_administrative_tag" "rule_tag" {
  location = { shared = {} }

  name = format("%s-pbf-tag", var.prefix)
}

resource "panos_administrative_tag" "group_tag" {
  location = { shared = {} }

  name = format("%s-pbf-group", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule-basic", var.prefix)

    from = {
      zone = ["trust"]
    }

    source_addresses = ["10.1.1.0/24"]
    destination_addresses = ["192.168.1.0/24"]
    applications = ["ssl", "web-browsing"]
    services = ["service-http", "service-https"]
    source_users = ["any"]

    action = {
      forward = {
        nexthop = {
          ip_address = "10.0.0.1"
        }
        egress_interface = "ethernet1/1"
        monitor = {
          ip_address = "10.0.0.1"
          disable_if_unreachable = true
          profile = "default"
        }
      }
    }

    enforce_symmetric_return = {
      enabled = true
      nexthop_address_list = [
        { name = "10.0.0.2" },
        { name = "10.0.0.3" }
      ]
    }

    tags = [resource.panos_administrative_tag.rule_tag.name]
    group_tag = resource.panos_administrative_tag.group_tag.name
    disabled = false
    description = "Basic PBF rule with all options"
    negate_source = false
    negate_destination = false
    active_active_device_binding = "both"
  }]
}
`

func TestAccPbfPolicy_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyBasicTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-basic", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("from").AtMapKey("zone"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("trust"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("source_addresses"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("10.1.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("destination_addresses"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("forward").AtMapKey("nexthop").AtMapKey("ip_address"),
						knownvalue.StringExact("10.0.0.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("forward").AtMapKey("egress_interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("description"),
						knownvalue.StringExact("Basic PBF rule with all options"),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("disabled"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("active_active_device_binding"),
						knownvalue.StringExact("both"),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("tags"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-pbf-tag", prefix)),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("group_tag"),
						knownvalue.StringExact(fmt.Sprintf("%s-pbf-group", prefix)),
					),
				},
			},
		},
	})
}

// Test 2: Discard action variant
const pbfPolicyActionDiscardTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule-discard", var.prefix)

    from = {
      zone = ["trust"]
    }

    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]

    action = {
      discard = {}
    }
  }]
}
`

func TestAccPbfPolicy_Action_Discard(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyActionDiscardTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-discard", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("discard"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

// Test 3: Forward-to-vsys action variant
const pbfPolicyActionForwardToVsysTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule-forward-vsys", var.prefix)

    from = {
      zone = ["trust"]
    }

    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]

    action = {
      forward_to_vsys = "vsys2"
    }
  }]
}
`

func TestAccPbfPolicy_Action_ForwardToVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyActionForwardToVsysTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-forward-vsys", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("forward_to_vsys"),
						knownvalue.StringExact("vsys2"),
					),
				},
			},
		},
	})
}

// Test 4: No-pbf action variant
const pbfPolicyActionNoPbfTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule-no-pbf", var.prefix)

    from = {
      zone = ["trust"]
    }

    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]

    action = {
      no_pbf = {}
    }
  }]
}
`

func TestAccPbfPolicy_Action_NoPbf(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyActionNoPbfTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-no-pbf", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("no_pbf"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

// Test 5: Forward action with FQDN nexthop
const pbfPolicyActionForwardNexthopFqdnTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule-fqdn", var.prefix)

    from = {
      zone = ["trust"]
    }

    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]

    action = {
      forward = {
        nexthop = {
          fqdn = "router.example.com"
        }
        egress_interface = "ethernet1/2"
      }
    }
  }]
}
`

func TestAccPbfPolicy_Action_Forward_NexthopFqdn(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyActionForwardNexthopFqdnTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-fqdn", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("forward").AtMapKey("nexthop").AtMapKey("fqdn"),
						knownvalue.StringExact("router.example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("action").AtMapKey("forward").AtMapKey("egress_interface"),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

// Test 6: Interface source variant
const pbfPolicyFromInterfaceTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
    name = format("%s-rule-interface", var.prefix)

    from = {
      interface = ["ethernet1/3"]
    }

    source_addresses = ["any"]
    destination_addresses = ["any"]
    services = ["any"]

    action = {
      forward = {
        nexthop = {
          ip_address = "10.0.0.1"
        }
        egress_interface = "ethernet1/3"
      }
    }
  }]
}
`

func TestAccPbfPolicy_From_Interface(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyFromInterfaceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule-interface", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_pbf_policy.policy",
						tfjsonpath.New("rules").AtSliceIndex(0).AtMapKey("from").AtMapKey("interface"),
						knownvalue.SetExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/3"),
						}),
					),
				},
			},
		},
	})
}

// Test 7: Rule ordering
const pbfPolicyOrderingTmpl = `
variable "prefix" { type = string }
variable "rule_names" { type = list(string) }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_pbf_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [
    for index, name in var.rule_names: {
      name = name

      from = {
        zone = ["trust"]
      }

      source_addresses = ["any"]
      destination_addresses = ["any"]
      services = ["any"]

      action = {
        forward = {
          nexthop = {
            ip_address = format("10.0.0.%s", index + 1)
          }
          egress_interface = "ethernet1/1"
        }
      }
    }
  ]
}
`

func TestAccPbfPolicyOrdering(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	rulesInitial := []string{"rule-1", "rule-2", "rule-3"}
	rulesReordered := []string{"rule-2", "rule-1", "rule-3"}

	stateExpectedRuleName := func(idx int, value string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(
			"panos_pbf_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixedPbf(prefix, value)),
		)
	}

	planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
		return plancheck.ExpectKnownValue(
			"panos_pbf_policy.policy",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixedPbf(prefix, value)),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: pbfPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefixPbf(prefix, rulesInitial)...),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-1"),
					stateExpectedRuleName(1, "rule-2"),
					stateExpectedRuleName(2, "rule-3"),
					ExpectServerPbfRulesCount(prefix, len(rulesInitial)),
					ExpectServerPbfRulesOrder(prefix, rulesInitial),
				},
			},
			{
				Config: pbfPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefixPbf(prefix, rulesInitial)...),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: pbfPolicyOrderingTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefixPbf(prefix, rulesReordered)...),
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
					ExpectServerPbfRulesOrder(prefix, rulesReordered),
				},
			},
		},
	})
}
