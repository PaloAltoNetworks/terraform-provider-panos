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

func TestAccOspfv3AuthRoutingProfile_Basic(t *testing.T) {
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
				Config: ospfv3AuthRoutingProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("spi"),
						knownvalue.StringExact("12345678"),
					),
					// Verify AH variant is set
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("ah"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"sha256": knownvalue.ObjectExact(map[string]knownvalue.Check{
								// Don't check key value - it will be hashed/encrypted
								"key": knownvalue.NotNull(),
							}),
							"md5":    knownvalue.Null(),
							"sha1":   knownvalue.Null(),
							"sha384": knownvalue.Null(),
							"sha512": knownvalue.Null(),
						}),
					),
					// Verify ESP variant is null (mutually exclusive with AH)
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("esp"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const ospfv3AuthRoutingProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_auth_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  spi = "12345678"
  ah = {
    sha256 = {
      key = "aaaaaaaa-bbbbbbbb-cccccccc-dddddddd-eeeeeeee-ffffffff-99999999-88888888"
    }
  }
}
`

func TestAccOspfv3AuthRoutingProfile_Ah_Md5(t *testing.T) {
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
				Config: ospfv3AuthRoutingProfile_Ah_Md5_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("spi"),
						knownvalue.StringExact("abcdef00"),
					),
					// Verify AH variant with MD5
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("ah"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"md5": knownvalue.ObjectExact(map[string]knownvalue.Check{
								// Don't check key value - it will be hashed/encrypted
								"key": knownvalue.NotNull(),
							}),
							"sha1":   knownvalue.Null(),
							"sha256": knownvalue.Null(),
							"sha384": knownvalue.Null(),
							"sha512": knownvalue.Null(),
						}),
					),
					// Verify ESP variant is null
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("esp"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const ospfv3AuthRoutingProfile_Ah_Md5_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_auth_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  spi = "abcdef00"
  ah = {
    md5 = {
      key = "11111111-22222222-33333333-44444444"
    }
  }
}
`

func TestAccOspfv3AuthRoutingProfile_Esp_AuthAndEncryption(t *testing.T) {
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
				Config: ospfv3AuthRoutingProfile_Esp_AuthAndEncryption_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("spi"),
						knownvalue.StringExact("aabbccdd"),
					),
					// Verify AH variant is null (mutually exclusive with ESP)
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("ah"),
						knownvalue.Null(),
					),
					// Verify ESP variant is set with authentication and encryption
					statecheck.ExpectKnownValue(
						"panos_ospfv3_auth_routing_profile.example",
						tfjsonpath.New("esp"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"authentication": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"sha256": knownvalue.ObjectExact(map[string]knownvalue.Check{
									// Don't check key value - it will be hashed/encrypted
									"key": knownvalue.NotNull(),
								}),
								"md5":    knownvalue.Null(),
								"none":   knownvalue.Null(),
								"sha1":   knownvalue.Null(),
								"sha384": knownvalue.Null(),
								"sha512": knownvalue.Null(),
							}),
							"encryption": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"algorithm": knownvalue.StringExact("aes-256-cbc"),
								// Don't check key value - it will be hashed/encrypted
								"key": knownvalue.NotNull(),
							}),
						}),
					),
				},
			},
		},
	})
}

const ospfv3AuthRoutingProfile_Esp_AuthAndEncryption_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospfv3_auth_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  spi = "aabbccdd"
  esp = {
    authentication = {
      sha256 = {
        key = "aaaaaaaa-bbbbbbbb-cccccccc-dddddddd-eeeeeeee-ffffffff-99999999-88888888"
      }
    }
    encryption = {
      algorithm = "aes-256-cbc"
      key = "11111111-22222222-33333333-44444444-55555555-66666666-77777777-88888888"
    }
  }
}
`
