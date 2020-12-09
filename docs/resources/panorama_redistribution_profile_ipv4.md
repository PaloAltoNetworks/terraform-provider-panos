---
page_title: "panos: panos_panorama_redistribution_profile_ipv4"
subcategory: "Panorama Networking"
---

# panos_panorama_redistribution_profile_ipv4

This resource allows you to add/update/delete Panorama IPv4 redistribution
profiles on a virtual router.


## Import Name

```
<template>:<template_stack>:<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_redistribution_profile_ipv4" "example" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_virtual_router.vr.name
    name = "example"
    priority = 1
    action = "redist"
    types = ["static"]
    interfaces = [panos_panorama_virtual_router.vr.interfaces]
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"
}

resource "panos_panorama_virtual_router" "vr" {
    template = panos_panorama_template.t.name
    name = "my virtual router"
    interfaces = ["ethernet1/2"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The redistribution profile's name.
* `template` - (Required) The template name.
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
