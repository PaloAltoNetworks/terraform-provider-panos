---
page_title: "panos: panos_dynamic_user_group"
subcategory: "Objects"
---

-> **Note:** Minimum PAN-OS version:  9.0.


# panos_dynamic_user_group

Gets info for a dynamic user group.


## Example Usage

```hcl
data "panos_dynamic_user_groups" "example" {}
```

## Argument Reference

NGFW:

* `vsys` - The vsys (default: `vsys1`).

Panorama:

* `device_group` - The device group location (default: `shared`)


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
