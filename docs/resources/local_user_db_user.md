---
page_title: "panos: panos_local_user_db_user"
subcategory: "Device"
---

# panos_local_user_db_user

This resource allows you to add/update/delete local user database users.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
resource "panos_local_user_db_user" "one" {
    name = "wu"
    password = "password"
    disabled = false

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


NGFW:

* `vsys` - The vsys (default: `shared`).


The following arguments are supported:

* `name` - (Required) The name.
* `password` - The password.
* `disabled` - (bool) If the user is disabled or not.
