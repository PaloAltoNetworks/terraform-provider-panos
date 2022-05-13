---
page_title: "panos: panos_local_user_db_group"
subcategory: "Device"
---

# panos_local_user_db_group

This resource allows you to add/update/delete local user database groups.


## PAN-OS

NGFW and Panorama.


## Import Name

```shell
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_local_user_db_group" "example" {
    name = "myGroup"
    users = [
        panos_local_user_db_user.one.name,
        panos_local_user_db_user.two.name,
    ]
}

resource "panos_local_user_db_user" "one" {
    name = "wu"
    password = "password"
}

resource "panos_local_user_db_user" "two" {
    name = "tang"
    password = "drowssap"
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
* `users` - List of users in this group.
