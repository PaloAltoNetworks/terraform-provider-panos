package provider_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	//"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccAdministrator_Password_Hashing(t *testing.T) {
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
				Config: panosAdministrators_Password_Hashing_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
					"password": config.StringVariable("initial"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("password"),
						knownvalue.StringExact("initial"),
					),
				},
			},
			{
				Config: panosAdministrators_Password_Hashing_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
					"password": config.StringVariable("updated"),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("password"),
						knownvalue.StringExact("updated"),
					),
				},
			},
		},
	})
}

const panosAdministrators_Password_Hashing_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }
variable "password" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}

resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = { panorama = {} }

  name = var.prefix

  password = var.password
}
`

func TestAccAdministrator_Basic(t *testing.T) {
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
				Config: panosAdministrators_Basic_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("name"),
						knownvalue.StringExact(prefix),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("client_certificate_only"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("disable_dns"),
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("traffic").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("Example Query"),
							"query": knownvalue.StringExact("addr.src in 10.0.0.0/8"),
						}),
					),
				},
			},
		},
	})
}

const panosAdministrators_Basic_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location

  name = var.prefix
  password = "admin123"

  client_certificate_only = false

  preferences = {
    disable_dns = true
    saved_log_query = {
      traffic = [
        {
          name = "Example Query"
          query = "addr.src in 10.0.0.0/8"
        }
      ]
    }
  }
}
`

func TestAccAdministrator_RoleBased_Custom(t *testing.T) {
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
				Config: panosAdministrators_RoleBased_Custom_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("custom").AtMapKey("profile"),
						knownvalue.StringExact(prefix),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_Custom_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}

resource "panos_admin_role" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  role = {
    vsys = {
      cli = "vsysreader"
    }
  }
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      custom = {
        profile = panos_admin_role.example.name
      }
    }
  }
}
`

func TestAccAdministrator_RoleBased_DeviceAdmin(t *testing.T) {
	t.Parallel()
	t.Skip("requires valid device references")

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
				Config: panosAdministrators_RoleBased_DeviceAdmin_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("device_admins"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("device1"),
							knownvalue.StringExact("device2"),
						}),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_DeviceAdmin_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      device_admins = ["device1", "device2"]
    }
  }
}
`

func TestAccAdministrator_RoleBased_DeviceReader(t *testing.T) {
	t.Parallel()
	t.Skip("requires valid device references")

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
				Config: panosAdministrators_RoleBased_DeviceReader_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("devicereader"),
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("device1"),
							knownvalue.StringExact("device2"),
						}),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_DeviceReader_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      devicereader = ["device1", "device2"]
    }
  }
}
`

func TestAccAdministrator_RoleBased_PanoramaAdmin(t *testing.T) {
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
				Config: panosAdministrators_RoleBased_PanoramaAdmin_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("panorama_admin"),
						knownvalue.StringExact("yes"),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_PanoramaAdmin_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      panorama_admin = "yes"
    }
  }
}
`

func TestAccAdministrator_RoleBased_SuperReader(t *testing.T) {
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
				Config: panosAdministrators_RoleBased_SuperReader_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("superreader"),
						knownvalue.StringExact("yes"),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_SuperReader_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      superreader = "yes"
    }
  }
}
`

func TestAccAdministrator_RoleBased_SuperUser(t *testing.T) {
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
				Config: panosAdministrators_RoleBased_SuperUser_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("superuser"),
						knownvalue.StringExact("yes"),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_SuperUser_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" example {
  location = { panorama = {} }
  name =  var.prefix
}

resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      superuser = "yes"
    }
  }
}
`

func TestAccAdministrator_RoleBased_VsysAdmin(t *testing.T) {
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
				Config: panosAdministrators_RoleBased_VsysAdmin_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("vsys_admins").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("localhost.localdomain"),
							"vsys": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("vsys1"),
							}),
						}),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_VsysAdmin_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name =  var.prefix
}

resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = { template = { name = "test-tmpl-commit" } }
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      vsys_admins = [
        {
          name = "localhost.localdomain"
          vsys = ["vsys1"]
        }
      ]
    }
  }
}
`

func TestAccAdministrator_RoleBased_VsysReader(t *testing.T) {
	t.Parallel()
	t.Skip("requires valid device references")

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
				Config: panosAdministrators_RoleBased_VsysReader_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("permissions").AtMapKey("role_based").AtMapKey("vsys_readers").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name": knownvalue.StringExact("localhost.localdomain"),
							"vsys": knownvalue.ListExact([]knownvalue.Check{
								knownvalue.StringExact("vsys1"),
							}),
						}),
					),
				},
			},
		},
	})
}

const panosAdministrators_RoleBased_VsysReader_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }

  name =  var.prefix
}


resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name     = var.prefix
  password = "admin123"

  permissions = {
    role_based = {
      vsys_readers = [
        {
          name = "localhost.localdomain"
          vsys = ["vsys1"]
        }
      ]
    }
  }
}
`

