
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

func TestAccDosProtectionProfile_Basic(t *testing.T) {
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
				Config: dosProtectionProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("disable_override"),
						knownvalue.StringExact("yes"),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("type"),
						knownvalue.StringExact("aggregate"),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("resource"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"sessions": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled":              knownvalue.Bool(true),
								"max_concurrent_limit": knownvalue.Int64Exact(1234),
							}),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	description = "test description"
	disable_override = "yes"
	type = "aggregate"
	resource = {
		sessions = {
			enabled = true
			max_concurrent_limit = 1234
		}
	}
}
`

func TestAccDosProtectionProfile_Classified(t *testing.T) {
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
				Config: dosProtectionProfile_Classified_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("type"),
						knownvalue.StringExact("classified"),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_Classified_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	type = "classified"
}
`

func TestAccDosProtectionProfile_ResourceSessions(t *testing.T) {
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
				Config: dosProtectionProfile_ResourceSessions_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("resource"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"sessions": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled":              knownvalue.Bool(true),
								"max_concurrent_limit": knownvalue.Int64Exact(1234),
							}),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_ResourceSessions_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	resource = {
		sessions = {
			enabled = true
			max_concurrent_limit = 1234
		}
	}
}
`

func TestAccDosProtectionProfile_FloodIcmp(t *testing.T) {
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
				Config: dosProtectionProfile_FloodIcmp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("flood"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"icmp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"activate_rate": knownvalue.Int64Exact(123),
									"alarm_rate":    knownvalue.Int64Exact(1234),
									"block": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duration": knownvalue.Int64Exact(12345),
									}),
									"maximal_rate": knownvalue.Int64Exact(123456),
								}),
							}),
							"icmpv6":   knownvalue.Null(),
							"other_ip": knownvalue.Null(),
							"tcp_syn":  knownvalue.Null(),
							"udp":      knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_FloodIcmp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	flood = {
		icmp = {
			enable = true
			red = {
				activate_rate = 123
				alarm_rate = 1234
				block = {
					duration = 12345
				}
				maximal_rate = 123456
			}
		}
	}
}
`

func TestAccDosProtectionProfile_FloodIcmpv6(t *testing.T) {
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
				Config: dosProtectionProfile_FloodIcmpv6_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("flood"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"icmpv6": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"activate_rate": knownvalue.Int64Exact(123),
									"alarm_rate":    knownvalue.Int64Exact(1234),
									"block": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duration": knownvalue.Int64Exact(12345),
									}),
									"maximal_rate": knownvalue.Int64Exact(123456),
								}),
							}),
							"icmp":     knownvalue.Null(),
							"other_ip": knownvalue.Null(),
							"tcp_syn":  knownvalue.Null(),
							"udp":      knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_FloodIcmpv6_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	flood = {
		icmpv6 = {
			enable = true
			red = {
				activate_rate = 123
				alarm_rate = 1234
				block = {
					duration = 12345
				}
				maximal_rate = 123456
			}
		}
	}
}
`

func TestAccDosProtectionProfile_FloodOtherIp(t *testing.T) {
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
				Config: dosProtectionProfile_FloodOtherIp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("flood"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"other_ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"activate_rate": knownvalue.Int64Exact(123),
									"alarm_rate":    knownvalue.Int64Exact(1234),
									"block": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duration": knownvalue.Int64Exact(12345),
									}),
									"maximal_rate": knownvalue.Int64Exact(123456),
								}),
							}),
							"icmp":     knownvalue.Null(),
							"icmpv6":   knownvalue.Null(),
							"tcp_syn":  knownvalue.Null(),
							"udp":      knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_FloodOtherIp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	flood = {
		other_ip = {
			enable = true
			red = {
				activate_rate = 123
				alarm_rate = 1234
				block = {
					duration = 12345
				}
				maximal_rate = 123456
			}
		}
	}
}
`

func TestAccDosProtectionProfile_FloodUdp(t *testing.T) {
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
				Config: dosProtectionProfile_FloodUdp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("flood"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"udp": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"activate_rate": knownvalue.Int64Exact(123),
									"alarm_rate":    knownvalue.Int64Exact(1234),
									"block": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duration": knownvalue.Int64Exact(12345),
									}),
									"maximal_rate": knownvalue.Int64Exact(123456),
								}),
							}),
							"icmp":     knownvalue.Null(),
							"icmpv6":   knownvalue.Null(),
							"other_ip": knownvalue.Null(),
							"tcp_syn":  knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_FloodUdp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	flood = {
		udp = {
			enable = true
			red = {
				activate_rate = 123
				alarm_rate = 1234
				block = {
					duration = 12345
				}
				maximal_rate = 123456
			}
		}
	}
}
`

func TestAccDosProtectionProfile_FloodTcpSyn_Red(t *testing.T) {
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
				Config: dosProtectionProfile_FloodTcpSyn_Red_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("flood"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"tcp_syn": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"activate_rate": knownvalue.Int64Exact(123),
									"alarm_rate":    knownvalue.Int64Exact(1234),
									"block": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duration": knownvalue.Int64Exact(12345),
									}),
									"maximal_rate": knownvalue.Int64Exact(123456),
								}),
								"syn_cookies": knownvalue.Null(),
							}),
							"icmp":     knownvalue.Null(),
							"icmpv6":   knownvalue.Null(),
							"other_ip": knownvalue.Null(),
							"udp":      knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_FloodTcpSyn_Red_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	flood = {
		tcp_syn = {
			enable = true
			red = {
				activate_rate = 123
				alarm_rate = 1234
				block = {
					duration = 12345
				}
				maximal_rate = 123456
			}
		}
	}
}
`

func TestAccDosProtectionProfile_FloodTcpSyn_SynCookies(t *testing.T) {
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
				Config: dosProtectionProfile_FloodTcpSyn_SynCookies_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_dos_protection_profile.example",
						tfjsonpath.New("flood"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"tcp_syn": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"syn_cookies": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"activate_rate": knownvalue.Int64Exact(123),
									"alarm_rate":    knownvalue.Int64Exact(1234),
									"block": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"duration": knownvalue.Int64Exact(12345),
									}),
									"maximal_rate": knownvalue.Int64Exact(123456),
								}),
								"red": knownvalue.Null(),
							}),
							"icmp":     knownvalue.Null(),
							"icmpv6":   knownvalue.Null(),
							"other_ip": knownvalue.Null(),
							"udp":      knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dosProtectionProfile_FloodTcpSyn_SynCookies_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_device_group" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_dos_protection_profile" "example" {
	depends_on = [panos_device_group.example]
	location = var.location
	name = var.prefix
	flood = {
		tcp_syn = {
			enable = true
			syn_cookies = {
				activate_rate = 123
				alarm_rate = 1234
				block = {
					duration = 12345
				}
				maximal_rate = 123456
			}
		}
	}
}
`
