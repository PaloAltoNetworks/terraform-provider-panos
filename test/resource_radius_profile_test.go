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

func TestAccRadiusProfile_Basic(t *testing.T) {
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
				Config: radiusProfile_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("retries"),
						knownvalue.Int64Exact(3),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("timeout"),
						knownvalue.Int64Exact(5),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("radius-server1"),
								"ip_address": knownvalue.StringExact("10.0.1.10"),
								"secret":     knownvalue.StringExact("secret123"),
								"port":       knownvalue.Int64Exact(1812),
							}),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  retries = 3
  timeout = 5
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    }
  ]
}
`

func TestAccRadiusProfile_Protocol_Chap(t *testing.T) {
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
				Config: radiusProfile_Protocol_Chap_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"chap":              knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"eap_ttls_with_pap": knownvalue.Null(),
							"pap":               knownvalue.Null(),
							"peap_mschapv2":     knownvalue.Null(),
							"peap_with_gtc":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_Protocol_Chap_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  protocol = {
    chap = {}
  }
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    }
  ]
}
`

func TestAccRadiusProfile_Protocol_EapTtlsWithPap(t *testing.T) {
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
				Config: radiusProfile_Protocol_EapTtlsWithPap_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"chap": knownvalue.Null(),
							"eap_ttls_with_pap": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"anonymous_outer_identity":   knownvalue.Bool(true),
								"radius_certificate_profile": knownvalue.StringExact("radius-cert-profile"),
							}),
							"pap":           knownvalue.Null(),
							"peap_mschapv2": knownvalue.Null(),
							"peap_with_gtc": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_Protocol_EapTtlsWithPap_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "radius-cert-profile"
  use_ocsp = true
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_certificate_profile.example]
  location = var.location

  name = var.prefix
  protocol = {
    eap_ttls_with_pap = {
      anonymous_outer_identity = true
      radius_certificate_profile = panos_certificate_profile.example.name
    }
  }
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    }
  ]
}
`

func TestAccRadiusProfile_Protocol_Pap(t *testing.T) {
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
				Config: radiusProfile_Protocol_Pap_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"chap":              knownvalue.Null(),
							"eap_ttls_with_pap": knownvalue.Null(),
							"pap":               knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"peap_mschapv2":     knownvalue.Null(),
							"peap_with_gtc":     knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_Protocol_Pap_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  protocol = {
    pap = {}
  }
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    }
  ]
}
`

func TestAccRadiusProfile_Protocol_PeapMschapv2(t *testing.T) {
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
				Config: radiusProfile_Protocol_PeapMschapv2_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"chap":              knownvalue.Null(),
							"eap_ttls_with_pap": knownvalue.Null(),
							"pap":               knownvalue.Null(),
							"peap_mschapv2": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"allow_password_change":      knownvalue.Bool(true),
								"anonymous_outer_identity":   knownvalue.Bool(true),
								"radius_certificate_profile": knownvalue.StringExact("radius-cert-profile"),
							}),
							"peap_with_gtc": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_Protocol_PeapMschapv2_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "radius-cert-profile"
  use_ocsp = true
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_certificate_profile.example]
  location = var.location

  name = var.prefix
  protocol = {
    peap_mschapv2 = {
      allow_password_change = true
      anonymous_outer_identity = true
      radius_certificate_profile = panos_certificate_profile.example.name
    }
  }
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    }
  ]
}
`

func TestAccRadiusProfile_Protocol_PeapWithGtc(t *testing.T) {
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
				Config: radiusProfile_Protocol_PeapWithGtc_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("protocol"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"chap":              knownvalue.Null(),
							"eap_ttls_with_pap": knownvalue.Null(),
							"pap":               knownvalue.Null(),
							"peap_mschapv2":     knownvalue.Null(),
							"peap_with_gtc": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"anonymous_outer_identity":   knownvalue.Bool(true),
								"radius_certificate_profile": knownvalue.StringExact("radius-cert-profile"),
							}),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_Protocol_PeapWithGtc_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_certificate_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = "radius-cert-profile"
  use_ocsp = true
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_certificate_profile.example]
  location = var.location

  name = var.prefix
  protocol = {
    peap_with_gtc = {
      anonymous_outer_identity = true
      radius_certificate_profile = panos_certificate_profile.example.name
    }
  }
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    }
  ]
}
`

func TestAccRadiusProfile_MultipleServers(t *testing.T) {
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
				Config: radiusProfile_MultipleServers_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_radius_profile.example",
						tfjsonpath.New("servers"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("radius-server1"),
								"ip_address": knownvalue.StringExact("10.0.1.10"),
								"secret":     knownvalue.StringExact("secret123"),
								"port":       knownvalue.Int64Exact(1812),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("radius-server2"),
								"ip_address": knownvalue.StringExact("10.0.1.11"),
								"secret":     knownvalue.StringExact("secret456"),
								"port":       knownvalue.Int64Exact(1813),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":       knownvalue.StringExact("radius-server3"),
								"ip_address": knownvalue.StringExact("10.0.1.12"),
								"secret":     knownvalue.StringExact("secret789"),
								"port":       knownvalue.Int64Exact(1812),
							}),
						}),
					),
				},
			},
		},
	})
}

const radiusProfile_MultipleServers_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_radius_profile" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  servers = [
    {
      name = "radius-server1"
      ip_address = "10.0.1.10"
      secret = "secret123"
      port = 1812
    },
    {
      name = "radius-server2"
      ip_address = "10.0.1.11"
      secret = "secret456"
      port = 1813
    },
    {
      name = "radius-server3"
      ip_address = "10.0.1.12"
      secret = "secret789"
      port = 1812
    }
  ]
}
`
