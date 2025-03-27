package provider_test

import (
	"context"
	"fmt"
	"testing"

	sdkErrors "github.com/PaloAltoNetworks/pango/errors"
	"github.com/PaloAltoNetworks/pango/objects/extdynlist"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccPanosExternalDynamicList_1(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList1Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.
							New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"imei":           knownvalue.Null(),
							"imsi":           knownvalue.Null(),
							"ip":             knownvalue.Null(),
							"predefined_ip":  knownvalue.Null(),
							"predefined_url": knownvalue.Null(),
							"url":            knownvalue.Null(),
							"domain": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile": knownvalue.StringExact("cert-profile"),
								"description":         knownvalue.StringExact("description"),
								"exception_list":      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("exception1"), knownvalue.StringExact("exception2")}),
								"expand_domain":       knownvalue.Bool(true),
								"url":                 knownvalue.StringExact("https://example.com/list.txt"),
								"auth": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"username": knownvalue.StringExact("user"),
									"password": knownvalue.StringExact("password"),
								}),
								"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"daily": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"at": knownvalue.StringExact("20"),
									}),
									"five_minute": knownvalue.Null(),
									"hourly":      knownvalue.Null(),
									"monthly":     knownvalue.Null(),
									"weekly":      knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosExternalDynamicList_2(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList2Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":         knownvalue.Null(),
							"imsi":           knownvalue.Null(),
							"ip":             knownvalue.Null(),
							"predefined_ip":  knownvalue.Null(),
							"predefined_url": knownvalue.Null(),
							"url":            knownvalue.Null(),
							"imei": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile": knownvalue.StringExact("cert-profile"),
								"description":         knownvalue.StringExact("description"),
								"exception_list":      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("123456789012345"), knownvalue.StringExact("012345678912345")}),
								"url":                 knownvalue.StringExact("https://example.com/list.txt"),
								"auth": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"username": knownvalue.StringExact("user"),
									"password": knownvalue.StringExact("password"),
								}),
								"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"daily":       knownvalue.Null(),
									"five_minute": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
									"hourly":      knownvalue.Null(),
									"monthly":     knownvalue.Null(),
									"weekly":      knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosExternalDynamicList_3(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList3Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":         knownvalue.Null(),
							"imei":           knownvalue.Null(),
							"ip":             knownvalue.Null(),
							"predefined_ip":  knownvalue.Null(),
							"predefined_url": knownvalue.Null(),
							"url":            knownvalue.Null(),
							"imsi": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile": knownvalue.StringExact("cert-profile"),
								"description":         knownvalue.StringExact("description"),
								"exception_list":      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("123456789012345"), knownvalue.StringExact("012345678912345")}),
								"url":                 knownvalue.StringExact("https://example.com/list.txt"),
								"auth": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"username": knownvalue.StringExact("user"),
									"password": knownvalue.StringExact("password"),
								}),
								"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"daily":       knownvalue.Null(),
									"five_minute": knownvalue.Null(),
									"hourly":      knownvalue.ObjectExact(map[string]knownvalue.Check{}),
									"monthly":     knownvalue.Null(),
									"weekly":      knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosExternalDynamicList_4(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList4Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":         knownvalue.Null(),
							"imei":           knownvalue.Null(),
							"imsi":           knownvalue.Null(),
							"predefined_ip":  knownvalue.Null(),
							"predefined_url": knownvalue.Null(),
							"url":            knownvalue.Null(),
							"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile": knownvalue.StringExact("cert-profile"),
								"description":         knownvalue.StringExact("description"),
								"exception_list":      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("1.1.1.1/32"), knownvalue.StringExact("1.1.1.2/32")}),
								"url":                 knownvalue.StringExact("https://example.com/list.txt"),
								"auth": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"username": knownvalue.StringExact("user"),
									"password": knownvalue.StringExact("password"),
								}),
								"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"daily":       knownvalue.Null(),
									"five_minute": knownvalue.Null(),
									"hourly":      knownvalue.Null(),
									"monthly": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"at":           knownvalue.StringExact("20"),
										"day_of_month": knownvalue.Int64Exact(6),
									}),
									"weekly": knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosExternalDynamicList_5(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList5Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":         knownvalue.Null(),
							"imei":           knownvalue.Null(),
							"imsi":           knownvalue.Null(),
							"ip":             knownvalue.Null(),
							"predefined_url": knownvalue.Null(),
							"url":            knownvalue.Null(),
							"predefined_ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"description":    knownvalue.StringExact("description"),
								"exception_list": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("1.1.1.1/32"), knownvalue.StringExact("1.1.1.2/32")}),
								"url":            knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosExternalDynamicList_6(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList6Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":        knownvalue.Null(),
							"imei":          knownvalue.Null(),
							"imsi":          knownvalue.Null(),
							"ip":            knownvalue.Null(),
							"predefined_ip": knownvalue.Null(),
							"url":           knownvalue.Null(),
							"predefined_url": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"description":    knownvalue.StringExact("description"),
								"exception_list": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("exception1"), knownvalue.StringExact("exception2")}),
								"url":            knownvalue.StringExact("panw-auth-portal-exclude-list"),
							}),
						}),
					),
				},
			},
		},
	})
}

