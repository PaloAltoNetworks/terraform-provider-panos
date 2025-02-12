package provider_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	sdkerrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/profiles/urlfiltering"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccUrlFilteringSecurityProfile(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccUrlFilteringSecurityProfileDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: urlFilteringSecurityProfile1Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-profile", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("description"),
						knownvalue.StringExact("description"),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("alert"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("games"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("allow"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("music"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("block"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("news"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("cloud_inline_cat"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("continue"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("travel"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("enable_container_page"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("local_inline_cat"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("log_container_page_only"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("log_http_hdr_referer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("log_http_hdr_user_agent"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("log_http_hdr_xff"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("safe_search_enforcement"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("credential_enforcement"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"log_severity": knownvalue.StringExact("medium"),
							"alert": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("games"),
							}),
							"allow": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("music"),
							}),
							"block": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("news"),
							}),
							"continue": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("travel"),
							}),
							"mode": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"disabled":           knownvalue.Null(),
								"domain_credentials": knownvalue.ObjectExact(nil),
								"group_mapping":      knownvalue.Null(),
								"ip_user":            knownvalue.Null(),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("http_header_insertion"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":             knownvalue.StringExact("header1"),
								"disable_override": knownvalue.StringExact("yes"),
								"type": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name": knownvalue.StringExact("Google Apps Access Control"),
										"domains": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.StringExact("example.com"),
										}),
										"headers": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":   knownvalue.StringExact("header1"),
												"header": knownvalue.StringExact("header1"),
												"log":    knownvalue.Bool(true),
												"value":  knownvalue.StringExact("value1"),
											}),
										}),
									}),
								}),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_url_filtering_security_profile.rules",
						tfjsonpath.New("mlav_category_exception"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-category", prefix)),
						}),
					),
				},
			},
			{
				Config: urlFilteringSecurityProfileCleanupTmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					UrlFilteringSecurityProfileExpectNoEntriesInLocation(prefix),
				},
			},
		},
	})
}

const urlFilteringSecurityProfile1Tmpl = `
variable prefix { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg1", var.prefix)
}

resource "panos_custom_url_category" "category" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-category", var.prefix)
}

resource "panos_url_filtering_security_profile" "rules" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-profile", var.prefix)
  description = "description"
  disable_override = "yes"

  alert = ["games"]
  allow = ["music"]
  block = ["news"]
  cloud_inline_cat = true
  continue = ["travel"]
  enable_container_page = true
  local_inline_cat = true
  log_container_page_only = true
  log_http_hdr_referer = true
  log_http_hdr_user_agent = true
  log_http_hdr_xff = true
  safe_search_enforcement = true
  credential_enforcement = {
    log_severity = "medium"
    alert = ["games"]
    allow = ["music"]
    block = ["news"]
    continue = ["travel"]
    mode = {
      domain_credentials = {}
    }
  }
  http_header_insertion = [{
    name = "header1"
    disable_override = "yes"
    type = [{
      name = "Google Apps Access Control"
      domains = ["example.com"]
      headers = [{
        name = "header1"
        header = "header1"
        log = true
        value = "value1"
      }]
    }]
  }]
  mlav_category_exception = [resource.panos_custom_url_category.category.name]
}
`

const urlFilteringSecurityProfileCleanupTmpl = `
variable prefix { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg1", var.prefix)
}
`

type urlFilteringSecurityProfileExpectNoEntriesInLocation struct {
	prefix string
}

func UrlFilteringSecurityProfileExpectNoEntriesInLocation(prefix string) *urlFilteringSecurityProfileExpectNoEntriesInLocation {
	return &urlFilteringSecurityProfileExpectNoEntriesInLocation{
		prefix: prefix,
	}
}

func (o *urlFilteringSecurityProfileExpectNoEntriesInLocation) CheckState(ctx context.Context, req statecheck.CheckStateRequest, resp *statecheck.CheckStateResponse) {
	service := urlfiltering.NewService(sdkClient)
	location := urlfiltering.NewDeviceGroupLocation()
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

func testAccUrlFilteringSecurityProfileDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		service := urlfiltering.NewService(sdkClient)

		location := urlfiltering.NewDeviceGroupLocation()
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
