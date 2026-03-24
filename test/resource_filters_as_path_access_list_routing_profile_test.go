package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccFiltersAsPathAccessListRoutingProfile_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: filtersAsPathAccessListRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_as_path_access_list_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_as_path_access_list_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("AS Path Access List for testing"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_as_path_access_list_routing_profile.example",
						tfjsonpath.New("aspath_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":         knownvalue.StringExact("1"),
								"action":       knownvalue.StringExact("deny"),
								"aspath_regex": knownvalue.StringExact("^65000_"),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAsPathAccessListRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_as_path_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "AS Path Access List for testing"

  aspath_entries = [
    {
      name = "1"
      action = "deny"
      aspath_regex = "^65000_"
    }
  ]
}
`

func TestAccFiltersAsPathAccessListRoutingProfile_Action_Permit(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: filtersAsPathAccessListRoutingProfile_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_as_path_access_list_routing_profile.example",
						tfjsonpath.New("aspath_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersAsPathAccessListRoutingProfile_Action_Permit_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_as_path_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  aspath_entries = [
    {
      name = "1"
      action = "permit"
      aspath_regex = "^65100_"
    }
  ]
}
`

func TestAccFiltersAsPathAccessListRoutingProfile_MultipleEntries(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: filtersAsPathAccessListRoutingProfile_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_as_path_access_list_routing_profile.example",
						tfjsonpath.New("aspath_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":         knownvalue.StringExact("1"),
								"action":       knownvalue.StringExact("permit"),
								"aspath_regex": knownvalue.StringExact("^65000_"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":         knownvalue.StringExact("2"),
								"action":       knownvalue.StringExact("deny"),
								"aspath_regex": knownvalue.StringExact("_65100$"),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAsPathAccessListRoutingProfile_MultipleEntries_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_as_path_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  aspath_entries = [
    {
      name = "1"
      action = "permit"
      aspath_regex = "^65000_"
    },
    {
      name = "2"
      action = "deny"
      aspath_regex = "_65100$"
    }
  ]
}
`
