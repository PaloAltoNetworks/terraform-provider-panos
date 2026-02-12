package provider_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCustomUrlCategory(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomUrlCategoryResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"type":   config.StringVariable("URL List"),
					"list":   config.ListVariable(config.StringVariable("example.com")),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.category",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-category", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.category",
						tfjsonpath.New("type"),
						knownvalue.StringExact("URL List"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.category",
						tfjsonpath.New("list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("example.com"),
						}),
					),
				},
			},
			{
				Config: testAccCustomUrlCategoryResourceTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"type":   config.StringVariable("Category Match"),
					"list":   config.ListVariable(config.StringVariable("unknown")),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.category",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-category", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.category",
						tfjsonpath.New("type"),
						knownvalue.StringExact("Category Match"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_url_category.category",
						tfjsonpath.New("list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("unknown"),
						}),
					),
				},
			},
		},
	})
}

const testAccCustomUrlCategoryResourceTmpl = `
variable prefix { type = string }
variable type { type = string }
variable list { type = list(string) }

resource "panos_custom_url_category" "category" {
  location = { shared = {} }

  name = format("%s-category", var.prefix)
  type = var.type
  list = var.list
}
`

func TestAccCustomUrlCategoryImportShared(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	name := fmt.Sprintf("test-acc-import-shared-%s", nameSuffix)

	importStateGenerateID := func(state *terraform.State) (string, error) {
		importState := map[string]any{
			"location": map[string]any{
				"shared": map[string]any{},
			},
			"name": name,
		}

		marshalled, err := json.Marshal(importState)
		if err != nil {
			return "", fmt.Errorf("Failed to marshal import state into JSON: %w", err)
		}

		return base64.StdEncoding.EncodeToString(marshalled), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomUrlCategoryImportSharedConfig,
				ConfigVariables: map[string]config.Variable{
					"name": config.StringVariable(name),
				},
			},
			{
				Config: testAccCustomUrlCategoryImportSharedStep,
				ConfigVariables: map[string]config.Variable{
					"name": config.StringVariable(name),
				},
				ResourceName:      "panos_custom_url_category.imported",
				ImportState:       true,
				ImportStateIdFunc: importStateGenerateID,
			},
		},
	})
}

const testAccCustomUrlCategoryImportSharedConfig = `
variable "name" { type = string }

resource "panos_custom_url_category" "test" {
  location = { shared = {} }
  name     = var.name
  type     = "URL List"
  list     = ["example.com", "test.com"]
  description = "Test category for import"
}
`

const testAccCustomUrlCategoryImportSharedStep = `
variable "name" { type = string }

resource "panos_custom_url_category" "imported" {
  location    = { shared = {} }
  name        = var.name
  type        = "URL List"
  list        = ["example.com", "test.com"]
  description = "Test category for import"
}
`

func TestAccCustomUrlCategoryImportDeviceGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	dgName := fmt.Sprintf("dg-%s", nameSuffix)
	name := fmt.Sprintf("cat-%s", nameSuffix)

	importStateGenerateID := func(state *terraform.State) (string, error) {
		importState := map[string]any{
			"location": map[string]any{
				"device_group": map[string]any{
					"name":            dgName,
					"panorama_device": "localhost.localdomain",
				},
			},
			"name": name,
		}

		marshalled, err := json.Marshal(importState)
		if err != nil {
			return "", fmt.Errorf("Failed to marshal import state into JSON: %w", err)
		}

		return base64.StdEncoding.EncodeToString(marshalled), nil
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCustomUrlCategoryImportDeviceGroupConfig,
				ConfigVariables: map[string]config.Variable{
					"dg_name": config.StringVariable(dgName),
					"name":    config.StringVariable(name),
				},
			},
			{
				Config: testAccCustomUrlCategoryImportDeviceGroupStep,
				ConfigVariables: map[string]config.Variable{
					"dg_name": config.StringVariable(dgName),
					"name":    config.StringVariable(name),
				},
				ResourceName:      "panos_custom_url_category.imported",
				ImportState:       true,
				ImportStateIdFunc: importStateGenerateID,
			},
		},
	})
}

const testAccCustomUrlCategoryImportDeviceGroupConfig = `
variable "dg_name" { type = string }
variable "name" { type = string }

resource "panos_device_group" "test" {
  location = { panorama = {} }
  name     = var.dg_name
}

resource "panos_custom_url_category" "test" {
  location = {
    device_group = {
      name            = panos_device_group.test.name
      panorama_device = "localhost.localdomain"
    }
  }
  name        = var.name
  type        = "URL List"
  list        = ["device-group.example.com"]
  description = "Test DG category for import"
}
`

const testAccCustomUrlCategoryImportDeviceGroupStep = `
variable "dg_name" { type = string }
variable "name" { type = string }

resource "panos_device_group" "test" {
  location = { panorama = {} }
  name     = var.dg_name
}

resource "panos_custom_url_category" "imported" {
  location = {
    device_group = {
      name            = panos_device_group.test.name
      panorama_device = "localhost.localdomain"
    }
  }
  name        = var.name
  type        = "URL List"
  list        = ["device-group.example.com"]
  description = "Test DG category for import"
}
`
