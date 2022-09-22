---
page_title: "panos: panos_userid_login"
subcategory: "User-ID"
---

# panos_userid_login

Login a specific user to a given IP address.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_userid_login" "example" {
    ip = "10.2.3.4"
    user = "user1"

    lifecycle {
        create_before_destroy = true
    }
}
```


## Argument Reference

The following arguments are supported:

* `vsys` - The vsys location (default: `vsys1`).
* `user` - (Required) The user.
* `ip` - (Required) The IP address.
