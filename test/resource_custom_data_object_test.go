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

func TestAccCustomDataObject_Basic(t *testing.T) {
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
				Config: customDataObject_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_data_object.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_data_object.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_data_object.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
				},
			},
		},
	})
}

const customDataObject_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_data_object" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	description = "test description"
	disable_override = "yes"
}
`

func TestAccCustomDataObject_PatternType_FileProperties(t *testing.T) {
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
				Config: customDataObject_PatternType_FileProperties_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_data_object.example",
						tfjsonpath.New("pattern_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"file_properties": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"pattern": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":           knownvalue.StringExact("test-pattern"),
										"file_type":      knownvalue.StringExact("pdf"),
										"file_property":  knownvalue.StringExact("panav-rsp-pdf-dlp-author"),
										"property_value": knownvalue.StringExact("author"),
									}),
								}),
							}),
							"predefined": knownvalue.Null(),
							"regex":      knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customDataObject_PatternType_FileProperties_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_data_object" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	pattern_type = {
		file_properties = {
			pattern = [
				{
					name = "test-pattern"
					file_type = "pdf"
					file_property = "panav-rsp-pdf-dlp-author"
					property_value = "author"
				}
			]
		}
	}
}
`

func TestAccCustomDataObject_PatternType_Predefined(t *testing.T) {
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
				Config: customDataObject_PatternType_Predefined_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_data_object.example",
						tfjsonpath.New("pattern_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"file_properties": knownvalue.Null(),
							"predefined": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"pattern": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":      knownvalue.StringExact("ABA-Routing-Number"),
										"file_type": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("xlsx")}),
									}),
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":      knownvalue.StringExact("credit-card-numbers"),
										"file_type": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("text/html")}),
									}),
								}),
							}),
							"regex": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customDataObject_PatternType_Predefined_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_data_object" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	pattern_type = {
		predefined = {
			pattern = [
				{
					name = "ABA-Routing-Number"
					file_type = ["xlsx"]
				},
				{
					name = "credit-card-numbers",
					file_type = ["text/html"]
				},
			]
		}
	}
}
`

func TestAccCustomDataObject_PatternType_Regex(t *testing.T) {
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
				Config: customDataObject_PatternType_Regex_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_data_object.example",
						tfjsonpath.New("pattern_type"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"file_properties": knownvalue.Null(),
							"predefined":      knownvalue.Null(),
							"regex": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"pattern": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":      knownvalue.StringExact("test-pattern"),
										"file_type": knownvalue.Null(),
										"regex":     knownvalue.StringExact("test-regex"),
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

const customDataObject_PatternType_Regex_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_data_object" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	pattern_type = {
		regex = {
			pattern = [
				{
					name = "test-pattern"
					regex = "test-regex"
				}
			]
		}
	}
}
`
