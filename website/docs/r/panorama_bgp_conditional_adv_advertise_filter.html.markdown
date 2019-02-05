---
layout: "panos"
page_title: "panos: panos_panorama_bgp_conditional_adv_advertise_filter"
sidebar_current: "docs-panos-panorama-resource-bgp-conditional-adv-advertise-filter"
description: |-
  Manages a Panorama advertise filter for a BGP conditional advertisement.
---

# panos_panorama_bgp_conditional_adv_advertise_filter

This resource allows you to add/update/delete a Panorama advertise filter for a
BGP conditional advertisement.

~> **Note:** A BGP conditional advertisement is valid only if there is at least
one non-exist filter and one advertise filter attached.  This filter must be paired
with the other in order for the configuration to be valid.


## Import Name

```
<template>:<template_stack>:<virtual_router>:<bgp_conditional_adv>:<name>
```


## Example Usage

```hcl
data "panos_system_info" "x" {}

resource "panos_panorama_bgp_conditional_adv_advertise_filter" "example" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    bgp_conditional_adv = "${panos_panorama_bgp_conditional_adv.ca.name}"
    name = "af"
    route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
    address_prefixes = ["192.168.1.0/24"]
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"
}

resource "panos_panorama_bgp_conditional_adv" "ca" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    name = "example"
}

resource "panos_panorama_bgp_conditional_adv_non_exist_filter" "af" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    bgp_conditional_adv = "${panos_panorama_bgp_conditional_adv.ca.name}"
    name = "nef"
    route_table = "${data.panos_system_info.x.version_major >= 8 ? "unicast" : ""}"
    address_prefixes = ["192.168.2.0/24"]
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
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this filter to.
* `bgp_conditional_adv` - (Required) The BGP conditional advertisement to add
  this filter to.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `as_path_regex` - (Optional) AS path to match.
* `community_regex` - (Optional) Community to match.
* `extended_community_regex` - (Optional) Extended community to match.
* `med` - (Optional) Match MED.
* `route_table` - (Optional, PAN-OS 8.0+) Route table to match rule.  Valid
  values are `unicast`, `multicast`, or `both`.  As of PAN-OS 8.1, there doesn't
  seem to be a way to configure this in the GUI, it is always set to `unicast`.
  Thus, if you're running this resource against PAN-OS 8.0+, the appropriate
  thing to do is set this value to `unicast` as well to match the GUI functionality.
* `address_prefixes` - (Optional) List of matching address prefixes.
* `next_hops` - (Optional) List of next hop attributes.
* `from_peers` - (Optional) List of peers that advertised the route entry.
