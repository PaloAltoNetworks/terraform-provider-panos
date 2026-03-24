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

func TestAccFiltersCommunityListRoutingProfile_Large_Basic(t *testing.T) {
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
				Config: filtersCommunityListRoutingProfile_Large_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Large community list for testing"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"large": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"large_entries": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":                      knownvalue.StringExact("1"),
										"action":                    knownvalue.StringExact("deny"),
										"large_community_regexes": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.StringExact("^1000:.*:.*"),
										}),
									}),
								}),
							}),
							"extended": knownvalue.Null(),
							"regular":  knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersCommunityListRoutingProfile_Large_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_community_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "Large community list for testing"

  type = {
    large = {
      large_entries = [
        {
          name = "1"
          action = "deny"
          large_community_regexes = ["^1000:.*:.*"]
        }
      ]
    }
  }
}
`

func TestAccFiltersCommunityListRoutingProfile_Large_Action_Permit(t *testing.T) {
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
				Config: filtersCommunityListRoutingProfile_Large_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("large").AtMapKey("large_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersCommunityListRoutingProfile_Large_Action_Permit_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_community_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    large = {
      large_entries = [
        {
          name = "1"
          action = "permit"
          large_community_regexes = ["^2000:.*:.*"]
        }
      ]
    }
  }
}
`

func TestAccFiltersCommunityListRoutingProfile_Large_MultipleEntries(t *testing.T) {
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
				Config: filtersCommunityListRoutingProfile_Large_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("large").AtMapKey("large_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("1"),
								"action": knownvalue.StringExact("permit"),
								"large_community_regexes": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("^1000:.*:.*"),
									knownvalue.StringExact("^2000:.*:.*"),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("2"),
								"action": knownvalue.StringExact("deny"),
								"large_community_regexes": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("^3000:.*:.*"),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersCommunityListRoutingProfile_Large_MultipleEntries_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_community_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    large = {
      large_entries = [
        {
          name = "1"
          action = "permit"
          large_community_regexes = ["^1000:.*:.*", "^2000:.*:.*"]
        },
        {
          name = "2"
          action = "deny"
          large_community_regexes = ["^3000:.*:.*"]
        }
      ]
    }
  }
}
`
