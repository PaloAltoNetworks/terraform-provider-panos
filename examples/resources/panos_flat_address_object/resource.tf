# Creates an address object in shared (NGFW and Panorama).
resource "panos_flat_address_object" "example1" {
  shared = true

  name        = "sharedObj"
  description = "Made by Terraform"
  ip_netmask  = "10.1.1.0/24"
}

# Creates an address object in the "foo" device group (Panorama only).
resource "panos_flat_address_object" "example2" {
  device_group = "foo"

  name        = "example fqdn"
  description = "Made by Terraform"
  fqdn        = "example.com"
}

# Creates an address object in vsys1 (NGFW only).
resource "panos_flat_address_object" "example3" {
  vsys = "vsys1"

  name        = "bond"
  description = "Shaken not stirred"
  ip_netmask  = "10.0.0.7"
}
