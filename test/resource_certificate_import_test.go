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
			{
				Config: certificateImport_Local_PEM_CertificateWithKeyAndPassphrase_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"certificate1": config.StringVariable(certPemInitial),
					"private_key1": config.StringVariable(privateKeyPemEncryptedInitial),
					"passphrase1":  config.StringVariable("paloalto"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("name"),
						knownvalue.StringExact(fmt.Sprintf("%s-cert1", prefix)),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("certificate"),
						knownvalue.StringExact(certPemInitial),
					),
					statecheck.ExpectKnownValue(
						"panos_certificate_import.example1",
						tfjsonpath.New("local").AtMapKey("pem").AtMapKey("private_key"),
						knownvalue.StringExact(privateKeyPemEncryptedInitial),
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

const certificateImport_Local_PEM_CertificateWithKeyAndPassphrase_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "certificate1" { type = string }
variable "private_key1" { type = string }
variable "passphrase1" { type = string }

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
      passphrase = var.passphrase1
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

const certPemInitial = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`

const certPemUpdated = `-----BEGIN CERTIFICATE-----
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
-----END CERTIFICATE-----`

const privateKeyPemInitial = `-----BEGIN PRIVATE KEY-----
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
7xFraUnOHD26eQ5+DBups9jL+Y+Hgw==
-----END PRIVATE KEY-----`

const privateKeyPemEncryptedInitial = `-----BEGIN ENCRYPTED PRIVATE KEY-----
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
-----END ENCRYPTED PRIVATE KEY-----`

const privateKeyPemUpdated = `-----BEGIN PRIVATE KEY-----
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
oTGtUd5Icd7xkTm7dOJbUXfXuOy8Bps=
-----END PRIVATE KEY-----`

const certKeyPkcs12Initial = `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`

const certKeyPkcs12Updated = `MIIRTwIBAzCCEQUGCSqGSIb3DQEHAaCCEPYEghDyMIIQ7jCCBtoGCSqGSIb3DQEHBqCCBsswggbHAgEAMIIGwAYJKoZIhvcNAQcBMF8GCSqGSIb3DQEFDTBSMDEGCSqGSIb3DQEFDDAkBBAx/fJgTgYCpFMBcFiD93MFAgIIADAMBggqhkiG9w0CCQUAMB0GCWCGSAFlAwQBKgQQ8M4V7QoVKHs+zjAkshQph4CCBlCL+qd3eoPN/UVPJGcbNfhBmvgm49NcpFics34M0j3ltCvZ+mXQxu9f4CEAhVvcCen2UTZ1mulMYVQ7VYZ0NYv+u7oxakJqf33JwgGCcd+PzjGLFGtj/9ICf2cbgluS5tUn5gxBOxA3VQCca4vGlwSjhMp5m1MH96aMSEPcEWzGBrakUruXxjpknV/K58B1mWh/oNppsalbRyalYo4kYExM6mw6jN/fxnSKR5vDMbTRgdKV/cKtnsaoQLnQRTqMFAGGlUi2VT4dnRy4lPu7TdC/Vetzw2tF1mBpvGZwmuPxGI4IFqZa6WVS7eeIG6pm9HcgbARJb0B8bztavGwJxFFAQ7n32celmPfusHtj3ldzlQMBl6/aAWbR+y53TR4wSHLtg80vOc+axrW2pIb5/JwMTr++eHOMsugxVhaeJ086imVGJpFYd9iAffzrk9MEWCFXcKyHCf6ZElMZO1E4v+sAJOczF1Tifm/KAkss/U9u6gpPwo7Eh/bo33RpJO76jJFl/6JQ9YHPzZUUczEMKvydK6yZsjZJZhSGbUxiCe50Q1Mu8zq6shHcWoStfdFYCZIuBFIKW9bi0HOhNy0ZWLWlQDoktLtaJuAgYO3FyMQ4lBU8RE9y+EkQKQZhb/+Y3Vb/duu0/nurKWtuB3Qjxb5+2eoLJKaX3P7ooPniCuYIzWkcZmmCXnlgfeSb8KlrQwOtT0bMtL7jGLROy90YKMedaN9CslGs1A/BgHOI8ymKghG4qyV+6njejQsXYtDyf6Q8CS/QRhSWb4wErvXDCsifKpEbsH0Zf025ELQmx17vizAVdZblOD/tQrLAGCVLNXV1MzsIa4t7cmetnI0RTl5mw8ewsQW5HU4HTmUpMFWa3ayezNIflIuSD8//ct6/gOsZHV25S2KsBHpb21K8XNzG6cvUXKZdSWLq1pLLAzel64CVA+TwE4fBEmrAw0YEznJ9FkJ/IQmu7fFnxIr8g2bKhUbG+qUIqC5D+Pp9uJ8IXD23d/Q1pm163Zl/o8HXTkRuy3A+1oS8JQwexj8MhZsVH4Oc6exT8uKZS2SXETlovp0apjZ4cJg9B+d7PR7WJf46tbranbqTX5RdxBtZBWdf0QuDz8gGNv2xoo0rn//pTlFn696JMn7/H0KDLf716Lta6YwdovvaPCnvqxygMzdbPr5TuXzNRbSvGECdcH/3OhxAZN9tgpor5TSl0ly8mOmU2pO9N3Gkxib3K6ohBVkSHr5kPi3TJV53Y5RKqq1yUdsVE8CNhb5F6ShJAjeYCpVs+u5qSOp0a55iz286cFO2I/Dy9962m+pHArhbT5hc5bmezKbG6w6kx72xBk3CnOXylXaoJLQe43H2rNcNiLlUlXxxNugPO/NnTOxhmdGHFgPbLOET2qQm+/xGY06Mmqfcb5tIkFKNAUQb6R4eHMZk5RJexFJJM+249vQmjhGTb5tZo1DgDYJR36eTbLNUi5gpRYAbs+3FNB2lby9iHRhhqEo3fBS+ONHiN61fTLu3vRZQlJCzyupAQNkUDRZRVDFH2YdXxHhG5dCFRntc/rtd6CfJwZEIv+YPUv1cCSHXSD0U9NL69g84NGHO2uZQfrU9ShyyU/uoXYLr8Z0oHx3EKADS1mQMhW1hQsSu9qC8Wj9vxMXHj1tbqRv65HDhNiLBjOt4Gty6tEnKQjROjQGqO5M3f07fwSUm2rGw4c3Va3/QkYr0FjvfDsyrhRutGQABmA6wtF1jQqy99cT6ReUwh3aujFQ9NPegRmddOej8Ye0xEDaOC+/FHH3rwIg27ULjNnvqsV6GyoiTOgoXoESY4qhGJFpdeAyHVpo6mhwf9d6PDqwCChzWweu2ZvFpTjtjv71XkcQuvP6bVN16sUbZXnDvgQgXJ/vqOFKclnVIACmqJpj0ZF9JIp2cQ2yZcO8VkKvOBMtJ8iiS+jUpVOU69DzTOlEFUUsXnsIFR/tZoPWFACfmLBd9KI64Z8ntUJESkzs/E58S8sGHSHw1vIqV1i7VreWJhINKMY2VxGx7RsVjRfAg2SuJ5Efx3diT+Ugrxkq2BT2iu1L9xg5KiKRg+FWGAXhj5MovRrAXu7iQJS6yjAfKelDRp8WR5viW0xA+wOCAX309h+i9F0bVCB4asEXuFEssNBUKCulSxxo0pTCCCgwGCSqGSIb3DQEHAaCCCf0Eggn5MIIJ9TCCCfEGCyqGSIb3DQEMCgECoIIJuTCCCbUwXwYJKoZIhvcNAQUNMFIwMQYJKoZIhvcNAQUMMCQEEIgAKG3W7syyBDjoxEoIDkkCAggAMAwGCCqGSIb3DQIJBQAwHQYJYIZIAWUDBAEqBBC+y690TDfGyTL2yjDgoJrmBIIJUHT77ca1wdx5jndYBRZq9FtX9EDNrWu1855QX77Zu8lk2ayaXIDU9tEsg3UGWxXhv3aALfO5WyppWfuRzNQG464nIy8l76Nwa9oOnrpsbZn/V2r4xNfobsCrrxWxo+7sTSYaWx2H02Fkeu11jl3s+ifH2r34JZAmcPF6x+6PLxAhkOUZc+vOqw1KeCOFgjFnwdHtRWfwnMzTgENUr2/fRerLukum8ooDmqeb3qJzinjjavAUtUJNPO/EVE08qI5nd6ojEjY+BkWywinuVzwLlRqo7UxV3xcDOUkGyzhxzOMwooEVjUcWun+FOIIAMRCPI1XgSUdHDVN1xpgTBF0f7Ai8g8vIAC64dvyPr+yz7CMe+4VEJlxth47eW3M5qucGJ2renHq/O05f1TWhR33lC0F9wtelLjRdhkdKff+a+kygzagCQPiw3pM6lgnJwKyLNPfFR+RQIJ/tgdin+yr0B2mfeKK3qBNLvjTTrHUc1ubBp5oyCNnCcKWz9sKBV7x5TGvN3NKs/rmKqyxj8MjdEni7RWPM0TTHtPZct9c3Z8v59sWRrMha10BFxr6D4xKcGAQLlQKT7yv+kdl4c4paoTpaE6EgcRW4WzcyW6uxxAagCC4ATep+xUQ+xxLscElgiUNZa+PsLR2HLrbDAXJqkNHUqhd7kNyUUiUCckCZQZST0OnyTrsvNIfIJ3pGFSFeFmEftonz6MOb3yv/mmDcajgIrKuaM78FGbld7oYPnpOzQqrnGE4GMsordvaj38P9ztvxEJF7y1+0y2Ac6s0AoOoWasU4Q5SOWKSvgyYCaOmkuoeKTXIbRCiFwcGVW6OV5CZFEqFu5bKAAh97r4SfPNanaKDglyelxhyoY2QgRcYH3xA8ROQEOl4V/MEjg/H4G8cz9f+LAbslpCBmL3ECEt4X47j+KQpzn4L7KhdCAjclXGdpilmRi8WUoc1GcpnZZjSSE7IvzDbkIC/040z9PPG+cCsiuZ9aJSlhsbIP+QiC8LEgNkRAaHbNyzfWnCAk90kw9vLHuSg55m3zCMrFvcM3Rz6y0AoP0YjGC+obW78kWMZuKcKN3nf+w7KUJx9MnkuF6v+pJojRG6JG6CdQlUi9xZBvE8c7QIyKt6EPgh/KBVvpMztpax7dvrgDjezCJZMp0KYnG6s8RyU/Y0Q+gstVvaHnn1BMErjxjvbaFGkceX18qGpbjd9jHWWO+TiV5pBDTJbmRldKGOObnvjNtovBr9MknOORmDWaqC3ecZ4Z+VoD7mfCZoBks/mw4vPghBUDJNGCeEX5fpGPhet/dHF0M8ENYfHbYJVDjcbz2WgQjKee7B5kUsnMou2CMsGF42AIlLa8R83nJlKn6j/r+nsj/dDnGMerz1yV29zr92cSaUz6CuMBiS8iLShUa783+Lq/uh4TzUjUH/6jGK0y5MsSwLXYs5uNpNwGXFA9z+e3LOUDFJtYkIr+jem10Bu0yEfXR/Ey8z83A9WmFXcM2+K9CmYySQsLekiNq3zcrylEC0ECqNMpU1Qo2Y+t1jxTWnk5F0Utd3z+310b+r3bYJtPDNB/ir6n/OONY7H1xlCbmeR5+e/thZlXYLJvB1SW0eQuPWip4XNA7t3Cc+BB0ybNCeq/NBsmw85/IE8b8SGRfo0gxFfIud2hL7sUJ3hOlBiNa0ocL9Jxf+jmLqhelmu+5UpXZrG6h6PbS9R/wPTPqUr4hDQzC+a3rRUVbG+MDqhHPH5EtQjma82sMSri/OYhD8ErVneliIk5Keqa4ycbqKF5nLZS5fkysSrwnV1tvtYfWZkkL6m1M0ErSmfJPAeVJQ+Io3GWC77IW3uPyGGwezGmFmjkRRA3HiKuzD7Gk2KnETfLP5TwyGrJOmH6HvFdhiKWj5SElR9KGXkZKibypZToHV9LH48whJjA+kHraBDLjjnvGj7lxUt+vfLRs5Chm/s9Jtne0sVxaKfMzQ/2b7vUJ6gOFjeekj+Le0YeP2JJTVtlEL1pWkoKmxgkOFBIhjY3Orc81+zs6oWtX5bXKxG7PM5MUoJ2oAk9bLdcKATWf/zPPib3VbIg3hsm1d9dYv4vvwE525C0unn3Hsbg8Hwvy4O/48IMeY76+Zr+koNn1QTH2pDIkaYAZPsihVHwWUspbSymAJgJlcyA42pzorfMz7Pn0EFhzH5ukMd5+qUR/gryNA4TRBTL+ht3KdPsq/jhxHEP6jZSDMO/4qrMxeKoLgknb0VVpTPxxOLvKtTlOANJ0QcAPaKULx9mnqOLzEirNDMjfUhNYBT7QqxSDxRnMf6jgFPEkR9a4bSGhrkyx+SGI303B5MYG9jpJjXPBlAGpaiHLjLg2WEuv1jGE1F5fwWDgm/Hbavn41iU0bk02pmdeQWx5j4rYijnrV0278qHu8LyGRS90d5m2IdPAEnGXQyUZiDfd4Hj+tANkYr/YJh8ccLkYUXzUQ6lwao5XHBB6+BOYMeh07IgJ4CljcFKQ0V2tIwcgSDCyI3ijXi3xc8Ww37kmz1vlrjHhdq+dXJIynsOixULUaWXXuHJheMezPqQ52qOe8uQlSeuiXdHAzWzc97z5R2u+4i62LaUoAmPq0oM7/wp266xPOM6KkKpqrOzF43iZuVrlxRcjihdt96s5tYvArGXIfCNE6Qtk4M4vpjZj6VWppUsBRR6yYosUX/gLfQxga3Ltd2xhc6qSbW6Oc4WrcqUvlcpYp7NI56HTHNG3n8GvkahWGL9CTgLYMfGsgN8TtPsMA1xSr4mJypmaK+lSCfj0opJWA76jWlUVTjmFTxIGfNUEkAGp6gFi7B4KBygS5ga1N+9VK1bShgNmLtWbwiBdSqktwdAzfmqFkmJLQMa+7btFV/TukuT8VJyb0wZef/Gu7X+hYWr4WJPlmCrzkYfYqsgfZMY80KUet/uv0FTg9Xaund0rQdoJpxQDJ0WmL1Mr49PWiycvQ6T/9qZFrYisXdGQDjRByIDH1R+N/HPGXn2Mohii02T1y1KjpEYqHgBFqypQ6TwEyrqNcmMPD++ocxhsW2Ldg3SABvI/SeKOJr7R3Jv+5fOnAEbFeeZkSmj7eCqBPNTznQFWOZ0cWIc+k9SUwwT9CZ0SAvwf+RbhfcyBhfJPZSm1Fs0vF3CdGQFvwG3bh0maHzqkx4Zd4RGsYsBod/m1XSqQQYLRdGuSBuU7ChXMSUwIwYJKoZIhvcNAQkVMRYEFKuze4PGkjz8qsvYivOoWp1DY08QMEEwMTANBglghkgBZQMEAgEFAAQgs+8VO3qFJQm3118TzIzLLIe8OybbcdcU/FW8Rz8X0aIECGqGgqIqU/n3AgIIAA==
`
