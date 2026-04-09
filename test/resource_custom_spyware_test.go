package provider_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCustomSpyware_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(strconv.Itoa(spywareName)),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("bugtraq"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("bugtraq1"),
							knownvalue.StringExact("bugtraq2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("comment"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("cve"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("cve1"),
							knownvalue.StringExact("cve2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("direction"),
						knownvalue.StringExact("both"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("reference"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ref1"),
							knownvalue.StringExact("ref2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("severity"),
						knownvalue.StringExact("critical"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("threatname"),
						knownvalue.StringExact("threatname"),
					),
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("vendor"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("vendor1"),
							knownvalue.StringExact("vendor2"),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	bugtraq = ["bugtraq1", "bugtraq2"]
	comment = "comment"
	cve = ["cve1", "cve2"]
	direction = "both"
	disable_override = "yes"
	reference = ["ref1", "ref2"]
	severity = "critical"
	threatname = "threatname"
	vendor = ["vendor1", "vendor2"]
}
`

func TestAccCustomSpyware_DefaultActionAlert(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionAlert_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert":        knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"allow":        knownvalue.Null(),
							"block_ip":     knownvalue.Null(),
							"drop":         knownvalue.Null(),
							"reset_both":   knownvalue.Null(),
							"reset_client": knownvalue.Null(),
							"reset_server": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionAlert_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		alert = {}
	}
}
`

func TestAccCustomSpyware_DefaultActionAllow(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionAllow_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert":        knownvalue.Null(),
							"allow":        knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"block_ip":     knownvalue.Null(),
							"drop":         knownvalue.Null(),
							"reset_both":   knownvalue.Null(),
							"reset_client": knownvalue.Null(),
							"reset_server": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionAllow_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		allow = {}
	}
}
`

func TestAccCustomSpyware_DefaultActionBlockIp(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionBlockIp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert": knownvalue.Null(),
							"allow": knownvalue.Null(),
							"block_ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"duration": knownvalue.Int64Exact(120),
								"track_by": knownvalue.StringExact("source-and-destination"),
							}),
							"drop":         knownvalue.Null(),
							"reset_both":   knownvalue.Null(),
							"reset_client": knownvalue.Null(),
							"reset_server": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionBlockIp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		block_ip = {
			duration = 120
			track_by = "source-and-destination"
		}
	}
}
`

func TestAccCustomSpyware_DefaultActionDrop(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionDrop_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert":        knownvalue.Null(),
							"allow":        knownvalue.Null(),
							"block_ip":     knownvalue.Null(),
							"drop":         knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"reset_both":   knownvalue.Null(),
							"reset_client": knownvalue.Null(),
							"reset_server": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionDrop_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		drop = {}
	}
}
`

func TestAccCustomSpyware_DefaultActionResetBoth(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionResetBoth_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert":        knownvalue.Null(),
							"allow":        knownvalue.Null(),
							"block_ip":     knownvalue.Null(),
							"drop":         knownvalue.Null(),
							"reset_both":   knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"reset_client": knownvalue.Null(),
							"reset_server": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionResetBoth_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		reset_both = {}
	}
}
`

func TestAccCustomSpyware_DefaultActionResetClient(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionResetClient_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert":        knownvalue.Null(),
							"allow":        knownvalue.Null(),
							"block_ip":     knownvalue.Null(),
							"drop":         knownvalue.Null(),
							"reset_both":   knownvalue.Null(),
							"reset_client": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"reset_server": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionResetClient_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		reset_client = {}
	}
}
`

func TestAccCustomSpyware_DefaultActionResetServer(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_DefaultActionResetServer_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("default_action"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"alert":        knownvalue.Null(),
							"allow":        knownvalue.Null(),
							"block_ip":     knownvalue.Null(),
							"drop":         knownvalue.Null(),
							"reset_both":   knownvalue.Null(),
							"reset_client": knownvalue.Null(),
							"reset_server": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_DefaultActionResetServer_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	default_action = {
		reset_server = {}
	}
}
`

func TestAccCustomSpyware_SignatureCombination(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_SignatureCombination_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("signature"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"combination": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"order_free": knownvalue.Bool(true),
								"time_attribute": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"interval":  knownvalue.Int64Exact(120),
									"threshold": knownvalue.Int64Exact(10),
									"track_by":  knownvalue.StringExact("source"),
								}),
								"and_condition": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name": knownvalue.StringExact("and1"),
										"or_condition": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":      knownvalue.StringExact("or1"),
												"threat_id": knownvalue.StringExact("10004"),
											}),
										}),
									}),
								}),
							}),
							"standard": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const customSpyware_SignatureCombination_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	signature = {
		combination = {
			order_free = true
			time_attribute = {
				interval = 120
				threshold = 10
				track_by = "source"
			}
			and_condition = [
				{
					name = "and1"
					or_condition = [
						{
							name = "or1"
							threat_id = "10004"
						}
					]
				}
			]
		}
	}
}
`

func TestAccCustomSpyware_SignatureStandard(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	spywareName := acctest.RandIntRange(6900001, 7000000)

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
				Config: customSpyware_SignatureStandard_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"spyware_name": config.IntegerVariable(spywareName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_custom_spyware.example",
						tfjsonpath.New("signature"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"combination": knownvalue.Null(),
							"standard": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":       knownvalue.StringExact("standard1"),
									"comment":    knownvalue.StringExact("comment"),
									"scope":      knownvalue.StringExact("session"),
									"order_free": knownvalue.Bool(true),
									"and_condition": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name": knownvalue.StringExact("and1"),
											"or_condition": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name": knownvalue.StringExact("or1"),
													"operator": knownvalue.ObjectExact(map[string]knownvalue.Check{
														"less_than": knownvalue.ObjectExact(map[string]knownvalue.Check{
															"value":     knownvalue.Int64Exact(10),
															"context":   knownvalue.StringExact("diameter-req-avp-code"),
															"qualifier": knownvalue.Null(),
														}),
														"equal_to":      knownvalue.Null(),
														"greater_than":  knownvalue.Null(),
														"pattern_match": knownvalue.Null(),
													}),
												}),
											}),
										}),
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

const customSpyware_SignatureStandard_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "spyware_name" { type = number }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_custom_spyware" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.spyware_name
	threatname = "threatname"
	severity = "critical"
	signature = {
		standard = [
			{
				name = "standard1"
				comment = "comment"
				scope = "session"
				order_free = true
				and_condition = [
					{
						name = "and1"
						or_condition = [
							{
								name = "or1"
								operator = {
									less_than = {
										value = 10
										context = "diameter-req-avp-code"
									}
								}
							}
						]
					}
				]
			}
		]
	}
}`
