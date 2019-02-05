---
layout: "panos"
page_title: "panos: panos_panorama_bgp_aggregate_advertise_filter"
sidebar_current: "docs-panos-panorama-resource-bgp-aggregate-advertise-filter"
description: |-
  Manages a Panorama route advertise filter for a BGP address aggregation rule.
---

# panos_panorama_bgp_aggregate_advertise_filter

This resource allows you to add/update/delete a Panorama route advertise filter for a
BGP address aggregation rule.


## Import Name

```
<template>:<template_stack>:<virtual_router>:<bgp_aggregate>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_bgp_aggregate_advertise_filter" "example" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp_aggregate.ag.virtual_router}"
    bgp_aggregate = "${panos_panorama_bgp_aggregate.ag.name}"
    name = "my advertise filter"
    as_path_regex = "*42*"
    med = "443"
    address_prefix {
        prefix = "10.1.1.0/24"
        exact = true
    }
    address_prefix {
        prefix = "10.1.2.0/24"
    }
}

resource "panos_panorama_template" "t" {
    name = "my template"
}

resource "panos_panorama_bgp_aggregate" "ag" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    name = "addyAgg1"
    prefix = "192.168.1.0/24"
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
* `bgp_aggregate` - (Required) The BGP address aggregation rule.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `as_path_regex` - (Optional) AS path to match.
* `community_regex` - (Optional) Community to match.
* `extended_community_regex` - (Optional) Extended community to match.
* `med` - (Optional) Match MED.
* `route_table` - (Optional, PAN-OS 8.0+) Route table to match rule.  Valid
  values are `unicast`, `multicast`, or `both`.
* `address_prefix` - (Optional, repeatable) Matching address prefix definition
  (see below).
* `next_hops` - (Optional) List of next hop attributes.
* `from_peers` - (Optional) List of peers that advertised the route entry.

Each `address_prefix` section offers the following params:

* `prefix` - (Required) Address prefix.
* `exact` - (Optional, bool) Match exact prefix length.
