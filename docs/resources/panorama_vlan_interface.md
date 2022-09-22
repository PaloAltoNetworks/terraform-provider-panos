---
page_title: "panos: panos_panorama_vlan_interface"
subcategory: "Network"
---

# panos_panorama_vlan_interface

This resource allows you to add/update/delete Panorama VLAN interfaces
for templates.


## PAN-OS

Panorama


## Import Name

```shell
<template>:<template_stack>:<vsys>:<name>
```


## Example Usage

```hcl
resource "panos_panorama_vlan_interface" "example" {
    template = panos_panorama_template.t.name
    name = "vlan.17"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
    comment = "Configured for internal traffic"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "t" {
    name = "my template"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The interface's name.  Must start with `vlan.`.
* `template` - (Required) The template name.
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
