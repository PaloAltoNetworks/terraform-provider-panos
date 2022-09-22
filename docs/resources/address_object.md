---
page_title: "panos: panos_address_object"
subcategory: "Objects"
---

# panos_address_object

This resource allows you to add/update/delete address objects.

If you need hundreds or thousands of address objects, you could consider using
[`panos_address_objects`](address_objects.html) instead.


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_address_object`


## Import Name

NGFW:

```shell
<vsys>:<name>
```

Panorama:

```shell
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_address_object" "example" {
    name = "localnet"
    value = "192.168.80.0/24"
    description = "The 192.168.80 network"
    tags = [
        "internal",
        "dmz",
    ]

    lifecycle {
        create_before_destroy = true
    }
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
* `type` - The type of address object.  This can be `ip-netmask`
  (default), `ip-range`, `fqdn`, or `ip-wildcard` (PAN-OS 9.0+).
* `value` - (Required) The address object's value.  This can take various
  forms depending on what type of address object this is, but can be something
  like `192.168.80.150` or `192.168.80.0/24`.
* `description` - The address object's description.
* `tags` - (list) List of administrative tags.
