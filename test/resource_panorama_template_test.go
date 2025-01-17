package provider_test

import (
	"context"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/panorama/template"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosTemplate_RequiredInputs(t *testing.T) {
	t.Parallel()

	resourceName := "acc_test_template"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("%s-%s", resourceName, nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckPanoramaTemplateDestroy(templateName),
		Steps: []resource.TestStep{
			{
				Config: makePanosTemplateConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_template."+resourceName,
						tfjsonpath.New("name"),
						knownvalue.StringExact(templateName),
					),
				},
			},
		},
	})
}

func makePanosTemplateConfig(label string) string {
	configTpl := `
    variable "template_name" { type = string }

    resource "panos_template" "%s" {
        name = var.template_name

        location = {
            panorama = {
                panorama_device = "localhost.localdomain"
            }
        }
    }
    `
	return fmt.Sprintf(configTpl, label)
}

func testAccCheckPanoramaTemplateDestroy(name string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := template.NewService(sdkClient)
		location := template.NewPanoramaLocation()
		ctx := context.TODO()

		reply, err := api.Read(ctx, *location, name, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading template entry via sdk: %v", err)
		}

		if reply != nil {
			if reply.EntryName() == name {
				return fmt.Errorf("template object still exists: %s", name)
			}
		}

		return nil
	}
}
