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

func TestAccLogExportSchedule_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: logExportSchedule_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Basic log export schedule"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("log_type"),
						knownvalue.StringExact("traffic"),
					),
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("start_time"),
						knownvalue.StringExact("03:30"),
					),
				},
			},
		},
	})
}

const logExportSchedule_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_log_export_schedule" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  description = "Basic log export schedule"
  enable = true
  log_type = "traffic"
  start_time = "03:30"
}
`

func TestAccLogExportSchedule_Protocol_Ftp(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: logExportSchedule_Protocol_Ftp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ftp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"hostname":     knownvalue.StringExact("ftp.example.com"),
								"passive_mode": knownvalue.Bool(true),
								"password":     knownvalue.StringExact("ftppass123"),
								"path":         knownvalue.StringExact("/logs/export"),
								"port":         knownvalue.Int64Exact(21),
								"username":     knownvalue.StringExact("ftpuser"),
							}),
							"scp": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const logExportSchedule_Protocol_Ftp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_log_export_schedule" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  log_type = "threat"
  enable = true
  start_time = "02:00"

  protocol = {
    ftp = {
      hostname = "ftp.example.com"
      passive_mode = true
      password = "ftppass123"
      path = "/logs/export"
      port = 21
      username = "ftpuser"
    }
  }
}
`

func TestAccLogExportSchedule_Protocol_Scp(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: logExportSchedule_Protocol_Scp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_log_export_schedule.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"scp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"hostname": knownvalue.StringExact("scp.example.com"),
								"password": knownvalue.StringExact("scppass123"),
								"path":     knownvalue.StringExact("/var/logs/export"),
								"port":     knownvalue.Int64Exact(22),
								"username": knownvalue.StringExact("scpuser"),
							}),
							"ftp": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const logExportSchedule_Protocol_Scp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_log_export_schedule" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  log_type = "wildfire"
  enable = true
  start_time = "04:15"

  protocol = {
    scp = {
      hostname = "scp.example.com"
      password = "scppass123"
      path = "/var/logs/export"
      port = 22
      username = "scpuser"
    }
  }
}
`
