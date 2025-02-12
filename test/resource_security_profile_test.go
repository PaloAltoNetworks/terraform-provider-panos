package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/profiles/secgroup"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccSecurityProfileGroup(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccSecurityProfileGroupDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: securityProfileGroupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-sec-group", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("data_filtering"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("test-profile1"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("file_blocking"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("basic file blocking"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("spyware"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("url_filtering"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("virus"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("vulnerability"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_security_profile_group.group",
						tfjsonpath.New("wildfire_analysis"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("default"),
						}),
					),
				},
			},
			{
				Config: securityProfileGroupCleanupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					SecurityProfileGroupExpectNoEntriesInLocation(prefix),
				},
			},
		},
	})
}

const securityProfileGroupTmpl = `
variable prefix { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg1", var.prefix)
}

resource "panos_security_profile_group" "group" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-sec-group", var.prefix)

  disable_override = "yes"
  data_filtering = ["test-profile1"]
  file_blocking = ["basic file blocking"]
  #gtp = ["default"]
  #sctp = ["default"]
  spyware = ["default"]
  url_filtering = ["default"]
  virus = ["default"]
  vulnerability = ["default"]
  wildfire_analysis = ["default"]
}
`

const securityProfileGroupCleanupTmpl = `
variable prefix { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg1", var.prefix)
}
`

type securityProfileGroupExpectNoEntriesInLocation struct {
	prefix string
}

func SecurityProfileGroupExpectNoEntriesInLocation(prefix string) *securityProfileGroupExpectNoEntriesInLocation {
	return &securityProfileGroupExpectNoEntriesInLocation{
		prefix: prefix,
	}
}

func (o *securityProfileGroupExpectNoEntriesInLocation) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := secgroup.NewService(sdkClient)
	location := secgroup.NewDeviceGroupLocation()
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

func testAccSecurityProfileGroupDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		service := secgroup.NewService(sdkClient)

		location := secgroup.NewDeviceGroupLocation()
		location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg1", prefix)

		ctx := context.TODO()
		entries, err := service.List(ctx, *location, "get", "", "")
		if err != nil && !sdkerrors.IsObjectNotFound(err) {
			return fmt.Errorf("failed to list existing entries via sdk: %w", err)
		}

		var leftEntries []string
		for _, elt := range entries {
			if strings.HasPrefix(elt.Name, prefix) {
				leftEntries = append(leftEntries, elt.Name)
			}
		}

		if len(leftEntries) > 0 {
			err := fmt.Errorf("terraform failed to remove entries from the server")
			delErr := service.Delete(ctx, *location, leftEntries...)
			if delErr != nil {
				return errors.Join(err, delErr)
			}
		}

		return nil
	}
}
