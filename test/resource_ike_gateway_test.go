package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccIkeGateway_Basic(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_Basic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-gw1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("comment"),
						knownvalue.StringExact("ike gateway comment"),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("certificate").AtMapKey("allow_id_payload_mismatch"),
						knownvalue.Bool(true),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_ike_gateway.example",
					// 	tfjsonpath.New("local_address").AtMapKey("ip"),
					// 	knownvalue.StringExact("10.0.0.1/32"),
					// ),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("peer_address").AtMapKey("ip"),
						knownvalue.StringExact("10.10.0.1/32"),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("disabled"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("ipv6"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("local_id").AtMapKey("id"),
						knownvalue.StringExact("local.id.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("peer_id").AtMapKey("id"),
						knownvalue.StringExact("peer.id.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("certificate").AtMapKey("strict_validation_revocation"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("protocol_common").AtMapKey("fragmentation").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("protocol_common").AtMapKey("nat_traversal").AtMapKey("enable"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_PresharedKey(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_PresharedKey,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("pre_shared_key").AtMapKey("key"),
						knownvalue.StringExact("supersecret"),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_PeerAddressFqdn(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_PeerAddressFqdn,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("peer_address").AtMapKey("fqdn"),
						knownvalue.StringExact("peer.address.com"),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_PeerAddressDynamic(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_PeerAddressDynamic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("peer_address").AtMapKey("dynamic"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{}),
					),
				},
			},
		},
	})
}

// func TestAccIkeGateway_LocalAddressFloatingIp(t *testing.T) {
// 	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
// 	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

// 	resource.ParallelTest(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: ikeGatewayConfig_LocalAddressFloatingIp,
// 				ConfigVariables: map[string]config.Variable{
// 					"prefix": config.StringVariable(prefix),
// 				},
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					statecheck.ExpectKnownValue(
// 						"panos_ike_gateway.example",
// 						tfjsonpath.New("local_address").AtMapKey("floating_ip"),
// 						knownvalue.StringExact("1.1.1.1"),
// 					),
// 				},
// 			},
// 		},
// 	})
// }

func TestAccIkeGateway_Protocol(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_Protocol,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("protocol").AtMapKey("version"),
						knownvalue.StringExact("ikev2-preferred"),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("protocol_common").AtMapKey("passive_mode"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_AuthCertProfile(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_AuthCertProfile,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("certificate").AtMapKey("certificate_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert-prof", prefix)),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_AuthLocalCert(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_AuthLocalCert,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("certificate").AtMapKey("local_certificate").AtMapKey("name"),
						knownvalue.StringExact("my-cert"),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_ProtocolIkev1(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_ProtocolIkev1,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("protocol").AtMapKey("ikev1").AtMapKey("exchange_mode"),
						knownvalue.StringExact("aggressive"),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_ProtocolIkev2(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_ProtocolIkev2,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("protocol").AtMapKey("ikev2").AtMapKey("require_cookie"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestAccIkeGateway_PlaintextValueMissingRejected(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_PlaintextValueMissing,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile(`The attribute at path.+`),
			},
		},
	})
}

const ikeGatewayConfig_Basic = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  comment = "ike gateway comment"
  disabled = true
  ipv6 = true

  authentication = {
    certificate = {
      allow_id_payload_mismatch = true
	  strict_validation_revocation = true
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
    #ip = panos_ethernet_interface.example.layer3.ips.0.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }

  local_id = {
    id   = "local.id.com"
    type = "fqdn"
  }

  peer_id = {
    id   = "peer.id.com"
    type = "fqdn"
  }

  protocol_common = {
    fragmentation = {
      enable = true
    }
    nat_traversal = {
      enable = true
    }
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_PresharedKey = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    pre_shared_key = {
      key = "supersecret"
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_PeerAddressFqdn = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
      allow_id_payload_mismatch = true
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    fqdn = "peer.address.com"
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_PeerAddressDynamic = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
      allow_id_payload_mismatch = true
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    dynamic = {}
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

// const ikeGatewayConfig_LocalAddressFloatingIp = `
// variable "prefix" { type = string }

// resource "panos_ike_gateway" "example" {
//   location = { template = { name = panos_template.example.name } }

//   name    = format("%s-gw1", var.prefix)

//   authentication = {
//     certificate = {
//       allow_id_payload_mismatch = true
//     }
//   }

//   local_address = {
//     interface = panos_ethernet_interface.example.name
// 	floating_ip = "1.1.1.1"
//   }

//   peer_address = {
//     ip = "10.10.0.1/32"
//   }
// }

// resource "panos_ethernet_interface" "example" {
//   location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

//   name = "ethernet1/1"

//   layer3 = {
//     ips = [{ name = "10.0.0.1/32" }]
//   }
// }

// resource "panos_template" "example" {
//    location = { panorama = {} }
//    name     = format("%s-tmpl", var.prefix)
// }
// `

const ikeGatewayConfig_Protocol = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
      allow_id_payload_mismatch = true
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }

  protocol = {
    version = "ikev2-preferred"
  }

  protocol_common = {
    passive_mode = true
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_AuthCertProfile = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
      certificate_profile = panos_certificate_profile.example.name
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name     = format("%s-cert-prof", var.prefix)
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_AuthLocalCert = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
	  local_certificate = {
	    name = "my-cert"
	  }
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_ProtocolIkev1 = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {}
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }

  protocol = {
    ikev1 = {
	  exchange_mode = "aggressive"
	}
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_ProtocolIkev2 = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {}
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }

  protocol = {
    ikev2 = {
	  require_cookie = true
	}
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`

const ikeGatewayConfig_PlaintextValueMissing = `
variable "prefix" { type = string }

resource "panos_ike_gateway" "example" {
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    pre_shared_key = {
      key = "[PLAINTEXT-VALUE-MISSING]"
    }
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name, vsys = "vsys1" } }

  name = "ethernet1/1"

  layer3 = {
    ips = [{ name = "10.0.0.1/32" }]
  }
}

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}
`
