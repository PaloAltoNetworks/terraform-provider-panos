---
page_title: "panos: panos_ethernet_interface"
subcategory: "Firewall Networking"
---

# panos_ethernet_interface

This resource allows you to add/update/delete ethernet interfaces.


## Import Name

```
<vsys>:<name>
```


## Example Usage

```hcl
# Configure a bare-bones ethernet interface.
resource "panos_ethernet_interface" "example" {
    name = "ethernet1/3"
    vsys = "vsys1"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}
```

```hcl
# Configure a DHCP ethernet interface for vsys1 to use.
resource "panos_ethernet_interface" "example" {
    name = "ethernet1/4"
    vsys = "vsys1"
    mode = "layer3"
    enable_dhcp = true
    create_dhcp_default_route = true
    dhcp_default_route_metric = 10
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The ethernet interface's name.  This should be something
  like `ethernet1/X`.
* `vsys` - (Required) The vsys that will use this interface.  This should be
  something like `vsys1` or `vsys3`.
* `mode` - (Required) The interface mode.  This can be any of the following
  values: `layer3`, `layer2`, `virtual-wire`, `tap`, `ha`, `decrypt-mirror`,
  or `aggregate-group`.
* `static_ips` - (Optional) List of static IPv4 addresses to set for this data
  interface.
* `enable_dhcp` - (Optional) Set to `true` to enable DHCP on this interface.
* `create_dhcp_default_route` - (Optional) Set to `true` to create a DHCP
  default route.
* `dhcp_default_route_metric` - (Optional) The metric for the DHCP default
  route.
* `ipv6_enabled` - (Optional) Set to `true` to enable IPv6.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
* `adjust_tcp_mss` - (Optional) Adjust TCP MSS (default: false).
* `netflow_profile` - (Optional) The netflow profile.
* `lldp_enabled` - (Optional) Enable LLDP (default: false).
* `lldp_profile` - (Optional) LLDP profile.
* `lldp_ha_passive_pre_negotiation` - (bool) LLDP HA passive pre-negotiation.
* `lacp_ha_passive_pre_negotiation` - (bool) LACP HA passive pre-negotiation.
* `link_speed` - (Optional) Link speed.  This can be any of the following:
  `10`, `100`, `1000`, or `auto`.
* `link_duplex` - (Optional) Link duplex setting.  This can be `full`, `half`,
  or `auto`.
* `link_state` - (Optional) The link state.  This can be `up`, `down`, or
  `auto`.
* `aggregate_group` - (Optional) The aggregate group (applicable for
  physical firewalls only).
* `comment` - (Optional) The interface comment.
* `lacp_port_priority` - (int) LACP port priority.
* `ipv4_mss_adjust` - (Optional, PAN-OS 7.1+) The IPv4 MSS adjust value.
* `ipv6_mss_adjust` - (Optional, PAN-OS 7.1+) The IPv6 MSS adjust value.
* `decrypt_forward` - (Optional, PAN-OS 8.1+) Enable decrypt forwarding.
* `rx_policing_rate` - (Optional, PAN-OS 8.1+) Receive policing rate in Mbps.
* `tx_policing_rate` - (Optional, PAN-OS 8.1+) Transmit policing rate in Mbps.
* `dhcp_send_hostname_enable` - (Optional, PAN-OS 9.0+) For DHCP layer3 interfaces:
  enable sending the firewall or a custom hostname to DHCP server
* `dhcp_send_hostname_value` - (Optional, PAN-OS 9.0+) For DHCP layer3 interfaces:
  the interface hostname.  Leaving this unspecified with `dhcp_send_hostname_enable`
  set means to send the system hostname.
