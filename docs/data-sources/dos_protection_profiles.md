---
page_title: "panos: panos_dos_protection_profiles"
subcategory: "Objects"
---

# panos_dos_protection_profiles

Gets the list of DOS protection profiles.


## Example Usage

```hcl
data "panos_dos_protection_profiles" "example" {}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
