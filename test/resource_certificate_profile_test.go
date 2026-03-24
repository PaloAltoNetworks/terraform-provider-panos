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

func TestAccCertificateProfile_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: certificateProfile_Basic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("block_expired_certificate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("block_timeout_certificate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("block_unauthenticated_certificate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("block_unknown_certificate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("certificate_status_timeout"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("crl_receive_timeout"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("domain"),
						knownvalue.StringExact("example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("ocsp_exclude_nonce"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("ocsp_receive_timeout"),
						knownvalue.Int64Exact(10),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("use_crl"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("use_ocsp"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("username_field"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"subject":     knownvalue.StringExact("common-name"),
							"subject_alt": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const certificateProfile_Basic = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  block_expired_certificate = true
  block_timeout_certificate = true
  block_unauthenticated_certificate = true
  block_unknown_certificate = true
  certificate_status_timeout = 10
  crl_receive_timeout = 10
  domain = "example.com"
  ocsp_exclude_nonce = true
  ocsp_receive_timeout = 10
  use_crl = true
  use_ocsp = true
  username_field = {
    subject = "common-name"
  }
}
`

func TestAccCertificateProfile_UsernameFieldSubjectAltEmail(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: certificateProfile_UsernameFieldSubjectAltEmail,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("username_field"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"subject":     knownvalue.Null(),
							"subject_alt": knownvalue.StringExact("email"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const certificateProfile_UsernameFieldSubjectAltEmail = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  username_field = {
    subject_alt = "email"
  }
}
`

func TestAccCertificateProfile_UsernameFieldSubjectAltPrincipalName(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: certificateProfile_UsernameFieldSubjectAltPrincipalName,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("username_field"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"subject":     knownvalue.Null(),
							"subject_alt": knownvalue.StringExact("principal-name"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const certificateProfile_UsernameFieldSubjectAltPrincipalName = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  username_field = {
    subject_alt = "principal-name"
  }
}
`

func TestAccCertificateProfile_WithCertificate(t *testing.T) {
	t.Parallel()
	t.Skip("This test is skipped because it requires CA and OCSP certificates to be present on the device.")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: certificateProfile_WithCertificate,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("certificate"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                    knownvalue.StringExact("ca_cert_name"),
								"default_ocsp_url":        knownvalue.StringExact("http://ocsp.example.com"),
								"ocsp_verify_certificate": knownvalue.StringExact("ocsp_signer_cert"),
								"template_name":           knownvalue.StringExact("cert_template"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const certificateProfile_WithCertificate = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  certificate = [{
    name = "ca_cert_name"
    default_ocsp_url = "http://ocsp.example.com"
    ocsp_verify_certificate = "ocsp_signer_cert"
    template_name = "cert_template"
  }]
}
`

func TestAccCertificateProfile_TemplateLocation(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: certificateProfile_TemplateLocation,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const certificateProfile_TemplateLocation = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix
}
`
