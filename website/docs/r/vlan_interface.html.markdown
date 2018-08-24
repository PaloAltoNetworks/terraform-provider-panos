---
layout: "panos"
page_title: "panos: panos_vlan_interface"
sidebar_current: "docs-panos-resource-vlan-interface"
description: |-
  Manages vlan interfaces.
---

# panos_vlan_interface

This resource allows you to add/update/delete vlan interfaces.

## Example Usage

```hcl
resource "panos_vlan_interface" "example" {
    name = "vlan.17"
    vsys = "vsys1"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.  Must start with `vlan.`.
* `vsys` - (Optional) The vsys that will use this interface (default: `vsys1`).
* `comment` - (Optional) The interface comment.
* `netflow_profile` - (Optional) The netflow profile.
* `static_ips` - (Optional) List of static IPv4 addresses to set for this data
  interface.
* `enable_dhcp` - (Optional) Set to `true` to enable DHCP on this interface.
* `create_dhcp_default_route` - (Optional) Set to `true` to create a DHCP
  default route.
* `dhcp_default_route_metric` - (Optional) The metric for the DHCP default
  route.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
* `adjust_tcp_mss` - (Optional) Adjust TCP MSS (default: false).
* `ipv4_mss_adjust` - (Optional, PAN-OS 8.0+) The IPv4 MSS adjust value.
* `ipv6_mss_adjust` - (Optional, PAN-OS 8.0+) The IPv6 MSS adjust value.
