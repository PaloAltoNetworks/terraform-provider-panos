---
page_title: "panos: panos_file_blocking_security_profile"
subcategory: "Objects"
---

-> **NOTE:** Minimum PAN-OS version: 8.0.


# panos_file_blocking_security_profile

Manages file_blocking security profiles.


## Import Name

NGFW:

```shell
<vsys>:<name>
```

Panorama:

```shell
<device_group>:<name>
```


## Example Usage

```hcl
resource "panos_file_blocking_security_profile" "example"
    name = "example"
    description = "made by Terraform"
    rule {
        name = "foo"
        applications = ["bbc-streaming"]
        file_types = ["ogg"]
    }

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

NGFW:

* `vsys` - (Optional) The vsys location (default: `vsys1`).

Panorama:

* `device_group` - (Optional) The device group location (default: `shared`)

The following arguments are supported:

* `name` - (Required) The name.
* `description` - The description.
* `rule` - (repeatable) Rule list spec, as defined below.

`rule` supports the following arguments:

* `name` - (Required) Name.
* `applications` - (list) List of applications.
* `file_types` - (list) List of file types.
* `direction` - The direction.  Valid values are `both` (default),
  `upload`, or `download`.
* `action` - The action to take.  Valid values are `alert` (default),
  `block`, or `continue`.
