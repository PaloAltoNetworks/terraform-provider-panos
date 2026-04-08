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

func TestAccSchedule_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: schedule_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
				},
			},
		},
	})
}

const schedule_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_schedule" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  disable_override = "yes"
}
`

func TestAccSchedule_ScheduleType_NonRecurring(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: schedule_ScheduleType_NonRecurring_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("schedule_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"non_recurring": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("2025/01/01@00:00-2025/01/31@23:59"),
								knownvalue.StringExact("2025/12/24@00:00-2025/12/26@23:59"),
							}),
							"recurring": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const schedule_ScheduleType_NonRecurring_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_schedule" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  schedule_type = {
    non_recurring = ["2025/01/01@00:00-2025/01/31@23:59", "2025/12/24@00:00-2025/12/26@23:59"]
  }
}
`

func TestAccSchedule_ScheduleType_Recurring_Daily(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: schedule_ScheduleType_Recurring_Daily_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("schedule_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"non_recurring": knownvalue.Null(),
							"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"daily": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.StringExact("09:00-17:00"),
									knownvalue.StringExact("18:00-22:00"),
								}),
								"weekly": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const schedule_ScheduleType_Recurring_Daily_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_schedule" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  schedule_type = {
    recurring = {
      daily = ["09:00-17:00", "18:00-22:00"]
    }
  }
}
`

func TestAccSchedule_ScheduleType_Recurring_Weekly(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: schedule_ScheduleType_Recurring_Weekly_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("schedule_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"non_recurring": knownvalue.Null(),
							"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"daily": knownvalue.Null(),
								"weekly": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"monday":    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("08:00-12:00"), knownvalue.StringExact("13:00-17:00")}),
									"tuesday":   knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("08:00-12:00"), knownvalue.StringExact("13:00-17:00")}),
									"wednesday": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("08:00-12:00"), knownvalue.StringExact("13:00-17:00")}),
									"thursday":  knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("08:00-12:00"), knownvalue.StringExact("13:00-17:00")}),
									"friday":    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("08:00-12:00"), knownvalue.StringExact("13:00-17:00")}),
									"saturday":  knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10:00-14:00")}),
									"sunday":    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10:00-14:00")}),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const schedule_ScheduleType_Recurring_Weekly_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_schedule" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  schedule_type = {
    recurring = {
      weekly = {
        monday    = ["08:00-12:00", "13:00-17:00"]
        tuesday   = ["08:00-12:00", "13:00-17:00"]
        wednesday = ["08:00-12:00", "13:00-17:00"]
        thursday  = ["08:00-12:00", "13:00-17:00"]
        friday    = ["08:00-12:00", "13:00-17:00"]
        saturday  = ["10:00-14:00"]
        sunday    = ["10:00-14:00"]
      }
    }
  }
}
`

func TestAccSchedule_ScheduleType_Recurring_Weekly_SingleDay(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: schedule_ScheduleType_Recurring_Weekly_SingleDay_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_schedule.example",
						tfjsonpath.New("schedule_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"non_recurring": knownvalue.Null(),
							"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"daily": knownvalue.Null(),
								"weekly": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"monday":    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("09:00-17:00")}),
									"tuesday":   knownvalue.Null(),
									"wednesday": knownvalue.Null(),
									"thursday":  knownvalue.Null(),
									"friday":    knownvalue.Null(),
									"saturday":  knownvalue.Null(),
									"sunday":    knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const schedule_ScheduleType_Recurring_Weekly_SingleDay_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_schedule" "example" {
  depends_on = [panos_device_group.example]
  location = var.location

  name = var.prefix
  schedule_type = {
    recurring = {
      weekly = {
        monday = ["09:00-17:00"]
      }
    }
  }
}
`
