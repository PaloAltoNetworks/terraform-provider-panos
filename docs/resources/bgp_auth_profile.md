---
page_title: "panos: panos_bgp_auth_profile"
subcategory: "Firewall Networking"
---

# panos_bgp_auth_profile

This resource allows you to add/update/delete a BGP auth profile.


## Example Usage

```hcl
resource "panos_bgp_auth_profile" "example" {
    virtual_router = panos_bgp.conf.virtual_router
    name = "prof1"
    secret = "secret"
}

resource "panos_bgp" "conf" {
    virtual_router = panos_virtual_router.rtr.name
    router_id = "5.5.5.5"
    as_number = "42"
}

resource "panos_virtual_router" "rtr" {
    name = "my virtual router"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this BGP
  auth profile to.
* `name` - (Required) The name.
* `secret` - (Optional) Shared secret for the TCP MD5 authentication.
