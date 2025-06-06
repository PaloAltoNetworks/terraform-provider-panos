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
