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

func TestAccEmailServerProfile_Basic(t *testing.T) {
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
				Config: emailServerProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("traffic"),
						knownvalue.StringExact("traffic-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("threat"),
						knownvalue.StringExact("threat-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("system"),
						knownvalue.StringExact("system-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("url"),
						knownvalue.StringExact("url-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("data"),
						knownvalue.StringExact("data-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("wildfire"),
						knownvalue.StringExact("wildfire-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("escaping").AtMapKey("escape_character"),
						knownvalue.StringExact("\\"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("escaping").AtMapKey("escaped_characters"),
						knownvalue.StringExact(`'"`),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("smtp-server"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("display_name"),
						knownvalue.StringExact("SMTP Server"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("from"),
						knownvalue.StringExact("alerts@example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("to"),
						knownvalue.StringExact("admin@example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("and_also_to"),
						knownvalue.StringExact("ops@example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("gateway"),
						knownvalue.StringExact("smtp.example.com"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("protocol"),
						knownvalue.StringExact("SMTP"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("port"),
						knownvalue.Int64Exact(587),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("authentication_type"),
						knownvalue.StringExact("Login"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("username"),
						knownvalue.StringExact("smtpuser"),
					),
					// password is hashed (hashing.type: solo) - cannot verify exact value
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("password"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("certificate_profile"),
						knownvalue.StringExact("None"),
					),
				},
			},
		},
	})
}

const emailServerProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_email_server_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  format = {
    traffic  = "traffic-fmt"
    threat   = "threat-fmt"
    system   = "system-fmt"
    url      = "url-fmt"
    data     = "data-fmt"
    wildfire = "wildfire-fmt"
    escaping = {
      escape_character   = "\\"
      escaped_characters = "'\""
    }
  }

  servers = [
    {
      name                 = "smtp-server"
      display_name         = "SMTP Server"
      from                 = "alerts@example.com"
      to                   = "admin@example.com"
      and_also_to          = "ops@example.com"
      gateway              = "smtp.example.com"
      protocol             = "SMTP"
      port                 = 587
      authentication_type  = "Login"
      username             = "smtpuser"
      password             = "SecurePassword123!"
      certificate_profile  = "None"
    }
  ]
}
`

func TestAccEmailServerProfile_TlsProtocol(t *testing.T) {
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
				Config: emailServerProfile_TlsProtocol_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.StringExact("tls-server"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("protocol"),
						knownvalue.StringExact("TLS"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("tls_version"),
						knownvalue.StringExact("1.1"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("port"),
						knownvalue.Int64Exact(465),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers").AtSliceIndex(0).AtMapKey("authentication_type"),
						knownvalue.StringExact("Plain"),
					),
				},
			},
		},
	})
}

const emailServerProfile_TlsProtocol_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_email_server_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  servers = [
    {
      name                = "tls-server"
      from                = "alerts@example.com"
      to                  = "admin@example.com"
      gateway             = "smtp.example.com"
      protocol            = "TLS"
      tls_version         = "1.1"
      port                = 465
      authentication_type = "Plain"
    }
  ]
}
`

func TestAccEmailServerProfile_MultipleServers(t *testing.T) {
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
				Config: emailServerProfile_MultipleServers_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("smtp-server"),
								"display_name":        knownvalue.Null(),
								"from":                knownvalue.StringExact("alerts@example.com"),
								"to":                  knownvalue.StringExact("admin@example.com"),
								"and_also_to":         knownvalue.Null(),
								"gateway":             knownvalue.StringExact("smtp.example.com"),
								"protocol":            knownvalue.StringExact("SMTP"),
								"port":                knownvalue.Int64Exact(25),
								"tls_version":         knownvalue.StringExact("1.2"),
								"authentication_type": knownvalue.StringExact("Auto"),
								"certificate_profile": knownvalue.StringExact("None"),
								"username":            knownvalue.Null(),
								"password":            knownvalue.Null(),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                knownvalue.StringExact("tls-server"),
								"display_name":        knownvalue.Null(),
								"from":                knownvalue.StringExact("alerts@example.com"),
								"to":                  knownvalue.StringExact("admin@example.com"),
								"and_also_to":         knownvalue.Null(),
								"gateway":             knownvalue.StringExact("tls.example.com"),
								"protocol":            knownvalue.StringExact("TLS"),
								"port":                knownvalue.Int64Exact(587),
								"tls_version":         knownvalue.StringExact("1.2"),
								"authentication_type": knownvalue.StringExact("Auto"),
								"certificate_profile": knownvalue.StringExact("None"),
								"username":            knownvalue.Null(),
								"password":            knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const emailServerProfile_MultipleServers_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_email_server_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  servers = [
    {
      name                = "smtp-server"
      from                = "alerts@example.com"
      to                  = "admin@example.com"
      gateway             = "smtp.example.com"
      protocol            = "SMTP"
      port                = 25
      authentication_type = "Auto"
    },
    {
      name                = "tls-server"
      from                = "alerts@example.com"
      to                  = "admin@example.com"
      gateway             = "tls.example.com"
      protocol            = "TLS"
      port                = 587
      tls_version         = "1.2"
      authentication_type = "Auto"
    }
  ]
}
`

func TestAccEmailServerProfile_FormatOnly(t *testing.T) {
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
				Config: emailServerProfile_FormatOnly_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("auth"),
						knownvalue.StringExact("auth-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("config"),
						knownvalue.StringExact("config-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("correlation"),
						knownvalue.StringExact("correlation-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("data"),
						knownvalue.StringExact("data-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("decryption"),
						knownvalue.StringExact("decryption-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("globalprotect"),
						knownvalue.StringExact("globalprotect-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("gtp"),
						knownvalue.StringExact("gtp-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("hip_match"),
						knownvalue.StringExact("hip-match-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("iptag"),
						knownvalue.StringExact("iptag-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("sctp"),
						knownvalue.StringExact("sctp-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("system"),
						knownvalue.StringExact("system-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("threat"),
						knownvalue.StringExact("threat-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("traffic"),
						knownvalue.StringExact("traffic-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("tunnel"),
						knownvalue.StringExact("tunnel-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("url"),
						knownvalue.StringExact("url-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("userid"),
						knownvalue.StringExact("userid-fmt"),
					),
					statecheck.ExpectKnownValue(
						"panos_email_server_profile.example",
						tfjsonpath.New("format").AtMapKey("wildfire"),
						knownvalue.StringExact("wildfire-fmt"),
					),
				},
			},
		},
	})
}

const emailServerProfile_FormatOnly_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_email_server_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix

  format = {
    auth          = "auth-fmt"
    config        = "config-fmt"
    correlation   = "correlation-fmt"
    data          = "data-fmt"
    decryption    = "decryption-fmt"
    globalprotect = "globalprotect-fmt"
    gtp           = "gtp-fmt"
    hip_match     = "hip-match-fmt"
    iptag         = "iptag-fmt"
    sctp          = "sctp-fmt"
    system        = "system-fmt"
    threat        = "threat-fmt"
    traffic       = "traffic-fmt"
    tunnel        = "tunnel-fmt"
    url           = "url-fmt"
    userid        = "userid-fmt"
    wildfire      = "wildfire-fmt"
  }
}
`
