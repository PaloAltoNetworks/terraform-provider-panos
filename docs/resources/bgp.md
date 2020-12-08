---
page_title: "panos: panos_bgp"
subcategory: "Firewall Networking"
---

# panos_bgp

This resource allows you to add/update/delete a virtual router's
BGP configuration.

**Important Note:**  When it comes to BGP configuration, PAN-OS requires that
BGP itself first be configured before you can add other BGP sub-config, such
as dampening profiles or peer groups.  Since every BGP resource must reference a
virtual router, the key to accomplishing this is by pointing the `virtual_router`
param for each BGP sub-config to `panos_bgp.foo.virtual_router` instead
of `panos_virtual_router.bar.name`.


## Import Name

```
<virtual_router>
```


## Example Usage

```hcl
resource "panos_bgp" "example" {
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
  configuration to.
* `enable` - (Optional, bool) Enable BGP or not (default: `true`).
* `router_id` - (Optional) Router ID of this BGP instance.
* `as_number` - (Optional) Local AS number.
* `bfd_profile` - (Optional, PAN-OS 7.1+) BFD configuration.
* `reject_default_route` - (Optional, bool) Do not learn default route from
  BGP (default: `true`).
* `install_route` - (Optional, bool) Populate BGP learned route to global
  route table.
* `aggregate_med` - (Optional, bool) Aggregate route only if they have
  same MED attributes (default: `true`).
* `default_local_preference` - (Optional) Default local preference (default:
  `"100"`).
* `as_format` - (Optional) AS format.  Valid values are `"2-byte"` (default)
  or `"4-byte"`.
* `always_compare_med` - (Optional, bool) Always compare MEDs.
* `deterministic_med_comparison` - (Optional, bool) Deterministic MED
  comparison (default: `true`).
* `ecmp_multi_as` - (Optional, bool) Support multiple AS in ECMP.
* `enforce_first_as` - (Optional, bool) Enforce First AS for EBGP (default:
  `true`).
* `enable_graceful_restart` - (Optional, bool) Enable graceful restart
  (default: `true`).
* `stale_route_time` - (Optional, int) Time to remove stale routes after
  peer restart, in seconds (default: `120`).
* `local_restart_time` - (Optional, int) Local restart time to advertise to
  peer, in seconds (default: `120`).
* `max_peer_restart_time` - (Optional, int) Maximum of peer restart time
  accepted, in seconds (default: `120`).
* `reflector_cluster_id` - (Optional) Route reflector cluster ID.
* `confederation_member_as` - (Optional) Confederation requires
  member-AS number.
* `allow_redistribute_default_route` - (Optional, bool) Allow redistribute
  default route to BGP.
