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

func TestAccGlobalProtectPortal_ExternalGateway_Fqdn(t *testing.T) {
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
				Config: globalProtectPortal_ExternalGateway_Fqdn_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"agent_user_override_key": knownvalue.StringExact("abcd"),
							"configs": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":                                knownvalue.StringExact("config1"),
									"save_user_credentials":               knownvalue.StringExact("1"),
									"portal_2fa":                          knownvalue.Bool(true),
									"internal_gateway_2fa":                knownvalue.Bool(true),
									"auto_discovery_external_gateway_2fa": knownvalue.Bool(true),
									"manual_only_gateway_2fa":             knownvalue.Bool(true),
									"refresh_config":                      knownvalue.Bool(true),
									"mdm_address":                         knownvalue.StringExact("mdm.example.com"),
									"mdm_enrollment_port":                 knownvalue.StringExact("443"),
									"source_user":                         knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("user1"), knownvalue.StringExact("user2")}),
									"third_party_vpn_clients":             knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("client1"), knownvalue.StringExact("client2")}),
									"os":                                  knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("windows"), knownvalue.StringExact("mac")}),
									"gateways": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"external": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"cutoff_time": knownvalue.Int64Exact(10),
											"list": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name":          knownvalue.StringExact("external-gateway1"),
													"fqdn":          knownvalue.StringExact("external.example.com"),
													"ip":            knownvalue.Null(),
													"priority_rule": knownvalue.Null(),
													"manual":        knownvalue.Bool(true),
												}),
											}),
										}),
										"internal": knownvalue.Null(),
									}),
									"internal_host_detection": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"ip_address": knownvalue.StringExact("192.168.1.1"),
										"hostname":   knownvalue.StringExact("internal.example.com"),
									}),
									"internal_host_detection_v6": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"ip_address": knownvalue.StringExact("2001:db8::1"),
										"hostname":   knownvalue.StringExact("internal-v6.example.com"),
									}),
									"agent_ui": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"passcode":                    knownvalue.StringExact("123456"),
										"uninstall_password":          knownvalue.StringExact("uninstall123"),
										"agent_user_override_timeout": knownvalue.Int64Exact(60),
										"max_agent_user_overrides":    knownvalue.Int64Exact(5),
										"welcome_page":                knownvalue.Null(),
									}),
									"hip_collection": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"max_wait_time":    knownvalue.Int64Exact(30),
										"collect_hip_data": knownvalue.Bool(true),
										"exclusion": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"category": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name": knownvalue.StringExact("category1"),
													"vendor": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":    knownvalue.StringExact("vendor1"),
															"product": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("product1"), knownvalue.StringExact("product2")}),
														}),
													}),
												}),
											}),
										}),
										"certificate_profile": knownvalue.Null(),
										"custom_checks": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"windows": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"registry_key": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectExact(map[string]knownvalue.Check{
														"name": knownvalue.StringExact("hip-reg-key1"),
														"registry_value": knownvalue.ListExact([]knownvalue.Check{
															knownvalue.StringExact("hip-reg-value1"),
															knownvalue.StringExact("hip-reg-value2"),
														}),
													}),
												}),
												"process_list": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.StringExact("hip-process1"),
													knownvalue.StringExact("hip-process2"),
												}),
											}),
											"mac_os": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"plist": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.ObjectExact(map[string]knownvalue.Check{
														"name": knownvalue.StringExact("hip-plist1"),
														"key": knownvalue.ListExact([]knownvalue.Check{
															knownvalue.StringExact("hip-key1"),
															knownvalue.StringExact("hip-key2"),
														}),
													}),
												}),
												"process_list": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.StringExact("hip-mac-process1"),
													knownvalue.StringExact("hip-mac-process2"),
												}),
											}),
											"linux": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"process_list": knownvalue.ListExact([]knownvalue.Check{
													knownvalue.StringExact("hip-linux-process1"),
													knownvalue.StringExact("hip-linux-process2"),
												}),
											}),
										}),
									}),
									"agent_config": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
									"gp_app_config": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"config": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":  knownvalue.StringExact("app-config1"),
												"value": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("app-value1"), knownvalue.StringExact("app-value2")}),
											}),
										}),
									}),
									"authentication_override": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"generate_cookie": knownvalue.Bool(true),
										"accept_cookie": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"cookie_lifetime": knownvalue.ObjectExact(map[string]knownvalue.Check{
												"lifetime_in_hours":   knownvalue.Int64Exact(24),
												"lifetime_in_days":    knownvalue.Null(),
												"lifetime_in_minutes": knownvalue.Null(),
											}),
										}),
										"cookie_encrypt_decrypt_cert": knownvalue.Null(),
									}),
									"machine_account_exists_with_serialno": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"no":  knownvalue.ObjectExact(map[string]knownvalue.Check{}),
										"yes": knownvalue.Null(),
									}),
									"certificate":        knownvalue.Null(),
									"client_certificate": knownvalue.Null(),
									"custom_checks": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"criteria": knownvalue.ObjectExact(map[string]knownvalue.Check{
											"registry_key": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name":               knownvalue.StringExact("reg1"),
													"default_value_data": knownvalue.StringExact("default_data"),
													"negate":             knownvalue.Bool(false),
													"registry_value": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":       knownvalue.StringExact("value1"),
															"value_data": knownvalue.StringExact("data1"),
															"negate":     knownvalue.Bool(false),
														}),
													}),
												}),
											}),
											"plist": knownvalue.ListExact([]knownvalue.Check{
												knownvalue.ObjectExact(map[string]knownvalue.Check{
													"name":   knownvalue.StringExact("plist1"),
													"negate": knownvalue.Bool(false),
													"key": knownvalue.ListExact([]knownvalue.Check{
														knownvalue.ObjectExact(map[string]knownvalue.Check{
															"name":   knownvalue.StringExact("key1"),
															"value":  knownvalue.StringExact("value1"),
															"negate": knownvalue.Bool(false),
														}),
													}),
												}),
											}),
										}),
									}),
								}),
							}),
							"root_ca": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ExternalGateway_Fqdn_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  client_config = {
    agent_user_override_key = "abcd"
    configs = [
      {
        name = "config1"
        save_user_credentials = "1"
        portal_2fa = true
        internal_gateway_2fa = true
        auto_discovery_external_gateway_2fa = true
        manual_only_gateway_2fa = true
        refresh_config = true
        mdm_address = "mdm.example.com"
        mdm_enrollment_port = "443"
        source_user = ["user1", "user2"]
        third_party_vpn_clients = ["client1", "client2"]
        os = ["windows", "mac"]
        gateways = {
          external = {
            cutoff_time = 10
            list = [
              {
                name = "external-gateway1"
                fqdn = "external.example.com"
                priority = "1"
                manual = true
              }
            ]
          }
        }
        internal_host_detection = {
          ip_address = "192.168.1.1"
          hostname = "internal.example.com"
        }
        internal_host_detection_v6 = {
          ip_address = "2001:db8::1"
          hostname = "internal-v6.example.com"
        }
        agent_ui = {
          passcode = "123456"
          uninstall_password = "uninstall123"
          agent_user_override_timeout = 60
          max_agent_user_overrides = 5
        }
        hip_collection = {
          max_wait_time = 30
          collect_hip_data = true
          exclusion = {
            category = [
              {
                name = "category1"
                vendor = [
                  {
                    name = "vendor1"
                    product = ["product1", "product2"]
                  }
                ]
              }
            ]
          }
          custom_checks = {
            windows = {
              registry_key = [
                {
                  name = "hip-reg-key1"
                  registry_value = ["hip-reg-value1", "hip-reg-value2"]
                }
              ]
              process_list = ["hip-process1", "hip-process2"]
            }
            mac_os = {
              plist = [
                {
                  name = "hip-plist1"
                  key = ["hip-key1", "hip-key2"]
                }
              ]
              process_list = ["hip-mac-process1", "hip-mac-process2"]
            }
            linux = {
              process_list = ["hip-linux-process1", "hip-linux-process2"]
            }
          }
        }
        agent_config = {}
        gp_app_config = {
          config = [
            {
              name = "app-config1"
              value = ["app-value1", "app-value2"]
            }
          ]
        }
        authentication_override = {
          generate_cookie = true
          accept_cookie = {
            cookie_lifetime = {
              lifetime_in_hours = 24
            }
          }
        }
        machine_account_exists_with_serialno = {
          no = {}
        }
        certificate = null
        custom_checks = {
          criteria = {
            registry_key = [
              {
                name = "reg1"
                default_value_data = "default_data"
                negate = false
                registry_value = [
                  {
                    name = "value1"
                    value_data = "data1"
                    negate = false
                  }
                ]
              }
            ]
            plist = [
              {
                name = "plist1"
                negate = false
                key = [
                  {
                    name = "key1"
                    value = "value1"
                    negate = false
                  }
                ]
              }
            ]
          }
        }
      }
    ]
  }
}
`

func TestAccGlobalProtectPortal_ExternalGateway_Ip(t *testing.T) {
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
				Config: globalProtectPortal_ExternalGateway_Ip_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config").AtMapKey("configs").AtSliceIndex(0).AtMapKey("gateways").AtMapKey("external").AtMapKey("list").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("external-gateway1"),
							"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"ipv4": knownvalue.StringExact("192.0.2.1"),
								"ipv6": knownvalue.StringExact("2001:db8::1"),
							}),
							"fqdn":          knownvalue.Null(),
							"priority_rule": knownvalue.Null(),
							"manual":        knownvalue.Bool(true),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ExternalGateway_Ip_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  client_config = {
    configs = [
      {
        name = "config1"
        gateways = {
          external = {
            list = [
              {
                name = "external-gateway1"
                ip = {
                  ipv4 = "192.0.2.1"
                  ipv6 = "2001:db8::1"
                }
                priority = "1"
                manual = true
              }
            ]
          }
        }
      }
    ]
  }
}
`

