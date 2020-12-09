---
page_title: "panos: panos_panorama_bgp_peer_group"
subcategory: "Panorama Networking"
---

# panos_panorama_bgp_peer_group

This resource allows you to add/update/delete a Panorama BGP peer group.


## Import Name

```
<template>:<template_stack>:<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_bgp_peer_group" "example" {
    template = panos_panorama_template.t.name
    virtual_router = panos_panorama_bgp.conf.virtual_router
    name = "myName"
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
    name = "myTemplate"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this BGP
  peer group to.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `aggregated_confed_as_path` - (Optional, bool) The peers understand aggregated confederation AS path (default: `true`).
* `soft_reset_with_stored_info` - (Optional, bool) Soft reset with stored info.
* `type` - (Optional) Peer group type.  Valid options are `ebgp` (default),
  `ebgp-confed`, `ibgp`, or `ibgp-confed`.
* `export_next_hop` - (Optional) Export next hop.  Valid values are
  `original`, `use-self`, or `resolve`.
* `import_next_hop` - (Optional) Import next hop.  Valid values are
  `original`, `use-peer`, or the empty string.
* `remove_private_as` - (Optional, bool) Remove private AS when exporting
  route.  Only available for `type=ebgp`.
