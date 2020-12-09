---
page_title: "panos: panos_panorama_layer3_subinterface"
subcategory: "Panorama Networking"
---

# panos_panorama_layer3_subinterface

This resource allows you to add/update/delete Panorama layer3 subinterfaces.


## Import Name

```
<template>::<interface_type>:<parent_interface>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_layer3_subinterface" "example" {
    template = panos_panorama_template.tmpl.name
    parent_interface = panos_panorama_ethernet_interface.e.name
    vsys = "vsys1"
    name = "ethernet1/5.5"
    tag = 5
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"
}

resource "panos_panorama_ethernet_interface" "e" {
    template = panos_panorama_template.tmpl.name
    name = "ethernet1/5"
    vsys = "vsys1"
    mode = "layer3"
}

resource "panos_panorama_template" "tmpl" {
    name = "myTemplate"
}
```

## Argument Reference

The following arguments are supported:

* `template` - (Required) The template name.
* `interface_type` - (Optional) The interface type.  Valid values are `ethernet` (default)
  or `aggregate-ethernet`.
* `parent_interface` - (Required) The name of the parent interface.
* `vsys` - (Required) The vsys that will use this interface.  This should be
  something like `vsys1` or `vsys3`.
* `name` - (Required) The interface's name.
* `tag` - (Optional, int) The interface's tag.
* `static_ips` - (Optional) List of static IPv4 addresses.
* `ipv6_enabled` - (Optional, bool) Set to `true` to enable IPv6.
* `ipv6_interface_id` - (Optional) The IPv6 interface ID.
* `management_profile` - (Optional) The management profile.
* `mtu` - (Optional) The MTU.
* `adjust_tcp_mss` - (Optional) Adjust TCP MSS (default: false).
* `ipv4_mss_adjust` - (Optional) The IPv4 MSS adjust value.
* `ipv6_mss_adjust` - (Optional) The IPv6 MSS adjust value.
* `netflow_profile` - (Optional) The netflow profile.
* `enable_dhcp` - (Optional, bool) Set to `true` to enable DHCP.
* `create_dhcp_default_route` - (Optional) Set to `true` to create a DHCP
  default route.
* `dhcp_default_route_metric` - (Optional) The metric for the DHCP default
  route.
* `comment` - (Optional) The interface comment.
* `decrypt_forward` - (Optional, bool, PAN-OS 8.1+) Set to `true` to enable decrypt forward.
* `dhcp_send_hostname_enable` - (Optional, PAN-OS 9.0+) For DHCP layer3 interfaces:
  enable sending the firewall or a custom hostname to DHCP server
* `dhcp_send_hostname_value` - (Optional, PAN-OS 9.0+) For DHCP layer3 interfaces:
  the interface hostname.  Leaving this unspecified with `dhcp_send_hostname_enable`
  set means to send the system hostname.
