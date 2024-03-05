# Traditional provider example.
provider "panos" {
  hostname = "10.1.1.1"
  username = "admin"
  password = "secret"
}

# Local inspection mode provider example.
provider "panos" {
    config_file = file("/tmp/candidate-config.xml")
    panos_version = "10.2.0"
}

terraform {
  required_providers {
    panos = {
      source  = "paloaltonetworks/terraform-provider-panos"
      version = "2.0.0"
    }
  }
}
