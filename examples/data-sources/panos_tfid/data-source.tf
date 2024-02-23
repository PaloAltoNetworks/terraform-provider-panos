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

# Example of how to create the ID for panos_nested_address_object
# named "foo" in shared.
data "panos_tfid" "example2" {
  name     = "foo"
  location = "shared"
}
