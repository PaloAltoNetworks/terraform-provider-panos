---
page_title: "panos: panos_address_objects"
subcategory: "Objects"
---

# panos_address_objects

Gets the list of address objects.


## Example Usage

```hcl
data "panos_address_objects" "example" {}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys to put the address object into (default:
  `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
