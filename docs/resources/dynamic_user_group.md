---
page_title: "panos: panos_dynamic_user_group"
subcategory: "Objects"
---

-> **Note:** Minimum PAN-OS version:  9.0.


# panos_dynamic_user_group

Manages dynamic user groups.


## Import Name

NGFW:

```
<vsys>:<name>
```

Panorama:

```
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_dynamic_user_group" "example" {
    name = "example"
    description = "made by Terraform"
    filter = "'tomato'"
}
```

## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys to put the address object into (default:
  `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) Name.
* `description` - Description.
* `filter` - (Required) The filter.
* `tags` - (list) List of administrative tags.
