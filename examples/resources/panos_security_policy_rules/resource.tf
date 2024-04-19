# Create a single drop-all rule in vsys1 at the end of the security policy.
resource "panos_security_policy_rules" "example1" {
  location = {
    vsys = {}
  }
  position = {
    last = true
  }
  rules = [
    {
      name                  = "drop-all"
      action                = "drop"
      applications          = ["any"]
      categories            = ["any"]
      destination_addresses = ["any"]
      destination_devices   = ["any"]
      destination_zones     = ["any"]
      disabled              = true
      services              = ["application-default"]
      source_addresses      = ["any"]
      source_devices        = ["any"]
      source_users          = ["any"]
      source_zones          = ["any"]
    },
  ]
}

# Create a group of two rules in pre-rulebase of device group "foo"
resource "panos_security_policy_rules" "example2" {
  location = {
    device_group = {
      name = "foo"
    }
  }
  position = {}
  rules = [
    {
      name                  = "bluey"
      action                = "allow"
      applications          = ["any"]
      categories            = ["any"]
      destination_addresses = ["any"]
      destination_devices   = ["any"]
      destination_zones     = ["any"]
      disabled              = true
      services              = ["application-default"]
      source_addresses      = ["any"]
      source_devices        = ["any"]
      source_users          = ["any"]
      source_zones          = ["any"]
    },
    {
      name                  = "bingo"
      action                = "drop"
      applications          = ["any"]
      categories            = ["any"]
      destination_addresses = ["any"]
      destination_devices   = ["any"]
      destination_zones     = ["any"]
      disabled              = true
      services              = ["application-default"]
      source_addresses      = ["any"]
      source_devices        = ["any"]
      source_users          = ["any"]
      source_zones          = ["any"]
    },
  ]
}
