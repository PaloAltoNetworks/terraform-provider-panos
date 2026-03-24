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

func TestAccFiltersAccessListRoutingProfile_Ipv6_Basic(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv6_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("IPv6 ACL test"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("10"),
								"action": knownvalue.StringExact("deny"),
								"source_address": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv6_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "IPv6 ACL test"

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv6_SourceAddress_Any(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv6_SourceAddress_Any_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("source_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"address": knownvalue.StringExact("any"),
							"entry":   knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv6_SourceAddress_Any_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
          source_address = {
            address = "any"
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv6_SourceAddress_Entry(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv6_SourceAddress_Entry_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("source_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"address": knownvalue.Null(),
							"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"address":     knownvalue.StringExact("2001:db8::/32"),
								"exact_match": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv6_SourceAddress_Entry_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
          source_address = {
            entry = {
              address = "2001:db8::/32"
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv6_SourceAddress_Entry_ExactMatch(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv6_SourceAddress_Entry_ExactMatch_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("source_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"address": knownvalue.Null(),
							"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"address":     knownvalue.StringExact("2001:db8::/64"),
								"exact_match": knownvalue.Bool(true),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv6_SourceAddress_Entry_ExactMatch_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
          source_address = {
            entry = {
              address = "2001:db8::/64"
              exact_match = true
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv6_Action_Permit(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv6_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv6_Action_Permit_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "permit"
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv6_MultipleEntries(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv6_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("10"),
								"action": knownvalue.StringExact("permit"),
								"source_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.StringExact("any"),
									"entry":   knownvalue.Null(),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("20"),
								"action": knownvalue.StringExact("deny"),
								"source_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.Null(),
									"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"address":     knownvalue.StringExact("2001:db8:1::/48"),
										"exact_match": knownvalue.Bool(true),
									}),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("30"),
								"action": knownvalue.StringExact("permit"),
								"source_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.Null(),
									"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"address":     knownvalue.StringExact("fd00::/8"),
										"exact_match": knownvalue.Null(),
									}),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv6_MultipleEntries_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_access_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
        },
        {
          name = "20"
          action = "deny"
          source_address = {
            entry = {
              address = "2001:db8:1::/48"
              exact_match = true
            }
          }
        },
        {
          name = "30"
          action = "permit"
          source_address = {
            entry = {
              address = "fd00::/8"
            }
          }
        }
      ]
    }
  }
}
`
