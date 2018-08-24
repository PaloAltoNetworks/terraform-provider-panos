---
layout: "panos"
page_title: "panos: panos_panorama_static_route_ipv4"
sidebar_current: "docs-panos-panorama-resource-static-route-ipv4"
description: |-
  Manages Panorama IPv4 static routes.
---

# panos_panorama_static_route_ipv4

This resource allows you to add/update/delete Panorama IPv4 static routes on a
virtual router for either a template or a template stack.

## Example Usage

```hcl
resource "panos_panorama_static_route_ipv4" "example" {
    name = "localnet"
    virtual_router = "${panos_panorama_virtual_router.vr1.name}"
    template = "template1"
    destination = "10.1.7.0/32"
    next_hop = "10.1.7.4"
}

resource "panos_panorama_virtual_router" "vr1" {
    name = "my virtual router"
    template = "template1"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `name` - (Required) The address object's name.
* `virtual_router` - (Required) The virtual router to add the static
  route to.
* `destination` - (Required) Destination IP address / prefix.
* `interface` - (Optional) Interface to use.
* `type` - (Optional) The next hop type.  Valid values are `ip-address` (the
  default), `discard`, `next-vr`, or an empty string for `None`.
* `next_hop` - (Optional) The value for the `type` setting.
* `admin_distance` - (Optional) The admin distance.
* `metric` - (Optional, int) Metric value / path cost (default: `10`).
* `route_table` - (Optional) Target routing table to install the route.  Valid
  values are `unicast` (the default), `no install`, `multicast`, or `both`.
* `bfd_profile` - (Optional, PAN-OS 7.1+) BFD configuration.
