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

func TestAccDynamicUpdates_1(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUpdatesConfig1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// antivirus
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("threshold"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("action"),
						knownvalue.StringExact("download-and-install"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					// app_profile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("threshold"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("action"),
						knownvalue.StringExact("download-and-install"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					// global_protect_clientless_vpn
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					// global_protect_datafile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					// statistics_service
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("url_reports"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("application_reports"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("file_identification_reports"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("health_performance_reports"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("passive_dns_monitoring"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("threat_prevention_information"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("threat_prevention_pcap"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("statistics_service").
							AtMapKey("threat_prevention_reports"),
						knownvalue.Bool(true),
					),
					// threats
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("threshold"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("new_app_threshold"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("disable_new_content"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("action"),
						knownvalue.StringExact("download-and-install"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("daily").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					// wf_private
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_15_mins").
							AtMapKey("action"),
						knownvalue.StringExact("download-and-install"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_15_mins").
							AtMapKey("at"),
						knownvalue.Int64Exact(10),
					),
					// wildfire
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_15_mins").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_15_mins").
							AtMapKey("action"),
						knownvalue.StringExact("download-and-install"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_15_mins").
							AtMapKey("at"),
						knownvalue.Int64Exact(10),
					),
				},
			},
		},
	})
}

func TestAccDynamicUpdates_2(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUpdatesConfig2,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// antivirus
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("at"),
						knownvalue.Int64Exact(20),
					),
					// app_profile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					// global_protect_clientless_vpn
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("at"),
						knownvalue.Int64Exact(20),
					),
					// global_protect_datafile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("at"),
						knownvalue.Int64Exact(20),
					),
					// threats
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("threshold"),
						knownvalue.Int64Exact(15),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("new_app_threshold"),
						knownvalue.Int64Exact(15),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("disable_new_content"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("hourly").
							AtMapKey("at"),
						knownvalue.Int64Exact(20),
					),
					// wf_private
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_30_mins").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_30_mins").
							AtMapKey("at"),
						knownvalue.Int64Exact(10),
					),
					// wildfire
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_30_mins").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_30_mins").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_30_mins").
							AtMapKey("at"),
						knownvalue.Int64Exact(10),
					),
				},
			},
		},
	})

}

func TestAccDynamicUpdates_3(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUpdatesConfig3,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// antivirus
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					// app_profile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("app_profile").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("day_of_week"),
						knownvalue.StringExact("monday"),
					),
					// global_protect_clientless_vpn
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					// global_protect_datafile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					// threats
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					// wf_private
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_hour").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_hour").
							AtMapKey("at"),
						knownvalue.Int64Exact(10),
					),
					// wildfire
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_hour").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_hour").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_hour").
							AtMapKey("at"),
						knownvalue.Int64Exact(10),
					),
				},
			},
		},
	})
}

func TestAccDynamicUpdates_4(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUpdatesConfig4,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// antivirus
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("anti_virus").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("day_of_week"),
						knownvalue.StringExact("monday"),
					),
					// global_protect_clientless_vpn
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_clientless_vpn").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("day_of_week"),
						knownvalue.StringExact("monday"),
					),
					// global_protect_datafile
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("global_protect_datafile").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("day_of_week"),
						knownvalue.StringExact("monday"),
					),
					// threats
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("at"),
						knownvalue.StringExact("20:10"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("threats").
							AtMapKey("recurring").
							AtMapKey("weekly").
							AtMapKey("day_of_week"),
						knownvalue.StringExact("monday"),
					),
					// wf_private
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_5_mins").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("every_5_mins").
							AtMapKey("at"),
						knownvalue.Int64Exact(4),
					),
					// wildfire
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_min").
							AtMapKey("sync_to_peer"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("every_min").
							AtMapKey("action"),
						knownvalue.StringExact("download-only"),
					),
				},
			},
		},
	})
}

func TestAccDynamicUpdates_5(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUpdatesConfig5,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// antivirus
					// global_protect_clientless_vpn
					// global_protect_datafile
					// threats
					// wf_private
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wf_private").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
					// wildfire
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("none"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
				},
			},
		},
	})
}

func TestAccDynamicUpdates_6(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: dynamicUpdatesConfig6,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// antivirus
					// global_protect_clientless_vpn
					// global_protect_datafile
					// threats
					// wf_private
					// wildfire
					statecheck.ExpectKnownValue(
						"panos_dynamic_updates.updates",
						tfjsonpath.New("update_schedule").
							AtMapKey("wildfire").
							AtMapKey("recurring").
							AtMapKey("real_time"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
				},
			},
		},
	})
}

