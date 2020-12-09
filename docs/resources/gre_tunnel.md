---
page_title: "panos: panos_gre_tunnel"
subcategory: "Firewall Networking"
---

# panos_gre_tunnel

This resource allows you to add/update/delete GRE tunnels.

*Minimum PAN-OS version*: 9.0

## Import Name

```
<name>
```

## Example Usage

```hcl
resource "panos_gre_tunnel" "example" {
    name = "myGreTunnel"
    interface = panos_ethernet_interface.ei.name
    local_address_value = panos_ethernet_interface.ei.static_ips.0
    peer_address = "192.168.1.1"
    tunnel_interface = panos_tunnel_interface.ti.name
    ttl = 42
}

resource "panos_ethernet_interface" "ei" {
    name = "ethernet1/1"
    vsys = "vsys1"
    mode = "layer3"
    static_ips = ["10.1.1.1/24"]
}

resource "panos_tunnel_interface" "ti" {
    name = "tunnel.7"
    vsys = "vsys1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The GRE tunnel name.
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
