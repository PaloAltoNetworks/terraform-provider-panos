package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	category "github.com/PaloAltoNetworks/pango/objects/profiles/customurlcategory"

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
		CheckDestroy:             testAccCustomUrlCategoryDestroy(prefix),
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
  location = { shared = true }

  name = format("%s-category", var.prefix)
  type = var.type
  list = var.list
}
`

func testAccCustomUrlCategoryDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := category.NewService(sdkClient)
		ctx := context.TODO()

		location := category.NewSharedLocation()

		entries, err := api.List(ctx, *location, "get", "", "")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("listing interface management entries via sdk: %v", err)
		}

		var leftEntries []string
		for _, elt := range entries {
			if strings.HasPrefix(elt.Name, prefix) {
				leftEntries = append(leftEntries, elt.Name)
			}
		}

		if len(leftEntries) > 0 {
			err := fmt.Errorf("terraform failed to remove entries from the server")
			delErr := api.Delete(ctx, *location, leftEntries...)
			if delErr != nil {
				return errors.Join(err, delErr)
			}
		}

		return nil
	}
}