const dynamicUpdatesConfig1 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_dynamic_updates" "updates" {
  location = { template = { name = panos_template.template.name } }

  update_schedule = {
    anti_virus = {
      recurring = {
        sync_to_peer = true
        threshold = 10
        daily = { action = "download-and-install", at = "20:10" }
      }
    }

    app_profile = {
      recurring = {
        sync_to_peer = true
        threshold = 10
        daily = { action = "download-and-install", at = "20:10" }
      }
    }

    global_protect_clientless_vpn = {
      recurring = {
        daily = { action = "download-and-install", at = "20:10" }
      }
    }

    global_protect_datafile = {
      recurring = {
        daily = { action = "download-and-install", at = "20:10" }
      }
    }

    statistics_service = {
      url_reports                 = true
      application_reports         = true
      file_identification_reports = true
      health_performance_reports  = true
      passive_dns_monitoring      = true

      threat_prevention_information = true
      threat_prevention_pcap        = true
      threat_prevention_reports     = true
    }

    threats = {
      recurring = {
        sync_to_peer      = true
        threshold         = 10
        new_app_threshold = 10

        daily = {
          disable_new_content = true
          action = "download-and-install"
          at = "20:10"
        }
      }
    }

    wf_private = {
      recurring = {
        sync_to_peer     = true
        every_15_mins = { action = "download-and-install", at = 10 }
      }
    }

    wildfire = {
      recurring = {
        every_15_mins = {
          sync_to_peer = true
          action       = "download-and-install"
          at = 10
        }
      }
    }
  }
}
`

const dynamicUpdatesConfig2 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_dynamic_updates" "updates" {
  location = { template = { name = panos_template.template.name } }

  update_schedule = {
    anti_virus = {
      recurring = {
        sync_to_peer = false
        hourly = { action = "download-only", at = 20 }
      }
    }

    app_profile = {
      recurring = {
        sync_to_peer = false
        none = {}
      }
    }

    global_protect_clientless_vpn = {
      recurring = {
        hourly = { action = "download-only", at = 20 }
      }
    }

    global_protect_datafile = {
      recurring = {
        hourly = { action = "download-only", at = 20 }
      }
    }

    threats = {
      recurring = {
        sync_to_peer      = false
        threshold         = 15
        new_app_threshold = 15

        hourly = {
          disable_new_content = true
          action = "download-only"
          at = 20
        }
      }
    }

    wf_private = {
      recurring = {
        sync_to_peer     = true
        every_30_mins = { action = "download-only", at = 10 }
      }
    }

    wildfire = {
      recurring = {
        every_30_mins = {
          sync_to_peer = true
          action       = "download-only"
          at = 10
        }
      }
    }
  }
}
`

const dynamicUpdatesConfig3 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_dynamic_updates" "updates" {
  location = { template = { name = panos_template.template.name } }

  update_schedule = {
    anti_virus = {
      recurring = {
        none = {}
      }
    }

    app_profile = {
      recurring = {
        weekly = { action = "download-only", at = "20:10", day_of_week = "monday" }
      }
    }

    global_protect_clientless_vpn = {
      recurring = {
        none = {}
      }
    }

    global_protect_datafile = {
      recurring = {
        none = {}
      }
    }

    threats = {
      recurring = {
        none = {}
      }
    }

    wf_private = {
      recurring = {
        every_hour = { action = "download-only", at = 10 }
      }
    }

    wildfire = {
      recurring = {
        every_hour = {
          sync_to_peer = true
          action       = "download-only"
          at = 10
        }
      }
    }
  }
}
`

const dynamicUpdatesConfig4 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_dynamic_updates" "updates" {
  location = { template = { name = panos_template.template.name } }

  update_schedule = {
    anti_virus = {
      recurring = {
        weekly = { action = "download-only", at = "20:10", day_of_week = "monday" }
      }
    }

    global_protect_clientless_vpn = {
      recurring = {
        weekly = { action = "download-only", at = "20:10", day_of_week = "monday" }
      }
    }

    global_protect_datafile = {
      recurring = {
        weekly = { action = "download-only", at = "20:10", day_of_week = "monday" }
      }
    }

    threats = {
      recurring = {
        weekly = {
          disable_new_content = true
          action = "download-only"
          at = "20:10", day_of_week = "monday"
        }
      }
    }

    wf_private = {
      recurring = {
        every_5_mins = { action = "download-only", at = 4 }
      }
    }

    wildfire = {
      recurring = {
        every_min = {
          sync_to_peer = true
          action       = "download-only"
        }
      }
    }
  }
}
`

const dynamicUpdatesConfig5 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_dynamic_updates" "updates" {
  location = { template = { name = panos_template.template.name } }

  update_schedule = {
    wf_private = {
      recurring = {
        none = {}
      }
    }

    wildfire = {
      recurring = {
        none = {}
      }
    }
  }
}
`

const dynamicUpdatesConfig6 = `
variable "prefix" { type = string }

locals {
  template_name = format("%s-tmpl", var.prefix)
}

resource "panos_template" "template" {
  location = { panorama = {} }

  name = local.template_name
}

resource "panos_dynamic_updates" "updates" {
  location = { template = { name = panos_template.template.name } }

  update_schedule = {
    wildfire = {
      recurring = {
        real_time = {}
      }
    }
  }
}
`
