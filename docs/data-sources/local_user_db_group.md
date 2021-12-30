---
page_title: "panos: panos_local_user_db_group"
subcategory: "Device"
---

# panos_local_user_db_group

Retrieve information on the specified local user database group.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_local_user_db_group" "example" {
    name = "myGroup"
}
```

## Argument Reference

Panorama:

* `template` - The template.
* `template_stack` - The template stack.


NGFW:

* `vsys` - (Optional) The vsys (default: `shared`).


The following arguments are supported:

* `name` - (Required) The name.


## Attribute Reference

The following attributes are supported:

* `users` - List of users in this group.
