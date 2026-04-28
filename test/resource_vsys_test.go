package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccPanosVsys_ConflictWithTemplateDefaultVsys verifies that creating a
// panos_vsys resource for "vsys1" fails when the template already has
// default_vsys = "vsys1" (because the template's PostCreate hook already
// created that vsys entry).
func TestAccPanosVsys_ConflictWithTemplateDefaultVsys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":   config.StringVariable(prefix),
		"location": location,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          vsys_ConflictWithTemplateDefaultVsys_Tmpl,
				ConfigVariables: configVars,
				ExpectError:     regexp.MustCompile(`.`),
			},
		},
	})
}

const vsys_ConflictWithTemplateDefaultVsys_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "test" {
  location = var.location
  name     = "${var.prefix}-template"
  default_vsys = "vsys1"
}

resource "panos_vsys" "vsys1" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }
  name = "vsys1"
}
`

// TestAccPanosVsys_WithTemplate verifies that creating a panos_vsys resource
// for "vsys2" succeeds when the template has default_vsys = "vsys1" (a
// different vsys name, so no conflict).
func TestAccPanosVsys_WithTemplate(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"panorama": config.ObjectVariable(map[string]config.Variable{}),
	})

	configVars := map[string]config.Variable{
		"prefix":   config.StringVariable(prefix),
		"location": location,
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:          vsys_WithTemplate_Tmpl,
				ConfigVariables: configVars,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_vsys.vsys2",
						tfjsonpath.New("name"),
						knownvalue.StringExact("vsys2"),
					),
				},
			},
		},
	})
}

const vsys_WithTemplate_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "test" {
  location = var.location
  name     = "${var.prefix}-template"
  default_vsys = "vsys1"
}

resource "panos_vsys" "vsys2" {
  location = {
    template = {
      name = panos_template.test.name
    }
  }
  name = "vsys2"
}
`
