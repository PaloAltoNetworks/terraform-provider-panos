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

func TestAccTacacsPlusProfile_Basic(t *testing.T) {
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
				Config: tacacsPlusProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("CHAP"),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.StringExact("server1"),
								"address": knownvalue.StringExact("192.168.1.10"),
								"secret":  knownvalue.StringExact("secret123"),
								"port":    knownvalue.Int64Exact(49),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("timeout"),
						knownvalue.Int64Exact(5),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("use_single_connection"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestAccTacacsPlusProfile_Protocol_PAP(t *testing.T) {
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
				Config: tacacsPlusProfile_Protocol_PAP_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.StringExact("PAP"),
					),
				},
			},
		},
	})
}

const tacacsPlusProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_tacacs_plus_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name     = var.prefix
  protocol = "CHAP"

  servers = [
    {
      name    = "server1"
      address = "192.168.1.10"
      secret  = "secret123"
      port    = 49
    }
  ]

  timeout               = 5
  use_single_connection = true
}
`

func TestAccTacacsPlusProfile_MultipleServers(t *testing.T) {
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
				Config: tacacsPlusProfile_MultipleServers_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.StringExact("server1"),
								"address": knownvalue.StringExact("192.168.1.10"),
								"secret":  knownvalue.StringExact("secret123"),
								"port":    knownvalue.Int64Exact(49),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.StringExact("server2"),
								"address": knownvalue.StringExact("192.168.1.11"),
								"secret":  knownvalue.StringExact("secret456"),
								"port":    knownvalue.Int64Exact(50),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.StringExact("server3"),
								"address": knownvalue.StringExact("tacacs.example.com"),
								"secret":  knownvalue.StringExact("secret789"),
								"port":    knownvalue.Int64Exact(49),
							}),
						}),
					),
				},
			},
		},
	})
}

const tacacsPlusProfile_Protocol_PAP_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_tacacs_plus_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name     = var.prefix
  protocol = "PAP"

  servers = [
    {
      name    = "server1"
      address = "192.168.1.10"
      secret  = "secret123"
    }
  ]
}
`

func TestAccTacacsPlusProfile_ServerDefaults(t *testing.T) {
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
				Config: tacacsPlusProfile_ServerDefaults_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":    knownvalue.StringExact("server1"),
								"address": knownvalue.StringExact("192.168.1.10"),
								"secret":  knownvalue.StringExact("secret123"),
								"port":    knownvalue.Int64Exact(49),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_tacacs_plus_profile.example",
						tfjsonpath.New("timeout"),
						knownvalue.Int64Exact(3),
					),
				},
			},
		},
	})
}

const tacacsPlusProfile_MultipleServers_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_tacacs_plus_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  servers = [
    {
      name    = "server1"
      address = "192.168.1.10"
      secret  = "secret123"
      port    = 49
    },
    {
      name    = "server2"
      address = "192.168.1.11"
      secret  = "secret456"
      port    = 50
    },
    {
      name    = "server3"
      address = "tacacs.example.com"
      secret  = "secret789"
      port    = 49
    }
  ]
}
`

const tacacsPlusProfile_ServerDefaults_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_tacacs_plus_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  servers = [
    {
      name    = "server1"
      address = "192.168.1.10"
      secret  = "secret123"
    }
  ]
}
`
