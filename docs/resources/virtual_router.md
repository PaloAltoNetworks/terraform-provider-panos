---
page_title: "panos: panos_virtual_router"
subcategory: "Network"
---

# panos_virtual_router

This resource allows you to add/update/delete virtual routers.

**Note** - The `default` virtual router may be configured with this resource,
however it will not be deleted from the firewall.  It will only be unexported
from the vsys that it is currently imported in, and any interfaces imported
into the virtual router will be removed.

This resource has some overlap with the
[`panos_virtual_router_entry`](virtual_router_entry)
resource.  If you want to use this resource with the other one, then make
sure that your `panos_virtual_router` spec does not define the
`interfaces` field.


## PAN-OS

NGFW and Panorama.


## Aliases

* `panos_panorama_virtual_router`


## Import Name

```
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_virtual_router" "example" {
    name = "my virtual router"
    static_dist = 15
    interfaces = [
        panos_ethernet_interface.e1.name,
        panos_ethernet_interface.e2.name,
    ]
}

resource "panos_ethernet_interface" "e1" {
    vsys = "vsys1"
    name = "ethernet1/1"
    mode = "layer3"
}

resource "panos_ethernet_interface" "e2" {
    vsys = "vsys1"
    name = "ethernet1/2"
    mode = "layer3"
}
```

## Argument Reference

Panorama (currently only templates can have virtual routers):

* `template` - The template.
* `template_stack` - The template stack.


Panorama / NGFW:

* `vsys` - The vsys to import the virtual router into.


The following arguments are supported:

* `name` - (Required) The virtual router's name.
* `interfaces` - List of interfaces that should use this virtual
  router.  If you intend to use
  [`panos_virtual_router_entry`](virtual_router_entry.html) then
  leave this field undefined.
* `static_dist` - (int) Admin distance - Static (default: `10`).
* `static_ipv6_dist` - (int) Admin distance - Static IPv6 (default: `10`).
* `ospf_int_dist` - (int) Admin distance - OSPF Int (default: `30`).
* `ospf_ext_dist` - (int) Admin distance - OSPF Ext (default: `110`).
* `ospfv3_int_dist` - (int) Admin distance - OSPFv3 Int (default: `30`).
* `ospfv3_ext_dist` - (int) Admin distance - OSPFv3 Ext (default: `110`).
* `ibgp_dist` - (int) Admin distance - IBGP (default: `200`).
* `ebgp_dist` - (int) Admin distance - EBGP (default: `20`).
* `rip_dist` - (int) Admin distance - RIP (default: `120`).
* `enable_ecmp` - (bool) Enable ECMP.
* `ecmp_max_path` - (int) Maximum number of ECMP paths supported.
* `ecmp_symmetric_return` - (bool) Allows return packets to egress out
  of the ingress interface of the flow.
* `ecmp_strict_source_path` - (bool) Force VPN traffic to exit interface
  that the source-ip belongs to.
* `ecmp_load_balance_method` - Load balancing algorithm.  Valid
  values are `ip-modulo`, `ip-hash`, `weighted-round-robin`, or
  `balanced-round-robin`.
* `ecmp_hash_source_only` - (bool) For `ecmp_load_balance_method` = `ip-hash`:
  Only use source address for hash.
* `ecmp_hash_use_port` - (bool) For `ecmp_load_balance_method` = `ip-hash`:
  Use source/destination port for hash.
* `ecmp_hash_seed` - (int) For `ecmp_load_balance_method` = `ip-hash`:
  User-specified hash seed.
* `ecmp_weighted_round_robin_interfaces` - (Map of ints) For `ecmp_load_balance_method` =
  `weighted-round-robin`: Interface weight used in weighted ECMP load balancing.
