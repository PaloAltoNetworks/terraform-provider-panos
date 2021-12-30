---
page_title: "panos: panos_custom_url_categories"
subcategory: "Objects"
---

# panos_custom_url_categories

Gets the list of custom URL categories.


## PAN-OS

NGFW and Panorama


## Example Usage

```hcl
data "panos_custom_url_categories" "example" {}
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
