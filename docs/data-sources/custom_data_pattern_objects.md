---
page_title: "panos: panos_custom_data_pattern_objects"
subcategory: "Objects"
---

# panos_custom_data_pattern_objects

Gets the list of custom data pattern objects.


## Example Usage

```hcl
data "panos_custom_data_pattern_objects" "example" {}
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
