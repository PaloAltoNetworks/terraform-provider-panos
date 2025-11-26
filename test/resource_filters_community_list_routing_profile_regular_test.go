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

func TestAccFiltersCommunityListRoutingProfile_Regular_Basic(t *testing.T) {
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
				Config: filtersCommunityListRoutingProfile_Regular_Basic_Tmpl,
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
						knownvalue.StringExact("Regular community list for testing"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"regular": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"regular_entries": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":   knownvalue.StringExact("1"),
										"action": knownvalue.StringExact("deny"),
										"communities": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.StringExact("100:200"),
										}),
									}),
								}),
							}),
							"extended": knownvalue.Null(),
							"large":    knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersCommunityListRoutingProfile_Regular_Basic_Tmpl = `
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
  description = "Regular community list for testing"

  type = {
    regular = {
      regular_entries = [
        {
          name = "1"
          action = "deny"
          communities = ["100:200"]
        }
      ]
    }
  }
}
`

func TestAccFiltersCommunityListRoutingProfile_Regular_Action_Permit(t *testing.T) {
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
				Config: filtersCommunityListRoutingProfile_Regular_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("regular").AtMapKey("regular_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersCommunityListRoutingProfile_Regular_Action_Permit_Tmpl = `
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
    regular = {
      regular_entries = [
        {
          name = "1"
          action = "permit"
          communities = ["300:400"]
        }
      ]
    }
  }
}
`

func TestAccFiltersCommunityListRoutingProfile_Regular_MultipleEntries(t *testing.T) {
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
				Config: filtersCommunityListRoutingProfile_Regular_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_community_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("regular").AtMapKey("regular_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("1"),
								"action": knownvalue.StringExact("permit"),
								"communities": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("100:200"),
									knownvalue.StringExact("300:400"),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("2"),
								"action": knownvalue.StringExact("deny"),
								"communities": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("500:600"),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersCommunityListRoutingProfile_Regular_MultipleEntries_Tmpl = `
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
    regular = {
      regular_entries = [
        {
          name = "1"
          action = "permit"
          communities = ["100:200", "300:400"]
        },
        {
          name = "2"
          action = "deny"
          communities = ["500:600"]
        }
      ]
    }
  }
}
`
