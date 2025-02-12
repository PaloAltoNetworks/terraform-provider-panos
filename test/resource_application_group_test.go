package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/application/group"

	//	"github.com/hashicorp/terraform-plugin-testing/compare"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type applicationGroupExpectNoEntriesInLocation struct {
	prefix string
}

func ApplicationGroupExpectNoEntriesInLocation(prefix string) *applicationGroupExpectNoEntriesInLocation {
	return &applicationGroupExpectNoEntriesInLocation{
		prefix: prefix,
	}
}

func (o *applicationGroupExpectNoEntriesInLocation) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := group.NewService(sdkClient)
	location := group.NewDeviceGroupLocation()
	location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg1", o.prefix)
	objects, err := service.List(ctx, *location, "get", "", "")
	if err != nil && !sdkerrors.IsObjectNotFound(err) {
		resp.Error = fmt.Errorf("failed to query server for entries: %w", err)
		return
	}

	var dangling []string
	for _, elt := range objects {
		if strings.HasPrefix(elt.Name, o.prefix) {
			dangling = append(dangling, elt.Name)
		}
	}

	if len(dangling) > 0 {
		resp.Error = fmt.Errorf("delete of the resource didn't remove it from the server")
	}
}

func TestAccPanosApplicationGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccCheckPanosApplicationGroupDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: applicationGroupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_application_group.group",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-appgroup", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_application_group.group",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_application_group.group",
						tfjsonpath.New("members"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("amazon-echo"),
							knownvalue.StringExact("amazon-alexa"),
						}),
					),
				},
			},
			{
				Config: applicationGroupCleanupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					ApplicationGroupExpectNoEntriesInLocation(prefix),
				},
			},
		},
	})
}

const applicationGroupTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg1", var.prefix)
}

resource "panos_application_group" "group" {
  location = { device_group = { name = panos_device_group.dg.name }}

  name = format("%s-appgroup", var.prefix)

  disable_override = "yes"
  members = ["amazon-echo", "amazon-alexa"]
}
`

const applicationGroupCleanupTmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg1", var.prefix)
}
`

func testAccCheckPanosApplicationGroupDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := group.NewService(sdkClient)
		location := group.NewDeviceGroupLocation()
		location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg1", prefix)

		ctx := context.TODO()
		existing, err := api.List(ctx, *location, "get", "", "")
		if err != nil && !sdkerrors.IsObjectNotFound(err) {
			return err
		}

		var dangling []string
		for _, elt := range existing {
			if strings.HasPrefix(elt.Name, prefix) {
				dangling = append(dangling, elt.Name)
			}
		}

		if len(dangling) > 0 {
			err = fmt.Errorf("Some entries were left after terraform teardown")
			deleteErr := api.Delete(ctx, *location, dangling...)
			if deleteErr != nil {
				err = errors.Join(err, deleteErr)
			}
			return err
		}

		return nil
	}
}
