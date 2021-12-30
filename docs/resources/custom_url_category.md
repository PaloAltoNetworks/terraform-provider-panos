---
page_title: "panos: panos_custom_url_category"
subcategory: "Objects"
---

# panos_custom_url_category

This resource allows you to add/update/delete custom URL categories.

This resource has some overlap with the
[`panos_custom_url_category_entry`](custom_url_category_entry.html) resource.  If
you want to use this resource with the other one, then make sure that your `sites`
param is left undefined.


## PAN-OS

NGFW and Panorama.


## Import Name

```
<device_group>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_custom_url_category" "example" {
    name = "myCustomCategory"
    description = "Made by Terraform"
    sites = [
        "foo.org",
        "bar.com",
        "example.com",
    ]
    type = data.panos_system_info.x.version_major >= 9 ? "URL List" : ""
}

data "panos_system_info" "x" {}
```


## Argument Reference

NGFW:

* `vsys` - The vsys (default: `vsys1`).

Panorama:

* `device_group` - The device group (default: `shared`).

The following arguments are supported:

* `name` - (Required) The name.
* `description` - The description.
* `sites` - (list) The site list.  Leave this undefined if you intend to manage
  the site listing with
  [`panos_custom_url_category_entry`](custom_url_category_entry.html) resources.
* `type` - (PAN-OS 9.0+) The custom URL category type.
