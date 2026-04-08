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

func TestAccPanosSamlIdpProfile_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosSamlIdpProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"cert_name":   config.StringVariable(certName),
					"certificate": config.StringVariable(certPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("certificate"),
						knownvalue.StringExact(certName),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("entity_id"),
						knownvalue.StringExact("my-entity-id"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("sso_url"),
						knownvalue.StringExact("https://my-sso-url.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("slo_url"),
						knownvalue.StringExact("https://my-slo-url.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("max_clock_skew"),
						knownvalue.Int64Exact(120),
					),
				},
			},
		},
	})
}

const panosSamlIdpProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "cert_name" { type = string }
variable "certificate" { type = string }

resource "panos_template" "example" {
    location = { panorama = {} }
    name     = format("%s-tmpl", var.prefix)
}

resource "panos_certificate_import" "cert" {
	location = {
		template = {
			name = panos_template.example.name
		}
	}
    name = var.cert_name
    local = {
        pem = {
            certificate = var.certificate
        }
    }
}

resource "panos_saml_idp_profile" "example" {
    location = {
        template = {
            name = panos_template.example.name
        }
    }
    name     = var.prefix

	certificate    = panos_certificate_import.cert.name
    entity_id      = "my-entity-id"
    sso_url        = "https://my-sso-url.com"
    slo_url        = "https://my-slo-url.com"
    max_clock_skew = 120
}
`

func TestAccPanosSamlIdpProfile_AllBoolFlags(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosSamlIdpProfile_AllBoolFlags_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"cert_name":   config.StringVariable(certName),
					"certificate": config.StringVariable(certPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("admin_use_only"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("validate_idp_certificate"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("want_auth_requests_signed"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const panosSamlIdpProfile_AllBoolFlags_Tmpl = `
variable "prefix" { type = string }
variable "cert_name" { type = string }
variable "certificate" { type = string }

resource "panos_template" "example" {
    location = { panorama = {} }
    name     = format("%s-tmpl", var.prefix)
}

resource "panos_certificate_import" "cert" {
	location = {
		template = {
			name = panos_template.example.name
		}
	}
    name = var.cert_name
    local = {
        pem = {
            certificate = var.certificate
        }
    }
}

resource "panos_saml_idp_profile" "example" {
    location = {
        template = {
            name = panos_template.example.name
        }
    }
    name     = var.prefix

	certificate    = panos_certificate_import.cert.name
    entity_id      = "my-entity-id"
    sso_url        = "https://my-sso-url.com"
	admin_use_only = true
	validate_idp_certificate = true
	want_auth_requests_signed = true
}
`

func TestAccPanosSamlIdpProfile_AttributeNameImports(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosSamlIdpProfile_AttributeNameImports_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"cert_name":   config.StringVariable(certName),
					"certificate": config.StringVariable(certPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("attribute_name_access_domain_import"),
						knownvalue.StringExact("access_domain"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("attribute_name_admin_role_import"),
						knownvalue.StringExact("admin_role"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("attribute_name_usergroup_import"),
						knownvalue.StringExact("usergroup"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("attribute_name_username_import"),
						knownvalue.StringExact("username"),
					),
				},
			},
		},
	})
}

const panosSamlIdpProfile_AttributeNameImports_Tmpl = `
variable "prefix" { type = string }
variable "cert_name" { type = string }
variable "certificate" { type = string }

resource "panos_template" "example" {
    location = { panorama = {} }
    name     = format("%s-tmpl", var.prefix)
}

resource "panos_certificate_import" "cert" {
	location = {
		template = {
			name = panos_template.example.name
		}
	}
    name = var.cert_name
    local = {
        pem = {
            certificate = var.certificate
        }
    }
}

resource "panos_saml_idp_profile" "example" {
    location = {
        template = {
            name = panos_template.example.name
        }
    }
    name     = var.prefix

	certificate    = panos_certificate_import.cert.name
    entity_id      = "my-entity-id"
    sso_url        = "https://my-sso-url.com"
	attribute_name_access_domain_import = "access_domain"
	attribute_name_admin_role_import = "admin_role"
	attribute_name_usergroup_import = "usergroup"
	attribute_name_username_import = "username"
}
`

func TestAccPanosSamlIdpProfile_RedirectBindings(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosSamlIdpProfile_RedirectBindings_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"cert_name":   config.StringVariable(certName),
					"certificate": config.StringVariable(certPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("slo_bindings"),
						knownvalue.StringExact("redirect"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("sso_bindings"),
						knownvalue.StringExact("redirect"),
					),
				},
			},
		},
	})
}

const panosSamlIdpProfile_RedirectBindings_Tmpl = `
variable "prefix" { type = string }
variable "cert_name" { type = string }
variable "certificate" { type = string }

resource "panos_template" "example" {
    location = { panorama = {} }
    name     = format("%s-tmpl", var.prefix)
}

resource "panos_certificate_import" "cert" {
	location = {
		template = {
			name = panos_template.example.name
		}
	}
    name = var.cert_name
    local = {
        pem = {
            certificate = var.certificate
        }
    }
}

resource "panos_saml_idp_profile" "example" {
    location = {
        template = {
            name = panos_template.example.name
        }
    }
    name     = var.prefix

	certificate    = panos_certificate_import.cert.name
    entity_id      = "my-entity-id"
    sso_url        = "https://my-sso-url.com"
	slo_bindings   = "redirect"
	sso_bindings   = "redirect"
}
`

func TestAccPanosSamlIdpProfile_Update(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
	certName := fmt.Sprintf("%s-cert", prefix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosSamlIdpProfile_Update_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"cert_name":   config.StringVariable(certName),
					"certificate": config.StringVariable(certPemInitial),
					"sso_url": config.StringVariable("https://my-sso-url.com"),
					"max_clock_skew": config.IntegerVariable(120),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("sso_url"),
						knownvalue.StringExact("https://my-sso-url.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("max_clock_skew"),
						knownvalue.Int64Exact(120),
					),
				},
			},
			{
				Config: panosSamlIdpProfile_Update_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"cert_name":   config.StringVariable(certName),
					"certificate": config.StringVariable(certPemInitial),
					"sso_url": config.StringVariable("https://new-sso-url.com"),
					"max_clock_skew": config.IntegerVariable(60),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("sso_url"),
						knownvalue.StringExact("https://new-sso-url.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_saml_idp_profile.example",
						tfjsonpath.New("max_clock_skew"),
						knownvalue.Int64Exact(60),
					),
				},
			},
		},
	})
}

const panosSamlIdpProfile_Update_Tmpl = `
variable "prefix" { type = string }
variable "cert_name" { type = string }
variable "certificate" { type = string }
variable "sso_url" { type = string }
variable "max_clock_skew" { type = number }

resource "panos_template" "example" {
    location = { panorama = {} }
    name     = format("%s-tmpl", var.prefix)
}

resource "panos_certificate_import" "cert" {
	location = {
		template = {
			name = panos_template.example.name
		}
	}
    name = var.cert_name
    local = {
        pem = {
            certificate = var.certificate
        }
    }
}

resource "panos_saml_idp_profile" "example" {
    location = {
        template = {
            name = panos_template.example.name
        }
    }
    name     = var.prefix

	certificate    = panos_certificate_import.cert.name
    entity_id      = "my-entity-id"
    sso_url        = var.sso_url
	max_clock_skew = var.max_clock_skew
}
`