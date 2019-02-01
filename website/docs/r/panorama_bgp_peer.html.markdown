---
layout: "panos"
page_title: "panos: panos_panorama_bgp_peer"
sidebar_current: "docs-panos-panorama-resource-bgp-peer"
description: |-
  Manages a Panorama BGP peer.
---

# panos_panorama_bgp_peer

This resource allows you to add/update/delete a Panorama BGP peer.


## Example Usage

```hcl
data "panos_system_info" "x" {}

// Peer definition that will work starting from PAN-OS 6.1.
resource "panos_panorama_bgp_peer" "example" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    bgp_peer_group = "${panos_panorama_bgp_peer_group.pg.name}"
    name = "peer1"
    peer_as = "${panos_panorama_bgp.conf.as_number}"
    local_address_interface = "${panos_panorama_ethernet_interface.e.name}"
    local_address_ip = "${panos_panorama_ethernet_interface.e.static_ips.0}"
    peer_address_ip = "5.6.7.8"
    max_prefixes = "unlimited"
    bfd_profile = "${
        data.panos_system_info.x.version_major >= 7 ? 
            data.panos_system_info.x.version_minor >= 1 ? "None" : ""
        : ""
    }"
    address_family_type = "${data.panos_system_info.x.version_major >= 8 ? "ipv4" : ""}"
    reflector_client = "${data.panos_system_info.x.version_major >= 8 ? "non-client" : ""}"
    min_route_advertisement_interval = "${
        data.panos_system_info.x.version_major >= 8 ? 
            data.panos_system_info.x.version_minor >= 1 ? 30 : 0
        : 0
    }"
}

resource "panos_panorama_bgp_peer_group" "pg" {
    template = "${panos_panorama_template.t.name}"
    virtual_router = "${panos_panorama_bgp.conf.virtual_router}"
    name = "myName"
    type = "ibgp"
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
    interfaces = ["${panos_panorama_ethernet_interface.e.name}"]
}

resource "panos_panorama_ethernet_interface" "e" {
    template = "${panos_panorama_template.t.name}"
    name = "ethernet1/5"
    mode = "layer3"
    static_ips = ["192.168.1.1/24"]
}

resource "panos_panorama_template" "t" {
    name = "myTemplate"
}
```

## Argument Reference

One and only one of the following must be specified:

* `template` - The template name.
* `template_stack` - The template stack name.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router to add this BGP
  peer to.
* `bgp_peer_group` - (Required) The BGP peer group to put this peer into.
* `name` - (Required) The name.
* `enable` - (Optional, bool) Enable or not (default: `true`).
* `peer_as` - (Optional) Peer AS number.
* `local_address_interface` - (Required) Interface to accept BGP session.
* `local_address_ip` - (Optional) Specify exact IP address if interface has
  multiple addresses.
* `peer_address_ip` - (Required) Peer IP address configuration.
* `reflector_client` - (Optional) This peer is reflector client.  Valid
  values are `non-client`, `client`, or `meshed-client`.
* `peering_type` - (Optional) Peering type that affects NOPEER
  community value handling.  Valid values are `unspecified` (default) or
  `bilateral`.
* `max_prefixes` - (Optional) Maximum of prefixes to receive from the
  peer.  This can be a number such as `"5000"` (default) or `unlimited`.
* `auth_profile` - (Optional) Auth profile.
* `keep_alive_interval` - (Optional, int) Keep alive interval, in
  seconds (default: `30`).
* `multi_hop` - (Optional, int) IP TTL value used for sending BGP packet.
* `open_delay_time` - (Optional, int) Open delay time, in seconds.
* `hold_time` - (Optional, int) Hold time, in seconds.
* `idle_hold_time` - (Optional, int) Idle hold time, in seconds.
* `allow_incoming_connections` - (Optional, bool) Allow incoming connections
  (default: `true`).
* `incoming_connections_remote_port` - (Optional, int) Restrict remote port for
  incoming BGP connections.
* `allow_outgoing_connections` - (Optional, bool) Allow outgoing connections
  (default: `true`).
* `outgoing_connections_local_port` - (Optional, int) Use specific local
  port for outgoing BGP connections.
* `bfd_profile` - (Optional, PAN-OS 7.1+) BFD profile.  This can be a specific
  BFD profile name, `None` (disables BFD), or `Inherit-vr-global-setting`.
* `enable_mp_bgp` - (Optional, bool, PAN-OS 8.0+) Enable MP BGP.
* `address_family_type` - (Optional, PAN-OS 8.0+) Set the AFI for this
  peer.  Valid values are `ipv4` or `ipv6`.
* `subsequent_address_family_unicast` - (Optional, bool, PAN-OS 8.0+) Enable
  unicast subsequent address family for this peer.
* `subsequent_address_family_multicast` - (Optional, bool, PAN-OS 8.0+) Enable
  multicast subsequent address family for this peer.
* `enable_sender_side_loop_detection` - (Optional, bool, PAN-OS 8.0+) Enable
  sender side loop detection.
* `min_route_advertisement_interval` - (Optional, int, PAN-OS 8.1+) Minimum
  route advertisement interval, in seconds.
