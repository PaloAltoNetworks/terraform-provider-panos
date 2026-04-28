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

func TestAccOspfAuthRoutingProfile_Password(t *testing.T) {
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
				Config: ospfAuthRoutingProfile_Password_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					// Password is hashed, so we can't check exact value
					// Just verify it's not null since we set it
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("password"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

const ospfAuthRoutingProfile_Password_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospf_auth_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  password = "test123"
}
`

func TestAccOspfAuthRoutingProfile_Md5(t *testing.T) {
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
				Config: ospfAuthRoutingProfile_Md5_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					// Check MD5 list has 2 entries
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5"),
						knownvalue.ListSizeExact(2),
					),
					// Check first MD5 entry name
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("1"),
					),
					// Key is hashed, so verify it's not null
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5").AtSliceIndex(0).AtMapKey("key"),
						knownvalue.NotNull(),
					),
					// Check preferred flag
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5").AtSliceIndex(0).AtMapKey("preferred"),
						knownvalue.Bool(true),
					),
					// Check second MD5 entry name
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("2"),
					),
					// Key is hashed, so verify it's not null
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5").AtSliceIndex(1).AtMapKey("key"),
						knownvalue.NotNull(),
					),
					// Check preferred flag (should be false/null for second entry)
					statecheck.ExpectKnownValue(
						"panos_ospf_auth_routing_profile.example",
						tfjsonpath.New("md5").AtSliceIndex(1).AtMapKey("preferred"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const ospfAuthRoutingProfile_Md5_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ospf_auth_routing_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  md5 = [
    {
      name = "1"
      key = "mykey123"
      preferred = true
    },
    {
      name = "2"
      key = "anotherkey456"
    }
  ]
}
`
