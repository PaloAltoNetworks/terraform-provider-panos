---
page_title: "panos: panos_data_filtering_security_profiles"
subcategory: "Objects"
---

# panos_data_filtering_security_profiles

Gets the list of data filtering security profiles.


## Example Usage

```hcl
data "panos_data_filtering_security_profiles" "example" {}
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
