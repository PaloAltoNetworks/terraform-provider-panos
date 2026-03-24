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

func TestAccFiltersAccessListRoutingProfile_Ipv4_Basic(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_Basic_Tmpl,
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
						knownvalue.StringExact("IPv4 ACL test"),
					),
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("10"),
								"action": knownvalue.StringExact("deny"),
								"source_address": knownvalue.Null(),
								"destination_address": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv4_Basic_Tmpl = `
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
  description = "IPv4 ACL test"

  type = {
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv4_SourceAddress_Any(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_SourceAddress_Any_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("source_address"),
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

const filtersAccessListRoutingProfile_Ipv4_SourceAddress_Any_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
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


func TestAccFiltersAccessListRoutingProfile_Ipv4_SourceAddress_Entry(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_SourceAddress_Entry_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("source_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"address": knownvalue.Null(),
							"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"address":  knownvalue.StringExact("192.168.1.0"),
								"wildcard": knownvalue.StringExact("0.0.0.255"),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv4_SourceAddress_Entry_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          source_address = {
            entry = {
              address = "192.168.1.0"
              wildcard = "0.0.0.255"
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv4_DestinationAddress_Any(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_DestinationAddress_Any_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("destination_address"),
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

const filtersAccessListRoutingProfile_Ipv4_DestinationAddress_Any_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          destination_address = {
            address = "any"
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv4_DestinationAddress_Entry(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_DestinationAddress_Entry_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("destination_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"address": knownvalue.Null(),
							"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"address":  knownvalue.StringExact("10.0.0.0"),
								"wildcard": knownvalue.StringExact("0.255.255.255"),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv4_DestinationAddress_Entry_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "deny"
          destination_address = {
            entry = {
              address = "10.0.0.0"
              wildcard = "0.255.255.255"
            }
          }
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv4_Action_Permit(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_Action_Permit_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries").AtSliceIndex(0).AtMapKey("action"),
						knownvalue.StringExact("permit"),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv4_Action_Permit_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
        }
      ]
    }
  }
}
`

func TestAccFiltersAccessListRoutingProfile_Ipv4_MultipleEntries(t *testing.T) {
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
				Config: filtersAccessListRoutingProfile_Ipv4_MultipleEntries_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_filters_access_list_routing_profile.example",
						tfjsonpath.New("type").AtMapKey("ipv4").AtMapKey("ipv4_entries"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("10"),
								"action": knownvalue.StringExact("permit"),
								"source_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.StringExact("any"),
									"entry":   knownvalue.Null(),
								}),
								"destination_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.Null(),
									"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"address":  knownvalue.StringExact("192.168.0.0"),
										"wildcard": knownvalue.StringExact("0.0.255.255"),
									}),
								}),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("20"),
								"action": knownvalue.StringExact("deny"),
								"source_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.Null(),
									"entry": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"address":  knownvalue.StringExact("10.0.0.0"),
										"wildcard": knownvalue.StringExact("0.255.255.255"),
									}),
								}),
								"destination_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"address": knownvalue.StringExact("any"),
									"entry":   knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const filtersAccessListRoutingProfile_Ipv4_MultipleEntries_Tmpl = `
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
    ipv4 = {
      ipv4_entries = [
        {
          name = "10"
          action = "permit"
          source_address = {
            address = "any"
          }
          destination_address = {
            entry = {
              address = "192.168.0.0"
              wildcard = "0.0.255.255"
            }
          }
        },
        {
          name = "20"
          action = "deny"
          source_address = {
            entry = {
              address = "10.0.0.0"
              wildcard = "0.255.255.255"
            }
          }
          destination_address = {
            address = "any"
          }
        }
      ]
    }
  }
}
`
