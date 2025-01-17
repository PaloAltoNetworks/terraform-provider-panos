package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	servicegroup "github.com/PaloAltoNetworks/pango/objects/service/group"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosServiceGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckPanosServiceGroupDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: testAccPanosServiceGroupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
					"groups": config.MapVariable(map[string]config.Variable{
						"group1": config.ObjectVariable(map[string]config.Variable{
							"tags":    config.ListVariable(config.StringVariable(fmt.Sprintf("%s-tag", prefix))),
							"members": config.ListVariable(config.StringVariable(fmt.Sprintf("%s-svc", prefix))),
						}),
						"group2": config.ObjectVariable(map[string]config.Variable{
							"members": config.ListVariable(),
						}),
					}),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_service_group.group1",
						tfjsonpath.New("tags"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-tag", prefix)),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_service_group.group1",
						tfjsonpath.New("members"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-svc", prefix)),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_service_group.group2",
						tfjsonpath.New("members"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const testAccPanosServiceGroupTmpl = `
variable "prefix" { type = string }
variable "groups" {
  type = map(object({
    tags = optional(list(string)),
    members = optional(list(string)),
  }))
}

resource "panos_service" "svc" {
  location = { shared = true }

  name = format("%s-svc", var.prefix)
  protocol = { tcp = { source_port = 80, destination_port = 443 }}
}

resource "panos_administrative_tag" "tag" {
  location = { shared = true }

  name = format("%s-tag", var.prefix)
}

resource "panos_service_group" "group1" {
  depends_on = [
    resource.panos_service.svc,
    resource.panos_administrative_tag.tag
  ]
  location = { shared = true }

  name = format("%s-group1", var.prefix)
  members = var.groups["group1"].members
  tags = var.groups["group1"].tags
}

resource "panos_service_group" "group2" {
  location = { shared = true }

  name = format("%s-group2", var.prefix)
  members = var.groups["group2"].members
}
`

func testAccCheckPanosServiceGroupDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		svc := servicegroup.NewService(sdkClient)
		ctx := context.TODO()
		location := servicegroup.NewSharedLocation()

		entries, err := svc.List(ctx, *location, "get", "", "")
		if err != nil && !sdkerrors.IsObjectNotFound(err) {
			return fmt.Errorf("Failed to list service group entries: %w", err)
		}

		err = nil
		var dangling []string
		for _, elt := range entries {
			if strings.HasPrefix(elt.Name, prefix) {
				dangling = append(dangling, elt.Name)
				err = errors.Join(err, fmt.Errorf("service group entry not deleted properly: %s", elt.Name))
			}
		}

		if len(dangling) > 0 {
			deleteErr := svc.Delete(ctx, *location, dangling...)
			if deleteErr != nil && !sdkerrors.IsObjectNotFound(deleteErr) {
				err = errors.Join(err, fmt.Errorf("failed to delete service group entries: %w", deleteErr))
			}
		}

		return err
	}
}
