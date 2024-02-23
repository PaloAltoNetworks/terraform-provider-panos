provider "panos" {
  hostname = "10.1.1.1"
  username = "admin"
  password = "secret"
}

terraform {
  required_providers {
    panos = {
      source  = "paloaltonetworks/terraform-provider-panos"
      version = "2.0.0"
    }
  }
}