const globalProtectPortal_CookieLifetime_Days_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  name = var.prefix
  location = var.location

  client_config = {
    configs = [
      {
        name = "config1"
        authentication_override = {
          generate_cookie = true
          accept_cookie = {
            cookie_lifetime = {
              lifetime_in_days = 30
            }
          }
        }
      }
    ]
  }
}
`

func TestAccGlobalProtectPortal_CookieLifetime_Days(t *testing.T) {
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
				Config: globalProtectPortal_CookieLifetime_Days_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config").AtMapKey("configs").AtSliceIndex(0).AtMapKey("authentication_override").AtMapKey("accept_cookie").AtMapKey("cookie_lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"lifetime_in_days":    knownvalue.Int64Exact(30),
							"lifetime_in_hours":   knownvalue.Null(),
							"lifetime_in_minutes": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_CookieLifetime_Minutes_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  client_config = {
    configs = [
      {
        name = "config1"
        authentication_override = {
          generate_cookie = true
          accept_cookie = {
            cookie_lifetime = {
              lifetime_in_minutes = 45
            }
          }
        }
      }
    ]
  }
}
`

func TestAccGlobalProtectPortal_CookieLifetime_Minutes(t *testing.T) {
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
				Config: globalProtectPortal_CookieLifetime_Minutes_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config").AtMapKey("configs").AtSliceIndex(0).AtMapKey("authentication_override").AtMapKey("accept_cookie").AtMapKey("cookie_lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"lifetime_in_minutes": knownvalue.Int64Exact(45),
							"lifetime_in_days":    knownvalue.Null(),
							"lifetime_in_hours":   knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_MachineAccount_Yes_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  client_config = {
    configs = [
      {
        name = "config1"
        machine_account_exists_with_serialno = {
          yes = {}
        }
      }
    ]
  }
}
`

func TestAccGlobalProtectPortal_MachineAccount_Yes(t *testing.T) {
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
				Config: globalProtectPortal_MachineAccount_Yes_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config").AtMapKey("configs").AtSliceIndex(0).AtMapKey("machine_account_exists_with_serialno"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"no":  knownvalue.Null(),
							"yes": knownvalue.ObjectExact(map[string]knownvalue.Check{}),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ClientCertificateLocal_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  client_config = {
    configs = [
      {
        name = "config1"
        client_certificate = {
          local = "local-cert"
        }
      }
    ]
  }
}
`

