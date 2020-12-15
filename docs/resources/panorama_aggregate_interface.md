---
page_title: "panos: panos_panorama_aggregate_interface"
subcategory: "Panorama Networking"
---

# panos_panorama_aggregate_interface

This resource allows you to add/update/delete Panorama aggregate ethernet interfaces.


## Import Name

```
<template>::<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_aggregate_interface" "example" {
    template = panos_panorama_template.t1.name
    vsys = "vsys1"
    name = "ae5"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}

resource "panos_panorama_template" "t1" {
    name = "myTemplate"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.
* `templaet` - (Required) The template where the interface should be.
* `vsys` - (Required) The vsys that will use this interface.  This should be
  something like `vsys1` or `vsys3`.
* `mode` - (Required) The interface mode.  Valid values are `layer3` (default),
  `layer2`, `virtual-wire`, `ha`, or `decrypt-mirror`.
* `netflow_profile` - (Optional) The netflow profile.
* `mtu` - (Optional) The MTU.
* `adjust_tcp_mss` - (Optional) Adjust TCP MSS (default: false).
* `ipv4_mss_adjust` - (Optional) The IPv4 MSS adjust value.
* `ipv6_mss_adjust` - (Optional) The IPv6 MSS adjust value.
* `enable_untagged_subinterface` - (Optional, bool) Set to `true` to enable
  untagged subinterfaces.
* `static_ips` - (Optional) List of static IPv4 addresses.
* `ipv6_enabled` - (Optional, bool) Set to `true` to enable IPv6.
* `ipv6_interface_id` - (Optional) The IPv6 interface ID.
* `management_profile` - (Optional) The management profile.
* `enable_dhcp` - (Optional, bool) Set to `true` to enable DHCP.
* `create_dhcp_default_route` - (Optional) Set to `true` to create a DHCP
  default route.
* `dhcp_default_route_metric` - (Optional) The metric for the DHCP default
  route.
* `lacp_enable` - (bool) Enable LACP.
* `lacp_fast_failover` - (bool) Enable LACP fast failover.
* `lacp_mode` - LACP mode.  Valid values are `active` or `passive`.
* `lacp_transmission_rate` - LACP transmission rate.  Valid values are `fast` or `slow`.
* `lacp_system_priority` - (int) LACP system priority.
* `lacp_max_ports` - (int) LACP max ports.
* `lacp_ha_passive_pre_negotiation` - (bool) LACP HA passive pre-negotiation.
* `lacp_ha_enable_same_system_mac` - (bool) LACP HA enable same system MAC.
* `lacp_ha_same_system_mac_address` - LACP HA same system MAC address.
* `lldp_enable` - (bool) Enable LLDP.
* `lldp_profile` - LLDP profile name.
* `lldp_ha_passive_pre_negotiation` - (bool) LLDP HA passive pre-negotiation.
* `comment` - (Optional) The interface comment.
* `decrypt_forward` - (Optional, bool, PAN-OS 8.1+) Set to `true` to enable decrypt forward.
* `dhcp_send_hostname_enable` - (Optional, PAN-OS 9.0+) For DHCP layer3 interfaces:
  enable sending the firewall or a custom hostname to DHCP server
* `dhcp_send_hostname_value` - (Optional, PAN-OS 9.0+) For DHCP layer3 interfaces:
  the interface hostname.  Leaving this unspecified with `dhcp_send_hostname_enable`
  set means to send the system hostname.
