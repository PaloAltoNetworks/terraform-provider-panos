---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "panos_addresses Resource - panos"
subcategory: Objects
description: |-
  
---

# panos_addresses (Resource)



## Example Usage

```terraform
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

# Provider function to get the address values

# Example 1: Get the value of a single address object.
output "foo_value" {
  value = provider::panos::address_value(panos_addresses.example.addresses.foo)
}

# Example 2: Transform all the address objects into a map of values.
output "address_values" {
  value = { for k, v in panos_addresses.example.addresses : k => provider::panos::address_value(panos_addresses.example.addresses[k]) }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `addresses` (Attributes Map) (see [below for nested schema](#nestedatt--addresses))
- `location` (Attributes) The location of this object. (see [below for nested schema](#nestedatt--location))

<a id="nestedatt--addresses"></a>
### Nested Schema for `addresses`

Optional:

- `description` (String) The description.
- `disable_override` (String) disable object override in child device groups
- `fqdn` (String) The FQDN value.
- `ip_netmask` (String) The IP netmask value.
- `ip_range` (String) The IP range value.
- `ip_wildcard` (String) The IP wildcard value.
- `tags` (List of String) The administrative tags.


<a id="nestedatt--location"></a>
### Nested Schema for `location`

Optional:

- `device_group` (Attributes) Located in a specific Device Group (see [below for nested schema](#nestedatt--location--device_group))
- `shared` (Attributes) Panorama shared object (see [below for nested schema](#nestedatt--location--shared))
- `vsys` (Attributes) Located in a specific Virtual System (see [below for nested schema](#nestedatt--location--vsys))

<a id="nestedatt--location--device_group"></a>
### Nested Schema for `location.device_group`

Optional:

- `name` (String) Device Group name
- `panorama_device` (String) Panorama device name


<a id="nestedatt--location--shared"></a>
### Nested Schema for `location.shared`


<a id="nestedatt--location--vsys"></a>
### Nested Schema for `location.vsys`

Optional:

- `name` (String) The Virtual System name
- `ngfw_device` (String) The NGFW device name

## Import

Import is supported using the following syntax:

```shell
# Addresses can be imported by providing the following base64 encoded object as the ID
# {
#   location = {
#     device_group = {
#       name = "example-device-group"
#       panorama_device = "localhost.localdomain"
#     }
#   }
#
#   names = [
#     "foo",
#     "bar"
#   ]
# }
terraform import panos_addresses.example $(echo '{"location":{"device_group":{"name":"example-device-group","panorama_device":"localhost.localdomain"}},"names":["foo","bar"]}' | base64)
```