func TestAccGlobalProtectPortal_ClientCertificateLocal(t *testing.T) {
	t.Parallel()
	t.Skip("Missing certificate resource support")

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
				Config: globalProtectPortal_ClientCertificateLocal_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config").AtMapKey("configs").AtSliceIndex(0).AtMapKey("client_certificate"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"local": knownvalue.StringExact("local-cert"),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ClientCertificateScep_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  client_config = {
    configs = [
      {
        name = "config1"
        client_certificate = {
          scep = "scep-profile"
        }
      }
    ]
  }
}
`

func TestAccGlobalProtectPortal_ClientCertificateScep(t *testing.T) {
	t.Parallel()
	t.Skip("missing scep profile resource")

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
				Config: globalProtectPortal_ClientCertificateScep_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("client_config").AtMapKey("configs").AtSliceIndex(0).AtMapKey("client_certificate"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"scep": knownvalue.StringExact("scep-profile"),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ClientlessVpn_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "zone_location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_zone" "example" {
  depends_on = [panos_template.example]

  location = var.zone_location
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  location = var.location

  name = var.prefix

  clientless_vpn = {
    crypto_settings = {
      server_cert_verification = {
        block_expired_certificate = true
        block_timeout_cert = true
        block_unknown_cert = true
        block_untrusted_issuer = true
      }
      ssl_protocol = {
        auth_algo_md5 = true
        auth_algo_sha1 = true
        auth_algo_sha256 = true
        auth_algo_sha384 = true
        enc_algo_3des = true
        enc_algo_aes_128_cbc = true
        enc_algo_aes_128_gcm = true
        enc_algo_aes_256_cbc = true
        enc_algo_aes_256_gcm = true
        enc_algo_rc4 = true
        keyxchg_algo_dhe = true
        keyxchg_algo_ecdhe = true
        keyxchg_algo_rsa = true
        max_version = "max"
        min_version = "tls1-0"
      }
    }
    hostname = "clientless.example.com"
    inactivity_logout = {
      minutes = 30
    }
    login_lifetime = {
      hours = 3
    }
    max_user = 1000
    proxy_server_setting = [
      {
        name = "proxy1"
        domains = ["domain1.com", "domain2.com"]
        use_proxy = true
        proxy_server = {
          server = "proxy.example.com"
          port = 8080
          user = "proxyuser"
          password = "proxypass"
        }
      }
    ]
    rewrite_exclude_domain_list = ["exclude1.com", "exclude2.com"]
    security_zone = panos_zone.example.name
  }
}
`

