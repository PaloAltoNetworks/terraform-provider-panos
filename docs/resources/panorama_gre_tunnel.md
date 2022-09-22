---
page_title: "panos: panos_panorama_gre_tunnel"
subcategory: "Network"
---

# panos_panorama_gre_tunnel

This resource allows you to add/update/delete Panorama GRE tunnels.


## Minimum PAN-OS Version

9.0


## PAN-OS

Panorama


## Import Name

```shell
<template>::<name>
```

## Example Usage

```hcl
resource "panos_panorama_gre_tunnel" "example" {
    template = panos_panorama_template.tmpl.name
    name = "myGreTunnel"
    interface = panos_panorama_ethernet_interface.ei.name
    local_address_value = panos_panorama_ethernet_interface.ei.static_ips.0
    peer_address = "192.168.1.1"
    tunnel_interface = panos_panorama_tunnel_interface.ti.name
    ttl = 42

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_template" "tmpl" {
    name = "My Template"

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_ethernet_interface" "ei" {
    template = panos_panorama_template.tmpl.name
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]

    lifecycle {
        create_before_destroy = true
    }
}

resource "panos_panorama_tunnel_interface" "ti" {
    template = panos_panorama_template.tmpl.name
    name = "tunnel.7"
    vsys = "vsys1"

    lifecycle {
        create_before_destroy = true
    }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The GRE tunnel name.
* `template` - (Required) The template name.
* `interface` - (Required) Interface to terminate tunnel.
* `local_address_type` - (Optional) Type of local address.  Valid values are
  `ip` (default) or `floating-ip`.
* `local_address_value` - (Required) IP address value.
* `peer_address` - (Required) Peer IP address.
* `tunnel_interface` - (Required) Tunnel interface to apply the GRE tunnel to.
* `ttl` - (Optional, int) Time to live.
* `copy_tos` - (Optional, bool) Copy IP TOS bits from inner packet to GRE packet.
* `enable_keep_alive` - (Optional, bool) Enable tunnel monitoring.
* `keep_alive_interval` - (Optional, int) Keep alive interval.
* `keep_alive_retry` - (Optional, int) Keep alive retry.
* `keep_alive_hold_timer` - (Optional, int) Keep alive hold timer.
* `disabled` - (Optional, bool) Disable the GRE tunnel.
