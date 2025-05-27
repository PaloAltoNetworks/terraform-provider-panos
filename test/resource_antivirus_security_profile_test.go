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

func TestAccAntivirusSecurityProfile(t *testing.T) {
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
				Config: panosAntivirusSecurityProfileTmpl1,
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

const panosAntivirusSecurityProfileTmpl1 = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_device_group" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_antivirus_security_profile" "example" {
  location = var.location

  name = var.prefix
  description = "Example antivirus security profile"
  disable_override = "no"

  #application_exceptions = [{
  #  name   = "app_exception_1"
  #  action = "alert"
  #}]

  #decoders = [{
  #  name            = "decoder_1"
  #  action          = "drop"
  #  wildfire_action = "alert"
  #  ml_action       = "reset-client"
  #}]

  #machine_learning_models = [{
  #  name   = "ml_model_1"
  #  action = "enable"
  #}]

  #machine_learning_exceptions = [{
  #  name        = "ml_exception_1"
  #  filename    = "example.exe"
  #  description = "Example ML exception"
  #}]

  packet_capture = true

  #threat_exceptions = [{
  #  name = "threat_exception_1"
  #}]

  wfrt_hold_mode = false
}
`
