---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "address_value function - panos"
subcategory: 
description: |-
  Return value of a given address resource
---

# function: address_value

Given an address object resource, return its value.

## Example Usage

```terraform
# Example 1: Get the value of a single address object.
output "foo_value" {
  value = provider::panos::address_value(panos_addresses.example.addresses.foo)
}

# Example 2: Transform all the address objects into a map of values.
output "address_values" {
  value = { for k, v in panos_addresses.example.addresses : k => provider::panos::address_value(panos_addresses.example.addresses[k]) }
}

resource "panos_addresses" "example" {
  location = {
    device_group = {
      name = panos_device_group.example.name
    }
  }

  addresses = {
    "foo" = {
      description = "foo example"
      ip_netmask  = "1.1.1.1"
    }
    "bar" = {
      description = "bar example"
      ip_netmask  = "2.2.2.2"
    }
  }
}

resource "panos_device_group" "example" {
  location = {
    panorama = {}
  }

  name = "example-device-group"
}
```

## Signature

<!-- signature generated by tfplugindocs -->
```text
address_value(address object) string
```

## Arguments

<!-- arguments generated by tfplugindocs -->
1. `address` (Object) address resource to get value from
