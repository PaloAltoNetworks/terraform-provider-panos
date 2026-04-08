
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

func TestAccLldpProfile_Basic(t *testing.T) {
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
				Config: lldpProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("transmit-receive"),
					),
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("snmp_syslog_notification"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("option_tlvs"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"port_description":   knownvalue.Bool(true),
							"system_name":        knownvalue.Bool(true),
							"system_description": knownvalue.Bool(true),
							"system_capabilities": knownvalue.Bool(true),
							"management_address": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const lldpProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_lldp_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  mode = "transmit-receive"
  snmp_syslog_notification = true
  option_tlvs = {
    port_description = true
    system_name = true
    system_description = true
    system_capabilities = true
  }
}
`

func TestAccLldpProfile_ManagementAddress(t *testing.T) {
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
				Config: lldpProfile_ManagementAddress_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("option_tlvs"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"management_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled": knownvalue.Bool(true),
								"iplist":  knownvalue.Null(),
							}),
							"port_description":    knownvalue.Null(),
							"system_name":         knownvalue.Null(),
							"system_description":  knownvalue.Null(),
							"system_capabilities": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const lldpProfile_ManagementAddress_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_lldp_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  option_tlvs = {
    management_address = {
      enabled = true
    }
  }
}
`

func TestAccLldpProfile_ManagementAddress_Ipv4(t *testing.T) {
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
				Config: lldpProfile_ManagementAddress_Ipv4_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("option_tlvs"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"management_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled": knownvalue.Bool(true),
								"iplist": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":      knownvalue.StringExact("ip1"),
										"interface": knownvalue.StringExact("ethernet1/1"),
										"ipv4":      knownvalue.StringExact("192.168.1.1"),
										"ipv6":      knownvalue.Null(),
									}),
								}),
							}),
							"port_description":    knownvalue.Null(),
							"system_name":         knownvalue.Null(),
							"system_description":  knownvalue.Null(),
							"system_capabilities": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const lldpProfile_ManagementAddress_Ipv4_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name} }

  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "192.168.1.1" }]
  }
}

resource "panos_lldp_profile" "example" {
  depends_on = [panos_ethernet_interface.example]
  location = var.location

  name = var.prefix
  option_tlvs = {
    management_address = {
      enabled = true
      iplist = [
        {
          name = "ip1"
          interface = panos_ethernet_interface.example.name
          ipv4 = "192.168.1.1"
        }
      ]
    }
  }
}
`

func TestAccLldpProfile_ManagementAddress_Ipv6(t *testing.T) {
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
				Config: lldpProfile_ManagementAddress_Ipv6_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("option_tlvs"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"management_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enabled": knownvalue.Bool(true),
								"iplist": knownvalue.ListExact([]knownvalue.Check{
									knownvalue.ObjectExact(map[string]knownvalue.Check{
										"name":      knownvalue.StringExact("ip1"),
										"interface": knownvalue.StringExact("ethernet1/1"),
										"ipv4":      knownvalue.Null(),
										"ipv6":      knownvalue.StringExact("2001:db8::1"),
									}),
								}),
							}),
							"port_description":    knownvalue.Null(),
							"system_name":         knownvalue.Null(),
							"system_description":  knownvalue.Null(),
							"system_capabilities": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const lldpProfile_ManagementAddress_Ipv6_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name} }

  name = "ethernet1/1"
  layer3 = {
    ipv6 = {
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_lldp_profile" "example" {
  depends_on = [panos_ethernet_interface.example]
  location = var.location

  name = var.prefix
  option_tlvs = {
    management_address = {
      enabled = true
      iplist = [
        {
          name = "ip1"
          interface = panos_ethernet_interface.example.name
          ipv6 = "2001:db8::1"
        }
      ]
    }
  }
}
`

func TestAccLldpProfile_Mode_TransmitOnly(t *testing.T) {
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
				Config: lldpProfile_Mode_TransmitOnly_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("transmit-only"),
					),
				},
			},
		},
	})
}

const lldpProfile_Mode_TransmitOnly_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_lldp_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  mode = "transmit-only"
}
`

func TestAccLldpProfile_Mode_ReceiveOnly(t *testing.T) {
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
				Config: lldpProfile_Mode_ReceiveOnly_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_lldp_profile.example",
						tfjsonpath.New("mode"),
						knownvalue.StringExact("receive-only"),
					),
				},
			},
		},
	})
}

const lldpProfile_Mode_ReceiveOnly_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_lldp_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  mode = "receive-only"
}
`


