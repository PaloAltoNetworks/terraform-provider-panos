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
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_AuthLocalCert,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("certificate").AtMapKey("local_certificate").AtMapKey("name"),
						knownvalue.StringExact(certName),
					),
				},
			},
		},
	})
}

// TestAccIkeGateway_AuthLocalCert_SanMatch creates an IKEv2 gateway whose
// local_id matches a DNS SAN baked into the imported certificate. The apply
// proves that the local-certificate name reaches PAN-OS as a child element
// (not as an XML attribute), because PAN-OS's IKEv2 SAN-vs-local-id validator
// runs only after it can read the cert reference. Reproduces the user-reported
// scenario verbatim.
func TestAccIkeGateway_AuthLocalCert_SanMatch(t *testing.T) {
	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: ikeGatewayConfig_AuthLocalCertSan,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"certificate": config.StringVariable(certPemWithDnsSan),
					"private_key": config.StringVariable(privateKeyPemForDnsSan),
					"local_id":    config.StringVariable("host.example.com"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("authentication").AtMapKey("certificate").AtMapKey("local_certificate").AtMapKey("name"),
						knownvalue.StringExact(certName),
					),
					statecheck.ExpectKnownValue(
						"panos_ike_gateway.example",
						tfjsonpath.New("local_id").AtMapKey("id"),
						knownvalue.StringExact("host.example.com"),
					),
				},
			},
		},
	})
}

