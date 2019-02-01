---
layout: "panos"
page_title: "panos: panos_panorama_bgp_dampening_profile"
sidebar_current: "docs-panos-panorama-resource-bgp-dampening-profile"
description: |-
  Manages a Panorama BGP dampening profile.
---

# panos_panorama_bgp_dampening_profile

This resource allows you to add/update/delete a Panorama BGP dampening profile.


## Example Usage

```hcl
resource "panos_panorama_bgp_dampening_profile" "example" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    name = "myDampeningProfile"
}

resource "panos_panorama_bgp" "conf" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_virtual_router.rtr.name}"
    router_id = "5.5.5.5"
    as_number = "42"
}

resource "panos_panorama_virtual_router" "rtr" {
    template = "${panos_panorama_template.t.name}"
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
