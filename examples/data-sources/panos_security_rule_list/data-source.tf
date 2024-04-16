# Get all security rules from candidate config in vsys1.
data "panos_security_rule_list" "example1" {
  location = {
    vsys = {}
  }
}

# Get all post-rulebase security rules from device group "foobar" whose name
# ends with "DMZ" from running config.
data "panos_security_rule_list" "example2" {
  location = {
    device_group = {
      name     = "foobar"
      rulebase = "post-rulebase"
    }
  }
  filter = {
    config = "running"
    value  = "name ends-with 'DMZ'"
  }
}
