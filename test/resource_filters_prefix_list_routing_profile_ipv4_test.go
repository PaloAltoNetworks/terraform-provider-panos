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

func TestAccFiltersPrefixListRoutingProfile_Ipv4_Basic(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv4_Basic_Tmpl,
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
						knownvalue.StringExact("Test prefix list description"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ipv4": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ipv4_entries": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":   knownvalue.StringExact("10"),
										"action": knownvalue.StringExact("deny"),
										"prefix": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"network": knownvalue.Null(),
											"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"network":                knownvalue.StringExact("10.0.0.0/8"),
												"greater_than_or_equal": knownvalue.Null(),
												"less_than_or_equal":    knownvalue.Null(),
											}),
										}),
									}),
								}),
							}),
							"ipv6": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const filtersPrefixListRoutingProfile_Ipv4_Basic_Tmpl = `
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
  description = "Test prefix list description"

  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          prefix = {
            entry = {
              network = "10.0.0.0/8"
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersPrefixListRoutingProfile_Ipv4_Prefix_Network(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv4_Prefix_Network_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("prefix"),
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

const filtersPrefixListRoutingProfile_Ipv4_Prefix_Network_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
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

func TestAccFiltersPrefixListRoutingProfile_Ipv4_Prefix_Entry(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv4_Prefix_Entry_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("prefix"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"network": knownvalue.Null(),
							"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"network":                knownvalue.StringExact("192.168.0.0/16"),
								"greater_than_or_equal": knownvalue.Int64Exact(24),
								"less_than_or_equal":    knownvalue.Int64Exact(28),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersPrefixListRoutingProfile_Ipv4_Prefix_Entry_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          prefix = {
            entry = {
              network = "192.168.0.0/16"
              greater_than_or_equal = 24
              less_than_or_equal = 28
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersPrefixListRoutingProfile_Ipv4_Action_Permit(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv4_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersPrefixListRoutingProfile_Ipv4_Action_Permit_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
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

func TestAccFiltersPrefixListRoutingProfile_Ipv4_MultipleEntries(t *testing.T) {
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
				Config: filtersPrefixListRoutingProfile_Ipv4_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_prefix_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries"),
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
										"network":                knownvalue.StringExact("10.0.0.0/8"),
										"greater_than_or_equal": knownvalue.Int64Exact(16),
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
										"network":                knownvalue.StringExact("172.16.0.0/12"),
										"greater_than_or_equal": knownvalue.Null(),
										"less_than_or_equal":    knownvalue.Int64Exact(24),
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

const filtersPrefixListRoutingProfile_Ipv4_MultipleEntries_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
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
              network = "10.0.0.0/8"
              greater_than_or_equal = 16
            }
          }
        },
        {
          name = "30"
          action = "permit"
          prefix = {
            entry = {
              network = "172.16.0.0/12"
              less_than_or_equal = 24
            }
          }
        }
      ]
    }
  }
}
`

