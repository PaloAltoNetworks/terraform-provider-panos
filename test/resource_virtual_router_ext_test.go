package provider_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// TestAccPanosVirtualRouter_BgpExport_AsPath_Remove verifies that when only
// the as_path.remove variant is configured in BGP export rules (without explicitly
// setting others to null), and after removing default values from the spec, the other
// variants (prepend, remove_and_prepend, none) remain null in state.
func TestAccPanosVirtualRouter_BgpExport_AsPath_Remove(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("acc-vrouter-bgp-%s", nameSuffix)
	routerName := fmt.Sprintf("test-vr-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosVirtualRouter_BgpExport_AsPath_Remove_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"router_name":   config.StringVariable(routerName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(routerName),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"prepend":            knownvalue.Null(),  // No default value after removing from spec
							"remove_and_prepend": knownvalue.Null(),  // No default value after removing from spec
						}),
					),
				},
			},
		},
	})
}

// TestAccPanosVirtualRouter_BgpExport_AsPath_Remove_ExplicitNull verifies that
// when explicitly setting prepend and remove_and_prepend to null alongside
// the remove variant, AND after removing default values from the spec, the fields
// remain null in state. This demonstrates that removing spec defaults allows explicit
// nulls to work correctly.
func TestAccPanosVirtualRouter_BgpExport_AsPath_Remove_ExplicitNull(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("acc-vrouter-bgp-%s", nameSuffix)
	routerName := fmt.Sprintf("test-vr-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosVirtualRouter_BgpExport_AsPath_Remove_ExplicitNull_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"router_name":   config.StringVariable(routerName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact(routerName),
					),
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.ObjectExact(map[string]knownvalue.Check{}),
							"prepend":            knownvalue.Null(),  // Now stays null after removing spec defaults
							"remove_and_prepend": knownvalue.Null(),  // Now stays null after removing spec defaults
						}),
					),
				},
			},
		},
	})
}

// TestAccPanosVirtualRouter_BgpExport_AsPath_Prepend_ExplicitNull verifies that
// when setting prepend to a non-default value (3) and explicitly setting other
// variants to null, we check if the other variants remain null or get defaults.
func TestAccPanosVirtualRouter_BgpExport_AsPath_Prepend_ExplicitNull(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("acc-vrouter-bgp-%s", nameSuffix)
	routerName := fmt.Sprintf("test-vr-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosVirtualRouter_BgpExport_AsPath_Prepend_ExplicitNull_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"router_name":   config.StringVariable(routerName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.Null(),
							"prepend":            knownvalue.Int64Exact(3),
							"remove_and_prepend": knownvalue.Null(),  // Check if this stays null
						}),
					),
				},
			},
		},
	})
}

// TestAccPanosVirtualRouter_BgpExport_AsPath_Prepend_Only verifies that when
// setting only prepend to a non-default value in a single rule, the other
// variants (none, remove, remove_and_prepend) remain null.
func TestAccPanosVirtualRouter_BgpExport_AsPath_Prepend_Only(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("acc-vrouter-bgp-%s", nameSuffix)
	routerName := fmt.Sprintf("test-vr-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosVirtualRouter_BgpExport_AsPath_Prepend_Only_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"router_name":   config.StringVariable(routerName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.Null(),
							"prepend":            knownvalue.Int64Exact(3),
							"remove_and_prepend": knownvalue.Null(),
						}),
					),
				},
			},
		},
	})
}

// TestAccPanosVirtualRouter_BgpExport_AsPath_MultipleVariants tests what happens
// when attempting to set multiple mutually exclusive variants simultaneously.
// This should be rejected by either Terraform validation or the provider.
func TestAccPanosVirtualRouter_BgpExport_AsPath_MultipleVariants(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("acc-vrouter-bgp-%s", nameSuffix)
	routerName := fmt.Sprintf("test-vr-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosVirtualRouter_BgpExport_AsPath_MultipleVariants_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"router_name":   config.StringVariable(routerName),
				},
				ExpectError: regexp.MustCompile(`Invalid Attribute Combination`),
			},
		},
	})
}

// TestAccPanosVirtualRouter_BgpExport_AsPath_Prepend verifies that the
// as_path prepend and remove_and_prepend variants can be configured with
// non-default values and that only the configured variant appears in state
// while others remain null.
func TestAccPanosVirtualRouter_BgpExport_AsPath_Prepend(t *testing.T) {
	t.Parallel()

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	templateName := fmt.Sprintf("acc-vrouter-bgp-%s", nameSuffix)
	routerName := fmt.Sprintf("test-vr-%s", nameSuffix)

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosVirtualRouter_BgpExport_AsPath_Prepend_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"template_name": config.StringVariable(templateName),
					"router_name":   config.StringVariable(routerName),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Rule 0: prepend variant with value 3
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(0).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.Null(),
							"prepend":            knownvalue.Int64Exact(3),
							"remove_and_prepend": knownvalue.Null(),
						}),
					),
					// Rule 1: remove_and_prepend variant with value 2
					statecheck.ExpectKnownValue(
						"panos_virtual_router.test",
						tfjsonpath.New("protocol").
							AtMapKey("bgp").
							AtMapKey("policy").
							AtMapKey("export").
							AtMapKey("rules").
							AtSliceIndex(1).
							AtMapKey("action").
							AtMapKey("allow").
							AtMapKey("update").
							AtMapKey("as_path"),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"none":               knownvalue.Null(),
							"remove":             knownvalue.Null(),
							"prepend":            knownvalue.Null(),
							"remove_and_prepend": knownvalue.Int64Exact(2),
						}),
					),
				},
			},
		},
	})
}

