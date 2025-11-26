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

func TestAccSslTlsServiceProfile_Basic(t *testing.T) {
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
				Config: sslTlsServiceProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("certificate"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("protocol_settings"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const sslTlsServiceProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "cert" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_tls_service_profile" "example" {
  depends_on = [panos_certificate_import.cert]
  location = var.location

  name = var.prefix
  certificate = panos_certificate_import.cert.name
}
`

func TestAccSslTlsServiceProfile_ProtocolSettings_EncryptionAlgorithms(t *testing.T) {
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
				Config: sslTlsServiceProfile_ProtocolSettings_EncryptionAlgorithms_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("protocol_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"allow_algorithm_3des":       knownvalue.Null(),
							"allow_algorithm_aes_128_cbc": knownvalue.Bool(true),
							"allow_algorithm_aes_128_gcm": knownvalue.Bool(true),
							"allow_algorithm_aes_256_cbc": knownvalue.Bool(true),
							"allow_algorithm_aes_256_gcm": knownvalue.Bool(true),
							"allow_algorithm_rc4":         knownvalue.Null(),
							"allow_algorithm_dhe":         knownvalue.Null(),
							"allow_algorithm_ecdhe":       knownvalue.Null(),
							"allow_algorithm_rsa":         knownvalue.Null(),
							"allow_authentication_sha1":   knownvalue.Null(),
							"allow_authentication_sha256": knownvalue.Null(),
							"allow_authentication_sha384": knownvalue.Null(),
							"max_version":                 knownvalue.StringExact("tls1-2"),
							"min_version":                 knownvalue.StringExact("tls1-0"),
						}),
					),
				},
			},
		},
	})
}

const sslTlsServiceProfile_ProtocolSettings_EncryptionAlgorithms_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "cert" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_tls_service_profile" "example" {
  depends_on = [panos_certificate_import.cert]
  location = var.location

  name = var.prefix
  certificate = panos_certificate_import.cert.name

  protocol_settings = {
    allow_algorithm_aes_128_cbc = true
    allow_algorithm_aes_128_gcm = true
    allow_algorithm_aes_256_cbc = true
    allow_algorithm_aes_256_gcm = true
    max_version = "tls1-2"
  }
}
`

func TestAccSslTlsServiceProfile_ProtocolSettings_KeyExchangeAlgorithms(t *testing.T) {
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
				Config: sslTlsServiceProfile_ProtocolSettings_KeyExchangeAlgorithms_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("protocol_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"allow_algorithm_3des":       knownvalue.Null(),
							"allow_algorithm_aes_128_cbc": knownvalue.Null(),
							"allow_algorithm_aes_128_gcm": knownvalue.Null(),
							"allow_algorithm_aes_256_cbc": knownvalue.Null(),
							"allow_algorithm_aes_256_gcm": knownvalue.Null(),
							"allow_algorithm_rc4":         knownvalue.Null(),
							"allow_algorithm_dhe":         knownvalue.Bool(true),
							"allow_algorithm_ecdhe":       knownvalue.Bool(true),
							"allow_algorithm_rsa":         knownvalue.Bool(true),
							"allow_authentication_sha1":   knownvalue.Null(),
							"allow_authentication_sha256": knownvalue.Null(),
							"allow_authentication_sha384": knownvalue.Null(),
							"max_version":                 knownvalue.StringExact("tls1-2"),
							"min_version":                 knownvalue.StringExact("tls1-0"),
						}),
					),
				},
			},
		},
	})
}

const sslTlsServiceProfile_ProtocolSettings_KeyExchangeAlgorithms_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "cert" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_tls_service_profile" "example" {
  depends_on = [panos_certificate_import.cert]
  location = var.location

  name = var.prefix
  certificate = panos_certificate_import.cert.name

  protocol_settings = {
    allow_algorithm_dhe = true
    allow_algorithm_ecdhe = true
    allow_algorithm_rsa = true
    max_version = "tls1-2"
  }
}
`

func TestAccSslTlsServiceProfile_ProtocolSettings_AuthenticationAlgorithms(t *testing.T) {
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
				Config: sslTlsServiceProfile_ProtocolSettings_AuthenticationAlgorithms_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("protocol_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"allow_algorithm_3des":       knownvalue.Null(),
							"allow_algorithm_aes_128_cbc": knownvalue.Null(),
							"allow_algorithm_aes_128_gcm": knownvalue.Null(),
							"allow_algorithm_aes_256_cbc": knownvalue.Null(),
							"allow_algorithm_aes_256_gcm": knownvalue.Null(),
							"allow_algorithm_rc4":         knownvalue.Null(),
							"allow_algorithm_dhe":         knownvalue.Null(),
							"allow_algorithm_ecdhe":       knownvalue.Null(),
							"allow_algorithm_rsa":         knownvalue.Null(),
							"allow_authentication_sha1":   knownvalue.Bool(true),
							"allow_authentication_sha256": knownvalue.Bool(true),
							"allow_authentication_sha384": knownvalue.Bool(true),
							"max_version":                 knownvalue.StringExact("tls1-2"),
							"min_version":                 knownvalue.StringExact("tls1-0"),
						}),
					),
				},
			},
		},
	})
}

const sslTlsServiceProfile_ProtocolSettings_AuthenticationAlgorithms_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "cert" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_tls_service_profile" "example" {
  depends_on = [panos_certificate_import.cert]
  location = var.location

  name = var.prefix
  certificate = panos_certificate_import.cert.name

  protocol_settings = {
    allow_authentication_sha1 = true
    allow_authentication_sha256 = true
    allow_authentication_sha384 = true
    max_version = "tls1-2"
  }
}
`