func TestAccPanosExternalDynamicList_7(t *testing.T) {
	t.Parallel()
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		CheckDestroy:             testAccExternalDynamicListCheckDestroy(prefix),
		Steps: []resource.TestStep{
			{
				Config: externalDynamicList7Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-list", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("no"),
					),
					statecheck.ExpectKnownValue(
						"panos_external_dynamic_list.list",
						tfjsonpath.New("type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"domain":         knownvalue.Null(),
							"imei":           knownvalue.Null(),
							"imsi":           knownvalue.Null(),
							"ip":             knownvalue.Null(),
							"predefined_ip":  knownvalue.Null(),
							"predefined_url": knownvalue.Null(),
							"url": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile": knownvalue.StringExact("cert-profile"),
								"description":         knownvalue.StringExact("description"),
								"exception_list":      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("exception1"), knownvalue.StringExact("exception2")}),
								"url":                 knownvalue.StringExact("https://example.com/list.txt"),
								"auth": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"username": knownvalue.StringExact("user"),
									"password": knownvalue.StringExact("password"),
								}),
								"recurring": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"daily":       knownvalue.Null(),
									"five_minute": knownvalue.Null(),
									"hourly":      knownvalue.Null(),
									"monthly":     knownvalue.Null(),
									"weekly": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"at":          knownvalue.StringExact("20"),
										"day_of_week": knownvalue.StringExact("friday"),
									}),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const externalDynamicList1Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    domain = {
      certificate_profile = "cert-profile"
      description = "description"
      exception_list = ["exception1", "exception2"]
      expand_domain = true
      url = "https://example.com/list.txt"
      auth = {
        username = "user"
        password = "password"
      }
      recurring = {
        daily = {
          at = "20"
        }
      }
    }
  }
}
`

const externalDynamicList2Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    imei = {
      certificate_profile = "cert-profile"
      description = "description"
      exception_list = ["123456789012345", "012345678912345"]
      url = "https://example.com/list.txt"
      auth = {
        username = "user"
        password = "password"
      }
      recurring = {
        five_minute = {}
      }
    }
  }
}
`

const externalDynamicList3Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    imsi = {
      certificate_profile = "cert-profile"
      description = "description"
      exception_list = ["123456789012345", "012345678912345"]
      url = "https://example.com/list.txt"
      auth = {
        username = "user"
        password = "password"
      }
      recurring = {
        hourly = {}
      }
    }
  }
}
`

const externalDynamicList4Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    ip = {
      certificate_profile = "cert-profile"
      description = "description"
      exception_list = ["1.1.1.1/32", "1.1.1.2/32"]
      url = "https://example.com/list.txt"
      auth = {
        username = "user"
        password = "password"
      }
      recurring = {
        monthly = {
          at = "20"
          day_of_month = 6
        }
      }
    }
  }
}
`

const externalDynamicList5Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    predefined_ip = {
      description = "description"

      exception_list = ["1.1.1.1/32", "1.1.1.2/32"]
      #url = "https://example.com/list2.txt"
    }
  }
}
`

const externalDynamicList6Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    predefined_url = {
      description = "description"
      exception_list = ["exception1", "exception2"]
      url = "panw-auth-portal-exclude-list"
    }
  }
}
`

const externalDynamicList7Tmpl = `
variable "prefix" { type = string }

resource "panos_device_group" "dg" {
  location = { panorama = {} }

  name = format("%s-dg", var.prefix)
}

resource "panos_external_dynamic_list" "list" {
  location = { device_group = { name = panos_device_group.dg.name } }

  name = format("%s-list", var.prefix)

  disable_override = "no"
  type = {
    url = {
      certificate_profile = "cert-profile"
      description = "description"
      exception_list = ["exception1", "exception2"]
      url = "https://example.com/list.txt"
      auth = {
        username = "user"
        password = "password"
      }
      recurring = {
        weekly = {
          at = "20"
          day_of_week = "friday"
        }
      }
    }
  }
}
`

func testAccExternalDynamicListCheckDestroy(prefix string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		api := extdynlist.NewService(sdkClient)
		location := extdynlist.NewDeviceGroupLocation()

		location.DeviceGroup.DeviceGroup = fmt.Sprintf("%s-dg", prefix)

		ctx := context.TODO()

		entry := fmt.Sprintf("%s-list", prefix)

		reply, err := api.Read(ctx, *location, entry, "show")
		if err != nil && !sdkErrors.IsObjectNotFound(err) {
			return fmt.Errorf("reading external dynamic list entry %s via sdk: %w", entry, err)
		}

		if reply != nil {
			if reply.EntryName() == entry {
				return fmt.Errorf("external dynamic list object still exists: %s", entry)
			}
		}

		return nil
	}
}
