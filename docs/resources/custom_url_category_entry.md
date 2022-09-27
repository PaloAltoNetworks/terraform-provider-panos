---
page_title: "panos: panos_custom_url_category_entry"
subcategory: "Objects"
---

# panos_custom_url_category_entry

This resource allows you to add/update/delete sites associated with a
custom URL category.

This resource has some overlap with the
[`panos_custom_url_category`](custom_url_category.html) resource.  If
you want to use this resource with the other one, then make sure that your `sites`
param is left undefined in the `panos_custom_url_category` definition.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_custom_url_category_entry" "example" {
    device_group = panos_custom_url_category.x.device_group
    vsys = panos_custom_url_category.x.vsys
    custom_url_category = panos_custom_url_category.x.name
    site = "example.com"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_custom_url_category" "x" {
    name = "myCustomCategory"
    description = "Made by Terraform"
    type = data.panos_system_info.x.version_major >= 9 ? "URL List" : ""

    lifecycle {
        create_before_destroy = true
    }
}

data "panos_system_info" "x" {}
```


## Argument Reference

NGFW:

* `vsys` - The vsys (default: `vsys1`).

Panorama:

* `device_group` - The device group (default: `shared`).

The following arguments are supported:

* `custom_url_category` - (Required) The custom URL category name.
* `site` - (Required) The site to add to the specified custom URL category.
