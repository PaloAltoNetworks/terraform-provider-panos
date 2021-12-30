---
page_title: "panos: panos_dynamic_user_group"
subcategory: "Objects"
---

-> **Note:** Minimum PAN-OS version:  9.0.


# panos_dynamic_user_group

Gets info for a dynamic user group.


## Example Usage

```hcl
data "panos_dynamic_user_group" "example" {
    name = panos_dynamic_user_group.x.name
}

resource "panos_dynamic_user_group" "x" {
    name = "example"
    description = "made by Terraform"
    filter = "'tomato'"
}
```

## Argument Reference

NGFW:

* `vsys` - The vsys (default: `vsys1`).

Panorama:

* `device_group` - The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

* `description` - Description.
* `filter` - The filter.
* `tags` - (list) List of administrative tags.
