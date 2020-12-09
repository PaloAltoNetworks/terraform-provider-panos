---
page_title: "panos: panos_panorama_address_object"
subcategory: "Panorama Objects"
---

# panos_panorama_address_object

This resource allows you to add/update/delete address objects on Panorama.


## Import Name

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_address_object" "example" {
    name = "localnet"
    value = "192.168.80.0/24"
    description = "The 192.168.80 network"
    tags = ["internal", "dmz"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The address object's name.
* `device_group` - (Optional) The device group to put the address object
  into (default: `shared`).
* `type` - (Optional) The type of address object.  This can be `ip-netmask`
  (default), `ip-range`, `fqdn`, or `ip-wildcard` (PAN-OS 9.0+).
* `value` - (Required) The address object's value.  This can take various
  forms depending on what type of address object this is, but can be something
  like `192.168.80.150` or `192.168.80.0/24`.
* `description` - (Optional) The address object's description.
* `tags` - (Optional) List of administrative tags.