const panosVirtualRouter_BgpExport_AsPath_Remove_ExplicitNull_Tmpl = `
variable "template_name" { type = string }
variable "router_name" { type = string }

resource "panos_template" "template" {
    name = var.template_name

    location = {
        panorama = {
            panorama_device = "localhost.localdomain"
        }
    }
}

resource "panos_virtual_router" "test" {
    location = {
        template = {
            name = panos_template.template.name
        }
    }

    name = var.router_name

    protocol = {
        bgp = {
            enable    = true
            router_id = "10.0.0.1"
            local_as  = 65001

            policy = {
                export = {
                    rules = [
            {
                name   = "test-export-rule"
                enable = true

                match = {
                    as_path = {
                        regex = "^65001_"
                    }
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                remove = {}
                                prepend = null
                                remove_and_prepend = null
                            }
                        }
                    }
                }
            }
                    ]
                }
            }
        }
    }
}
`

const panosVirtualRouter_BgpExport_AsPath_Prepend_Only_Tmpl = `
variable "template_name" { type = string }
variable "router_name" { type = string }

resource "panos_template" "template" {
    name = var.template_name

    location = {
        panorama = {
            panorama_device = "localhost.localdomain"
        }
    }
}

resource "panos_virtual_router" "test" {
    location = {
        template = {
            name = panos_template.template.name
        }
    }

    name = var.router_name

    protocol = {
        bgp = {
            enable    = true
            router_id = "10.0.0.1"
            local_as  = 65001

            policy = {
                export = {
                    rules = [
            {
                name   = "test-export-prepend"
                enable = true

                match = {
                    address_prefix = [
                        {
                            name = "10.0.0.0/8"
                            exact = false
                        }
                    ]
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                prepend = 3
                            }
                        }
                    }
                }
            }
                    ]
                }
            }
        }
    }
}
`

const panosVirtualRouter_BgpExport_AsPath_MultipleVariants_Tmpl = `
variable "template_name" { type = string }
variable "router_name" { type = string }

resource "panos_template" "template" {
    name = var.template_name

    location = {
        panorama = {
            panorama_device = "localhost.localdomain"
        }
    }
}

resource "panos_virtual_router" "test" {
    location = {
        template = {
            name = panos_template.template.name
        }
    }

    name = var.router_name

    protocol = {
        bgp = {
            enable    = true
            router_id = "10.0.0.1"
            local_as  = 65001

            policy = {
                export = {
                    rules = [
            {
                name   = "test-export-multiple"
                enable = true

                match = {
                    address_prefix = [
                        {
                            name = "10.0.0.0/8"
                            exact = false
                        }
                    ]
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                none = {}
                                prepend = 3
                                remove_and_prepend = 2
                            }
                        }
                    }
                }
            }
                    ]
                }
            }
        }
    }
}
`

const panosVirtualRouter_BgpExport_AsPath_Prepend_ExplicitNull_Tmpl = `
variable "template_name" { type = string }
variable "router_name" { type = string }

resource "panos_template" "template" {
    name = var.template_name

    location = {
        panorama = {
            panorama_device = "localhost.localdomain"
        }
    }
}

resource "panos_virtual_router" "test" {
    location = {
        template = {
            name = panos_template.template.name
        }
    }

    name = var.router_name

    protocol = {
        bgp = {
            enable    = true
            router_id = "10.0.0.1"
            local_as  = 65001

            policy = {
                export = {
                    rules = [
            {
                name   = "test-export-rule"
                enable = true

                match = {
                    address_prefix = [
                        {
                            name = "10.0.0.0/8"
                            exact = false
                        }
                    ]
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                none = null
                                remove = null
                                prepend = 3
                                remove_and_prepend = null
                            }
                        }
                    }
                }
            }
                    ]
                }
            }
        }
    }
}
`

const panosVirtualRouter_BgpExport_AsPath_Remove_Tmpl = `
variable "template_name" { type = string }
variable "router_name" { type = string }

resource "panos_template" "template" {
    name = var.template_name

    location = {
        panorama = {
            panorama_device = "localhost.localdomain"
        }
    }
}

resource "panos_virtual_router" "test" {
    location = {
        template = {
            name = panos_template.template.name
        }
    }

    name = var.router_name

    protocol = {
        bgp = {
            enable    = true
            router_id = "10.0.0.1"
            local_as  = 65001

            policy = {
                export = {
                    rules = [
            {
                name   = "test-export-rule"
                enable = true

                match = {
                    as_path = {
                        regex = "^65001_"
                    }
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                remove = {}
                            }
                        }
                    }
                }
            }
                    ]
                }
            }
        }
    }
}
`

const panosVirtualRouter_BgpExport_AsPath_Prepend_Tmpl = `
variable "template_name" { type = string }
variable "router_name" { type = string }

resource "panos_template" "template" {
    name = var.template_name

    location = {
        panorama = {
            panorama_device = "localhost.localdomain"
        }
    }
}

resource "panos_virtual_router" "test" {
    location = {
        template = {
            name = panos_template.template.name
        }
    }

    name = var.router_name

    protocol = {
        bgp = {
            enable    = true
            router_id = "10.0.0.1"
            local_as  = 65001

            policy = {
                export = {
                    rules = [
            {
                name   = "test-export-prepend"
                enable = true

                match = {
                    address_prefix = [
                        {
                            name = "10.0.0.0/8"
                            exact = false
                        }
                    ]
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                prepend = 3
                            }
                        }
                    }
                }
            },
            {
                name   = "test-export-remove-prepend"
                enable = true

                match = {
                    as_path = {
                        regex = "^65002_"
                    }
                }

                action = {
                    allow = {
                        update = {
                            as_path = {
                                remove_and_prepend = 2
                            }
                        }
                    }
                }
            }
                    ]
                }
            }
        }
    }
}
`
