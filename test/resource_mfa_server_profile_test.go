package provider_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccMfaServerProfile_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_Basic,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.StringExact("test-cert-profile"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("okta-adaptive-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-api-host"),
								"value": knownvalue.StringExact("api.okta.example.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-baseuri"),
								"value": knownvalue.StringExact("/api/v1"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-token"),
								"value": knownvalue.StringExact("test-token-123"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-org"),
								"value": knownvalue.StringExact("test-org"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_Basic = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = "test-cert-profile"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_cert_profile = panos_certificate_profile.example.name
  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [
    {
      name = "okta-api-host"
      value = "api.okta.example.com"
    },
    {
      name = "okta-baseuri"
      value = "/api/v1"
    },
    {
      name = "okta-token"
      value = "test-token-123"
    },
    {
      name = "okta-org"
      value = "test-org"
    },
    {
      name = "okta-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_Minimal(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_Minimal,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("okta-adaptive-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_Minimal = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [
    {
      name  = "okta-api-host"
      value = "api.okta.example.com"
    },
    {
      name  = "okta-baseuri"
      value = "/api/v1"
    },
    {
      name  = "okta-token"
      value = "test-token"
    },
    {
      name  = "okta-org"
      value = "test-org"
    },
    {
      name  = "okta-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_MultipleConfigurations(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_MultipleConfigurations,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("duo-security-v2"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-integration-key"),
								"value": knownvalue.StringExact("DIxxxxxxxxxxxxxxxxxx"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-secret-key"),
								"value": knownvalue.StringExact("secretxxxxxxxxxxxxxxxxx"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-api-host"),
								"value": knownvalue.StringExact("api-xxxxxxxx.duosecurity.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-baseuri"),
								"value": knownvalue.StringExact("https://api-xxxxxxxx.duosecurity.com"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_MultipleConfigurations = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "duo-security-v2"
  mfa_config = [
    {
      name = "duo-integration-key"
      value = "DIxxxxxxxxxxxxxxxxxx"
    },
    {
      name = "duo-secret-key"
      value = "secretxxxxxxxxxxxxxxxxx"
    },
    {
      name = "duo-api-host"
      value = "api-xxxxxxxx.duosecurity.com"
    },
    {
      name = "duo-timeout"
      value = "30"
    },
    {
      name = "duo-baseuri"
      value = "https://api-xxxxxxxx.duosecurity.com"
    }
  ]
}
`

func TestAccMfaServerProfile_MaxLengthValues(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate long but valid accesskey value (alphanumeric)
	maxConfigValue := "abcdef" + strings.Repeat("0123456789ABCDEF", 7) + "xyz"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_MaxLengthValues,
				ConfigVariables: map[string]config.Variable{
					"prefix":           config.StringVariable(prefix),
					"max_config_value": config.StringVariable(maxConfigValue),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("rsa-securid-access-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-api-host"),
								"value": knownvalue.StringExact("api.rsa.example.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-baseuri"),
								"value": knownvalue.StringExact("https://tenant.rsa.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-accesskey"),
								"value": knownvalue.StringExact(maxConfigValue),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-accessid"),
								"value": knownvalue.StringExact("RSAID123"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-assurancepolicyid"),
								"value": knownvalue.StringExact("policy-123"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_MaxLengthValues = `
variable "prefix" { type = string }
variable "max_config_value" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "rsa-securid-access-v1"
  mfa_config = [
    {
      name = "rsa-api-host"
      value = "api.rsa.example.com"
    },
    {
      name = "rsa-baseuri"
      value = "https://tenant.rsa.com"
    },
    {
      name = "rsa-accesskey"
      value = var.max_config_value
    },
    {
      name = "rsa-accessid"
      value = "RSAID123"
    },
    {
      name = "rsa-assurancepolicyid"
      value = "policy-123"
    },
    {
      name = "rsa-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_EmptyConfigurations(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_EmptyConfigurations,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("okta-adaptive-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-api-host"),
								"value": knownvalue.StringExact("api.okta.example.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-baseuri"),
								"value": knownvalue.StringExact("/api/v1"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-token"),
								"value": knownvalue.StringExact("test-token"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-org"),
								"value": knownvalue.StringExact("test-org"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_EmptyConfigurations = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [
    {
      name  = "okta-api-host"
      value = "api.okta.example.com"
    },
    {
      name  = "okta-baseuri"
      value = "/api/v1"
    },
    {
      name  = "okta-token"
      value = "test-token"
    },
    {
      name  = "okta-org"
      value = "test-org"
    },
    {
      name  = "okta-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_CertificateProfileOnly(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_CertificateProfileOnly,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_CertificateProfileOnly = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-cert"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_cert_profile = panos_certificate_profile.example.name
}
`

func TestAccMfaServerProfile_AllParameters(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_AllParameters,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.StringExact("comprehensive-cert-profile"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("ping-identity-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("ping-api-host"),
								"value": knownvalue.StringExact("api.pingidentity.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("ping-baseuri"),
								"value": knownvalue.StringExact("https://tenant.pingone.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("ping-token"),
								"value": knownvalue.StringExact("secret789xyz"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("ping-org-alias"),
								"value": knownvalue.StringExact("12345678-1234-1234-1234-123456789abc"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("ping-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_AllParameters = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = "comprehensive-cert-profile"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_cert_profile = panos_certificate_profile.example.name
  mfa_vendor_type = "ping-identity-v1"
  mfa_config = [
    {
      name = "ping-api-host"
      value = "api.pingidentity.com"
    },
    {
      name = "ping-baseuri"
      value = "https://tenant.pingone.com"
    },
    {
      name = "ping-token"
      value = "secret789xyz"
    },
    {
      name = "ping-org-alias"
      value = "12345678-1234-1234-1234-123456789abc"
    },
    {
      name = "ping-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_ConfigurationEdgeCases(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_ConfigurationEdgeCases,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("duo-security-v2"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-api-host"),
								"value": knownvalue.StringExact("api-special.duosecurity.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-integration-key"),
								"value": knownvalue.StringExact("DI!@#$%^&*()_+-="),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-secret-key"),
								"value": knownvalue.StringExact("  secret with spaces  "),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-baseuri"),
								"value": knownvalue.StringExact("https://example.com/path?query=value&other=123"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_ConfigurationEdgeCases = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "duo-security-v2"
  mfa_config = [
    {
      name = "duo-api-host"
      value = "api-special.duosecurity.com"
    },
    {
      name = "duo-integration-key"
      value = "DI!@#$%^&*()_+-="
    },
    {
      name = "duo-secret-key"
      value = "  secret with spaces  "
    },
    {
      name = "duo-baseuri"
      value = "https://example.com/path?query=value&other=123"
    },
    {
      name = "duo-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_AllVendorTypes(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_AllVendorTypes,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.duo",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-duo", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.duo",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("duo-security-v2"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.okta",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-okta", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.okta",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("okta-adaptive-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.ping",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-ping", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.ping",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("ping-identity-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.rsa",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-rsa", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.rsa",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("rsa-securid-access-v1"),
					),
				},
			},
		},
	})
}

const mfaServerProfile_AllVendorTypes = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "duo" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-duo"
  mfa_vendor_type = "duo-security-v2"
  mfa_config = [
    { name = "duo-api-host", value = "api-a1b2c3d4.duosecurity.com" },
    { name = "duo-integration-key", value = "DIWXYZ1234567890ABCD" },
    { name = "duo-secret-key", value = "abcdefghijklmnopqrstuvwxyz0123456789ABCD" },
    { name = "duo-timeout", value = "60" },
    { name = "duo-baseuri", value = "https://api-a1b2c3d4.duosecurity.com" }
  ]
}

resource "panos_mfa_server_profile" "okta" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-okta"
  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [
    { name = "okta-api-host", value = "api.okta.example.com" },
    { name = "okta-baseuri", value = "https://tenant.okta.com" },
    { name = "okta-token", value = "00aaBBccDDeeFFggHHiiJJkkLLmmNNooP" },
    { name = "okta-org", value = "myorganization" },
    { name = "okta-timeout", value = "45" }
  ]
}

resource "panos_mfa_server_profile" "ping" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-ping"
  mfa_vendor_type = "ping-identity-v1"
  mfa_config = [
    { name = "ping-api-host", value = "api.pingidentity.com" },
    { name = "ping-baseuri", value = "https://tenant.pingone.com" },
    { name = "ping-token", value = "secret789xyz" },
    { name = "ping-org-alias", value = "12345678-1234-1234-1234-123456789abc" },
    { name = "ping-timeout", value = "30" }
  ]
}

resource "panos_mfa_server_profile" "rsa" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-rsa"
  mfa_vendor_type = "rsa-securid-access-v1"
  mfa_config = [
    { name = "rsa-api-host", value = "api.rsa.example.com" },
    { name = "rsa-baseuri", value = "https://tenant.rsa.com" },
    { name = "rsa-accesskey", value = "abcdef1234567890ABCDEF1234567890" },
    { name = "rsa-accessid", value = "RSAID123456" },
    { name = "rsa-assurancepolicyid", value = "policy-abc-123-def-456" },
    { name = "rsa-timeout", value = "90" }
  ]
}
`

func TestAccMfaServerProfile_LargeConfigurationList(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Build expected configuration list for state checks (using Okta valid keys)
	expectedConfigs := []knownvalue.Check{
		knownvalue.ObjectExact(map[string]knownvalue.Check{
			"name":  knownvalue.StringExact("okta-api-host"),
			"value": knownvalue.StringExact("api1.okta.com"),
		}),
		knownvalue.ObjectExact(map[string]knownvalue.Check{
			"name":  knownvalue.StringExact("okta-baseuri"),
			"value": knownvalue.StringExact("https://tenant1.okta.com"),
		}),
		knownvalue.ObjectExact(map[string]knownvalue.Check{
			"name":  knownvalue.StringExact("okta-token"),
			"value": knownvalue.StringExact("token123"),
		}),
		knownvalue.ObjectExact(map[string]knownvalue.Check{
			"name":  knownvalue.StringExact("okta-org"),
			"value": knownvalue.StringExact("myorg"),
		}),
		knownvalue.ObjectExact(map[string]knownvalue.Check{
			"name":  knownvalue.StringExact("okta-timeout"),
			"value": knownvalue.StringExact("30"),
		}),
	}

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_LargeConfigurationList,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("okta-adaptive-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact(expectedConfigs),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_LargeConfigurationList = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [
    { name = "okta-api-host", value = "api1.okta.com" },
    { name = "okta-baseuri", value = "https://tenant1.okta.com" },
    { name = "okta-token", value = "token123" },
    { name = "okta-org", value = "myorg" },
    { name = "okta-timeout", value = "30" }
  ]
}
`

func TestAccMfaServerProfile_ConfigurationNameEdgeCases(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	// Generate a long value (100+ characters)
	longValue := strings.Repeat("abcdefghij", 10)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_ConfigurationNameEdgeCases,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"long_value": config.StringVariable(longValue),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("rsa-securid-access-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-api-host"),
								"value": knownvalue.StringExact("api.rsa.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-baseuri"),
								"value": knownvalue.StringExact("https://tenant.rsa.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-accesskey"),
								"value": knownvalue.StringExact(longValue),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-accessid"),
								"value": knownvalue.StringExact("accessid123"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-assurancepolicyid"),
								"value": knownvalue.StringExact("policy456"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-timeout"),
								"value": knownvalue.StringExact("30"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_ConfigurationNameEdgeCases = `
variable "prefix" { type = string }
variable "long_value" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "rsa-securid-access-v1"
  mfa_config = [
    {
      name = "rsa-api-host"
      value = "api.rsa.com"
    },
    {
      name = "rsa-baseuri"
      value = "https://tenant.rsa.com"
    },
    {
      name = "rsa-accesskey"
      value = var.long_value
    },
    {
      name = "rsa-accessid"
      value = "accessid123"
    },
    {
      name = "rsa-assurancepolicyid"
      value = "policy456"
    },
    {
      name = "rsa-timeout"
      value = "30"
    }
  ]
}
`

func TestAccMfaServerProfile_OktaComplete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_OktaComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("okta-adaptive-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-api-host"),
								"value": knownvalue.StringExact("api.okta.example.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-baseuri"),
								"value": knownvalue.StringExact("https://tenant.okta.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-token"),
								"value": knownvalue.StringExact("00aaBBccDDeeFFggHHiiJJkkLLmmNNooP"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-org"),
								"value": knownvalue.StringExact("myorganization"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("okta-timeout"),
								"value": knownvalue.StringExact("45"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_OktaComplete = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-cert"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_cert_profile = panos_certificate_profile.example.name
  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [
    {
      name = "okta-api-host"
      value = "api.okta.example.com"
    },
    {
      name = "okta-baseuri"
      value = "https://tenant.okta.com"
    },
    {
      name = "okta-token"
      value = "00aaBBccDDeeFFggHHiiJJkkLLmmNNooP"
    },
    {
      name = "okta-org"
      value = "myorganization"
    },
    {
      name = "okta-timeout"
      value = "45"
    }
  ]
}
`

func TestAccMfaServerProfile_DuoComplete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_DuoComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("duo-security-v2"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-api-host"),
								"value": knownvalue.StringExact("api-a1b2c3d4.duosecurity.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-integration-key"),
								"value": knownvalue.StringExact("DIWXYZ1234567890ABCD"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-secret-key"),
								"value": knownvalue.StringExact("abcdefghijklmnopqrstuvwxyz0123456789ABCD"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-timeout"),
								"value": knownvalue.StringExact("60"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("duo-baseuri"),
								"value": knownvalue.StringExact("https://api-a1b2c3d4.duosecurity.com"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_DuoComplete = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-cert"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_cert_profile = panos_certificate_profile.example.name
  mfa_vendor_type = "duo-security-v2"
  mfa_config = [
    {
      name = "duo-api-host"
      value = "api-a1b2c3d4.duosecurity.com"
    },
    {
      name = "duo-integration-key"
      value = "DIWXYZ1234567890ABCD"
    },
    {
      name = "duo-secret-key"
      value = "abcdefghijklmnopqrstuvwxyz0123456789ABCD"
    },
    {
      name = "duo-timeout"
      value = "60"
    },
    {
      name = "duo-baseuri"
      value = "https://api-a1b2c3d4.duosecurity.com"
    }
  ]
}
`

func TestAccMfaServerProfile_RSAComplete(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_RSAComplete,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_cert_profile"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_vendor_type"),
						knownvalue.StringExact("rsa-securid-access-v1"),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("mfa_config"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-api-host"),
								"value": knownvalue.StringExact("api.rsa.example.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-baseuri"),
								"value": knownvalue.StringExact("https://tenant.rsa.com"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-accesskey"),
								"value": knownvalue.StringExact("abcdef1234567890ABCDEF1234567890"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-accessid"),
								"value": knownvalue.StringExact("RSAID123456"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-assurancepolicyid"),
								"value": knownvalue.StringExact("policy-abc-123-def-456"),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":  knownvalue.StringExact("rsa-timeout"),
								"value": knownvalue.StringExact("90"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_mfa_server_profile.example",
						tfjsonpath.New("location").AtMapKey("template").AtMapKey("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-tmpl", prefix)),
					),
				},
			},
		},
	})
}

const mfaServerProfile_RSAComplete = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_certificate_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = "${var.prefix}-cert"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_cert_profile = panos_certificate_profile.example.name
  mfa_vendor_type = "rsa-securid-access-v1"
  mfa_config = [
    {
      name = "rsa-api-host"
      value = "api.rsa.example.com"
    },
    {
      name = "rsa-baseuri"
      value = "https://tenant.rsa.com"
    },
    {
      name = "rsa-accesskey"
      value = "abcdef1234567890ABCDEF1234567890"
    },
    {
      name = "rsa-accessid"
      value = "RSAID123456"
    },
    {
      name = "rsa-assurancepolicyid"
      value = "policy-abc-123-def-456"
    },
    {
      name = "rsa-timeout"
      value = "90"
    }
  ]
}
`

// Negative validation tests

func TestAccMfaServerProfile_InvalidVendorType(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_InvalidVendorType,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile(`Invalid MFA Vendor Type`),
			},
		},
	})
}

const mfaServerProfile_InvalidVendorType = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "invalid-vendor-v1"
  mfa_config = [{
    name  = "some-key"
    value = "some-value"
  }]
}
`

func TestAccMfaServerProfile_MissingRequiredKeys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_MissingRequiredKeys,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile(`Missing Required Configuration Keys`),
			},
		},
	})
}

const mfaServerProfile_MissingRequiredKeys = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "okta-adaptive-v1"
  mfa_config = [{
    name  = "okta-api-host"
    value = "api.okta.example.com"
  }]
}
`

func TestAccMfaServerProfile_InvalidConfigKeys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_InvalidConfigKeys,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile(`Invalid Configuration Keys`),
			},
		},
	})
}

const mfaServerProfile_InvalidConfigKeys = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "duo-security-v2"
  mfa_config = [
    # All required Duo keys
    { name = "duo-api-host", value = "api-xxxxxxxx.duosecurity.com" },
    { name = "duo-integration-key", value = "DIxxxxxxxxxxxxxxxxxx" },
    { name = "duo-secret-key", value = "secretxxxxxxxxxxxxxxxxx" },
    { name = "duo-timeout", value = "30" },
    { name = "duo-baseuri", value = "https://api-xxxxxxxx.duosecurity.com" },
    # Invalid Okta keys
    { name = "okta-api-host", value = "api.okta.com" },
    { name = "okta-baseuri", value = "https://tenant.okta.com" },
    { name = "okta-org", value = "myorg" }
  ]
}
`

func TestAccMfaServerProfile_MissingAndInvalidKeys(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_MissingAndInvalidKeys,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				// Both errors should be reported
				ExpectError: regexp.MustCompile(`(Missing Required Configuration Keys|Invalid Configuration Keys)`),
			},
		},
	})
}

const mfaServerProfile_MissingAndInvalidKeys = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  mfa_vendor_type = "ping-identity-v1"
  mfa_config = [
    # Only 1 valid Ping key
    { name = "ping-api-host", value = "api.pingidentity.com" },
    # 2 invalid keys from other vendors
    { name = "duo-integration-key", value = "DIxxxxxxxxxxxxxxxxxx" },
    { name = "okta-org", value = "myorg" }
  ]
}
`

func TestAccMfaServerProfile_ConfigWithoutVendorType(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: mfaServerProfile_ConfigWithoutVendorType,
				ConfigVariables: map[string]config.Variable{
					"prefix": config.StringVariable(prefix),
				},
				ExpectError: regexp.MustCompile(`Configuration Requires Vendor Type`),
			},
		},
	})
}

const mfaServerProfile_ConfigWithoutVendorType = `
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = "${var.prefix}-tmpl"
}

resource "panos_mfa_server_profile" "example" {
  location = { template = { name = panos_template.example.name } }
  name = var.prefix

  # Config provided without vendor type should fail
  mfa_config = [{
    name  = "some-key"
    value = "some-value"
  }]
}
`
