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

func TestAccDefaultSecurityPolicy(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: defaultSecurityPolicyTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_default_security_policy.defaults",
						tfjsonpath.New("rules"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":             knownvalue.StringExact("interzone-default"),
								"action":           knownvalue.StringExact("drop"),
								"group_tag":        knownvalue.Null(),
								"icmp_unreachable": knownvalue.Null(),
								"log_end":          knownvalue.Null(),
								"log_setting":      knownvalue.Null(),
								"log_start":        knownvalue.Null(),
								"profile_setting":  knownvalue.Null(),
								"tag":              knownvalue.Null(),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":             knownvalue.StringExact("intrazone-default"),
								"action":           knownvalue.StringExact("drop"),
								"group_tag":        knownvalue.Null(),
								"icmp_unreachable": knownvalue.Null(),
								"log_end":          knownvalue.Null(),
								"log_setting":      knownvalue.Null(),
								"log_start":        knownvalue.Null(),
								"profile_setting":  knownvalue.Null(),
								"tag":              knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const defaultSecurityPolicyTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_default_security_policy" "defaults" {
  location = { device_group = { name = panos_device_group.dg.name } }
  rules = [
    {
      name = "interzone-default"
      action = "drop"
    },
    {
      name = "intrazone-default"
      action = "drop"
    },
  ]
}
`
