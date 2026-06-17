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

func TestAccSnmptrapProfile_Basic_V2c(t *testing.T) {
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
				Config: snmptrapProfile_Basic_V2c_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v2c").AtMapKey("servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("trap-v2c-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v2c").AtMapKey("servers").AtSliceIndex(0).AtMapKey("manager"),
						knownvalue.StringExact("192.168.1.100"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v2c").AtMapKey("servers").AtSliceIndex(0).AtMapKey("community"),
						knownvalue.StringExact("public"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v2c").AtMapKey("servers").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("trap-v2c-2"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v2c").AtMapKey("servers").AtSliceIndex(1).AtMapKey("manager"),
						knownvalue.StringExact("192.168.1.101"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v2c").AtMapKey("servers").AtSliceIndex(1).AtMapKey("community"),
						knownvalue.StringExact("private"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_Basic_V2c_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v2c = {
      servers = [
        {
          name      = "trap-v2c-1"
          manager   = "192.168.1.100"
          community = "public"
        },
        {
          name      = "trap-v2c-2"
          manager   = "192.168.1.101"
          community = "private"
        }
      ]
    }
  }
}
`

func TestAccSnmptrapProfile_Basic_V3(t *testing.T) {
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
				Config: snmptrapProfile_Basic_V3_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					// First server: all params with non-default auth/priv protocols
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("trap-v3-1"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("manager"),
						knownvalue.StringExact("10.0.0.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("user"),
						knownvalue.StringExact("snmpuser"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("engine_id"),
						knownvalue.StringExact("0x80001F8880A1B2C3D4E5F601"),
					),
					// authentication_password has hashing type solo - only check NotNull
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("authentication_password"),
						knownvalue.NotNull(),
					),
					// privacy_password has hashing type solo - only check NotNull
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("privacy_password"),
						knownvalue.NotNull(),
					),
					// Non-default auth protocol (SHA-256 instead of SHA)
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("authentication_protocol"),
						knownvalue.StringExact("SHA-256"),
					),
					// Non-default priv protocol (AES-192 instead of AES)
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("privacy_protocol"),
						knownvalue.StringExact("AES-192"),
					),
					// Second server: uses defaults for auth/priv protocols
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("name"),
						knownvalue.StringExact("trap-v3-2"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("manager"),
						knownvalue.StringExact("10.0.0.2"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("user"),
						knownvalue.StringExact("snmpuser2"),
					),
					// authentication_password has hashing type solo - only check NotNull
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("authentication_password"),
						knownvalue.NotNull(),
					),
					// privacy_password has hashing type solo - only check NotNull
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("privacy_password"),
						knownvalue.NotNull(),
					),
					// Default auth protocol (SHA)
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("authentication_protocol"),
						knownvalue.StringExact("SHA"),
					),
					// Default priv protocol (AES)
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(1).AtMapKey("privacy_protocol"),
						knownvalue.StringExact("AES"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_Basic_V3_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v3 = {
      servers = [
        {
          name                    = "trap-v3-1"
          manager                 = "10.0.0.1"
          user                    = "snmpuser"
          engine_id               = "0x80001F8880A1B2C3D4E5F601"
          authentication_password = "AuthPassword1!"
          privacy_password        = "PrivPassword1!"
          authentication_protocol = "SHA-256"
          privacy_protocol        = "AES-192"
        },
        {
          name                    = "trap-v3-2"
          manager                 = "10.0.0.2"
          user                    = "snmpuser2"
          authentication_password = "AuthPassword2!"
          privacy_password        = "PrivPassword2!"
        }
      ]
    }
  }
}
`

// --- Enum value coverage tests ---
// Each test exercises untested auth/priv protocol enum values.

func TestAccSnmptrapProfile_V3_AuthProto_SHA224(t *testing.T) {
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
				Config: snmptrapProfile_V3_AuthProto_SHA224_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("authentication_protocol"),
						knownvalue.StringExact("SHA-224"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("privacy_protocol"),
						knownvalue.StringExact("AES-256"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_V3_AuthProto_SHA224_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v3 = {
      servers = [
        {
          name                    = "sha224-server"
          manager                 = "10.0.0.1"
          user                    = "snmpuser"
          authentication_password = "AuthPassword1!"
          privacy_password        = "PrivPassword1!"
          authentication_protocol = "SHA-224"
          privacy_protocol        = "AES-256"
        }
      ]
    }
  }
}
`

func TestAccSnmptrapProfile_V3_AuthProto_SHA384(t *testing.T) {
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
				Config: snmptrapProfile_V3_AuthProto_SHA384_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("authentication_protocol"),
						knownvalue.StringExact("SHA-384"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_V3_AuthProto_SHA384_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v3 = {
      servers = [
        {
          name                    = "sha384-server"
          manager                 = "10.0.0.1"
          user                    = "snmpuser"
          authentication_password = "AuthPassword1!"
          privacy_password        = "PrivPassword1!"
          authentication_protocol = "SHA-384"
        }
      ]
    }
  }
}
`

func TestAccSnmptrapProfile_V3_AuthProto_SHA512(t *testing.T) {
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
				Config: snmptrapProfile_V3_AuthProto_SHA512_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("authentication_protocol"),
						knownvalue.StringExact("SHA-512"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_V3_AuthProto_SHA512_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v3 = {
      servers = [
        {
          name                    = "sha512-server"
          manager                 = "10.0.0.1"
          user                    = "snmpuser"
          authentication_password = "AuthPassword1!"
          privacy_password        = "PrivPassword1!"
          authentication_protocol = "SHA-512"
        }
      ]
    }
  }
}
`

// Legacy enum values from PR — not in our XML schema.
// These tests verify whether the device actually accepts MD5, DES, AES128.

func TestAccSnmptrapProfile_V3_Legacy_MD5(t *testing.T) {
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
				Config: snmptrapProfile_V3_Legacy_MD5_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("authentication_protocol"),
						knownvalue.StringExact("MD5"),
					),
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("privacy_protocol"),
						knownvalue.StringExact("DES"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_V3_Legacy_MD5_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v3 = {
      servers = [
        {
          name                    = "md5-server"
          manager                 = "10.0.0.1"
          user                    = "snmpuser"
          authentication_password = "AuthPassword1!"
          privacy_password        = "PrivPassword1!"
          authentication_protocol = "MD5"
          privacy_protocol        = "DES"
        }
      ]
    }
  }
}
`

func TestAccSnmptrapProfile_V3_Legacy_AES128(t *testing.T) {
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
				Config: snmptrapProfile_V3_Legacy_AES128_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_snmptrap_profile.example",
						tfjsonpath.New("version").AtMapKey("v3").AtMapKey("servers").AtSliceIndex(0).AtMapKey("privacy_protocol"),
						knownvalue.StringExact("AES128"),
					),
				},
			},
		},
	})
}

const snmptrapProfile_V3_Legacy_AES128_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_snmptrap_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  version = {
    v3 = {
      servers = [
        {
          name                    = "aes128-server"
          manager                 = "10.0.0.1"
          user                    = "snmpuser"
          authentication_password = "AuthPassword1!"
          privacy_password        = "PrivPassword1!"
          privacy_protocol        = "AES128"
        }
      ]
    }
  }
}
`
