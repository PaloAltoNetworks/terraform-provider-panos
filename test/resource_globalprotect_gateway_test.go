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

func TestAccGlobalProtectGateway_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGatewayConfig,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("block_quarantined_devices"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("certificate_profile"),
						knownvalue.StringExact("cert-profile"),
					),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("client_auth"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                   knownvalue.StringExact("client-auth-1"),
								"os":                     knownvalue.StringExact("any"),
								"authentication_profile": knownvalue.StringExact("auth-profile"),
								"auto_retrieve_passcode": knownvalue.Bool(true),
								"username_label":         knownvalue.StringExact("Username"),
								"password_label":         knownvalue.StringExact("Password"),
								"authentication_message": knownvalue.StringExact("Enter login credentials"),
								"user_credential_or_client_cert_required": knownvalue.StringExact("yes"),
							}),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("log_success"),
						knownvalue.Bool(true),
					),
					// statecheck.ExpectKnownValue(
					// 	"panos_globalprotect_gateway.example",
					// 	tfjsonpath.New("remote_user_tunnel"),
					// 	knownvalue.StringExact("tunnel.1"),
					// ),
					// statecheck.ExpectKnownValue(
					// 	"panos_globalprotect_gateway.example",
					// 	tfjsonpath.New("satellite_tunnel"),
					// 	knownvalue.StringExact("tunnel.2"),
					// ),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("tunnel_mode"),
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

const globalProtectGatewayConfig = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_authentication_profile" "auth_profile" {
  location = { template = { name = panos_template.example.name } }
  name = "auth-profile"
  user_domain = "example.com"
}

resource "panos_certificate_profile" "cert_profile" {
  location = { template = { name = panos_template.example.name } }
  name = "cert-profile"
  username_field = {
    subject = "common-name"
  }
}

resource "panos_ethernet_interface" "eth" {
  location = { template = { name = panos_template.example.name } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
    ipv6 = {
      enabled = true
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_tunnel_interface" "tunnel1" {
  location = { template = { name = panos_template.example.name } }
  name = "tunnel.1"
}

resource "panos_tunnel_interface" "tunnel2" {
  location = { template = { name = panos_template.example.name } }
  name = "tunnel.2"
}

resource "panos_globalprotect_gateway" "example" {
  depends_on =  [panos_template.example]
  location = var.location
  name = var.prefix
  block_quarantined_devices = true
  certificate_profile = panos_certificate_profile.cert_profile.name
  local_address = {
    interface = panos_ethernet_interface.eth.name
    ip_address_family = "ipv4"
    ip = {}
  }
  client_auth = [
    {
      name = "client-auth-1"
      os = "any"
      authentication_profile = panos_authentication_profile.auth_profile.name
      auto_retrieve_passcode = true
      username_label = "Username"
      password_label = "Password"
      authentication_message = "Enter login credentials"
      user_credential_or_client_cert_required = "yes"
    }
  ]
  hip_notification = [
    {
      name = "hip-notification-1"
      match_message = {
        include_app_list = true
        show_notification_as = "pop-up-message"
        message = "HIP match message"
      }
      not_match_message = {
        show_notification_as = "system-tray-balloon"
        message = "HIP not match message"
      }
    }
  ]
  log_fail = true
  log_success = true
  #remote_user_tunnel = panos_tunnel_interface.tunnel1.name
  #satellite_tunnel = panos_tunnel_interface.tunnel2.name
  tunnel_mode = true
}
`

func TestAccGlobalProtectGateway_RemoteUserTunnelConfigs(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_RemoteUserTunnelConfigs_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("remote_user_tunnel_configs"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                          knownvalue.StringExact("config1"),
								"source_user":                   knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("any")}),
								"os":                            knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("any")}),
								"dns_server":                    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("8.8.8.8")}),
								"dns_suffix":                    knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("example.com")}),
								"ip_pool":                       knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.0/24")}),
								"authentication_server_ip_pool": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10.0.1.0/24")}),
								"authentication_override": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"generate_cookie":             knownvalue.Bool(true),
									"cookie_encrypt_decrypt_cert": knownvalue.Null(),
									"accept_cookie": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"cookie_lifetime": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"lifetime_in_hours":   knownvalue.Int64Exact(8),
											"lifetime_in_days":    knownvalue.Null(),
											"lifetime_in_minutes": knownvalue.Null(),
										}),
									}),
								}),
								"source_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"region":     knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("US")}),
									"ip_address": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("any")}),
								}),
								"split_tunneling": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"access_route":         knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("10.0.0.0/8")}),
									"exclude_access_route": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("192.168.1.0/24")}),
									"include_applications": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("any")}),
									"exclude_applications": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("none")}),
									"include_domains": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"list": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":  knownvalue.StringExact("domain1.com"),
												"ports": knownvalue.ListExact([]knownvalue.Check{knownvalue.Int64Exact(80), knownvalue.Int64Exact(443)}),
											}),
										}),
									}),
									"exclude_domains": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"list": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":  knownvalue.StringExact("domain2.com"),
												"ports": knownvalue.ListExact([]knownvalue.Check{knownvalue.Int64Exact(8080)}),
											}),
										}),
									}),
								}),
								"no_direct_access_to_local_network": knownvalue.Bool(true),
								"retrieve_framed_ip_address":        knownvalue.Bool(true),
							}),
						}),
					),
				},
			},
		},
	})
}

const globalProtectGateway_RemoteUserTunnelConfigs_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  remote_user_tunnel_configs = [
    {
      name = "config1"
      source_user = ["any"]
      os = ["any"]
      dns_server = ["8.8.8.8"]
      dns_suffix = ["example.com"]
      ip_pool = ["10.0.0.0/24"]
      authentication_server_ip_pool = ["10.0.1.0/24"]
      authentication_override = {
        generate_cookie = true
        accept_cookie = {
          cookie_lifetime = {
            lifetime_in_hours = 8
          }
        }
      }
      source_address = {
        region = ["US"]
        ip_address = ["any"]
      }
      split_tunneling = {
        access_route = ["10.0.0.0/8"]
        exclude_access_route = ["192.168.1.0/24"]
        include_applications = ["any"]
        exclude_applications = ["none"]
        include_domains = {
          list = [
            {
              name = "domain1.com"
              ports = [80, 443]
            }
          ]
        }
        exclude_domains = {
          list = [
            {
              name = "domain2.com"
              ports = [8080]
            }
          ]
        }
      }
      no_direct_access_to_local_network = true
      retrieve_framed_ip_address = true
    }
  ]
}
`

func TestAccGlobalProtectGateway_RemoteUserTunnelConfigs_CookieLifetime(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_RemoteUserTunnelConfigs_CookieLifetimeDays_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("remote_user_tunnel_configs").AtSliceIndex(0).AtMapKey("authentication_override").AtMapKey("accept_cookie").AtMapKey("cookie_lifetime").AtMapKey("lifetime_in_days"),
						knownvalue.Int64Exact(10),
					),
				},
			},
			{
				Config: globalProtectGateway_RemoteUserTunnelConfigs_CookieLifetimeMinutes_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("remote_user_tunnel_configs").AtSliceIndex(0).AtMapKey("authentication_override").AtMapKey("accept_cookie").AtMapKey("cookie_lifetime").AtMapKey("lifetime_in_minutes"),
						knownvalue.Int64Exact(30),
					),
				},
			},
		},
	})
}

