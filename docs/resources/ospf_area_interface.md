---
page_title: "panos: panos_ospf_area_interface"
subcategory: "Network"
---

# panos_ospf_area_interface

Manages an OSPF area interface.


## Import Name

NGFW:

```
<virtual_router>:<ospf_area>:<name>
```

Panorama:

```
<template>::<virtual_router>:<ospf_area>:<name>
```


## Example Usage

```hcl
# Panorama example.
resource "panos_ospf_area_interface" "example" {
    template = panos_ospf_area.x.template
    virtual_router = panos_ospf_area.x.virtual_router
    ospf_area = panos_ospf_area.x.name
    name = panos_panorama_ethernet_interface.x.name
    enable = true
    passive = true
}

resource "panos_panorama_template" "x" {
    name = "my template"
}

resource "panos_panorama_ethernet_interface" "x" {
    template = panos_panorama_template.x.name
    name = "ethernet1/3"
    mode = "layer3"
    vsys = "vsys1"
}
resource "panos_panorama_virtual_router" "x" {
    template = panos_panorama_template.x.name
    interfaces = [panos_panorama_ethernet_interface.x.name]
    name = "my virtual router"
}

resource "panos_ospf" "x" {
    template = panos_panorama_virtual_router.x.template
    virtual_router = panos_panorama_virtual_router.x.name
    enable = false
}

resource "panos_ospf_area" "x" {
    template = panos_ospf.x.template
    virtual_router = panos_ospf.x.virtual_router
    name = "10.1.2.3"
}
```


## Argument Reference

Panorama:

* `template` - (Optional, but Required for Panorama) The template location.

The following arguments are supported:

* `virtual_router` - (Required) The virtual router name.
* `ospf_area` - (Required) OSPF area name.
* `name` - (Required) Interface name.
* `enable` - (bool) Enable (default: `true`).
* `passive` - (bool) Passive.
* `link_type` - Link type.  Valid values are `broadcast` (default),
  `p2p`, or `p2mp`.
* `metric` - (int) Metric (default: `10`).
* `priority` - (int) Priority (default: `1`).
* `hello_interval` - (int) Hello interval in seconds (default: `10`).
* `dead_counts` - (int) Dead counts (default: `4`).
* `retransmit_interval` - (int) Retransmit interval in seconds (default: `5`).
* `transit_delay` - (int) Transit delay in seconds (default: `1`).
* `grace_restart_delay` - (int) Graceful restart hello delay in
  seconds (default: `10`).
* `auth_profile` - Auth profile.
* `neighbors` - (list) (p2mp) List of neighbors.
* `bfd_profile` - BFD profile.