func TestAccAdministrator_AuthenticationProfile(t *testing.T) {
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
				Config: panosAdministrator_AuthenticationProfile_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("authentication_profile"),
						knownvalue.StringExact("auth-profile"),
					),
				},
			},
		},
	})
}

const panosAdministrator_AuthenticationProfile_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_authentication_profile" "auth_profile" {
  depends_on = [panos_template.example]
  location = var.location
  name = "auth-profile"
  user_domain = "example.com"
}

resource "panos_administrator" "example" {
  depends_on = [panos_authentication_profile.auth_profile]
  location = var.location
  name = var.prefix
  password = "admin123"

  authentication_profile = panos_authentication_profile.auth_profile.name
}
`

func TestAccAdministrator_PasswordProfile(t *testing.T) {
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
				Config: panosAdministrator_PasswordProfile_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("password_profile"),
						knownvalue.StringExact("password-profile"),
					),
				},
			},
		},
	})
}

const panosAdministrator_PasswordProfile_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_password_profile" "pwd_profile" {
  depends_on = [panos_template.example]
  location = var.location
  name = "password-profile"

  password_change = {
    expiration_period = 90
    expiration_warning_period = 7
    post_expiration_admin_login_count = 3
    post_expiration_grace_period = 5
  }
}

resource "panos_administrator" "example" {
  depends_on = [panos_password_profile.pwd_profile]
  location = var.location
  name = var.prefix
  password = "admin123"

  password_profile = panos_password_profile.pwd_profile.name
}
`

func TestAccAdministrator_PublicKey(t *testing.T) {
	t.Parallel()
	t.Skip("requires valid SSH public key")

	nameSuffix := acctest.RandStringFromCharSet(6, acctest.CharSetAlphaNum)
	prefix := fmt.Sprintf("test-acc-%s", nameSuffix)

	location := config.ObjectVariable(map[string]config.Variable{
		"template": config.ObjectVariable(map[string]config.Variable{
			"name": config.StringVariable(prefix),
		}),
	})

	publicKey := "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDKqKT5TZ3Z example@example.com"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: panosAdministrator_PublicKey_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":     config.StringVariable(prefix),
					"location":   location,
					"public_key": config.StringVariable(publicKey),
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("public_key"),
						knownvalue.StringExact(publicKey),
					),
				},
			},
		},
	})
}

const panosAdministrator_PublicKey_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }
variable "public_key" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  password = "admin123"

  public_key = var.public_key
}
`

func TestAccAdministrator_Preferences_SavedLogQuery(t *testing.T) {
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
				Config: panosAdministrator_Preferences_SavedLogQuery_Tmpl,
				ConfigVariables: map[string]config.Variable{
					"prefix":   config.StringVariable(prefix),
					"location": location,
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("alarm").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("Alarm Query"),
							"query": knownvalue.StringExact("severity eq critical"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("auth").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("Auth Query"),
							"query": knownvalue.StringExact("( authproto eq 'RADIUS' )"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("config").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("Config Query"),
							"query": knownvalue.StringExact("( cmd eq 'set' )"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("threat").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("Threat Query"),
							"query": knownvalue.StringExact("severity eq high"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("decryption").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("Decryption Query"),
							"query": knownvalue.StringExact("( policy_name eq 'SSL-Decrypt' )"),
						}),
					),
					statecheck.ExpectKnownValue(
						"panos_administrator.example",
						tfjsonpath.New("preferences").AtMapKey("saved_log_query").AtMapKey("system").AtSliceIndex(0),
						knownvalue.ObjectExact(map[string]knownvalue.Check{
							"name":  knownvalue.StringExact("System Query"),
							"query": knownvalue.StringExact("( eventid eq 'general' )"),
						}),
					),
				},
			},
		},
	})
}

const panosAdministrator_Preferences_SavedLogQuery_Tmpl = `
variable "location" { type = any }
variable "prefix" { type = string }

resource "panos_template" "example" {
  location = { panorama = {} }
  name = var.prefix
}

resource "panos_administrator" "example" {
  depends_on = [panos_template.example]
  location = var.location
  name = var.prefix
  password = "admin123"

  preferences = {
    saved_log_query = {
      alarm = [
        {
          name = "Alarm Query"
          query = "severity eq critical"
        }
      ]
      auth = [
        {
          name = "Auth Query"
          query = "( authproto eq 'RADIUS' )"
        }
      ]
      config = [
        {
          name = "Config Query"
          query = "( cmd eq 'set' )"
        }
      ]
      threat = [
        {
          name = "Threat Query"
          query = "severity eq high"
        }
      ]
      decryption = [
        {
          name = "Decryption Query"
          query = "( policy_name eq 'SSL-Decrypt' )"
        }
      ]
      system = [
        {
          name = "System Query"
          query = "( eventid eq 'general' )"
        }
      ]
    }
  }
}
`
