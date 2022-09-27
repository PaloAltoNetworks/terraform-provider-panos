---
page_title: "panos: panos_file_blocking_security_profile"
subcategory: "Objects"
---

# panos_file_blocking_security_profile

Gets info on file_blocking security profiles.


## Example Usage

```hcl
data "panos_file_blocking_security_profile" "example" {
    name = panos_file_blocking_security_profile.x.name
}

resource "panos_file_blocking_security_profile" "x"
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


## Attribute Reference

The following attributes are supported:

* `description` - The description.
* `rule` - (repeatable) Rule list spec, as defined below.

`rule` supports the following arguments:

* `name` - Name.
* `applications` - (list) List of applications.
* `file_types` - (list) List of file types.
* `direction` - The direction.
* `action` - The action to take.