func TestAccGlobalProtectPortal_ClientlessVpn_Basic(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	zoneLocation := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectPortal_ClientlessVpn_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":        config.StringVariable(prefix),
					"location":      location,
					"zone_location": zoneLocation,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("clientless_vpn"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"crypto_settings": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"server_cert_verification": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"block_expired_certificate": knownvalue.Bool(true),
									"block_timeout_cert":        knownvalue.Bool(true),
									"block_unknown_cert":        knownvalue.Bool(true),
									"block_untrusted_issuer":    knownvalue.Bool(true),
								}),
								"ssl_protocol": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"auth_algo_md5":        knownvalue.Bool(true),
									"auth_algo_sha1":       knownvalue.Bool(true),
									"auth_algo_sha256":     knownvalue.Bool(true),
									"auth_algo_sha384":     knownvalue.Bool(true),
									"enc_algo_3des":        knownvalue.Bool(true),
									"enc_algo_aes_128_cbc": knownvalue.Bool(true),
									"enc_algo_aes_128_gcm": knownvalue.Bool(true),
									"enc_algo_aes_256_cbc": knownvalue.Bool(true),
									"enc_algo_aes_256_gcm": knownvalue.Bool(true),
									"enc_algo_rc4":         knownvalue.Bool(true),
									"keyxchg_algo_dhe":     knownvalue.Bool(true),
									"keyxchg_algo_ecdhe":   knownvalue.Bool(true),
									"keyxchg_algo_rsa":     knownvalue.Bool(true),
									"max_version":          knownvalue.StringExact("max"),
									"min_version":          knownvalue.StringExact("tls1-0"),
								}),
							}),
							"hostname": knownvalue.StringExact("clientless.example.com"),
							"inactivity_logout": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"minutes": knownvalue.Int64Exact(30),
								"hours":   knownvalue.Null(),
							}),
							"login_lifetime": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"hours":   knownvalue.Int64Exact(3),
								"minutes": knownvalue.Null(),
							}),
							"max_user": knownvalue.Int64Exact(1000),
							"proxy_server_setting": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":      knownvalue.StringExact("proxy1"),
									"domains":   knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("domain1.com"), knownvalue.StringExact("domain2.com")}),
									"use_proxy": knownvalue.Bool(true),
									"proxy_server": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"server":   knownvalue.StringExact("proxy.example.com"),
										"port":     knownvalue.Int64Exact(8080),
										"user":     knownvalue.StringExact("proxyuser"),
										"password": knownvalue.StringExact("proxypass"),
									}),
								}),
							}),
							"rewrite_exclude_domain_list": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("exclude1.com"), knownvalue.StringExact("exclude2.com")}),
							"security_zone":               knownvalue.StringExact(prefix),
							"apps_to_user_mapping":        knownvalue.Null(),
							"dns_proxy":                   knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ClientlessVpn_LogoutLoginVariants_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  location = var.location

  name = var.prefix

  clientless_vpn = {
    inactivity_logout = {
      hours = 2
    }
    login_lifetime = {
      minutes = 90
    }
    // Other required fields
    hostname = "clientless.example.com"
    max_user = 1000
  }
}
`

func TestAccGlobalProtectPortal_ClientlessVpn_LogoutLoginVariants(t *testing.T) {
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
				Config: globalProtectPortal_ClientlessVpn_LogoutLoginVariants_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("clientless_vpn").AtMapKey("inactivity_logout"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"hours":   knownvalue.Int64Exact(2),
							"minutes": knownvalue.Null(),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("clientless_vpn").AtMapKey("login_lifetime"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"minutes": knownvalue.Int64Exact(90),
							"hours":   knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_ClientlessVpn_AppsToUserMapping_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }
variable "app_location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_device_group" "example" {
  location = { panorama = {} }
  name = var.prefix

  templates = [panos_template.example.name]
}

resource "panos_application" "app1" {
  depends_on = [panos_device_group.example]

  location = var.app_location
  name = "${var.prefix}-app1"
  category = "business-systems"
  subcategory = "auth-service"
  technology = "client-server"
  risk = 1
}

resource "panos_application" "app2" {
  depends_on = [panos_device_group.example]

  location = var.app_location
  name = "${var.prefix}-app2"
  category = "collaboration"
  subcategory = "email"
  technology = "browser-based"
  risk = 2
}

resource "panos_globalprotect_portal" "example" {
  name = var.prefix
  location = var.location

  clientless_vpn = {
    apps_to_user_mapping = [
      {
        name = "app-mapping1"
        source_user = ["user1", "user2"]
        applications = [panos_application.app1.name, panos_application.app2.name]
        enable_custom_app_URL_address_bar = true
        display_global_protect_agent_download_link = true
      },
      {
        name = "app-mapping2"
        source_user = ["user3"]
        applications = [panos_application.app1.name]
        enable_custom_app_URL_address_bar = false
        display_global_protect_agent_download_link = false
      }
    ]
    // Other required fields
    hostname = "clientless.example.com"
    max_user = 1000
  }
}
`

