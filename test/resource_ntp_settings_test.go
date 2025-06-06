package provider_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccNtpSettings(t *testing.T) {
	location := config.ObjectVariable(map[string]config.Variable{
		"system": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ntpSettingsConfig1,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").AtMapKey("primary_ntp_server").AtMapKey("ntp_server_address"),
						knownvalue.StringExact("172.16.0.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").AtMapKey("secondary_ntp_server").AtMapKey("ntp_server_address"),
						knownvalue.StringExact("172.16.0.2"),
					),
				},
			},
			{
				Config: ntpSettingsConfig2,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").AtMapKey("primary_ntp_server").AtMapKey("authentication_type").AtMapKey("autokey"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
				},
			},
			{
				Config: ntpSettingsConfig3,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").
							AtMapKey("primary_ntp_server").
							AtMapKey("authentication_type").
							AtMapKey("symmetric_key").
							AtMapKey("key_id"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").
							AtMapKey("primary_ntp_server").
							AtMapKey("authentication_type").
							AtMapKey("symmetric_key").
							AtMapKey("algorithm").
							AtMapKey("sha1").
							AtMapKey("authentication_key"),
						knownvalue.StringExact("da39a3ee5e6b4b0d3255bfef95601890afd80709"),
					),
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").
							AtMapKey("secondary_ntp_server").
							AtMapKey("authentication_type").
							AtMapKey("symmetric_key").
							AtMapKey("key_id"),
						knownvalue.Int64Exact(1),
					),
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").
							AtMapKey("secondary_ntp_server").
							AtMapKey("authentication_type").
							AtMapKey("symmetric_key").
							AtMapKey("algorithm").
							AtMapKey("md5").
							AtMapKey("authentication_key"),
						knownvalue.StringExact("d41d8cd98f00b204e9800998ecf8427e"),
					),
				},
			},
			{
				Config: ntpSettingsConfig4,
				ConfigVariables: map[string]config.Variable{
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ntp_settings.settings",
						tfjsonpath.New("ntp_servers").
							AtMapKey("secondary_ntp_server").
							AtMapKey("authentication_type").
							AtMapKey("symmetric_key").
							AtMapKey("algorithm").
							AtMapKey("md5").
							AtMapKey("authentication_key"),
						knownvalue.StringExact("83d043db4fdfe6882fb7f01a09d92b11"),
					),
				},
			},
		},
	})
}

const ntpSettingsConfig1 = `
variable "location" { type = map }

resource "panos_ntp_settings" "settings" {
  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = { none = {}}
    }
    secondary_ntp_server = {
      ntp_server_address = "172.16.0.2"
      authentication_type = { none = {}}
    }
  }
}
`

const ntpSettingsConfig2 = `
variable "location" { type = map }

resource "panos_ntp_settings" "settings" {
  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = {
        autokey = {}
      }
    }
    secondary_ntp_server = {
      ntp_server_address = "172.16.0.2"
    }
  }
}
`

const ntpSettingsConfig3 = `
variable "location" { type = map }

resource "panos_ntp_settings" "settings" {
  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = {
        symmetric_key = {
          key_id = 1
          algorithm = {
            sha1 = {
              authentication_key = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
            }
          }
        }
      }
    }

    secondary_ntp_server = {
      ntp_server_address = "172.16.0.2"
      authentication_type = {
        symmetric_key = {
          key_id = 1
          algorithm = {
            md5 = {
              authentication_key = "d41d8cd98f00b204e9800998ecf8427e"
            }
          }
        }
      }
    }
  }
}
`

const ntpSettingsConfig4 = `
variable "location" { type = map }

resource "panos_ntp_settings" "settings" {
  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = {
        symmetric_key = {
          key_id = 1
          algorithm = {
            sha1 = {
              authentication_key = "da39a3ee5e6b4b0d3255bfef95601890afd80709"
            }
          }
        }
      }
    }

    secondary_ntp_server = {
      ntp_server_address = "172.16.0.2"
      authentication_type = {
        symmetric_key = {
          key_id = 1
          algorithm = {
            md5 = {
              authentication_key = "83d043db4fdfe6882fb7f01a09d92b11"
            }
          }
        }
      }
    }
  }
}
`
