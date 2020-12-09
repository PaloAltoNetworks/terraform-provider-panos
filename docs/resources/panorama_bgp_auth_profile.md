---
page_title: "panos: panos_panorama_bgp_auth_profile"
subcategory: "Panorama Networking"
---

# panos_panorama_bgp_auth_profile

This resource allows you to add/update/delete a Panorama BGP auth profile.


## Example Usage

```hcl
resource "panos_panorama_bgp_auth_profile" "example" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.conf.virtual_router
    name = "prof1"
    secret = "secret"
}

resource "panos_panorama_bgp" "conf" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_virtual_router.rtr.name
    router_id = "5.5.5.5"
    as_number = "42"
}

resource "panos_panorama_virtual_router" "rtr" {
    template = panos_panorama_template.t.name
    name = "my virtual router"
}

resource "panos_panorama_template" "t" {
    name = "my template"
}
```

## Argument Reference

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this BGP
  auth profile to.
* `name` - (Required) The name.
* `secret` - (Optional) Shared secret for the TCP MD5 authentication.