func TestAccGlobalProtectPortal_ClientlessVpn_AppsToUserMapping(t *testing.T) {
	t.Parallel()
	t.Skip("disabled, requires more testing around applications")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template_vsys": config.ObjectVariable(map[string]config.Variable{
			"template": config.StringVariable(prefix),
		}),
	})

	appLocation := config.ObjectVariable(map[string]config.Variable{
		"device_group": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: globalProtectPortal_ClientlessVpn_AppsToUserMapping_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":       config.StringVariable(prefix),
					"location":     location,
					"app_location": appLocation,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("clientless_vpn").AtMapKey("apps_to_user_mapping"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                              knownvalue.StringExact("app-mapping1"),
								"source_user":                       knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("user1"), knownvalue.StringExact("user2")}),
								"applications":                      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact(prefix + "-app1"), knownvalue.StringExact(prefix + "-app2")}),
								"enable_custom_app_URL_address_bar": knownvalue.Bool(true),
								"display_global_protect_agent_download_link": knownvalue.Bool(true),
							}),
							knownvalue.ObjectExact(map[string]knownvalue.Check{
								"name":                              knownvalue.StringExact("app-mapping2"),
								"source_user":                       knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("user3")}),
								"applications":                      knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact(prefix + "-app1")}),
								"enable_custom_app_URL_address_bar": knownvalue.Bool(false),
								"display_global_protect_agent_download_link": knownvalue.Bool(false),
							}),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_PortalConfig_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name} }

  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "192.168.1.1" }]
    ipv6 = {
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_globalprotect_portal" "example" {
  name = var.prefix
  location = var.location

  portal_config = {
    #certificate_profile = "portal-cert-profile"
    client_auth = [
      {
        name = "client-auth1"
        os = "Any"
        #authentication_profile = "auth-profile1"
        auto_retrieve_passcode = true
        username_label = "Username"
        password_label = "Password"
        authentication_message = "Enter login credentials"
        user_credential_or_client_cert_required = "no"
      }
    ]
    config_selection = {
      #certificate_profile = "config-cert-profile"
      custom_checks = {
        mac_os = {
          plist = [
            {
              name = "plist1"
              key = ["key1", "key2"]
            }
          ]
        }
        windows = {
          registry_key = [
            {
              name = "reg1"
              registry_value = ["value1", "value2"]
            }
          ]
        }
      }
    }
    #custom_help_page = "help.html"
    #custom_home_page = "home.html"
    #custom_login_page = "login.html"
    local_address = {
      interface = panos_ethernet_interface.example.name
      ip_address_family = "ipv4"
      ip = {
        ipv4 = "192.168.1.1"
        ipv6 = "2001:db8::1"
      }
    }
    log_fail = true
    #log_setting = "portal-log-setting"
    log_success = true
    #ssl_tls_service_profile = "portal-ssl-profile"
  }
}
`

func TestAccGlobalProtectPortal_PortalConfig_Basic(t *testing.T) {
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
				Config: globalProtectPortal_PortalConfig_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("portal_config"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"certificate_profile": knownvalue.Null(),
							"client_auth": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":                   knownvalue.StringExact("client-auth1"),
									"os":                     knownvalue.StringExact("Any"),
									"authentication_profile": knownvalue.Null(),
									"auto_retrieve_passcode": knownvalue.Bool(true),
									"username_label":         knownvalue.StringExact("Username"),
									"password_label":         knownvalue.StringExact("Password"),
									"authentication_message": knownvalue.StringExact("Enter login credentials"),
									"user_credential_or_client_cert_required": knownvalue.StringExact("no"),
								}),
							}),
							"config_selection": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"certificate_profile": knownvalue.Null(),
								"custom_checks": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"mac_os": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"plist": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name": knownvalue.StringExact("plist1"),
												"key":  knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("key1"), knownvalue.StringExact("key2")}),
											}),
										}),
									}),
									"windows": knownvalue.ObjectExact(map[string]knownvalue.Check{
										"registry_key": knownvalue.ListExact([]knownvalue.Check{
											knownvalue.ObjectExact(map[string]knownvalue.Check{
												"name":           knownvalue.StringExact("reg1"),
												"registry_value": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("value1"), knownvalue.StringExact("value2")}),
											}),
										}),
									}),
								}),
							}),
							"custom_help_page":  knownvalue.Null(),
							"custom_home_page":  knownvalue.StringExact("factory-default"),
							"custom_login_page": knownvalue.StringExact("factory-default"),
							"local_address": knownvalue.ObjectExact(map[string]knownvalue.Check{
								"interface":         knownvalue.StringExact("ethernet1/1"),
								"ip_address_family": knownvalue.StringExact("ipv4"),
								"ip": knownvalue.ObjectExact(map[string]knownvalue.Check{
									"ipv4": knownvalue.StringExact("192.168.1.1"),
									"ipv6": knownvalue.StringExact("2001:db8::1"),
								}),
								"floating_ip": knownvalue.Null(),
							}),
							"log_fail":                knownvalue.Bool(true),
							"log_setting":             knownvalue.Null(),
							"log_success":             knownvalue.Bool(true),
							"ssl_tls_service_profile": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_PortalConfig_FloatingIp_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_ethernet_interface" "example" {
  location = { template = { name = panos_template.example.name} }

  name = "ethernet1/1"
  layer3 = {
    ips = [{ name = "192.168.1.1" }]
    ipv6 = {
      addresses = [{ name = "2001:db8::1" }]
    }
  }
}

resource "panos_globalprotect_portal" "example" {
  name = var.prefix
  location = var.location

  portal_config = {
    local_address = {
      interface = panos_ethernet_interface.example.name
      ip_address_family = "ipv4_ipv6"
      floating_ip = {
        ipv4 = "192.168.1.10"
        ipv6 = "2001:db8::10"
      }
    }
  }
}
`

