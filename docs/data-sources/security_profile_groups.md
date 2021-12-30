---
page_title: "panos: panos_security_profile_groups"
subcategory: "Objects"
---

# panos_security_profile_groups

Gets the list of security profile groups.


## Example Usage

```hcl
data "panos_security_profile_groups" "example" {}
```


## Argument Reference

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)


NGFW:

* `vsys` - (Optional) The vsys (default: `vsys1`).


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
