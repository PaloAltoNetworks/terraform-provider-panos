# Note that IDs created by this data source may not exactly match IDs created
# by a given resource, but it will be compatible and a valid ID for
# doing state imports.

# Example of how to create the ID for panos_address_object named
# "foo" in vsys1.
#
# All variables should be specified for a given location, default value or not.
data "panos_tfid" "example1" {
  name     = "foo"
  location = "vsys"
  variables = {
    "name" : "vsys1",
    "ngfw_device" : "localhost.localdomain",
  }
}

# Example of how to create the ID for panos_address_object
# named "foo" in shared.
data "panos_tfid" "example2" {
  name     = "foo"
  location = "shared"
}

# Example of how to create the ID for panos_security_policy_rules
# stored in device group "foo".
data "panos_tfid" "example3" {
  location = "device_group"
  variables = {
    "name" : "foo",
    "panorama_device" : "localhost.localdomain",
    "rulebase" : "pre-rulebase",
  }
  rules = [
    {
      name = "bluey"
      uuid = "6d27e31b-0f89-4ac4-a5f5-c7346504b82f"
    },
    {
      name = "bingo"
      uuid = "c685eeea-1a89-47d0-a0ce-3e87dcffda77"
    },
  ]
}