func TestAccGlobalProtectPortal_PortalConfig_FloatingIp(t *testing.T) {
	t.Parallel()
	t.Skip("requires floating ip support")

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
				Config: globalProtectPortal_PortalConfig_FloatingIp_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("portal_config").AtMapKey("local_address").AtMapKey("floating_ip"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"ipv4": knownvalue.StringExact("192.168.1.10"),
							"ipv6": knownvalue.StringExact("2001:db8::10"),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_SatelliteConfig_Basic_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }


resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on = [panos_template.example]

  name = var.prefix
  location = var.location

  satellite_config = {
    configs = [
      {
        name = "satellite-config1"
        devices = ["device1", "device2"]
        source_user = ["user1", "user2"]
        gateways = [
          {
            name = "gateway1"
            ipv6_preferred = true
            priority = 1
            fqdn = "gateway1.example.com"
          }
        ]
        config_refresh_interval = 24
      }
    ]
    root_ca = ["root-ca1", "root-ca2"]
  }
}
`

func TestAccGlobalProtectPortal_SatelliteConfig_Basic(t *testing.T) {
	t.Parallel()
	t.Skip("missing support for satellite devices")

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
				Config: globalProtectPortal_SatelliteConfig_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("satellite_config"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"configs": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.ObjectExact(map[string]knownvalue.Check{
									"name":        knownvalue.StringExact("satellite-config1"),
									"devices":     knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("device1"), knownvalue.StringExact("device2")}),
									"source_user": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("user1"), knownvalue.StringExact("user2")}),
									"gateways": knownvalue.ListExact([]knownvalue.Check{
										knownvalue.ObjectExact(map[string]knownvalue.Check{
											"name":           knownvalue.StringExact("gateway1"),
											"ipv6_preferred": knownvalue.Bool(true),
											"priority":       knownvalue.Int64Exact(1),
											"fqdn":           knownvalue.StringExact("gateway1.example.com"),
										}),
									}),
									"config_refresh_interval": knownvalue.Int64Exact(24),
								}),
							}),
							"root_ca": knownvalue.ListExact([]knownvalue.Check{knownvalue.StringExact("root-ca1"), knownvalue.StringExact("root-ca2")}),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_SatelliteConfig_ClientCertificateLocal_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  name = var.prefix
  location = var.location

  satellite_config = {
    client_certificate = {
      local = {
        certificate_life_time = 30
        certificate_renewal_period = 7
        issuing_certificate = "issuing-cert"
        ocsp_responder = "ocsp-responder"
      }
    }
  }
}
`

