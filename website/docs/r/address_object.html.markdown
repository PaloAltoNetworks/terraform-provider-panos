---
layout: "panos"
page_title: "panos: panos_address_object"
sidebar_current: "docs-panos-resource-address-object"
description: |-
  Manages address objects.
---

# panos_address_object

This resource allows you to add/update/delete address objects.

## PAN-OS Support

* Firewall
* Panorama

## Example Usage

```hcl
# Firewall example.
resource "panos_address_object" "example" {
    name = "localnet"
    value = "192.168.80.0/24"
    description = "The 192.168.80 network"
    tags = ["internal", "dmz"]
}

# Panorama example.
resource "panos_address_object" "example" {
    name = "localnet"
    device_group = "My Group"
    value = "192.168.80.0/24"
    description = "The 192.168.80 network"
    tags = ["internal", "dmz"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The address object's name.
* `vsys` - (Optional) The firewall vsys to put the address object into
  (default: `vsys1`).
* `device_group` - (Optional) The Panorama device group to put the address
  object into (default: `shared`)
* `type` - (Optional) The type of address object.  This can be `ip-netmask`
  (default), `ip-range`, or `fqdn`.
* `value` - (Required) The address object's value.  This can take various
  forms depending on what type of address object this is, but can be something
  like `192.168.80.150` or `192.168.80.0/24`.
* `description` - (Optional) The address object's description.
* `tags` - (Optional) List of administrative tags.
