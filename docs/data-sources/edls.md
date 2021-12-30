---
page_title: "panos: panos_edls"
subcategory: "Objects"
---

# panos_edls

Gets the list of EDLs.


## Example Usage

```hcl
data "panos_edls" "example" {}
```

## Argument Reference

Panorama:

* `device_group` - (Optional) The device group (default: `shared`)

NGFW:

* `vsys` - (Optional) The vsys (default: `vsys1`).


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
