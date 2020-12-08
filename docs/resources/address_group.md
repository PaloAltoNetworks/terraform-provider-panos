---
page_title: "panos: panos_address_group"
subcategory: "Firewall Objects"
---

# panos_address_group

This resource allows you to add/update/delete address groups.

Address groups are either statically defined or dynamically defined, so only
`static_addresses` or `dynamic_match` should be defined within a given address
group.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
# Static group
resource "panos_address_group" "example1" {
    name = "static ntp grp"
    description = "My NTP servers"
    static_addresses = [
        panos_address_object.ao1.name,
        panos_address_object.ao2.name,
    ]
}

resource "panos_address_object" "ao1" {
    name = "ntp1"
    value = "10.0.0.1"
}

resource "panos_address_object" "ao2" {
    name = "ntp2"
    value = "10.0.0.2"
}
```

```hcl
# Dynamic group
resource "panos_address_group" "example2" {
    name = "dynamic grp"
    description = "My internal NTP servers"
    dynamic_match = "'internal' and 'ntp'"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The address group's name.
* `vsys` - (Optional) The vsys to put the address group into (default:
  `vsys1`).
* `static_addresses` - (Optional) The address objects to include in this
  statically defined address group.
* `dynamic_match` - (Optional) The IP tags to include in this DAG.
* `description` - (Optional) The address group's description.
* `tags` - (Optional) List of administrative tags.
