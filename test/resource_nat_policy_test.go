package provider_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"text/template"

	"github.com/PaloAltoNetworks/pango/policies/rules/nat"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type deviceType int

const (
	devicePanorama deviceType = iota
	deviceFirewall
)

var (
	UnexpectedRulesError = errors.New("exhaustive resource didn't delete existing rules")
	DanglingObjectsError = errors.New("some objects were not deleted by the provider")
)

type expectServerNatRulesOrder struct {
	Location  nat.Location
	Prefix    string
	RuleNames []string
}

func ExpectServerNatRulesOrder(prefix string, ruleNames []string) *expectServerNatRulesOrder {
	location := nat.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	return &expectServerNatRulesOrder{
		Location:  *location,
		Prefix:    prefix,
		RuleNames: ruleNames,
	}
}

func (o *expectServerNatRulesOrder) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := nat.NewService(sdkClient)

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

type expectServerNatRulesCount struct {
	Prefix   string
	Location nat.Location
	Count    int
}

func ExpectServerNatRulesCount(prefix string, count int) *expectServerNatRulesCount {
	location := nat.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

	return &expectServerNatRulesCount{
		Prefix:   prefix,
		Location: *location,
		Count:    count,
	}
}

func (o *expectServerNatRulesCount) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := nat.NewService(sdkClient)

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

const natPolicyExtendedResource1Tmpl = `
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

resource "panos_nat_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
      name = format("%s-rule1", var.prefix)
      source_zones = ["any"]
      source_addresses = ["any"]
      destination_zone = ["external"]
      destination_addresses = ["any"]

      source_translation = {
        dynamic_ip_and_port = {
          translated_address = ["1.1.1.1"]
        }
      }

      destination_translation = {
        translated_address = "1.1.1.1"
        translated_port = 443
        dns_rewrite = { direction = "reverse"}
      }
  }]
}
`

const natPolicyExtendedResource2Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl2", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg2", var.prefix)
  templates = [ resource.panos_template.template.name ]
}

resource "panos_ethernet_interface" "interface" {
  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "1.1.1.1" }]
  }
}

resource "panos_nat_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
      name = format("%s-rule2", var.prefix)
      source_zones = ["any"]
      source_addresses = ["any"]
      destination_zone = ["external"]
      destination_addresses = ["any"]

      source_translation = {
        dynamic_ip_and_port = {
          interface_address = {
            interface = resource.panos_ethernet_interface.interface.name
            ip        = "1.1.1.1"
          }
        }
      }

      dynamic_destination_translation = {
        translated_address = "1.1.1.1"
        translated_port = 443
        distribution = "least-sessions"
      }

      active_active_device_binding = "primary"
  }]
}
`

const natPolicyExtendedResource3Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl3", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg3", var.prefix)
  templates = [ resource.panos_template.template.name ]
}

resource "panos_ethernet_interface" "interface" {
  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "1.1.1.1" }]
  }
}

resource "panos_nat_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
      name = format("%s-rule3", var.prefix)
      source_zones = ["any"]
      source_addresses = ["any"]
      destination_zone = ["external"]
      destination_addresses = ["any"]

      source_translation = {
        dynamic_ip = {
          translated_address = ["172.16.0.1"]
          fallback = {
            translated_address = ["192.168.0.1"]
          }
        }
      }
  }]
}
`

const natPolicyExtendedResource4Tmpl = `
variable "prefix" { type = string }

resource "panos_template" "template" {
  location = { panorama = {} }

  name = format("%s-tmpl5", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg5", var.prefix)
  templates = [ resource.panos_template.template.name ]
}

resource "panos_ethernet_interface" "interface" {
  location = { template = { name = resource.panos_template.template.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "192.168.0.1" }]
  }
}

resource "panos_nat_policy" "policy" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [{
      name = format("%s-rule5", var.prefix)
      source_zones = ["any"]
      source_addresses = ["any"]
      destination_zone = ["external"]
      destination_addresses = ["any"]

      source_translation = {
        dynamic_ip = {
          translated_address = ["172.16.0.1"]
          fallback = {
            interface_address = {
              interface = resource.panos_ethernet_interface.interface.name
              ip = "192.168.0.1"
            }
          }
        }
      }
  }]
}
`

