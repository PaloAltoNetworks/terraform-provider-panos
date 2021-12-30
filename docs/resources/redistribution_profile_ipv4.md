---
page_title: "panos: panos_redistribution_profile_ipv4"
subcategory: "Network"
---

# panos_redistribution_profile_ipv4

This resource allows you to add/update/delete IPv4 redistribution profiles
on a virtual router.


## PAN-OS

NGFW


## Import Name

```
<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_redistribution_profile_ipv4" "example" {
    virtual_router = panos_virtual_router.vr.name
    name = "example"
    priority = 1
    action = "redist"
    types = ["static"]
    interfaces = [panos_virtual_router.vr.interfaces.0]
}

resource "panos_virtual_router" "vr" {
    name = "my virtual router"
    interfaces = ["ethernet1/2"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The redistribution profile's name.
* `virtual_router` - (Required) The virtual router to add the
  redistribution profile to.
* `priority` - (Required, int) The priority, integer from 1 to 255.
* `action` - (Optional) The action.  Valid values are `redist` (default) or
  `no-redist`.
* `types` - (Optional) The source types.  Valid values are `bgp`, `connect`,
  `ospf`, `rip`, and `static`.
* `interfaces` - (Optional) Specify candidate routes.
* `destinations` - (Optional) Specify candidate routes' next-hop addresses
  (subnet match).
* `next_hops` - (Optional) Specify candidate routes' next-hop addresses
  (subnet match).
* `ospf_path_types` - (Optional) OSPF path types.  Valid values are
  `intra-area`, `inter-area`, `ext-1`, and `ext-2`.
* `ospf_areas` - (Optional) OSPF areas.
* `ospf_tags` - (Optional) OSPF tags.
* `bgp_communities` - (Optional) BGP communities.
* `bgp_extended_communities` - (Optional) BGP extended communities.
