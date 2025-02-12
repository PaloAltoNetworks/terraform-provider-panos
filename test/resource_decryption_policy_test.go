package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/policies/rules/decryption"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type expectServerDecryptionRulesOrder struct {
	Location  decryption.Location
	Prefix    string
	RuleNames []string
}

func ExpectServerDecryptionRulesOrder(prefix string, ruleNames []string) *expectServerDecryptionRulesOrder {
	location := decryption.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)
	return &expectServerDecryptionRulesOrder{
		Location:  *location,
		Prefix:    prefix,
		RuleNames: ruleNames,
	}
}

func (o *expectServerDecryptionRulesOrder) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := decryption.NewService(sdkClient)

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

type expectServerDecryptionRulesCount struct {
	Prefix   string
	Location decryption.Location
	Count    int
}

func ExpectServerDecryptionRulesCount(prefix string, count int) *expectServerDecryptionRulesCount {
	location := decryption.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)
	return &expectServerDecryptionRulesCount{
		Location: *location,
		Prefix:   prefix,
		Count:    count,
	}
}

func (o *expectServerDecryptionRulesCount) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := decryption.NewService(sdkClient)

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

const decryptionPolicyExtendedResource1Tmpl = `
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
  location = { device_group = { name = panos_template.template.name } }

  name = format("%s-tag1", var.prefix)
}

resource "panos_decryption_policy" "policy" {
  location = { device_group = { name = panos_template.template.name } }

  rules = [{
    name = "rule"
    description = "description"

    action = "decrypt"
    #category = ["category1"]

    destination_addresses = ["any"]
    destination_zones = ["any"]
    destination_hip = ["any"]

    disabled = true
    group_tag = panos_administrative_tag.tag.name

    log_fail = true
    #log_setting = "setting"
    log_success = true

    negate_destination = true
    negate_source = true

    #packet_broker_profile = "pb-profile"
    #profile = "profile"

    services = ["any"]

    source_addresses = ["any"]
    source_zones = ["any"]
    source_hip = ["any"]
    source_user = ["any"]

    tag = [panos_administrative_tag.tag.name]

    target = {
      #devices = [{ name = "device1", vsys = [{ name = "vsys1"}] }]
      negate = true
      tags = [panos_administrative_tag.tag.name]
    }

    type = {
      ssh_proxy = {}
    }
  }]
}
`

const decryptionPolicyCleanupTmpl = `
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

func TestAccDecryptionPolicyExtended(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             decryptionPolicyCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: decryptionPolicyExtendedResource1Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_decryption_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":        knownvalue.StringExact("rule"),
							"description": knownvalue.StringExact("description"),
							"action":      knownvalue.StringExact("decrypt"),
							"category":    knownvalue.Null(),
							"destination_addresses": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"destination_zones": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"destination_hip": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"disabled":              knownvalue.Bool(true),
							"group_tag":             knownvalue.StringExact(fmt.Sprintf("%s-tag1", prefix)),
							"log_fail":              knownvalue.Bool(true),
							"log_setting":           knownvalue.Null(),
							"log_success":           knownvalue.Bool(true),
							"negate_destination":    knownvalue.Bool(true),
							"negate_source":         knownvalue.Bool(true),
							"packet_broker_profile": knownvalue.Null(),
							"profile":               knownvalue.Null(),
							"services": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"source_addresses": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"source_zones": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"source_hip": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"source_user": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("any"),
							}),
							"tag": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact(fmt.Sprintf("%s-tag1", prefix)),
							}),
							"target": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"devices": knownvalue.Null(),
								"negate":  knownvalue.Bool(true),
								"tags":    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact(fmt.Sprintf("%s-tag1", prefix))}),
							}),
							"type": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ssh_proxy":              knownvalue.ObjectExact(map[string]knownvalue.Check{}),
								"ssl_forward_proxy":      knownvalue.Null(),
								"ssl_inbound_inspection": knownvalue.Null(),
							}),
							"uuid": knownvalue.NotNull(),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosDecryptionPolicyOrdering(t *testing.T) {
	rulesInitial := []string{"rule-1", "rule-2", "rule-3"}
	rulesReordered := []string{"rule-2", "rule-1", "rule-3"}

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

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
			"panos_decryption_policy.rules",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
		return plancheck.ExpectKnownValue(
			"panos_decryption_policy.rules",
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             decryptionPolicyCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: decryptionPolicyOrderTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-1"),
					stateExpectedRuleName(1, "rule-2"),
					stateExpectedRuleName(2, "rule-3"),
					ExpectServerDecryptionRulesCount(prefix, len(rulesInitial)),
					ExpectServerDecryptionRulesOrder(prefix, rulesInitial),
				},
			},
			{
				Config: decryptionPolicyOrderTmpl,
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
				Config: decryptionPolicyOrderTmpl,
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
					ExpectServerDecryptionRulesOrder(prefix, rulesReordered),
				},
			},
		},
	})
}

const decryptionPolicyOrderTmpl = `
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

resource "panos_decryption_policy" "rules" {
  location = { device_group = { name = panos_device_group.dg.name } }

  rules = [
    for index, name in var.rule_names: {
      name = name

      source_zones = ["any"]
      source_addresses = ["any"]
      destination_zones = ["any"]
      destination_addresses = ["any"]

      type = { ssh_proxy = {} }
    }
  ]
}
`

func decryptionPolicyCheckDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		service := decryption.NewService(sdkClient)
		ctx := context.TODO()

		location := decryption.NewDeviceGroupLocation()
		location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

		rules, err := service.List(ctx, *location, "get", "", "")
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
			delErr := service.Delete(ctx, *location, danglingNames...)
			if delErr != nil {
				err = errors.Join(err, delErr)
			}

			return err
		}

		return nil
	}
}