func TestAccNatPolicyExtended(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: natPolicyExtendedResource1Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rule1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_zones"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_addresses"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_zone"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("external"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_addresses"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("any"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip_and_port").
							AtMapKey("translated_address"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("1.1.1.1"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_translation").
							AtMapKey("translated_address"),
						knownvalue.StringExact("1.1.1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_translation").
							AtMapKey("translated_port"),
						knownvalue.Int64Exact(443),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_translation").
							AtMapKey("dns_rewrite").
							AtMapKey("direction"),
						knownvalue.StringExact("reverse"),
					),
				},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: natPolicyExtendedResource2Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip_and_port").
							AtMapKey("translated_address"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip_and_port").
							AtMapKey("interface_address").
							AtMapKey("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip_and_port").
							AtMapKey("interface_address").
							AtMapKey("ip"),
						knownvalue.StringExact("1.1.1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("destination_translation"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("dynamic_destination_translation").
							AtMapKey("translated_address"),
						knownvalue.StringExact("1.1.1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("dynamic_destination_translation").
							AtMapKey("translated_port"),
						knownvalue.Int64Exact(443),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("dynamic_destination_translation").
							AtMapKey("distribution"),
						knownvalue.StringExact("least-sessions"),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("active_active_device_binding"),
						knownvalue.StringExact("primary"),
					),
				},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: natPolicyExtendedResource3Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip_and_port"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip").
							AtMapKey("translated_address"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("172.16.0.1"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip").
							AtMapKey("fallback").
							AtMapKey("translated_address"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.0.1"),
						}),
					),
				},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: natPolicyExtendedResource4Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip_and_port"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip").
							AtMapKey("translated_address"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("172.16.0.1"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip").
							AtMapKey("fallback").
							AtMapKey("interface_address").
							AtMapKey("interface"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_nat_policy.policy",
						tfjsonpath.New("rules").
							AtSliceIndex(0).
							AtMapKey("source_translation").
							AtMapKey("dynamic_ip").
							AtMapKey("fallback").
							AtMapKey("interface_address").
							AtMapKey("ip"),
						knownvalue.StringExact("192.168.0.1"),
					),
				},
			},
		},
	})
}

func TestAccPanosNatPolicyOrdering(t *testing.T) {
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

	device := devicePanorama

	sdkLocation, _ := natPolicyLocationByDeviceType(device, "pre-rulebase")

	stateExpectedRuleName := func(idx int, value string) statecheck.StateCheck {
		return statecheck.ExpectKnownValue(
			fmt.Sprintf("panos_nat_policy.%s", prefix),
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	planExpectedRuleName := func(idx int, value string) plancheck.PlanCheck {
		return plancheck.ExpectKnownValue(
			fmt.Sprintf("panos_nat_policy.%s", prefix),
			tfjsonpath.New("rules").AtSliceIndex(idx).AtMapKey("name"),
			knownvalue.StringExact(prefixed(value)),
		)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			natPolicyPreCheck(prefix, sdkLocation)
		},
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makeNatPolicyConfig(prefix),
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					stateExpectedRuleName(0, "rule-1"),
					stateExpectedRuleName(1, "rule-2"),
					stateExpectedRuleName(2, "rule-3"),
					ExpectServerNatRulesCount(prefix, len(rulesInitial)),
					ExpectServerNatRulesOrder(prefix, rulesInitial),
				},
			},
			{
				Config: makeNatPolicyConfig(prefix),
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(rulesInitial)...),
					"prefix":     config.StringVariable(prefix),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: makeNatPolicyConfig(prefix),
				ConfigVariables: map[string]config.Variable{
					"rule_names": config.ListVariable(withPrefix(rulesReordered)...),
					"prefix":     config.StringVariable(prefix),
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
					ExpectServerNatRulesOrder(prefix, rulesReordered),
				},
			},
		},
	})
}

const configTmpl = `
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

resource "panos_nat_policy" "{{ .ResourceName }}" {
  location = { device_group = { name = resource.panos_device_group.dg.name }}

  rules = [
    for index, name in var.rule_names: {
      name = name

      source_zones = ["any"]
      source_addresses = ["any"]
      destination_zone = ["external"]
      destination_addresses = ["any"]

      destination_translation = {
        translated_address = format("172.16.0.%s", index)
      }
    }
  ]
}
`

func makeNatPolicyConfig(prefix string) string {
	var buf bytes.Buffer
	tmpl := template.Must(template.New("").Parse(configTmpl))

	context := struct {
		ResourceName string
	}{
		ResourceName: prefix,
	}

	err := tmpl.Execute(&buf, context)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

func natPolicyLocationByDeviceType(typ deviceType, rulebase string) (nat.Location, config.Variable) {
	var sdkLocation nat.Location
	var cfgLocation config.Variable
	switch typ {
	case devicePanorama:
		sdkLocation = nat.Location{
			Shared: &nat.SharedLocation{
				Rulebase: rulebase,
			},
		}
		cfgLocation = config.ObjectVariable(map[string]config.Variable{
			"shared": config.ObjectVariable(map[string]config.Variable{
				"rulebase": config.StringVariable(rulebase),
			}),
		})
	case deviceFirewall:
		sdkLocation = nat.Location{
			Vsys: &nat.VsysLocation{
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

func natPolicyPreCheck(prefix string, location nat.Location) {
	service := nat.NewService(sdkClient)
	ctx := context.TODO()

	stringPointer := func(value string) *string { return &value }

	rules := []nat.Entry{
		{
			Name:        fmt.Sprintf("%s-rule0", prefix),
			Description: stringPointer("Rule 0"),
			From:        []string{"any"},
			To:          []string{"external"},
			Source:      []string{"any"},
			Destination: []string{"any"},
		},
		{
			Name:        fmt.Sprintf("%s-rule99", prefix),
			Description: stringPointer("Rule 99"),
			From:        []string{"any"},
			To:          []string{"external"},
			Source:      []string{"any"},
			Destination: []string{"any"},
		},
	}

	for _, elt := range rules {
		_, err := service.Create(ctx, location, &elt)
		if err != nil {
			panic(fmt.Sprintf("natPolicyPreCheck failed: %s", err))
		}
	}
}