func TestAccGlobalProtectPortal_SatelliteConfig_ClientCertificateLocal(t *testing.T) {
	t.Parallel()
	t.Skip("missing local certificates resources")

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
				Config: globalProtectPortal_SatelliteConfig_ClientCertificateLocal_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("satellite_config").AtMapKey("client_certificate").AtMapKey("local"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"certificate_life_time":      knownvalue.Int64Exact(30),
							"certificate_renewal_period": knownvalue.Int64Exact(7),
							"issuing_certificate":        knownvalue.StringExact("issuing-cert"),
							"ocsp_responder":             knownvalue.StringExact("ocsp-responder"),
						}),
					),
				},
			},
		},
	})
}

const globalProtectPortal_SatelliteConfig_ClientCertificateScep_Tmpl = `
variable "prefix" { type = string }
variable "location" { type = any }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_globalprotect_portal" "example" {
  depends_on =  [panos_template.example]
  name = var.prefix
  location = var.location

  satellite_config = {
    client_certificate = {
      scep = {
        certificate_renewal_period = 5
        scep = "scep-profile"
      }
    }
  }
}
`

func TestAccGlobalProtectPortal_SatelliteConfig_ClientCertificateScep(t *testing.T) {
	t.Parallel()
	t.Skip("missing scep profile resource")

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
				Config: globalProtectPortal_SatelliteConfig_ClientCertificateScep_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_globalprotect_portal.example",
						tfjsonpath.New("satellite_config").AtMapKey("client_certificate").AtMapKey("scep"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"certificate_renewal_period": knownvalue.Int64Exact(5),
							"scep":                       knownvalue.StringExact("scep-profile"),
						}),
					),
				},
			},
		},
	})
}
