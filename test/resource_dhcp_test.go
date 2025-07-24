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

func TestAccDhcp_Relay_Ip(t *testing.T) {
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
				Config: dhcp_Relay_Ip_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("relay").AtMapKey("ip"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"server": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("10.0.0.1"),
								knownvalue.StringExact("10.0.0.2"),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Relay_Ip_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  relay = {
    ip = {
      enabled = true
      server = ["10.0.0.1", "10.0.0.2"]
    }
  }
}
`

func TestAccDhcp_Relay_Ipv6(t *testing.T) {
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
				Config: dhcp_Relay_Ipv6_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("relay").AtMapKey("ipv6"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"enabled": knownvalue.Bool(true),
							"server": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":      knownvalue.StringExact("2001:db8::1"),
									"interface": knownvalue.StringExact("ethernet1/2"),
								}),
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":      knownvalue.StringExact("2001:db8::2"),
									"interface": knownvalue.Null(),
								}),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Relay_Ipv6_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_interface" "example2" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/2"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  relay = {
    ipv6 = {
      enabled = true
      server = [
        {
          name = "2001:db8::1"
          interface = panos_ethernet_interface.example2.name
        },
        {
          name = "2001:db8::2"
        }
      ]
    }
  }
}
`

func TestAccDhcp_Server_Option_Lease_Timeout(t *testing.T) {
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
				Config: dhcp_Server_Option_Lease_Timeout_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("lease"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"timeout":   knownvalue.Int64Exact(720),
							"unlimited": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Lease_Timeout_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      lease = {
        timeout = 720
      }
    }
  }
}
`

func TestAccDhcp_Server_Option_Lease_Unlimited(t *testing.T) {
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
				Config: dhcp_Server_Option_Lease_Unlimited_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("lease"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"unlimited": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"timeout":   knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Lease_Unlimited_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      lease = {
        unlimited = {}
      }
    }
  }
}
`

func TestAccDhcp_Server_Option_Dns(t *testing.T) {
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
				Config: dhcp_Server_Option_Dns_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("dns"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"primary":   knownvalue.StringExact("8.8.8.8"),
							"secondary": knownvalue.StringExact("8.8.4.4"),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Dns_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      dns = {
        primary = "8.8.8.8"
        secondary = "8.8.4.4"
      }
    }
  }
}
`

func TestAccDhcp_Server_Option_DnsSuffix(t *testing.T) {
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
				Config: dhcp_Server_Option_DnsSuffix_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("dns_suffix"),
						knownvalue.StringExact("example.com"),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_DnsSuffix_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      dns_suffix = "example.com"
    }
  }
}
`

func TestAccDhcp_Server_Option_Gateway(t *testing.T) {
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
				Config: dhcp_Server_Option_Gateway_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("gateway"),
						knownvalue.StringExact("192.168.1.1"),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Gateway_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      gateway = "192.168.1.1"
    }
  }
}
`

func TestAccDhcp_Server_Option_Inheritance(t *testing.T) {
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
				Config: dhcp_Server_Option_Inheritance_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("inheritance").AtMapKey("source"),
						knownvalue.StringExact("ethernet1/2"),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Inheritance_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_ethernet_interface" "example2" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/2"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      inheritance = {
        source = panos_ethernet_interface.example2.name
      }
    }
  }
}
`

func TestAccDhcp_Server_Option_Nis(t *testing.T) {
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
				Config: dhcp_Server_Option_Nis_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("nis"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"primary":   knownvalue.StringExact("192.168.1.10"),
							"secondary": knownvalue.StringExact("192.168.1.11"),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Nis_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      nis = {
        primary = "192.168.1.10"
        secondary = "192.168.1.11"
      }
    }
  }
}
`

func TestAccDhcp_Server_Option_Ntp(t *testing.T) {
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
				Config: dhcp_Server_Option_Ntp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("ntp"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"primary":   knownvalue.StringExact("192.168.1.20"),
							"secondary": knownvalue.StringExact("192.168.1.21"),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Ntp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      ntp = {
        primary = "192.168.1.20"
        secondary = "192.168.1.21"
      }
    }
  }
}
`

