---
page_title: "panos: panos_user_tag"
subcategory: "User-ID"
---

# panos_user_tag

Manages a specific set of tags for a single user.

This resource only manages the given tags for the given user.  Any
other tags associated with the user are left as-is.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_user_tag" "example1" {
    user = "user1"
    tags = [
        "tag1",
        "tag2",
    ]
}

# It is safe to have multiple resources target the same user.
resource "panos_user_tag" "example2" {
    user = "user1"
    tags = [
        "tag3",
    ]
}
```


## Argument Reference

The following arguments are supported:

* `vsys` - The vsys location (default: `vsys1`).
* `user` - (Required) The user.
* `tags` - (list) List of tags.
