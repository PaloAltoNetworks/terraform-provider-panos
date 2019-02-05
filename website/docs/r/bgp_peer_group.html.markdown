---
layout: "panos"
page_title: "panos: panos_bgp_peer_group"
sidebar_current: "docs-panos-resource-bgp-peer-group"
description: |-
  Manages a BGP peer group.
---

# panos_bgp_peer_group

This resource allows you to add/update/delete a BGP peer group.


## Import Name

```
<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_bgp_peer_group" "example" {
    virtual_router = "${panos_bgp.conf.virtual_router}"
    name = "myName"
}

resource "panos_bgp" "conf" {
    virtual_router = "${panos_virtual_router.rtr.name}"
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
