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

func TestAccDynamicUserGroup_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUserGroup_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Test dynamic user group"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("'tag1' or 'tag2'"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("tag1"),
							knownvalue.StringExact("tag2"),
						}),
					),
				},
			},
		},
	})
}

func TestAccDynamicUserGroup_DisableOverride(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUserGroup_DisableOverride_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
				},
			},
		},
	})
}

func TestAccDynamicUserGroup_Minimal(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUserGroup_Minimal_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("filter"),
						knownvalue.StringExact("'minimal'"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("description"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_user_group.example",
						tfjsonpath.New("tags"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const dynamicUserGroup_Minimal_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_dynamic_user_group" "example" {
  location = var.location

  name = var.prefix
  filter = "'minimal'"
}
`

const dynamicUserGroup_DisableOverride_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_dynamic_user_group" "example" {
  location = var.location

  name = var.prefix
  filter = "'tag1'"
  disable_override = "yes"
}
`

const dynamicUserGroup_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_administrative_tag" "tag1" {
  location = var.location
  name = "tag1"
}

resource "panos_administrative_tag" "tag2" {
  location = var.location
  name = "tag2"
}

resource "panos_dynamic_user_group" "example" {
  depends_on = [panos_administrative_tag.tag1, panos_administrative_tag.tag2]
  location = var.location

  name = var.prefix
  description = "Test dynamic user group"
  filter = "'tag1' or 'tag2'"
  tags = [panos_administrative_tag.tag1.name, panos_administrative_tag.tag2.name]
}
`
