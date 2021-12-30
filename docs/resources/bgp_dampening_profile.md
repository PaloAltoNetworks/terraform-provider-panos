---
page_title: "panos: panos_bgp_dampening_profile"
subcategory: "Network"
---

# panos_bgp_dampening_profile

This resource allows you to add/update/delete a BGP dampening profile.


## PAN-OS

NGFW


## Import Name

```
<virtual_router>:<name>
```


## Example Usage

```hcl
resource "panos_bgp_dampening_profile" "example" {
    virtual_router = panos_bgp.conf.virtual_router
    name = "myDampeningProfile"
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
  dampening profile to.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `cutoff` - (Optional, float) Cutoff threshold value (default: `1.25`).
* `reuse` - (Optional, float) Reuse threshold value (default: `0.5`).
* `max_hold_time` - (Optional, int) Maximum hold-down time, in
  seconds (default: `900`).
* `decay_half_life_reachable` - (Optional, int) Decay half-life while
  reachable, in seconds (default: `300`).
* `decay_half_life_unreachable` - (Optional, int) Decay half-life while
  unreachable, in seconds (default: `900`).