const globalProtectGateway_RemoteUserTunnelConfigs_CookieLifetimeDays_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  remote_user_tunnel_configs = [
    {
      name = "config1"
      authentication_override = {
        generate_cookie = true
        accept_cookie = {
          cookie_lifetime = {
            lifetime_in_days = 10
          }
        }
      }
    }
  ]
}
`

const globalProtectGateway_RemoteUserTunnelConfigs_CookieLifetimeMinutes_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  remote_user_tunnel_configs = [
    {
      name = "config1"
      authentication_override = {
        generate_cookie = true
        accept_cookie = {
          cookie_lifetime = {
            lifetime_in_minutes = 30
          }
        }
      }
    }
  ]
}
`

func TestAccGlobalProtectGateway_Roles(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_Roles_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("roles"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name": knownvalue.StringExact("default"),
								"login_lifetime": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"days":    knownvalue.Int64Exact(1),
									"hours":   knownvalue.Null(),
									"minutes": knownvalue.Null(),
								}),
								"inactivity_logout":           knownvalue.Int64Exact(180),
								"lifetime_notify_prior":       knownvalue.Int64Exact(30),
								"lifetime_notify_message":     knownvalue.StringExact("Your session will expire soon."),
								"inactivity_notify_prior":     knownvalue.Int64Exact(30),
								"inactivity_notify_message":   knownvalue.StringExact("Your session will time out soon."),
								"admin_logout_notify":         knownvalue.Bool(true),
								"admin_logout_notify_message": knownvalue.StringExact("You have been logged out."),
							}),
						}),
					),
				},
			},
		},
	})
}

