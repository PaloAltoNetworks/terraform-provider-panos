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

func TestAccSslDecrypt_Basic(t *testing.T) {
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
				Config: sslDecrypt_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_decrypt.example",
						tfjsonpath.New("disabled_ssl_exclude_cert_from_predefined"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("cert1"),
							knownvalue.StringExact("cert2"),
						}),
					),
				},
			},
		},
	})
}

const sslDecrypt_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ssl_decrypt" "example" {
  depends_on = [panos_template.example]
  location = var.location

  disabled_ssl_exclude_cert_from_predefined = ["cert1", "cert2"]
}
`

func TestAccSslDecrypt_ForwardCertificates(t *testing.T) {
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
				Config: sslDecrypt_ForwardCertificates_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_decrypt.example",
						tfjsonpath.New("forward_trust_certificate_rsa"),
						knownvalue.StringExact(fmt.Sprintf("%s-trust", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_ssl_decrypt.example",
						tfjsonpath.New("forward_untrust_certificate_rsa"),
						knownvalue.StringExact(fmt.Sprintf("%s-untrust", prefix)),
					),
				},
			},
		},
	})
}

const sslDecrypt_ForwardCertificates_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "trust" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-trust"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_certificate_import" "untrust" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-untrust"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_decrypt" "example" {
  depends_on = [panos_template.example]
  location = var.location

  forward_trust_certificate_rsa = panos_certificate_import.trust.name
  forward_untrust_certificate_rsa = panos_certificate_import.untrust.name
}
`

func TestAccSslDecrypt_SslExcludeCert(t *testing.T) {
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
				Config: sslDecrypt_SslExcludeCert_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_decrypt.example",
						tfjsonpath.New("ssl_exclude_cert"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":        knownvalue.StringExact("cert-exclude1"),
								"description": knownvalue.StringExact("First excluded certificate"),
								"exclude":     knownvalue.Bool(true),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":        knownvalue.StringExact("cert-exclude2"),
								"description": knownvalue.StringExact("Second excluded certificate"),
								"exclude":     knownvalue.Bool(false),
							}),
						}),
					),
				},
			},
		},
	})
}

const sslDecrypt_SslExcludeCert_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ssl_decrypt" "example" {
  depends_on = [panos_template.example]
  location = var.location

  ssl_exclude_cert = [
    {
      name = "cert-exclude1"
      description = "First excluded certificate"
      exclude = true
    },
    {
      name = "cert-exclude2"
      description = "Second excluded certificate"
      exclude = false
    }
  ]
}
`

func TestAccSslDecrypt_TrustedRootCa(t *testing.T) {
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
				Config: sslDecrypt_TrustedRootCa_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_decrypt.example",
						tfjsonpath.New("trusted_root_ca"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-trusted-ca1", prefix)),
							knownvalue.StringExact(fmt.Sprintf("%s-trusted-ca2", prefix)),
						}),
					),
				},
			},
		},
	})
}

const sslDecrypt_TrustedRootCa_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "trusted_ca1" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-trusted-ca1"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_certificate_import" "trusted_ca2" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-trusted-ca2"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_decrypt" "example" {
  depends_on = [panos_template.example]
  location = var.location

  trusted_root_ca = [
    panos_certificate_import.trusted_ca1.name,
    panos_certificate_import.trusted_ca2.name
  ]
}
`

func TestAccSslDecrypt_RootCaExcludeList(t *testing.T) {
	t.Parallel()
	t.Skip("root_ca_exclude_list requires actual CA certificates, not regular certificates")

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
				Config: sslDecrypt_RootCaExcludeList_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
					"private_key": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_ssl_decrypt.example",
						tfjsonpath.New("root_ca_exclude_list"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact(fmt.Sprintf("%s-root-ca1", prefix)),
							knownvalue.StringExact(fmt.Sprintf("%s-root-ca2", prefix)),
						}),
					),
				},
			},
		},
	})
}

const sslDecrypt_RootCaExcludeList_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "private_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_import" "root_ca1" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-root-ca1"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_certificate_import" "root_ca2" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-root-ca2"

  local = {
    pem = {
      certificate = var.certificate
      private_key = var.private_key
    }
  }
}

resource "panos_ssl_decrypt" "example" {
  depends_on = [panos_template.example]
  location = var.location

  root_ca_exclude_list = [
    panos_certificate_import.root_ca1.name,
    panos_certificate_import.root_ca2.name
  ]
}
`
