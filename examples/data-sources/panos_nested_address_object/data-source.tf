# Returns candidate config info about "foo" in vsys1.
data "panos_nested_address_object" "example1" {
  location = {
    vsys = {}
  }
  name = "foo"
}

# Returns running config for "bar" in the "baz" device group.
data "panos_nested_address_object" "example2" {
  filter = {
    config = "running"
  }
  location = {
    device_group = {
      name = "baz"
    }
  }
  name = "bar"
}
