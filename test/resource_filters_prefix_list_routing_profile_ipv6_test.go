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

func TestAccFiltersPrefixListRoutingProfile_Ipv6_Basic(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv6_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Test IPv6 prefix list"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ipv4": knownvalue.Null(),
							"ipv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ipv6_entries": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":   knownvalue.StringExact("10"),
										"action": knownvalue.StringExact("deny"),
										"prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"network": knownvalue.Null(),
											"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"network":                knownvalue.StringExact("2001:db8::/32"),
												"greater_than_or_equal": knownvalue.Null(),
												"less_than_or_equal":    knownvalue.Null(),
											}),
										}),
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

const filtersPrefixListRoutingProfile_Ipv6_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "Test IPv6 prefix list"

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "2001:db8::/32"
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersPrefixListRoutingProfile_Ipv6_Prefix_Network(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv6_Prefix_Network_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("prefix"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"network": knownvalue.StringExact("any"),
							"entry":   knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersPrefixListRoutingProfile_Ipv6_Prefix_Network_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          prefix = {
            network = "any"
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersPrefixListRoutingProfile_Ipv6_Prefix_Entry(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv6_Prefix_Entry_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("prefix"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"network": knownvalue.Null(),
							"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"network":                knownvalue.StringExact("fd00::/8"),
								"greater_than_or_equal": knownvalue.Int64Exact(64),
								"less_than_or_equal":    knownvalue.Int64Exact(96),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersPrefixListRoutingProfile_Ipv6_Prefix_Entry_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          prefix = {
            entry = {
              network = "fd00::/8"
              greater_than_or_equal = 64
              less_than_or_equal = 96
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersPrefixListRoutingProfile_Ipv6_Action_Permit(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv6_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersPrefixListRoutingProfile_Ipv6_Action_Permit_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "permit"
          prefix = {
            network = "any"
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersPrefixListRoutingProfile_Ipv6_MultipleEntries(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv6_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv6").AtMapKey("ipv6_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("10"),
								"action": knownvalue.StringExact("permit"),
								"prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"network": knownvalue.StringExact("any"),
									"entry":   knownvalue.Null(),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("20"),
								"action": knownvalue.StringExact("deny"),
								"prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"network": knownvalue.Null(),
									"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"network":                knownvalue.StringExact("2001:db8::/32"),
										"greater_than_or_equal": knownvalue.Int64Exact(48),
										"less_than_or_equal":    knownvalue.Null(),
									}),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("30"),
								"action": knownvalue.StringExact("permit"),
								"prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"network": knownvalue.Null(),
									"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"network":                knownvalue.StringExact("fc00::/7"),
										"greater_than_or_equal": knownvalue.Null(),
										"less_than_or_equal":    knownvalue.Int64Exact(64),
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

const filtersPrefixListRoutingProfile_Ipv6_MultipleEntries_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_filters_prefix_list_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  type = {
    ipv6 = {
      ipv6_entries = [
        {
          name = "10"
          action = "permit"
          prefix = {
            network = "any"
          }
        },
        {
          name = "20"
          action = "deny"
          prefix = {
            entry = {
              network = "2001:db8::/32"
              greater_than_or_equal = 48
            }
          }
        },
        {
          name = "30"
          action = "permit"
          prefix = {
            entry = {
              network = "fc00::/7"
              less_than_or_equal = 64
            }
          }
        }
      ]
    }
  }
}
`
