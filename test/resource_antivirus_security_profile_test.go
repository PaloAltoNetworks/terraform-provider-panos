package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAntivirusSecurityProfile_Basic(t *testing.T) {
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
				Config: antivirusSecurityProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Example antivirus security profile"),
					),
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("packet_capture"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("wfrt_hold_mode"),
						knownvalue.Bool(false),
					),
				},
			},
		},
	})
}

const antivirusSecurityProfile_Basic_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  description = "Example antivirus security profile"
  disable_override = "no"
  packet_capture = true
  wfrt_hold_mode = false
}
`

const antivirusSecurityProfile_ApplicationExceptions_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  application_exceptions = [{
    name   = "panos-web-interface"
    action = "alert"
  }]
}
`

func TestAccAntivirusSecurityProfile_ApplicationExceptions(t *testing.T) {
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
				Config: antivirusSecurityProfile_ApplicationExceptions_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("application_exceptions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("panos-web-interface"),
								"action": knownvalue.StringExact("alert"),
							}),
						}),
					),
				},
			},
		},
	})
}

const antivirusSecurityProfile_Decoders_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  decoders = [{
    name = "http"
    action = "drop"
    wildfire_action = "alert"
    ml_action = "reset-client"
  }]
}
`

func TestAccAntivirusSecurityProfile_Decoders(t *testing.T) {
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
				Config: antivirusSecurityProfile_Decoders_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("decoders"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("http"),
								"action": knownvalue.StringExact("drop"),
								"wildfire_action": knownvalue.StringExact("alert"),
								"ml_action": knownvalue.StringExact("reset-client"),
							}),
						}),
					),
				},
			},
		},
	})
}

const antivirusSecurityProfile_MachineLearningModels_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  machine_learning_models = [
    {
      name = "Windows Executables"
      action = "enable(alert-only)"
    },
    {
      name = "PowerShell Script 2"
      action = "disable"
    },
    {
      name = "Executable Linked Format"
      action = "enable"
    }
  ]
}
`

func TestAccAntivirusSecurityProfile_MachineLearningModels(t *testing.T) {
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
				Config: antivirusSecurityProfile_MachineLearningModels_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("machine_learning_models"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("Windows Executables"),
								"action": knownvalue.StringExact("enable(alert-only)"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("PowerShell Script 2"),
								"action": knownvalue.StringExact("disable"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":   knownvalue.StringExact("Executable Linked Format"),
								"action": knownvalue.StringExact("enable"),
							}),
						}),
					),
				},
			},
		},
	})
}

const antivirusSecurityProfile_MachineLearningExceptions_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  machine_learning_exceptions = [{
    name = "ml_exception_1"
    filename = "example.exe"
    description = "Example ML exception"
  }]
}
`

func TestAccAntivirusSecurityProfile_MachineLearningExceptions(t *testing.T) {
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
				Config: antivirusSecurityProfile_MachineLearningExceptions_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("machine_learning_exceptions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("ml_exception_1"),
								"filename": knownvalue.StringExact("example.exe"),
								"description": knownvalue.StringExact("Example ML exception"),
							}),
						}),
					),
				},
			},
		},
	})
}

const antivirusSecurityProfile_ThreatExceptions_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  threat_exceptions = [{
    name = "20036500"
  }]
}
`

func TestAccAntivirusSecurityProfile_ThreatExceptions(t *testing.T) {
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
				Config: antivirusSecurityProfile_ThreatExceptions_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_antivirus_security_profile.example",
						tfjsonpath.New("threat_exceptions"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("20036500"),
							}),
						}),
					),
				},
			},
		},
	})
}

