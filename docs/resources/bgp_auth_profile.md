---
page_title: "panos: panos_bgp_auth_profile"
subcategory: "Network"
---

# panos_bgp_auth_profile

This resource allows you to add/update/delete a BGP auth profile.


## PAN-OS

NGFW


## Example Usage

```hcl
resource "panos_bgp_auth_profile" "example" {
    virtual_router = panos_bgp.conf.virtual_router
    name = "prof1"
    secret = "secret"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_bgp" "conf" {
    virtual_router = panos_virtual_router.rtr.name
    router_id = "5.5.5.5"
    as_number = "42"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_virtual_router" "rtr" {
    name = "my virtual router"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this BGP
  auth profile to.
* `name` - (Required) The name.
* `secret` - (Optional) Shared secret for the TCP MD5 authentication.