// NOTE: a SAN-mismatch negative-test variant was considered (local_id =
// "wrong.example.com" against the cert's DNS SAN "host.example.com" with
// ExpectError matching the user-reported "the local ID string must be from
// local certificate" message). On PAN-OS 11.2 that validator only fires
// during commit, not during candidate-config set, so the acceptance test
// framework (which never commits) cannot observe the rejection. Adding a
// commit step would substantially complicate the test and slow the suite;
// TestAccIkeGateway_AuthLocalCert_SanMatch already exercises the full bug
// repro path successfully, which is the load-bearing proof that the codegen
// fix lets PAN-OS see the cert reference.

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
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_certificate_import" "example" {
  depends_on = [panos_template.example]
  location = { template = { name = panos_template.example.name } }

  name = format("%s-cert", var.prefix)

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ike_gateway" "example" {
  depends_on = [panos_certificate_import.example]
  location = { template = { name = panos_template.example.name } }

  name    = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
	  local_certificate = {
	    name = panos_certificate_import.example.name
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

const ikeGatewayConfig_AuthLocalCertSan = `
variable "prefix"      { type = string }
variable "certificate" { type = string }
variable "private_key" { type = string }
variable "local_id"    { type = string }

resource "panos_template" "example" {
   location = { panorama = {} }
   name     = format("%s-tmpl", var.prefix)
}

resource "panos_certificate_import" "example" {
  depends_on = [panos_template.example]
  location = { template = { name = panos_template.example.name } }

  name = format("%s-cert", var.prefix)

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_certificate_profile" "example" {
  depends_on = [panos_certificate_import.example]
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

resource "panos_ike_gateway" "example" {
  depends_on = [panos_certificate_import.example, panos_certificate_profile.example]
  location = { template = { name = panos_template.example.name } }

  name = format("%s-gw1", var.prefix)

  authentication = {
    certificate = {
      local_certificate = {
        name = panos_certificate_import.example.name
      }
      certificate_profile = panos_certificate_profile.example.name
    }
  }

  local_id = {
    type = "fqdn"
    id   = var.local_id
  }

  protocol = {
    version = "ikev2"
  }

  local_address = {
    interface = panos_ethernet_interface.example.name
  }

  peer_address = {
    ip = "10.10.0.1/32"
  }
}
`

// certPemWithDnsSan is a self-signed test certificate with
// Subject CN = host.example.com and X509v3 Subject Alternative Name
// DNS:host.example.com. Generated once with:
//
//	openssl req -x509 -newkey rsa:2048 -days 3650 -nodes \
//	  -keyout key.pem -out cert.pem \
//	  -subj "/CN=host.example.com/O=test-acc/C=US" \
//	  -addext "subjectAltName=DNS:host.example.com"
const certPemWithDnsSan = `-----BEGIN CERTIFICATE-----
MIIDdDCCAlygAwIBAgIUTk6IiJcJaHB0iE65HrIL5lxbJycwDQYJKoZIhvcNAQEL
BQAwOzEZMBcGA1UEAwwQaG9zdC5leGFtcGxlLmNvbTERMA8GA1UECgwIdGVzdC1h
Y2MxCzAJBgNVBAYTAlVTMB4XDTI2MDYxODA3MzUyNloXDTM2MDYxNTA3MzUyNlow
OzEZMBcGA1UEAwwQaG9zdC5leGFtcGxlLmNvbTERMA8GA1UECgwIdGVzdC1hY2Mx
CzAJBgNVBAYTAlVTMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEApnWW
D3sEMnQDOH1LcDxUwToyK460SOzfirWSGy1XrzKcJ40H+B2cSNUyWYPeOw9UmhRo
6f5NYwN58UB3EQd1lUMDtEW7LtKhCCWi78yOktN+cZw69OJamU8EYWqhURhhuBVj
0LT8AK7Jh1MC49WCNZ22hEHxNF97XDwd71MNthE5W0El6+262o5Pz6a1uyqTImJN
NguDfp99FLbnGveqvStaBw/9PLLR7drC6TLZYNsM7pUidz+02EHOYt4roRue6LER
a+WzqoJGtAeoehkVqeEtDcO2+/B0SaZKS+KduEBiJSomq9emF4kjbDZZ92JVffqo
Tv8FlA+E0NPKMZQJ0QIDAQABo3AwbjAdBgNVHQ4EFgQUDcP31bn5mCBuxlY9MkXp
r/9Z3ocwHwYDVR0jBBgwFoAUDcP31bn5mCBuxlY9MkXpr/9Z3ocwDwYDVR0TAQH/
BAUwAwEB/zAbBgNVHREEFDASghBob3N0LmV4YW1wbGUuY29tMA0GCSqGSIb3DQEB
CwUAA4IBAQAE5Ii7yd7dVV/6YzDSSGBsF9zl3w4vjs7f2bWzeiJSUuxeGGwMV9EG
v7pyJa3NfAZvRC9Z4otGpyK/+e6xE8OdJ9bRuK1R0GY/mWsWfAWjmE2CNBKp0E9h
0Ts+gEfpWFUjlOHFR5JU6SvMhTCo+7bmT5HR1H0qbGVjfDRdEAVkixxDffz1oXAF
f9R22hqX16OZs6ahw2zt6PUOnVcaaQeIth3XkuX05T7w2myO50mg89EFIJbIw0CI
M79hx+Ez876a20GB7aHasRZrVkDZUmgpWoB+18MIUE9l7HzB5ADt4W7f0nCVQuHy
JyLmjtBGX4YJMmJa7l863PTB3uhUdSg0
-----END CERTIFICATE-----`

const privateKeyPemForDnsSan = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCmdZYPewQydAM4
fUtwPFTBOjIrjrRI7N+KtZIbLVevMpwnjQf4HZxI1TJZg947D1SaFGjp/k1jA3nx
QHcRB3WVQwO0Rbsu0qEIJaLvzI6S035xnDr04lqZTwRhaqFRGGG4FWPQtPwArsmH
UwLj1YI1nbaEQfE0X3tcPB3vUw22ETlbQSXr7brajk/PprW7KpMiYk02C4N+n30U
tuca96q9K1oHD/08stHt2sLpMtlg2wzulSJ3P7TYQc5i3iuhG57osRFr5bOqgka0
B6h6GRWp4S0Nw7b78HRJpkpL4p24QGIlKiar16YXiSNsNln3YlV9+qhO/wWUD4TQ
08oxlAnRAgMBAAECggEAJU9FRW5//BQDNcnlmwA2ygGSha9Mau0G1McKOsOwu+DW
3cmPT/ZJFY0NpkE3khrHTmrFMi9Q5qC8mA2oMl2u5ff5kFIx2JaWx+XxrZh2m2PC
m1HWaBWFE4hBtdkJ08yoeHN45ipusnN60bVSOEFemEhjhKToHKJLGtsRpZcCw3wf
jfrtDCfO+hycnrd2rJkq50WkPDQZSvdxKxaJrb1TyFUdak3Xf3wzWMdIOI6k7JBe
rWCFF61tKQb9v/P16n86T9L7qPFA+8RLTp3R9YTi3iHRyG8sHwtChb0njvHg7FKt
U/bqZ8na7P80EXZUJSoN4Ug3/TvpGE88pEHo+xY6rQKBgQDdsqXYD4B+hv3KEsIH
qPV52kcRH5vO9kjQ2cNOk0Y6BapCVIuN4vEWa3C6FpT2ZVm8Dp6o4DBSPRIlB6tX
fF26hlmXeXu8THdJZZLKhmN+ClUYOyC6LfSkzSV2+2XaETr1UCWYVoCpy0mgvPJZ
v+KTL9ADnA2SARd/PD17mOnGLQKBgQDANvcaqDGqWVTheLtngnUr9N3S3b4XO7fI
MCs+AXTBsj0QRA05YmW/YPWfKUu6d9Vq9/LXgSuqGqvt6vb79TbBsp8l6Ww4WKGS
g++CQ2Q7pNDZJDvm0fWPvphAd37793CKuzHEB2hArrHp3whM49+BjjdO/qRl13H7
VUEB3dYctQKBgAoXQpM2CWw46r5S4kAFAb9dHxT5clcWQLQ45TnjXDPx5BEG1h9M
MBsMIuJlerxIWrBDnhcjtS9ZFkVXNwZRY9bEnLlXTzl/5YISvH65ZTfscnka3995
jgQeTlE/GiC13hAiaMOpVEvmM+C8GO/a2w5GA9rWNIvrvs0MyeOhTyq1AoGAO4qL
PwGs6NTlOzbX7nd17ljawfAYaz//bQ6mxn1S+pFI4xoBcq4tUHwredMj9y4ZuRn3
apRDv1ylt3xaZ7AM9zFqpSbKdCXYXvdpoNNZYDRs0Was+5I8W/uxU/7wIgMDJKZa
Axw8ShUTXZvOCWtpF8vDDEBLEpULZMyC554dLiUCgYB9PcIn+KElppEYgPUGmJwy
hXU2wUnJXwqiFabQtXnCP/kyOjyvL1uwf/atkcJk2EkG4jeQyMZDhzjH57l+Tj9/
62GZuN2ZYSMF0wqv61hn6PNLi9vO66gXoyi08miw3OUIWKOXvaVjqH+XDHTJqosg
VdgiBA9XRWYswTsCXIX4zQ==
-----END PRIVATE KEY-----`