func TestAccDhcp_Server_Option_Pop3Server(t *testing.T) {
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
				Config: dhcp_Server_Option_Pop3Server_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("pop3_server"),
						knownvalue.StringExact("192.168.1.30"),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Pop3Server_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      pop3_server = "192.168.1.30"
    }
  }
}
`

func TestAccDhcp_Server_Option_SmtpServer(t *testing.T) {
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
				Config: dhcp_Server_Option_SmtpServer_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("smtp_server"),
						knownvalue.StringExact("192.168.1.25"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("mode"),
						knownvalue.StringExact("enabled"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("probe_ip"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_SmtpServer_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    mode = "enabled"
    probe_ip = true
    option = {
      smtp_server = "192.168.1.25"
    }
  }
}
`

func TestAccDhcp_Server_Option_SubnetMask(t *testing.T) {
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
				Config: dhcp_Server_Option_SubnetMask_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("subnet_mask"),
						knownvalue.StringExact("255.255.255.0"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("reserved"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":        knownvalue.StringExact("192.168.1.250"),
								"mac":         knownvalue.StringExact("aa:bb:cc:dd:ee:ff"),
								"description": knownvalue.StringExact("Reserved IP for Printer"),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_SubnetMask_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      subnet_mask = "255.255.255.0"
    }
    reserved = [
      {
        name = "192.168.1.250"
        mac = "aa:bb:cc:dd:ee:ff"
        description = "Reserved IP for Printer"
      }
    ]
  }
}
`

func TestAccDhcp_Server_Option_UserDefined_Ip(t *testing.T) {
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
				Config: dhcp_Server_Option_UserDefined_Ip_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("user_defined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("custom_ip_option"),
								"code": knownvalue.Int64Exact(43),
								"ip":   knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.1"), knownvalue.StringExact("10.0.0.2")}),

								"inherited":               knownvalue.Bool(false),
								"vendor_class_identifier": knownvalue.StringExact("Custom VCI"),
								"ascii":                   knownvalue.Null(),
								"hex":                     knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_UserDefined_Ip_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      user_defined = [
        {
          name = "custom_ip_option"
          code = 43
          ip = ["10.0.0.1", "10.0.0.2"]
          inherited = false
          vendor_class_identifier = "Custom VCI"
        }
      ]
    }
  }
}
`

func TestAccDhcp_Server_Option_UserDefined_Ascii(t *testing.T) {
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
				Config: dhcp_Server_Option_UserDefined_Ascii_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("user_defined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("custom_ascii_option"),
								"code":                    knownvalue.Int64Exact(201),
								"ascii":                   knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("custom option")}),
								"inherited":               knownvalue.Bool(false),
								"vendor_class_identifier": knownvalue.Null(),
								"ip":                      knownvalue.Null(),
								"hex":                     knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_UserDefined_Ascii_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      user_defined = [
        {
          name = "custom_ascii_option"
          code = 201
          ascii = ["custom option"]
          inherited = false
        }
      ]
    }
  }
}
`

func TestAccDhcp_Server_Option_UserDefined_Hex(t *testing.T) {
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
				Config: dhcp_Server_Option_UserDefined_Hex_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("user_defined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("custom_hex_option"),
								"code":                    knownvalue.Int64Exact(202),
								"hex":                     knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("0A0B0C")}),
								"inherited":               knownvalue.Bool(false),
								"vendor_class_identifier": knownvalue.Null(),
								"ip":                      knownvalue.Null(),
								"ascii":                   knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_UserDefined_Hex_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      user_defined = [
        {
          name = "custom_hex_option"
          code = 202
          hex = ["0A0B0C"]
          inherited = false
        }
      ]
    }
  }
}
`

