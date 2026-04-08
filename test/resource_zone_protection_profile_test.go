
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

func TestAccZoneProtectionProfile_Basic(t *testing.T) {
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
				Config: zoneProtectionProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("asymmetric_path"),
						knownvalue.StringExact("bypass"),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("description"),
						knownvalue.StringExact("test description"),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_icmp_error"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_icmp_frag"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_icmp_large_packet"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_icmp_ping_zero_id"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_ip_frag"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_ip_spoof"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_loose_source_routing"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_malformed_option"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_overlapping_tcp_segment_mismatch"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_record_route"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_security"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_stream_id"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_strict_source_routing"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_tcp_split_handshake"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_tcp_syn_with_data"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_tcp_synack_with_data"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_timestamp"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("discard_unknown_option"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  description = "test description"
  asymmetric_path = "bypass"
  discard_icmp_error = true
  discard_icmp_frag = true
  discard_icmp_large_packet = true
  discard_icmp_ping_zero_id = true
  discard_ip_frag = true
  discard_ip_spoof = true
  discard_loose_source_routing = true
  discard_malformed_option = true
  discard_overlapping_tcp_segment_mismatch = true
  discard_record_route = true
  discard_security = true
  discard_stream_id = true
  discard_strict_source_routing = true
  discard_tcp_split_handshake = true
  discard_tcp_syn_with_data = true
  discard_tcp_synack_with_data = true
  discard_timestamp = true
  discard_unknown_option = true
}
`

func TestAccZoneProtectionProfile_Flood_Icmp(t *testing.T) {
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
				Config: zoneProtectionProfile_Flood_Icmp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("flood").AtMapKey("icmp"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"activate_rate": knownvalue.Int64Exact(100),
								"alarm_rate":    knownvalue.Int64Exact(200),
								"maximal_rate":  knownvalue.Int64Exact(300),
							}),
						}),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Flood_Icmp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  flood = {
    icmp = {
      enable = true
      red = {
        activate_rate = 100
        alarm_rate    = 200
        maximal_rate  = 300
      }
    }
  }
}
`

func TestAccZoneProtectionProfile_Flood_Icmpv6(t *testing.T) {
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
				Config: zoneProtectionProfile_Flood_Icmpv6_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("flood").AtMapKey("icmpv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"activate_rate": knownvalue.Int64Exact(100),
								"alarm_rate":    knownvalue.Int64Exact(200),
								"maximal_rate":  knownvalue.Int64Exact(300),
							}),
						}),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Flood_Icmpv6_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  flood = {
    icmpv6 = {
      enable = true
      red = {
        activate_rate = 100
        alarm_rate    = 200
        maximal_rate  = 300
      }
    }
  }
}
`

func TestAccZoneProtectionProfile_Flood_OtherIp(t *testing.T) {
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
				Config: zoneProtectionProfile_Flood_OtherIp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("flood").AtMapKey("other_ip"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"activate_rate": knownvalue.Int64Exact(100),
								"alarm_rate":    knownvalue.Int64Exact(200),
								"maximal_rate":  knownvalue.Int64Exact(300),
							}),
						}),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Flood_OtherIp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  flood = {
    other_ip = {
      enable = true
      red = {
        activate_rate = 100
        alarm_rate    = 200
        maximal_rate  = 300
      }
    }
  }
}
`

func TestAccZoneProtectionProfile_Flood_TcpSyn_Red(t *testing.T) {
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
				Config: zoneProtectionProfile_Flood_TcpSyn_Red_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("flood").AtMapKey("tcp_syn"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"activate_rate": knownvalue.Int64Exact(100),
								"alarm_rate":    knownvalue.Int64Exact(200),
								"maximal_rate":  knownvalue.Int64Exact(300),
							}),
							"syn_cookies": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Flood_TcpSyn_Red_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  flood = {
    tcp_syn = {
      enable = true
      red = {
        activate_rate = 100
        alarm_rate    = 200
        maximal_rate  = 300
      }
    }
  }
}
`

func TestAccZoneProtectionProfile_Flood_TcpSyn_SynCookies(t *testing.T) {
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
				Config: zoneProtectionProfile_Flood_TcpSyn_SynCookies_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("flood").AtMapKey("tcp_syn"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"red":    knownvalue.Null(),
							"syn_cookies": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"activate_rate": knownvalue.Int64Exact(100),
								"alarm_rate":    knownvalue.Int64Exact(200),
								"maximal_rate":  knownvalue.Int64Exact(300),
							}),
						}),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Flood_TcpSyn_SynCookies_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  flood = {
    tcp_syn = {
      enable = true
      syn_cookies = {
        activate_rate = 100
        alarm_rate    = 200
        maximal_rate  = 300
      }
    }
  }
}
`

func TestAccZoneProtectionProfile_Flood_Udp(t *testing.T) {
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
				Config: zoneProtectionProfile_Flood_Udp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone_protection_profile.example",
						tfjsonpath.New("flood").AtMapKey("udp"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enable": knownvalue.Bool(true),
							"red": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"activate_rate": knownvalue.Int64Exact(100),
								"alarm_rate":    knownvalue.Int64Exact(200),
								"maximal_rate":  knownvalue.Int64Exact(300),
							}),
						}),
					),
				},
			},
		},
	})
}

const zoneProtectionProfile_Flood_Udp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone_protection_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  flood = {
    udp = {
      enable = true
      red = {
        activate_rate = 100
        alarm_rate    = 200
        maximal_rate  = 300
      }
    }
  }
}
`

