---
page_title: "panos: panos_user_tag"
subcategory: "NGFW User-ID"
---


# panos_user_tag

Gets information on user tags.

**NOTE:** This is for Firewall only.


## Example Usage

```hcl
# Get all users and tags.
data "panos_user_tag" "example1" {}

# Get all tags for the given user.
data "panos_user_tag" "example2" {
    user = "user1"
}
```


## Argument Reference

The following arguments are supported:

* `vsys` - The vsys location (default: `vsys1`).
* `user` - The user.


## Attribute Reference

The following attributes are supported:

* `entries` - A list of entries specs, as defined below.

`entries` supports the following attributes:

* `user` - The user.
* `tags` - (list) List of tags.
