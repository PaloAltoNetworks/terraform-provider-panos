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

func TestAccLogForwarding(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: logForwardingResource1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("enhanced_application_logging"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-match1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("action_desc"),
						knownvalue.StringExact("action description"),
					),
					// match_list[0].actions[0]
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(0).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-action1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(0).
							AtMapKey("type").
							AtMapKey("integration").
							AtMapKey("action"),
						knownvalue.StringExact("Azure-Security-Center-Integration"),
					),
					// match_list[0].actions[1]
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(1).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-action2", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(1).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("action"),
						knownvalue.StringExact("add-tag"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(1).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("registration").
							AtMapKey("localhost"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(1).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("tags"),
						knownvalue.ListSizeExact(0),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(1).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("target"),
						knownvalue.StringExact("source-address"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(1).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("timeout"),
						knownvalue.Int64Exact(3600),
					),
					// match_list[0].actions[2]
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(2).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-action3", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(2).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("action"),
						knownvalue.StringExact("remove-tag"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(2).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("registration").
							AtMapKey("panorama"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(2).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("target"),
						knownvalue.StringExact("user"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(2).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tag1", prefix)),
							knownvalue.StringExact(fmt.Sprintf("%s-tag2", prefix)),
						}),
					),
					// match_list[0].actions[3]
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(3).
							AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-action4", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(3).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("action"),
						knownvalue.StringExact("add-tag"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(3).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("target"),
						knownvalue.StringExact("source-address"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_forwarding_profile.profile",
						tfjsonpath.New("match_list").
							AtSliceIndex(0).
							AtMapKey("actions").
							AtSliceIndex(3).
							AtMapKey("type").
							AtMapKey("tagging").
							AtMapKey("registration").
							AtMapKey("remote").
							AtMapKey("http_profile"),
						knownvalue.StringExact("http-profile"),
					),
				},
			},
		},
	})
}

const logForwardingResource1 = `
variable "prefix" { type = string }

locals {
  device_group = format("%s-dg", var.prefix)
}

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = local.device_group
}

resource "panos_log_forwarding_profile" "profile" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-profile", var.prefix)

  disable_override = "no"

  enhanced_application_logging = true

  match_list = [{
    name = format("%s-match1", var.prefix)

    action_desc = "action description"
    actions     = [
      {
        name = format("%s-action1", var.prefix)
        type = { integration = { action = "Azure-Security-Center-Integration" } }
      },
      {
        name = format("%s-action2", var.prefix)
        type = {
          tagging = {
            action = "add-tag"
            registration = { localhost = {} }
            tags = []
            target = "source-address"
            timeout = 3600
          }
        }
      },
      {
        name = format("%s-action3", var.prefix)
        type = {
          tagging = {
            action = "remove-tag"
            registration = { panorama = {} }
            target = "user"
            tags = [format("%s-tag1", var.prefix), format("%s-tag2", var.prefix)]
          }
        }
      },
      {
        name = format("%s-action4", var.prefix)
        type = {
          tagging = {
            action = "add-tag"
            target = "source-address"
            registration = { remote = { http_profile = "http-profile" } }
          }
        }
      },
    ]
    filter = "All Logs"
    log_type = "traffic"
    quarantine = true
    send_email = ["test@example.com"]
    send_http = ["http://example.com"]
    send_snmptrap = ["trap.example.com"]
    send_syslog = ["syslog.example.com"]
    send_to_panorama = true
  }]
}
`