func TestAccSslTlsServiceProfile_ProtocolSettings_TlsVersions(t *testing.T) {
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
				Config: sslTlsServiceProfile_ProtocolSettings_TlsVersions_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("protocol_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"allow_algorithm_3des":       knownvalue.Null(),
							"allow_algorithm_aes_128_cbc": knownvalue.Null(),
							"allow_algorithm_aes_128_gcm": knownvalue.Null(),
							"allow_algorithm_aes_256_cbc": knownvalue.Null(),
							"allow_algorithm_aes_256_gcm": knownvalue.Null(),
							"allow_algorithm_rc4":         knownvalue.Null(),
							"allow_algorithm_dhe":         knownvalue.Null(),
							"allow_algorithm_ecdhe":       knownvalue.Null(),
							"allow_algorithm_rsa":         knownvalue.Null(),
							"allow_authentication_sha1":   knownvalue.Null(),
							"allow_authentication_sha256": knownvalue.Null(),
							"allow_authentication_sha384": knownvalue.Null(),
							"max_version":                 knownvalue.StringExact("tls1-2"),
							"min_version":                 knownvalue.StringExact("tls1-1"),
						}),
					),
				},
			},
		},
	})
}

const sslTlsServiceProfile_ProtocolSettings_TlsVersions_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "cert" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_tls_service_profile" "example" {
  depends_on = [panos_certificate_import.cert]
  location = var.location

  name = var.prefix
  certificate = panos_certificate_import.cert.name

  protocol_settings = {
    min_version = "tls1-1"
    max_version = "tls1-2"
  }
}
`

func TestAccSslTlsServiceProfile_ProtocolSettings_Complete(t *testing.T) {
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
				Config: sslTlsServiceProfile_ProtocolSettings_Complete_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_tls_service_profile.example",
						tfjsonpath.New("protocol_settings"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"allow_algorithm_3des":       knownvalue.Null(),
							"allow_algorithm_aes_128_cbc": knownvalue.Bool(true),
							"allow_algorithm_aes_128_gcm": knownvalue.Bool(true),
							"allow_algorithm_aes_256_cbc": knownvalue.Bool(true),
							"allow_algorithm_aes_256_gcm": knownvalue.Bool(true),
							"allow_algorithm_rc4":         knownvalue.Null(),
							"allow_algorithm_dhe":         knownvalue.Bool(true),
							"allow_algorithm_ecdhe":       knownvalue.Bool(true),
							"allow_algorithm_rsa":         knownvalue.Bool(true),
							"allow_authentication_sha1":   knownvalue.Bool(true),
							"allow_authentication_sha256": knownvalue.Bool(true),
							"allow_authentication_sha384": knownvalue.Bool(true),
							"max_version":                 knownvalue.StringExact("tls1-2"),
							"min_version":                 knownvalue.StringExact("tls1-1"),
						}),
					),
				},
			},
		},
	})
}

const sslTlsServiceProfile_ProtocolSettings_Complete_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "cert" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_tls_service_profile" "example" {
  depends_on = [panos_certificate_import.cert]
  location = var.location

  name = var.prefix
  certificate = panos_certificate_import.cert.name

  protocol_settings = {
    allow_algorithm_aes_128_cbc = true
    allow_algorithm_aes_128_gcm = true
    allow_algorithm_aes_256_cbc = true
    allow_algorithm_aes_256_gcm = true
    allow_algorithm_dhe = true
    allow_algorithm_ecdhe = true
    allow_algorithm_rsa = true
    allow_authentication_sha1 = true
    allow_authentication_sha256 = true
    allow_authentication_sha384 = true
    min_version = "tls1-1"
    max_version = "tls1-2"
  }
}
`
