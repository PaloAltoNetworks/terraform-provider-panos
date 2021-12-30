---
page_title: "panos: panos_custom_url_category"
subcategory: "Objects"
---

# panos_custom_url_category

This data source gets info on the given custom URL category.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_custom_url_category" "example" {
    name = "myCustomCategory"
}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group (default: `shared`).

The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

* `description` - The description.
* `sites` - (list) The site list.
* `type` - The custom URL category type.