const globalProtectGateway_Roles_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  roles = [
    {
      name = "default"
      login_lifetime = {
        days = 1
      }
      inactivity_logout = 180
      lifetime_notify_prior = 30
      lifetime_notify_message = "Your session will expire soon."
      inactivity_notify_prior = 30
      inactivity_notify_message = "Your session will time out soon."
      admin_logout_notify = true
      admin_logout_notify_message = "You have been logged out."
    }
  ]
}
`

func TestAccGlobalProtectGateway_Roles_LoginLifetime(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_Roles_LoginLifetimeHours_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("roles").AtSliceIndex(0).AtMapKey("login_lifetime").AtMapKey("hours"),
						knownvalue.Int64Exact(10),
					),
				},
			},
			{
				Config: globalProtectGateway_Roles_LoginLifetimeMinutes_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("roles").AtSliceIndex(0).AtMapKey("login_lifetime").AtMapKey("minutes"),
						knownvalue.Int64Exact(120),
					),
				},
			},
		},
	})
}

const globalProtectGateway_Roles_LoginLifetimeHours_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  roles = [
    {
      name = "default"
      login_lifetime = {
        hours = 10
      }
    }
  ]
}
`

const globalProtectGateway_Roles_LoginLifetimeMinutes_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  roles = [
    {
      name = "default"
      login_lifetime = {
        minutes = 120
      }
    }
  ]
}
`

func TestAccGlobalProtectGateway_SecurityRestrictions(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_SecurityRestrictions_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("security_restrictions"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"disallow_automatic_restoration": knownvalue.Bool(true),
							"source_ip_enforcement": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"enable": knownvalue.Bool(true),
								"custom": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"source_ipv4_netmask": knownvalue.Int64Exact(24),
									"source_ipv6_netmask": knownvalue.Int64Exact(64),
								}),
								"default": knownvalue.Null(),
							}),
						}),
					),
				},
			},
		},
	})
}

const globalProtectGateway_SecurityRestrictions_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  security_restrictions = {
    disallow_automatic_restoration = true
    source_ip_enforcement = {
      enable = true
      custom = {
        source_ipv4_netmask = 24
        source_ipv6_netmask = 64
      }
    }
  }
}
`

func TestAccGlobalProtectGateway_SecurityRestrictions_SourceIpEnforcement(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_SecurityRestrictions_SourceIpEnforcementDefault_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("security_restrictions").AtMapKey("source_ip_enforcement").AtMapKey("custom"),
						knownvalue.Null(),
					),
				},
			},
		},
	})
}

const globalProtectGateway_SecurityRestrictions_SourceIpEnforcementDefault_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  security_restrictions = {
    source_ip_enforcement = {
      enable = true
      default = {}
    }
  }
}
`

func TestAccGlobalProtectGateway_LocalAddress_Ip_Ipv4(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_LocalAddress_Ip_Ipv4_Config,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":         knownvalue.StringExact("ethernet1/1"),
							"ip_address_family": knownvalue.StringExact("ipv4"),
							"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ipv4": knownvalue.StringExact("10.0.0.1/24"),
								"ipv6": knownvalue.Null(),
							}),
							"floating_ip": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccGlobalProtectGateway_LocalAddress_Ip_Ipv6(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_LocalAddress_Ip_Ipv6_Config,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":         knownvalue.StringExact("ethernet1/1"),
							"ip_address_family": knownvalue.StringExact("ipv6"),
							"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ipv4": knownvalue.Null(),
								"ipv6": knownvalue.StringExact("2001:db8::1"),
							}),
							"floating_ip": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

func TestAccGlobalProtectGateway_LocalAddress_Ip_Ipv4_Ipv6(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectGateway_LocalAddress_Ip_Ipv4_Ipv6_Config,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_gateway.example",
						tfjsonpath.New("local_address"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"interface":         knownvalue.StringExact("ethernet1/1"),
							"ip_address_family": knownvalue.StringExact("ipv4_ipv6"),
							"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ipv4": knownvalue.StringExact("10.0.0.1/24"),
								"ipv6": knownvalue.StringExact("2001:db8::1"),
							}),
							"floating_ip": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectGateway_LocalAddress_Ip_Ipv4_Config = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "eth" {
  location = { template = { name = panos_template.example.name } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
    ipv6 = {
      enabled = true
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  local_address = {
    interface = panos_ethernet_interface.eth.name
    ip_address_family = "ipv4"
    ip = {
      ipv4 = "10.0.0.1/24"
    }
  }
}
`

const globalProtectGateway_LocalAddress_Ip_Ipv6_Config = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "eth" {
  location = { template = { name = panos_template.example.name } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
    ipv6 = {
      enabled = true
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  local_address = {
    interface = panos_ethernet_interface.eth.name
    ip_address_family = "ipv6"
    ip = {
      ipv6 = "2001:db8::1"
    }
  }
}
`

const globalProtectGateway_LocalAddress_Ip_Ipv4_Ipv6_Config = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "eth" {
  location = { template = { name = panos_template.example.name } }
  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "10.0.0.1/24" }]
    ipv6 = {
      enabled = true
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_globalprotect_gateway" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  local_address = {
    interface = panos_ethernet_interface.eth.name
    ip_address_family = "ipv4_ipv6"
    ip = {
      ipv4 = "10.0.0.1/24"
      ipv6 = "2001:db8::1"
    }
  }
}
`
