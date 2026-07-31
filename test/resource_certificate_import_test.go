package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccCertificateImport_Local_PEM_Certificate(t *testing.T) {
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
				Config: certificateImport_Local_PEM_Certificate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemInitial),
					),
				},
			},
			{
				Config: certificateImport_Local_PEM_Certificate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemUpdated),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			{
				Config: certificateImport_Local_PEM_Certificate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemUpdated),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemUpdated),
					),
				},
			},
		},
	})
}

func TestAccCertificateImport_Local_PEM_CertificateWithKey(t *testing.T) {
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
				Config: certificateImport_Local_PEM_CertificateWithKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"certificate1": config.StringVariable(certPemInitial),
					"private_key1": config.StringVariable(privateKeyPemInitial),
					"certificate2": config.StringVariable(certPemUpdated),
					"private_key2": config.StringVariable(privateKeyPemUpdated),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example2",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert2", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemInitial),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("private_key"),
						knownvalue.StringExact(privateKeyPemInitial),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example2",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemUpdated),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example2",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("private_key"),
						knownvalue.StringExact(privateKeyPemUpdated),
					),
				},
			},
			{
				Config: certificateImport_Local_PEM_Certificate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certPemUpdated),
					"private_key": config.StringVariable(privateKeyPemUpdated),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			{
				Config: certificateImport_Local_PEM_CertificateWithKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"certificate1": config.StringVariable(certPemUpdated),
					"private_key1": config.StringVariable(privateKeyPemUpdated),
					"certificate2": config.StringVariable(certPemInitial),
					"private_key2": config.StringVariable(privateKeyPemInitial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example2",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert2", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemUpdated),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("private_key"),
						knownvalue.StringExact(privateKeyPemUpdated),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example2",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemInitial),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example2",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("private_key"),
						knownvalue.StringExact(privateKeyPemInitial),
					),
				},
			},
		},
	})
}

func TestAccCertificateImport_Local_PKCS12_Certificate(t *testing.T) {
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
				Config: certificateImport_Local_PKCS12_Certificate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certKeyPkcs12Initial),
					"passphrase":  config.StringVariable("paloalto"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pkcs12").AtMapKey("certificate"),
						knownvalue.StringExact(certKeyPkcs12Initial),
					),
				},
			},
			{
				Config: certificateImport_Local_PKCS12_Certificate_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certKeyPkcs12Updated),
					"passphrase":  config.StringVariable("paloalto"),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
		},
	})
}

// NOTE: This test is commented out because of a bug in the provider's Read function for vsys locations.
// func TestAccCertificateImport_Vsys_Local_PEM_Certificate(t *testing.T) {
// 	t.Parallel()
//
// 	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
// 	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)
//
// 	location := config.ObjectVariable(map[string]config.Variable{
// 		"vsys": config.ObjectVariable(map[string]config.Variable{
// 			"name": config.StringVariable("vsys1"),
// 		}),
// 	})
//
// 	resource.Test(t, resource.TestCase{
// 		PreCheck:                 func() { testAccPreCheck(t) },
// 		ProtoV6ProviderFactories: testAccProviders,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: certificateImport_Vsys_Local_PEM_Certificate_Tmpl,
// 				ConfigVariables: map[string]config.Variable{
// 					"prefix":      config.StringVariable(prefix),
// 					"location":    location,
// 					"certificate": config.StringVariable(certPemInitial),
// 				},
// 				ConfigStateChecks: []statecheck.StateCheck{
// 					statecheck.ExpectKnownValue(
// 						"panos_certificate_import.example1",
// 						tfjsonpath.New("name"),
// 						knownvalue.StringExact(prefix),
// 					),
// 					statecheck.ExpectKnownValue(
// 						"panos_certificate_import.example1",
// 						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
// 						knownvalue.StringExact(certPemInitial),
// 					),
// 				},
// 			},
// 		},
// 	})
// }

func TestAccCertificateImport_Local_PKCS12_CertificateWithKey(t *testing.T) {
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
				Config: certificateImport_Local_PKCS12_CertificateWithKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certKeyPkcs12Initial),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").
							AtMapKey("pkcs12").
							AtMapKey("certificate"),
						knownvalue.StringExact(certKeyPkcs12Initial),
					),
				},
			},
			{
				Config: certificateImport_Local_PKCS12_CertificateWithKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certKeyPkcs12Updated),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectNonEmptyPlan(),
					},
				},
			},
			{
				Config: certificateImport_Local_PKCS12_CertificateWithKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certKeyPkcs12Updated),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").
							AtMapKey("pkcs12").
							AtMapKey("certificate"),
						knownvalue.StringExact(certKeyPkcs12Updated),
					),
				},
			},
			{
				Config: certificateImport_Local_PKCS12_CertificateWithKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":      config.StringVariable(prefix),
					"location":    location,
					"certificate": config.StringVariable(certKeyPkcs12Updated),
				},
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

const certificateImport_Local_PEM_Certificate_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_certificate_import" "example1" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  local = {
    pem = {
      certificate = var.certificate
    }
  }
}
`

const certificateImport_Local_PEM_CertificateWithKey_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate1" { type = string }
variable "private_key1" { type = string }

variable "certificate2" { type = string }
variable "private_key2" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_certificate_import" "example1" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert1"

  local = {
    pem = {
      certificate = var.certificate1
      private_key = var.private_key1
    }
  }
}

resource "panos_certificate_import" "example2" {
  depends_on = [panos_template.example]
  location = var.location

  name = "${var.prefix}-cert2"

  local = {
    pem = {
      certificate = var.certificate2
      private_key = var.private_key2
    }
  }
}
`

const certificateImport_Local_PKCS12_Certificate_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }
variable "passphrase" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_certificate_import" "example1" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  local = {
    pkcs12 = {
      certificate = var.certificate
      passphrase = var.passphrase
    }
  }
}
`

const certificateImport_Local_PKCS12_CertificateWithKey_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name = var.prefix
}

resource "panos_certificate_import" "example1" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  local = {
    pkcs12 = {
      certificate = var.certificate
      passphrase = "paloalto"
    }
  }
}
`
