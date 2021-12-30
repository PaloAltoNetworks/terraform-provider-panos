---
page_title: "panos: panos_virtual_router"
subcategory: "Network"
---

# panos_virtual_router

Retrieve information on the specified virtual router.


## PAN-OS

NGFW and Panorama.


## Example Usage

```hcl
data "panos_virtual_router" "example" {
    name = "my virtual router"
}
```


## Argument Reference

Panorama (currently only templates can have virtual routers):

* `template` - The template.
* `template_stack` - The template stack.

The following arguments are supported:

* `name` - (Required) The virtual router's name.


## Attribute Reference

The following attributes are supported:

* `interfaces` - List of interfaces.
* `static_dist` - (int) Admin distance - Static.
* `static_ipv6_dist` - (int) Admin distance - Static IPv6.
* `ospf_int_dist` - (int) Admin distance - OSPF Int.
* `ospf_ext_dist` - (int) Admin distance - OSPF Ext.
* `ospfv3_int_dist` - (int) Admin distance - OSPFv3 Int.
* `ospfv3_ext_dist` - (int) Admin distance - OSPFv3 Ext.
* `ibgp_dist` - (int) Admin distance - IBGP.
* `ebgp_dist` - (int) Admin distance - EBGP.
* `rip_dist` - (int) Admin distance - RIP.
* `enable_ecmp` - (bool) Enable ECMP.
* `ecmp_max_path` - (int) Maximum number of ECMP paths supported.
* `ecmp_symmetric_return` - (bool) Allows return packets to egress out
  of the ingress interface of the flow.
* `ecmp_strict_source_path` - (bool) Force VPN traffic to exit interface
  that the source-ip belongs to.
* `ecmp_load_balance_method` - Load balancing algorithm.
* `ecmp_hash_source_only` - (bool) For `ecmp_load_balance_method` = `ip-hash`:
  Only use source address for hash.
* `ecmp_hash_use_port` - (bool) For `ecmp_load_balance_method` = `ip-hash`:
  Use source/destination port for hash.
* `ecmp_hash_seed` - (int) For `ecmp_load_balance_method` = `ip-hash`:
  User-specified hash seed.
* `ecmp_weighted_round_robin_interfaces` - (Map of ints) For `ecmp_load_balance_method` =
  `weighted-round-robin`: Interface weight used in weighted ECMP load balancing.
