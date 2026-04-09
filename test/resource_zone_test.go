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

func TestAccZone_Basic(t *testing.T) {
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
				Config: zone_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("enable_device_identification"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("enable_user_identification"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("enable_packet_buffer_protection"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("log_setting"),
						knownvalue.StringExact("log-settings-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("zone_protection_profile"),
						knownvalue.StringExact("zpp-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("device_acl").AtMapKey("include_list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.1.1.1"),
							knownvalue.StringExact("10.1.1.2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("device_acl").AtMapKey("exclude_list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.1.1.3"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("user_acl").AtMapKey("include_list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.2.1.1"),
							knownvalue.StringExact("10.2.1.2"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("user_acl").AtMapKey("exclude_list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("10.2.1.3"),
						}),
					),
				},
			},
		},
	})
}

const zone_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_zone_protection_profile" "zpp" {
	location = var.location
	name = "zpp-1"
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_zone_protection_profile.zpp]
	location = var.location
	name = var.prefix
	enable_device_identification = true
	enable_user_identification = true
	network = {
		layer3 = []
		enable_packet_buffer_protection = true
		log_setting = "log-settings-1"
		zone_protection_profile = "zpp-1"
	}
	device_acl = {
		include_list = ["10.1.1.1", "10.1.1.2"]
		exclude_list = ["10.1.1.3"]
	}
	user_acl = {
		include_list = ["10.2.1.1", "10.2.1.2"]
		exclude_list = ["10.2.1.3"]
	}
}
`

func TestAccZone_NetworkLayer2(t *testing.T) {
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
				Config: zone_NetworkLayer2_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("layer2"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
							knownvalue.StringExact("ethernet1/2"),
						}),
					),
				},
			},
		},
	})
}

const zone_NetworkLayer2_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_ethernet_interface" "eth1" {
	location = var.location
	name = "ethernet1/1"
	layer2 = {}
}

resource "panos_ethernet_interface" "eth2" {
	location = var.location
	name = "ethernet1/2"
	layer2 = {}
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_ethernet_interface.eth1, panos_ethernet_interface.eth2]
	location = var.location
	name = var.prefix
	network = {
		layer2 = ["ethernet1/1", "ethernet1/2"]
	}
}
`

func TestAccZone_NetworkLayer3(t *testing.T) {
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
				Config: zone_NetworkLayer3_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("layer3"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
							knownvalue.StringExact("ethernet1/2"),
						}),
					),
				},
			},
		},
	})
}

const zone_NetworkLayer3_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_ethernet_interface" "eth1" {
	location = var.location
	name = "ethernet1/1"
	layer3 = {}
}

resource "panos_ethernet_interface" "eth2" {
	location = var.location
	name = "ethernet1/2"
	layer3 = {}
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_ethernet_interface.eth1, panos_ethernet_interface.eth2]
	location = var.location
	name = var.prefix
	network = {
		layer3 = ["ethernet1/1", "ethernet1/2"]
	}
}
`

func TestAccZone_NetworkTap(t *testing.T) {
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
				Config: zone_NetworkTap_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("tap"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
						}),
					),
				},
			},
		},
	})
}

const zone_NetworkTap_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_ethernet_interface" "eth1" {
	location = var.location
	name = "ethernet1/1"
	tap = {}
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_ethernet_interface.eth1]
	location = var.location
	name = var.prefix
	network = {
		tap = ["ethernet1/1"]
	}
}
`

func TestAccZone_NetworkTunnel(t *testing.T) {
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
				Config: zone_NetworkTunnel_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("tunnel"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
				},
			},
		},
	})
}

const zone_NetworkTunnel_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_tunnel_interface" "tun1" {
	location = var.location
	name = "tunnel.1"
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_tunnel_interface.tun1]
	location = var.location
	name = var.prefix
	network = {
		tunnel = {}
	}
}
`

func TestAccZone_NetworkVirtualWire(t *testing.T) {
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
				Config: zone_NetworkVirtualWire_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("virtual_wire"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("ethernet1/1"),
						}),
					),
				},
			},
		},
	})
}

const zone_NetworkVirtualWire_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_ethernet_interface" "eth1" {
	location = var.location
	name = "ethernet1/1"
	virtual_wire = {}
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_ethernet_interface.eth1]
	location = var.location
	name = var.prefix
	network = {
		virtual_wire = ["ethernet1/1"]
	}
}
`

func TestAccZone_NetworkNetInspection(t *testing.T) {
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
				Config: zone_NetworkNetInspection_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("net_inspection"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const zone_NetworkNetInspection_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example]
	location = var.location
	name = var.prefix
	network = {
		layer3 = []
		net_inspection = true
	}
}
`

func TestAccZone_Layer3WithTunnelInterface(t *testing.T) {
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
				Config: zone_Layer3WithTunnelInterface_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_zone.example",
						tfjsonpath.New("network").AtMapKey("layer3"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("tunnel.1"),
						}),
					),
				},
			},
		},
	})
}

const zone_Layer3WithTunnelInterface_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
	location = { panorama = {} }
	name = var.prefix
}

resource "panos_tunnel_interface" "tun1" {
	location = var.location
	name = "tunnel.1"
}

resource "panos_zone" "example" {
	depends_on = [panos_template.example, panos_tunnel_interface.tun1]
	location = var.location
	name = var.prefix
	network = {
		layer3 = ["tunnel.1"]
	}
}
`