func TestAccDhcp_Server_Option_Wins(t *testing.T) {
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
				Config: dhcp_Server_Option_Wins_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact("ethernet1/1"),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("ip_pool"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("192.168.1.0/24"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("wins"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"primary":   knownvalue.StringExact("192.168.1.10"),
							"secondary": knownvalue.StringExact("192.168.1.11"),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_Wins_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      wins = {
        primary = "192.168.1.10"
        secondary = "192.168.1.11"
      }
    }
  }
}
`

func TestAccDhcp_Server_Mode_Disabled(t *testing.T) {
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
				Config: dhcp_Server_Mode_Disabled_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("mode"),
						knownvalue.StringExact("disabled"),
					),
				},
			},
		},
	})
}

const dhcp_Server_Mode_Disabled_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    mode = "disabled"
  }
}
`

func TestAccDhcp_Server_Mode_Auto(t *testing.T) {
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
				Config: dhcp_Server_Mode_Auto_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("mode"),
						knownvalue.StringExact("auto"),
					),
				},
			},
		},
	})
}

const dhcp_Server_Mode_Auto_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    mode = "auto"
  }
}
`

func TestAccDhcp_Server_ProbeIp_False(t *testing.T) {
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
				Config: dhcp_Server_ProbeIp_False_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("probe_ip"),
						knownvalue.Bool(false),
					),
				},
			},
		},
	})
}

const dhcp_Server_ProbeIp_False_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    probe_ip = false
  }
}
`

func TestAccDhcp_Server_Option_UserDefined_Inherited(t *testing.T) {
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
				Config: dhcp_Server_Option_UserDefined_Inherited_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("user_defined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("custom_inherited_option"),
								"code":                    knownvalue.Int64Exact(43),
								"ip":                      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.1")}),
								"inherited":               knownvalue.Bool(true),
								"vendor_class_identifier": knownvalue.StringExact("some-vci"),
								"ascii":                   knownvalue.Null(),
								"hex":                     knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_UserDefined_Inherited_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      user_defined = [
        {
          name = "custom_inherited_option"
          code = 43
          ip = ["10.0.0.1"]
          inherited = true
          vendor_class_identifier = "some-vci"
        }
      ]
    }
  }
}
`

func TestAccDhcp_Server_Option_UserDefined_Ascii_Vci(t *testing.T) {
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
				Config: dhcp_Server_Option_UserDefined_Ascii_Vci_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("user_defined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("custom_ascii_vci_option"),
								"code":                    knownvalue.Int64Exact(43),
								"ascii":                   knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("custom option")}),
								"inherited":               knownvalue.Bool(false),
								"vendor_class_identifier": knownvalue.StringExact("Custom VCI ASCII"),
								"ip":                      knownvalue.Null(),
								"hex":                     knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_UserDefined_Ascii_Vci_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      user_defined = [
        {
          name = "custom_ascii_vci_option"
          code = 43
          ascii = ["custom option"]
          inherited = false
          vendor_class_identifier = "Custom VCI ASCII"
        }
      ]
    }
  }
}
`

func TestAccDhcp_Server_Option_UserDefined_Hex_Vci(t *testing.T) {
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
				Config: dhcp_Server_Option_UserDefined_Hex_Vci_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_dhcp.example",
						tfjsonpath.New("server").AtMapKey("option").AtMapKey("user_defined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("custom_hex_vci_option"),
								"code":                    knownvalue.Int64Exact(43),
								"hex":                     knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("0A0B0C")}),
								"inherited":               knownvalue.Bool(false),
								"vendor_class_identifier": knownvalue.StringExact("Custom VCI HEX"),
								"ip":                      knownvalue.Null(),
								"ascii":                   knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const dhcp_Server_Option_UserDefined_Hex_Vci_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = {
    panorama = {}
  }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = "ethernet1/1"
  layer3 = {}
}

resource "panos_dhcp" "example" {
  location = var.location
  name = panos_ethernet_interface.example.name

  server = {
    ip_pool = ["192.168.1.0/24"]
    option = {
      user_defined = [
        {
          name = "custom_hex_vci_option"
          code = 43
          hex = ["0A0B0C"]
          inherited = false
          vendor_class_identifier = "Custom VCI HEX"
        }
      ]
    }
  }
}
`
