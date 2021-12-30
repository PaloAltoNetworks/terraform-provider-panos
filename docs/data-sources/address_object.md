---
page_title: "panos: panos_address_object"
subcategory: "Objects"
---

# panos_address_object

Gets information on an address object.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_address_object" "example" {
    name = panos_address_object.x.name
}

resource "panos_address_object" "x" {
    name = "localnet"
    value = "192.168.80.0/24"
    description = "The 192.168.80 network"
    tags = [
        "internal",
        "dmz",
    ]
}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys to put the address object into (default:
  `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The address object's name.


## Attribute Reference

The following attributes are supported:

* `type` - The type of address object.
* `value` - The address object's value.
* `description` - The address object's description.
* `tags` - (list) List of administrative tags.
