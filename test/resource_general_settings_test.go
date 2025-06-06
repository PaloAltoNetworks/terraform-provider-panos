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

const generalSettings_Location_Tmpl = `
variable "create_template" { type = any }
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  count = var.create_template ? 1 : 0
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}


resource "panos_general_settings" "example" {
  depends_on = [panos_template.example]

  location = var.location

  hostname = var.prefix
}
`

func TestAccGeneralSettings_Location_Template(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	createTemplate := config.BoolVariable(true)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: generalSettings_Location_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"create_template": createTemplate,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_general_settings.example",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

func TestAccGeneralSettings_Location_System(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"system": config.ObjectVariable(map[string]config.Variable{}),
	})

	createTemplate := config.BoolVariable(false)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: generalSettings_Location_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"create_template": createTemplate,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_general_settings.example",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const generalSettings_Sanity_Delete_Initial_Tmpl = `
variable "create_template" { type = bool }
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  count = var.create_template ? 1 : 0
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}


resource "panos_general_settings" "example" {
  depends_on = [panos_template.example]

  location = var.location

  hostname = var.prefix
}

resource "panos_ntp_settings" "example" {
  depends_on = [panos_template.example]

  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = { none = {} }
    }
  }
}
`

const generalSettings_Sanity_Delete_Updated_Tmpl = `
variable "create_template" { type = bool }
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  count = var.create_template ? 1 : 0
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

resource "panos_ntp_settings" "example" {
  depends_on = [panos_template.example]

  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = { none = {} }
    }
  }
}
`

const generalSettings_Sanity_Delete_Final_Tmpl = `
variable "create_template" { type = bool }
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  count = var.create_template ? 1 : 0
  location = { panorama = {} }

  name = format("%s-tmpl", var.prefix)
}

data "panos_general_settings" "example" {
  location = var.location
}

resource "panos_ntp_settings" "example" {
  depends_on = [panos_template.example]

  location = var.location

  ntp_servers = {
    primary_ntp_server = {
      ntp_server_address = "172.16.0.1"
      authentication_type = { none = {} }
    }
  }
}
`

func TestAccGeneralSettings_Sanity_Delete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"ngfw_device":     config.StringVariable("localhost.localdomain"),
			"panorama_device": config.StringVariable("localhost.localdomain"),
			"name":            config.StringVariable(prefix),
		}),
	})

	createTemplate := config.BoolVariable(true)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: generalSettings_Sanity_Delete_Initial_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"create_template": createTemplate,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_general_settings.example",
						tfjsonpath.New("hostname"),
						knownvalue.StringExact(prefix),
					),
				},
			},
			{
				Config: generalSettings_Sanity_Delete_Updated_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"create_template": createTemplate,
				},
			},
			{
				Config: generalSettings_Sanity_Delete_Final_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":          config.StringVariable(prefix),
					"location":        location,
					"create_template": createTemplate,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.panos_general_settings.example",
						tfjsonpath.New("hostname"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}
