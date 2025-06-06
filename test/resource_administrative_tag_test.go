package provider_test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAdministrativeTag(t *testing.T) {
	resourceName := "test_tag"
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"shared": config.ObjectVariable(map[string]config.Variable{}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makeAdministrativeTagConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"location": location,
					"tag_name": config.StringVariable(fmt.Sprintf("%s-tag1-nocolor", prefix)),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						fmt.Sprintf("panos_administrative_tag.%s", resourceName),
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tag1-nocolor", prefix)),
					),
				},
			},
		},
	})

	colorValue := "color1"
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: makeAdministrativeTagConfig(resourceName),
				ConfigVariables: map[string]config.Variable{
					"location": location,
					"tag_name": config.StringVariable(fmt.Sprintf("%s-tag1-color", prefix)),
					"color":    config.StringVariable(colorValue),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						fmt.Sprintf("panos_administrative_tag.%s", resourceName),
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tag1-color", prefix)),
					),
					statecheck.ExpectKnownValue(
						fmt.Sprintf("panos_administrative_tag.%s", resourceName),
						tfjsonpath.New("color"),
						knownvalue.StringExact(colorValue)),
				},
			},
		},
	})
}

const resourceTmpl = `
variable "location" { type = map }
variable "tag_name" { type = string }
variable "color" {
  type = string
  default = null
}

resource "panos_administrative_tag" "{{ .ResourceName }}" {
  location = var.location

  name  = var.tag_name
  color = var.color
}
`

func makeAdministrativeTagConfig(resourceName string) string {
	var buf bytes.Buffer
	tmpl := template.Must(template.New("").Parse(resourceTmpl))

	context := struct {
		ResourceName string
	}{
		ResourceName: resourceName,
	}

	err := tmpl.Execute(&buf, context)
	if err != nil {
		panic(err)
	}

	return buf.String()
}
