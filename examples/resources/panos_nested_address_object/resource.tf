# Creates an address object in shared (NGFW and Panorama).
resource "panos_nested_address_object" "example1" {
  location = {
    shared = true
  }

  name        = "sharedObj"
  description = "Made by Terraform"
  ip_netmask  = "10.1.1.0/24"
}

# Creates an address object in the "foo" device group (Panorama only).
resource "panos_nested_address_object" "example2" {
  location = {
    device_group = {
      name = "foo"
    }
  }

  name        = "example fqdn"
  description = "Made by Terraform"
  fqdn        = "example.com"
}

# Creates an address object in vsys1 (NGFW only).
resource "panos_nested_address_object" "example3" {
  location = {
    vsys = {}
  }

  name        = "bond"
  description = "Shaken not stirred"
  ip_netmask  = "10.0.0.7"
}
