---
page_title: "panos: panos_application_objects"
subcategory: "Objects"
---

# panos_application_objects

Gets the list of application objects.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_application_objects" "example" {}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)


## Attribute Reference

The following attributes are supported:

* `total` - (int) The number of items present.
* `listing` - (list) A list of the items present.
