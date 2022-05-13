---
page_title: "panos: panos_panorama_address_group"
subcategory: "Objects"
---

# panos_panorama_address_group

This resource allows you to add/update/delete Panorama address groups.

Address groups are either statically defined or dynamically defined, so only
`static_addresses` or `dynamic_match` should be defined within a given address
group.


## PAN-OS

Panorama


## Import Name

```shell
<device_group>:<name>
```


## Example Usage

```hcl
# Static group
resource "panos_panorama_address_group" "example" {
    name = "static ntp grp"
    description = "My NTP servers"
    static_addresses = [
        panos_panorama_address_object.o1.name,
        panos_panorama_address_object.o2.name,
    ]
}

resource "panos_panorama_address_object" "o1" {
    name = "ntp1"
    value = "192.168.1.1"
}

resource "panos_panorama_address_object" "o2" {
    name = "ntp2"
    value = "192.168.1.1"
}
```

```hcl
# Dynamic group
resource "panos_panorama_address_group" "example" {
    name = "dynamic grp"
    description = "My internal NTP servers"
    dynamic_match = "'internal' and 'ntp'"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The address group's name.
* `device_group` - (Optional) The device group to put the address group into
  (default: `shared`).
* `static_addresses` - (Optional) The address objects to include in this
  statically defined address group.
* `dynamic_match` - (Optional) The IP tags to include in this DAG.
* `description` - (Optional) The address group's description.
* `tags` - (Optional) List of administrative tags.
